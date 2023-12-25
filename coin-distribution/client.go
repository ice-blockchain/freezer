// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"math/big"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
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

func handleRPCError(ctx context.Context, target error) (retryAfter time.Duration) {
	var sysErr *syscall.Errno
	if errors.As(target, &sysErr) {
		return time.Second
	}

	var opErr *net.OpError
	if errors.As(target, &opErr) {
		return time.Second * 5
	}

	var httpErr *rpc.HTTPError
	if errors.As(target, &httpErr) {
		if httpErr.StatusCode == http.StatusTooManyRequests {
			return time.Hour
		} else if httpErr.StatusCode >= http.StatusInternalServerError {
			return time.Minute
		}

		return 0
	}

	if errors.Is(target, core.ErrIntrinsicGas) {
		log.Error(errors.Wrap(sendEthereumGasLimitTooLowSlackMessage(ctx, target.Error()), "failed to send slack message"))

		return time.Minute * 10
	}

	for _, ethErr := range []error{
		core.ErrNonceTooLow,
		core.ErrNonceMax,
		core.ErrInsufficientFundsForTransfer,
		core.ErrMaxInitCodeSizeExceeded,
		core.ErrInsufficientFunds,
		core.ErrTxTypeNotSupported,
		core.ErrSenderNoEOA,
		core.ErrBlobFeeCapTooLow,
	} {
		if errors.Is(target, ethErr) {
			return 0
		}
	}

	return time.Minute
}

func maybeRetryRPCRequest[T any](ctx context.Context, fn func() (T, error)) (val T, err error) {
main:
	for attempt := 1; ctx.Err() == nil; attempt++ {
		val, err = fn()
		if err == nil {
			return val, nil
		}

		retryAfter := handleRPCError(ctx, err)
		if retryAfter > 0 {
			log.Error(errors.Wrapf(err, "failed to call ethereum RPC (attempt %v), retrying after %v", attempt, retryAfter.String()))
		} else {
			log.Error(errors.Wrapf(err, "failed to call ethereum RPC (attempt %v), unrecoverable error", attempt))

			return val, multierror.Append(errClientUncoverable, err)
		}

		retryTimer := time.NewTimer(retryAfter)
		select {
		case <-ctx.Done():
			retryTimer.Stop()

			break main

		case <-retryTimer.C:
			retryTimer.Stop()
		}
	}

	return val, multierror.Append(err, ctx.Err())
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

func (ec *ethClientImpl) CreateTransactionOpts(ctx context.Context, gasPrice, chanID *big.Int, gasLimit uint64) *bind.TransactOpts {
	opts, err := bind.NewKeyedTransactorWithChainID(ec.Key, chanID)
	log.Panic(errors.Wrap(err, "failed to create transaction options")) //nolint:revive,nolintlint //.
	opts.Context = ctx
	opts.Value = big.NewInt(0)
	opts.GasLimit = gasLimit
	opts.GasPrice = gasPrice

	return opts
}

func (ec *ethClientImpl) Airdrop(ctx context.Context, chanID *big.Int, gas gasGetter, recipients []common.Address, amounts []*big.Int) (string, error) {
	fn := func() (string, error) {
		gasPrice, gasLimit, err := gas.GetGasOptions(ctx)
		if err != nil {
			return "", errors.Wrap(err, "failed to get gas options")
		}

		opts := ec.CreateTransactionOpts(ctx, gasPrice, chanID, gasLimit)
		tx, err := ec.AirdropToWallets(opts, recipients, amounts)
		if err != nil {
			return "", err //nolint:wrapcheck //.
		}

		return tx.Hash().String(), nil
	}

	return maybeRetryRPCRequest(ctx, fn)
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
