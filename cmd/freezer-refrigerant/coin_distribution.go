// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	coindistribution "github.com/ice-blockchain/freezer/coin-distribution"
	"github.com/ice-blockchain/wintr/server"
)

func (s *service) setupCoinDistributionRoutes(router *server.Router) {
	router.
		Group("/v1w").
		POST("/getCoinDistributionsForReview", server.RootHandler(s.GetCoinDistributionsForReview)).
		POST("/reviewDistributions", server.RootHandler(s.ReviewCoinDistributions))
}

// GetCoinDistributionsForReview godoc
//
//	@Schemes
//	@Description	Fetches data of pending coin distributions for review.
//	@Tags			CoinDistribution
//	@Accept			json
//	@Produce		json
//	@Param			Authorization				header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			x_client_type				query		string	false	"the type of the client calling this API. I.E. `web`"
//	@Param			cursor						query		uint64	true	"current cursor to fetch data from"	default(0)
//	@Param			limit						query		uint64	false	"count of records in response, 5000 by default"
//	@Param			createdAtOrderBy			query		string	false	"if u want to order by createdAt"								Enums(asc,desc)
//	@Param			iceOrderBy					query		string	false	"if u want to order by ice amount"								Enums(asc,desc)
//	@Param			usernameOrderBy				query		string	false	"if u want to order by username lexicographically"				Enums(asc,desc)
//	@Param			referredByUsernameOrderBy	query		string	false	"if u want to order by referredByUsername lexicographically"	Enums(asc,desc)
//	@Param			usernameKeyword				query		string	false	"if u want to find usernames starting with keyword"
//	@Param			referredByUsernameKeyword	query		string	false	"if u want to find referredByUsernames starting with keyword"
//	@Success		200							{object}	coindistribution.CoinDistributionsForReview
//	@Failure		401							{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403							{object}	server.ErrorResponse	"if not allowed"
//	@Failure		422							{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500							{object}	server.ErrorResponse
//	@Failure		504							{object}	server.ErrorResponse	"if request times out"
//	@Router			/getCoinDistributionsForReview [POST].
func (s *service) GetCoinDistributionsForReview( //nolint:gocritic // .
	ctx context.Context,
	req *server.Request[coindistribution.GetCoinDistributionsForReviewArg, coindistribution.CoinDistributionsForReview],
) (*server.Response[coindistribution.CoinDistributionsForReview], *server.Response[server.ErrorResponse]) {
	if req.AuthenticatedUser.Role != adminRole {
		return nil, server.Forbidden(errors.Errorf("insufficient role: %v, admin role required", req.AuthenticatedUser.Role))
	}
	if req.Data.Limit == 0 {
		req.Data.Limit = defaultDistributionLimit
	}
	if req.Data.CreatedAtOrderBy != "" && !strings.EqualFold(req.Data.CreatedAtOrderBy, "desc") && !strings.EqualFold(req.Data.CreatedAtOrderBy, "asc") {
		return nil, server.UnprocessableEntity(errors.Errorf("`createdAtOrderBy` has to be `asc` or `desc`"), "invalid params")
	}
	if req.Data.IceOrderBy != "" && !strings.EqualFold(req.Data.IceOrderBy, "desc") && !strings.EqualFold(req.Data.IceOrderBy, "asc") {
		return nil, server.UnprocessableEntity(errors.Errorf("`iceOrderBy` has to be `asc` or `desc`"), "invalid params")
	}
	if req.Data.UsernameOrderBy != "" && !strings.EqualFold(req.Data.UsernameOrderBy, "desc") && !strings.EqualFold(req.Data.UsernameOrderBy, "asc") {
		return nil, server.UnprocessableEntity(errors.Errorf("`usernameOrderBy` has to be `asc` or `desc`"), "invalid params")
	}
	if req.Data.ReferredByUsernameOrderBy != "" && !strings.EqualFold(req.Data.ReferredByUsernameOrderBy, "desc") && !strings.EqualFold(req.Data.ReferredByUsernameOrderBy, "asc") { //nolint:lll // .
		return nil, server.UnprocessableEntity(errors.Errorf("`referredByUsernameOrderBy` has to be `asc` or `desc`"), "invalid params")
	}
	resp, err := s.coinDistributionRepository.GetCoinDistributionsForReview(ctx, req.Data)
	if err != nil {
		return nil, server.Unexpected(errors.Wrapf(err, "failed to GetCoinDistributionsForReview for %#v", req.Data))
	}

	return server.OK(resp), nil
}

// ReviewCoinDistributions godoc
//
//	@Schemes
//	@Description	Reviews Coin Distributions.
//	@Tags			CoinDistribution
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			x_client_type	query	string	false	"the type of the client calling this API. I.E. `web`"
//	@Param			decision		query	string	true	"the decision for the current coin distributions"	Enums(approve,deny)
//	@Success		200				"OK"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/reviewDistributions [POST].
func (s *service) ReviewCoinDistributions( //nolint:gocritic // .
	ctx context.Context,
	req *server.Request[struct {
		Decision string `form:"decision" required:"true" swaggerignore:"true" enums:"approve,deny"`
	}, any],
) (*server.Response[any], *server.Response[server.ErrorResponse]) {
	if req.AuthenticatedUser.Role != adminRole {
		return nil, server.Forbidden(errors.Errorf("insufficient role: %v, admin role required", req.AuthenticatedUser.Role))
	}
	if !strings.EqualFold(req.Data.Decision, "approve") && !strings.EqualFold(req.Data.Decision, "deny") {
		return nil, server.UnprocessableEntity(errors.Errorf("`decision` has to be `approve` or `deny`"), "invalid params")
	}
	if err := s.coinDistributionRepository.ReviewCoinDistributions(ctx, req.AuthenticatedUser.UserID, req.Data.Decision); err != nil {
		return nil, server.Unexpected(errors.Wrapf(err, "failed to ReviewCoinDistributions for adminUserID:%v,decision:%v", req.AuthenticatedUser.UserID, req.Data.Decision)) //nolint:lll // .
	}

	return server.OK[any](), nil
}
