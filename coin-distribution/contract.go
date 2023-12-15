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
		GetCoinDistributionsForReview(ctx context.Context, cursor, limit uint64) (updatedCursor uint64, distributions []*PendingReview, err error)
		CheckHealth(ctx context.Context) error
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
		db *storage.DB
	}
	config struct {
		StartHours  int   `yaml:"startHours"`
		EndHours    int   `yaml:"endHours"`
		Workers     int64 `yaml:"workers"`
		BatchSize   int64 `yaml:"batchSize"`
		Development bool  `yaml:"development"`
	}
	coinDistribution struct {
		*PendingReview
		Day        stdlibtime.Time
		InternalID uint64
	}
)
