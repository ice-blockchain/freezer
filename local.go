// SPDX-License-Identifier: ice License 1.0

package main

import (
	"flag"

	"github.com/ice-blockchain/freezer/tokenomics/fixture"
	"github.com/ice-blockchain/freezer/tokenomics/seeding"
	serverauthfixture "github.com/ice-blockchain/wintr/auth/fixture"
	"github.com/ice-blockchain/wintr/log"
)

//nolint:gochecknoglobals // Because those are flags
var (
	generateAuth   = flag.String("generateAuth", "", "generate a new auth for a random user, with the specified role")
	startSeeding   = flag.Bool("startSeeding", false, "whether to start seeding a remote database or not")
	startLocalType = flag.String("type", "all", "the strategy to use to spin up the local environment")
)

func main() {
	flag.Parse()
	if generateAuth != nil && *generateAuth != "" {
		userID, token := testingAuthorization(*generateAuth)
		log.Info("UserID")
		log.Info("=================================================================================")
		log.Info(userID)
		log.Info("Authorization Bearer Token")
		log.Info("=================================================================================")
		log.Info(token)

		return
	}
	if *startSeeding {
		seeding.StartSeeding()

		return
	}

	fixture.StartLocalTestEnvironment(fixture.StartLocalTestEnvironmentType(*startLocalType))
}

func testingAuthorization(role string) (userID, token string) {
	return serverauthfixture.CreateUser(role)
}
