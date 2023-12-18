// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	stdlibtime "time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
)

func (d *databaseConfig) MustDisable(ctx context.Context) {
	for err := d.Disable(ctx); err != nil; err = d.Disable(ctx) {
		log.Error(err, "failed to disable coinDistributer")
		stdlibtime.Sleep(stdlibtime.Millisecond)
	}
}

func (d *databaseConfig) IsEnabled(ctx context.Context) bool {
	reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()
	val, err := storage.Get[struct {
		Enabled bool
	}](reqCtx, d.DB, `SELECT value::bool as enabled FROM global WHERE key = 'coin_distributer_enabled'`)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to check if coinDistributer is enabled"))

		return false
	}

	return val.Enabled
}

func (d *databaseConfig) Disable(ctx context.Context) error {
	reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()
	rows, err := storage.Exec(reqCtx, d.DB, `UPDATE global SET value = 'false' WHERE key = 'coin_distributer_enabled'`)
	if err != nil {
		return errors.Wrap(err, "failed to disable coinDistributer")
	}
	if rows == 0 {
		return errors.Wrap(storage.ErrNotFound, "failed to disable coinDistributer")
	}

	return nil
}

func init() {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
}

func MustStartCoinDistribution(ctx context.Context, _ context.CancelFunc) Client {
	cfg.EnsureValid()
	eth := mustNewEthClient(ctx, cfg.Ethereum.RPC, cfg.Ethereum.PrivateKey, cfg.Ethereum.ContractAddress)

	cd := mustCreateCoinDistributionFromConfig(ctx, &cfg, eth)
	cd.MustStart(ctx, nil, nil)

	return cd
}

func mustCreateCoinDistributionFromConfig(ctx context.Context, conf *config, ethClient ethClient) *coinDistributer {
	db := storage.MustConnect(ctx, ddl, applicationYamlKey)
	cd := &coinDistributer{
		Client:    ethClient,
		Processor: newCoinProcessor(ethClient, db, conf),
		Tracker:   newCoinTracker(ethClient, db, conf),
		DB:        db,
	}

	return cd
}

func (cd *coinDistributer) MustStart(ctx context.Context, notifyProcessed chan<- *batch, notifyTracked chan<- []*string) {
	cd.Tracker.Start(ctx, notifyTracked)
	cd.Processor.Start(ctx, notifyProcessed)
}

func (cd *coinDistributer) Close() error {
	return multierror.Append( //nolint:wrapcheck //.
		errors.Wrap(cd.Processor.Close(), "failed to close processor"),
		errors.Wrap(cd.Tracker.Close(), "failed to close tracker"),
		errors.Wrap(cd.Client.Close(), "failed to close eth client"),
		errors.Wrap(cd.DB.Close(), "failed to close db"),
	).ErrorOrNil()
}

func (cd *coinDistributer) CheckHealth(ctx context.Context) error {
	return errors.Wrap(cd.DB.Ping(ctx), "[health-check] failed to ping DB")
}
