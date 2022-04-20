// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ICE-Blockchain/freezer/economy/fixture"
	"github.com/ICE-Blockchain/wintr/log"
)

func main() {
	cleanUP := fixture.TestSetup()
	defer cleanUP()
	defer log.Info("stopping test environment, locally...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Info("started test environment, locally")
	<-quit
}
