// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/time"
)

func (s *balanceRecalculationTriggerStreamSource) processPreviousIncompleteT0TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration,
) {
	if elapsedDuration == 0 ||
		details.T0UserID == details.UUserID ||
		details.T0UserID == "" {
		return
	}
	if isPositiveMining := details.LastMiningEndedAt.After(*now.Time); isPositiveMining { // This means that the previous one was negative.
		s.slashT0TotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, true)
	} else {
		s.mintT0TotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, true)
	}
}

func (s *balanceRecalculationTriggerStreamSource) rollbackT0TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	if details.T0UserID == details.UUserID ||
		details.T0UserID == "" ||
		details.LastMiningEndedAt.Before(*now.Time) {
		return
	}
	negativeBalance := s.getBalance(true, details.t0TypeDetail(), balancesByPK)
	if negativeBalance == nil || negativeBalance.Amount.IsZero() {
		return
	}
	if details.RollbackUsedAt != nil {
		positiveBalance := s.getOrInitBalance(false, details.t0TypeDetail(), details.UUserID, balancesByPK)
		positiveBalance.add(negativeBalance.Amount)
	}
	negativeBalance.Amount = coin.ZeroICEFlakes()
}

func (s *balanceRecalculationTriggerStreamSource) processT0TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance,
	details *BalanceRecalculationDetails,
	now *time.Time,
	elapsedDuration stdlibtime.Duration,
) {
	if details.T0UserID == details.UUserID ||
		details.T0UserID == "" {
		return
	}
	defer func() {
		s.getBalance(false, lastXMiningSessionsT0TypeDetail, balancesByPK).Amount = nil
	}()
	if isPositiveMining := details.LastMiningEndedAt.After(*now.Time); isPositiveMining {
		s.mintT0TotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, false)
	} else {
		s.slashT0TotalNoPreStakingBonusBalanceType(balancesByPK, details, now, elapsedDuration, false)
	}
}

func (s *balanceRecalculationTriggerStreamSource) mintT0TotalNoPreStakingBonusBalanceType( //nolint:revive // Nope.
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, previous bool,
) {
	if details.T0 == 0 {
		return
	}
	params := &userMiningRateRecalculationParameters{T0: details.T0}
	mintedAmount := s.calculateMintedStandardCoins(details.BaseMiningRate, params, elapsedDuration, true)
	positiveBalance := s.getOrInitBalance(false, details.t0TypeDetail(), details.UUserID, balancesByPK)
	positiveBalance.add(mintedAmount)
	positiveTotalThisMiningSessionBalance := s.getOrInitBalance(false, s.t0ThisDurationDegradationReferenceTypeDetail(details, now), details.UUserID, balancesByPK)
	positiveTotalThisMiningSessionBalance.add(mintedAmount)
	if previous {
		degradationReference := s.getOrInitBalance(false, degradationT0T1T2TotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
		degradationReference.add(s.getBalance(false, lastXMiningSessionsT0TypeDetail, balancesByPK).Amount)
		degradationReference.add(mintedAmount)
	}
}

//nolint:revive // Prefer decoupling.
func (s *balanceRecalculationTriggerStreamSource) slashT0TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time, elapsedDuration stdlibtime.Duration, previous bool,
) {
	positiveBalance := s.getOrInitBalance(false, details.t0TypeDetail(), details.UUserID, balancesByPK)
	if positiveBalance.Amount.IsZero() {
		return
	}
	aggressive := details.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	var referenceAmount *coin.ICEFlake
	if aggressive {
		referenceAmount = s.getBalance(false, details.t0AggressiveDegradationReferenceTypeDetail(), balancesByPK).Amount
	} else {
		referenceAmount = s.getBalance(false, lastXMiningSessionsT0TypeDetail, balancesByPK).Amount
	}
	negativeThisDuration := s.getOrInitBalance(true, s.lastXMiningSessionsThisDurationTypeDetail(previous), details.UUserID, balancesByPK)
	slashedAmount := s.calculateDegradation(elapsedDuration, referenceAmount, aggressive)
	positiveBalance.subtract(slashedAmount)
	negativeThisDuration.add(slashedAmount)
	if details.RollbackUsedAt == nil || (previous && details.RollbackUsedAt.Equal(*details.LastMiningStartedAt.Time)) {
		negativeBalance := s.getOrInitBalance(true, details.t0TypeDetail(), details.UUserID, balancesByPK)
		negativeBalance.add(slashedAmount)
	}
}

func (s *balanceRecalculationTriggerStreamSource) processDegradationForT0TotalNoPreStakingBonusBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, now *time.Time,
) {
	if details.T0UserID == details.UUserID ||
		details.T0UserID == "" {
		return
	}
	isPositiveMining := details.LastMiningEndedAt.After(*now.Time)
	s.processLastXPositiveMiningSessions(balancesByPK, false, details.t0TypeDetail()+"/&", lastXMiningSessionsT0TypeDetail, details.UUserID)
	if isPositiveMining {
		degradationReference := s.getOrInitBalance(false, degradationT0T1T2TotalReferenceBalanceTypeDetail, details.UUserID, balancesByPK)
		degradationReference.add(s.getBalance(false, lastXMiningSessionsT0TypeDetail, balancesByPK).Amount)
	}

	aggressive := details.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time)
	referenceBalance := s.getBalance(false, details.t0AggressiveDegradationReferenceTypeDetail(), balancesByPK)
	if !isPositiveMining && aggressive && (referenceBalance == nil || referenceBalance.Amount.IsNil()) {
		positiveBalance := s.getOrInitBalance(false, details.t0TypeDetail(), details.UUserID, balancesByPK)
		referenceBalance = s.getOrInitBalance(false, details.t0AggressiveDegradationReferenceTypeDetail(), details.UUserID, balancesByPK)
		referenceBalance.add(positiveBalance.Amount)
	}
	if isPositiveMining && referenceBalance != nil && !referenceBalance.Amount.IsZero() {
		referenceBalance.Amount = coin.ZeroICEFlakes()
	}
}

func (s *balanceRecalculationTriggerStreamSource) t0ThisDurationDegradationReferenceTypeDetail(details *BalanceRecalculationDetails, now *time.Time) string { //nolint:lll // .
	return fmt.Sprintf("%v/&%v", details.t0TypeDetail(), s.lastXMiningSessionsCollectingIntervalDateFormat(now))
}

const (
	lastXMiningSessionsT0TypeDetail = t0BalanceTypeDetail + "/0"
)

func (d *BalanceRecalculationDetails) t0TypeDetail() string {
	return fmt.Sprintf("%v_%v", t0BalanceTypeDetail, d.T0UserID)
}

func (d *BalanceRecalculationDetails) t0AggressiveDegradationReferenceTypeDetail() string {
	return fmt.Sprintf("%v_", d.t0TypeDetail())
}

func (d *BalanceRecalculationDetails) t0Changed(typeDetail string) bool {
	if !strings.HasPrefix(typeDetail, t0BalanceTypeDetail+"_") {
		return false
	}
	userID := strings.Replace(typeDetail, t0BalanceTypeDetail+"_", "", 1)
	userID = strings.Replace(userID, "_", "", 1)
	userID = strings.Split(userID, "/")[0]

	return d.T0UserID != userID
}
