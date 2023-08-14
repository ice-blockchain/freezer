// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	_ "embed"
	"io"
	stdlibtime "time"

	"github.com/pkg/errors"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	extrabonusnotifier "github.com/ice-blockchain/freezer/extra-bonus-notifier"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/multimedia/picture"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

const (
	MaxPreStakingYears = 5
)

const (
	PositiveMiningRateType MiningRateType = "positive"
	NegativeMiningRateType MiningRateType = "negative"
	NoneMiningRateType     MiningRateType = "none"
)

var (
	ErrNotFound                                        = errors.New("not found")
	ErrRelationNotFound                                = errors.New("relationship not found")
	ErrDuplicate                                       = errors.New("duplicate")
	ErrNegativeMiningProgressDecisionRequired          = errors.New("you have negative mining progress, please decide what to do with it")
	ErrRaceCondition                                   = errors.New("race condition")
	ErrGlobalRankHidden                                = errors.New("global rank is hidden")
	ErrDecreasingPreStakingAllocationOrYearsNotAllowed = errors.New("decreasing pre-staking allocation or years not allowed")
	PreStakingBonusesPerYear                           = map[uint8]float64{
		1: 35,
		2: 70,
		3: 115,
		4: 170,
		5: 250,
	}
	PreStakingYearsByPreStakingBonuses = map[float64]uint8{
		35:  1,
		70:  2,
		115: 3,
		170: 4,
		250: 5,
	}
)

