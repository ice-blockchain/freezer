// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	stdlibtime "time"

	"github.com/hashicorp/go-multierror"
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
	proc.WG.Add(1)
	go func() {
		defer proc.WG.Done()
		proc.Controller(ctx, notify)
	}()
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

func (proc *coinProcessor) BatchMarkAccepted(ctx context.Context, data *batch, txHash string) error {
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

func (proc *coinProcessor) BatchMarkRejected(ctx context.Context, data *batch) error {
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

func (proc *coinProcessor) BatchPrepareFetch(ctx context.Context) (*batch, error) { //nolint:funlen //.
	const stmt = `
with records as (
	select
		user_id
	from
		pending_coin_distributions
	where
		eth_status = 'NEW'
	order by
		created_at ASC
	limit $1
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

	result, err := storage.ExecMany[batchRecord](ctx, proc.DB, stmt, batchSize)
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

func (proc *coinProcessor) DeleteTransactions(ctx context.Context, hash string) error {
	const stmt = `delete from pending_coin_distributions where eth_status = 'ACCEPTED' and eth_tx = $1`

	r, err := storage.Exec(ctx, proc.DB, stmt, hash)
	if err != nil {
		return errors.Wrap(err, "failed to delete transactions")
	}

	log.Info(fmt.Sprintf("transaction: %v: deleted: %v", hash, r))

	return nil
}

func (proc *coinProcessor) RejectTransaction(ctx context.Context, hash string) error {
	const stmt = `
update pending_coin_distributions
set
	eth_status = 'REJECTED'
where
	eth_status = 'ACCEPTED' and
	eth_tx = $1
`

	r, err := storage.Exec(ctx, proc.DB, stmt, hash)
	if err != nil {
		return errors.Wrap(err, "failed to update transactions")
	}

	log.Info(fmt.Sprintf("transaction: %v: rejected: %v", hash, r))

	return nil
}

func (proc *coinProcessor) GetGasOptions(ctx context.Context) (price *big.Int, limit uint64, err error) {
	gasOverride, err := proc.GetGasPriceOverride(ctx)
	if err != nil {
		return nil, 0, err
	}

	if gasOverride == 0 {
		price, err = proc.GetGasPrice(ctx)
		if err != nil {
			return nil, 0, err
		}
	} else {
		price = big.NewInt(0).SetUint64(gasOverride)
	}

	limit, err = proc.GetGasLimit(ctx)
	if err != nil {
		return nil, 0, err
	}

	return price, limit, nil
}

func (proc *coinProcessor) Distribute(ctx context.Context, data *batch) (string, error) {
	recipients, amounts := data.Prepare()
	for recordNum := range data.Records {
		log.Info(fmt.Sprintf("batch %v: distributing %v iceflakes to address %v for user %q",
			data.ID,
			data.Records[recordNum].Iceflakes,
			data.Records[recordNum].EthAddress,
			data.Records[recordNum].UserID,
		))
	}

	txHash, err := proc.Client.Airdrop(ctx, big.NewInt(proc.Conf.Ethereum.ChainID), proc, recipients, amounts)
	if err != nil {
		log.Error(errors.Wrapf(err, "batch %v: failed to run contract", data.ID))

		return "", errors.Wrapf(err, "failed to run contract on batch %v", data.ID)
	}

	log.Info(fmt.Sprintf("batch %v: transaction hash: %v", data.ID, txHash))

	return txHash, nil
}

func (proc *coinProcessor) Do(ctx context.Context) (*batch, error) {
	data, err := proc.BatchPrepareFetch(ctx)
	if err != nil {
		return nil, err
	}

	txHash, err := proc.Distribute(ctx, data)
	if err != nil {
		err = errors.Wrapf(err, "failed to distribute batch")
		log.Error(err)
		if err2 := proc.BatchMarkRejected(ctx, data); err2 != nil {
			log.Error(errors.Wrapf(err2, "failed to mark batch %v as rejected", data.ID))
		}

		return data, err
	}

	data.TX = txHash
	if err = proc.BatchMarkAccepted(ctx, data, txHash); err != nil {
		log.Error(errors.Wrapf(err, "failed to mark batch %v as accepted", data.ID))

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

func (proc *coinProcessor) GetAction(ctx context.Context) workerAction {
	switch {
	case !proc.IsEnabled(ctx):
		return workerActionDisabled

	case proc.IsOnDemandMode(ctx):
		return workerActionOnDemand

	case proc.isBlocked():
		return workerActionBlocked
	}

	return workerActionRun
}

func (proc *coinProcessor) maybeSendMessage(ctx context.Context, key string, f func(context.Context) error) {
	var lastMessageSentAt time.Time

	err := databaseGetValue(ctx, proc.DB, key, &lastMessageSentAt)
	if err != nil {
		log.Error(errors.Wrapf(err, "failed to get %v", key))

		return
	}

	now := time.Now()
	year, month, day := now.Date()
	sentYear, sentMonth, sentDay := lastMessageSentAt.Date()
	if (year == sentYear && month == sentMonth && day == sentDay) || lastMessageSentAt.After(*now.Time) {
		log.Debug(fmt.Sprintf("%v: message was already sent: %v (now %v)", key, lastMessageSentAt, now))

		return
	}

	err = f(ctx)
	if err != nil {
		log.Error(errors.Wrapf(err, "failed to send message"))

		return
	}

	err = databaseSetValue(ctx, proc.DB, key, now)
	if err != nil {
		log.Error(errors.Wrapf(err, "failed to set %v", key))
	}
}

func (proc *coinProcessor) Controller(ctx context.Context, notify chan<- *batch) {
	const tickInternal = stdlibtime.Minute

	log.Info("controller started")
	defer log.Info("controller stopped")

	if proc.HasPendingTransactions(ctx, ethApiStatusAccepted) {
		log.Info("controller: waiting for all accepted transactions to finish")
		err := proc.WaitForAllAcceptedTransactions(ctx, notify)
		if err != nil {
			log.Error(errors.Wrapf(err, "failed to WaitForAllAcceptedTransactions"))
			proc.MustDisable(err.Error())
		}
	}

	ticker := stdlibtime.NewTicker(tickInternal)
	defer ticker.Stop()

	signals := make(chan struct{}, 1)
	signals <- struct{}{}

	go func() {
		for range ticker.C {
			select {
			case signals <- struct{}{}:
			default:
			}
		}
	}()

	prevAction := workerActionBlocked
	for {
		select {
		case <-ctx.Done():
			log.Info(fmt.Sprintf("controller: context: %v", ctx.Err()))

			return

		case <-proc.CancelSignal:
			log.Info("controller: exit signal")

			return

		case <-signals:
			action := proc.GetAction(ctx)
			if action == workerActionDisabled || action == workerActionBlocked {
				log.Info(fmt.Sprintf("controller: disabled or blocked (%v)", action))
				if prevAction == workerActionRun {
					proc.maybeSendMessage(ctx, configKeyCoinDistributerMsgOffline, func(ctx context.Context) (err error) {
						if proc.HasPendingTransactions(ctx, ethApiStatusNew) {
							err = errors.Wrap(sendCoinDistributerHasUnfinishedWork(ctx), "failed to sendCoinDistributerHasUnfinishedWork")
						}

						return multierror.Append(err,
							errors.Wrap(sendCoinDistributerIsNowOfflineSlackMessage(ctx),
								"failed to sendCoinDistributerIsNowOfflineSlackMessage"),
						)
					})
				}
				prevAction = action

				continue
			}

			if action == workerActionRun {
				log.Info("controller: unblocked")
				proc.maybeSendMessage(ctx, configKeyCoinDistributerMsgOnline, sendCoinDistributerIsNowOnlineSlackMessage)
			} else if action == workerActionOnDemand {
				log.Info("controller: on demand mode trigger")
				log.Error(errors.Wrapf(proc.DisableOnDemand(ctx), "failed to DisableOnDemand"))
			}
			prevAction = action

			if !proc.HasPendingTransactions(ctx, ethApiStatusNew) {
				log.Info("controller: no pending transactions")

				continue
			}

			log.Info(fmt.Sprintf("controller: running action %v", action))
			log.Error(errors.Wrap(sendCoinDistributerStartedProcessingSlackMessage(ctx),
				"failed to send DistributerStartedProcessingSlackMessage"))
			err := proc.RunDistribution(ctx, action == workerActionOnDemand, notify)
			if err != nil {
				log.Error(errors.Wrapf(err, "controller: action %v failed with error %v", action, err))
				proc.MustDisable(err.Error())

				continue
			}

			if !proc.HasPendingTransactions(ctx, ethApiStatusNew) {
				proc.maybeSendMessage(
					ctx,
					configKeyCoinDistributerMsgFinished,
					sendAllCurrentCoinDistributionsWereCommittedInEthereumSlackMessage)
			}

			log.Info(fmt.Sprintf("controller: action %v finished", action))
		}
	}
}

func (proc *coinProcessor) GetNextPendingTransaction(ctx context.Context) (string, error) {
	const stmt = `
select
	eth_tx
from
	pending_coin_distributions
where
	eth_status = 'ACCEPTED'
order by
	created_at ASC
limit 1
`

	val, err := storage.Get[string](ctx, proc.DB, stmt)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			err = nil
		}

		return "", errors.Wrap(err, "failed to get next pending transaction")
	}

	return *val, nil
}

func (proc *coinProcessor) WaitForAllAcceptedTransactions(ctx context.Context, notify chan<- *batch) error {
	for ctx.Err() == nil {
		txHash, err := proc.GetNextPendingTransaction(ctx)
		if err != nil {
			return err
		}

		if txHash == "" {
			return nil
		}

		status, err := proc.WaitForTransaction(ctx, txHash)
		if err != nil {
			return err
		}

		switch status {
		case ethTxStatusSuccessful:
			err = proc.DeleteTransactions(ctx, txHash)

		case ethTxStatusFailed:
			proc.MustDisable(fmt.Sprintf("accepted transaction %v failed", txHash))
			err = proc.RejectTransaction(ctx, txHash)
		}

		if err != nil {
			return errors.Wrapf(err, "failed to update transaction status: %v", txHash)
		}

		sendNotify(notify, &batch{TX: txHash, Status: status, Records: []*batchRecord{
			{
				EthTX: &txHash,
			},
		}})
	}

	return ctx.Err()
}

func sleepWithContext(ctx context.Context, d stdlibtime.Duration) {
	tm := stdlibtime.NewTimer(d)
	defer tm.Stop()

	select {
	case <-ctx.Done():
	case <-tm.C:
	}
}

func (proc *coinProcessor) WaitForTransaction(ctx context.Context, hash string) (ethTxStatus, error) {
	start := time.Now()

	for ctx.Err() == nil {
		status, err := proc.Client.TransactionStatus(ctx, hash)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get transaction status: %v", hash)
		}

		if status == ethTxStatusPending {
			sleepDuration := stdlibtime.Second * 3

			if d := stdlibtime.Since(*start.Time); d > stdlibtime.Hour {
				sleepDuration = stdlibtime.Hour
				log.Warn(fmt.Sprintf("transaction %v is in PENDING state since %v (%v)", hash, start, d))
				if d > stdlibtime.Hour*24 {
					sendCoinDistributerTransactionStuck(ctx, hash, start)
				}
			}

			sleepWithContext(ctx, sleepDuration)

			continue
		}

		log.Info(fmt.Sprintf("transaction %v: status: %v, duration: %v", hash, status, stdlibtime.Since(*start.Time)))

		return status, nil
	}

	return "", ctx.Err()
}

func (proc *coinProcessor) RunDistribution(ctx context.Context, ondemand bool, notify chan<- *batch) error {
	for it := 1; ctx.Err() == nil; it++ {
		if !proc.IsEnabled(ctx) {
			log.Info(fmt.Sprintf("distribution: iteration %v: disabled", it))

			return nil
		} else if proc.isBlocked() && !ondemand {
			log.Info(fmt.Sprintf("distribution: iteration %v: blocked", it))

			return nil
		}

		log.Info(fmt.Sprintf("distribution: iteration %v", it))
		b, err := proc.Do(context.WithoutCancel(ctx))
		if err != nil {
			if errors.Is(err, errNotEnoughData) {
				err = nil
			}

			return err
		}

		status, err := proc.WaitForTransaction(ctx, b.TX)
		if err != nil {
			return err
		}

		b.Status = status
		sendNotify(notify, b)

		switch status {
		case ethTxStatusSuccessful:
			err = proc.DeleteTransactions(ctx, b.TX)

		case ethTxStatusFailed:
			proc.MustDisable(fmt.Sprintf("transaction %v failed", b.TX))

			return errors.Wrapf(proc.RejectTransaction(ctx, b.TX), "failed to reject transaction %v", b.TX)

		default:
			log.Panic(fmt.Sprintf("unexpected transaction status: %v (%v)", status, b.TX))
		}

		if err != nil {
			return errors.Wrapf(err, "failed to update transaction status: %v", b.TX)
		}
	}

	return ctx.Err()
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
