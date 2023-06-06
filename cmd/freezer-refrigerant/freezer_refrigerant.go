// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"
	"strconv"
	"sync"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/cmd/freezer-refrigerant/api"
	"github.com/ice-blockchain/freezer/miner"
	"github.com/ice-blockchain/freezer/tokenomics"
	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/server"
)

// @title						Tokenomics API
// @version					latest
// @description				API that handles everything related to write-only operations for user's tokenomics.
// @query.collection.format	multi
// @schemes					https
// @contact.name				ice.io
// @contact.url				https://ice.io
// @BasePath					/v1w
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	api.SwaggerInfo.Host = cfg.Host
	api.SwaggerInfo.Version = cfg.Version
	server.New(new(service), applicationYamlKey, swaggerRoot).ListenAndServe(ctx, cancel)
}

func (s *service) RegisterRoutes(router *server.Router) {
	s.setupTokenomicsRoutes(router)
}

func (s *service) Init(ctx context.Context, cancel context.CancelFunc) {
	s.tokenomicsProcessor = tokenomics.StartProcessor(ctx, cancel)
	s.wg = new(sync.WaitGroup)
	s.wg.Add(1)
	ctx, s.stopX = context.WithCancel(ctx)
	go func() {
		defer s.wg.Done()
		miner.MustStartMining(ctx)
	}()
}

func (s *service) Close(ctx context.Context) error {
	s.stopX()
	s.wg.Wait()
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "could not close processor because context ended")
	}

	return errors.Wrap(s.tokenomicsProcessor.Close(), "could not close service")
}

func (s *service) CheckHealth(ctx context.Context) error {
	log.Debug("checking health...", "package", "tokenomics")

	return errors.Wrap(s.tokenomicsProcessor.CheckHealth(ctx), "failed to check processor's health")
}

func contextWithHashCode[REQ, RESP any](ctx context.Context, req *server.Request[REQ, RESP]) context.Context {
	switch hashCode := req.AuthenticatedUser.Claims["hashCode"].(type) {
	case int:
		return tokenomics.ContextWithHashCode(ctx, uint64(hashCode))
	case int64:
		return tokenomics.ContextWithHashCode(ctx, uint64(hashCode))
	case uint64:
		return tokenomics.ContextWithHashCode(ctx, hashCode)
	case float64:
		return tokenomics.ContextWithHashCode(ctx, uint64(hashCode))
	case string:
		hc, err := strconv.ParseUint(hashCode, 10, 64)
		log.Error(err)

		return tokenomics.ContextWithHashCode(ctx, hc)
	default:
		return ctx
	}
}
