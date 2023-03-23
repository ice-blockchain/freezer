// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"fmt"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/time"
)

func (s *balanceRecalculationTriggerStreamSource) processPreviousIncompleteThisDurationTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, perDurationTypeDetail string,
) {
	if elapsedDuration == 0 {
		return
	}
	if isPositiveMining := details.LastMiningEndedAt.After(*now.Time); isPositiveMining { // This means that the previous one was negative.
		s.slashThisDurationTotalNoPreStakingBonusBalanceType(balancesByPK, details, perDurationTypeDetail, true)
	} else {
		s.mintThisDurationTotalNoPreStakingBonusBalanceType(balancesByPK, details, elapsedDuration, perDurationTypeDetail)
	}
}

func (s *balanceRecalculationTriggerStreamSource) processThisDurationTotalNoPreStakingBonusBalanceType( //nolint:revive // .
	balancesByPK map[string]*balance,
	aggregatedPendingBalances map[bool]*balance,
	details *BalanceRecalculationDetails,
	now *time.Time,
	elapsedDuration stdlibtime.Duration,
	perDurationTypeDetail string,
) {
	if aggregatedPendingBalances != nil && aggregatedPendingBalances[false] != nil {
		positiveBalance := s.getOrInitBalance(false, perDurationTypeDetail, details.UUserID, balancesByPK)
		positiveBalance.add(aggregatedPendingBalances[false].Amount)
	}
	if aggregatedPendingBalances != nil && aggregatedPendingBalances[true] != nil {
		negativeBalance := s.getOrInitBalance(true, perDurationTypeDetail, details.UUserID, balancesByPK)
		negativeBalance.add(aggregatedPendingBalances[true].Amount)
	}
	if isPositiveMining := details.LastMiningEndedAt.After(*now.Time); isPositiveMining {
		s.mintThisDurationTotalNoPreStakingBonusBalanceType(balancesByPK, details, elapsedDuration, perDurationTypeDetail)
	} else {
		s.slashThisDurationTotalNoPreStakingBonusBalanceType(balancesByPK, details, perDurationTypeDetail, false)
	}
}

func (s *balanceRecalculationTriggerStreamSource) mintThisDurationTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, elapsedDuration stdlibtime.Duration, perDurationTypeDetail string,
) {
	params := &userMiningRateRecalculationParameters{
		T0:         details.T0,
		T1:         details.T1,
		T2:         details.T2,
		ExtraBonus: details.ExtraBonus,
	}
	mintedAmount := s.calculateMintedStandardCoins(details.BaseMiningRate, params, elapsedDuration, false)
	positiveBalance := s.getOrInitBalance(false, perDurationTypeDetail, details.UUserID, balancesByPK)
	positiveBalance.add(mintedAmount)
}

func (s *balanceRecalculationTriggerStreamSource) slashThisDurationTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, perDurationTypeDetail string, previous bool,
) {
	totalPositiveLastXMiningSessions := s.getBalance(true, s.lastXMiningSessionsThisDurationTypeDetail(previous), balancesByPK)
	if totalPositiveLastXMiningSessions == nil || totalPositiveLastXMiningSessions.Amount.IsZero() {
		if totalPositiveLastXMiningSessions != nil {
			totalPositiveLastXMiningSessions.Amount = nil
		}

		return
	}
	negativeBalance := s.getOrInitBalance(true, perDurationTypeDetail, details.UUserID, balancesByPK)
	negativeBalance.add(totalPositiveLastXMiningSessions.Amount)
	totalPositiveLastXMiningSessions.Amount = nil
}

func (*balanceRecalculationTriggerStreamSource) lastXMiningSessionsThisDurationTypeDetail(previous bool) string {
	return fmt.Sprintf("t0+t1+t2+total/%v/0", previous)
}
