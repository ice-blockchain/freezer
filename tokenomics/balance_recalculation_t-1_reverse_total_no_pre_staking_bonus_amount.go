// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/time"
)

//nolint:dupl // Prefer decoupling.
func (s *balanceRecalculationTriggerStreamSource) calculateElapsedTMinus1ReverseDurations(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) (previousElapsedDuration, nowElapsedDuration stdlibtime.Duration) {
	if details.TMinus1LastMiningStartedAt == nil {
		return 0, 0
	}
	totalBalance := s.getBalance(false, details.reverseTMinus1TypeDetail(), balancesByPK)
	if totalBalance == nil || totalBalance.UpdatedAt == nil {
		return 0, now.Sub(*details.TMinus1LastMiningStartedAt.Time)
	}
	if details.TMinus1LastMiningEndedAt.Before(*now.Time) && totalBalance.UpdatedAt.Before(*details.TMinus1LastMiningEndedAt.Time) {
		previousElapsedDuration = details.TMinus1LastMiningEndedAt.Sub(*totalBalance.UpdatedAt.Time)
		nowElapsedDuration = now.Sub(*details.TMinus1LastMiningEndedAt.Time)
	}
	if details.TMinus1PreviousMiningEndedAt != nil &&
		details.TMinus1PreviousMiningEndedAt.Before(*totalBalance.UpdatedAt.Time) &&
		details.TMinus1LastMiningEndedAt.After(*now.Time) &&
		details.TMinus1LastMiningStartedAt.Before(*now.Time) &&
		totalBalance.UpdatedAt.Before(*details.TMinus1LastMiningStartedAt.Time) {
		previousElapsedDuration = details.TMinus1LastMiningStartedAt.Sub(*totalBalance.UpdatedAt.Time)
		nowElapsedDuration = now.Sub(*details.TMinus1LastMiningStartedAt.Time)
	}
	if nowElapsedDuration == 0 {
		nowElapsedDuration = now.Sub(*totalBalance.UpdatedAt.Time)
	}

	return previousElapsedDuration, nowElapsedDuration
}

//nolint:gocognit // Hard to improve.
func (s *balanceRecalculationTriggerStreamSource) processPreviousIncompleteTMinus1ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	if elapsedDuration == 0 ||
		details.TMinus1UserID == details.UUserID ||
		details.TMinus1UserID == details.T0UserID ||
		details.TMinus1UserID == "" ||
		details.TMinus1LastMiningEndedAt == nil {
		return
	}
	isPositiveMining := details.LastMiningEndedAt.After(*now.Time)
	isWithinCurrentTMinus1PositiveMiningSession := details.TMinus1LastMiningStartedAt.Before(*details.LastMiningEndedAt.Time) &&
		details.TMinus1LastMiningEndedAt.After(*details.LastMiningEndedAt.Time)
	wasPreviousTMinus1MiningPositive := isWithinCurrentTMinus1PositiveMiningSession ||
		(details.TMinus1PreviousMiningEndedAt != nil && details.TMinus1PreviousMiningEndedAt.After(*details.LastMiningEndedAt.Time))
	if wasPreviousTMinus1MiningPositive && !isPositiveMining {
		s.mintTMinus1ReverseTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration)
	} else if !wasPreviousTMinus1MiningPositive {
		s.slashTMinus1ReverseTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, true)
	}
}

func (s *balanceRecalculationTriggerStreamSource) rollbackTMinus1ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	if details.TMinus1UserID == details.UUserID ||
		details.TMinus1UserID == details.T0UserID ||
		details.TMinus1UserID == "" ||
		details.TMinus1LastMiningEndedAt == nil ||
		details.TMinus1LastMiningEndedAt.Before(*now.Time) {
		return
	}
	negativeBalance := s.getBalance(true, details.reverseTMinus1TypeDetail(), balancesByPK)
	if negativeBalance == nil || negativeBalance.Amount.IsZero() {
		return
	}
	if details.TMinus1RollbackUsedAt != nil {
		positiveBalance := s.getOrInitBalance(false, details.reverseTMinus1TypeDetail(), details.UUserID, balancesByPK)
		positiveBalance.add(negativeBalance.Amount)
	}
	negativeBalance.Amount = coin.ZeroICEFlakes()
}

func (s *balanceRecalculationTriggerStreamSource) processTMinus1ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	if details.TMinus1UserID == details.UUserID ||
		details.TMinus1UserID == details.T0UserID ||
		details.TMinus1UserID == "" ||
		details.TMinus1LastMiningEndedAt == nil {
		return
	}
	defer func() {
		s.getBalance(false, lastXMiningSessionsReverseTMinus1TypeDetail, balancesByPK).Amount = nil
	}()
	isPositiveMining := details.LastMiningEndedAt.After(*now.Time)
	isRefPositiveMining := details.TMinus1LastMiningEndedAt.After(*now.Time)
	if isPositiveMining && isRefPositiveMining {
		s.mintTMinus1ReverseTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration)
	} else if !isRefPositiveMining {
		s.slashTMinus1ReverseTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, false)
	}
}

