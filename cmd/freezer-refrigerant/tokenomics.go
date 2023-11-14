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
		POST("/tokenomics/:userId/mining-sessions", server.RootHandler(s.StartNewMiningSession)).
		POST("/tokenomics/:userId/extra-bonus-claims", server.RootHandler(s.ClaimExtraBonus)).
		PUT("/tokenomics/:userId/pre-staking", server.RootHandler(s.StartOrUpdatePreStaking))
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
//	@Router			/tokenomics/{userId}/mining-sessions [POST].
func (s *service) StartNewMiningSession( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[StartNewMiningSessionRequestBody, tokenomics.MiningSummary],
) (*server.Response[tokenomics.MiningSummary], *server.Response[server.ErrorResponse]) {
	ms := &tokenomics.MiningSummary{MiningSession: &tokenomics.MiningSession{UserID: &req.Data.UserID}}
	enhancedCtx := tokenomics.ContextWithClientType(contextWithHashCode(ctx, req), req.Data.XClientType)
	if err := s.tokenomicsProcessor.StartNewMiningSession(enhancedCtx, ms, req.Data.Resurrect, req.Data.SkipKYCStep); err != nil {
		err = errors.Wrapf(err, "failed to start a new mining session for userID:%v, data:%#v", req.Data.UserID, req.Data)
		switch {
		case errors.Is(err, tokenomics.ErrNegativeMiningProgressDecisionRequired):
			if tErr := terror.As(err); tErr != nil {
				return nil, server.Conflict(err, resurrectionDecisionRequiredErrorCode, tErr.Data)
			}

			fallthrough
		case errors.Is(err, tokenomics.ErrKYCRequired):
			if tErr := terror.As(err); tErr != nil {
				return nil, server.Conflict(err, kycStepRequiredErrorCode, tErr.Data)
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
//	@Router			/tokenomics/{userId}/extra-bonus-claims [POST].
func (s *service) ClaimExtraBonus( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[ClaimExtraBonusRequestBody, tokenomics.ExtraBonusSummary],
) (*server.Response[tokenomics.ExtraBonusSummary], *server.Response[server.ErrorResponse]) {
	resp := &tokenomics.ExtraBonusSummary{UserID: req.Data.UserID}
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
//	@Router			/tokenomics/{userId}/pre-staking [PUT].
func (s *service) StartOrUpdatePreStaking( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[StartOrUpdatePreStakingRequestBody, tokenomics.PreStakingSummary],
) (*server.Response[tokenomics.PreStakingSummary], *server.Response[server.ErrorResponse]) {
	const maxAllocation = 100
	if req.Data.Years > tokenomics.MaxPreStakingYears {
		req.Data.Years = tokenomics.MaxPreStakingYears
	}
	if req.Data.Allocation > maxAllocation {
		req.Data.Allocation = maxAllocation
	}
	st := &tokenomics.PreStakingSummary{
		PreStaking: &tokenomics.PreStaking{
			UserID:     req.Data.UserID,
			Years:      uint64(req.Data.Years),
			Allocation: float64(req.Data.Allocation),
		},
	}
	if err := s.tokenomicsProcessor.StartOrUpdatePreStaking(contextWithHashCode(ctx, req), st); err != nil {
		err = errors.Wrapf(err, "failed to StartOrUpdatePreStaking for %#v", req.Data)
		switch {
		case errors.Is(err, tokenomics.ErrDecreasingPreStakingAllocationOrYearsNotAllowed):
			return nil, server.ForbiddenWithCode(err, decreasingPreStakingAllocationOrYearsNotAllowedErrorCode)
		case errors.Is(err, tokenomics.ErrRelationNotFound):
			return nil, server.NotFound(err, userNotFoundErrorCode)
		}

		return nil, server.Unexpected(err)
	}

	return server.OK(st), nil
}
