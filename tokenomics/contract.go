// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	_ "embed"
	"io"
	"sync/atomic"
	stdlibtime "time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	extrabonusnotifier "github.com/ice-blockchain/freezer/extra-bonus-notifier"
	detailedCoinMetrics "github.com/ice-blockchain/freezer/tokenomics/detailed_coin_metrics"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
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
const (
	ArbitrumBlockchainNetworkType BlockchainNetworkType = "arbitrum"
	BNBBlockchainNetworkType      BlockchainNetworkType = "bnb"
	EthereumBlockchainNetworkType BlockchainNetworkType = "ethereum"
)

var (
	ErrInvalidMiningBoostUpgradeTX                     = errors.New("transaction for upgrading mining boost tier is invalid")
	ErrNotFound                                        = errors.New("not found")
	ErrRelationNotFound                                = errors.New("relationship not found")
	ErrDuplicate                                       = errors.New("duplicate")
	ErrNegativeMiningProgressDecisionRequired          = errors.New("you have negative mining progress, please decide what to do with it")
	ErrKYCRequired                                     = errors.New("user needs to complete one or more kyc steps or skip any of them(if allowed)")
	ErrMiningDisabled                                  = errors.New("mining is disabled")
	ErrRaceCondition                                   = errors.New("race condition")
	ErrGlobalRankHidden                                = errors.New("global rank is hidden")
	ErrDecreasingPreStakingAllocationOrYearsNotAllowed = errors.New("decreasing pre-staking allocation or years not allowed")
	PreStakingBonusesPerYear                           = map[uint8]float64{
		0: 0,
		1: 35,
		2: 70,
		3: 115,
		4: 170,
		5: 250,
	}
	PreStakingYearsByPreStakingBonuses = map[float64]uint8{
		0:   0,
		35:  1,
		70:  2,
		115: 3,
		170: 4,
		250: 5,
	}
)

