// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"github.com/ice-blockchain/wintr/time"
)

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

func wasReferredByUpdated(usr *user) *activeUserWithUpdatedReferredBy {
	if (usr.IDT0Old == 0 && usr.IDTMinus1Old == 0) ||
		usr.IDT0Old == usr.IDT0 || usr.IDT0Old == usr.IDT0*-1 ||
		usr.IDTMinus1 == usr.IDTMinus1Old || usr.IDTMinus1 == usr.IDTMinus1Old*-1 ||
		usr.MiningSessionSoloEndedAt.IsNil() || usr.MiningSessionSoloEndedAt.Before(*time.Now().Time) {
		return nil
	}

	return &activeUserWithUpdatedReferredBy{
		ID:            usr.ID,
		oldIDT0:       usr.IDT0Old,
		oldIDTMinus1:  usr.IDTMinus1Old,
		newIDT0:       usr.IDT0,
		newIDTMinus1:  usr.IDTMinus1,
		ActiveT1Count: usr.ActiveT1Referrals,
	}
}
