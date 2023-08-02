// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"github.com/ice-blockchain/wintr/time"
)

func changeT0AndTMinus1Referrals(usr *user, reset bool) (IDT0Changed, IDTMinus1Changed bool) {
	if usr.IDT0 <= 0 {
		usr.IDT0 *= -1
		usr.IDTMinus1 *= -1
		if reset {
			usr.BalanceT0, usr.BalanceForT0, usr.SlashingRateT0, usr.SlashingRateForT0 = 0, 0, 0, 0
			usr.BalanceForTMinus1, usr.SlashingRateForTMinus1 = 0, 0
			usr.ResurrectT0UsedAt, usr.ResurrectTMinus1UsedAt = new(time.Time), new(time.Time)
		}
		IDT0Changed = true
	} else if usr.IDTMinus1 <= 0 {
		usr.IDTMinus1 *= -1
		if reset {
			usr.BalanceForTMinus1, usr.SlashingRateForTMinus1 = 0, 0
			usr.ResurrectTMinus1UsedAt = new(time.Time)
		}
		IDTMinus1Changed = true
	} else {
		usr.IDT0 = 0
		usr.IDTMinus1 = 0
	}

	return IDT0Changed, IDTMinus1Changed
}

func didReferralJustStopMining(now *time.Time, before *user, t0Ref, tMinus1Ref *referral) *referralThatStoppedMining {
	if before == nil ||
		before.MiningSessionSoloEndedAt.IsNil() ||
		before.BalanceLastUpdatedAt.IsNil() ||
		before.MiningSessionSoloEndedAt.After(*now.Time) ||
		before.BalanceLastUpdatedAt.After(*before.MiningSessionSoloEndedAt.Time) {
		return nil
	}
	var idT0, idTminus1 int64
	if t0Ref != nil {
		idT0 = t0Ref.ID
	}
	if tMinus1Ref != nil {
		idTminus1 = tMinus1Ref.ID
	}

	return &referralThatStoppedMining{
		ID:              before.ID,
		IDT0:            idT0,
		IDTMinus1:       idTminus1,
		StoppedMiningAt: before.MiningSessionSoloEndedAt,
	}
}
