// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/wintr/time"
)

const (
	testMiningBase = 16
	testIDT0       = 42
	testIDTMinus1  = 69
)

var (
	testTime = time.New(stdlibtime.Date(2023, 1, 2, 3, 4, 5, 6, stdlibtime.UTC))
)

func newUser() *user {
	u := new(user)
	u.UserID = "test_user_id"
	u.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour)
	u.MiningSessionSoloEndedAt = timeDelta(23 * stdlibtime.Hour)

	return u
}

func newRef() *referral {
	r := new(referral)
	r.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour)
	r.MiningSessionSoloEndedAt = timeDelta(23 * stdlibtime.Hour)

	return r
}

func timeDelta(d stdlibtime.Duration) *time.Time {
	return time.New(testTime.Add(d))
}

func testSoloMiningNoExtraBonus(t *testing.T) {
	t.Run("No referrals", func(t *testing.T) {
		m := newUser()

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, testMiningBase, m.BalanceSolo)
	})
	t.Run("With T0", func(t *testing.T) {
		m := newUser()
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
	})

	t.Run("For tMinus1", func(t *testing.T) {
		m := newUser()
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, nil, ref)
		require.NotNil(t, m)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 0.8, m.BalanceForTMinus1)
	})

	t.Run("With T1", func(t *testing.T) {
		m := newUser()
		m.ActiveT1Referrals = 4
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, nil, ref)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 16, m.BalanceT1)
	})
	t.Run("With T0 + T1", func(t *testing.T) {
		m := newUser()
		m.ActiveT1Referrals = 4
		ref := newRef()
		tMinus1Ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, tMinus1Ref)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
		require.EqualValues(t, 16, m.BalanceT1)
	})

	t.Run("With T2", func(t *testing.T) {
		m := newUser()
		m.ActiveT2Referrals = 20

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 16, m.BalanceT2)
	})
	t.Run("With T0 + T1 + tMinus1 + T2", func(t *testing.T) {
		m := newUser()
		m.ActiveT1Referrals = 4
		m.ActiveT2Referrals = 20
		ref := newRef()
		tMinus1Ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, tMinus1Ref)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
		require.EqualValues(t, 16, m.BalanceT1)
		require.EqualValues(t, 16, m.BalanceT2)
		require.EqualValues(t, 0.8, m.BalanceForTMinus1)
	})
}

func testSoloMiningWithExtraBonus(t *testing.T) {
	t.Run("No referrals", func(t *testing.T) {
		m := newUser()
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		m.ExtraBonus = 100

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
	})

	t.Run("With T0", func(t *testing.T) {
		m := newUser()
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		m.ExtraBonus = 100
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
	})

	t.Run("For tMinus1", func(t *testing.T) {
		m := newUser()
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		m.ExtraBonus = 100
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, nil, ref)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 0.8, m.BalanceForTMinus1)
	})

	t.Run("With T1", func(t *testing.T) {
		m := newUser()
		m.ActiveT1Referrals = 4
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		m.ExtraBonus = 100

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 16, m.BalanceT1)
	})
	t.Run("With T0 + T1", func(t *testing.T) {
		m := newUser()
		ref := newRef()
		m.ActiveT1Referrals = 4
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		m.ExtraBonus = 100

		m, _ = mine(testMiningBase, testTime, m, ref, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
		require.EqualValues(t, 16, m.BalanceT1)
	})
	t.Run("With T2", func(t *testing.T) {
		m := newUser()
		m.ActiveT2Referrals = 20
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		m.ExtraBonus = 100

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 16, m.BalanceT2)
	})
	t.Run("With T0 + T1 + tMinus1 + T2", func(t *testing.T) {
		m := newUser()
		ref := newRef()
		refMinus := newRef()
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		m.ExtraBonus = 100
		m.ActiveT1Referrals = 4
		m.ActiveT2Referrals = 20

		m, _ = mine(testMiningBase, testTime, m, ref, refMinus)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
		require.EqualValues(t, 16, m.BalanceT1)
		require.EqualValues(t, 16, m.BalanceT2)
		require.EqualValues(t, 0.8, m.BalanceForTMinus1)
	})
}

