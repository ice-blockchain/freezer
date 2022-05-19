// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"github.com/ice-blockchain/freezer/economy"
	"github.com/ice-blockchain/wintr/server"
)

// Public API.
type (
	RequestStartMining struct {
		AuthenticatedUser server.AuthenticatedUser `json:"authenticatedUser" swaggerignore:"true"`
	}
	RequestStartStaking struct {
		AuthenticatedUser server.AuthenticatedUser `json:"authenticatedUser" swaggerignore:"true"`
		economy.Staking
	}
)

// Private API.

const (
	applicationYamlKey   = "cmd/freezer-refrigerant"
	miningInProgress     = "MINING_IN_PROGRESS"
	userNotFound         = "USER_NOT_FOUND"
	stakingAlradyEnabled = "STAKING_ALREADY_ENABLED"
)

//nolint:gochecknoglobals // Because its loaded once, at runtime.
var cfg config

type (
	// | service implements server.State and is responsible for managing the state and lifecycle of the package.
	service struct {
		economyProcessor economy.Processor
	}
	config struct {
		Host    string `yaml:"host"`
		Version string `yaml:"version"`
	}
)
