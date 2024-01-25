// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/time"
)

func mine(baseMiningRate float64, now *time.Time, usr *user, t0Ref, tMinus1Ref *referral) (updatedUser *user, shouldGenerateHistory, IDT0Changed bool, pendingAmountForTMinus1, pendingAmountForT0 float64) {
	if usr == nil || usr.MiningSessionSoloStartedAt.IsNil() || usr.MiningSessionSoloEndedAt.IsNil() {
		return nil, false, false, 0, 0
	}
	clonedUser1 := *usr
	updatedUser = &clonedUser1
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

		return nil, false, IDT0Changed, 0, 0
	}

	if updatedUser.BalanceLastUpdatedAt.IsNil() {
		updatedUser.BalanceLastUpdatedAt = updatedUser.MiningSessionSoloStartedAt
	} else {
		if updatedUser.BalanceLastUpdatedAt.Year() != now.Year() ||
			updatedUser.BalanceLastUpdatedAt.YearDay() != now.YearDay() ||
			updatedUser.BalanceLastUpdatedAt.Hour() != now.Hour() ||
			(cfg.Development && updatedUser.BalanceLastUpdatedAt.Minute() != now.Minute()) {
			shouldGenerateHistory = true
			updatedUser.BalanceTotalSlashed = 0
			updatedUser.BalanceTotalMinted = 0
			if updatedUser.MiningSessionSoloEndedAt.After(*now.Time) && updatedUser.isAbsoluteZero() {
				usr.BalanceLastUpdatedAt = updatedUser.MiningSessionSoloStartedAt
			}
		}
		if updatedUser.MiningSessionSoloEndedAt.After(*now.Time) && updatedUser.isAbsoluteZero() {
			updatedUser.BalanceLastUpdatedAt = updatedUser.MiningSessionSoloStartedAt
		}
	}

	var (
		mintedAmount        float64
		elapsedTimeFraction float64
		miningSessionRatio  float64
	)
	if timeSpent := now.Sub(*updatedUser.BalanceLastUpdatedAt.Time); cfg.Development {
		elapsedTimeFraction = timeSpent.Minutes()
		miningSessionRatio = 1
	} else {
		elapsedTimeFraction = timeSpent.Hours()
		miningSessionRatio = 24.
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
	needInstantSlashing := usr.WasQuizReset(updatedUser.BalanceLastUpdatedAt)
	if needInstantSlashing {
		pendingResurrectionForTMinus1 = 0
		pendingResurrectionForT0 = 0
		updatedUser.BalanceSolo = 0
		updatedUser.BalanceT1 = 0
		updatedUser.BalanceT2 = 0
		if tMinus1Ref != nil {
			updatedUser.BalanceForTMinus1 = 0
		}
		if t0Ref != nil {
			updatedUser.BalanceForT0 = 0
			updatedUser.BalanceT0 = 0
		}
	}
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

			if updatedUser.SlashingRateForT0 != 0 && (!needInstantSlashing) {
				updatedUser.SlashingRateForT0 = 0
			}
		}
		if tMinus1Ref != nil && !tMinus1Ref.MiningSessionSoloEndedAt.IsNil() && tMinus1Ref.MiningSessionSoloEndedAt.After(*now.Time) {
			updatedUser.BalanceForTMinus1 += 5 * baseMiningRate * elapsedTimeFraction / 100

			if updatedUser.SlashingRateForTMinus1 != 0 && (!needInstantSlashing) {
				updatedUser.SlashingRateForTMinus1 = 0
			}
		}
		if updatedUser.ActiveT1Referrals < 0 {
			updatedUser.ActiveT1Referrals = 0
		}
		if updatedUser.ActiveT2Referrals < 0 {
			updatedUser.ActiveT2Referrals = 0
		}
		t1Rate := (25 * float64(updatedUser.ActiveT1Referrals)) * baseMiningRate * elapsedTimeFraction / 100
		t2Rate := (5 * float64(updatedUser.ActiveT2Referrals)) * baseMiningRate * elapsedTimeFraction / 100
		updatedUser.BalanceT1 += t1Rate
		updatedUser.BalanceT2 += t2Rate
		mintedAmount += t1Rate + t2Rate

	} else {
		if updatedUser.SlashingRateSolo == 0 {
			updatedUser.SlashingRateSolo = updatedUser.BalanceSolo / 60. / miningSessionRatio
		}
		if unAppliedSoloPending != 0 {
			updatedUser.SlashingRateSolo += unAppliedSoloPending / 60. / miningSessionRatio
		}
		if updatedUser.SlashingRateSolo < 0 {
			updatedUser.SlashingRateSolo = 0
		}
	}

	if t0Ref != nil &&
		((!t0Ref.MiningSessionSoloEndedAt.IsNil() && t0Ref.MiningSessionSoloEndedAt.Before(*now.Time)) || updatedUser.MiningSessionSoloEndedAt.Before(*now.Time)) {
		if updatedUser.SlashingRateForT0 == 0 {
			updatedUser.SlashingRateForT0 = updatedUser.BalanceForT0 / 60. / miningSessionRatio
		}
		if updatedUser.SlashingRateT0 == 0 {
			updatedUser.SlashingRateT0 = updatedUser.BalanceT0 / 60. / miningSessionRatio
		}
	}
	if tMinus1Ref != nil &&
		((!tMinus1Ref.MiningSessionSoloEndedAt.IsNil() && tMinus1Ref.MiningSessionSoloEndedAt.Before(*now.Time)) || updatedUser.MiningSessionSoloEndedAt.Before(*now.Time)) {
		if updatedUser.SlashingRateForTMinus1 == 0 {
			updatedUser.SlashingRateForTMinus1 = updatedUser.BalanceForTMinus1 / 60. / miningSessionRatio
		}
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
	}
	if updatedUser.BalanceForTMinus1 < 0 {
		updatedUser.BalanceForTMinus1 = 0
	}

	if usr.BalanceTotalPreStaking+usr.BalanceTotalStandard == 0 {
		slashedAmount = 0
	}

	totalAmount := updatedUser.BalanceSolo + updatedUser.BalanceT0 + updatedUser.BalanceT1 + updatedUser.BalanceT2
	updatedUser.BalanceTotalStandard, updatedUser.BalanceTotalPreStaking = tokenomics.ApplyPreStaking(totalAmount, updatedUser.PreStakingAllocation, updatedUser.PreStakingBonus)
	mintedStandard, mintedPreStaking := tokenomics.ApplyPreStaking(mintedAmount, updatedUser.PreStakingAllocation, updatedUser.PreStakingBonus)
	slashedStandard, slashedPreStaking := tokenomics.ApplyPreStaking(slashedAmount, updatedUser.PreStakingAllocation, updatedUser.PreStakingBonus)
	updatedUser.BalanceTotalMinted += mintedStandard + mintedPreStaking
	updatedUser.BalanceTotalSlashed += slashedStandard + slashedPreStaking
	updatedUser.BalanceLastUpdatedAt = now

	if needInstantSlashing {
		updatedUser.SlashingRateSolo = 0
		updatedUser.SlashingRateT0 = 0
		updatedUser.SlashingRateForT0 = 0
		updatedUser.SlashingRateForTMinus1 = 0
	}

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
		u.BalanceT1 == 0 &&
		u.BalanceT2 == 0 &&
		u.BalanceSoloPending-u.BalanceSoloPendingApplied == 0 &&
		u.BalanceForT0 == 0 &&
		u.BalanceForTMinus1 == 0
}
