// SPDX-License-Identifier: ice License 1.0

package miner

import (
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/tokenomics"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
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
	applicationYamlKey = "miner"
	requestDeadline    = 30 * stdlibtime.Second
)

// .
var (
	//nolint:gochecknoglobals // Singleton & global config mounted only during bootstrap.
	cfg struct {
		tokenomics.Config `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		Development       bool                     `yaml:"development"`
		Workers           int64                    `yaml:"workers"`
		BatchSize         int64                    `yaml:"batchSize"`
	}
)

type (
	user struct {
		tokenomics.MiningSessionSoloLastStartedAtField
		tokenomics.MiningSessionSoloStartedAtField
		tokenomics.MiningSessionSoloEndedAtField
		tokenomics.PreviousMiningSessionSoloEndedAtField
		tokenomics.ExtraBonusStartedAtField
		tokenomics.UserIDField
		UpdatedUser
		tokenomics.BalanceSoloPendingField
		tokenomics.BalanceT1PendingField
		tokenomics.BalanceT2PendingField
		tokenomics.ActiveT1ReferralsField
		tokenomics.ActiveT2ReferralsField
		tokenomics.PreStakingBonusField
		tokenomics.PreStakingAllocationField
		tokenomics.ExtraBonusField
	}

	UpdatedUser struct { // This is public only because we have to embed it, and it has to be if so.
		tokenomics.BalanceLastUpdatedAtField
		tokenomics.ResurrectSoloUsedAtField
		tokenomics.ResurrectT0UsedAtField
		tokenomics.ResurrectTMinus1UsedAtField
		tokenomics.DeserializedUsersKey
		tokenomics.IDT0Field
		tokenomics.IDTMinus1Field
		tokenomics.BalanceTotalStandardField
		tokenomics.BalanceTotalPreStakingField
		tokenomics.BalanceTotalMintedField
		tokenomics.BalanceTotalSlashedField
		tokenomics.BalanceSoloPendingAppliedField
		tokenomics.BalanceT1PendingAppliedField
		tokenomics.BalanceT2PendingAppliedField
		tokenomics.BalanceSoloField
		tokenomics.BalanceT0Field
		tokenomics.BalanceT1Field
		tokenomics.BalanceT2Field
		tokenomics.BalanceForT0Field
		tokenomics.BalanceForTMinus1Field
		tokenomics.SlashingRateSoloField
		tokenomics.SlashingRateT0Field
		tokenomics.SlashingRateT1Field
		tokenomics.SlashingRateT2Field
		tokenomics.SlashingRateForT0Field
		tokenomics.SlashingRateForTMinus1Field
	}

	referral struct {
		tokenomics.MiningSessionSoloStartedAtField
		tokenomics.MiningSessionSoloEndedAtField
		tokenomics.PreviousMiningSessionSoloEndedAtField
		tokenomics.ResurrectSoloUsedAtField
		tokenomics.DeserializedUsersKey
	}

	referralThatStoppedMining struct {
		StoppedMiningAt     *time.Time
		ID, IDT0, IDTMinus1 int64
	}

	miner struct {
		db storage.DB
		mb messagebroker.Client
	}
)
