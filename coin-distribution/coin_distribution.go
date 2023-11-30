// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
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
		println(workerNumber)
	}
}
