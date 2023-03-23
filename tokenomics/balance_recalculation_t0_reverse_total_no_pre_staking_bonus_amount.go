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
func (s *balanceRecalculationTriggerStreamSource) calculateElapsedT0ReverseDurations(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) (previousElapsedDuration, nowElapsedDuration stdlibtime.Duration) {
	if details.T0LastMiningStartedAt == nil {
		return 0, 0
	}
	totalBalance := s.getBalance(false, details.reverseT0TypeDetail(), balancesByPK)
	if totalBalance == nil || totalBalance.UpdatedAt == nil {
		return 0, now.Sub(*details.T0LastMiningStartedAt.Time)
	}
	if details.T0LastMiningEndedAt.Before(*now.Time) && totalBalance.UpdatedAt.Before(*details.T0LastMiningEndedAt.Time) {
		previousElapsedDuration = details.T0LastMiningEndedAt.Sub(*totalBalance.UpdatedAt.Time)
		nowElapsedDuration = now.Sub(*details.T0LastMiningEndedAt.Time)
	}
	if details.T0PreviousMiningEndedAt != nil &&
		details.T0PreviousMiningEndedAt.Before(*totalBalance.UpdatedAt.Time) &&
		details.T0LastMiningEndedAt.After(*now.Time) &&
		details.T0LastMiningStartedAt.Before(*now.Time) &&
		totalBalance.UpdatedAt.Before(*details.T0LastMiningStartedAt.Time) {
		previousElapsedDuration = details.T0LastMiningStartedAt.Sub(*totalBalance.UpdatedAt.Time)
		nowElapsedDuration = now.Sub(*details.T0LastMiningStartedAt.Time)
	}
	if nowElapsedDuration == 0 {
		nowElapsedDuration = now.Sub(*totalBalance.UpdatedAt.Time)
	}

	return previousElapsedDuration, nowElapsedDuration
}

func (s *balanceRecalculationTriggerStreamSource) processPreviousIncompleteT0ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	if elapsedDuration == 0 ||
		details.T0UserID == details.UUserID ||
		details.T0UserID == "" ||
		details.T0LastMiningEndedAt == nil {
		return
	}
	isPositiveMining := details.LastMiningEndedAt.After(*now.Time)
	isWithinCurrentT0PositiveMiningSession := details.T0LastMiningStartedAt.Before(*details.LastMiningEndedAt.Time) &&
		details.T0LastMiningEndedAt.After(*details.LastMiningEndedAt.Time)
	wasPreviousT0MiningPositive := isWithinCurrentT0PositiveMiningSession ||
		(details.T0PreviousMiningEndedAt != nil && details.T0PreviousMiningEndedAt.After(*details.LastMiningEndedAt.Time))
	if wasPreviousT0MiningPositive && !isPositiveMining {
		s.mintT0ReverseTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration)
	} else if !wasPreviousT0MiningPositive {
		s.slashT0ReverseTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, true)
	}
}

func (s *balanceRecalculationTriggerStreamSource) rollbackT0ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	if details.T0UserID == details.UUserID ||
		details.T0UserID == "" ||
		details.T0LastMiningEndedAt == nil ||
		details.T0LastMiningEndedAt.Before(*now.Time) {
		return
	}
	negativeBalance := s.getBalance(true, details.reverseT0TypeDetail(), balancesByPK)
	if negativeBalance == nil || negativeBalance.Amount.IsZero() {
		return
	}
	if details.T0RollbackUsedAt != nil {
		positiveBalance := s.getOrInitBalance(false, details.reverseT0TypeDetail(), details.UUserID, balancesByPK)
		positiveBalance.add(negativeBalance.Amount)
	}
	negativeBalance.Amount = coin.ZeroICEFlakes()
}

func (s *balanceRecalculationTriggerStreamSource) processT0ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	if details.T0UserID == details.UUserID ||
		details.T0UserID == "" ||
		details.T0LastMiningEndedAt == nil {
		return
	}
	defer func() {
		s.getBalance(false, lastXMiningSessionsReverseT0TypeDetail, balancesByPK).Amount = nil
	}()
	isPositiveMining := details.LastMiningEndedAt.After(*now.Time)
	isRefPositiveMining := details.T0LastMiningEndedAt.After(*now.Time)
	if isPositiveMining && isRefPositiveMining {
		s.mintT0ReverseTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration)
	} else if !isRefPositiveMining {
		s.slashT0ReverseTotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, false)
	}
}

