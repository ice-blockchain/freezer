// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"runtime"
	stdlibtime "time"

	"github.com/alitto/pond"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
)

func newCoinTracker(client ethClient, db *storage.DB, conf *config) *coinTracker {
	const workersMin = 4
	workersMax := runtime.NumCPU() * 4 //nolint:gomnd //.

	return &coinTracker{
		Client:         client,
		Conf:           conf,
		CancelSignal:   make(chan struct{}),
		Workers:        pond.New(workersMax, 0, pond.MinWorkers(workersMin)),
		databaseConfig: &databaseConfig{DB: db},
	}
}

func (ct *coinTracker) DeleteTransactions(ctx context.Context, hashes []string) error {
	const stmt = `delete from pending_coin_distributions where eth_status = 'ACCEPTED' and eth_tx = any($1)`

	r, err := storage.Exec(ctx, ct.DB, stmt, hashes)
	if err != nil {
		return errors.Wrap(err, "failed to delete transactions")
	}

	log.Info(fmt.Sprintf("deleted %d transaction(s): %v", r, hashes))

	return nil
}

func (ct *coinTracker) RejectTransactions(ctx context.Context, hashes []string) error {
	const stmt = `
update pending_coin_distributions
set
	eth_status = 'REJECTED'
where
	eth_status = 'ACCEPTED' and
	eth_tx = any($1)
`

	r, err := storage.Exec(ctx, ct.DB, stmt, hashes)
	if err != nil {
		return errors.Wrap(err, "failed to update transactions")
	}

	log.Info(fmt.Sprintf("rejected %d transaction(s): %v", r, hashes))

	return nil
}

func (ct *coinTracker) ProcessTransactionsHashes(ctx context.Context, hashes []*string) (err error) {
	statuses, err := ct.Client.TransactionsStatus(ctx, hashes)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction statuses")
	}

	if _, ok := statuses[ethTxStatusSuccessful]; ok {
		if deleteErr := ct.DeleteTransactions(ctx, statuses[ethTxStatusSuccessful]); deleteErr != nil {
			err = multierror.Append(err, deleteErr)
		}
		delete(statuses, ethTxStatusSuccessful)
	}

	if _, ok := statuses[ethTxStatusFailed]; ok {
		if rejectErr := ct.RejectTransactions(ctx, statuses[ethTxStatusFailed]); rejectErr != nil {
			err = multierror.Append(err, rejectErr)
		}
		log.Info("some transactions failed, disabling coinDistributer")
		ct.MustDisable(ctx, fmt.Sprintf("some transactions failed: %v", err.Error()))
		delete(statuses, ethTxStatusFailed)
	}

	for status, hashes := range statuses {
		log.Warn(fmt.Sprintf("found %d unexpected transaction(s) with status %v: %v", len(hashes), status, hashes))
	}

	return err //nolint:wrapcheck //.
}

func (ct *coinTracker) StartChecker(ctx context.Context, notify chan<- []*string) {
	const checkInterval = stdlibtime.Minute

	ticker := stdlibtime.NewTicker(checkInterval)
	defer ticker.Stop()

	signals := make(chan struct{}, 1)
	signals <- struct{}{}

	go func() {
		for range ticker.C {
			select {
			case signals <- struct{}{}:
			default:
				log.Warn("checker: signal channel is full")
			}
		}
	}()

	hadWork := false
	for {
		select {
		case <-ctx.Done():
			log.Info(fmt.Sprintf("stopping checker due to context: %v", ctx.Err()))

			return

		case <-ct.CancelSignal:
			log.Info("stopping checker due to cancel signal")

			return

		case <-signals:
			if hadWork && !ct.HasPendingTransactions(ctx, ethApiStatusAccepted) {
				hadWork = false
				log.Error(sendAllCurrentCoinDistributionsWereCommittedInEthereumSlackMessage(ctx),
					"failed to sendAllCurrentCoinDistributionsWereCommittedInEthereumSlackMessage")
			}

			hasWork, err := ct.Do(ctx, notify)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to check accepted transactions"))
			} else {
				hadWork = hadWork || hasWork
			}
		}
	}
}

func (ct *coinTracker) Do(ctx context.Context, notify chan<- []*string) (submitted bool, err error) {
	const limit = 100
	var offset uint

	for ctx.Err() == nil {
		hashes, err := ct.FetchTransactions(ctx, limit, offset)
		if err != nil {
			return false, errors.Wrap(err, "failed to fetch transactions")
		}
		if len(hashes) == 0 {
			log.Debug("no transactions found to check")

			return submitted, nil
		}

		ct.Workers.Submit(func() {
			err = ct.ProcessTransactionsHashes(ctx, hashes)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to process transactions"))

				return
			}

			sendNotify(notify, hashes)
		})

		submitted = true
		offset += uint(len(hashes))
	}

	return submitted, nil
}

func (ct *coinTracker) FetchTransactions(ctx context.Context, limit, offset uint) ([]*string, error) {
	const stmt = `
select
	distinct on (eth_tx) eth_tx
from
	pending_coin_distributions
where
	eth_status = 'ACCEPTED' and
	eth_tx is not null
order by eth_tx, created_at asc
limit $1
offset $2
`

	rows, err := storage.ExecMany[string](ctx, ct.DB, stmt, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch transactions")
	} else if len(rows) == 0 {
		return []*string{}, nil
	}

	log.Info(fmt.Sprintf("found %d accepted transaction(s) with limit %v and offset %v", len(rows), limit, offset))

	return rows, nil
}

func (ct *coinTracker) Start(ctx context.Context, notify chan<- []*string) {
	go ct.StartChecker(ctx, notify)
}

func (ct *coinTracker) Close() error {
	close(ct.CancelSignal)
	ct.Workers.StopAndWait()

	return nil
}