func testSoloMiningWithPreStaking(t *testing.T) {
	t.Run("No referrals", func(t *testing.T) {
		m := newUser()
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
	})
	t.Run("With T0", func(t *testing.T) {
		m := newUser()
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
	})
	t.Run("For tMinus1", func(t *testing.T) {
		m := newUser()
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, nil, ref)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 0.8, m.BalanceForTMinus1)
	})
	t.Run("With T1", func(t *testing.T) {
		m := newUser()
		m.ActiveT1Referrals = 4
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 16, m.BalanceT1)
	})
	t.Run("With T0 + T1", func(t *testing.T) {
		m := newUser()
		m.ActiveT1Referrals = 4
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
		require.EqualValues(t, 16, m.BalanceT1)
	})
	t.Run("With T2", func(t *testing.T) {
		m := newUser()
		m.ActiveT2Referrals = 20
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 16, m.BalanceT2)
	})
	t.Run("With T0 + T1 + tMinus1 + T2", func(t *testing.T) {
		m := newUser()
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		m.ActiveT1Referrals = 4
		m.ActiveT2Referrals = 20
		ref := newRef()
		refMinus := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, refMinus)
		require.NotNil(t, m)
		require.EqualValues(t, 16, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
		require.EqualValues(t, 16, m.BalanceT1)
		require.EqualValues(t, 16, m.BalanceT2)
		require.EqualValues(t, 0.8, m.BalanceForTMinus1)
	})
}

func testSoloMiningWithPreStakingAndExtraBonus(t *testing.T) {
	t.Run("No referrals", func(t *testing.T) {
		m := newUser()
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		m.ExtraBonus = 100
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
	})
	t.Run("With T0", func(t *testing.T) {
		m := newUser()
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		m.ExtraBonus = 100
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
	})
	t.Run("For tMinus1", func(t *testing.T) {
		m := newUser()
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		m.ExtraBonus = 100
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, nil, ref)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 0.8, m.BalanceForTMinus1)
	})
	t.Run("With T1", func(t *testing.T) {
		m := newUser()
		m.ActiveT1Referrals = 4
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		m.ExtraBonus = 100
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 16, m.BalanceT1)
	})
	t.Run("With T0 + T1", func(t *testing.T) {
		m := newUser()
		m.ActiveT1Referrals = 4
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		m.ExtraBonus = 100
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		ref := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
		require.EqualValues(t, 16, m.BalanceT1)
	})
	t.Run("With T2", func(t *testing.T) {
		m := newUser()
		m.ActiveT2Referrals = 20
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		m.ExtraBonus = 100
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 16, m.BalanceT2)
	})
	t.Run("With T0 + T1 + tMinus1 + T2", func(t *testing.T) {
		m := newUser()
		m.ExtraBonusStartedAt = timeDelta(stdlibtime.Hour)
		m.ExtraBonus = 100
		m.PreStakingBonus = 200
		m.PreStakingAllocation = 50
		m.ActiveT1Referrals = 4
		m.ActiveT2Referrals = 20
		ref := newRef()
		refMinus := newRef()

		m, _ = mine(testMiningBase, testTime, m, ref, refMinus)
		require.NotNil(t, m)
		require.EqualValues(t, 32, m.BalanceSolo)
		require.EqualValues(t, 4, m.BalanceT0)
		require.EqualValues(t, 4, m.BalanceForT0)
		require.EqualValues(t, 16, m.BalanceT1)
		require.EqualValues(t, 16, m.BalanceT2)
		require.EqualValues(t, 0.8, m.BalanceForTMinus1)
	})
}

func testSoloMining(t *testing.T) {
	t.Parallel()

	t.Run("No extra bonus", testSoloMiningNoExtraBonus)
	t.Run("With extra bonus", testSoloMiningWithExtraBonus)
	t.Run("With Pre-staking", testSoloMiningWithPreStaking)
	t.Run("With extra bonus + Pre-staking", testSoloMiningWithPreStakingAndExtraBonus)
}

func testNegativeMiningSoloSlashing(t *testing.T) {
	m := newUser()
	m.BalanceLastUpdatedAt = timeDelta(-stdlibtime.Hour)
	m.MiningSessionSoloStartedAt = timeDelta(25 * stdlibtime.Hour)
	m.MiningSessionSoloEndedAt = timeDelta(-stdlibtime.Hour)
	m.BalanceSolo = 1440
	m.BalanceT0 = 1440
	m.BalanceForT0 = 1440
	m.BalanceT1 = 1440
	m.BalanceT2 = 1440
	m.BalanceForTMinus1 = 1440
	m.IDT0 = testIDT0

	m, _ = mine(testMiningBase, testTime, m, nil, nil)
	require.NotNil(t, m)

	require.EqualValues(t, 1, m.SlashingRateSolo)
	require.EqualValues(t, 1, m.SlashingRateT0)
	require.EqualValues(t, 1, m.SlashingRateT1)
	require.EqualValues(t, 1, m.SlashingRateT2)
	require.EqualValues(t, 0, m.SlashingRateForT0)
	require.EqualValues(t, 0, m.SlashingRateForTMinus1)

	require.EqualValues(t, 1439, m.BalanceSolo)
	require.EqualValues(t, 1439, m.BalanceT0)
	require.EqualValues(t, 1439, m.BalanceT1)
	require.EqualValues(t, 1439, m.BalanceT2)
}

