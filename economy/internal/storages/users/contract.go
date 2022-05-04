package users

import (
	"time"

	"github.com/framey-io/go-tarantool"
)

// Private API.

type (
	UserID = string

	userEconomy struct {
		_msgpack            struct{} `msgpack:",asArray"`
		UserID              UserID
		ProfilePictureURL   string
		Balance             float64
		StakingPercentage   float64
		HashCode            uint64
		LastMiningStartedAt uint64
		StakingYears        uint64
		CreatedAt           uint64
		UpdatedAt           uint64
		BalanceUpdatedAt    uint64
	}

	// | usersSource is responsible for processing new messages of sourceUser type, transforming it and storing it in the db as user type.
	usersSource struct {
		db tarantool.Connector
	}

	referredBy struct {
		_msgpack struct{} `msgpack:",asArray"`
		UserID   UserID
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
