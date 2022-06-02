package balances

import (
	"time"

	"github.com/framey-io/go-tarantool"
	"github.com/ice-blockchain/wintr/coin"
)

// Private API.

type (
	UserID = string

	// TODO: add description.
	stakingInfo struct {
		Bonus      uint64
		Allocation uint64
	}

	// | balanceSource is responsible for processing messages from balance update topic, updating balance information at database.
	balanceSource struct {
		db  tarantool.Connector
		cfg *config
	}

	// TODO: add description.
	userEconomy struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack             struct{} `msgpack:",asArray"`
		UserID               UserID
		LastMiningStartedAt  *time.Time
		ElapsedNanoseconds   uint64
		T0Referrals          uint64
		T1Referrals          uint64
		T2Referrals          uint64
		BaseHourlyMiningRate *coin.ICEFlake
		StakingInfo          string
		Balances             string
	}

	// TODO: add description.
	balance struct {
		Amount    *coin.ICEFlake
		UpdatedAt uint64
	}

	config struct {
		MessageBroker struct {
			Topics []struct {
				Name       string `yaml:"name" json:"name"`
				Partitions uint64 `yaml:"partitions" json:"partitions"`
				Partition  uint64 `yaml:"partition" json:"partition"`
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

const (
	applicationYamlKey = "economy"
	base10             = 10
	bitSize64          = 64
)
