// SPDX-License-Identifier: ice License 1.0

package main

import (
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	coindistribution "github.com/ice-blockchain/freezer/coin-distribution"
	"github.com/ice-blockchain/freezer/tokenomics"
)

// Public API.

type (
	StartNewMiningSessionRequestBody struct {
		// Specify this if you want to resurrect the user.
		// `true` recovers all the lost balance, `false` deletes it forever, `null/undefined` does nothing. Default is `null/undefined`.
		Resurrect        *bool  `json:"resurrect" example:"true"`
		UserID           string `uri:"userId" swaggerignore:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		XClientType      string `form:"x_client_type" swaggerignore:"true" required:"false" example:"web"`
		Authorization    string `header:"Authorization" swaggerignore:"true" required:"true" example:"some token"`
		XAccountMetadata string `header:"X-Account-Metadata" swaggerignore:"true" required:"false" example:"some token"`
		// Specify this if you want to skip one or more specific KYC steps before starting a new mining session or extending an existing one.
		// Some KYC steps are not skippable.
		SkipKYCSteps []users.KYCStep `json:"skipKYCSteps" example:"0,1"`
	}
	ClaimExtraBonusRequestBody struct {
		UserID string `uri:"userId" swaggerignore:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	StartOrUpdatePreStakingRequestBody struct {
		Years      *uint8 `json:"years" required:"true" maximum:"5" example:"1"`
		Allocation *uint8 `json:"allocation" required:"true" maximum:"100" example:"100"`
		UserID     string `uri:"userId" swaggerignore:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	InitializeMiningBoostUpgradeRequestBody struct {
		MiningBoostLevelIndex *uint8 `json:"miningBoostLevelIndex" required:"true" example:"0"`
		UserID                string `uri:"userId" swaggerignore:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	FinalizeMiningBoostUpgradeRequestBody struct {
		UserID  string                           `uri:"userId" swaggerignore:"true" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		Network tokenomics.BlockchainNetworkType `json:"network" required:"true" example:"ethereum" enums:"arbitrum,bnb,ethereum"`
		TXHash  string                           `json:"txHash" required:"true" example:"0xf75c78ab01ee4641be46794756f46137dea03a4980126dce4f2df933cccb34ea"`
	}
)

// Private API.

const (
	applicationYamlKey = "cmd/freezer-refrigerant"
	swaggerRootSuffix  = "/tokenomics/w"

	adminRole = "admin"
)

// Values for server.ErrorResponse#Code.
const (
	userNotFoundErrorCode                         = "USER_NOT_FOUND"
	prestakingDisabled                            = "PRESTAKING_DISABLED"
	miningInProgressErrorCode                     = "MINING_IN_PROGRESS"
	raceConditionErrorCode                        = "RACE_CONDITION"
	resurrectionDecisionRequiredErrorCode         = "RESURRECTION_DECISION_REQUIRED"
	kycStepsRequiredErrorCode                     = "KYC_STEPS_REQUIRED"
	miningDisabledErrorCode                       = "MINING_DISABLED"
	noExtraBonusAvailableErrorCode                = "NO_EXTRA_BONUS_AVAILABLE"
	extraBonusAlreadyClaimedErrorCode             = "EXTRA_BONUS_ALREADY_CLAIMED"
	noPendingMiningBoostUpgradeFoundErrorCode     = "NO_PENDING_MINING_BOOST_UPGRADE_FOUND"
	invalidMiningBoostUpgradeTransactionErrorCode = "INVALID_MINING_BOOST_UPGRADE_TRANSACTION"
	transactionAlreadyUsed                        = "TRANSACTION_ALREADY_USED"

	defaultDistributionLimit = 5000

	tokeroTenant = "tokero"
)

// .
var (
	//nolint:gochecknoglobals // Because its loaded once, at runtime.
	cfg                    config
	errMiningDisabled      = errors.New("mining disabled")
	errMiningBoostDisabled = errors.New("mining boost disabled")
)

type (
	// | service implements server.State and is responsible for managing the state and lifecycle of the package.
	service struct {
		tokenomicsProcessor        tokenomics.Processor
		coinDistributionRepository coindistribution.Repository
	}
	config struct {
		Host    string `yaml:"host"`
		Version string `yaml:"version"`
		Tenant  string `yaml:"tenant"`
	}
)
