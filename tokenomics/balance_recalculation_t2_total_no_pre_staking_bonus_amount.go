// SPDX-License-Identifier: ice License 1.0

//nolint:dupl // .
package tokenomics

import (
	"fmt"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/time"
)

func (s *balanceRecalculationTriggerStreamSource) processPreviousIncompleteT2TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	if elapsedDuration == 0 {
		return
	}
	if isPositiveMining := details.LastMiningEndedAt.After(*now.Time); isPositiveMining { // This means that the previous one was negative.
		s.slashT2TotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, true)
	} else {
		s.mintT2TotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, true)
	}
}

func (s *balanceRecalculationTriggerStreamSource) rollbackT2TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	if details.LastMiningEndedAt.Before(*now.Time) {
		return
	}
	negativeBalance := s.getBalance(true, t2BalanceTypeDetail, balancesByPK)
	if negativeBalance == nil || negativeBalance.Amount.IsZero() {
		return
	}
	if details.RollbackUsedAt != nil {
		positiveAmount := s.getOrInitBalance(false, t2BalanceTypeDetail, details.UUserID, balancesByPK)
		positiveAmount.add(negativeBalance.Amount)
	}
	negativeBalance.Amount = coin.ZeroICEFlakes()
}

func (s *balanceRecalculationTriggerStreamSource) processT2TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance,
	aggregatedPendingBalances map[bool]*balance,
	details *BalanceRecalculationDetails,
	now *time.Time,
	elapsedDuration stdlibtime.Duration,
) {
	defer func() {
		s.getBalance(false, lastXMiningSessionsT2TypeDetail, balancesByPK).Amount = nil
	}()
	if aggregatedPendingBalances != nil && aggregatedPendingBalances[false] != nil {
		positiveBalance := s.getOrInitBalance(false, t2BalanceTypeDetail, details.UUserID, balancesByPK)
		positiveBalance.subtract(aggregatedPendingBalances[false].Amount)
	}
	if aggregatedPendingBalances != nil && aggregatedPendingBalances[true] != nil {
		negativeBalance := s.getOrInitBalance(true, t2BalanceTypeDetail, details.UUserID, balancesByPK)
		negativeBalance.subtract(aggregatedPendingBalances[true].Amount)
	}
	if isPositiveMining := details.LastMiningEndedAt.After(*now.Time); isPositiveMining {
		s.mintT2TotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, false)
	} else {
		s.slashT2TotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, false)
	}
}

func (s *balanceRecalculationTriggerStreamSource) mintT2TotalNoPreStakingBonusBalanceType( //nolint:revive // Nope.
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, previous bool,
) {
	if details.T2 == 0 {
		return
	}
	params := &userMiningRateRecalculationParameters{T2: details.T2}
	mintedAmount := s.calculateMintedStandardCoins(details.BaseMiningRate, params, elapsedDuration, true)
	positiveBalance := s.getOrInitBalance(false, t2BalanceTypeDetail, details.UUserID, balancesByPK)
	positiveBalance.add(mintedAmount)
	positiveTotalThisMiningSessionBalance := s.getOrInitBalance(false, s.t2ThisDurationDegradationReferenceTypeDetail(now), details.UUserID, balancesByPK)
	positiveTotalThisMiningSessionBalance.add(mintedAmount)
	if previous {
		degradationReference := s.getOrInitBalance(false, degradationT0T1T2TotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
		degradationReference.add(s.getBalance(false, lastXMiningSessionsT2TypeDetail, balancesByPK).Amount)
		degradationReference.add(mintedAmount)
	}
}

//nolint:revive // Prefer decoupling.
func (s *balanceRecalculationTriggerStreamSource) slashT2TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, previous bool,
) {
	positiveBalance := s.getOrInitBalance(false, t2BalanceTypeDetail, details.UUserID, balancesByPK)
	if positiveBalance.Amount.IsZero() {
		return
	}
	aggressive := details.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	var referenceAmount *coin.ICEFlake
	if aggressive {
		referenceAmount = s.getBalance(false, aggressiveDegradationT2ReferenceBalanceTypeDetail, balancesByPK).Amount
	} else {
		referenceAmount = s.getBalance(false, lastXMiningSessionsT2TypeDetail, balancesByPK).Amount
	}
	negativeThisDuration := s.getOrInitBalance(true, s.lastXMiningSessionsThisDurationTypeDetail(previous), details.UUserID, balancesByPK)
	slashedAmount := s.calculateDegradation(elapsedDuration, referenceAmount, aggressive)
	positiveBalance.subtract(slashedAmount)
	negativeThisDuration.add(slashedAmount)
	if details.RollbackUsedAt == nil || (previous && details.RollbackUsedAt.Equal(*details.LastMiningStartedAt.Time)) {
		negativeBalance := s.getOrInitBalance(true, t2BalanceTypeDetail, details.UUserID, balancesByPK)
		negativeBalance.add(slashedAmount)
	}
}

func (s *balanceRecalculationTriggerStreamSource) processDegradationForT2TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	isPositiveMining := details.LastMiningEndedAt.After(*now.Time)
	s.processLastXPositiveMiningSessions(balancesByPK, false, t2BalanceTypeDetail+"/&", lastXMiningSessionsT2TypeDetail, details.UUserID)
	if isPositiveMining {
		degradationReference := s.getOrInitBalance(false, degradationT0T1T2TotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
		degradationReference.add(s.getBalance(false, lastXMiningSessionsT2TypeDetail, balancesByPK).Amount)
	}

	aggressive := details.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	referenceBalance := s.getBalance(false, aggressiveDegradationT2ReferenceBalanceTypeDetail, balancesByPK)
	if !isPositiveMining && aggressive && (referenceBalance == nil || referenceBalance.Amount.IsNil()) {
		positiveBalance := s.getOrInitBalance(false, t2BalanceTypeDetail, details.UUserID, balancesByPK)
		referenceBalance = s.getOrInitBalance(false, aggressiveDegradationT2ReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
		referenceBalance.add(positiveBalance.Amount)
	}
	if isPositiveMining && referenceBalance != nil && !referenceBalance.Amount.IsZero() {
		referenceBalance.Amount = coin.ZeroICEFlakes()
	}
}

func (s *balanceRecalculationTriggerStreamSource) t2ThisDurationDegradationReferenceTypeDetail(now *time.Time) string {
	return fmt.Sprintf("%v/&%v", t2BalanceTypeDetail, s.lastXMiningSessionsCollectingIntervalDateFormat(now))
}

const (
	lastXMiningSessionsT2TypeDetail = t2BalanceTypeDetail + "/0"
)
