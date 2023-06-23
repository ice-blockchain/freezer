// SPDX-License-Identifier: ice License 1.0

package extrabonusnotifier

import (
	"github.com/ice-blockchain/wintr/time"
)

func CalculateExtraBonus(
	newsSeen, extraBonusDaysClaimNotAvailable, extraBonusIndex uint16,
	now, extraBonusLastClaimAvailableAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt *time.Time,
) uint16 {
	const networkDelayDelta = 1.333
	var (
		newsSeenBonusValues            = cfg.ExtraBonuses.NewsSeenValues
		miningStreakValues             = cfg.ExtraBonuses.MiningStreakValues
		bonusPercentageRemaining       = 100 * (1 + extraBonusDaysClaimNotAvailable)
		miningStreak                   = calculateMiningStreak(now, miningSessionSoloStartedAt, miningSessionSoloEndedAt)
		flatBonusValue                 = cfg.ExtraBonuses.FlatValues[extraBonusIndex]
		firstDelayedClaimPenaltyWindow = int64(float64(cfg.ExtraBonuses.DelayedClaimPenaltyWindow.Nanoseconds()) * networkDelayDelta)
	)
	if flatBonusValue == 0 || extraBonusLastClaimAvailableAt.IsNil() {
		return 0
	}
	if delay := now.Sub(*extraBonusLastClaimAvailableAt.Time); delay.Nanoseconds() > firstDelayedClaimPenaltyWindow {
		bonusPercentageRemaining -= 25 * uint16(delay/cfg.ExtraBonuses.DelayedClaimPenaltyWindow)
	}
	if miningStreak >= uint64(len(miningStreakValues)) {
		miningStreak = uint64(len(miningStreakValues) - 1)
	}
	if newsSeen >= uint16(len(newsSeenBonusValues)) {
		newsSeen = uint16(len(newsSeenBonusValues) - 1)
	}

	return ((flatBonusValue + miningStreakValues[miningStreak] + newsSeenBonusValues[newsSeen]) * bonusPercentageRemaining) / 100
}

func calculateMiningStreak(now, start, end *time.Time) uint64 {
	if start.IsNil() || end.IsNil() || now.After(*end.Time) || now.Before(*start.Time) {
		return 0
	}

	return uint64(now.Sub(*start.Time) / cfg.MiningSessionDuration)
}
