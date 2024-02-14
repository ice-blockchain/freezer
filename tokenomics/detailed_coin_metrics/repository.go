// SPDX-License-Identifier: ice License 1.0

package detailed_coin_metrics //nolint:revive,nosnakecase,stylecheck //.

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	appCfg "github.com/ice-blockchain/wintr/config"
)

func loadConfig() *config {
	var cfg config

	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	if cfg.APIKey == "" {
		panic("api key is not set")
	}

	return &cfg
}

func New() Repository {
	return newRepository(loadConfig())
}

func newRepository(conf *config) *repository {
	return &repository{
		APIClient: newAPIClient(conf.APIKey),
	}
}

func (r *repository) ReadDetails(ctx context.Context) (*Details, error) {
	data, err := r.APIClient.GetLatestQuote(ctx, iceSlug, targetCurrency)
	if err != nil {
		return nil, errors.Wrap(err, "cannot fetch data from API")
	}

	quote, ok := data[iceID]
	if !ok {
		panic(fmt.Sprintf("unexpected API response: %+v", data))
	}

	return &Details{
		CurrentPrice: quote.Quote[targetCurrency].Price,
		Volume24h:    quote.Quote[targetCurrency].Volume24h,
	}, nil
}
