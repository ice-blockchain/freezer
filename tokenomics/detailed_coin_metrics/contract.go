// SPDX-License-Identifier: ice License 1.0

package detailed_coin_metrics //nolint:revive,nosnakecase,stylecheck //.

import (
	"context"
	"errors"

	"github.com/ice-blockchain/wintr/time"
)

type (
	ReadRepository interface {
		ReadDetails(ctx context.Context) (*Details, error)
	}
	Repository interface {
		ReadRepository
	}
	Details struct {
		UpdatedAt    *time.Time `json:"updatedAt"`
		CurrentPrice float64    `json:"currentPrice"`
		Volume24h    float64    `json:"volume24h"`
	}
)

var ( //nolint:gofumpt //.
	ErrAPIFailed = errors.New("API call failed")
)

const (
	applicationYamlKey = "tokenomics/detailed-coin-metrics"
	iceSlug            = "ice-decentralized-future"
	iceID              = 27650
	targetCurrency     = "USD"
)

type (
	config struct {
		APIKey string `yaml:"api-key" mapstructure:"api-key"` //nolint:tagliatelle,tagalign //.
	}
	repository struct {
		APIClient apiClient
	}
	apiClient interface {
		GetLatestQuote(ctx context.Context, slug, currency string) (map[int]apiResponseQuoteData, error)
	}
	apiClientImpl struct {
		Key string
	}
	apiResponseStatus struct {
		ErrorMessage string `json:"error_message"` //nolint:tagliatelle //.
		ErrorCode    int    `json:"error_code"`    //nolint:tagliatelle //.
	}
	apiResponseQuoteCurrency struct {
		Price     float64 `json:"price"`
		Volume24h float64 `json:"volume_24h"` //nolint:tagliatelle //.
	}
	apiResponseQuoteData struct {
		Quote map[string]apiResponseQuoteCurrency `json:"quote"`
		Name  string                              `json:"name"`
		Slug  string                              `json:"slug"`
		ID    int                                 `json:"id"`
	}
)
