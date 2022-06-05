// SPDX-License-Identifier: BUSL-1.1

package usereconomy

import (
	"github.com/framey-io/go-tarantool"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/time"
)

// Private API.

const (
	tierLevel0 uint64 = 0
	tierLevel1        = 1
	tierLevel2        = 2

	balanceTypeStandard string = "standard"
	balanceTypeStaking         = "staking"
	balanceTypeTotal           = "total"
)

type (
	UserID      = string
	BalanceType = string
	TierLevel   = uint64

	// | userEconomy is the internal structure for deserialization from the DB.
	userEconomy struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack            struct{} `msgpack:",asArray"`
		LastMiningStartedAt *time.Time
		CreatedAt           *time.Time
		UpdatedAt           *time.Time
		UserID              UserID
		Username            string
		ProfilePictureURL   string
		HashCode            uint64
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

	// | userEconomySource is responsible for processing new messages of sourceUser type, transforming it and storing it in the db as user type.
	userEconomySource struct {
		db tarantool.Connector
	}

	tier struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack struct{} `msgpack:",asArray"`
		UserID   UserID
	}

	totalUsers struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack struct{} `msgpack:",asArray"`
		Key      string
		Value    uint64
	}

	balances struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack  struct{} `msgpack:",asArray"`
		UpdatedAt *time.Time
		Amount    *coin.ICEFlake
		UserID    UserID
		Type      string
		AmountW0  uint64
		AmountW1  uint64
		AmountW2  uint64
		AmountW3  uint64
	}
)
