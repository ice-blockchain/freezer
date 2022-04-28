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
		PATCH("/economy/start-mining", server.RootHandler(newRequestStartMining, s.StartMining)).
		PATCH("/economy/start-staking", server.RootHandler(newRequestStartStaking, s.StartStaking))
}

// StartMining godoc
// @Schemes
// @Description  Starts or resumes the mining for the authenticated user.
// @Tags         Economy
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string           true  "Insert your access token"  default(Bearer <Add access token here>)
// @Success      200            "OK"
// @Failure      400            {object}  server.ErrorResponse  "if validations fail"
// @Failure      401            {object}  server.ErrorResponse  "if not authorized"
// @Failure      409            {object}  server.ErrorResponse  "if mining is in progress"
// @Failure      422            {object}  server.ErrorResponse  "if syntax fails"
// @Failure      500            {object}  server.ErrorResponse
// @Failure      504            {object}  server.ErrorResponse  "if request times out"
// @Router       /economy/start-mining [PATCH].
func (s *service) StartMining(ctx context.Context, r server.ParsedRequest) server.Response {
	req := r.(*RequestStartMining)

	err := s.economyRepository.StartMining(ctx, req.AuthenticatedUser.ID)
	if err != nil {
		if errors.Is(err, economy.ErrMiningInProgress) {
			return server.Response{
				Code: http.StatusConflict,
				Data: server.ErrorResponse{
					Error: err.Error(),
					Code:  "MINING_IN_PROGRESS",
				}.Fail(err),
			}
		}

		return server.Unexpected(err)
	}

	return server.OK()
}

func newRequestStartMining() server.ParsedRequest {
	return new(RequestStartMining)
}

func (req *RequestStartMining) SetAuthenticatedUser(user server.AuthenticatedUser) {
	if req.AuthenticatedUser.ID == "" {
		req.AuthenticatedUser = user
	}
}

func (req *RequestStartMining) GetAuthenticatedUser() server.AuthenticatedUser {
	return req.AuthenticatedUser
}

func (req *RequestStartMining) Validate() *server.Response {
	return nil
}

func (req *RequestStartMining) Bindings(c *gin.Context) []func(obj interface{}) error {
	return []func(obj interface{}) error{server.ShouldBindAuthenticatedUser(c)}
}

// StartStaking godoc
// @Schemes
// @Description  Starts staking for the authenticated user.
// @Tags         Economy
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        request        body    economy.Staking  true  "Request params"
// @Success      200            "OK"
// @Failure      400            {object}  server.ErrorResponse  "if validations fail"
// @Failure      401            {object}  server.ErrorResponse  "if not authorized"
// @Failure      404            {object}  server.ErrorResponse  "user not found"
// @Failure      409            {object}  server.ErrorResponse  "if staking is already enabled"
// @Failure      422            {object}  server.ErrorResponse  "if syntax fails"
// @Failure      500            {object}  server.ErrorResponse
// @Failure      504            {object}  server.ErrorResponse  "if request times out"
// @Router       /economy/start-staking [PATCH].
func (s *service) StartStaking(ctx context.Context, r server.ParsedRequest) server.Response {
	req := r.(*RequestStartStaking)

	//nolint:nolintlint // TODO implement me.

	return server.OK(req)
}

func newRequestStartStaking() server.ParsedRequest {
	return new(RequestStartStaking)
}

func (req *RequestStartStaking) SetAuthenticatedUser(user server.AuthenticatedUser) {
	if req.AuthenticatedUser.ID == "" {
		req.AuthenticatedUser = user
	}
}

func (req *RequestStartStaking) GetAuthenticatedUser() server.AuthenticatedUser {
	return req.AuthenticatedUser
}

func (req *RequestStartStaking) Validate() *server.Response {
	if req.Years == 0 || req.Years > 5 {
		return invalidProperties(errors.Errorf("years `%v` is not allowed, only natural numbers within [1..5] allowed", req.Years))
	}
	if req.Percentage <= 0.0 || req.Percentage > 100.0 {
		return invalidProperties(errors.Errorf("percentage `%v` is not allowed, only within (0.0, 100.0] allowed", req.Percentage))
	}

	return nil
}

func (req *RequestStartStaking) Bindings(c *gin.Context) []func(obj interface{}) error {
	return []func(obj interface{}) error{c.ShouldBindJSON, server.ShouldBindAuthenticatedUser(c)}
}

func invalidProperties(err error) *server.Response {
	return &server.Response{
		Data: server.ErrorResponse{
			Error: err.Error(),
			Code:  "INVALID_PROPERTIES",
		}.Fail(err),
		Code: http.StatusBadRequest,
	}
}
