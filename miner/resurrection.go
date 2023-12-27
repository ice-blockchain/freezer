// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/time"
)

func resurrect(now *time.Time, usr *user, t0Ref, tMinus1Ref *referral) (pendingResurrectionForTMinus1, pendingResurrectionForT0 float64) {
	if !usr.ResurrectSoloUsedAt.IsNil() && usr.ResurrectSoloUsedAt.After(*now.Time) {
		var resurrectDelta float64
		if timeSpent := usr.MiningSessionSoloStartedAt.Sub(*usr.MiningSessionSoloPreviouslyEndedAt.Time); cfg.Development {
			resurrectDelta = timeSpent.Minutes()
		} else {
			resurrectDelta = timeSpent.Hours()
		}

		usr.BalanceSolo += usr.SlashingRateSolo * resurrectDelta
		usr.BalanceT0 += usr.SlashingRateT0 * resurrectDelta
		mintedAmount := (usr.SlashingRateSolo + usr.SlashingRateT0) * resurrectDelta
		mintedStandard, mintedPreStaking := tokenomics.ApplyPreStaking(mintedAmount, usr.PreStakingAllocation, usr.PreStakingBonus)
		usr.BalanceTotalMinted += mintedStandard + mintedPreStaking

		usr.SlashingRateSolo, usr.SlashingRateT0 = 0, 0
		usr.ResurrectSoloUsedAt = now
	} else {
		usr.ResurrectSoloUsedAt = nil
	}

	if t0Ref != nil && !t0Ref.ResurrectSoloUsedAt.IsNil() && usr.ResurrectT0UsedAt.IsNil() {
		var resurrectDelta float64
		if timeSpent := t0Ref.MiningSessionSoloStartedAt.Sub(*t0Ref.MiningSessionSoloPreviouslyEndedAt.Time); cfg.Development {
			resurrectDelta = timeSpent.Minutes()
		} else {
			resurrectDelta = timeSpent.Hours()
		}

		amount := usr.SlashingRateForT0 * resurrectDelta
		usr.BalanceForT0 += amount
		pendingResurrectionForT0 += amount

		usr.SlashingRateForT0 = 0
		usr.ResurrectT0UsedAt = now
	} else {
		usr.ResurrectT0UsedAt = nil
	}

	if tMinus1Ref != nil && !tMinus1Ref.ResurrectSoloUsedAt.IsNil() && usr.ResurrectTMinus1UsedAt.IsNil() {
		var resurrectDelta float64
		if timeSpent := tMinus1Ref.MiningSessionSoloStartedAt.Sub(*tMinus1Ref.MiningSessionSoloPreviouslyEndedAt.Time); cfg.Development {
			resurrectDelta = timeSpent.Minutes()
		} else {
			resurrectDelta = timeSpent.Hours()
		}

		amount := usr.SlashingRateForTMinus1 * resurrectDelta
		usr.BalanceForTMinus1 += amount
		pendingResurrectionForTMinus1 += amount

		usr.SlashingRateForTMinus1 = 0
		usr.ResurrectTMinus1UsedAt = now
	} else {
		usr.ResurrectTMinus1UsedAt = nil
	}

	if usr.MiningSessionSoloEndedAt.After(*now.Time) {
		usr.SlashingRateSolo, usr.SlashingRateT0 = 0, 0
	}
	if usr.SlashingRateForT0 > 0 && (t0Ref == nil || t0Ref.MiningSessionSoloEndedAt.IsNil() || (t0Ref.MiningSessionSoloEndedAt.After(*now.Time) && usr.MiningSessionSoloEndedAt.After(*now.Time))) {
		usr.SlashingRateForT0 = 0
	}

	if usr.SlashingRateForTMinus1 > 0 && (tMinus1Ref == nil || tMinus1Ref.MiningSessionSoloEndedAt.IsNil() || (tMinus1Ref.MiningSessionSoloEndedAt.After(*now.Time) && usr.MiningSessionSoloEndedAt.After(*now.Time))) { //nolint:lll // .
		usr.SlashingRateForTMinus1 = 0
	}

	return pendingResurrectionForTMinus1, pendingResurrectionForT0
}
