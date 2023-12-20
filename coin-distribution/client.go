// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	coindistribution "github.com/ice-blockchain/freezer/coin-distribution/internal"
	"github.com/ice-blockchain/wintr/log"
)

func mustNewEthClient(ctx context.Context, endpoint, privateKey, contract string) *ethClientImpl {
	key, err := crypto.HexToECDSA(privateKey)
	log.Panic(errors.Wrap(err, "failed to parse private key")) //nolint:revive,nolintlint //.

	rpcClient, err := ethclient.DialContext(ctx, endpoint)
	log.Panic(errors.Wrap(err, "failed to connect to ethereum RPC")) //nolint:revive,nolintlint //.

	distributor, err := coindistribution.NewCoindistribution(common.HexToAddress(contract), rpcClient)
	log.Panic(errors.Wrap(err, "failed to create contract instance")) //nolint:revive,nolintlint //.

	return &ethClientImpl{
		RPC:        rpcClient,
		AirDropper: distributor,
		Key:        key,
		Mutex:      new(sync.Mutex),
	}
}

func maybeRetryRPCRequest[T any](ctx context.Context, fn func() (T, error)) (val T, err error) {
	var httpErr *rpc.HTTPError

	for attempt := 1; ctx.Err() == nil; attempt++ {
		val, err = fn()
		if err == nil {
			return val, nil
		}

		log.Error(errors.Wrapf(err, "failed to call ethereum RPC (attempt %v)", attempt))
		if errors.As(err, &httpErr) {
			switch httpErr.StatusCode {
			case http.StatusInternalServerError, http.StatusTooManyRequests:
			default:
				return val, err
			}
		}

		// In case any other error occurred (network timeout, dns, etc), retry after a second.
		time.Sleep(time.Second)
	}

	return val, err
}

func (ec *ethClientImpl) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return maybeRetryRPCRequest(ctx, func() (*big.Int, error) {
		return ec.RPC.SuggestGasPrice(ctx) //nolint:wrapcheck //.
	})
}

func (ec *ethClientImpl) AirdropToWallets(opts *bind.TransactOpts, recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	// AirdropToWallets() uses PendingNonceAt() to get the next nonce internally which is not thread-safe.
	ec.Mutex.Lock()
	defer ec.Mutex.Unlock()

	return ec.AirDropper.AirdropToWallets(opts, recipients, amounts) //nolint:wrapcheck //.
}

func (ec *ethClientImpl) CreateTransactionOpts(ctx context.Context, gasPrice, chanID *big.Int) *bind.TransactOpts {
	opts, err := bind.NewKeyedTransactorWithChainID(ec.Key, chanID)
	log.Panic(errors.Wrap(err, "failed to create transaction options")) //nolint:revive,nolintlint //.
	opts.Context = ctx
	opts.Value = big.NewInt(0)
	opts.GasLimit = gasLimit
	opts.GasPrice = gasPrice

	return opts
}

func (ec *ethClientImpl) Airdrop(ctx context.Context, gasPrice, chanID *big.Int, recipients []common.Address, amounts []*big.Int) (string, error) {
	return maybeRetryRPCRequest(ctx, func() (string, error) {
		tx, err := ec.AirDropper.AirdropToWallets(ec.CreateTransactionOpts(ctx, gasPrice, chanID), recipients, amounts)
		if err != nil {
			return "", err //nolint:wrapcheck //.
		}

		return tx.Hash().String(), nil
	})
}

func (ec *ethClientImpl) TransactionsStatus(ctx context.Context, hashes []*string) (statuses map[ethTxStatus][]string, err error) { //nolint:funlen //.
	elements := make([]rpc.BatchElem, len(hashes)) //nolint:makezero //.
	results := make([]*types.Receipt, len(hashes)) //nolint:makezero //.
	for elementIdx := range elements {
		elements[elementIdx] = rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []any{*hashes[elementIdx]},
			Result: &results[elementIdx],
		}
	}

	if _, batchErr := maybeRetryRPCRequest(ctx, func() (bool, error) {
		return true, ec.RPC.Client().BatchCallContext(ctx, elements) //nolint:wrapcheck //.
	}); batchErr != nil {
		return nil, batchErr
	}

	statuses = make(map[ethTxStatus][]string)
	for elementIdx := range elements {
		receipt := results[elementIdx]
		if receipt == nil {
			// Transaction is not mined yet.
			continue
		} else if elements[elementIdx].Error != nil {
			err = multierror.Append(err, elements[elementIdx].Error)

			continue
		}

		if receipt.Status == types.ReceiptStatusSuccessful {
			statuses[ethTxStatusSuccessful] = append(statuses[ethTxStatusSuccessful], *hashes[elementIdx])
		} else {
			statuses[ethTxStatusFailed] = append(statuses[ethTxStatusFailed], *hashes[elementIdx])
		}
	}

	return statuses, err //nolint:wrapcheck //.
}

func (ec *ethClientImpl) Close() error {
	ec.RPC.Close()

	return nil
}
