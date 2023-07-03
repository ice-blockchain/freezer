// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_didReferralJustStopMining(t *testing.T) {
	t.Parallel()

	t.Run("EmptyData", func(t *testing.T) {
		x := didReferralJustStopMining(testTime, nil, nil, nil)
		require.Nil(t, x)
	})

	t.Run("Stopped", func(t *testing.T) {
		before := newUser()
		before.BalanceLastUpdatedAt = timeDelta(-time.Hour * 2)
		before.MiningSessionSoloEndedAt = timeDelta(-time.Hour)

		x := didReferralJustStopMining(testTime, before, nil, nil)
		require.NotNil(t, x)
		require.NotNil(t, x.StoppedMiningAt.Time)
	})

	t.Run("Full parameters list", func(t *testing.T) {
		before := newUser()
		before.BalanceLastUpdatedAt = timeDelta(-time.Hour * 2)
		before.MiningSessionSoloEndedAt = timeDelta(-time.Hour)

		t0Ref := newRef()
		tMinus1Ref := newRef()

		x := didReferralJustStopMining(testTime, before, t0Ref, tMinus1Ref)
		require.NotNil(t, x)
		require.NotNil(t, x.StoppedMiningAt.Time)
		require.Equal(t, before.ID, x.ID)
		require.Equal(t, t0Ref.ID, x.IDT0)
		require.Equal(t, tMinus1Ref.ID, x.IDTMinus1)
	})
}