func (s *balanceRecalculationTriggerStreamSource) mintT0ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	params := &userMiningRateRecalculationParameters{T0: 1}
	mintedAmount := s.calculateMintedStandardCoins(details.BaseMiningRate, params, elapsedDuration, true)
	positiveBalance := s.getOrInitBalance(false, details.reverseT0TypeDetail(), details.UUserID, balancesByPK)
	positiveBalance.add(mintedAmount)
	positiveTotalThisMiningSessionBalance := s.getOrInitBalance(false, s.reverseT0ThisDurationDegradationReferenceTypeDetail(details, now), details.UUserID, balancesByPK) //nolint:lll // .
	positiveTotalThisMiningSessionBalance.add(mintedAmount)
}

//nolint:dupl,revive // Prefer decoupling.
func (s *balanceRecalculationTriggerStreamSource) slashT0ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, previous bool,
) {
	positiveBalance := s.getOrInitBalance(false, details.reverseT0TypeDetail(), details.UUserID, balancesByPK)
	if positiveBalance.Amount.IsZero() {
		return
	}
	aggressive := details.T0LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	var referenceAmount *coin.ICEFlake
	if aggressive {
		referenceAmount = s.getBalance(false, details.reverseT0AggressiveDegradationReferenceTypeDetail(), balancesByPK).Amount
	} else {
		referenceAmount = s.getBalance(false, lastXMiningSessionsReverseT0TypeDetail, balancesByPK).Amount
	}
	slashedAmount := s.calculateDegradation(elapsedDuration, referenceAmount, aggressive)
	positiveBalance.subtract(slashedAmount)
	if details.T0RollbackUsedAt == nil || (previous && details.T0LastMiningEndedAt.After(*now.Time) && details.T0RollbackUsedAt.Equal(*details.T0LastMiningStartedAt.Time)) { //nolint:lll // .
		negativeBalance := s.getOrInitBalance(true, details.reverseT0TypeDetail(), details.UUserID, balancesByPK)
		negativeBalance.add(slashedAmount)
	}
}

func (s *balanceRecalculationTriggerStreamSource) processDegradationForT0ReverseTotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	if details.T0UserID == details.UUserID ||
		details.T0UserID == "" ||
		details.T0LastMiningEndedAt == nil {
		return
	}
	isPositiveMining := details.T0LastMiningEndedAt.After(*now.Time)
	s.processLastXPositiveMiningSessions(balancesByPK, false, details.reverseT0TypeDetail()+"/&", lastXMiningSessionsReverseT0TypeDetail, details.UUserID)

	aggressive := details.T0LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	referenceBalance := s.getBalance(false, details.reverseT0AggressiveDegradationReferenceTypeDetail(), balancesByPK)
	if !isPositiveMining && aggressive && (referenceBalance == nil || referenceBalance.Amount.IsNil()) {
		positiveBalance := s.getOrInitBalance(false, details.reverseT0TypeDetail(), details.UUserID, balancesByPK)
		referenceBalance = s.getOrInitBalance(false, details.reverseT0AggressiveDegradationReferenceTypeDetail(), details.UUserID, balancesByPK)
		referenceBalance.add(positiveBalance.Amount)
	}
	if isPositiveMining && referenceBalance != nil && !referenceBalance.Amount.IsZero() {
		referenceBalance.Amount = coin.ZeroICEFlakes()
	}
}

func (s *balanceRecalculationTriggerStreamSource) reverseT0ThisDurationDegradationReferenceTypeDetail(details *BalanceRecalculationDetails, now *time.Time) string { //nolint:lll // .
	return fmt.Sprintf("%v/&%v", details.reverseT0TypeDetail(), s.lastXMiningSessionsCollectingIntervalDateFormat(now))
}

const (
	lastXMiningSessionsReverseT0TypeDetail = reverseT0BalanceTypeDetail + "/0"
)

func (d *BalanceRecalculationDetails) reverseT0TypeDetail() string {
	return fmt.Sprintf("%v_%v", reverseT0BalanceTypeDetail, d.T0UserID)
}

func (d *BalanceRecalculationDetails) reverseT0AggressiveDegradationReferenceTypeDetail() string {
	return fmt.Sprintf("%v_", d.reverseT0TypeDetail())
}

func (d *BalanceRecalculationDetails) reverseT0Changed(typeDetail string) bool {
	if !strings.HasPrefix(typeDetail, reverseT0BalanceTypeDetail+"_") {
		return false
	}
	userID := strings.Replace(typeDetail, reverseT0BalanceTypeDetail+"_", "", 1)
	userID = strings.Replace(userID, "_", "", 1)
	userID = strings.Split(userID, "/")[0]

	return d.T0UserID != userID
}
