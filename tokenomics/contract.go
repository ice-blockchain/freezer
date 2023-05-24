// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	_ "embed"
	"io"
	stdlibtime "time"

	"github.com/pkg/errors"

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
	PreStakingBonusesPerYear                           = map[uint8]uint16{
		1: 35,
		2: 70,
		3: 115,
		4: 170,
		5: 250,
	}
	PreStakingYearsByPreStakingBonuses = map[uint16]uint8{
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
		Bonus    int64   `json:"bonus" example:"120"`
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
		Bonus uint64 `json:"bonus,omitempty" example:"100"`
	}
	PreStaking struct {
		UserID     string `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Years      uint64 `json:"years,omitempty" example:"1"`
		Allocation uint64 `json:"allocation,omitempty" example:"100"`
	}
	MiningRateBonuses struct {
		T1         uint64 `json:"t1,omitempty" example:"100"`
		T2         uint64 `json:"t2,omitempty" example:"200"`
		PreStaking uint64 `json:"preStaking,omitempty" example:"300"`
		Extra      uint64 `json:"extra,omitempty" example:"300"`
		Total      uint64 `json:"total,omitempty" example:"300"`
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
		UserID              string `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		AvailableExtraBonus uint16 `json:"availableExtraBonus,omitempty" example:"2"`
	}
	RankingSummary struct {
		GlobalRank uint64 `json:"globalRank,omitempty" example:"12333"`
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

// DB fields.
type (
	BalanceLastUpdatedAtField struct {
		BalanceLastUpdatedAt *time.Time `redis:"balance_last_updated_at,omitempty"`
	}
	MiningSessionSoloLastStartedAtField struct {
		MiningSessionSoloLastStartedAt *time.Time `redis:"mining_session_solo_last_started_at,omitempty"`
	}
	MiningSessionSoloStartedAtField struct {
		MiningSessionSoloStartedAt *time.Time `redis:"mining_session_solo_started_at,omitempty"`
	}
	MiningSessionSoloEndedAtField struct {
		MiningSessionSoloEndedAt *time.Time `redis:"mining_session_solo_ended_at,omitempty"`
	}
	MiningSessionSoloPreviouslyEndedAtField struct {
		MiningSessionSoloPreviouslyEndedAt *time.Time `redis:"mining_session_solo_previously_ended_at,omitempty"`
	}
	ExtraBonusStartedAtField struct {
		ExtraBonusStartedAt *time.Time `redis:"extra_bonus_started_at,omitempty"`
	}
	ResurrectSoloUsedAtField struct {
		ResurrectSoloUsedAt *time.Time `redis:"resurrect_solo_used_at,omitempty"`
	}
	ResurrectT0UsedAtField struct {
		ResurrectT0UsedAt *time.Time `redis:"resurrect_t0_used_at,omitempty"`
	}
	ResurrectTMinus1UsedAtField struct {
		ResurrectTMinus1UsedAt *time.Time `redis:"resurrect_tminus1_used_at,omitempty"`
	}
	MiningSessionSoloDayOffLastAwardedAtField struct {
		MiningSessionSoloDayOffLastAwardedAt *time.Time `redis:"mining_session_solo_day_off_last_awarded_at,omitempty"`
	}
	ExtraBonusLastClaimAvailableAtField struct {
		ExtraBonusLastClaimAvailableAt *time.Time `redis:"extra_bonus_last_claim_available_at,omitempty"`
	}
	UserIDField struct {
		UserID string `redis:"user_id"`
	}
	ProfilePictureNameField struct {
		ProfilePictureName string `redis:"profile_picture_name,omitempty"`
	}
	UsernameField struct {
		Username string `redis:"username,omitempty"`
	}
	MiningBlockchainAccountAddressField struct {
		MiningBlockchainAccountAddress string `redis:"mining_blockchain_account_address,omitempty"`
	}
	BlockchainAccountAddressField struct {
		BlockchainAccountAddress string `redis:"blockchain_account_address,omitempty"`
	}
	BalanceTotalStandardField struct {
		BalanceTotalStandard float64 `redis:"balance_total_standard"`
	}
	BalanceTotalPreStakingField struct {
		BalanceTotalPreStaking float64 `redis:"balance_total_pre_staking"`
	}
	BalanceTotalMintedField struct {
		BalanceTotalMinted float64 `redis:"balance_total_minted"`
	}
	BalanceTotalSlashedField struct {
		BalanceTotalSlashed float64 `redis:"balance_total_slashed"`
	}
	BalanceSoloPendingField struct {
		BalanceSoloPending float64 `redis:"balance_solo_pending,omitempty"`
	}
	BalanceT1PendingField struct {
		BalanceT1Pending float64 `redis:"balance_t1_pending,omitempty"`
	}
	BalanceT2PendingField struct {
		BalanceT2Pending float64 `redis:"balance_t2_pending,omitempty"`
	}
	BalanceSoloPendingAppliedField struct {
		BalanceSoloPendingApplied float64 `redis:"balance_solo_pending_applied,omitempty"`
	}
	BalanceT1PendingAppliedField struct {
		BalanceT1PendingApplied float64 `redis:"balance_t1_pending_applied,omitempty"`
	}
	BalanceT2PendingAppliedField struct {
		BalanceT2PendingApplied float64 `redis:"balance_t2_pending_applied,omitempty"`
	}
	BalanceSoloField struct {
		BalanceSolo float64 `redis:"balance_solo"`
	}
	BalanceT0Field struct {
		BalanceT0 float64 `redis:"balance_t0"`
	}
	BalanceT1Field struct {
		BalanceT1 float64 `redis:"balance_t1"`
	}
	BalanceT2Field struct {
		BalanceT2 float64 `redis:"balance_t2"`
	}
	BalanceForT0Field struct {
		BalanceForT0 float64 `redis:"balance_for_t0"`
	}
	BalanceForTMinus1Field struct {
		BalanceForTMinus1 float64 `redis:"balance_for_tminus1"`
	}
	SlashingRateSoloField struct {
		SlashingRateSolo float64 `redis:"slashing_rate_solo"`
	}
	SlashingRateT0Field struct {
		SlashingRateT0 float64 `redis:"slashing_rate_t0"`
	}
	SlashingRateT1Field struct {
		SlashingRateT1 float64 `redis:"slashing_rate_t1"`
	}
	SlashingRateT2Field struct {
		SlashingRateT2 float64 `redis:"slashing_rate_t2"`
	}
	SlashingRateForT0Field struct {
		SlashingRateForT0 float64 `redis:"slashing_rate_for_t0"`
	}
	SlashingRateForTMinus1Field struct {
		SlashingRateForTMinus1 float64 `redis:"slashing_rate_for_tminus1"`
	}
	DeserializedUsersKey struct {
		HistoryPart string `redis:"-"`
		ID          int64  `redis:"-"`
	}
	IDT0Field struct {
		IDT0 int64 `redis:"id_t0,omitempty"`
	}
	IDTMinus1Field struct {
		IDTMinus1 int64 `redis:"id_tminus1,omitempty"`
	}
	IDT0ResettableField struct {
		IDT0 int64 `redis:"id_t0"`
	}
	IDTMinus1ResettableField struct {
		IDTMinus1 int64 `redis:"id_tminus1"`
	}
	ActiveT1ReferralsField struct {
		ActiveT1Referrals int32 `redis:"active_t1_referrals,omitempty"`
	}
	ActiveT2ReferralsField struct {
		ActiveT2Referrals int32 `redis:"active_t2_referrals,omitempty"`
	}
	PreStakingBonusField struct {
		PreStakingBonus uint16 `redis:"pre_staking_bonus,omitempty"`
	}
	PreStakingAllocationField struct {
		PreStakingAllocation uint16 `redis:"pre_staking_allocation,omitempty"`
	}
	ExtraBonusField struct {
		ExtraBonus uint16 `redis:"extra_bonus,omitempty"`
	}
	NewsSeenField struct {
		NewsSeen uint16 `redis:"news_seen"`
	}
	ExtraBonusDaysClaimNotAvailableField struct {
		ExtraBonusDaysClaimNotAvailable uint16 `redis:"extra_bonus_days_claim_not_available"`
	}
	UTCOffsetField struct {
		UTCOffset int16 `redis:"utc_offset"`
	}
	HideRankingField struct {
		HideRanking bool `redis:"hide_ranking"`
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
)

type (
	balance struct {
		UpdatedAt   *time.Time `json:"updatedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		UserID      string     `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		TypeDetail  string     `json:"typeDetail,omitempty" example:"/2022-01-03"`
		Type        string     `json:"type,omitempty" example:"1"`
		Amount      float64    `json:"amount,omitempty" example:"1,235.777777777"`
		HashCode    int64      `json:"hashCode,omitempty" example:"11"`
		WorkerIndex int16      `json:"workerIndex,omitempty" example:"11"`
		Negative    bool       `json:"negative,omitempty" example:"false"`
	}
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
		cfg                 *Config
		extraBonusStartDate *time.Time
		shutdown            func() error
		db                  storage.DB
		mb                  messagebroker.Client
		pictureClient       picture.Client
	}

	processor struct {
		*repository
	}

	Config struct {
		messagebroker.Config    `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		AdoptionMilestoneSwitch struct {
			ActiveUserMilestones []struct {
				Users          uint64  `yaml:"users"`
				BaseMiningRate float64 `yaml:"baseMiningRate"`
			} `yaml:"activeUserMilestones"`
			ConsecutiveDurationsRequired uint64              `yaml:"consecutiveDurationsRequired"`
			Duration                     stdlibtime.Duration `yaml:"duration"`
		} `yaml:"adoptionMilestoneSwitch"`
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
		RollbackNegativeMining struct {
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