func testNegativeMiningT0Slashing(t *testing.T) {
	m := newUser()
	m.BalanceLastUpdatedAt = timeDelta(-stdlibtime.Hour)
	m.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour)
	m.MiningSessionSoloEndedAt = timeDelta(23 * stdlibtime.Hour)
	m.BalanceT0 = 1440
	m.BalanceForT0 = 1440
	m.IDT0 = testIDT0

	ref := newRef()
	ref.MiningSessionSoloStartedAt = timeDelta(-25 * stdlibtime.Hour)
	ref.MiningSessionSoloEndedAt = timeDelta(-stdlibtime.Hour)

	m, _ = mine(testMiningBase, testTime, m, ref, nil)
	require.NotNil(t, m)

	require.EqualValues(t, 0, m.SlashingRateSolo)
	require.EqualValues(t, 0, m.SlashingRateT0)
	require.EqualValues(t, 1, m.SlashingRateForT0)
	require.EqualValues(t, 1439, m.BalanceForT0)
}

func testNegativeMiningT0SlashingSoloSlashing(t *testing.T) {
	m := newUser()
	m.BalanceLastUpdatedAt = timeDelta(-stdlibtime.Hour)
	m.MiningSessionSoloStartedAt = timeDelta(-25 * stdlibtime.Hour)
	m.MiningSessionSoloEndedAt = timeDelta(-stdlibtime.Hour)
	m.BalanceSolo = 1440
	m.BalanceT0 = 1440
	m.BalanceForT0 = 1440
	m.BalanceT1 = 1440
	m.BalanceT2 = 1440
	m.IDT0 = testIDT0

	ref := newRef()
	ref.MiningSessionSoloStartedAt = timeDelta(-25 * stdlibtime.Hour)
	ref.MiningSessionSoloEndedAt = timeDelta(-stdlibtime.Hour)

	m, _ = mine(testMiningBase, testTime, m, ref, nil)
	require.NotNil(t, m)

	require.EqualValues(t, 1, m.SlashingRateSolo)
	require.EqualValues(t, 1, m.SlashingRateT0)
	require.EqualValues(t, 1, m.SlashingRateT1)
	require.EqualValues(t, 1, m.SlashingRateT2)
	require.EqualValues(t, 1, m.SlashingRateForT0)

	require.EqualValues(t, 1439, m.BalanceSolo)
	require.EqualValues(t, 1439, m.BalanceT0)
	require.EqualValues(t, 1439, m.BalanceT1)
	require.EqualValues(t, 1439, m.BalanceT2)
	require.EqualValues(t, 1439, m.BalanceForT0)
}

func testNegativeMiningT1minusSlashingSoloMining(t *testing.T) {
	m := newUser()
	m.BalanceLastUpdatedAt = timeDelta(-stdlibtime.Hour)
	m.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour)
	m.MiningSessionSoloEndedAt = timeDelta(23 * stdlibtime.Hour)
	m.BalanceForTMinus1 = 1440
	m.IDT0 = testIDT0
	m.IDTMinus1 = testIDTMinus1

	ref := newRef()
	ref.MiningSessionSoloStartedAt = timeDelta(-25 * stdlibtime.Hour)
	ref.MiningSessionSoloEndedAt = timeDelta(-stdlibtime.Hour)

	m, _ = mine(testMiningBase, testTime, m, nil, ref)
	require.NotNil(t, m)

	require.EqualValues(t, 0, m.SlashingRateSolo)
	require.EqualValues(t, 1, m.SlashingRateForTMinus1)
	require.EqualValues(t, 1439, m.BalanceForTMinus1)
}

func testNegativeMining(t *testing.T) {
	t.Parallel()

	t.Run("Solo slashing", testNegativeMiningSoloSlashing)
	t.Run("For T0 slashing while Solo is mining", testNegativeMiningT0Slashing)
	t.Run("For T0 slashing while Solo is slashing also", testNegativeMiningT0SlashingSoloSlashing)
	t.Run("For T1Minus slashing while Solo is mining", testNegativeMiningT1minusSlashingSoloMining)
}