type (
	MiningRateType string
	Miner          struct {
		Balance           string `json:"balance,omitempty" example:"12345.6334"`
		UserID            string `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Username          string `json:"username,omitempty" example:"jdoe"`
		ProfilePictureURL string `json:"profilePictureUrl,omitempty" example:"https://somecdn.com/p1.jpg"`
		balance           float64
	}
	BalanceSummary struct {
		Balances[string]
	}
	Balances[DENOM ~float64 | ~string] struct {
		Total                          DENOM  `json:"total,omitempty" swaggertype:"string" example:"1,243.02"`
		BaseFactor                     DENOM  `json:"baseFactor,omitempty" swaggerignore:"true" swaggertype:"string" example:"1,243.02"`
		Standard                       DENOM  `json:"standard,omitempty" swaggertype:"string" example:"1,243.02"`
		PreStaking                     DENOM  `json:"preStaking,omitempty" swaggertype:"string" example:"1,243.02"`
		TotalNoPreStakingBonus         DENOM  `json:"totalNoPreStakingBonus,omitempty" swaggertype:"string" example:"1,243.02"`
		T1                             DENOM  `json:"t1,omitempty" swaggertype:"string" example:"1,243.02"`
		T2                             DENOM  `json:"t2,omitempty" swaggertype:"string" example:"1,243.02"`
		TotalReferrals                 DENOM  `json:"totalReferrals,omitempty" swaggertype:"string" example:"1,243.02"`
		UserID                         string `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		miningBlockchainAccountAddress string
	}
	BalanceHistoryBalanceDiff struct {
		Amount   string  `json:"amount" example:"1,243.02"`
		amount   float64 //nolint:revive // That's intended.
		Bonus    float64 `json:"bonus" example:"120.00"`
		Negative bool    `json:"negative" example:"true"`
	}
	BalanceHistoryEntry struct {
		Time       stdlibtime.Time            `json:"time" swaggertype:"string" example:"2022-01-03T16:20:52.156534Z"`
		Balance    *BalanceHistoryBalanceDiff `json:"balance"`
		TimeSeries []*BalanceHistoryEntry     `json:"timeSeries"`
	}
	AdoptionSummary struct {
		Milestones       []*Adoption[string] `json:"milestones"`
		TotalActiveUsers uint64              `json:"totalActiveUsers" example:"11"`
	}
	AdoptionSnapshot struct {
		*Adoption[float64]
		Before *Adoption[float64] `json:"before,omitempty"`
	}
	Adoption[DENOM ~string | ~float64] struct {
		AchievedAt       *time.Time `json:"achievedAt,omitempty" redis:"achieved_at" example:"2022-01-03T16:20:52.156534Z"`
		BaseMiningRate   DENOM      `json:"baseMiningRate,omitempty" redis:"base_mining_rate" swaggertype:"string" example:"1,243.02"`
		Milestone        uint64     `json:"milestone,omitempty" redis:"milestone" example:"1"`
		TotalActiveUsers uint64     `json:"totalActiveUsers,omitempty" redis:"total_active_users" example:"1"`
	}
	PreStakingSummary struct {
		*PreStaking
		Bonus float64 `json:"bonus,omitempty" example:"100.00"`
	}
	PreStaking struct {
		UserID     string  `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Years      uint64  `json:"years,omitempty" example:"1"`
		Allocation float64 `json:"allocation,omitempty" example:"100.00"`
	}
	MiningRateBonuses struct {
		T1         float64 `json:"t1,omitempty" example:"100.00"`
		T2         float64 `json:"t2,omitempty" example:"200.00"`
		PreStaking float64 `json:"preStaking,omitempty" example:"300.00"`
		Extra      float64 `json:"extra,omitempty" example:"300.00"`
		Total      float64 `json:"total,omitempty" example:"300.00"`
	}
	MiningRateSummary[DENOM ~string | ~float64] struct {
		Bonuses *MiningRateBonuses `json:"bonuses,omitempty"`
		Amount  DENOM              `json:"amount,omitempty" example:"1,234,232.001" swaggertype:"string"`
	}
	MiningRates[T float64 | *MiningRateSummary[string]] struct {
		Total                          T              `json:"total,omitempty"`
		TotalNoPreStakingBonus         T              `json:"totalNoPreStakingBonus,omitempty"`
		PositiveTotalNoPreStakingBonus T              `json:"positiveTotalNoPreStakingBonus,omitempty"`
		Standard                       T              `json:"standard,omitempty"`
		PreStaking                     T              `json:"preStaking,omitempty"`
		Base                           T              `json:"base,omitempty"`
		Type                           MiningRateType `json:"type,omitempty"`
		UserID                         string         `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	MiningSummary struct {
		MiningRates   *MiningRates[*MiningRateSummary[string]] `json:"miningRates,omitempty"`
		MiningSession *MiningSession                           `json:"miningSession,omitempty"`
		ExtraBonusSummary
		MiningStreak                uint64 `json:"miningStreak,omitempty"  example:"2"`
		RemainingFreeMiningSessions uint64 `json:"remainingFreeMiningSessions,omitempty" example:"1"`
	}
	MiningSession struct {
		LastNaturalMiningStartedAt    *time.Time          `json:"lastNaturalMiningStartedAt,omitempty" example:"2022-01-03T16:20:52.156534Z" swaggerignore:"true"`
		StartedAt                     *time.Time          `json:"startedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		EndedAt                       *time.Time          `json:"endedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		PreviouslyEndedAt             *time.Time          `json:"previouslyEndedAt,omitempty" swaggerignore:"true" example:"2022-01-03T16:20:52.156534Z"`
		ResettableStartingAt          *time.Time          `json:"resettableStartingAt,omitempty" example:"2022-01-03T16:20:52.156534Z" `
		WarnAboutExpirationStartingAt *time.Time          `json:"warnAboutExpirationStartingAt,omitempty" example:"2022-01-03T16:20:52.156534Z" `
		Free                          *bool               `json:"free,omitempty" example:"true"`
		UserID                        *string             `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Extension                     stdlibtime.Duration `json:"extension,omitempty" swaggerignore:"true" example:"24h"`
		MiningStreak                  uint64              `json:"miningStreak,omitempty" swaggerignore:"true" example:"11"`
	}
	ExtraBonusSummary struct {
		UserID              string  `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		AvailableExtraBonus float64 `json:"availableExtraBonus,omitempty" example:"2.00"`
	}
	RankingSummary struct {
		GlobalRank uint64 `json:"globalRank" example:"12333"`
	}
	ReadRepository interface {
		GetBalanceSummary(ctx context.Context, userID string) (*BalanceSummary, error)
		GetRankingSummary(ctx context.Context, userID string) (*RankingSummary, error)
		GetTopMiners(ctx context.Context, keyword string, limit, offset uint64) (topMiners []*Miner, nextOffset uint64, err error)
		GetMiningSummary(ctx context.Context, userID string) (*MiningSummary, error)
		GetPreStakingSummary(ctx context.Context, userID string) (*PreStakingSummary, error)
		GetBalanceHistory(ctx context.Context, userID string, start, end *time.Time, utcOffset stdlibtime.Duration, limit, offset uint64) ([]*BalanceHistoryEntry, error) //nolint:lll // .
		GetAdoptionSummary(context.Context) (*AdoptionSummary, error)
	}
	WriteRepository interface {
		StartNewMiningSession(ctx context.Context, ms *MiningSummary, rollbackNegativeMiningProgress *bool) error
		ClaimExtraBonus(ctx context.Context, ebs *ExtraBonusSummary) error
		StartOrUpdatePreStaking(context.Context, *PreStakingSummary) error
	}
	Repository interface {
		io.Closer
		CheckHealth(context.Context) error

		ReadRepository
		WriteRepository
	}
	Processor interface {
		Repository
	}
)

// Private API.

const (
	applicationYamlKey                  = "tokenomics"
	dayFormat, hourFormat, minuteFormat = "2006-01-02", "2006-01-02T15", "2006-01-02T15:04"
	totalActiveUsersGlobalKey           = "TOTAL_ACTIVE_USERS"
	requestingUserIDCtxValueKey         = "requestingUserIDCtxValueKey"
	userHashCodeCtxValueKey             = "userHashCodeCtxValueKey"
	requestDeadline                     = 25 * stdlibtime.Second

	floatToStringFormatter = "%.2f"
)

type (
	usersTableSource struct {
		*processor
	}

	miningSessionsTableSource struct {
		*processor
	}

	completedTasksSource struct {
		*processor
	}

	viewedNewsSource struct {
		*processor
	}

	deviceMetadataTableSource struct {
		*processor
	}

	repository struct {
		cfg                           *Config
		extraBonusStartDate           *time.Time
		extraBonusIndicesDistribution map[uint16]map[uint16]uint16
		shutdown                      func() error
		db                            storage.DB
		dwh                           dwh.Client
		mb                            messagebroker.Client
		pictureClient                 picture.Client
	}

	processor struct {
		*repository
	}

	Config struct {
		AdoptionMilestoneSwitch struct {
			ActiveUserMilestones []struct {
				Users          uint64  `yaml:"users"`
				BaseMiningRate float64 `yaml:"baseMiningRate"`
			} `yaml:"activeUserMilestones"`
			ConsecutiveDurationsRequired uint64              `yaml:"consecutiveDurationsRequired"`
			Duration                     stdlibtime.Duration `yaml:"duration"`
			Disabled                     bool                `yaml:"disabled"`
		} `yaml:"adoptionMilestoneSwitch"`
		messagebroker.Config                `mapstructure:",squash"`
		extrabonusnotifier.ExtraBonusConfig `mapstructure:",squash"`
		RollbackNegativeMining              struct {
			Available struct {
				After stdlibtime.Duration `yaml:"after"`
				Until stdlibtime.Duration `yaml:"until"`
			} `yaml:"available"`
		} `yaml:"rollbackNegativeMining"`
		MiningSessionDuration struct {
			Min                      stdlibtime.Duration `yaml:"min"`
			Max                      stdlibtime.Duration `yaml:"max"`
			WarnAboutExpirationAfter stdlibtime.Duration `yaml:"warnAboutExpirationAfter"`
		} `yaml:"miningSessionDuration"`
		ReferralBonusMiningRates struct {
			T0 uint16 `yaml:"t0"`
			T1 uint32 `yaml:"t1"`
			T2 uint32 `yaml:"t2"`
		} `yaml:"referralBonusMiningRates"`
		ConsecutiveNaturalMiningSessionsRequiredFor1ExtraFreeArtificialMiningSession struct {
			Min uint64 `yaml:"min"`
			Max uint64 `yaml:"max"`
		} `yaml:"consecutiveNaturalMiningSessionsRequiredFor1ExtraFreeArtificialMiningSession"`
		GlobalAggregationInterval struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		} `yaml:"globalAggregationInterval"`
	}
)
