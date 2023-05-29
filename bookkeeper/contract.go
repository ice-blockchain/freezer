// SPDX-License-Identifier: ice License 1.0

package bookkeeper

import (
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/tokenomics"
)

// Private API.

const (
	applicationYamlKey       = "bookkeeper"
	parentApplicationYamlKey = "tokenomics"
	requestDeadline          = 30 * stdlibtime.Second
)

// .
var (
	//nolint:gochecknoglobals // Singleton & global config mounted only during bootstrap.
	cfg struct {
		tokenomics.Config `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		Workers           int64                    `yaml:"workers"`
		BatchSize         int64                    `yaml:"batchSize"`
	}
)

type (
	bookkeeper struct{}
)
