// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	stdlibtime "time"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func newCoinProcessor(client ethClient, db *storage.DB, conf *config) *coinProcessor {
	proc := &coinProcessor{
		Client:         client,
		Conf:           conf,
		WG:             new(sync.WaitGroup),
		CancelSignal:   make(chan struct{}),
		databaseConfig: &databaseConfig{DB: db},
	}
	proc.gasPriceCache.mu = new(sync.RWMutex)
	proc.gasPriceCache.time = time.New(stdlibtime.Time{})

	return proc
}

func (proc *coinProcessor) Start(ctx context.Context, notify chan<- *batch) {
	proc.WG.Add(int(proc.Conf.Workers))

	log.Info(fmt.Sprintf("starting [%d] worker(s) ...", proc.Conf.Workers))

	for workerNumber := int64(0); workerNumber < proc.Conf.Workers; workerNumber++ {
		log.Info(fmt.Sprintf("starting worker [%v]", workerNumber))
		go func(wn int64) {
			defer proc.WG.Done()
			proc.Worker(ctx, notify, wn)
		}(workerNumber)
	}
}

func (proc *coinProcessor) GetGasPrice(ctx context.Context) (value *big.Int, err error) { //nolint:funlen //.
	const retryAttempts = 3

	proc.gasPriceCache.mu.RLock()
	if proc.gasPriceCache.price != nil && stdlibtime.Since(*proc.gasPriceCache.time.Time) < gasPriceCacheTTL {
		proc.gasPriceCache.mu.RUnlock()

		return proc.gasPriceCache.price, nil
	}
	proc.gasPriceCache.mu.RUnlock()

	proc.gasPriceCache.mu.Lock()
	defer proc.gasPriceCache.mu.Unlock()
	if proc.gasPriceCache.price != nil && stdlibtime.Since(*proc.gasPriceCache.time.Time) < gasPriceCacheTTL {
		return proc.gasPriceCache.price, nil
	}

	for attempt := 1; attempt <= retryAttempts; attempt++ {
		value, err = proc.Client.SuggestGasPrice(ctx)
		if err == nil {
			break
		}

		log.Error(errors.Wrapf(err, "failed to get gas price (attempt %v of %v)", attempt, retryAttempts))
		stdlibtime.Sleep(stdlibtime.Second)
	}

	if value == nil {
		return nil, errors.Wrap(err, "failed to get gas price")
	}

	if value != proc.gasPriceCache.price {
		log.Info(fmt.Sprintf("gas price was updated from %v to %v", proc.gasPriceCache.price, value.String()))
	}

	proc.gasPriceCache.price = value
	proc.gasPriceCache.time = time.Now()

	return value, nil
}

func (proc *coinProcessor) BatchMarkAccepted(ctx context.Context, _ int64, data *batch, txHash string) error {
	const stmt = `
update pending_coin_distributions
set
	eth_status = 'ACCEPTED',
	eth_tx = $1
where
	eth_status = 'PENDING' and
	user_id = ANY($2)
`

	_, err := storage.Exec(ctx, proc.DB, stmt, txHash, data.Users())
	data.SetAccepted(txHash)

	return errors.Wrapf(err, "failed to mark batch %v with TX %v as accepted", data.ID, txHash)
}

func (proc *coinProcessor) BatchMarkRejected(ctx context.Context, _ int64, data *batch) error {
	const stmt = `
update pending_coin_distributions
set
	eth_status = 'REJECTED'
where
	eth_status = 'PENDING' and
	user_id = ANY($1)
`
	_, err := storage.Exec(ctx, proc.DB, stmt, data.Users())
	data.SetStatus(ethApiStatusRejected)

	return errors.Wrapf(err, "failed to mark batch %v with as rejected", data.ID)
}

func (proc *coinProcessor) BatchPrepareFetch(ctx context.Context, workerNumber int64) (*batch, error) { //nolint:funlen //.
	const stmt = `
with records as (
	select
		user_id
	from
		pending_coin_distributions
	where
		eth_status = 'NEW' and
		(internal_id % 10) = $1
	order by
		created_at ASC
	limit $2
	for update skip locked
)
update pending_coin_distributions up
set
	eth_status = 'PENDING'
from
	records
where
	up.user_id = records.user_id
returning up.*
`

	result, err := storage.ExecMany[batchRecord](ctx, proc.DB, stmt, workerNumber, batchSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch pending coin distributions")
	} else if len(result) == 0 {
		return nil, errNotEnoughData
	}

	return &batch{
		ID:      ulid.Make().String(),
		Records: result,
	}, nil
}