func (s *balanceRecalculationTriggerStreamSource) mintTMinus1ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	params := &userMiningRateRecalculationParameters{T2: 1}
	mintedAmount := s.calculateMintedStandardCoins(details.BaseMiningRate, params, elapsedDuration, true)
	positiveBalance := s.getOrInitBalance(false, details.reverseTMinus1TypeDetail(), details.UUserID, balancesByPK)
	positiveBalance.add(mintedAmount)
	positiveTotalThisMiningSessionBalance := s.getOrInitBalance(false, s.reverseTMinus1ThisDurationDegradationReferenceTypeDetail(details, now), details.UUserID, balancesByPK) //nolint:lll // .
	positiveTotalThisMiningSessionBalance.add(mintedAmount)
}

//nolint:dupl,revive // Prefer decoupling.
func (s *balanceRecalculationTriggerStreamSource) slashTMinus1ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, previous bool,
) {
	positiveBalance := s.getOrInitBalance(false, details.reverseTMinus1TypeDetail(), details.UUserID, balancesByPK)
	if positiveBalance.Amount.IsZero() {
		return
	}
	aggressive := details.TMinus1LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	var referenceAmount *coin.ICEFlake
	if aggressive {
		referenceAmount = s.getBalance(false, details.reverseTMinus1AggressiveDegradationReferenceTypeDetail(), balancesByPK).Amount
	} else {
		referenceAmount = s.getBalance(false, lastXMiningSessionsReverseTMinus1TypeDetail, balancesByPK).Amount
	}
	slashedAmount := s.calculateDegradation(elapsedDuration, referenceAmount, aggressive)
	positiveBalance.subtract(slashedAmount)
	if details.TMinus1RollbackUsedAt == nil || (previous && details.TMinus1LastMiningEndedAt.After(*now.Time) && details.TMinus1RollbackUsedAt.Equal(*details.TMinus1LastMiningStartedAt.Time)) { //nolint:lll // .
		negativeBalance := s.getOrInitBalance(true, details.reverseTMinus1TypeDetail(), details.UUserID, balancesByPK)
		negativeBalance.add(slashedAmount)
	}
}

func (s *balanceRecalculationTriggerStreamSource) processDegradationForTMinus1ReverseTotalNoPreStakingBonusBalanceType( //nolint:gocognit // Barely.
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	if details.TMinus1UserID == details.UUserID ||
		details.TMinus1UserID == details.T0UserID ||
		details.TMinus1UserID == "" ||
		details.TMinus1LastMiningEndedAt == nil {
		return
	}
	isPositiveMining := details.TMinus1LastMiningEndedAt.After(*now.Time)
	s.processLastXPositiveMiningSessions(balancesByPK, false, details.reverseTMinus1TypeDetail()+"/&", lastXMiningSessionsReverseTMinus1TypeDetail, details.UUserID) //nolint:lll // .

	aggressive := details.TMinus1LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	referenceBalance := s.getBalance(false, details.reverseTMinus1AggressiveDegradationReferenceTypeDetail(), balancesByPK)
	if !isPositiveMining && aggressive && (referenceBalance == nil || referenceBalance.Amount.IsNil()) {
		positiveBalance := s.getOrInitBalance(false, details.reverseTMinus1TypeDetail(), details.UUserID, balancesByPK)
		referenceBalance = s.getOrInitBalance(false, details.reverseTMinus1AggressiveDegradationReferenceTypeDetail(), details.UUserID, balancesByPK)
		referenceBalance.add(positiveBalance.Amount)
	}
	if isPositiveMining && referenceBalance != nil && !referenceBalance.Amount.IsZero() {
		referenceBalance.Amount = coin.ZeroICEFlakes()
	}
}

func (s *balanceRecalculationTriggerStreamSource) reverseTMinus1ThisDurationDegradationReferenceTypeDetail(details *BalanceRecalculationDetails, now *time.Time) string { //nolint:lll // .
	return fmt.Sprintf("%v/&%v", details.reverseTMinus1TypeDetail(), s.lastXMiningSessionsCollectingIntervalDateFormat(now))
}

const (
	lastXMiningSessionsReverseTMinus1TypeDetail = reverseTMinus1BalanceTypeDetail + "/0"
)

func (d *BalanceRecalculationDetails) reverseTMinus1TypeDetail() string {
	return fmt.Sprintf("%v_%v", reverseTMinus1BalanceTypeDetail, d.TMinus1UserID)
}

func (d *BalanceRecalculationDetails) reverseTMinus1AggressiveDegradationReferenceTypeDetail() string {
	return fmt.Sprintf("%v_", d.reverseTMinus1TypeDetail())
}

func (d *BalanceRecalculationDetails) reverseTMinus1Changed(typeDetail string) bool {
	if !strings.HasPrefix(typeDetail, reverseTMinus1BalanceTypeDetail+"_") {
		return false
	}
	userID := strings.Replace(typeDetail, reverseTMinus1BalanceTypeDetail+"_", "", 1)
	userID = strings.Replace(userID, "_", "", 1)
	userID = strings.Split(userID, "/")[0]

	return d.TMinus1UserID != userID
}
