// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	_ "embed"
	"io"
	tm "time"

	"github.com/framey-io/go-tarantool"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

var (
	ErrNotFound              = storage.ErrNotFound
	ErrMiningInProgress      = errors.New("mining in progress")
	ErrStakingAlreadyEnabled = errors.New("staking already enabled")
)

type (
	UserID               = string
	TotalUsers           = uint64
	BaseHourlyMiningRate = *coin.ICEFlake
	UserEconomy          struct {
		LastMiningStartedAt *time.Time                          `json:"lastMiningStartedAt" example:"2022-01-03T16:20:52.156534Z"`
		HourlyMiningRate    *coin.ICEFlake                      `json:"hourlyMiningRate" example:"232"`
		Adoption            map[TotalUsers]BaseHourlyMiningRate `json:"adoption"`
		Balance             Balance                             `json:"balance"`
		CurrentTotalUsers   TotalUsers                          `json:"currentTotalUsers" example:"1000000"`
		Staking             Staking                             `json:"staking"`
		GlobalRank          uint64                              `json:"globalRank" example:"1000"`
	}
	EstimatedEarnings struct {
		StandardHourlyMiningRate *coin.ICEFlake `json:"standardHourlyMiningRate" swaggertype:"string" example:"12.123456789"`
		StakingHourlyMiningRate  *coin.ICEFlake `json:"stakingHourlyMiningRate" swaggertype:"string" example:"12.123456789"`
	}
	Staking struct {
		Years      uint64 `json:"years" example:"1"`
		Percentage uint64 `json:"percentage" example:"200"`
	}
	Balance struct {
		Total     *coin.ICEFlake  `json:"total" example:"232"`
		Referrals ReferralBalance `json:"referrals"`
	}
	ReferralBalance struct {
		T0 *coin.ICEFlake `json:"t0" example:"232"`
		T1 *coin.ICEFlake `json:"t1" example:"232"`
		T2 *coin.ICEFlake `json:"t2" example:"232"`
	}
	GetEstimatedEarningsArg struct {
		T1ActiveReferrals uint64 `form:"t1" example:"20"`
		T2ActiveReferrals uint64 `form:"t2" example:"20"`
		T0ActiveReferee   bool   `form:"t0" example:"true"`
		StakingYears      uint8  `form:"stakingYears" example:"1"`
		StakingAllocation uint8  `form:"stakingAllocation" example:"100"`
	}
	TopMiner struct {
		Balance           *coin.ICEFlake `json:"balance" swaggertype:"string" example:"12.123456789"`
		UserID            string         `json:"userId" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Username          string         `json:"username" example:"jdoe"`
		ProfilePictureURL string         `json:"profilePictureURL" example:"https://somecdn.com/p1.jpg"`
	}
	AdoptionMilestone struct {
		HourlyMiningRate *coin.ICEFlake `json:"hourlyMiningRate" swaggertype:"string" example:"12.123456789"`
		Users            UserCounter    `json:"users"`
		Achieved         bool           `json:"achieved" example:"true"`
	}
	UserCounter struct {
		Total  uint64 `json:"total,omitempty" example:"1000000000"`
		Active uint64 `json:"active,omitempty" example:"1000000000"`
	}
	Adoption struct {
		Adoption []*AdoptionMilestone `json:"adoption"`
		Users    UserCounter          `json:"users"`
	}
	DailyUserGrowth struct {
		Year  int         `json:"year" example:"2022"`
		Month int         `json:"month" example:"12"`
		Day   int         `json:"day" example:"31"`
		Users UserCounter `json:"users"`
	}
	UserStats struct {
		UserGrowth []*DailyUserGrowth `json:"userGrowth"`
		Users      UserCounter        `json:"users"`
	}
	Days            = uint16
	GetTopMinersArg struct {
		Keyword string `form:"keyword" example:"ab"`
		Limit   uint64 `form:"limit" example:"20"`
		Offset  uint64 `form:"offset" example:"0"`
	}
	Repository interface {
		io.Closer
		ReadRepository
	}
	Processor interface {
		Repository
		WriteRepository
		CheckHealth(context.Context) error
	}
	// ReadRepository manages the database operations related to `users_economy`.
	ReadRepository interface {
		GetUserEconomy(context.Context, string, bool) (*UserEconomy, error)
		GetTopMiners(context.Context, *GetTopMinersArg) ([]*TopMiner, error)
		GetEstimatedEarnings(context.Context, *GetEstimatedEarningsArg) (*EstimatedEarnings, error)
		GetAdoption(context.Context) (*Adoption, error)
		GetUserStats(context.Context, Days) (*UserStats, error)
	}
	// WriteRepository manage the database operations related to `user_economy`.
	WriteRepository interface {
		StartMining(context.Context, UserID) error
		StartStaking(context.Context, UserID, Staking) error
	}

	// | MiningStarted is structure to deserialize from the DB and to hold notification message.
	MiningStarted struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack            struct{}   `msgpack:",asArray"`
		LastMiningStartedAt *time.Time `json:"ts"`
	}

	// | StakingEnabled is structure to hold notification message sent to message broker.
	StakingEnabled struct {
		TS *time.Time `json:"ts"`
		Staking
	}
)

const (
	applicationYamlKey                = "economy"
	balanceTypeStaking                = "staking"
	balanceTypeStandard               = "standard"
	base10                            = 10
	bitSize64                         = 64
	miningDuration                    = 24 * tm.Hour
	percentage100                     = 100
	balancesUpdateMillisecondsTicker  = 100 * tm.Millisecond
	sendUpdateBalancesMessageDeadline = 30 * tm.Second
)

var (
	//go:embed DDL.lua
	ddl string
	//nolint:gochecknoglobals // Because its loaded once, at runtime.
	cfg config
)

type (
	// | userEconomy is the internal (UserEconomy) structure for deserialization from the DB
	// because it cannot deserialize time.Time or map/json structures properly.
	// !! Order of fields is crucial, so do not change it !!
	userEconomySummary struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack                    struct{} `msgpack:",asArray"`
		LastMiningStartedAt         *time.Time
		StakingBalanceUpdatedAt     *time.Time
		Balance                     *coin.ICEFlake
		StakingBalance              *coin.ICEFlake
		BaseHourlyMiningRate        *coin.ICEFlake
		T0Amount                    *coin.ICEFlake
		T1Amount                    *coin.ICEFlake
		T2Amount                    *coin.ICEFlake
		UserID                      string
		Username                    string
		ProfilePictureURL           string
		Adoptions                   string
		HashCode                    uint64
		T0Count                     uint64
		T1Count                     uint64
		T2Count                     uint64
		GlobalRank                  uint64
		StakingPercentageAllocation uint64
		StakingYears                uint64
		CurrentTotalUsers           uint64
		StakingPercentageBonus      uint64
	}

	// | stakingAlreadyEnabled is the internal structure for deserialization from the DB.
	stakingAlreadyEnabled struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack struct{} `msgpack:",asArray"`
		Value    bool
	}

	// | userBalance is the internal structure for deserialization from the DB.
	userBalance struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack struct{} `msgpack:",asArray"`
		Balance  *coin.ICEFlake
	}

	// | staking is the internal structure for deserialization from the DB.
	staking struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack   struct{} `msgpack:",asArray"`
		CreatedAt  *time.Time
		UpdatedAt  *time.Time
		UserID     UserID
		Percentage uint64
		Years      uint64
	}

	// | repository implements the public API that this package exposes.
	repository struct {
		close func() error
		ReadRepository
	}
	// | processor implements the processing API that this package exposes.
	processor struct {
		close func() error
		ReadRepository
		WriteRepository
		mb     messagebroker.Client
		ticker *tickerManager
	}
	economy struct {
		db tarantool.Connector
		mb messagebroker.Client
	}
	// | ticker manager allows gracefully close the ticker.
	tickerManager struct {
		mb     messagebroker.Client
		cfg    *config
		closed bool
	}
	// | config holds the configuration of this package mounted from `application.yaml`.
	config struct {
		MessageBroker struct {
			ConsumingTopics []string `yaml:"consumingTopics"`
			Topics          []struct {
				Name string `yaml:"name" json:"name"`
			} `yaml:"topics"`
		} `yaml:"messageBroker"`
		Rates struct {
			Tier0 uint64
			Tier1 uint64
			Tier2 uint64
		} `yaml:"rates"`
		InactivityHoursDeadline uint64 `yaml:"inactivityHoursDeadline"`
	}
)
