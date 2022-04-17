// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"

	"github.com/ICE-Blockchain/freezer/economy"

	"github.com/ICE-Blockchain/wintr/server"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	appCfg "github.com/ICE-Blockchain/wintr/config"
	"github.com/ICE-Blockchain/wintr/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	log.Info("Starting app...", "version", cfg.Version)
	srv := server.New(new(service), applicationYamlKey, "")
	srv.ListenAndServe(ctx, cancel)
}

func (s *service) RegisterRoutes(_ *gin.Engine) {
}

func (s *service) Init(ctx context.Context, cancel context.CancelFunc) {
	s.economyProcessor = economy.StartProcessor(ctx, cancel)
}

func (s *service) Close(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "could not close economy processor because context ended")
	}

	return errors.Wrap(s.economyProcessor.Close(), "could not close economy processor")
}

func (s *service) CheckHealth(ctx context.Context, req *server.RequestCheckHealth) server.Response {
	log.Debug("checking health...", "package", "economy")

	if err := s.economyProcessor.CheckHealth(ctx); err != nil {
		return server.Unexpected(err)
	}

	return server.OK(req)
}