func (proc *coinProcessor) Distribute(ctx context.Context, num int64, data *batch) (string, error) {
	recipients, amounts := data.Prepare()

	for recordNum := range data.Records {
		log.Info(fmt.Sprintf("worker [%v]: batch %v: distributing %v iceflakes to address %v for user %q",
			num,
			data.ID,
			data.Records[recordNum].Iceflakes,
			data.Records[recordNum].EthAddress,
			data.Records[recordNum].UserID,
		))
	}

	price, err := proc.GetGasPrice(ctx)
	if err != nil {
		log.Error(err, fmt.Sprintf("worker [%v]: batch %v: failed to get gas price", num, data.ID))

		return "", err
	}

	txHash, err := proc.Client.Airdrop(ctx, price, big.NewInt(proc.Conf.Ethereum.ChainID), recipients, amounts)
	if err != nil {
		log.Error(err, fmt.Sprintf("worker [%v]: batch %v: failed to run contract", num, data.ID))

		return "", errors.Wrapf(err, "failed to run contract on batch %v", data.ID)
	}

	log.Info(fmt.Sprintf("worker [%v]: batch %v: transaction hash: %v", num, data.ID, txHash))

	return txHash, nil
}

func (proc *coinProcessor) Do(ctx context.Context, num int64) (*batch, error) {
	data, err := proc.BatchPrepareFetch(ctx, num)
	if err != nil {
		return nil, err
	}

	txHash, err := proc.Distribute(ctx, num, data)
	if err != nil {
		log.Error(err, fmt.Sprintf("worker [%v]: failed to distribute batch", num))
		if err = proc.BatchMarkRejected(ctx, num, data); err != nil {
			log.Error(errors.Wrapf(err, "worker [%v]: failed to mark batch %v as rejected", num, data.ID))
		}
		proc.MustDisable(ctx)

		return data, err
	}

	if err = proc.BatchMarkAccepted(ctx, num, data, txHash); err != nil {
		log.Error(errors.Wrapf(err, "worker [%v]: failed to mark batch %v as accepted", num, data.ID))

		return data, err
	}

	return data, nil
}

func sendNotify[DataType any, DestType chan<- DataType](dest DestType, data DataType) {
	if dest == nil {
		return
	}

	select {
	case dest <- data:
	default:
	}
}

func (proc *coinProcessor) Worker(ctx context.Context, notify chan<- *batch, num int64) { //nolint:funlen //.
	const tickInternal = stdlibtime.Second * 30

	log.Info(fmt.Sprintf("started worker [%v] with internal %s", num, tickInternal))
	defer log.Info(fmt.Sprintf("worker [%v]: stopped", num))

	ticker := stdlibtime.NewTicker(tickInternal)
	defer ticker.Stop()

	signals := make(chan struct{}, 1)
	signals <- struct{}{}

	go func() {
		for range ticker.C {
			select {
			case signals <- struct{}{}:
			default:
				log.Warn(fmt.Sprintf("worker [%v]: signal channel is full", num))
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info(fmt.Sprintf("worker [%v]: %v", num, ctx.Err()))

			return

		case <-proc.CancelSignal:
			log.Info(fmt.Sprintf("worker [%v]: exit signal", num))

			return

		case <-signals:
			if !proc.IsEnabled(ctx) || proc.isBlocked() {
				log.Info(fmt.Sprintf("worker [%v]: coinDistributer is disabled or blocked", num))

				continue
			}

			data, err := proc.Do(ctx, num)
			if data != nil {
				sendNotify(notify, data)
			}
			if err != nil {
				if !errors.Is(err, errNotEnoughData) {
					log.Error(errors.Wrapf(err, "worker [%v]: failed to process batch %v", num, data.ID))
				}

				continue
			}

			select {
			case signals <- struct{}{}:
			default:
			}
		}
	}
}

// isInTimeWindow checks if current hour is in time window [startHour, endHour].
func isInTimeWindow(currentHour, startHour, endHour int) bool {
	for _, v := range []int{currentHour, startHour, endHour} {
		if v < 0 || v > 23 {
			log.Panic(fmt.Sprintf("invalid hour: %v", v))
		}
	}

	if startHour < endHour {
		return currentHour >= startHour && currentHour <= endHour
	}

	return currentHour >= startHour || currentHour <= endHour
}

func (proc *coinProcessor) isBlocked() bool {
	return !isInTimeWindow(time.Now().Hour(), proc.Conf.StartHours, proc.Conf.EndHours)
}

func (proc *coinProcessor) Close() error {
	close(proc.CancelSignal)

	log.Info("waiting for workers to stop ...")
	proc.WG.Wait()

	return nil
}
