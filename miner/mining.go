// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"math"
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/time"
)

func mine(now *time.Time, usr *user, t0Ref, tMinus1Ref *referral) (updatedUser *user, shouldGenerateHistory, IDT0Changed bool, pendingAmountForTMinus1, pendingAmountForT0 float64) {
	if usr == nil || usr.MiningSessionSoloStartedAt.IsNil() || usr.MiningSessionSoloEndedAt.IsNil() {
		return nil, false, false, 0, 0
	}
	clonedUser1 := *usr
	updatedUser = &clonedUser1
	if isTenantInDistributionMode() {
		return updatedUser, false, false, 0, 0
	}

	pendingResurrectionForTMinus1, pendingResurrectionForT0 := resurrect(now, updatedUser, t0Ref, tMinus1Ref)
	IDT0Changed, _ = changeT0AndTMinus1Referrals(updatedUser)
	if updatedUser.MiningSessionSoloEndedAt.Before(*now.Time) && updatedUser.isAbsoluteZero() {
		if updatedUser.BalanceT1Pending-updatedUser.BalanceT1PendingApplied != 0 ||
			updatedUser.BalanceT2Pending-updatedUser.BalanceT2PendingApplied != 0 {
			updatedUser.BalanceT1PendingApplied = updatedUser.BalanceT1Pending
			updatedUser.BalanceT2PendingApplied = updatedUser.BalanceT2Pending
			updatedUser.BalanceLastUpdatedAt = now

			return updatedUser, false, IDT0Changed, 0, 0
		}
		if updatedUser.BalanceT1 > 0 || updatedUser.BalanceT2 > 0 {
			updatedUser.BalanceTotalStandard, updatedUser.BalanceTotalPreStaking = 0, 0
			updatedUser.BalanceT1 = 0
			updatedUser.BalanceT2 = 0
			updatedUser.BalanceLastUpdatedAt = now

			return updatedUser, false, IDT0Changed, 0, 0
		}

		return nil, false, IDT0Changed, 0, 0
	}

	var miningSessionRatio float64
	var miningPeriod stdlibtime.Duration
	if cfg.Development {
		miningPeriod = 1 * stdlibtime.Minute
		miningSessionRatio = 1
	} else {
		miningPeriod = 1 * stdlibtime.Hour
		miningSessionRatio = 24.
	}
	if updatedUser.MiningSessionSoloEndedAt.Before(*now.Time) && (updatedUser.reachedSlashingFloor() || updatedUser.slashingDisabled()) {
		fullSlashingDuration := stdlibtime.Duration(cfg.SlashingDaysCount * int64(miningSessionRatio) * int64(miningPeriod))
		shouldGenerateHistory = !updatedUser.BalanceLastUpdatedAt.IsNil() && (updatedUser.BalanceLastUpdatedAt.Year() != now.Year() ||
			updatedUser.BalanceLastUpdatedAt.YearDay() != now.YearDay() ||
			(cfg.Development && updatedUser.BalanceLastUpdatedAt.Minute() != now.Minute())) &&
			((updatedUser.slashingDisabled() && updatedUser.BalanceLastUpdatedAt.After(*updatedUser.MiningSessionSoloEndedAt.Time) &&
				updatedUser.BalanceLastUpdatedAt.Sub(*updatedUser.MiningSessionSoloEndedAt.Time) < cfg.MiningSessionDuration.Min) ||
				(updatedUser.reachedSlashingFloor() &&
					now.After(updatedUser.MiningSessionSoloEndedAt.Time.Add(fullSlashingDuration)) &&
					now.Sub(*updatedUser.MiningSessionSoloEndedAt.Time) < stdlibtime.Duration((cfg.SlashingDaysCount+1)*int64(miningSessionRatio)*int64(miningPeriod))))

		if shouldGenerateHistory {
			updatedUser.BalanceTotalSlashed = 0
			updatedUser.BalanceTotalMinted = 0
		}

		if shouldGenerateHistory ||
			(!updatedUser.ReferralsCountChangeGuardUpdatedAt.IsNil() &&
				!updatedUser.MiningSessionSoloStartedAt.IsNil() &&
				updatedUser.ReferralsCountChangeGuardUpdatedAt.Equal(*updatedUser.MiningSessionSoloStartedAt.Time)) {
			updatedUser.BalanceLastUpdatedAt = now
			// We need to update ReferralsCountChangeGuardUpdatedAt last time to avoid ErrDuplicate on next sessions
			return updatedUser, shouldGenerateHistory, IDT0Changed, 0, 0
		}

		return nil, shouldGenerateHistory, IDT0Changed, 0, 0
	}

	if updatedUser.BalanceLastUpdatedAt.IsNil() {
		updatedUser.BalanceLastUpdatedAt = updatedUser.MiningSessionSoloStartedAt
	} else {
		if updatedUser.BalanceLastUpdatedAt.Year() != now.Year() ||
			updatedUser.BalanceLastUpdatedAt.YearDay() != now.YearDay() ||
			(cfg.Development && updatedUser.BalanceLastUpdatedAt.Minute() != now.Minute()) {
			shouldGenerateHistory = true
			updatedUser.BalanceTotalSlashed = 0
			updatedUser.BalanceTotalMinted = 0
		}
		userWasInSlashingAndReachedFloor := (updatedUser.reachedSlashingFloor() && usr.SlashingRateSolo > 0)
		if updatedUser.MiningSessionSoloEndedAt.After(*now.Time) && (updatedUser.isAbsoluteZero() || userWasInSlashingAndReachedFloor) {
			updatedUser.BalanceLastUpdatedAt = updatedUser.MiningSessionSoloStartedAt
		}
	}

	var (
		mintedAmount        float64
		elapsedTimeFraction float64
	)
	if timeSpent := now.Sub(*updatedUser.BalanceLastUpdatedAt.Time); cfg.Development {
		elapsedTimeFraction = timeSpent.Minutes()
	} else {
		elapsedTimeFraction = timeSpent.Hours()
	}

	unAppliedSoloPending := updatedUser.BalanceSoloPending - updatedUser.BalanceSoloPendingApplied
	unAppliedT1Pending := updatedUser.BalanceT1Pending - updatedUser.BalanceT1PendingApplied
	unAppliedT2Pending := updatedUser.BalanceT2Pending - updatedUser.BalanceT2PendingApplied
	updatedUser.BalanceSoloPendingApplied = updatedUser.BalanceSoloPending
	updatedUser.BalanceT1PendingApplied = updatedUser.BalanceT1Pending
	updatedUser.BalanceT2PendingApplied = updatedUser.BalanceT2Pending
	if unAppliedSoloPending == 0 {
		updatedUser.BalanceSoloPending = 0
		updatedUser.BalanceSoloPendingApplied = 0
	}
	if unAppliedT1Pending == 0 {
		updatedUser.BalanceT1Pending = 0
		updatedUser.BalanceT1PendingApplied = 0
	}
	if unAppliedT2Pending == 0 {
		updatedUser.BalanceT2Pending = 0
		updatedUser.BalanceT2PendingApplied = 0
	}

	baseMiningRate := updatedUser.baseMiningRate(now)
	maxT1Referrals := (*cfg.miningBoostLevels.Load())[len(*cfg.miningBoostLevels.Load())-1].MaxT1Referrals
	if updatedUser.MiningSessionSoloEndedAt.After(*now.Time) {
		if !updatedUser.ExtraBonusStartedAt.IsNil() && now.Before(updatedUser.ExtraBonusStartedAt.Add(cfg.ExtraBonuses.Duration)) {
			rate := (100 + float64(updatedUser.ExtraBonus)) * baseMiningRate * elapsedTimeFraction / 100.
			updatedUser.BalanceSolo += rate
			mintedAmount += rate
		} else {
			rate := baseMiningRate * elapsedTimeFraction
			updatedUser.BalanceSolo += rate
			mintedAmount += rate
		}
		if t0Ref != nil && !t0Ref.MiningSessionSoloEndedAt.IsNil() && t0Ref.MiningSessionSoloEndedAt.After(*now.Time) {
			rate := 25 * baseMiningRate * elapsedTimeFraction / 100
			updatedUser.BalanceForT0 += rate
			updatedUser.BalanceT0 += rate
			mintedAmount += rate

			if updatedUser.SlashingRateForT0 != 0 {
				updatedUser.SlashingRateForT0 = 0
			}
		}
		if tMinus1Ref != nil && !tMinus1Ref.MiningSessionSoloEndedAt.IsNil() && tMinus1Ref.MiningSessionSoloEndedAt.After(*now.Time) {
			updatedUser.BalanceForTMinus1 += 5 * baseMiningRate * elapsedTimeFraction / 100

			if updatedUser.SlashingRateForTMinus1 != 0 {
				updatedUser.SlashingRateForTMinus1 = 0
			}
		}
		if updatedUser.ActiveT1Referrals < 0 {
			updatedUser.ActiveT1Referrals = 0
		}
		if updatedUser.ActiveT2Referrals < 0 {
			updatedUser.ActiveT2Referrals = 0
		}
		activeT1Referrals := int32(0)
		if updatedUser.MiningBoostLevelIndex != nil {
			if updatedUser.IsVerified() && updatedUser.VerifiedT1Referrals >= maxT1Referrals {
				activeT1Referrals = int32((*cfg.miningBoostLevels.Load())[int(*updatedUser.MiningBoostLevelIndex)].MaxT1Referrals)
			} else {
				activeT1Referrals = int32(math.Min(float64((*cfg.miningBoostLevels.Load())[int(*updatedUser.MiningBoostLevelIndex)].MaxT1Referrals), float64(updatedUser.ActiveT1Referrals)))
			}
		} else {
			if cfg.T1ReferralsAllowedWithoutAnyMiningBoostLevel {
				activeT1Referrals = updatedUser.ActiveT1Referrals
			}
		}
		t1Rate := (25 * float64(activeT1Referrals)) * baseMiningRate * elapsedTimeFraction / 100
		t2Rate := (5 * float64(updatedUser.ActiveT2Referrals)) * baseMiningRate * elapsedTimeFraction / 100
		updatedUser.BalanceT1 += t1Rate
		updatedUser.BalanceT2 += t2Rate
		mintedAmount += t1Rate + t2Rate

	} else {
		if !updatedUser.slashingDisabled() {
			if updatedUser.SlashingRateSolo == 0 {
				updatedUser.SlashingRateSolo = updatedUser.BalanceSolo / float64(cfg.SlashingDaysCount) / miningSessionRatio
			}
			if unAppliedSoloPending != 0 {
				updatedUser.SlashingRateSolo += unAppliedSoloPending / float64(cfg.SlashingDaysCount) / miningSessionRatio
			}
			if updatedUser.SlashingRateSolo < 0 {
				updatedUser.SlashingRateSolo = 0
			}
		}
	}

	if t0Ref != nil {
		if updatedUser.SlashingRateForT0 == 0 && !t0Ref.MiningSessionSoloEndedAt.IsNil() && t0Ref.MiningSessionSoloEndedAt.Before(*now.Time) && !t0Ref.slashingDisabled() && !t0Ref.reachedSlashingFloor() {
			updatedUser.SlashingRateForT0 = updatedUser.BalanceForT0 / float64(cfg.SlashingDaysCount) / miningSessionRatio
		}
		if updatedUser.SlashingRateT0 == 0 && !updatedUser.MiningSessionSoloEndedAt.IsNil() && updatedUser.MiningSessionSoloEndedAt.Before(*now.Time) && !updatedUser.slashingDisabled() && !updatedUser.reachedSlashingFloor() {
			updatedUser.SlashingRateT0 = updatedUser.BalanceT0 / float64(cfg.SlashingDaysCount) / miningSessionRatio
		}
	}
	if tMinus1Ref != nil {
		if updatedUser.SlashingRateForTMinus1 == 0 && !tMinus1Ref.MiningSessionSoloEndedAt.IsNil() && tMinus1Ref.MiningSessionSoloEndedAt.Before(*now.Time) && !tMinus1Ref.slashingDisabled() && !tMinus1Ref.reachedSlashingFloor() {
			updatedUser.SlashingRateForTMinus1 = updatedUser.BalanceForTMinus1 / float64(cfg.SlashingDaysCount) / miningSessionRatio
		}
	}
	if updatedUser.BalanceT1WelcomeBonusPendingApplied < float64(min(maxT1Referrals, uint64(usr.TotalT1Referrals)))*cfg.WelcomeBonusV2Amount {
		if unAppliedT1WelcomeBonusPending := updatedUser.BalanceT1WelcomeBonusPending - updatedUser.BalanceT1WelcomeBonusPendingApplied; unAppliedT1WelcomeBonusPending == 0 {
			updatedUser.BalanceT1WelcomeBonusPending = 0
			updatedUser.BalanceT1WelcomeBonusPendingApplied = 0
		} else {
			unAppliedT1Pending += min(unAppliedT1WelcomeBonusPending, float64(min(maxT1Referrals, uint64(usr.TotalT1Referrals)))*cfg.WelcomeBonusV2Amount-updatedUser.BalanceT1WelcomeBonusPendingApplied)
			updatedUser.BalanceT1WelcomeBonusPendingApplied = min(updatedUser.BalanceT1WelcomeBonusPending, float64(min(maxT1Referrals, uint64(usr.TotalT1Referrals)))*cfg.WelcomeBonusV2Amount)
		}
	} else {
		updatedUser.BalanceT1WelcomeBonusPending = 0
		updatedUser.BalanceT1WelcomeBonusPendingApplied = 0
	}

	slashedAmount := (updatedUser.SlashingRateSolo + updatedUser.SlashingRateT0) * elapsedTimeFraction
	updatedUser.BalanceSolo -= updatedUser.SlashingRateSolo * elapsedTimeFraction

	pendingAmountForTMinus1 -= updatedUser.SlashingRateForTMinus1 * elapsedTimeFraction
	pendingAmountForT0 -= updatedUser.SlashingRateForT0 * elapsedTimeFraction

	updatedUser.BalanceForTMinus1 += pendingAmountForTMinus1
	updatedUser.BalanceForT0 += pendingAmountForT0
	updatedUser.BalanceT0 -= updatedUser.SlashingRateT0 * elapsedTimeFraction
	updatedUser.BalanceSolo += unAppliedSoloPending
	updatedUser.BalanceT1 += unAppliedT1Pending
	updatedUser.BalanceT2 += unAppliedT2Pending

	pendingAmountForTMinus1 += pendingResurrectionForTMinus1
	pendingAmountForT0 += pendingResurrectionForT0

	if unAppliedSoloPending < 0 {
		slashedAmount += -unAppliedSoloPending
	} else {
		mintedAmount += unAppliedSoloPending
	}
	if unAppliedT1Pending < 0 {
		slashedAmount += -unAppliedT1Pending
	} else {
		mintedAmount += unAppliedT1Pending
	}
	if unAppliedT2Pending < 0 {
		slashedAmount += -unAppliedT2Pending
	} else {
		mintedAmount += unAppliedT2Pending
	}
	if updatedUser.BalanceSolo < 0 {
		updatedUser.BalanceSolo = 0
	}
	if updatedUser.BalanceT0 < 0 {
		updatedUser.BalanceT0 = 0
	}
	if updatedUser.BalanceT1 < 0 {
		updatedUser.BalanceT1 = 0
	}
	if updatedUser.BalanceT2 < 0 {
		updatedUser.BalanceT2 = 0
	}
	if updatedUser.BalanceForT0 < 0 {
		updatedUser.BalanceForT0 = 0
		pendingAmountForT0 = 0
	}
	if updatedUser.BalanceForTMinus1 < 0 {
		updatedUser.BalanceForTMinus1 = 0
		pendingAmountForTMinus1 = 0
	}

	if usr.BalanceTotalPreStaking+usr.BalanceTotalStandard == 0 {
		slashedAmount = 0
	}
	if updatedUser.WelcomeBonusV2Applied == nil || !*updatedUser.WelcomeBonusV2Applied {
		updatedUser.BalanceSolo += cfg.WelcomeBonusV2Amount - 10
		trueVal := model.FlexibleBool(true)
		updatedUser.WelcomeBonusV2Applied = &trueVal
	} else {
		updatedUser.WelcomeBonusV2Applied = nil
	}

	totalAmount := updatedUser.BalanceSolo + updatedUser.BalanceT0 + updatedUser.BalanceT1 + updatedUser.BalanceT2
	updatedUser.BalanceTotalStandard, updatedUser.BalanceTotalPreStaking = tokenomics.ApplyPreStaking(totalAmount, updatedUser.PreStakingAllocation, updatedUser.PreStakingBonus)
	mintedStandard, mintedPreStaking := tokenomics.ApplyPreStaking(mintedAmount, updatedUser.PreStakingAllocation, updatedUser.PreStakingBonus)
	slashedStandard, slashedPreStaking := tokenomics.ApplyPreStaking(slashedAmount, updatedUser.PreStakingAllocation, updatedUser.PreStakingBonus)
	updatedUser.BalanceTotalMinted += mintedStandard + mintedPreStaking
	updatedUser.BalanceTotalSlashed += slashedStandard + slashedPreStaking
	updatedUser.BalanceLastUpdatedAt = now

	return updatedUser, shouldGenerateHistory, IDT0Changed, pendingAmountForTMinus1, pendingAmountForT0
}

