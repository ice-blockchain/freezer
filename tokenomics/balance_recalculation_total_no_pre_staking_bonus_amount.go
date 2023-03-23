// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"fmt"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/time"
)

func (s *balanceRecalculationTriggerStreamSource) processPreviousIncompleteTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	if elapsedDuration == 0 {
		return
	}
	if isPositiveMining := details.LastMiningEndedAt.After(*now.Time); isPositiveMining { // This means that the previous one was negative.
		s.slashTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, true)
	} else {
		s.mintTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, true)
	}
}

func (s *balanceRecalculationTriggerStreamSource) rollbackTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	if details.LastMiningEndedAt.Before(*now.Time) {
		return
	}
	negativeBalance := s.getBalance(true, "", balancesByPK)
	if negativeBalance == nil || negativeBalance.Amount.IsZero() {
		return
	}
	if details.RollbackUsedAt != nil {
		positiveBalance := s.getOrInitBalance(false, "", details.UUserID, balancesByPK)
		positiveBalance.add(negativeBalance.Amount)
	}
	negativeBalance.Amount = coin.ZeroICEFlakes()
}

func (s *balanceRecalculationTriggerStreamSource) processTotalNoPreStakingBonusBalanceType( //nolint:funlen,gocognit // .
	balancesByPK map[string]*balance,
	aggregatedPendingBalances map[bool]*balance,
	details *BalanceRecalculationDetails,
	now *time.Time,
	elapsedDuration stdlibtime.Duration,
) {
	defer func() {
		s.getBalance(false, lastXMiningSessionsTypeDetail, balancesByPK).Amount = nil
	}()
	isAggressiveDegradation := details.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	isPositiveMining := details.LastMiningEndedAt.After(*now.Time)
	if aggregatedPendingBalances != nil && aggregatedPendingBalances[false] != nil {
		positiveBalance := s.getOrInitBalance(false, "", details.UUserID, balancesByPK)
		positiveBalance.add(aggregatedPendingBalances[false].Amount)
		if !isPositiveMining && isAggressiveDegradation {
			referenceBalance := s.getOrInitBalance(false, aggressiveDegradationTotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
			referenceBalance.add(aggregatedPendingBalances[false].Amount)
		}
		positiveTotalThisMiningSessionBalance := s.getOrInitBalance(!isPositiveMining && !isAggressiveDegradation, s.thisDurationDegradationReferenceTypeDetail(now), details.UUserID, balancesByPK) //nolint:lll // .
		positiveTotalThisMiningSessionBalance.add(aggregatedPendingBalances[false].Amount)
	}
	if isPositiveMining {
		s.mintTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, false)
	} else {
		s.slashTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, false)
	}
	if aggregatedPendingBalances != nil && aggregatedPendingBalances[true] != nil {
		positiveBalance := s.getOrInitBalance(false, "", details.UUserID, balancesByPK)
		positiveBalance.subtract(aggregatedPendingBalances[true].Amount)
		if !isPositiveMining && isAggressiveDegradation {
			referenceBalance := s.getOrInitBalance(false, aggressiveDegradationTotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
			referenceBalance.subtract(aggregatedPendingBalances[true].Amount)
		}
		positiveTotalThisMiningSessionBalance := s.getOrInitBalance(!isPositiveMining && !isAggressiveDegradation, s.thisDurationDegradationReferenceTypeDetail(now), details.UUserID, balancesByPK) //nolint:lll // .
		positiveTotalThisMiningSessionBalance.subtract(aggregatedPendingBalances[true].Amount)
	}
}

func (s *balanceRecalculationTriggerStreamSource) mintTotalNoPreStakingBonusBalanceType( //nolint:revive // Nope.
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, previous bool,
) {
	params := &userMiningRateRecalculationParameters{ExtraBonus: details.ExtraBonus}
	mintedAmount := s.calculateMintedStandardCoins(details.BaseMiningRate, params, elapsedDuration, false)
	positiveBalance := s.getOrInitBalance(false, "", details.UUserID, balancesByPK)
	positiveBalance.add(mintedAmount)
	positiveTotalThisMiningSessionBalance := s.getOrInitBalance(false, s.thisDurationDegradationReferenceTypeDetail(now), details.UUserID, balancesByPK)
	positiveTotalThisMiningSessionBalance.add(mintedAmount)
	if previous {
		degradationReference := s.getOrInitBalance(false, degradationT0T1T2TotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
		degradationReference.add(s.getBalance(false, lastXMiningSessionsTypeDetail, balancesByPK).Amount)
		degradationReference.add(mintedAmount)
	}
}

//nolint:revive // Not a problem here.
func (s *balanceRecalculationTriggerStreamSource) slashTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, previous bool,
) {
	positiveBalance := s.getOrInitBalance(false, "", details.UUserID, balancesByPK)
	if positiveBalance.Amount.IsZero() {
		return
	}
	var slashedAmount *coin.ICEFlake
	if details.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.Available.Until).Before(*now.Time) {
		slashedAmount = positiveBalance.Amount
	} else if aggressive := details.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time); aggressive {
		slashedAmount = s.calculateDegradation(elapsedDuration, s.getBalance(false, aggressiveDegradationTotalReferenceBalanceTypeDetail, balancesByPK).Amount, aggressive) //nolint:lll // .
	} else {
		slashedAmount = s.calculateDegradation(elapsedDuration, s.getBalance(false, lastXMiningSessionsTypeDetail, balancesByPK).Amount, aggressive)
	}
	negativeThisDuration := s.getOrInitBalance(true, s.lastXMiningSessionsThisDurationTypeDetail(previous), details.UUserID, balancesByPK)
	positiveBalance.subtract(slashedAmount)
	negativeThisDuration.add(slashedAmount)
	if details.RollbackUsedAt == nil || (previous && details.RollbackUsedAt.Equal(*details.LastMiningStartedAt.Time)) {
		negativeBalance := s.getOrInitBalance(true, "", details.UUserID, balancesByPK)
		negativeBalance.add(slashedAmount)
	}
}

func (s *balanceRecalculationTriggerStreamSource) processDegradationForTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	isPositiveMining := details.LastMiningEndedAt.After(*now.Time)
	s.processLastXPositiveMiningSessions(balancesByPK, isPositiveMining, "/&", lastXMiningSessionsTypeDetail, details.UUserID)
	if isPositiveMining {
		degradationReference := s.getOrInitBalance(false, degradationT0T1T2TotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
		degradationReference.add(s.getBalance(false, lastXMiningSessionsTypeDetail, balancesByPK).Amount)
	}

	aggressive := details.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	referenceBalance := s.getBalance(false, aggressiveDegradationTotalReferenceBalanceTypeDetail, balancesByPK)
	if !isPositiveMining && aggressive && (referenceBalance == nil || referenceBalance.Amount.IsNil()) {
		positiveBalance := s.getOrInitBalance(false, "", details.UUserID, balancesByPK)
		referenceBalance = s.getOrInitBalance(false, aggressiveDegradationTotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
		referenceBalance.add(positiveBalance.Amount)
	}
	if isPositiveMining && referenceBalance != nil && !referenceBalance.Amount.IsZero() {
		referenceBalance.Amount = coin.ZeroICEFlakes()
	}
}

func (s *balanceRecalculationTriggerStreamSource) thisDurationDegradationReferenceTypeDetail(now *time.Time) string {
	return fmt.Sprintf("/&%v", s.lastXMiningSessionsCollectingIntervalDateFormat(now))
}

const (
	lastXMiningSessionsTypeDetail = "/0"
)
