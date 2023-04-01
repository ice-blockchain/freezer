// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	_ "embed"
	"io"
	"sync"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/go-tarantool-client"
	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
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
	ErrNotFound                                        = storagev2.ErrNotFound
	ErrRelationNotFound                                = storagev2.ErrRelationNotFound
	ErrDuplicate                                       = storagev2.ErrDuplicate
	ErrNegativeMiningProgressDecisionRequired          = errors.New("you have negative mining progress, please decide what to do with it")
	ErrRaceCondition                                   = errors.New("race condition")
	ErrGlobalRankHidden                                = errors.New("global rank is hidden")
	ErrDecreasingPreStakingAllocationOrYearsNotAllowed = errors.New("decreasing pre-staking allocation or years not allowed")
)

type (
	MiningRateType    string
	AddBalanceCommand struct {
		*Balances[coin.ICEFlake]
		Negative *bool  `json:"negative,omitempty" example:"false"`
		EventID  string `json:"eventId,omitempty" example:"some unique id"`
	}
	Miner struct {
		_msgpack          struct{}  `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
		Balance           *coin.ICE `json:"balance,omitempty" example:"12345.6334"`
		UserID            string    `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Username          string    `json:"username,omitempty" example:"jdoe"`
		ProfilePictureURL string    `json:"profilePictureUrl,omitempty" example:"https://somecdn.com/p1.jpg"`
	}
	BalanceSummary struct {
		Balances[coin.ICE]
	}
	Balances[DENOM coin.ICEFlake | coin.ICE] struct {
		Total                          *DENOM `json:"total,omitempty" swaggertype:"string" example:"1,243.02"`
		BaseFactor                     *DENOM `json:"baseFactor,omitempty" swaggerignore:"true" swaggertype:"string" example:"1,243.02"`
		Standard                       *DENOM `json:"standard,omitempty" swaggertype:"string" example:"1,243.02"`
		PreStaking                     *DENOM `json:"preStaking,omitempty" swaggertype:"string" example:"1,243.02"`
		TotalNoPreStakingBonus         *DENOM `json:"totalNoPreStakingBonus,omitempty" swaggertype:"string" example:"1,243.02"`
		T1                             *DENOM `json:"t1,omitempty" swaggertype:"string" example:"1,243.02"`
		T2                             *DENOM `json:"t2,omitempty" swaggertype:"string" example:"1,243.02"`
		TotalReferrals                 *DENOM `json:"totalReferrals,omitempty" swaggertype:"string" example:"1,243.02"`
		UserID                         string `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		miningBlockchainAccountAddress string
	}
	BalanceHistoryBalanceDiff struct {
		Amount   *coin.ICE      `json:"amount" swaggertype:"string" example:"1,243.02"`
		amount   *coin.ICEFlake //nolint:revive // That's intended.
		Bonus    int64          `json:"bonus" example:"120"`
		Negative bool           `json:"negative" example:"true"`
	}
	BalanceHistoryEntry struct {
		Time       stdlibtime.Time            `json:"time" swaggertype:"string" example:"2022-01-03T16:20:52.156534Z"`
		Balance    *BalanceHistoryBalanceDiff `json:"balance"`
		TimeSeries []*BalanceHistoryEntry     `json:"timeSeries"`
	}
	AdoptionSummary struct {
		Milestones       []*Adoption[coin.ICE] `json:"milestones"`
		TotalActiveUsers uint64                `json:"totalActiveUsers" example:"11"`
	}
	AdoptionSnapshot struct {
		*Adoption[coin.ICEFlake]
		Before *Adoption[coin.ICEFlake] `json:"before,omitempty"`
	}
	Adoption[DENOM coin.ICEFlake | coin.ICE] struct {
		_msgpack         struct{}   `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
		AchievedAt       *time.Time `json:"achievedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		BaseMiningRate   *DENOM     `json:"baseMiningRate,omitempty" swaggertype:"string" example:"1,243.02"`
		Milestone        uint64     `json:"milestone,omitempty" example:"1"`
		TotalActiveUsers uint64     `json:"totalActiveUsers,omitempty" example:"1"`
	}
	PreStakingSummary struct {
		*PreStaking
		Bonus uint64 `json:"bonus,omitempty" example:"100"`
	}
	PreStakingSnapshot struct {
		*PreStakingSummary
		Before *PreStakingSummary `json:"before,omitempty"`
	}
	PreStaking struct {
		CreatedAt   *time.Time `json:"createdAt,omitempty" swaggerignore:"true" example:"2022-01-03T16:20:52.156534Z"`
		UserID      string     `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Years       uint64     `json:"years,omitempty" example:"1"`
		Allocation  uint64     `json:"allocation,omitempty" example:"100"`
		HashCode    int64      `json:"-" example:"11"`
		WorkerIndex int16      `json:"-" example:"11"`
	}
	MiningRateBonuses struct {
		T1         uint64 `json:"t1,omitempty" example:"100"`
		T2         uint64 `json:"t2,omitempty" example:"200"`
		PreStaking uint64 `json:"preStaking,omitempty" example:"300"`
		Extra      uint64 `json:"extra,omitempty" example:"300"`
		Total      uint64 `json:"total,omitempty" example:"300"`
	}
	MiningRateSummary[DENOM coin.ICEFlake | coin.ICE] struct {
		Amount  *DENOM             `json:"amount,omitempty" example:"1,234,232.001" swaggertype:"string"`
		Bonuses *MiningRateBonuses `json:"bonuses,omitempty"`
	}
	MiningRates[T coin.ICEFlake | MiningRateSummary[coin.ICE]] struct {
		_msgpack                       struct{}       `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
		Total                          *T             `json:"total,omitempty"`
		TotalNoPreStakingBonus         *T             `json:"totalNoPreStakingBonus,omitempty"`
		PositiveTotalNoPreStakingBonus *T             `json:"positiveTotalNoPreStakingBonus,omitempty"`
		Standard                       *T             `json:"standard,omitempty"`
		PreStaking                     *T             `json:"preStaking,omitempty"`
		Base                           *T             `json:"base,omitempty"`
		Type                           MiningRateType `json:"type,omitempty"`
		UserID                         string         `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	MiningSummary struct {
		MiningRates   *MiningRates[MiningRateSummary[coin.ICE]] `json:"miningRates,omitempty"`
		MiningSession *MiningSession                            `json:"miningSession,omitempty"`
		ExtraBonusSummary
		MiningStreak                uint64 `json:"miningStreak,omitempty"  example:"2"`
		RemainingFreeMiningSessions uint64 `json:"remainingFreeMiningSessions,omitempty" example:"1"`
	}
	MiningSession struct {
		_msgpack                      struct{}   `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
		LastNaturalMiningStartedAt    *time.Time `json:"lastNaturalMiningStartedAt,omitempty" example:"2022-01-03T16:20:52.156534Z" swaggerignore:"true"`
		StartedAt                     *time.Time `json:"startedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		EndedAt                       *time.Time `json:"endedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		ResettableStartingAt          *time.Time `json:"resettableStartingAt,omitempty" example:"2022-01-03T16:20:52.156534Z" `
		WarnAboutExpirationStartingAt *time.Time `json:"warnAboutExpirationStartingAt,omitempty" example:"2022-01-03T16:20:52.156534Z" `
		Free                          *bool      `json:"free,omitempty" example:"true"`
		UserID                        *string    `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		MiningStreak                  uint64     `json:"miningStreak,omitempty" swaggerignore:"true" example:"11"`
	}
	ExtraBonusSummary struct {
		UserID              string `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		AvailableExtraBonus uint64 `json:"availableExtraBonus,omitempty" example:"2"`
		ExtraBonusIndex     uint64 `json:"extraBonusIndex,omitempty" swaggerignore:"true" example:"1"`
	}
	RankingSummary struct {
		GlobalRank uint64 `json:"globalRank,omitempty" example:"12333"`
	}
	FreeMiningSessionStarted struct {
		StartedAt                   *time.Time `json:"startedAt,omitempty"`
		EndedAt                     *time.Time `json:"endedAt,omitempty"`
		UserID                      string     `json:"userId,omitempty" `
		ID                          string     `json:"id,omitempty"`
		RemainingFreeMiningSessions uint64     `json:"remainingFreeMiningSessions,omitempty"`
		MiningStreak                uint64     `json:"miningStreak,omitempty" example:"11"`
	}
	ReadRepository interface {
		GetBalanceSummary(ctx context.Context, userID string) (*BalanceSummary, error)
		GetRankingSummary(ctx context.Context, userID string) (*RankingSummary, error)
		GetTopMiners(ctx context.Context, keyword string, limit, offset uint64) ([]*Miner, error)
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

		ReadRepository
		WriteRepository
	}
	Processor interface {
		Repository
		CheckHealth(context.Context) error
	}
)

// Private API.

const (
	applicationYamlKey                                         = "tokenomics"
	dayFormat, hourFormat, minuteFormat                        = "2006-01-02", "2006-01-02T15", "2006-01-02T15:04"
	totalActiveUsersGlobalKey                                  = "TOTAL_ACTIVE_USERS"
	requestingUserIDCtxValueKey                                = "requestingUserIDCtxValueKey"
	userHashCodeCtxValueKey                                    = "userHashCodeCtxValueKey"
	percentage100                                              = uint64(100)
	registrationICEFlakeBonusAmount                            = 10 * uint64(coin.Denomination)
	lastAdoptionMilestone                                      = 6
	miningRatesRecalculationBatchSize                          = 10
	balanceRecalculationBatchSize                              = 10
	extraBonusProcessingBatchSize                              = 100
	maxICEBlockchainConcurrentOperations                       = 10000
	balanceCalculationProcessingSeedingStreamEmitFrequency     = 0 * stdlibtime.Second
	refreshMiningRatesProcessingSeedingStreamEmitFrequency     = 0 * stdlibtime.Second
	blockchainBalanceSynchronizationSeedingStreamEmitFrequency = 0 * stdlibtime.Second
	extraBonusProcessingSeedingStreamEmitFrequency             = 0 * stdlibtime.Second
	requestDeadline                                            = 25 * stdlibtime.Second
)

const (
	t0BalanceTypeDetail                                  = "t0"
	t1BalanceTypeDetail                                  = "t1"
	t2BalanceTypeDetail                                  = "t2"
	degradationT0T1T2TotalReferenceBalanceTypeDetail     = "@&"
	aggressiveDegradationTotalReferenceBalanceTypeDetail = "_"
	aggressiveDegradationT1ReferenceBalanceTypeDetail    = t1BalanceTypeDetail + "_"
	aggressiveDegradationT2ReferenceBalanceTypeDetail    = t2BalanceTypeDetail + "_"
	reverseT0BalanceTypeDetail                           = "&" + t0BalanceTypeDetail
	reverseTMinus1BalanceTypeDetail                      = "&t-1"
)

const (
	totalNoPreStakingBonusBalanceType balanceType = iota
	pendingXBalanceType
)

// .
var (
	//go:embed DDL.lua
	ddl string
	//go:embed DDL.sql
	ddlV2 string
)

type (
	balanceType                           int8
	userMiningRateRecalculationParameters struct {
		UserID                                                        users.UserID
		T0, T1, T2, ExtraBonus, PreStakingAllocation, PreStakingBonus uint64
	}
	user struct {
		CreatedAt                      *time.Time `json:"createdAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		UpdatedAt                      *time.Time `json:"updatedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		RollbackUsedAt                 *time.Time `json:"rollbackUsedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		LastNaturalMiningStartedAt     *time.Time `json:"lastNaturalMiningStartedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		LastMiningStartedAt            *time.Time `json:"lastMiningStartedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		LastMiningEndedAt              *time.Time `json:"lastMiningEndedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		PreviousMiningStartedAt        *time.Time `json:"previousMiningStartedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		PreviousMiningEndedAt          *time.Time `json:"previousMiningEndedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		LastFreeMiningSessionAwardedAt *time.Time `json:"lastFreeMiningSessionAwardedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		UserID                         string     `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		ReferredBy                     string     `json:"referredBy,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Username                       string     `json:"username,omitempty" example:"jdoe"`
		FirstName                      string     `json:"firstName,omitempty" example:"John"`
		LastName                       string     `json:"lastName,omitempty" example:"Doe"`
		ProfilePictureURL              string     `json:"profilePictureUrl,omitempty" example:"https://somecdn.com/p1.jpg"`
		MiningBlockchainAccountAddress string     `json:"miningBlockchainAccountAddress,omitempty" example:"0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		BlockchainAccountAddress       string     `json:"blockchainAccountAddress,omitempty" example:"0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		HashCode                       int64      `json:"hashCode,omitempty" example:"1234567890"`
		HideRanking                    bool       `json:"hideRanking,omitempty" example:"false"`
		Verified                       bool       `json:"verified,omitempty" example:"false"`
	}
	balance struct {
		UpdatedAt   *time.Time     `json:"updatedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		Amount      *coin.ICEFlake `json:"amount,omitempty" example:"1,235.777777777"`
		UserID      string         `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		TypeDetail  string         `json:"typeDetail,omitempty" example:"/2022-01-03"`
		Type        balanceType    `json:"type,omitempty" example:"1"`
		Negative    bool           `json:"negative,omitempty" example:"false"`
		HashCode    int64          `json:"hashCode,omitempty" example:"11"`
		WorkerIndex int16          `json:"workerIndex,omitempty" example:"11"`
	}
	miningSummary struct {
		LastNaturalMiningStartedAt                    *time.Time     `json:"lastNaturalMiningStartedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		LastMiningStartedAt                           *time.Time     `json:"lastMiningStartedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		LastMiningEndedAt                             *time.Time     `json:"lastMiningEndedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		PreviousMiningStartedAt                       *time.Time     `json:"previousMiningStartedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		PreviousMiningEndedAt                         *time.Time     `json:"previousMiningEndedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		LastFreeMiningSessionAwardedAt                *time.Time     `json:"lastFreeMiningSessionAwardedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		NegativeTotalNoPreStakingBonusBalanceAmount   *coin.ICEFlake `json:"negativeTotalNoPreStakingBonusBalanceAmount,omitempty" example:"1,235.777777777"`
		NegativeTotalT0NoPreStakingBonusBalanceAmount *coin.ICEFlake `json:"negativeTotalT0NoPreStakingBonusBalanceAmount,omitempty" example:"1,235.777777777"`
		NegativeTotalT1NoPreStakingBonusBalanceAmount *coin.ICEFlake `json:"negativeTotalT1NoPreStakingBonusBalanceAmount,omitempty" example:"1,235.777777777"`
		NegativeTotalT2NoPreStakingBonusBalanceAmount *coin.ICEFlake `json:"negativeTotalT2NoPreStakingBonusBalanceAmount,omitempty" example:"1,235.777777777"`
		MiningStreak                                  uint64         `json:"miningStreak,omitempty" example:"11"`
		PreStakingYears                               uint64         `json:"preStakingYears,omitempty" example:"11"`
		PreStakingAllocation                          uint64         `json:"preStakingAllocation,omitempty" example:"11"`
		PreStakingBonus                               uint64         `json:"preStakingBonus,omitempty" example:"11"`
	}
	deviceMetadata struct {
		_msgpack struct{}        `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
		Before   *deviceMetadata `json:"before,omitempty"`
		UserID   string          `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		TZ       string          `json:"tz,omitempty" example:"+03:00"`
	}
	viewedNews struct {
		_msgpack struct{} `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
		UserID   string   `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		NewsID   string   `json:"newsId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}

	usersTableSource struct {
		*processor
	}

	globalTableSource struct {
		*processor
	}

	miningSessionsTableSource struct {
		*processor
	}

	addBalanceCommandsSource struct {
		*processor
	}

	viewedNewsSource struct {
		*processor
	}

	deviceMetadataTableSource struct {
		*processor
	}

	balanceRecalculationTriggerStreamSource struct {
		*processor
	}

	miningRatesRecalculationTriggerStreamSource struct {
		*processor
	}

	blockchainBalanceSynchronizationTriggerStreamSource struct {
		*processor
	}

	extraBonusProcessingTriggerStreamSource struct {
		*processor
	}

	repository struct {
		cfg           *config
		shutdown      func() error
		db            tarantool.Connector
		dbV2          *storagev2.DB
		mb            messagebroker.Client
		pictureClient picture.Client
	}

	processor struct {
		*repository
		streamsDoneWg *sync.WaitGroup
		cancelStreams context.CancelFunc
	}

	config struct {
		messagebroker.Config    `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		AdoptionMilestoneSwitch struct {
			ActiveUserMilestones         []uint64            `yaml:"activeUserMilestones"`
			ConsecutiveDurationsRequired uint64              `yaml:"consecutiveDurationsRequired"`
			Duration                     stdlibtime.Duration `yaml:"duration"`
		} `yaml:"adoptionMilestoneSwitch"`
		ExtraBonuses struct {
			FlatValues                []uint64            `yaml:"flatValues"`
			NewsSeenValues            []uint64            `yaml:"newsSeenValues"`
			MiningStreakValues        []uint64            `yaml:"miningStreakValues"`
			Duration                  stdlibtime.Duration `yaml:"duration"`
			UTCOffsetDuration         stdlibtime.Duration `yaml:"utcOffsetDuration" mapstructure:"utcOffsetDuration"`
			ClaimWindow               stdlibtime.Duration `yaml:"claimWindow"`
			DelayedClaimPenaltyWindow stdlibtime.Duration `yaml:"delayedClaimPenaltyWindow"`
			AvailabilityWindow        stdlibtime.Duration `yaml:"availabilityWindow"`
			TimeToAvailabilityWindow  stdlibtime.Duration `yaml:"timeToAvailabilityWindow"`
		} `yaml:"extraBonuses"`
		RollbackNegativeMining struct {
			Available struct {
				After stdlibtime.Duration `yaml:"after"`
				Until stdlibtime.Duration `yaml:"until"`
			} `yaml:"available"`
			LastXMiningSessionsCollectingInterval stdlibtime.Duration `yaml:"lastXMiningSessionsCollectingInterval" mapstructure:"lastXMiningSessionsCollectingInterval"`
			AggressiveDegradationStartsAfter      stdlibtime.Duration `yaml:"aggressiveDegradationStartsAfter"`
		} `yaml:"rollbackNegativeMining"`
		MiningSessionDuration struct {
			Min                      stdlibtime.Duration `yaml:"min"`
			Max                      stdlibtime.Duration `yaml:"max"`
			WarnAboutExpirationAfter stdlibtime.Duration `yaml:"warnAboutExpirationAfter"`
		} `yaml:"miningSessionDuration"`
		ReferralBonusMiningRates struct {
			T0 uint64 `yaml:"t0"`
			T1 uint64 `yaml:"t1"`
			T2 uint64 `yaml:"t2"`
		} `yaml:"referralBonusMiningRates"`
		ConsecutiveNaturalMiningSessionsRequiredFor1ExtraFreeArtificialMiningSession struct {
			Min uint64 `yaml:"min"`
			Max uint64 `yaml:"max"`
		} `yaml:"consecutiveNaturalMiningSessionsRequiredFor1ExtraFreeArtificialMiningSession"`
		GlobalAggregationInterval struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		} `yaml:"globalAggregationInterval"`
		WorkerCount int16 `yaml:"workerCount"`
	}
)
