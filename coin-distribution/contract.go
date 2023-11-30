// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	_ "embed"
	"io"
	"sync"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
)

// Public API.

type (
	Client interface {
		io.Closer
		CheckHealth(context.Context) error
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
	config struct {
		StartHours  int   `yaml:"startHours"`
		EndHours    int   `yaml:"endHours"`
		Workers     int64 `yaml:"workers"`
		BatchSize   int64 `yaml:"batchSize"`
		Development bool  `yaml:"development"`
	}
)
