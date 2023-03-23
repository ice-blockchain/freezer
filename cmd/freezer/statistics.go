// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/server"
)

func (s *service) setupStatisticsRoutes(router *server.Router) {
	router.
		Group("/v1r").
		GET("/tokenomics-statistics/top-miners", server.RootHandler(s.GetTopMiners)).
		GET("/tokenomics-statistics/adoption", server.RootHandler(s.GetAdoption))
}

// GetTopMiners godoc
//
//	@Schemes
//	@Description	Returns the paginated leaderboard with top miners.
//	@Tags			Statistics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			keyword			query		string	false	"a keyword to look for in the user's username or firstname/lastname"
//	@Param			limit			query		uint64	false	"max number of elements to return. Default is `10`."
//	@Param			offset			query		uint64	false	"number of elements to skip before starting to fetch data"
//	@Success		200				{array}		tokenomics.Miner
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/tokenomics-statistics/top-miners [GET].
func (s *service) GetTopMiners( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetTopMinersArg, []*tokenomics.Miner],
) (*server.Response[[]*tokenomics.Miner], *server.Response[server.ErrorResponse]) {
	const defaultLimit, maxLimit = 10, 1000
	if req.Data.Limit == 0 {
		req.Data.Limit = defaultLimit
	}
	if req.Data.Limit > maxLimit {
		req.Data.Limit = maxLimit
	}
	resp, err := s.tokenomicsRepository.GetTopMiners(ctx, req.Data.Keyword, req.Data.Limit, req.Data.Offset)
	if err != nil {
		return nil, server.Unexpected(errors.Wrapf(err, "failed to get top miners for userID:%v & req:%#v", req.AuthenticatedUser.UserID, req.Data))
	}

	return server.OK(&resp), nil
}

// GetAdoption godoc
//
//	@Schemes
//	@Description	Returns the current adoption information.
//	@Tags			Statistics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	tokenomics.AdoptionSummary
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/tokenomics-statistics/adoption [GET].
func (s *service) GetAdoption( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetAdoptionArg, tokenomics.AdoptionSummary],
) (*server.Response[tokenomics.AdoptionSummary], *server.Response[server.ErrorResponse]) {
	resp, err := s.tokenomicsRepository.GetAdoptionSummary(ctx)
	if err != nil {
		return nil, server.Unexpected(errors.Wrapf(err, "failed to get adoption summary for userID:%v", req.AuthenticatedUser.UserID))
	}

	return server.OK(resp), nil
}
