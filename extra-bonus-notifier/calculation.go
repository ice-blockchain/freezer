// SPDX-License-Identifier: ice License 1.0

package extrabonusnotifier

import (
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/time"
)

func CalculateExtraBonus(
	newsSeen, extraBonusDaysClaimNotAvailable, extraBonusIndex uint16,
	now, extraBonusLastClaimAvailableAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt *time.Time,
) float64 {
	const networkDelayDelta = 1.333
	var (
		newsSeenBonusValues            = cfg.ExtraBonuses.NewsSeenValues
		miningStreakValues             = cfg.ExtraBonuses.MiningStreakValues
		bonusPercentageRemaining       = float64(100 * (1 + extraBonusDaysClaimNotAvailable))
		miningStreak                   = model.CalculateMiningStreak(now, miningSessionSoloStartedAt, miningSessionSoloEndedAt, cfg.MiningSessionDuration)
		flatBonusValue                 = cfg.ExtraBonuses.FlatValues[extraBonusIndex]
		firstDelayedClaimPenaltyWindow = int64(float64(cfg.ExtraBonuses.DelayedClaimPenaltyWindow.Nanoseconds()) * networkDelayDelta)
	)
	if flatBonusValue == 0 || extraBonusLastClaimAvailableAt.IsNil() {
		return 0
	}
	if delay := now.Sub(*extraBonusLastClaimAvailableAt.Time); delay.Nanoseconds() > firstDelayedClaimPenaltyWindow {
		bonusPercentageRemaining -= 25 * float64(delay/cfg.ExtraBonuses.DelayedClaimPenaltyWindow)
	}
	if miningStreak >= uint64(len(miningStreakValues)) {
		miningStreak = uint64(len(miningStreakValues) - 1)
	}
	if newsSeen >= uint16(len(newsSeenBonusValues)) {
		newsSeen = uint16(len(newsSeenBonusValues) - 1)
	}

	return (float64(flatBonusValue+miningStreakValues[miningStreak]+newsSeenBonusValues[newsSeen]) * bonusPercentageRemaining) / 100
}
