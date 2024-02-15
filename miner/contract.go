// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	stdlibtime "time"

	"github.com/ice-blockchain/eskimo/kyc/quiz"
	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	coindistribution "github.com/ice-blockchain/freezer/coin-distribution"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

type (
	Client interface {
		io.Closer
		CheckHealth(context.Context) error
	}
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
	cfg config
)

type (
	user struct {
		model.MiningSessionSoloLastStartedAtField
		model.MiningSessionSoloStartedAtField
		model.MiningSessionSoloEndedAtField
		model.MiningSessionSoloPreviouslyEndedAtField
		model.ExtraBonusStartedAtField
		model.ReferralsCountChangeGuardUpdatedAtField
		model.KYCState
		model.MiningBlockchainAccountAddressField
		model.CountryField
		model.UsernameField
		model.LatestDeviceField
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
		model.UserIDField
		model.ExtraBonusLastClaimAvailableAtField
		model.BalanceLastUpdatedAtField
		model.ResurrectSoloUsedAtField
		model.ResurrectT0UsedAtField
		model.ResurrectTMinus1UsedAtField
		model.SoloLastEthereumCoinDistributionProcessedAtField
		model.ForT0LastEthereumCoinDistributionProcessedAtField
		model.ForTMinus1LastEthereumCoinDistributionProcessedAtField
		model.BalanceSoloEthereumPendingField
		model.BalanceT0EthereumPendingField
		model.BalanceT1EthereumPendingField
		model.BalanceT2EthereumPendingField
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
		model.BalanceSoloEthereumField
		model.BalanceT0EthereumField
		model.BalanceT1EthereumField
		model.BalanceT2EthereumField
		model.BalanceForT0EthereumField
		model.BalanceForTMinus1EthereumField
		model.BalanceSoloEthereumMainnetRewardPoolContributionField
		model.BalanceT0EthereumMainnetRewardPoolContributionField
		model.BalanceT1EthereumMainnetRewardPoolContributionField
		model.BalanceT2EthereumMainnetRewardPoolContributionField
		model.SlashingRateSoloField
		model.SlashingRateT0Field
		model.SlashingRateT1Field
		model.SlashingRateT2Field
		model.SlashingRateForT0Field
		model.SlashingRateForTMinus1Field
		model.ExtraBonusDaysClaimNotAvailableField
		model.KYCQuizDisabledResettableField
		model.KYCQuizCompletedResettableField
	}
	referralUpdated struct {
		model.DeserializedUsersKey
		model.IDT0Field
		model.IDTMinus1Field
	}

	referral struct {
		model.MiningSessionSoloStartedAtField
		model.MiningSessionSoloEndedAtField
		model.MiningSessionSoloPreviouslyEndedAtField
		model.ResurrectSoloUsedAtField
		model.UserIDField
		model.CountryField
		model.UsernameField
		model.MiningBlockchainAccountAddressField
		model.KYCState
		model.IDT0Field
		model.DeserializedUsersKey
		model.BalanceTotalStandardField
		model.BalanceSoloEthereumField
		model.BalanceT0EthereumField
		model.BalanceT1EthereumField
		model.BalanceT2EthereumField
		model.PreStakingAllocationField
		model.PreStakingBonusField
	}

	referralCountGuardUpdatedUser struct {
		model.ReferralsCountChangeGuardUpdatedAtField
		model.DeserializedUsersKey
	}

	referralThatStoppedMining struct {
		StoppedMiningAt     *time.Time
		ID, IDT0, IDTMinus1 int64
	}

	prestakingResettableUpdatedUser struct {
		model.PreStakingAllocationResettableField
		model.PreStakingBonusResettableField
		model.DeserializedUsersKey
	}

	miner struct {
		coinDistributionStartedSignaler             chan struct{}
		coinDistributionEndedSignaler               chan struct{}
		stopCoinDistributionCollectionWorkerManager chan struct{}
		coinDistributionWorkerMX                    *sync.Mutex
		coinDistributionRepository                  coindistribution.Repository
		quizRepository                              quiz.ReadRepository
		mb                                          messagebroker.Client
		db                                          storage.DB
		dwhClient                                   dwh.Client
		cancel                                      context.CancelFunc
		telemetry                                   *telemetry
		wg                                          *sync.WaitGroup
		extraBonusStartDate                         *time.Time
		extraBonusIndicesDistribution               map[uint16]map[uint16]uint16
	}
	config struct {
		disableAdvancedTeam                *atomic.Pointer[[]string]
		coinDistributionCollectorStartedAt *atomic.Pointer[time.Time]
		coinDistributionCollectorSettings  *atomic.Pointer[coindistribution.CollectorSettings]
		tokenomics.Config                  `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		EthereumDistributionFrequency      struct {
			Min stdlibtime.Duration `yaml:"min"`
			Max stdlibtime.Duration `yaml:"max"`
		} `yaml:"ethereumDistributionFrequency" mapstructure:"ethereumDistributionFrequency"`
		MainnetRewardPoolContributionEthAddress string  `yaml:"mainnetRewardPoolContributionEthAddress" mapstructure:"mainnetRewardPoolContributionEthAddress"`
		MainnetRewardPoolContributionPercentage float64 `yaml:"mainnetRewardPoolContributionPercentage" mapstructure:"mainnetRewardPoolContributionPercentage"`
		Workers                                 int64   `yaml:"workers"`
		BatchSize                               int64   `yaml:"batchSize"`
		Development                             bool    `yaml:"development"`
	}
)
