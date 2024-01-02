// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"os"
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/time"
)

func pointerToString[T any](v *T) string {
	if v == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%v", *v)
}

func TestFullCoinDistribution(t *testing.T) { //nolint:paralleltest,funlen //.
	const testUserName = "testUser"

	maybeSkipTest(t)

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
	conf.Ethereum.ChainID = 97
	conf.Ethereum.RPC = rpc
	conf.Ethereum.PrivateKey = privateKey

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
	defer cd.Close()

	chBatches := make(chan *batch, 1)
	cd.MustStart(context.TODO(), chBatches)

	t.Logf("waiting for batch to be processed")
	processedBatch := <-chBatches
	t.Logf("batch: %+v processed: status %v", processedBatch, processedBatch.Status)
	for i := range processedBatch.Records {
		t.Logf("record: %v processed: %v", pointerToString(processedBatch.Records[i].EthTX), processedBatch.Records[i].EthStatus)
		require.Equal(t, ethApiStatusAccepted, processedBatch.Records[i].EthStatus)
	}
	require.Equal(t, ethTxStatusSuccessful, processedBatch.Status)
}

func TestDatabaseSetGetValues(t *testing.T) {
	var boolValue bool

	maybeSkipTest(t)

	db := storage.MustConnect(context.TODO(), ddl, applicationYamlKey)
	defer db.Close()

	err := databaseSetValue(context.TODO(), db, configKeyCoinDistributerEnabled, false)
	require.NoError(t, err)

	err = databaseGetValue(context.TODO(), db, configKeyCoinDistributerEnabled, &boolValue)
	require.NoError(t, err)
	require.False(t, boolValue)

	err = databaseSetValue(context.TODO(), db, configKeyCoinDistributerEnabled, true)
	require.NoError(t, err)

	err = databaseGetValue(context.TODO(), db, configKeyCoinDistributerEnabled, &boolValue)
	require.NoError(t, err)
	require.True(t, boolValue)

	testTime := time.New(stdlibtime.Date(2021, 1, 2, 3, 4, 5, 0, stdlibtime.UTC))
	err = databaseSetValue(context.TODO(), db, configKeyCoinDistributerMsgOnline, testTime)
	require.NoError(t, err)

	var timeValue time.Time
	err = databaseGetValue(context.TODO(), db, configKeyCoinDistributerMsgOnline, &timeValue)
	require.NoError(t, err)
	require.Equal(t, testTime, &timeValue)
}

func TestCoinDistributionWaitOK(t *testing.T) { //nolint:paralleltest,funlen //.
	const (
		testUserName = "testUserOK"
		testTxOK     = "0xAABBCCDDEE"
	)

	maybeSkipTest(t)

	cl := &mockedDummyEthClient{}
	conf := new(config)

	t.Run("AddPendingEntry", func(t *testing.T) {
		db := storage.MustConnect(context.TODO(), ddl, applicationYamlKey)
		defer db.Close()

		helperTruncatePendingTransactions(context.TODO(), t, db)

		const stmt = `
		INSERT INTO pending_coin_distributions
			(created_at, day, internal_id, iceflakes, user_id, eth_address, eth_status, eth_tx)
		VALUES (now(), CURRENT_DATE, 1, '10000000000000000000000000'::uint256, $1, $2, 'ACCEPTED', $3)
		ON CONFLICT (user_id,day) DO NOTHING
		`

		_, err := storage.Exec(context.TODO(), db, stmt, testUserName, "0x1234", testTxOK)
		require.NoError(t, err)
	})

	cd := mustCreateCoinDistributionFromConfig(context.TODO(), conf, cl)
	require.NotNil(t, cd)
	defer cd.Close()

	chBatches := make(chan *batch, 1)
	cd.MustStart(context.TODO(), chBatches)

	t.Logf("waiting for check for pending transaction")
	processedBatch := <-chBatches
	t.Logf("batch: %+v processed: status %v", processedBatch, processedBatch.Status)
	require.Equal(t, ethTxStatusSuccessful, processedBatch.Status)
	require.Equal(t, testTxOK, processedBatch.TX)
	require.Len(t, processedBatch.Records, 1)
	require.Equal(t, testTxOK, *processedBatch.Records[0].EthTX)
}

func TestCoinDistributionWaitFailed(t *testing.T) { //nolint:paralleltest,funlen //.
	const (
		testUserName = "testUserOK"
		testTxFailed = "0xAABBCCDDEE"
	)

	maybeSkipTest(t)

	cl := &mockedDummyEthClient{txErr: map[string]error{testTxFailed: nil}}
	conf := new(config)

	t.Run("AddPendingEntry", func(t *testing.T) {
		db := storage.MustConnect(context.TODO(), ddl, applicationYamlKey)
		defer db.Close()

		helperTruncatePendingTransactions(context.TODO(), t, db)

		const stmt = `
		INSERT INTO pending_coin_distributions
			(created_at, day, internal_id, iceflakes, user_id, eth_address, eth_status, eth_tx)
		VALUES (now(), CURRENT_DATE, 1, '10000000000000000000000000'::uint256, $1, $2, 'ACCEPTED', $3)
		ON CONFLICT (user_id,day) DO NOTHING
		`

		_, err := storage.Exec(context.TODO(), db, stmt, testUserName, "0x1234", testTxFailed)
		require.NoError(t, err)
	})

	cd := mustCreateCoinDistributionFromConfig(context.TODO(), conf, cl)
	require.NotNil(t, cd)
	defer cd.Close()

	chBatches := make(chan *batch, 1)
	cd.MustStart(context.TODO(), chBatches)

	t.Logf("waiting for check for pending transaction")
	processedBatch := <-chBatches
	t.Logf("batch: %+v processed: status %v", processedBatch, processedBatch.Status)
	require.Equal(t, ethTxStatusFailed, processedBatch.Status)
	require.Equal(t, testTxFailed, processedBatch.TX)
	require.Len(t, processedBatch.Records, 1)
	require.Equal(t, testTxFailed, *processedBatch.Records[0].EthTX)
	require.False(t, cd.Processor.IsEnabled(context.TODO()))
}
