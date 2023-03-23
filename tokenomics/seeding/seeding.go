// SPDX-License-Identifier: ice License 1.0

//go:build !test

package seeding

import (
	"fmt"
	"os"
	"strings"
	stdlibtime "time"

	"github.com/ice-blockchain/go-tarantool-client"
	"github.com/ice-blockchain/wintr/log"
)

func StartSeeding() {
	before := stdlibtime.Now()
	db := dbConnector()
	defer func() {
		log.Panic(db.Close()) //nolint:revive // It doesnt really matter.
		log.Info(fmt.Sprintf("seeding finalized in %v", stdlibtime.Since(before).String()))
	}()
	log.Info("TODO: implement seeding")
}

func cleanUpWorkerSpaces(db tarantool.Connector) { //nolint:deadcode,unused // .
	tables := []string{
		"balance_recalculation_worker_",
		"extra_bonus_processing_worker_",
		"blockchain_balance_synchronization_worker_",
		"mining_rates_recalculation_worker_",
		"balance_recalculation_worker_",
		"pre_stakings_",
		"mining_sessions_dlq_",
	}
	for _, table := range tables {
		for i := 0; i < 1000; i++ {
			sql := fmt.Sprintf(`delete from %[1]v%[2]v where 1=1`, table, i)
			_, err := db.PrepareExecute(sql, map[string]any{})
			log.Panic(err)
		}
	}
}

func dbConnector() tarantool.Connector {
	parts := strings.Split(os.Getenv("MASTER_DB_INSTANCE_ADDRESS"), "@")
	userAndPass := strings.Split(parts[0], ":")
	opts := tarantool.Opts{
		Timeout:       20 * stdlibtime.Second, //nolint:gomnd // It doesnt matter here.
		Reconnect:     stdlibtime.Millisecond,
		MaxReconnects: 10, //nolint:gomnd // It doesnt matter here.
		User:          userAndPass[0],
		Pass:          userAndPass[1],
	}
	db, err := tarantool.Connect(parts[1], opts)
	log.Panic(err)

	return db
}
