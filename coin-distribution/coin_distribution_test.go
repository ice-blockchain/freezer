// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
)

func TestFullCoinDistribution(t *testing.T) { //nolint:paralleltest,funlen //.
	const testUserName = "testUser"

	rpc, privateKey, contractAddr, targetAddr :=
		os.Getenv("TEST_ETH_ENDPOINT_RPC"),
		os.Getenv("TEST_ETH_PRIVATE_KEY"),
		os.Getenv("TEST_ETH_CONTRACT_ADDRESS"),
		os.Getenv("ETH_TARGET_ADDRESS")
	if rpc == "" || privateKey == "" || contractAddr == "" || targetAddr == "" {
		t.Skip("skip full coin distribution test")
	}

	cl := mustNewEthClient(context.TODO(), rpc, privateKey, contractAddr)
	require.NotNil(t, cl)
	defer cl.Close()

	conf := new(config)
	conf.Ethereum.ContractAddress = contractAddr
	conf.Ethereum.ChainID = 5
	conf.Ethereum.RPC = rpc
	conf.Ethereum.PrivateKey = privateKey
	conf.BatchSize = 1
	conf.Workers = 2

	t.Run("AddPendingEntry", func(t *testing.T) {
		db := storage.MustConnect(context.TODO(), ddl, applicationYamlKey)
		defer db.Close()

		helperTruncatePendingTransactions(context.TODO(), t, db)

		const stmt = `
		INSERT INTO pending_coin_distributions
			(created_at, day, internal_id, iceflakes, user_id, eth_address)
		VALUES (now(), CURRENT_DATE, 1, '10000000000000000000000000'::uint256, $1, $2)
		ON CONFLICT (user_id,day) DO NOTHING
		`

		_, err := storage.Exec(context.TODO(), db, stmt, testUserName, targetAddr)
		require.NoError(t, err)
	})

	cd := mustCreateCoinDistributionFromConfig(context.TODO(), conf, cl)
	require.NotNil(t, cd)

	chBatches := make(chan *batch, 1)
	chTracker := make(chan []*string, 1)
	cd.MustStart(context.TODO(), chBatches, chTracker)

	var processedBatch *batch
	for {
		select {
		case b := <-chBatches:
			t.Logf("batch: %+v processed", b)
			for i := range b.Records {
				t.Logf("record: %+v processed: %v", b.Records[i].EthTX, b.Records[i].EthStatus)
				require.Equal(t, ethApiStatusAccepted, b.Records[i].EthStatus)
			}
			processedBatch = b

		case txs := <-chTracker:
			require.NotNil(t, processedBatch)
			found := false
			for i := range txs {
				t.Logf("transaction: %v processed", *txs[i])
				for recordNum := range processedBatch.Records {
					if processedBatch.Records[recordNum].EthTX != nil && *processedBatch.Records[recordNum].EthTX == *txs[i] {
						t.Logf("transaction: %+v found in batch", *txs[i])
						found = true
					}
				}
			}

			require.True(t, found)

			return
		}
	}
}
