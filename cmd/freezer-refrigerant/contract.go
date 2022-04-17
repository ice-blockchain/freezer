// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"github.com/ICE-Blockchain/freezer/economy"
)

// Private API.

const applicationYamlKey = "cmd/freezer-refrigerant"

//nolint:gochecknoglobals // Because its loaded once, at runtime.
var cfg config

type (
	// | service implements server.State and is responsible for managing the state and lifecycle of the package.
	service struct {
		economyProcessor economy.Processor
	}
	config struct {
		Version string `yaml:"version"`
	}
)
