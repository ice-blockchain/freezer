// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	_ "embed"
	"io"
	"sync"

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
		Workers     int64 `yaml:"workers"`
		BatchSize   int64 `yaml:"batchSize"`
		Development bool  `yaml:"development"`
	}
)
