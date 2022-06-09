// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/economy"
	"github.com/ice-blockchain/wintr/server"
)

func (s *service) setupEconomyRoutes(router *gin.Engine) {
	router.
		Group("/v1").
		GET("/economy/user-economy/:userId", server.RootHandler(newRequestGetUserEconomy, s.GetUserEconomy)).
		GET("/economy/top-miners", server.RootHandler(newRequestGetTopMiners, s.GetTopMiners)).
		GET("/economy/estimated-earnings", server.RootHandler(newRequestGetEstimatedEarnings, s.GetEstimatedEarnings)).
		GET("/economy/adoption", server.RootHandler(newRequestGetAdoption, s.GetAdoption)).
		GET("/economy/user-stats", server.RootHandler(newRequestGetUserStats, s.GetUserStats))
}

// GetUserEconomy godoc
// @Schemes
// @Description  Returns the user's personal economy
// @Tags         Economy
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        userId         path      string  true  "ID of the user"
// @Success      200            {object}  economy.UserEconomy
// @Failure      400            {object}  server.ErrorResponse  "if validations fail"
// @Failure      401            {object}  server.ErrorResponse  "if not authorized"
// @Failure      404            {object}  server.ErrorResponse  "if not found"
// @Failure      422            {object}  server.ErrorResponse  "if syntax fails"
// @Failure      500            {object}  server.ErrorResponse
// @Failure      504            {object}  server.ErrorResponse  "if request times out"
// @Router       /economy/user-economy/{userId} [GET].
func (s *service) GetUserEconomy(ctx context.Context, r server.ParsedRequest) server.Response {
	req := r.(*RequestGetUserEconomy)

	// If true user is trying to get own personal economy, otherwise another user economy.
	ownEconomy := req.AuthenticatedUser.ID == req.UserID
	ue, err := s.economyRepository.GetUserEconomy(ctx, req.UserID, ownEconomy)
	if err != nil {
		if errors.Is(err, economy.ErrNotFound) {
			return userNotFound(err)
		}

		return server.Unexpected(err)
	}

	return server.OK(ue)
}

func newRequestGetUserEconomy() server.ParsedRequest {
	return new(RequestGetUserEconomy)
}

func (req *RequestGetUserEconomy) SetAuthenticatedUser(user server.AuthenticatedUser) {
	if req.AuthenticatedUser.ID == "" {
		req.AuthenticatedUser = user
	}
}

func (req *RequestGetUserEconomy) GetAuthenticatedUser() server.AuthenticatedUser {
	return req.AuthenticatedUser
}

func (req *RequestGetUserEconomy) Validate() *server.Response {
	return server.RequiredStrings(map[string]string{"userId": req.UserID})
}

func (req *RequestGetUserEconomy) Bindings(c *gin.Context) []func(obj interface{}) error {
	return []func(obj interface{}) error{c.ShouldBindUri, server.ShouldBindAuthenticatedUser(c)}
}

// GetTopMiners godoc
// @Schemes
// @Description  Returns the paginated leaderboard with top miners.
// @Tags         Economy
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        keyword        query     string  false  "a keyword to look for"
// @Param        limit          query     uint64  false  "max number of elements to return"
// @Param        offset         query     uint64  false  "number of elements to skip before starting to fetch data"
// @Success      200            {array}   economy.TopMiner
// @Failure      400            {object}  server.ErrorResponse  "if validations fail"
// @Failure      401            {object}  server.ErrorResponse  "if not authorized"
// @Failure      422            {object}  server.ErrorResponse  "if syntax fails"
// @Failure      500            {object}  server.ErrorResponse
// @Failure      504            {object}  server.ErrorResponse  "if request times out"
// @Router       /economy/top-miners [GET].
func (s *service) GetTopMiners(ctx context.Context, r server.ParsedRequest) server.Response {
	resp, err := s.economyRepository.GetTopMiners(ctx, &r.(*RequestGetTopMiners).GetTopMinersArg)
	if err != nil {
		return server.Unexpected(err)
	}

	return server.OK(resp)
}