func updateT0AndTMinus1ReferralsForUserHasNeverMined(usr *user) (updatedUser *referralUpdated) {
	if usr.IDT0 < 0 && (usr.MiningSessionSoloLastStartedAt.IsNil() || usr.MiningSessionSoloEndedAt.IsNil()) &&
		usr.BalanceLastUpdatedAt.IsNil() {
		if IDT0Changed, _ := changeT0AndTMinus1Referrals(usr); IDT0Changed {
			return &referralUpdated{
				DeserializedUsersKey: usr.DeserializedUsersKey,
				IDT0Field:            usr.IDT0Field,
				IDTMinus1Field:       usr.IDTMinus1Field,
			}
		}
	}

	return nil
}

func (u *user) isAbsoluteZero() bool {
	return u.BalanceSolo == 0 &&
		u.BalanceT0 == 0 &&
		u.BalanceSoloPending-u.BalanceSoloPendingApplied == 0 &&
		u.BalanceForT0 == 0 &&
		u.BalanceForTMinus1 == 0
}

func (u *user) reachedSlashingFloor() bool {
	return (u.BalanceSolo + u.BalanceT0 + u.BalanceT1 + u.BalanceT2) <= cfg.SlashingFloor
}

func (ref *referral) reachedSlashingFloor() bool {
	return (ref.BalanceSolo + ref.BalanceT0 + ref.BalanceT1 + ref.BalanceT2) <= cfg.SlashingFloor
}

func (u *user) slashingDisabled() bool {
	if u == nil || u.MiningBoostLevelIndex == nil {
		return false
	}

	return (*cfg.miningBoostLevels.Load())[*u.MiningBoostLevelIndex].SlashingDisabled
}

func (ref *referral) slashingDisabled() bool {
	if ref == nil || ref.MiningBoostLevelIndex == nil {
		return false
	}

	return (*cfg.miningBoostLevels.Load())[*ref.MiningBoostLevelIndex].SlashingDisabled
}

func (u *user) baseMiningRate(now *time.Time) float64 {
	if u == nil {
		return 0
	}

	return cfg.BaseMiningRate(now, u.CreatedAt)
}
