// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/cmd/freezer/api"
	"github.com/ice-blockchain/freezer/economy"
	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/server"
)

//nolint:godot // Because those are comments parsed by swagger
// @title                    Economy API
// @version                  latest
// @description              API that handles everything related to write-only operations for user's economy.
// @query.collection.format  multi
// @schemes                  https
// @contact.name             ICE
// @contact.url              https://ice.io
// @BasePath                 /v1
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	api.SwaggerInfo.Host = cfg.Host
	api.SwaggerInfo.Version = cfg.Version
	srv := server.New(new(service), applicationYamlKey, "/economy")
	srv.ListenAndServe(ctx, cancel)
}

func (s *service) RegisterRoutes(engine *gin.Engine) {
	s.setupEconomyRoutes(engine)
}

func (s *service) Init(ctx context.Context, cancel context.CancelFunc) {
	s.economyProcessor = economy.StartProcessor(ctx, cancel)
	s.economyRepository = economy.New(ctx, cancel)
}

func (s *service) Close(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "could not close economy processor because context ended")
	}

	err1 := s.economyRepository.Close()
	err2 := s.economyProcessor.Close()

	if err1 != nil && err2 != nil {
		return multierror.Append(err1, err2)
	}
	var err error
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}

	return errors.Wrap(err, "could not close service")
}

func (s *service) CheckHealth(ctx context.Context, req *server.RequestCheckHealth) server.Response {
	log.Debug("checking health...", "package", "economy")

	if err := s.economyProcessor.CheckHealth(ctx); err != nil {
		return server.Unexpected(err)
	}

	return server.OK(req)
}
