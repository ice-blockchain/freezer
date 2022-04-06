// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"io"
	"time"
)

// Public API.

type (
	TotalUsers           = uint64
	BaseHourlyMiningRate = float64
	UserEconomy          struct {
		Balance             Balance                             `json:"balance"`
		HourlyMiningRate    float64                             `json:"hourlyMiningRate" example:"232.5"`
		GlobalRank          uint64                              `json:"globalRank" example:"1000"`
		CurrentTotalUsers   TotalUsers                          `json:"currentTotalUsers" example:"1000000"`
		Adoption            map[TotalUsers]BaseHourlyMiningRate `json:"adoption"`
		LastMiningStartedAt time.Time                           `json:"lastMiningStartedAt" example:"2022-01-03T16:20:52.156534Z"`
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
		ProfilePictureURL string  `json:"profilePictureURL" example:"https://somecdn.com/p1.jpg"`
		Balance           float64 `json:"balance" example:"232.5"`
	}
	Repository interface {
		io.Closer
	}
)

// Private API.

type (
	repository struct{}
)
