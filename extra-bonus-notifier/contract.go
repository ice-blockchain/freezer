// SPDX-License-Identifier: ice License 1.0

package extrabonusnotifier

import (
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/model"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

type (
	User struct {
		model.ExtraBonusStartedAtField
		model.UserIDField
		UpdatedUser
		model.UTCOffsetField
	}
	UpdatedUser struct {
		model.ExtraBonusLastClaimAvailableAtField
		model.DeserializedUsersKey
		model.ExtraBonusDaysClaimNotAvailableResettableField
		ExtraBonusIndex uint16 `redis:"-"`
	}
	ExtraBonusAvailable struct {
		UserID          string `json:"userId,omitempty"`
		ExtraBonusIndex uint16 `json:"extraBonusIndex,omitempty"`
	}
	ExtraBonusConfig struct {
		ExtraBonuses struct {
			FlatValues                []uint16            `yaml:"flatValues"`
			NewsSeenValues            []uint16            `yaml:"newsSeenValues"`
			MiningStreakValues        []uint16            `yaml:"miningStreakValues"`
			Duration                  stdlibtime.Duration `yaml:"duration"`
			UTCOffsetDuration         stdlibtime.Duration `yaml:"utcOffsetDuration" mapstructure:"utcOffsetDuration"`
			ClaimWindow               stdlibtime.Duration `yaml:"claimWindow"`
			DelayedClaimPenaltyWindow stdlibtime.Duration `yaml:"delayedClaimPenaltyWindow"`
			AvailabilityWindow        stdlibtime.Duration `yaml:"availabilityWindow"`
			TimeToAvailabilityWindow  stdlibtime.Duration `yaml:"timeToAvailabilityWindow"`
		} `yaml:"extraBonuses"`
	}
)

// Private API.

const (
	applicationYamlKey       = "extra-bonus-notifier"
	parentApplicationYamlKey = "tokenomics"
	requestDeadline          = 30 * stdlibtime.Second
)

// .
var (
	//nolint:gochecknoglobals // Singleton & global config mounted only during bootstrap.
	cfg struct {
		messagebrokerConfig   `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		ExtraBonusConfig      `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		MiningSessionDuration stdlibtime.Duration      `yaml:"miningSessionDuration"`
		Workers               int64                    `yaml:"workers"`
		BatchSize             int64                    `yaml:"batchSize"`
		Chunks                uint16                   `yaml:"chunks"`
	}
)

type (
	messagebrokerConfig = messagebroker.Config
	extraBonusNotifier  struct {
		mb                            messagebroker.Client
		extraBonusStartDate           *time.Time
		extraBonusIndicesDistribution map[uint16]map[uint16]uint16
	}
)
