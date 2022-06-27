// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"github.com/ice-blockchain/freezer/economy"
	"github.com/ice-blockchain/wintr/server"
)

// Public API.

type (
	RequestGetUserEconomy struct {
		AuthenticatedUser server.AuthenticatedUser `json:"authenticatedUser" swaggerignore:"true"`
		UserID            string                   `uri:"userId" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	RequestGetTopMiners struct {
		AuthenticatedUser server.AuthenticatedUser `json:"authenticatedUser" swaggerignore:"true"`
		economy.GetTopMinersArg
	}
	RequestGetEstimatedEarnings struct {
		AuthenticatedUser server.AuthenticatedUser `json:"authenticatedUser" swaggerignore:"true"`
		economy.GetEstimatedEarningsArg
	}
	RequestGetAdoption struct {
		AuthenticatedUser server.AuthenticatedUser `json:"authenticatedUser" swaggerignore:"true"`
	}
	RequestGetUserStats struct {
		AuthenticatedUser server.AuthenticatedUser `json:"authenticatedUser" swaggerignore:"true"`
		// Default is `7`, limited to 30.
		LastNoOfDays uint16 `form:"lastNoOfDays" example:"7"`
	}
)

// Private API.

const (
	applicationYamlKey  = "cmd/freezer"
	defaultLastNoOfDays = 7
	maximumLastNoOfDays = 30
)

//nolint:gochecknoglobals // Because its loaded once, at runtime.
var cfg config

type (
	// | service implements server.State and is responsible for managing the state and lifecycle of the package.
	service struct {
		economyRepository economy.Repository
	}
	config struct {
		Host              string `yaml:"host"`
		Version           string `yaml:"version"`
		DefaultPagination struct {
			Limit    uint64 `yaml:"limit"`
			MaxLimit uint64 `yaml:"maxLimit"`
		} `yaml:"defaultPagination"`
	}
)
