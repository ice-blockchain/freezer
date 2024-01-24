// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/server"
	"github.com/ice-blockchain/wintr/time"
)

func (s *service) setupTokenomicsRoutes(router *server.Router) {
	router.
		Group("/v1r").
		GET("/tokenomics/:userId/mining-summary", server.RootHandler(s.GetMiningSummary)).
		GET("/tokenomics/:userId/pre-staking-summary", server.RootHandler(s.GetPreStakingSummary)).
		GET("/tokenomics/:userId/balance-summary", server.RootHandler(s.GetBalanceSummary)).
		GET("/tokenomics/:userId/balance-history", server.RootHandler(s.GetBalanceHistory)).
		GET("/tokenomics/:userId/ranking-summary", server.RootHandler(s.GetRankingSummary))
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
//	@Router			/tokenomics/{userId}/mining-summary [GET].
func (s *service) GetMiningSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetMiningSummaryArg, tokenomics.MiningSummary],
) (*server.Response[tokenomics.MiningSummary], *server.Response[server.ErrorResponse]) {
	mining, err := s.tokenomicsRepository.GetMiningSummary(contextWithHashCode(ctx, req), req.Data.UserID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's mining summary for userID:%v", req.Data.UserID)
		if errors.Is(err, tokenomics.ErrRelationNotFound) {
			return nil, server.NotFound(err, userNotFoundErrorCode)
		}

		return nil, server.Unexpected(err)
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
//	@Router			/tokenomics/{userId}/pre-staking-summary [GET].
func (s *service) GetPreStakingSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetPreStakingSummaryArg, tokenomics.PreStakingSummary],
) (*server.Response[tokenomics.PreStakingSummary], *server.Response[server.ErrorResponse]) {
	preStaking, err := s.tokenomicsRepository.GetPreStakingSummary(contextWithHashCode(ctx, req), req.Data.UserID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's pre-staking summary for userID:%v", req.Data.UserID)
		switch {
		case errors.Is(err, tokenomics.ErrNotFound):
			return nil, server.NotFound(err, userPreStakingNotEnabledErrorCode)
		case errors.Is(err, tokenomics.ErrPrestakingDisabled):
			return nil, server.ForbiddenWithCode(err, prestakingDisabledForUser)
		default:
			return nil, server.Unexpected(err)
		}
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
//	@Router			/tokenomics/{userId}/balance-summary [GET].
func (s *service) GetBalanceSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetBalanceSummaryArg, tokenomics.BalanceSummary],
) (*server.Response[tokenomics.BalanceSummary], *server.Response[server.ErrorResponse]) {
	balance, err := s.tokenomicsRepository.GetBalanceSummary(contextWithHashCode(ctx, req), req.Data.UserID)
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
//	@Router			/tokenomics/{userId}/balance-history [GET].
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
	hist, err := s.tokenomicsRepository.GetBalanceHistory(contextWithHashCode(ctx, req), req.Data.UserID, startDate, endDate, utcOffset, req.Data.Limit, req.Data.Offset) //nolint:lll // .
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
//	@Router			/tokenomics/{userId}/ranking-summary [GET].
func (s *service) GetRankingSummary( //nolint:gocritic // False negative.
	ctx context.Context,
	req *server.Request[GetRankingSummaryArg, tokenomics.RankingSummary],
) (*server.Response[tokenomics.RankingSummary], *server.Response[server.ErrorResponse]) {
	ranking, err := s.tokenomicsRepository.GetRankingSummary(contextWithHashCode(ctx, req), req.Data.UserID)
	if err != nil {
		err = errors.Wrapf(err, "failed to get user's ranking summary for userID:%v", req.Data.UserID)
		if errors.Is(err, tokenomics.ErrGlobalRankHidden) {
			return nil, server.ForbiddenWithCode(err, globalRankHiddenErrorCode)
		}

		return nil, server.Unexpected(err)
	}

	return server.OK(ranking), nil
}
