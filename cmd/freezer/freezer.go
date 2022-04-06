// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"

	"github.com/ICE-Blockchain/freezer/cmd/freezer/api"

	"github.com/ICE-Blockchain/freezer/economy"

	"github.com/ICE-Blockchain/wintr/server"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	appCfg "github.com/ICE-Blockchain/wintr/config"
	"github.com/ICE-Blockchain/wintr/log"
)

//nolint:godot // Because those are comments parsed by swagger
// @title                    Economy API
// @version                  latest
// @description              API that handles everything related to user's economy including balance, mining rate and adoption rate.
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
	s.economyRepository = economy.New(ctx, cancel)
}

func (s *service) Close(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "could not close repository because context ended")
	}

	return errors.Wrap(s.economyRepository.Close(), "could not close repository")
}

func (s *service) CheckHealth(ctx context.Context, r *server.RequestCheckHealth) server.Response {
	log.Debug("checking health...", "package", "economy")
	// TODO to be implemented

	return server.OK(r)
}