type (
	BlockchainNetworkType     string
	PendingMiningBoostUpgrade struct {
		ExpiresAt      *time.Time `json:"expiresAt" example:"2022-01-03T16:20:52.156534Z"`
		ICEPrice       string     `json:"icePrice" example:"1234.1234"`
		PaymentAddress string     `json:"paymentAddress" example:"UQBLoASuwnSQVdsw4vZzkhCsN3bruqh68trjf03kHoooMc2k"`
	}
	MiningBoostLevel struct {
		ICEPrice                   string  `json:"icePrice" example:"1234.1234" mapstructure:"-"`
		icePrice                   float64 `json:"-" example:"1234.1234" mapstructure:"-"`
		MaxT1Referrals             uint64  `json:"maxT1Referrals" example:"5" mapstructure:"maxT1Referrals"`
		MiningSessionLengthSeconds uint32  `json:"miningSessionLengthSeconds" example:"86400" mapstructure:"miningSessionLengthSeconds"`
		MiningRateBonus            uint16  `json:"miningRateBonus" example:"100" mapstructure:"miningRateBonus"`
		SlashingDisabled           bool    `json:"slashingDisabled" example:"false" mapstructure:"slashingDisabled"`
	}
	MiningBoostSummary struct {
		CurrentLevelIndex *uint8              `json:"currentLevelIndex,omitempty" example:"0"`
		Levels            []*MiningBoostLevel `json:"levels"`
	}
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
		Total                              DENOM  `json:"total,omitempty" swaggertype:"string" example:"1,243.02"`
		BaseFactor                         DENOM  `json:"baseFactor,omitempty" swaggerignore:"true" swaggertype:"string" example:"1,243.02"`
		Standard                           DENOM  `json:"standard,omitempty" swaggertype:"string" example:"1,243.02"`
		PreStaking                         DENOM  `json:"preStaking,omitempty" swaggertype:"string" example:"1,243.02"`
		TotalNoPreStakingBonus             DENOM  `json:"totalNoPreStakingBonus,omitempty" swaggertype:"string" example:"1,243.02"`
		T1                                 DENOM  `json:"t1,omitempty" swaggertype:"string" example:"1,243.02"`
		T2                                 DENOM  `json:"t2,omitempty" swaggertype:"string" example:"1,243.02"`
		TotalReferrals                     DENOM  `json:"totalReferrals,omitempty" swaggertype:"string" example:"1,243.02"`
		TotalMiningBlockchain              DENOM  `json:"totalMiningBlockchain,omitempty" swaggertype:"string" example:"1,243.02"`
		TotalMainnetRewardPoolContribution DENOM  `json:"totalMainnetRewardPoolContribution,omitempty" swaggertype:"string" example:"1,243.02"`
		UserID                             string `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		miningBlockchainAccountAddress     string
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
	TotalCoins struct {
		Total      float64 `json:"total" example:"111111.2423"`
		Blockchain float64 `json:"blockchain" example:"111111.2423"`
		Standard   float64 `json:"standard" example:"111111.2423"`
		PreStaking float64 `json:"preStaking" example:"111111.2423"`
	}
	TotalCoinsTimeSeriesDataPoint struct {
		Date stdlibtime.Time `json:"date" example:"2022-01-03T16:20:52.156534Z"`
		TotalCoins
	}
	BlockchainDetails struct {
		Timestamp *time.Time `json:"-" redis:"timestamp"`
		MarketCap float64    `json:"marketCap" example:"111111.2423" redis:"market_cap"`
		detailedCoinMetrics.Details
	}
	TotalCoinsSummary struct {
		BlockchainDetails *BlockchainDetails               `json:"blockchainDetails"`
		TimeSeries        []*TotalCoinsTimeSeriesDataPoint `json:"timeSeries"`
		TotalCoins
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
		Bonus float64 `json:"bonus" example:"100.00"`
	}
	PreStaking struct {
		UserID     string  `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Years      uint64  `json:"years" example:"1"`
		Allocation float64 `json:"allocation" example:"100.00"`
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
		MiningStreak                uint64        `json:"miningStreak,omitempty"  example:"2"`
		RemainingFreeMiningSessions uint64        `json:"remainingFreeMiningSessions,omitempty" example:"1"`
		KYCStepBlocked              users.KYCStep `json:"kycStepBlocked,omitempty" example:"2"`
		MiningStarted               bool          `json:"miningStarted,omitempty" example:"true"`
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
		MiningBoostLevel              uint8               `json:"miningBoostLevel,omitempty" swaggerignore:"true" example:"1"`
	}
	ExtraBonusSummary struct {
		UserID              string  `json:"userId,omitempty" swaggerignore:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		AvailableExtraBonus float64 `json:"availableExtraBonus,omitempty" example:"2.00"`
	}
	RankingSummary struct {
		GlobalRank uint64 `json:"globalRank" example:"12333"`
	}
	IceStats struct {
		CirculatingSupply     float64 `json:"circulatingSupply"`
		TotalSupply           float64 `json:"totalSupply"`
		Price                 float64 `json:"price"`
		MarketCap             float64 `json:"marketCap"`
		TradingVolume24       float64 `json:"24hTradingVolume"`
		FullyDilutedMarketCap float64 `json:"fullyDilutedMarketCap"`
	}
	ReadRepository interface {
		GetMiningBoostSummary(ctx context.Context, userID string) (*MiningBoostSummary, error)
		GetBalanceSummary(ctx context.Context, userID string) (*BalanceSummary, error)
		GetTotalCoinsSummary(ctx context.Context, days uint64, utcOffset stdlibtime.Duration) (*TotalCoinsSummary, error)
		GetRankingSummary(ctx context.Context, userID string) (*RankingSummary, error)
		GetTopMiners(ctx context.Context, keyword string, limit, offset uint64) (topMiners []*Miner, nextOffset uint64, err error)
		GetMiningSummary(ctx context.Context, userID string) (*MiningSummary, error)
		GetPreStakingSummary(ctx context.Context, userID string) (*PreStakingSummary, error)
		GetBalanceHistory(ctx context.Context, userID string, start, end *time.Time, utcOffset stdlibtime.Duration, limit, offset uint64) ([]*BalanceHistoryEntry, error) //nolint:lll // .
		GetAdoptionSummary(ctx context.Context, userID string) (*AdoptionSummary, error)
	}
	WriteRepository interface {
		StartNewMiningSession(ctx context.Context, ms *MiningSummary, rollbackNegativeMiningProgress *bool, skipKYCSteps []users.KYCStep) error
		ClaimExtraBonus(ctx context.Context, ebs *ExtraBonusSummary) error
		StartOrUpdatePreStaking(context.Context, *PreStakingSummary) error

		InitializeMiningBoostUpgrade(ctx context.Context, miningBoostLevelIndex uint8, userID string) (*PendingMiningBoostUpgrade, error)
		FinalizeMiningBoostUpgrade(ctx context.Context, network BlockchainNetworkType, txHash, userID string) (*PendingMiningBoostUpgrade, error)
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
	clientTypeCtxValueKey               = "clientTypeCtxValueKey"
	userHashCodeCtxValueKey             = "userHashCodeCtxValueKey"
	authorizationCtxValueKey            = "authorizationCtxValueKey"
	xAccountMetadataCtxValueKey         = "xAccountMetadataCtxValueKey"
	requestDeadline                     = 25 * stdlibtime.Second

	floatToStringFormatter = "%.2f"

	daysCountToInitCoinsCacheOnStartup     = 90
	routinesCountToInitCoinsCacheOnStartup = 1
	totalCoinStatsCacheLockKey             = "totalCoinStatsCache"
	totalCoinStatsCacheLockDuration        = 1 * stdlibtime.Minute

	totalCoinStatsDetailsLockKey      = "totalCoinStatsDetails"
	totalCoinStatsDetailsLockDuration = 1 * stdlibtime.Minute
	totalCoinStatsDetailsKey          = "totalCoinStatsDetailsData"
	miningBoostPricePrecision         = 4 // 4 digits after floating point.

	tokeroTenant = "tokero"
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
		cfg                               *Config
		extraBonusStartDate               *time.Time
		livenessLoadDistributionStartDate *time.Time
		extraBonusIndicesDistribution     map[uint16]map[uint16]uint16
		shutdown                          func() error
		db                                storage.DB
		globalDB                          *storagev2.DB
		dwh                               dwh.Client
		mb                                messagebroker.Client
		pictureClient                     picture.Client
	}

	processor struct {
		*repository
	}

	kycConfigJSON struct {
		DynamicDistributionSocialKYC []*struct {
			DisabledVersions   []string      `json:"disabledVersions"`
			ForceKYCForUserIds []string      `json:"forceKYCForUserIds"`
			KYCStep            users.KYCStep `json:"step"`
			EnabledMobile      bool          `json:"enabledMobile"`
			EnabledWeb         bool          `json:"enabledWeb"`
		} `json:"dynamic-distribution-kyc"`
		FaceAuth struct {
			DisabledVersions   []string `json:"disabledVersions"`
			ForceKYCForUserIds []string `json:"forceKYCForUserIds"`
			Enabled            bool     `json:"enabled"`
		} `json:"face-auth"`
		Social1KYC struct {
			DisabledVersions   []string   `json:"disabledVersions"`
			ForceKYCForUserIds []string   `json:"forceKYCForUserIds"`
			EnabledMobile      bool       `json:"enabledMobile"`
			EnabledWeb         bool       `json:"enabledWeb"`
			StartDate          *time.Time `json:"startDate"`
			Duration           string     `json:"duration"`
			KYCDuration        stdlibtime.Duration
		} `json:"social1-kyc"`
		QuizKYC struct {
			DisabledVersions   []string `json:"disabledVersions"`
			ForceKYCForUserIds []string `json:"forceKYCForUserIds"`
			Enabled            bool     `json:"enabled"`
		} `json:"quiz-kyc"`
		Social2KYC struct {
			DisabledVersions   []string   `json:"disabledVersions"`
			ForceKYCForUserIds []string   `json:"forceKYCForUserIds"`
			EnabledMobile      bool       `json:"enabledMobile"`
			EnabledWeb         bool       `json:"enabledWeb"`
			StartDate          *time.Time `json:"startDate"`
			Duration           string     `json:"duration"`
			KYCDuration        stdlibtime.Duration
		} `json:"social2-kyc"`
		WebFaceAuth struct {
			Enabled bool `json:"enabled"`
		} `json:"web-face-auth"`
		WebQuizKYC struct {
			Enabled bool `json:"enabled"`
		} `json:"web-quiz-kyc"`
	}

	blockchainCoinStatsJSON struct {
		CoinsAddedHistory []*struct {
			Date       *time.Time `json:"date"`
			CoinsAdded float64    `json:"coinsAdded"`
		} `json:"coinsAddedHistory"`
	}
	Config struct {
		disableAdvancedTeam     *atomic.Pointer[[]string]
		kycConfigJSON           *atomic.Pointer[kycConfigJSON]
		blockchainCoinStatsJSON *atomic.Pointer[blockchainCoinStatsJSON]
		MiningBoost             struct {
			icePrice                      *atomic.Pointer[float64]                      `yaml:"-" mapstructure:"-" json:"-"`
			levels                        *atomic.Pointer[[]*MiningBoostLevel]          `yaml:"-" mapstructure:"-" json:"-"`
			networkEndpointCurrentLBIndex map[BlockchainNetworkType]*atomic.Uint64      `yaml:"-" mapstructure:"-" json:"-"`
			networkClients                map[BlockchainNetworkType][]*ethclient.Client `yaml:"-" mapstructure:"-" json:"-"`
			NetworkEndpoints              map[BlockchainNetworkType][]string            `yaml:"networkEndpoints" mapstructure:"networkEndpoints"`
			ContractAddresses             map[BlockchainNetworkType]string              `yaml:"contractAddresses" mapstructure:"contractAddresses"`
			Levels                        map[float64]*MiningBoostLevel                 `yaml:"levels" mapstructure:"levels"`
			SessionLength                 stdlibtime.Duration                           `yaml:"sessionLength" mapstructure:"sessionLength"`
			PriceDelta                    uint8                                         `yaml:"priceDelta" mapstructure:"priceDelta"`
		} `yaml:"mining-boost" mapstructure:"mining-boost"`
		BlockchainCoinStatsJSONURL          string `yaml:"blockchain-coin-stats-json-url" mapstructure:"blockchain-coin-stats-json-url"`
		extrabonusnotifier.ExtraBonusConfig `mapstructure:",squash"`
		Adoption                            struct {
			StartingBaseMiningRate    float64             `yaml:"startingBaseMiningRate" mapstructure:"startingBaseMiningRate"`
			DurationBetweenMilestones stdlibtime.Duration `yaml:"durationBetweenMilestones" mapstructure:"durationBetweenMilestones"`
			Milestones                uint8               `yaml:"milestones" mapstructure:"milestones"`
		} `yaml:"adoption"`
		messagebroker.Config `mapstructure:",squash"`
		KYC                  struct {
			RequireQuizOnlyOnSpecificDayOfWeek *int                `yaml:"require-quiz-only-on-specific-day-of-week" mapstructure:"require-quiz-only-on-specific-day-of-week"` //nolint:lll // .
			TryResetKYCStepsURL                string              `yaml:"try-reset-kyc-steps-url" mapstructure:"try-reset-kyc-steps-url"`
			FaceAuthAvailabilityURL            string              `yaml:"face-auth-availability-url" mapstructure:"face-auth-availability-url"`
			ConfigJSONURL                      string              `yaml:"config-json-url" mapstructure:"config-json-url"`
			FaceRecognitionDelay               stdlibtime.Duration `yaml:"face-recognition-delay" mapstructure:"face-recognition-delay"`
			Social1Delay                       stdlibtime.Duration `yaml:"social1-delay" mapstructure:"social1-delay"`
			Social2Delay                       stdlibtime.Duration `yaml:"social2-delay" mapstructure:"social2-delay"`
			DynamicSocialDelay                 stdlibtime.Duration `yaml:"dynamic-social-delay" mapstructure:"dynamic-social-delay"`
			QuizDelay                          stdlibtime.Duration `yaml:"quiz-delay" mapstructure:"quiz-delay"`
		} `yaml:"kyc" mapstructure:"kyc"`
		MiningSessionDuration struct {
			Min                      stdlibtime.Duration `yaml:"min"`
			Max                      stdlibtime.Duration `yaml:"max"`
			WarnAboutExpirationAfter stdlibtime.Duration `yaml:"warnAboutExpirationAfter"`
		} `yaml:"miningSessionDuration"`
		RollbackNegativeMining struct {
			Available struct {
				After stdlibtime.Duration `yaml:"after"`
				Until stdlibtime.Duration `yaml:"until"`
			} `yaml:"available"`
		} `yaml:"rollbackNegativeMining"`
		ConsecutiveNaturalMiningSessionsRequiredFor1ExtraFreeArtificialMiningSession struct {
			Min uint64 `yaml:"min"`
			Max uint64 `yaml:"max"`
		} `yaml:"consecutiveNaturalMiningSessionsRequiredFor1ExtraFreeArtificialMiningSession"`
		GlobalAggregationInterval struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		} `yaml:"globalAggregationInterval"`
		DetailedCoinMetrics struct {
			RefreshInterval stdlibtime.Duration `yaml:"refresh-interval" mapstructure:"refresh-interval"`
		} `yaml:"detailed-coin-metrics" mapstructure:"detailed-coin-metrics"`
		ReferralBonusMiningRates struct {
			T0 uint16 `yaml:"t0"`
			T1 uint32 `yaml:"t1"`
			T2 uint32 `yaml:"t2"`
		} `yaml:"referralBonusMiningRates"`
		AdminUsers                                   []string `yaml:"adminUsers" mapstructure:"adminUsers"`
		Tenant                                       string   `yaml:"tenant"`
		DefaultReferralName                          string   `yaml:"defaultReferralName"`
		SlashingFloor                                float64  `yaml:"slashingFloor" mapstructure:"slashingFloor"`
		WelcomeBonusV2Amount                         float64  `yaml:"welcomeBonusV2Amount" mapstructure:"welcomeBonusV2Amount"`
		T1ReferralsAllowedWithoutAnyMiningBoostLevel bool     `yaml:"t1ReferralsAllowedWithoutAnyMiningBoostLevel" mapstructure:"t1ReferralsAllowedWithoutAnyMiningBoostLevel"`
		TasksV2Enabled                               bool     `yaml:"tasksV2Enabled" mapstructure:"tasksV2Enabled"`
	}
)

var (
	//go:embed globalDDL.sql
	globalDDL string
)
