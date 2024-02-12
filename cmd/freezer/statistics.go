// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/server"
)

func (s *service) setupStatisticsRoutes(router *server.Router) {
	router.
		Group("/v1r").
		GET("/tokenomics-statistics/top-miners", server.RootHandler(s.GetTopMiners)).
		GET("/tokenomics-statistics/adoption", server.RootHandler(s.GetAdoption)).
		GET("/tokenomics-statistics/total-coins", server.RootHandler(s.GetTotalCoins))
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
//	@Header			200				{integer}	X-Next-Offset			"if this value is 0, pagination stops, if not, use it in the `offset` query param for the next call. "
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
	resp, nextOffset, err := s.tokenomicsRepository.GetTopMiners(ctx, req.Data.Keyword, req.Data.Limit, req.Data.Offset)
	if err != nil {
		return nil, server.Unexpected(errors.Wrapf(err, "failed to get top miners for userID:%v & req:%#v", req.AuthenticatedUser.UserID, req.Data))
	}

	return &server.Response[[]*tokenomics.Miner]{
		Code:    http.StatusOK,
		Data:    &resp,
		Headers: map[string]string{"X-Next-Offset": strconv.FormatUint(nextOffset, 10)},
	}, nil
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

// GetTotalCoins godoc
//
//	@Schemes
//	@Description	Returns statistics about total coins, with an usecase breakdown.
//	@Tags			Statistics
//	@Accept			json
//	@Produce		json
//	@Param			days	query		uint64	false	"number of days in the past to look for. Defaults to 3. Max is 90."
//	@Param			tz		query		string	false	"Timezone in format +04:30 or -03:45"
//	@Success		200		{object}	tokenomics.TotalCoinsSummary
//	@Failure		400		{object}	server.ErrorResponse	"if validations failed"
//	@Failure		401		{object}	server.ErrorResponse	"if not authorized"
//	@Failure		422		{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500		{object}	server.ErrorResponse
//	@Failure		504		{object}	server.ErrorResponse	"if request times out"
//	@Router			/tokenomics-statistics/total-coins [GET].
func (s *service) GetTotalCoins( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetTotalCoinsArg, tokenomics.TotalCoinsSummary],
) (*server.Response[tokenomics.TotalCoinsSummary], *server.Response[server.ErrorResponse]) {
	const defaultDays, maxDays = 3, 90
	if req.Data.Days == 0 {
		req.Data.Days = defaultDays
	}
	if req.Data.Days > maxDays {
		req.Data.Days = maxDays
	}
	if req.Data.TZ == "" {
		req.Data.TZ = "+00:00"
	}
	utcOffset, err := stdlibtime.ParseDuration(strings.Replace(req.Data.TZ+"m", ":", "h", 1))
	if err != nil {
		return nil, server.UnprocessableEntity(errors.Wrapf(err, "invalid timezone:`%v`", req.Data.TZ), invalidPropertiesErrorCode)
	}
	resp, err := s.tokenomicsRepository.GetTotalCoinsSummary(ctx, req.Data.Days, utcOffset)
	if err != nil {
		return nil, server.Unexpected(errors.Wrapf(err, "failed to GetTotalCoinsSummary for userID:%v,req:%#v", req.AuthenticatedUser.UserID, req.Data))
	}

	return server.OK(resp), nil
}
