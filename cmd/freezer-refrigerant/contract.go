// SPDX-License-Identifier: ice License 1.0

package main

import (
	"sync"

	"github.com/ice-blockchain/freezer/tokenomics"
)

// Public API.

type (
	StartNewMiningSessionRequestBody struct {
		// Specify this if you want to resurrect the user.
		// `true` recovers all the lost balance, `false` deletes it forever, `null/undefined` does nothing. Default is `null/undefined`.
		Resurrect *bool  `json:"resurrect" example:"true"`
		UserID    string `uri:"userId" swaggerignore:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	ClaimExtraBonusRequestBody struct {
		UserID string `uri:"userId" swaggerignore:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	StartOrUpdatePreStakingRequestBody struct {
		UserID     string `uri:"userId" swaggerignore:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Years      uint8  `json:"years" required:"true" maximum:"5" example:"1"`
		Allocation uint8  `json:"allocation" required:"true" maximum:"100" example:"100"`
	}
)

// Private API.

const (
	applicationYamlKey = "cmd/freezer-refrigerant"
	swaggerRoot        = "/tokenomics/w"
)

// Values for server.ErrorResponse#Code.
const (
	userNotFoundErrorCode                                    = "USER_NOT_FOUND"
	decreasingPreStakingAllocationOrYearsNotAllowedErrorCode = "DECREASING_PRE_STAKING_ALLOCATION_OR_YEARS_NOT_ALLOWED"
	miningInProgressErrorCode                                = "MINING_IN_PROGRESS"
	raceConditionErrorCode                                   = "RACE_CONDITION"
	resurrectionDecisionRequiredErrorCode                    = "RESURRECTION_DECISION_REQUIRED"
	noExtraBonusAvailableErrorCode                           = "NO_EXTRA_BONUS_AVAILABLE"
	extraBonusAlreadyClaimedErrorCode                        = "EXTRA_BONUS_ALREADY_CLAIMED"
)

type (
	// | service implements server.State and is responsible for managing the state and lifecycle of the package.
	service struct {
		tokenomicsProcessor tokenomics.Processor
		wg                  *sync.WaitGroup
	}
	config struct {
		Host    string `yaml:"host"`
		Version string `yaml:"version"`
	}
)