func testMiningResurrectT0(t *testing.T) {
	m := newUser()
	m.SlashingRateForT0 = 10
	m.IDT0 = testIDT0

	t0Ref := new(referral)
	t0Ref.MiningSessionSoloStartedAt = timeDelta(0)
	t0Ref.MiningSessionSoloPreviouslyEndedAt = timeDelta(-24 * 10 * stdlibtime.Hour)
	t0Ref.ResurrectSoloUsedAt = timeDelta(stdlibtime.Hour)

	m, _ = mine(testMiningBase, testTime, m, t0Ref, nil)
	require.NotNil(t, m)

	require.EqualValues(t, 2400, m.BalanceForT0)
	require.EqualValues(t, 0, m.SlashingRateForT0)
}

func testMiningResurrectT0ResetSlashing(t *testing.T) {
	m := newUser()
	m.SlashingRateForT0 = 10
	m.IDT0 = testIDT0

	t0Ref := new(referral)
	t0Ref.MiningSessionSoloStartedAt = timeDelta(0)
	t0Ref.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour)
	t0Ref.MiningSessionSoloPreviouslyEndedAt = timeDelta(-24 * 10 * stdlibtime.Hour)
	m.ResurrectT0UsedAt = timeDelta(stdlibtime.Hour)

	m, _ = mine(testMiningBase, testTime, m, t0Ref, nil)
	require.NotNil(t, m)
	require.EqualValues(t, 0, m.SlashingRateForT0)
}

func testMiningResurrectTMinus1ResetSlashing(t *testing.T) {
	m := newUser()
	m.SlashingRateForTMinus1 = 10
	m.IDTMinus1 = testIDTMinus1

	ref := new(referral)
	ref.MiningSessionSoloStartedAt = timeDelta(0)
	ref.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour)
	ref.MiningSessionSoloPreviouslyEndedAt = timeDelta(-24 * 10 * stdlibtime.Hour)
	m.ResurrectTMinus1UsedAt = timeDelta(stdlibtime.Hour)

	m, _ = mine(testMiningBase, testTime, m, nil, ref)
	require.NotNil(t, m)
	require.EqualValues(t, 0, m.SlashingRateForTMinus1)
}

func testMiningResurrectSolo(t *testing.T) {
	m := newUser()
	m.SlashingRateSolo = 10
	m.MiningSessionSoloStartedAt = timeDelta(0)
	m.MiningSessionSoloPreviouslyEndedAt = timeDelta(-24 * 10 * stdlibtime.Hour)
	m.ResurrectSoloUsedAt = timeDelta(stdlibtime.Hour)

	m, _ = mine(testMiningBase, testTime, m, nil, nil)
	require.NotNil(t, m)
	require.EqualValues(t, 2400, m.BalanceSolo)
	require.EqualValues(t, 0, m.SlashingRateSolo)
}

func testMiningResurrectT1(t *testing.T) {
	m := newUser()
	m.SlashingRateForTMinus1 = 10
	m.IDTMinus1 = testIDTMinus1
	m.IDT0 = testIDT0

	ref := new(referral)
	ref.MiningSessionSoloStartedAt = timeDelta(0)
	ref.MiningSessionSoloPreviouslyEndedAt = timeDelta(-24 * 10 * stdlibtime.Hour)
	ref.ResurrectSoloUsedAt = timeDelta(stdlibtime.Hour)

	m, _ = mine(testMiningBase, testTime, m, nil, ref)
	require.NotNil(t, m)
	require.EqualValues(t, 2400, m.BalanceForTMinus1)
	require.EqualValues(t, 0, m.SlashingRateForTMinus1)
}

func testMiningResurrect(t *testing.T) {
	t.Parallel()

	t.Run("Solo", testMiningResurrectSolo)
	t.Run("T0", testMiningResurrectT0)
	t.Run("T1", testMiningResurrectT1)

	t.Run("T0_ResetSlashing", testMiningResurrectT0ResetSlashing)
	t.Run("T1_ResetSlashing", testMiningResurrectTMinus1ResetSlashing)
}

func Test_BalancePositive(t *testing.T) {
	t.Parallel()

	t.Run("Solo mining", testSoloMining)
	t.Run("Negative mining", testNegativeMining)
	t.Run("Resurrect", testMiningResurrect)
}

func Test_MinerNil(t *testing.T) {
	t.Parallel()

	m, h := mine(testMiningBase, testTime, nil, nil, nil)
	require.Nil(t, m)
	require.Nil(t, h)
}