func newRequestGetTopMiners() server.ParsedRequest {
	return new(RequestGetTopMiners)
}

func (req *RequestGetTopMiners) SetAuthenticatedUser(user server.AuthenticatedUser) {
	if req.AuthenticatedUser.ID == "" {
		req.AuthenticatedUser = user
	}
}

func (req *RequestGetTopMiners) GetAuthenticatedUser() server.AuthenticatedUser {
	return req.AuthenticatedUser
}

func (req *RequestGetTopMiners) Validate() *server.Response {
	if req.Limit == 0 {
		req.Limit = cfg.DefaultPagination.Limit
	}

	if req.Limit > cfg.DefaultPagination.MaxLimit {
		req.Limit = cfg.DefaultPagination.MaxLimit
	}

	return nil
}

func (req *RequestGetTopMiners) Bindings(c *gin.Context) []func(obj interface{}) error {
	return []func(obj interface{}) error{server.ShouldBindAuthenticatedUser(c)}
}

// GetEstimatedEarnings godoc
// @Schemes
// @Description  Returns estimated earnings based on the provided parameters.
// @Tags         Economy
// @Accept       json
// @Produce      json
// @Param        Authorization      header    string  true   "Insert your access token"  default(Bearer <Add access token here>)
// @Param        t0                 query     bool    false  "if the user that referred you should be active or not"
// @Param        t1                 query     uint64  false  "number of t1 active referrals you desire"
// @Param        t2                 query     uint64  false  "number of t2 active referrals you desire"
// @Param        stakingYears       query     uint8   false  "number of years you want to enable staking for"
// @Param        stakingAllocation  query     uint8   false  "the percentage [0..100] of your balance you want to stake"
// @Success      200                {object}  economy.EstimatedEarnings
// @Failure      400                {object}  server.ErrorResponse  "if validations fail"
// @Failure      401                {object}  server.ErrorResponse  "if not authorized"
// @Failure      422                {object}  server.ErrorResponse  "if syntax fails"
// @Failure      500                {object}  server.ErrorResponse
// @Failure      504                {object}  server.ErrorResponse  "if request times out"
// @Router       /economy/estimated-earnings [GET].
func (s *service) GetEstimatedEarnings(ctx context.Context, r server.ParsedRequest) server.Response {
	resp, err := s.economyRepository.GetEstimatedEarnings(ctx, &r.(*RequestGetEstimatedEarnings).GetEstimatedEarningsArg)
	if err != nil {
		return server.Unexpected(errors.Wrapf(err, "failed to get estimated earnings for userID:%v", r.(*RequestGetEstimatedEarnings).AuthenticatedUser.ID))
	}

	return server.OK(resp)
}

func newRequestGetEstimatedEarnings() server.ParsedRequest {
	return new(RequestGetEstimatedEarnings)
}

func (req *RequestGetEstimatedEarnings) SetAuthenticatedUser(user server.AuthenticatedUser) {
	if req.AuthenticatedUser.ID == "" {
		req.AuthenticatedUser = user
	}
}

func (req *RequestGetEstimatedEarnings) GetAuthenticatedUser() server.AuthenticatedUser {
	return req.AuthenticatedUser
}

func (req *RequestGetEstimatedEarnings) Validate() *server.Response {
	//nolint:gomnd // Not a magic number, its 100%.
	if req.StakingAllocation > 100 {
		return server.BadRequest(errors.Errorf("staking allocation has to be within [0,100], %v is invalid", req.StakingAllocation), "MISSING_PROPERTIES")
	}

	return nil
}

func (req *RequestGetEstimatedEarnings) Bindings(c *gin.Context) []func(obj interface{}) error {
	return []func(obj interface{}) error{c.ShouldBindQuery, server.ShouldBindAuthenticatedUser(c)}
}

