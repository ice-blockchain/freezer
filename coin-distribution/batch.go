// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ice-blockchain/wintr/log"
)

func (r *batchRecord) Address() common.Address {
	return common.HexToAddress(r.EthAddress)
}

func (r *batchRecord) Amount() *big.Int {
	const base = 10

	value, ok := big.NewInt(0).SetString(r.Iceflakes, base)
	if !ok {
		log.Panic(fmt.Sprintf("failed to parse amount %q of user %q", r.Iceflakes, r.UserID))
	}

	return value
}

func (b *batch) Prepare() ([]common.Address, []*big.Int) {
	users := make(map[common.Address]*big.Int, len(b.Records))
	for idx := range b.Records {
		addr := b.Records[idx].Address()
		amount := b.Records[idx].Amount()
		if prev, ok := users[addr]; ok {
			amount = prev.Add(prev, amount)
		}
		users[addr] = amount
	}

	addresses := make([]common.Address, 0, len(users))
	amounts := make([]*big.Int, 0, len(users))
	for addr, amount := range users {
		addresses = append(addresses, addr)
		amounts = append(amounts, amount)
	}

	return addresses, amounts
}

func (b *batch) Users() []string {
	users := make([]string, len(b.Records)) //nolint:makezero //.
	for idx := range b.Records {
		users[idx] = b.Records[idx].UserID
	}

	return users
}

func (b *batch) SetStatus(status ethApiStatus) {
	for idx := range b.Records {
		b.Records[idx].EthStatus = status
	}
}

func (b *batch) SetAccepted(txHash string) {
	for idx := range b.Records {
		b.Records[idx].EthStatus = ethApiStatusAccepted
		b.Records[idx].EthTX = &txHash
	}
}
