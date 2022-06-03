// SPDX-License-Identifier: BUSL-1.1

package usereconomy

import (
	"github.com/framey-io/go-tarantool"
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
	UserID = string

	userEconomy struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack            struct{} `msgpack:",asArray"`
		UserID              UserID
		Username            string
		ProfilePictureURL   string
		HashCode            uint64
		LastMiningStartedAt *time.Time
		CreatedAt           *time.Time
		UpdatedAt           *time.Time
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

	referralEarnings struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack  struct{} `msgpack:",asArray"`
		UserID    UserID
		Type      string
		UpdatedAt *time.Time
	}
)
