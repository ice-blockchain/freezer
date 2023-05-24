// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testDidANewDayOffJustStartEmpty(t *testing.T) {
	t.Helper()
	require.Nil(t, didANewDayOffJustStart(nil, nil))
	require.Nil(t, didANewDayOffJustStart(testTime, nil))

	m := newUser()
	require.Nil(t, didANewDayOffJustStart(testTime, m))

	m.BalanceLastUpdatedAt = testTime
	m.MiningSessionSoloLastStartedAt = timeDelta(-time.Hour * 25)
	require.Nil(t, didANewDayOffJustStart(testTime, m))
}

func testDidANewDayOffJustStartDayOff(t *testing.T) {
	t.Helper()
	m := newUser()
	m.BalanceLastUpdatedAt = timeDelta(-time.Hour * 2)
	m.MiningSessionSoloLastStartedAt = timeDelta(-time.Hour * 25)

	started := didANewDayOffJustStart(testTime, m)
	require.NotNil(t, started)
}

func Test_didANewDayOffJustStart(t *testing.T) {
	t.Parallel()

	t.Run("Empty", testDidANewDayOffJustStartEmpty)
	t.Run("DayOff", testDidANewDayOffJustStartDayOff)
}
