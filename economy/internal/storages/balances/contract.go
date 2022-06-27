// SPDX-License-Identifier: BUSL-1.1

package balances

import (
	"time"

	"github.com/framey-io/go-tarantool"

	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
)

// Private API.

//nolint:gochecknoglobals // Because its loaded once, at runtime.
var cfg config

type (
	UserID      = string
	BalanceType = string

	// | balanceSource is responsible for processing messages from balance update topic, updating balance information at database.
	balanceDistributedBatchProcessingStreamSource struct {
		db tarantool.Connector
		mb messagebroker.Client
	}

	// | userEconomy is the internal structure for deserialization from the DB.
	userEconomy struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack             struct{} `msgpack:",asArray"`
		UserID               UserID
		BaseHourlyMiningRate *coin.ICEFlake
		Balances             string
		Bonus                uint64
		Allocation           uint64
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

	generalBalanceDivider         uint64 = 3600000000000
	stakedHourlyMiningRateDivider uint64 = 10000
	t0StandardDivider             uint64 = 1440000000000000
	t1StandardDivider             uint64 = 1440000000000000
	t2StandardDivider             uint64 = 7200000000000000
	t0StakingDivider              uint64 = 144000000000000000
	t1StakingDivider              uint64 = 144000000000000000
	t2StakingDivider              uint64 = 720000000000000000

	userBalanceMessageDeadline = 30 * time.Second
)
