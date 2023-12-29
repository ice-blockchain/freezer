// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"sync"
	"syscall"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/wintr/log"
)

type (
	mockedDummyEthClient struct {
		dropErr error
		gas     int64
	}
	mockedAirDropper struct {
		errBefore int
	}
	mockedGasGetter struct {
		val int64
	}
)

func (m *mockedDummyEthClient) SuggestGasPrice(context.Context) (*big.Int, error) {
	if m.gas == 0 {
		m.gas = rand.Int63n(10_000) + 1 //nolint:gosec //.
	}

	m.gas += rand.Int63n(1_000) + 1 //nolint:gosec //.

	return big.NewInt(m.gas), nil
}

func (m *mockedDummyEthClient) Airdrop(context.Context, *big.Int, gasGetter, []common.Address, []*big.Int) (string, error) {
	if m.dropErr != nil {
		return "", m.dropErr
	}

	return fmt.Sprintf("%10d", rand.Int63n(10_000_000_000)), nil //nolint:gosec //.
}

func (*mockedDummyEthClient) Close() error {
	return nil
}

func (*mockedDummyEthClient) TransactionsStatus(context.Context, []*string) (map[ethTxStatus][]string, error) {
	return nil, nil //nolint:nilnil //.
}

func (*mockedDummyEthClient) TransactionStatus(context.Context, string) (ethTxStatus, error) {
	return ethTxStatusSuccessful, nil
}

func (m *mockedAirDropper) AirdropToWallets(opts *bind.TransactOpts, _ []common.Address, _ []*big.Int) (*types.Transaction, error) {
	if m.errBefore > 0 {
		m.errBefore--

		log.Info(fmt.Sprintf("airdropper: error(s) left: %v", m.errBefore))

		return nil, &net.OpError{Err: syscall.ECONNRESET}
	}

	log.Info(fmt.Sprintf("airdropper: gas price %v, limit %v", opts.GasPrice.String(), opts.GasLimit))

	return types.NewTransaction(
			0,
			common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
			big.NewInt(0),
			0,
			big.NewInt(0),
			nil,
		),
		nil
}

func (m *mockedGasGetter) GetGasOptions(context.Context) (*big.Int, uint64, error) {
	m.val++

	log.Info(fmt.Sprintf("gas getter: %v", m.val))

	return big.NewInt(m.val), uint64(m.val), nil
}

func TestGasPriceUpdateDuringRetry(t *testing.T) {
	t.Parallel()

	const errCount = 3

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	dropper := &mockedAirDropper{errBefore: errCount}

	impl := new(ethClientImpl)
	impl.Mutex = new(sync.Mutex)
	impl.Key = privateKey
	impl.AirDropper = dropper
	gasGetter := new(mockedGasGetter)

	_, err = impl.Airdrop(context.TODO(), big.NewInt(1), gasGetter, []common.Address{{1}}, []*big.Int{big.NewInt(1)})
	require.NoError(t, err)

	require.Zero(t, dropper.errBefore)
	require.Equal(t, errCount+1, int(gasGetter.val))
}
