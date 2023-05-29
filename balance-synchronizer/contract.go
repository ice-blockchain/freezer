// SPDX-License-Identifier: ice License 1.0

package balancesynchronizer

import (
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
)

// Public API.

type (
	BalanceUpdated struct {
		UserID     string  `json:"userId,omitempty"`
		Standard   float64 `json:"standard,omitempty"`
		PreStaking float64 `json:"preStaking,omitempty"`
	}
)

// Private API.

const (
	applicationYamlKey       = "balance-synchronizer"
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
	user struct {
		model.UserIDField
		model.MiningBlockchainAccountAddressField
		model.DeserializedUsersKey
		model.BalanceTotalStandardField
		model.BalanceTotalPreStakingField
	}

	balanceSynchronizer struct {
		mb messagebroker.Client
	}
)
