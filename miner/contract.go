// SPDX-License-Identifier: ice License 1.0

package miner

import (
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

type (
	DayOffStarted struct {
		StartedAt                   *time.Time `json:"startedAt,omitempty"`
		EndedAt                     *time.Time `json:"endedAt,omitempty"`
		UserID                      string     `json:"userId,omitempty" `
		ID                          string     `json:"id,omitempty"`
		RemainingFreeMiningSessions uint64     `json:"remainingFreeMiningSessions,omitempty"`
		MiningStreak                uint64     `json:"miningStreak,omitempty"`
	}
)

// Private API.

const (
	applicationYamlKey       = "miner"
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
		Development       bool                     `yaml:"development"`
	}
)

type (
	user struct {
		model.MiningSessionSoloLastStartedAtField
		model.MiningSessionSoloStartedAtField
		model.MiningSessionSoloEndedAtField
		model.MiningSessionSoloPreviouslyEndedAtField
		model.ExtraBonusStartedAtField
		model.UserIDField
		UpdatedUser
		model.BalanceSoloPendingField
		model.BalanceT1PendingField
		model.BalanceT2PendingField
		model.ActiveT1ReferralsField
		model.ActiveT2ReferralsField
		model.PreStakingBonusField
		model.PreStakingAllocationField
		model.ExtraBonusField
		model.UTCOffsetField
	}

	UpdatedUser struct { // This is public only because we have to embed it, and it has to be if so.
		model.ExtraBonusLastClaimAvailableAtField
		model.BalanceLastUpdatedAtField
		model.ResurrectSoloUsedAtField
		model.ResurrectT0UsedAtField
		model.ResurrectTMinus1UsedAtField
		model.DeserializedUsersKey
		model.IDT0Field
		model.IDTMinus1Field
		model.BalanceTotalStandardField
		model.BalanceTotalPreStakingField
		model.BalanceTotalMintedField
		model.BalanceTotalSlashedField
		model.BalanceSoloPendingAppliedField
		model.BalanceT1PendingAppliedField
		model.BalanceT2PendingAppliedField
		model.BalanceSoloField
		model.BalanceT0Field
		model.BalanceT1Field
		model.BalanceT2Field
		model.BalanceForT0Field
		model.BalanceForTMinus1Field
		model.SlashingRateSoloField
		model.SlashingRateT0Field
		model.SlashingRateT1Field
		model.SlashingRateT2Field
		model.SlashingRateForT0Field
		model.SlashingRateForTMinus1Field
		model.ExtraBonusDaysClaimNotAvailableField
	}

	referral struct {
		model.MiningSessionSoloStartedAtField
		model.MiningSessionSoloEndedAtField
		model.MiningSessionSoloPreviouslyEndedAtField
		model.ResurrectSoloUsedAtField
		model.DeserializedUsersKey
	}

	referralThatStoppedMining struct {
		StoppedMiningAt     *time.Time
		ID, IDT0, IDTMinus1 int64
	}

	miner struct {
		mb messagebroker.Client
	}
)
