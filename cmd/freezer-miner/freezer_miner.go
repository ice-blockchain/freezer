// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/miner"
	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const pkgName = "cmd/freezer-miner"

	var cfg struct{ Version string }
	appCfg.MustLoadFromKey(pkgName, &cfg)

	log.Info(fmt.Sprintf("starting version `%v`...", cfg.Version))

	server.New(new(service), pkgName, "").ListenAndServe(ctx, cancel)
}

type (
	// | service implements server.State and is responsible for managing the state and lifecycle of the package.
	service struct{ miner miner.Client }
)

func (s *service) RegisterRoutes(_ *server.Router) {}

func (s *service) Init(ctx context.Context, cancel context.CancelFunc) {
	s.miner = miner.MustStartMining(ctx, cancel)
}

func (s *service) Close(_ context.Context) error {
	return errors.Wrap(s.miner.Close(), "could not close service")
}

func (s *service) CheckHealth(ctx context.Context) error {
	log.Debug("checking health...", "package", "miner")

	return errors.Wrap(s.miner.CheckHealth(ctx), "failed to check miner's health")
}
