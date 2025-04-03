// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/server"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

type (
	GetMiningSummaryArg struct {
		UserID string `uri:"userId" required:"true" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	GetMiningBoostSummaryArg struct {
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
	GetAdoptionArg   struct{}
	GetTotalCoinsArg struct {
		TZ   string `form:"tz" example:"+4:30" allowUnauthorized:"true"`
		Days uint64 `form:"days" example:"7"`
	}
)

// Private API.

// Values for server.ErrorResponse#Code.
const (
	userPreStakingNotEnabledErrorCode = "PRE_STAKING_NOT_ENABLED"
	globalRankHiddenErrorCode         = "GLOBAL_RANK_HIDDEN"
	invalidPropertiesErrorCode        = "INVALID_PROPERTIES"
)

func (s *service) registerReadRoutes(router *server.Router) {
	s.setupTokenomicsReadRoutes(router)
	s.setupStatisticsRoutes(router)
}

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
//	@Router			/v1r/tokenomics-statistics/top-miners [GET].
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
	resp, nextOffset, err := s.tokenomicsProcessor.GetTopMiners(ctx, req.Data.Keyword, req.Data.Limit, req.Data.Offset)
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
//	@Router			/v1r/tokenomics-statistics/adoption [GET].
func (s *service) GetAdoption( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetAdoptionArg, tokenomics.AdoptionSummary],
) (*server.Response[tokenomics.AdoptionSummary], *server.Response[server.ErrorResponse]) {
	resp, err := s.tokenomicsProcessor.GetAdoptionSummary(ctx, req.AuthenticatedUser.UserID)
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
//	@Router			/v1r/tokenomics-statistics/total-coins [GET].
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
	resp, err := s.tokenomicsProcessor.GetTotalCoinsSummary(ctx, req.Data.Days, utcOffset)
	if err != nil {
		return nil, server.Unexpected(errors.Wrapf(err, "failed to GetTotalCoinsSummary for userID:%v,req:%#v", req.AuthenticatedUser.UserID, req.Data))
	}

	return server.OK(resp), nil
}

func (s *service) setupTokenomicsReadRoutes(router *server.Router) {
	router.
		Group("/v1r").
		GET("/tokenomics/:userId/mining-boost-summary", server.RootHandler(s.GetMiningBoostSummary)).
		GET("/tokenomics/:userId/mining-summary", server.RootHandler(s.GetMiningSummary)).
		GET("/tokenomics/:userId/pre-staking-summary", server.RootHandler(s.GetPreStakingSummary)).
		GET("/tokenomics/:userId/balance-summary", server.RootHandler(s.GetBalanceSummary)).
		GET("/tokenomics/:userId/balance-history", server.RootHandler(s.GetBalanceHistory)).
		GET("/tokenomics/:userId/ranking-summary", server.RootHandler(s.GetRankingSummary))
}

// GetMiningBoostSummary godoc
//
//	@Schemes
//	@Description	Returns the mining boost related information.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string	true	"ID of the user"
//	@Success		200				{object}	tokenomics.MiningBoostSummary
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		404				{object}	server.ErrorResponse	"if not found"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1r/tokenomics/{userId}/mining-boost-summary [GET].
func (s *service) GetMiningBoostSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetMiningBoostSummaryArg, tokenomics.MiningBoostSummary],
) (*server.Response[tokenomics.MiningBoostSummary], *server.Response[server.ErrorResponse]) {
	summary, err := s.tokenomicsProcessor.GetMiningBoostSummary(ctx, req.Data.UserID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's mining boost summary for userID:%v", req.Data.UserID)
		if errors.Is(err, tokenomics.ErrRelationNotFound) {
			return nil, server.NotFound(err, userNotFoundErrorCode)
		}

		return nil, server.Unexpected(err)
	}

	return server.OK(summary), nil
}

// GetMiningSummary godoc
//
//	@Schemes
//	@Description	Returns the mining related information.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string	true	"ID of the user"
//	@Success		200				{object}	tokenomics.MiningSummary
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		404				{object}	server.ErrorResponse	"if not found"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1r/tokenomics/{userId}/mining-summary [GET].
func (s *service) GetMiningSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetMiningSummaryArg, tokenomics.MiningSummary],
) (*server.Response[tokenomics.MiningSummary], *server.Response[server.ErrorResponse]) {
	mining, err := s.tokenomicsProcessor.GetMiningSummary(contextWithHashCode(ctx, req), req.Data.UserID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's mining summary for userID:%v", req.Data.UserID)
		if errors.Is(err, tokenomics.ErrRelationNotFound) {
			return nil, server.NotFound(err, userNotFoundErrorCode)
		}

		return nil, server.Unexpected(err)
	}
	if cfg.Tenant == tokeroTenant {
		mining.MiningRates.Type = tokenomics.NoneMiningRateType
		mining.MiningRates.Total.Amount = "0.00"
	}

	return server.OK(mining), nil
}

// GetPreStakingSummary godoc
//
//	@Schemes
//	@Description	Returns the pre-staking related information.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string	true	"ID of the user"
//	@Success		200				{object}	tokenomics.PreStakingSummary
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		404				{object}	server.ErrorResponse	"if not found"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1r/tokenomics/{userId}/pre-staking-summary [GET].
func (s *service) GetPreStakingSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetPreStakingSummaryArg, tokenomics.PreStakingSummary],
) (*server.Response[tokenomics.PreStakingSummary], *server.Response[server.ErrorResponse]) {
	preStaking, err := s.tokenomicsProcessor.GetPreStakingSummary(contextWithHashCode(ctx, req), req.Data.UserID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's pre-staking summary for userID:%v", req.Data.UserID)
		if errors.Is(err, tokenomics.ErrNotFound) {
			return nil, server.NotFound(err, userPreStakingNotEnabledErrorCode)
		}

		return nil, server.Unexpected(err)
	}

	return server.OK(preStaking), nil
}

// GetBalanceSummary godoc
//
//	@Schemes
//	@Description	Returns the balance related information.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string	true	"ID of the user"
//	@Success		200				{object}	tokenomics.BalanceSummary
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1r/tokenomics/{userId}/balance-summary [GET].
func (s *service) GetBalanceSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetBalanceSummaryArg, tokenomics.BalanceSummary],
) (*server.Response[tokenomics.BalanceSummary], *server.Response[server.ErrorResponse]) {
	balance, err := s.tokenomicsProcessor.GetBalanceSummary(contextWithHashCode(ctx, req), req.Data.UserID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's balance summary for userID:%v", req.Data.UserID)

		return nil, server.Unexpected(err)
	}

	return server.OK(balance), nil
}

// GetBalanceHistory godoc
//
//	@Schemes
//	@Description	Returns the balance history for the provided params.
//	@Description	If `startDate` is after `endDate`, we go backwards in time: I.E. today, yesterday, etc.
//	@Description	If `startDate` is before `endDate`, we go forwards in time: I.E. today, tomorrow, etc.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string	true	"ID of the user"
//	@Param			startDate		query		string	false	"The start date in RFC3339 or ISO8601 formats. Default is `now` in UTC."
//	@Param			endDate			query		string	false	"The start date in RFC3339 or ISO8601 formats. Default is `end of day, relative to startDate`."
//	@Param			tz				query		string	false	"The user's timezone. I.E. `+03:00`, `-1:30`. Default is UTC."
//	@Param			limit			query		uint64	false	"max number of elements to return. Default is `24`."
//	@Param			offset			query		uint64	false	"number of elements to skip before starting to fetch data"
//	@Success		200				{array}		tokenomics.BalanceHistoryEntry
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if not allowed"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1r/tokenomics/{userId}/balance-history [GET].
func (s *service) GetBalanceHistory( //nolint:gocritic,funlen // False negative.
	ctx context.Context,
	req *server.Request[GetBalanceHistoryArg, []*tokenomics.BalanceHistoryEntry],
) (*server.Response[[]*tokenomics.BalanceHistoryEntry], *server.Response[server.ErrorResponse]) {
	const defaultLimit, maxLimit = 24, 1000
	if req.Data.Limit > maxLimit {
		req.Data.Limit = maxLimit
	}
	if req.Data.Limit == 0 {
		req.Data.Limit = defaultLimit
	}
	var startDate, endDate *time.Time
	if req.Data.StartDate == nil {
		startDate = time.Now()
	} else {
		startDate = time.New(*req.Data.StartDate)
	}
	if req.Data.EndDate == nil {
		endDate = time.New(startDate.Add(-1 * users.NanosSinceMidnight(startDate)))
	} else {
		endDate = time.New(*req.Data.EndDate)
	}
	if req.Data.TZ == "" {
		req.Data.TZ = "+00:00"
	}
	utcOffset, err := stdlibtime.ParseDuration(strings.Replace(req.Data.TZ+"m", ":", "h", 1))
	if err != nil {
		return nil, server.UnprocessableEntity(errors.Wrapf(err, "invalid timezone:`%v`", req.Data.TZ), invalidPropertiesErrorCode)
	}
	hist, err := s.tokenomicsProcessor.GetBalanceHistory(contextWithHashCode(ctx, req), req.Data.UserID, startDate, endDate, utcOffset, req.Data.Limit, req.Data.Offset) //nolint:lll // .
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's balance history for userID:%v, data:%#v", req.Data.UserID, req.Data)

		return nil, server.Unexpected(err)
	}

	return server.OK(&hist), nil
}

// GetRankingSummary godoc
//
//	@Schemes
//	@Description	Returns the ranking related information.
//	@Tags			Tokenomics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			userId			path		string	true	"ID of the user"
//	@Success		200				{object}	tokenomics.RankingSummary
//	@Failure		400				{object}	server.ErrorResponse	"if validations fail"
//	@Failure		401				{object}	server.ErrorResponse	"if not authorized"
//	@Failure		403				{object}	server.ErrorResponse	"if hidden by the user"
//	@Failure		422				{object}	server.ErrorResponse	"if syntax fails"
//	@Failure		500				{object}	server.ErrorResponse
//	@Failure		504				{object}	server.ErrorResponse	"if request times out"
//	@Router			/v1r/tokenomics/{userId}/ranking-summary [GET].
func (s *service) GetRankingSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetRankingSummaryArg, tokenomics.RankingSummary],
) (*server.Response[tokenomics.RankingSummary], *server.Response[server.ErrorResponse]) {
	ranking, err := s.tokenomicsProcessor.GetRankingSummary(contextWithHashCode(ctx, req), req.Data.UserID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's ranking summary for userID:%v", req.Data.UserID)
		if errors.Is(err, tokenomics.ErrGlobalRankHidden) {
			return nil, server.ForbiddenWithCode(err, globalRankHiddenErrorCode)
		}

		return nil, server.Unexpected(err)
	}

	return server.OK(ranking), nil
}
