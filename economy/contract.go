// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	_ "embed"
	"io"
	"time"

	"github.com/framey-io/go-tarantool"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
)

// Public API.

var (
	ErrNotFound         = storage.ErrNotFound
	ErrMiningInProgress = errors.New("mining in progress")
)

type (
	UserID               = string
	TotalUsers           = uint64
	BaseHourlyMiningRate = float64
	UserEconomy          struct {
		LastMiningStartedAt time.Time                           `json:"lastMiningStartedAt" example:"2022-01-03T16:20:52.156534Z"`
		Adoption            map[TotalUsers]BaseHourlyMiningRate `json:"adoption"`
		Balance             Balance                             `json:"balance"`
		CurrentTotalUsers   TotalUsers                          `json:"currentTotalUsers" example:"1000000"`
		Staking             Staking                             `json:"staking"`
		HourlyMiningRate    float64                             `json:"hourlyMiningRate" example:"232.5"`
		GlobalRank          uint64                              `json:"globalRank" example:"1000"`
	}
	Staking struct {
		Years      uint64  `json:"years" example:"1"`
		Percentage float64 `json:"percentage" example:"25.0"`
	}
	Balance struct {
		Total     float64         `json:"total" example:"232.5"`
		Referrals ReferralBalance `json:"referrals"`
	}
	ReferralBalance struct {
		T1 float64 `json:"t1" example:"232.5"`
		T2 float64 `json:"t2" example:"232.5"`
	}
	TopMiner struct {
		UserID            string  `json:"userId" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Username          string  `json:"username" example:"jdoe"`
		ProfilePictureURL string  `json:"profilePictureURL" example:"https://somecdn.com/p1.jpg"`
		Balance           float64 `json:"balance" example:"232.5"`
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
	}
	// WriteRepository manage the database operations related to `user_economy`.
	WriteRepository interface {
		StartMining(context.Context, UserID) error
	}
)

const (
	applicationYamlKey = "economy"
	base10             = 10
	bitSize64          = 64
	miningDuration     = 24 * time.Hour
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
		_msgpack            struct{} `msgpack:",asArray"`
		UserID              string
		ProfilePictureURL   string
		Adoptions           string
		Balance             float64
		StakingPercentage   float64
		HashCode            uint64
		LastMiningStartedAt uint64
		StakingYears        uint64
		CreatedAt           uint64
		UpdatedAt           uint64
		BalanceUpdatedAt    uint64
		T1Count             uint64
		T2Count             uint64
		GlobalRank          uint64
		T1EarningsSum       float64
		T2EarningsSum       float64
		CurrentTotalUsers   uint64
	}

	// | userEconomyLastMining is the internal structure for deserialization from the DB.
	userEconomyLastMining struct {
		_msgpack            struct{} `msgpack:",asArray"`
		LastMiningStartedAt uint64
	}

	// | miningStarted is internal structure to hold notification message.
	miningStarted struct {
		TS time.Time `json:"ts"`
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
	}
	economy struct {
		db tarantool.Connector
		mb messagebroker.Client
	}
	// | config holds the configuration of this package mounted from `application.yaml`.
	config struct {
		MessageBroker struct {
			Topics []struct {
				Name string `yaml:"name" json:"name"`
			} `yaml:"topics"`
		} `yaml:"messageBroker"`
		Rates struct {
			Tier1 float64
			Tier2 float64
		} `yaml:"rates"`
		InactivityHoursDeadline uint64 `yaml:"inactivityHoursDeadline"`
	}
)
