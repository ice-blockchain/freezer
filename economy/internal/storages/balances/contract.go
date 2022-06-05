package balances

import (
	"github.com/framey-io/go-tarantool"

	"github.com/ice-blockchain/wintr/coin"
)

// Private API.

type (
	UserID = string

	// | stakingInfo is the internal structure to parse balances.
	stakingInfo struct {
		Bonus      uint64
		Allocation uint64
	}

	// | balanceSource is responsible for processing messages from balance update topic, updating balance information at database.
	balanceSource struct {
		db  tarantool.Connector
		cfg *config
	}

	// | userEconomy is the internal structure for deserialization from the DB.
	userEconomy struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack             struct{} `msgpack:",asArray"`
		UserID               UserID
		BaseHourlyMiningRate *coin.ICEFlake
		StakingInfo          string
		Balances             string
		T0Referrals          uint64
		T1Referrals          uint64
		T2Referrals          uint64
	}

	// | balance is the internal structure to parse balances information.
	balance struct {
		Amount    *coin.ICEFlake
		UpdatedAt uint64
	}

	// | config holds the configuration of this package mounted from `application.yaml`.
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
	applicationYamlKey        = "economy"
	base10                    = 10
	bitSize64                 = 64
	percentage100      uint64 = 100
)