func Test_MinerPending(t *testing.T) {
	t.Parallel()

	t.Run("Apply", func(t *testing.T) {
		m := newUser()
		m.MiningSessionSoloEndedAt = timeDelta(-stdlibtime.Hour)
		m.BalanceT1Pending = 2
		m.BalanceT2Pending = 2
		m.BalanceT1PendingApplied = 1
		m.BalanceT2PendingApplied = 1

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.NotNil(t, m)
		require.EqualValues(t, 2.0, m.BalanceT1PendingApplied)
		require.EqualValues(t, 2.0, m.BalanceT2PendingApplied)
		require.EqualValues(t, 2.0, m.BalanceT1Pending)
		require.EqualValues(t, 2.0, m.BalanceT2Pending)
	})

	t.Run("Skip", func(t *testing.T) {
		m := newUser()
		m.MiningSessionSoloEndedAt = timeDelta(-stdlibtime.Hour)
		m.BalanceT1Pending = 1
		m.BalanceT2Pending = 1
		m.BalanceT1PendingApplied = 1
		m.BalanceT2PendingApplied = 1

		m, _ = mine(testMiningBase, testTime, m, nil, nil)
		require.Nil(t, m)
	})
}

func Test_MinerWithHistory(t *testing.T) {
	t.Parallel()

	m := newUser()
	m.BalanceLastUpdatedAt = time.New(testTime.Add(-stdlibtime.Hour * 2))

	m, h := mine(testMiningBase, testTime, m, nil, nil)
	require.NotNil(t, m)
	require.NotNil(t, h)

	require.EqualValues(t, float64(testMiningBase)*2, m.BalanceSolo)

	t.Logf("new:     %p", m)
	t.Logf("history: %p", h)

	// Pointers must be different.
	require.False(t, m == h)
}

func Test_MinerNegativeBalance(t *testing.T) {
	t.Parallel()

	m := newUser()
	m.MiningSessionSoloStartedAt = m.MiningSessionSoloEndedAt
	m.IDT0 = testIDT0
	m.IDTMinus1 = testIDTMinus1
	m.BalanceSolo = -1
	m.BalanceT0 = -2
	m.BalanceT1 = -3
	m.BalanceT2 = -4
	m.BalanceForT0 = -5
	m.BalanceForTMinus1 = -6
	m.ActiveT1Referrals = -7
	m.ActiveT2Referrals = -8

	m, _ = mine(testMiningBase, testTime, m, nil, nil)
	require.NotNil(t, m)

	require.Zero(t, m.BalanceSolo)
	require.Zero(t, m.BalanceT0)
	require.Zero(t, m.BalanceT1)
	require.Zero(t, m.BalanceT2)
	require.Zero(t, m.BalanceForT0)
	require.Zero(t, m.BalanceForTMinus1)
	require.Zero(t, m.ActiveT1Referrals)
	require.Zero(t, m.ActiveT2Referrals)
}

func testMinerPendingSlashingSolo(t *testing.T) {
	t.Parallel()

	m := newUser()
	m.MiningSessionSoloStartedAt = m.MiningSessionSoloEndedAt
	m.IDT0 = testIDT0
	m.IDTMinus1 = testIDTMinus1
	m.BalanceSoloPending = 1
	m.BalanceSoloPendingApplied = 3

	m, _ = mine(testMiningBase, testTime, m, nil, nil)
	require.NotNil(t, m)
	require.EqualValues(t, 1., m.BalanceSoloPendingApplied)
}

func testMinerPendingSlashingT1(t *testing.T) {
	t.Parallel()

	m := newUser()
	m.MiningSessionSoloStartedAt = m.MiningSessionSoloEndedAt
	m.IDT0 = testIDT0
	m.IDTMinus1 = testIDTMinus1
	m.BalanceT1Pending = 1
	m.BalanceT1PendingApplied = 3

	m, _ = mine(testMiningBase, testTime, m, nil, nil)
	require.NotNil(t, m)
	require.EqualValues(t, 1., m.BalanceT1PendingApplied)
}

func testMinerPendingSlashingT2(t *testing.T) {
	t.Parallel()

	m := newUser()
	m.MiningSessionSoloStartedAt = m.MiningSessionSoloEndedAt
	m.IDT0 = testIDT0
	m.IDTMinus1 = testIDTMinus1
	m.BalanceT2Pending = 1
	m.BalanceT2PendingApplied = 3

	m, _ = mine(testMiningBase, testTime, m, nil, nil)
	require.NotNil(t, m)
	require.EqualValues(t, 1., m.BalanceT2PendingApplied)
}

func Test_MinerPendingSlashing(t *testing.T) {
	t.Parallel()

	t.Run("Solo", testMinerPendingSlashingSolo)
	t.Run("T1", testMinerPendingSlashingT1)
	t.Run("T2", testMinerPendingSlashingT2)
}
