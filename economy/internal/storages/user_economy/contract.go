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
		BalanceUpdatedAt    *time.Time
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

	userSnapshot struct {
		User   *user
		Before *user
	}

	user struct {
		CreatedAt         time.Time  `json:"createdAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		UpdatedAt         time.Time  `json:"updatedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		DeletedAt         *time.Time `json:"deletedAt,omitempty" example:"2022-01-03T16:20:52.156534Z"`
		ID                string     `json:"id,omitempty" example:"226fcb86-fcce-458e-95f0-867e09c8c274"`
		Email             string     `form:"email,omitempty" json:"email" example:"jdoe@gmail.com"`
		FullName          string     `form:"fullName,omitempty" json:"fullName" example:"John Doe"`
		PhoneNumber       string     `form:"phoneNumber,omitempty" json:"phoneNumber" example:"+12099216581"`
		Username          string     `form:"username,omitempty" json:"username" example:"jdoe"`
		ReferredBy        string     `form:"referredBy,omitempty" json:"referredBy" example:"billy112"`
		ProfilePictureURL string     `json:"profilePictureURL,omitempty" example:"https://somecdn.com/p1.jpg"`
		// ISO 3166 country code.
		Country  string `json:"country" example:"us"`
		HashCode uint64 `json:"hashCode"`
	}
)
