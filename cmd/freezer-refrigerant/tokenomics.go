// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/server"
	"github.com/ice-blockchain/wintr/terror"
)

func (s *service) setupTokenomicsRoutes(router *server.Router) {
	router.
		Group("/v1w").
		PUT("/tokenomics/:userId/mining-boosts", server.RootHandler(s.InitializeMiningBoostUpgrade)).
		PATCH("/tokenomics/:userId/mining-boosts", server.RootHandler(s.FinalizeMiningBoostUpgrade)).
		POST("/tokenomics/:userId/mining-sessions", server.RootHandler(s.StartNewMiningSession)).
		POST("/tokenomics/:userId/extra-bonus-claims", server.RootHandler(s.ClaimExtraBonus)).
		PUT("/tokenomics/:userId/pre-staking", server.RootHandler(s.StartOrUpdatePreStaking))
}

// InitializeMiningBoostUpgrade godoc
//
//	@Schemes
//	@Description	Initializes the process to enable a new mining boost tier.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string									true	"ID of the user"
//	@Param			x_client_type	query		string									false	"the type of the client calling this API. I.E. `web`"
//	@Param			request			body		InitializeMiningBoostUpgradeRequestBody	true	"Request params"
//	@Success		200				{object}	tokenomics.PendingMiningBoostUpgrade
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		404				{object}	server.ErrorResponse	"if user not found"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1w/tokenomics/{userId}/mining-boosts [PUT].
func (s *service) InitializeMiningBoostUpgrade( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[InitializeMiningBoostUpgradeRequestBody, tokenomics.PendingMiningBoostUpgrade],
) (*server.Response[tokenomics.PendingMiningBoostUpgrade], *server.Response[server.ErrorResponse]) {
	if cfg.Tenant == tokeroTenant {
		return nil, server.Forbidden(errMiningBoostDisabled)
	}
	resp, err := s.tokenomicsProcessor.InitializeMiningBoostUpgrade(ctx, *req.Data.MiningBoostLevelIndex, req.Data.UserID)
	if err = errors.Wrapf(err, "failed to InitializeMiningBoostUpgrade for data:%#v", req.Data); err != nil {
		switch {
		case errors.Is(err, tokenomics.ErrRelationNotFound):
			return nil, server.NotFound(err, userNotFoundErrorCode)
		default:
			return nil, server.Unexpected(err)
		}
	}

	return server.OK(resp), nil
}

// FinalizeMiningBoostUpgrade godoc
//
//	@Schemes
//	@Description	Finalizes the process to enable a new mining boost tier.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string									true	"ID of the user"
//	@Param			x_client_type	query		string									false	"the type of the client calling this API. I.E. `web`"
//	@Param			request			body		FinalizeMiningBoostUpgradeRequestBody	true	"Request params"
//	@Success		200				{object}	tokenomics.PendingMiningBoostUpgrade
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		404				{object}	server.ErrorResponse	"if user not found or process was not initialized"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1w/tokenomics/{userId}/mining-boosts [PATCH].
func (s *service) FinalizeMiningBoostUpgrade( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[FinalizeMiningBoostUpgradeRequestBody, tokenomics.PendingMiningBoostUpgrade],
) (*server.Response[tokenomics.PendingMiningBoostUpgrade], *server.Response[server.ErrorResponse]) {
	if cfg.Tenant == tokeroTenant {
		return nil, server.Forbidden(errMiningBoostDisabled)
	}
	resp, err := s.tokenomicsProcessor.FinalizeMiningBoostUpgrade(ctx, req.Data.Network, req.Data.TXHash, req.Data.UserID)
	if err = errors.Wrapf(err, "failed to FinalizeMiningBoostUpgrade for data:%#v", req.Data); err != nil {
		switch {
		case errors.Is(err, tokenomics.ErrInvalidMiningBoostUpgradeTX):
			return nil, server.BadRequest(err, invalidMiningBoostUpgradeTransactionErrorCode)
		case errors.Is(err, tokenomics.ErrNotFound):
			return nil, server.NotFound(err, noPendingMiningBoostUpgradeFoundErrorCode)
		case errors.Is(err, tokenomics.ErrRelationNotFound):
			return nil, server.NotFound(err, userNotFoundErrorCode)
		case errors.Is(err, tokenomics.ErrDuplicate):
			return nil, server.Conflict(err, transactionAlreadyUsed)
		default:
			return nil, server.Unexpected(err)
		}
	}

	return server.OK(resp), nil
}

