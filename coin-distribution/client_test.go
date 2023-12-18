// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type (
	mockedDummyEthClient struct {
		dropErr error
		gas     int64
	}
	mockedEthSubscription struct{}
)

func (*mockedEthSubscription) Unsubscribe() {}

func (*mockedEthSubscription) Err() <-chan error {
	ch := make(chan error)
	close(ch)

	return ch
}

func (m *mockedDummyEthClient) SuggestGasPrice(context.Context) (*big.Int, error) {
	if m.gas == 0 {
		m.gas = rand.Int63n(10_000) + 1 //nolint:gosec //.
	}

	m.gas += rand.Int63n(1_000) + 1 //nolint:gosec //.

	return big.NewInt(m.gas), nil
}

func (m *mockedDummyEthClient) Airdrop(context.Context, *big.Int, *big.Int, []common.Address, []*big.Int) (string, error) {
	if m.dropErr != nil {
		return "", m.dropErr
	}

	return fmt.Sprintf("%10d", rand.Int63n(10_000_000_000)), nil //nolint:gosec //.
}

func (*mockedDummyEthClient) Close() error {
	return nil
}

func (*mockedDummyEthClient) ListenForBlocks(context.Context, chan<- *types.Header) (ethereum.Subscription, error) {
	return new(mockedEthSubscription), nil
}

func (*mockedDummyEthClient) TransactionsInBlockByHash(context.Context, common.Hash) (types.Transactions, error) {
	return nil, nil
}

func (*mockedDummyEthClient) TransactionStatus(context.Context, common.Hash) (successful bool, err error) {
	return true, nil
}

func (*mockedDummyEthClient) TransactionsStatus(context.Context, []*string) (map[ethTxStatus][]string, error) {
	return nil, nil //nolint:nilnil //.
}
