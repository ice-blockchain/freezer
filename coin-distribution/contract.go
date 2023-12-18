// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	_ "embed"
	"io"
	"sync"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

type (
	Client interface {
		io.Closer
		CheckHealth(ctx context.Context) error
	}

	Repository interface {
		io.Closer
		GetCoinDistributionsForReview(ctx context.Context, arg *GetCoinDistributionsForReviewArg) (*CoinDistributionsForReview, error)
		CheckHealth(ctx context.Context) error
		ReviewCoinDistributions(ctx context.Context, reviewerUserID string, decision string) error
		NotifyCoinDistributionCollectionCycleEnded(ctx context.Context) error
		GetCollectorStatus(ctx context.Context) (latestCollectingDate *time.Time, collectorEnabled bool, err error)
		CollectCoinDistributionsForReview(ctx context.Context, records []*ByEarnerForReview) error
	}

	CoinDistributionsForReview struct {
		Distributions []*PendingReview `json:"distributions"`
		Cursor        uint64           `json:"cursor" example:"5065"`
		TotalRows     uint64           `json:"totalRows" example:"5065"`
		TotalIce      float64          `json:"totalIce" example:"5065.3"`
	}

	GetCoinDistributionsForReviewArg struct {
		CreatedAtOrderBy          string `form:"createdAtOrderBy" example:"asc"`
		IceOrderBy                string `form:"iceOrderBy" example:"asc"`
		UsernameOrderBy           string `form:"usernameOrderBy" example:"asc"`
		ReferredByUsernameOrderBy string `form:"referredByUsernameOrderBy" example:"asc"`
		UsernameKeyword           string `form:"usernameKeyword" example:"jdoe"`
		ReferredByUsernameKeyword string `form:"referredByUsernameKeyword" example:"jdoe"`
		Cursor                    uint64 `form:"cursor" example:"5065"`
		Limit                     uint64 `form:"limit" example:"5000"`
	}

	PendingReview struct {
		CreatedAt          *time.Time `json:"time" swaggertype:"string" example:"2022-01-03T16:20:52.156534Z"`
		Iceflakes          string     `json:"iceflakes" swaggertype:"string" example:"100000000000000"`
		Username           string     `json:"username" swaggertype:"string" example:"myusername"`
		ReferredByUsername string     `json:"referredByUsername" swaggertype:"string" example:"myrefusername"`
		UserID             string     `json:"userId" swaggertype:"string" example:"12746386-03de-44d7-91c7-856fa66b6ed6"`
		EthAddress         string     `json:"ethAddress" swaggertype:"string" example:"0x43...."`
		Ice                float64    `json:"ice" db:"-" example:"1000"`
		IceInternal        int64      `json:"-" db:"ice" swaggerignore:"true"`
	}

	ByEarnerForReview struct {
		CreatedAt          *time.Time
		Username           string
		ReferredByUsername string
		UserID             string
		EarnerUserID       string
		EthAddress         string
		InternalID         int64
		Balance            float64
	}
)

// Private API.

const (
	applicationYamlKey = "coin-distribution"
	requestDeadline    = 25 * stdlibtime.Second
)

// .
var (
	//nolint:gochecknoglobals // Singleton & global config mounted only during bootstrap.
	cfg config
	//go:embed DDL.sql
	ddl string
)

type (
	coinDistributer struct {
		db     *storage.DB
		cancel context.CancelFunc
		wg     *sync.WaitGroup
	}
	repository struct {
		cfg *config
		db  *storage.DB
	}
	config struct {
		AlertSlackWebhook string `yaml:"alert-slack-webhook" mapstructure:"alert-slack-webhook"` //nolint:tagliatelle // .
		Environment       string `yaml:"environment" mapstructure:"environment"`
		ReviewURL         string `yaml:"review-url" mapstructure:"review-url"`
		StartHours        int    `yaml:"startHours"`
		EndHours          int    `yaml:"endHours"`
		Workers           int64  `yaml:"workers"`
		BatchSize         int64  `yaml:"batchSize"`
		Development       bool   `yaml:"development"`
	}
)