// StartNewMiningSession godoc
//
//	@Schemes
//	@Description	Starts a new mining session for the user, if not already in progress with another one.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string								true	"ID of the user"
//	@Param			x_client_type	query		string								false	"the type of the client calling this API. I.E. `web`"
//	@Param			request			body		StartNewMiningSessionRequestBody	true	"Request params"
//	@Success		201				{object}	tokenomics.MiningSummary
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		404				{object}	server.ErrorResponse	"if user not found"
//	@Failure		409				{object}	server.ErrorResponse	"if mining is in progress or if a decision about negative mining progress or kyc is required"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1w/tokenomics/{userId}/mining-sessions [POST].
func (s *service) StartNewMiningSession( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[StartNewMiningSessionRequestBody, tokenomics.MiningSummary],
) (*server.Response[tokenomics.MiningSummary], *server.Response[server.ErrorResponse]) {
	if cfg.Tenant == tokeroTenant {
		return nil, server.Forbidden(errMiningDisabled)
	}
	ms := &tokenomics.MiningSummary{MiningSession: &tokenomics.MiningSession{UserID: &req.Data.UserID}}
	ctx = contextWithHashCode(ctx, req)
	ctx = tokenomics.ContextWithClientType(ctx, req.Data.XClientType)
	ctx = tokenomics.ContextWithAuthorization(ctx, req.Data.Authorization)
	ctx = tokenomics.ContextWithXAccountMetadata(ctx, req.Data.XAccountMetadata)
	if err := s.tokenomicsProcessor.StartNewMiningSession(ctx, ms, req.Data.Resurrect, req.Data.SkipKYCSteps); err != nil {
		err = errors.Wrapf(err, "failed to start a new mining session for userID:%v, data:%#v", req.Data.UserID, req.Data)
		switch {
		case errors.Is(err, tokenomics.ErrNegativeMiningProgressDecisionRequired):
			if tErr := terror.As(err); tErr != nil {
				return nil, server.Conflict(err, resurrectionDecisionRequiredErrorCode, tErr.Data)
			}

			fallthrough
		case errors.Is(err, tokenomics.ErrKYCRequired):
			if tErr := terror.As(err); tErr != nil {
				return nil, server.Conflict(err, kycStepsRequiredErrorCode, tErr.Data)
			}

			fallthrough
		case errors.Is(err, tokenomics.ErrMiningDisabled):
			if tErr := terror.As(err); tErr != nil {
				return nil, server.ForbiddenWithCode(err, miningDisabledErrorCode, tErr.Data)
			}

			fallthrough
		case errors.Is(err, tokenomics.ErrRaceCondition):
			return nil, server.BadRequest(err, raceConditionErrorCode)
		case errors.Is(err, tokenomics.ErrDuplicate):
			return nil, server.Conflict(err, miningInProgressErrorCode)
		case errors.Is(err, tokenomics.ErrRelationNotFound):
			return nil, server.NotFound(err, userNotFoundErrorCode)
		}

		return nil, server.Unexpected(err)
	}

	return server.Created(ms), nil
}

// ClaimExtraBonus godoc
//
//	@Schemes
//	@Description	Claims an extra bonus for the user.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string	true	"ID of the user"
//	@Success		201				{object}	tokenomics.ExtraBonusSummary
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		404				{object}	server.ErrorResponse	"if user not found or no extra bonus available"
//	@Failure		409				{object}	server.ErrorResponse	"if already claimed"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1w/tokenomics/{userId}/extra-bonus-claims [POST].
func (s *service) ClaimExtraBonus( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[ClaimExtraBonusRequestBody, tokenomics.ExtraBonusSummary],
) (*server.Response[tokenomics.ExtraBonusSummary], *server.Response[server.ErrorResponse]) {
	resp := &tokenomics.ExtraBonusSummary{UserID: req.Data.UserID}
	if true {
		return nil, server.Forbidden(errors.New("disabled"))
	}
	if err := s.tokenomicsProcessor.ClaimExtraBonus(contextWithHashCode(ctx, req), resp); err != nil {
		err = errors.Wrapf(err, "failed to claim extra bonus for userID:%v", req.Data.UserID)
		switch {
		case errors.Is(err, tokenomics.ErrNotFound):
			return nil, server.NotFound(err, noExtraBonusAvailableErrorCode)
		case errors.Is(err, tokenomics.ErrDuplicate):
			return nil, server.Conflict(err, extraBonusAlreadyClaimedErrorCode)
		}

		return nil, server.Unexpected(err)
	}

	return server.Created(resp), nil
}

// StartOrUpdatePreStaking godoc
//
//	@Schemes
//	@Description	Starts or updates pre-staking for the user.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string								true	"ID of the user"
//	@Param			request			body		StartOrUpdatePreStakingRequestBody	true	"Request params"
//	@Success		200				{object}	tokenomics.PreStakingSummary
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		404				{object}	server.ErrorResponse	"user not found"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1w/tokenomics/{userId}/pre-staking [PUT].
func (s *service) StartOrUpdatePreStaking( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[StartOrUpdatePreStakingRequestBody, tokenomics.PreStakingSummary],
) (*server.Response[tokenomics.PreStakingSummary], *server.Response[server.ErrorResponse]) {
	const maxAllocation = 100
	if *req.Data.Years > tokenomics.MaxPreStakingYears {
		*req.Data.Years = tokenomics.MaxPreStakingYears
	}
	if *req.Data.Allocation > maxAllocation {
		*req.Data.Allocation = maxAllocation
	}
	allocation := float64(*req.Data.Allocation)
	st := &tokenomics.PreStakingSummary{
		PreStaking: &tokenomics.PreStaking{
			UserID:     req.Data.UserID,
			Years:      uint64(*req.Data.Years),
			Allocation: allocation,
		},
	}

	if true {
		return nil, server.ForbiddenWithCode(errors.Errorf("Endpoint disabled"), prestakingDisabled)
	}

	if err := s.tokenomicsProcessor.StartOrUpdatePreStaking(contextWithHashCode(ctx, req), st); err != nil {
		err = errors.Wrapf(err, "failed to StartOrUpdatePreStaking for %#v", req.Data)
		if errors.Is(err, tokenomics.ErrRelationNotFound) {
			return nil, server.NotFound(err, userNotFoundErrorCode)
		}

		return nil, server.Unexpected(err)
	}

	return server.OK(st), nil
}
