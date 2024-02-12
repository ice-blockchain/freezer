// SPDX-License-Identifier: ice License 1.0

package main

import (
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/tokenomics"
)

// Public API.

type (
	GetMiningSummaryArg struct {
		UserID string `uri:"userId" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	GetPreStakingSummaryArg struct {
		UserID string `uri:"userId" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	GetBalanceSummaryArg struct {
		UserID string `uri:"userId" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	GetBalanceHistoryArg struct {
		// The start date in RFC3339 or ISO8601 formats. Default is `now` in UTC.
		StartDate *stdlibtime.Time `form:"startDate" swaggertype:"string" example:"2022-01-03T16:20:52.156534Z"`
		// The start date in RFC3339 or ISO8601 formats. Default is `end of day, relative to startDate`.
		EndDate *stdlibtime.Time `form:"endDate" swaggertype:"string" example:"2022-01-03T16:20:52.156534Z"`
		UserID  string           `uri:"userId" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		TZ      string           `form:"tz" example:"-03:00"`
		// Default is 24.
		Limit  uint64 `form:"limit" maximum:"1000" example:"24"`
		Offset uint64 `form:"offset" example:"0"`
	}
	GetRankingSummaryArg struct {
		UserID string `uri:"userId" allowForbiddenGet:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	GetTopMinersArg struct {
		Keyword string `form:"keyword" example:"jdoe"`
		// Default is 10.
		Limit  uint64 `form:"limit" maximum:"1000" example:"10"`
		Offset uint64 `form:"offset" example:"0"`
	}
	GetAdoptionArg      struct{}
	GetCoinsDetailedArg struct{}
	GetTotalCoinsArg    struct {
		TZ   string `form:"tz" example:"+4:30" allowUnauthorized:"true"`
		Days uint64 `form:"days" example:"7"`
	}
)

// Private API.

const (
	applicationYamlKey = "cmd/freezer"
	swaggerRoot        = "/tokenomics/r"
)

// Values for server.ErrorResponse#Code.
const (
	userNotFoundErrorCode             = "USER_NOT_FOUND"
	userPreStakingNotEnabledErrorCode = "PRE_STAKING_NOT_ENABLED"
	globalRankHiddenErrorCode         = "GLOBAL_RANK_HIDDEN"
	invalidPropertiesErrorCode        = "INVALID_PROPERTIES"
)

type (
	// | service implements server.State and is responsible for managing the state and lifecycle of the package.
	service struct {
		tokenomicsRepository tokenomics.Repository
	}
	config struct {
		Host    string `yaml:"host"`
		Version string `yaml:"version"`
	}
)
