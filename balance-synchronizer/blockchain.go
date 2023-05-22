// SPDX-License-Identifier: ice License 1.0

package balancesynchronizer

import (
	"context"
	"strconv"
	"strings"

	"github.com/ice-blockchain/wintr/coin"
)

type (
	blockchainMessage struct { // TODO: delete this and use the actual one.
		AccountAddress     string
		ICEFlake           string
		PreStakingICEFlake string
	}
)

func shouldSynchronizeBlockchainAccount(iteration uint64, usr *user) *blockchainMessage {
	if usr.MiningBlockchainAccountAddress == "" || iteration%(uint64(usr.ID)%100) != 0 {
		return nil
	}
	var standard, preStaking *coin.ICEFlake
	if standardParts := strings.Split(strconv.FormatFloat(usr.BalanceTotalStandard, 'f', 9, 64), "."); len(standardParts) == 1 {
		standard = coin.UnsafeParseAmount(standardParts[0]).MultiplyUint64(coin.Denomination)
	} else if len(standardParts) == 2 {
		standard = coin.UnsafeParseAmount(standardParts[0]).MultiplyUint64(coin.Denomination).Add(coin.UnsafeParseAmount(standardParts[1]))
	}
	if preStakingParts := strings.Split(strconv.FormatFloat(usr.BalanceTotalPreStaking, 'f', 9, 64), "."); len(preStakingParts) == 1 {
		preStaking = coin.UnsafeParseAmount(preStakingParts[0]).MultiplyUint64(coin.Denomination)
	} else if len(preStakingParts) == 2 {
		preStaking = coin.UnsafeParseAmount(preStakingParts[0]).MultiplyUint64(coin.Denomination).Add(coin.UnsafeParseAmount(preStakingParts[1]))
	}

	return &blockchainMessage{
		AccountAddress:     usr.MiningBlockchainAccountAddress,
		ICEFlake:           coin.ZeroICEFlakes().Add(standard).String(),
		PreStakingICEFlake: coin.ZeroICEFlakes().Add(preStaking).String(),
	}
}

func (bs *balanceSynchronizer) synchronizeBlockchainAccounts(ctx context.Context, msgs []*blockchainMessage) error {
	if len(msgs) == 0 || ctx.Err() != nil {
		return nil
	}

	return nil
}
