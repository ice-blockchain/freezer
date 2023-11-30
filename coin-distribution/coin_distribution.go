// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"sync"
	stdlibtime "time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func init() {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
}

func MustStartCoinDistribution(ctx context.Context, cancel context.CancelFunc) Client {
	cd := &coinDistributer{
		db: storage.MustConnect(context.Background(), ddl, applicationYamlKey),
		wg: new(sync.WaitGroup),
	}
	cd.wg.Add(int(cfg.Workers))
	cd.cancel = cancel

	for workerNumber := int64(0); workerNumber < cfg.Workers; workerNumber++ {
		go func(wn int64) {
			defer cd.wg.Done()
			cd.distributeCoins(ctx, wn)
		}(workerNumber)
	}

	return cd
}

func (cd *coinDistributer) Close() error {
	cd.cancel()
	cd.wg.Wait()

	return multierror.Append(
		errors.Wrap(cd.db.Close(), "failed to close db"),
	).ErrorOrNil()
}

func (cd *coinDistributer) CheckHealth(ctx context.Context) error {
	return errors.Wrap(cd.db.Ping(ctx), "[health-check] failed to ping DB")
}

func (cd *coinDistributer) distributeCoins(ctx context.Context, workerNumber int64) {
	for ctx.Err() == nil {
		if !cd.isEnabled(ctx) {
			log.Info(fmt.Sprintf("coinDistributer[%v] is disabled", workerNumber))
			stdlibtime.Sleep(requestDeadline)

			continue
		}
		if currentHour := time.Now().Hour() + 1; (cfg.StartHours < cfg.EndHours && (currentHour < cfg.StartHours || currentHour > cfg.EndHours)) ||
			(cfg.StartHours > cfg.EndHours && (currentHour < cfg.StartHours && currentHour > cfg.EndHours)) {
			log.Info(fmt.Sprintf("coinDistributer[%v] is disabled until %v-%v", workerNumber, cfg.StartHours, cfg.EndHours))
			stdlibtime.Sleep(requestDeadline)

			continue
		}
		reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
		err := storage.DoInTransaction(reqCtx, cd.db, func(conn storage.QueryExecer) error {
			// Logic here

			return nil
		})
		cancel()
		if err == nil {
			log.Info(fmt.Sprintf("TODO: add stuff here"))
		} else {
			log.Error(errors.Wrapf(err, "TODO: add stuff here"))
		}

		if false { // if call to ethereum failed or if transaction commit failed
			// Send Slack message with as much info as possible

			for err = cd.Disable(ctx); err != nil; err = cd.Disable(ctx) {
			}
		}
	}
}

func (cd *coinDistributer) isEnabled(rooCtx context.Context) bool {
	ctx, cancel := context.WithTimeout(rooCtx, requestDeadline)
	defer cancel()
	val, err := storage.Get[struct {
		Enabled bool
	}](ctx, cd.db, `SELECT value::bool as enabled FROM pending_coin_distribution_configurations WHERE key = 'enabled'`)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to check if coinDistributer is enabled"))

		return false
	}

	return val.Enabled
}

func (cd *coinDistributer) Disable(rooCtx context.Context) error {
	ctx, cancel := context.WithTimeout(rooCtx, requestDeadline)
	defer cancel()
	rows, err := storage.Exec(ctx, cd.db, `UPDATE pending_coin_distribution_configurations SET value = 'false' WHERE key = 'enabled'`)
	if err != nil {
		return errors.Wrap(err, "failed to disable coinDistributer")
	}
	if rows == 0 {
		return errors.Wrap(storage.ErrNotFound, "failed to disable coinDistributer")
	}

	return nil
}
