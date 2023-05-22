// SPDX-License-Identifier: ice License 1.0

package extrabonusnotifier

import (
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/tokenomics"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

type (
	ExtraBonusAvailable struct {
		UserID          string `json:"userId,omitempty"`
		ExtraBonusIndex uint16 `json:"extraBonusIndex,omitempty"`
	}
)

// Private API.

const (
	applicationYamlKey = "extra-bonus-notifier"
	requestDeadline    = 30 * stdlibtime.Second
)

// .
var (
	//nolint:gochecknoglobals // Singleton & global config mounted only during bootstrap.
	cfg struct {
		tokenomics.Config `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		Workers           int64                    `yaml:"workers"`
		BatchSize         int64                    `yaml:"batchSize"`
		Chunks            uint16                   `yaml:"chunks"`
	}
)

type (
	user struct {
		*UpdatedUser
		tokenomics.ExtraBonusStartedAtField
		tokenomics.UserIDField
		tokenomics.UTCOffsetField
	}
	UpdatedUser struct {
		tokenomics.ExtraBonusLastClaimAvailableAtField
		tokenomics.DeserializedUsersKey
		tokenomics.ExtraBonusDaysClaimNotAvailableField
		extraBonusIndex uint16 `redis:"-"`
	}

	extraBonusNotifier struct {
		db                            storage.DB
		mb                            messagebroker.Client
		extraBonusStartDate           *time.Time
		extraBonusIndicesDistribution map[uint16]map[uint16]uint16
	}
)
