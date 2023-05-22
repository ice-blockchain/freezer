// SPDX-License-Identifier: ice License 1.0

package miner

import "github.com/ice-blockchain/wintr/time"

func changeT0AndTMinus1Referrals(usr *user) {
	if usr.IDT0 <= 0 {
		usr.IDT0 *= -1
		usr.IDTMinus1 *= -1
		usr.BalanceT0, usr.BalanceForT0, usr.SlashingRateT0, usr.SlashingRateForT0 = 0, 0, 0, 0
		usr.BalanceForTMinus1, usr.SlashingRateForTMinus1 = 0, 0
		usr.ResurrectT0UsedAt, usr.ResurrectTMinus1UsedAt = new(time.Time), new(time.Time)
	} else if usr.IDTMinus1 <= 0 {
		usr.IDTMinus1 *= -1
		usr.BalanceForTMinus1, usr.SlashingRateForTMinus1 = 0, 0
		usr.ResurrectTMinus1UsedAt = new(time.Time)
	} else {
		usr.IDT0 = 0
		usr.IDTMinus1 = 0
	}
}

func didReferralJustStopMining(now *time.Time, before, after *user) *referralThatStoppedMining {
	if before == nil ||
		after == nil ||
		before.MiningSessionSoloEndedAt.IsNil() ||
		before.BalanceLastUpdatedAt.IsNil() ||
		before.MiningSessionSoloEndedAt.After(*now.Time) ||
		before.BalanceLastUpdatedAt.After(*before.MiningSessionSoloEndedAt.Time) {
		return nil
	}

	return &referralThatStoppedMining{
		ID:              before.ID,
		IDT0:            after.IDT0,
		IDTMinus1:       after.IDTMinus1,
		StoppedMiningAt: before.MiningSessionSoloEndedAt,
	}
}
