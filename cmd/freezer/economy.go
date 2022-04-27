// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ICE-Blockchain/freezer/economy"
	"github.com/ICE-Blockchain/wintr/server"
)

func (s *service) setupEconomyRoutes(router *gin.Engine) {
	router.
		Group("/v1").
		GET("/economy/user-economy/:userId", server.RootHandler(newRequestGetUserEconomy, s.GetUserEconomy)).
		GET("/economy/top-miners", server.RootHandler(newRequestGetTopMiners, s.GetTopMiners))
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
// @Param        Authorization  header    string  true   "Insert your access token"  default(Bearer <Add access token here>)
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
	// Req := r.(*RequestGetTopMiners).

	//nolint:nolintlint // TODO implement me.

	return server.OK([]*economy.TopMiner{{
		UserID:            "bogus",
		ProfilePictureURL: "bogus",
		Balance:           1.2, //nolint:gomnd // It will be implementated later.
	}})
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
	return nil
}

func (req *RequestGetTopMiners) Bindings(c *gin.Context) []func(obj interface{}) error {
	return []func(obj interface{}) error{server.ShouldBindAuthenticatedUser(c)}
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