// GetAdoption godoc
// @Schemes
// @Description  Returns the current adoption information.
// @Tags         Economy
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Insert your access token"  default(Bearer <Add access token here>)
// @Success      200            {object}  economy.Adoption
// @Failure      400            {object}  server.ErrorResponse  "if validations fail"
// @Failure      401            {object}  server.ErrorResponse  "if not authorized"
// @Failure      422            {object}  server.ErrorResponse  "if syntax fails"
// @Failure      500            {object}  server.ErrorResponse
// @Failure      504            {object}  server.ErrorResponse  "if request times out"
// @Router       /economy/adoption [GET].
func (s *service) GetAdoption(ctx context.Context, r server.ParsedRequest) server.Response {
	resp, err := s.economyRepository.GetAdoption(ctx)
	if err != nil {
		return server.Unexpected(errors.Wrapf(err, "failed to get current adoption for userID:%v", r.(*RequestGetAdoption).AuthenticatedUser.ID))
	}

	return server.OK(resp)
}

func newRequestGetAdoption() server.ParsedRequest {
	return new(RequestGetAdoption)
}

func (req *RequestGetAdoption) SetAuthenticatedUser(user server.AuthenticatedUser) {
	if req.AuthenticatedUser.ID == "" {
		req.AuthenticatedUser = user
	}
}

func (req *RequestGetAdoption) GetAuthenticatedUser() server.AuthenticatedUser {
	return req.AuthenticatedUser
}

func (req *RequestGetAdoption) Validate() *server.Response {
	return nil
}

func (req *RequestGetAdoption) Bindings(c *gin.Context) []func(obj interface{}) error {
	return []func(obj interface{}) error{server.ShouldBindAuthenticatedUser(c)}
}

// GetUserStats godoc
// @Schemes
// @Description  Returns statistics about the user population.
// @Tags         Economy
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Insert your access token"  default(Bearer <Add access token here>)
// @Param        lastNoOfDays   query     uint16  false  "number of days in the past to look for. Defaults to 7."
// @Success      200            {object}  economy.UserStats
// @Failure      400            {object}  server.ErrorResponse  "if validations fail"
// @Failure      401            {object}  server.ErrorResponse  "if not authorized"
// @Failure      422            {object}  server.ErrorResponse  "if syntax fails"
// @Failure      500            {object}  server.ErrorResponse
// @Failure      504            {object}  server.ErrorResponse  "if request times out"
// @Router       /economy/user-stats [GET].
func (s *service) GetUserStats(ctx context.Context, r server.ParsedRequest) server.Response {
	resp, err := s.economyRepository.GetUserStats(ctx, r.(*RequestGetUserStats).LastNoOfDays)
	if err != nil {
		return server.Unexpected(errors.Wrapf(err, "failed to get user stats for userID:%v", r.(*RequestGetUserStats).AuthenticatedUser.ID))
	}

	return server.OK(resp)
}

func newRequestGetUserStats() server.ParsedRequest {
	return new(RequestGetUserStats)
}

func (req *RequestGetUserStats) SetAuthenticatedUser(user server.AuthenticatedUser) {
	if req.AuthenticatedUser.ID == "" {
		req.AuthenticatedUser = user
	}
}

func (req *RequestGetUserStats) GetAuthenticatedUser() server.AuthenticatedUser {
	return req.AuthenticatedUser
}

func (req *RequestGetUserStats) Validate() *server.Response {
	if req.LastNoOfDays == 0 {
		req.LastNoOfDays = defaultLastNoOfDays
	}

	return nil
}

func (req *RequestGetUserStats) Bindings(c *gin.Context) []func(obj interface{}) error {
	return []func(obj interface{}) error{c.ShouldBindQuery, server.ShouldBindAuthenticatedUser(c)}
}

func userNotFound(err error) server.Response {
	return server.Response{
		Data: server.ErrorResponse{
			Error: err.Error(),
			Code:  "USER_NOT_FOUND",
		}.Fail(err),
		Code: http.StatusNotFound,
	}
}
