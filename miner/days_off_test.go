// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"fmt"
	"regexp"
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/wintr/time"
)

func testDidANewDayOffJustStartEmpty(t *testing.T) {
	t.Helper()
	require.Nil(t, didANewDayOffJustStart(nil, nil))
	require.Nil(t, didANewDayOffJustStart(testTime, nil))

	m := newUser()
	require.Nil(t, didANewDayOffJustStart(testTime, m))

	m.BalanceLastUpdatedAt = testTime
	m.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 25)
	require.Nil(t, didANewDayOffJustStart(testTime, m))
}

func testDidANewDayOffJustStartDayOff(t *testing.T) {
	t.Helper()
	m := newUser()
	m.BalanceLastUpdatedAt = timeDelta(-stdlibtime.Hour * 2)
	m.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 25)

	started := didANewDayOffJustStart(testTime, m)
	require.NotNil(t, started)
	require.Regexp(t, regexp.MustCompile(fmt.Sprintf("%v~[0-9]+", m.UserID)), started.ID)
	require.Equal(t, m.UserID, started.UserID)
	require.EqualValues(t, 0, started.MiningStreak)
	require.EqualValues(t, 0, started.RemainingFreeMiningSessions)
}

func testDidANewDayOffJustStartDayOff_MiningSessionSoloStartedAtIsNil(t *testing.T) {
	t.Helper()
	m := newUser()

	m.BalanceLastUpdatedAt = time.New(testTime.Add(-1 * stdlibtime.Hour))
	m.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloStartedAt = nil
	m.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 49)

	started := didANewDayOffJustStart(testTime, m)
	require.Nil(t, started)
}

func testDidANewDayOffJustStartDayOff_MiningSessionSoloLastStartedAtIsNil(t *testing.T) {
	t.Helper()
	m := newUser()

	m.BalanceLastUpdatedAt = time.New(testTime.Add(-1 * stdlibtime.Hour))
	m.MiningSessionSoloLastStartedAt = nil
	m.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 49)

	started := didANewDayOffJustStart(testTime, m)
	require.Nil(t, started)
}

func testDidANewDayOffJustStartDayOff_MiningSessionSoloEndedAtIsNil(t *testing.T) {
	t.Helper()
	m := newUser()

	m.BalanceLastUpdatedAt = time.New(testTime.Add(-1 * stdlibtime.Hour))
	m.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloEndedAt = nil

	started := didANewDayOffJustStart(testTime, m)
	require.Nil(t, started)
}

func testDidANewDayOffJustStartDayOff_BalanceLastUpdatedAtIsNil(t *testing.T) {
	t.Helper()
	m := newUser()

	m.BalanceLastUpdatedAt = nil
	m.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 49)

	started := didANewDayOffJustStart(testTime, m)
	require.Nil(t, started)
}

func testDidANewDayOffJustStartDayOff_MiningSessionLastStartedAfterNow(t *testing.T) {
	t.Helper()
	m := newUser()

	m.BalanceLastUpdatedAt = testTime
	m.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 23)
	m.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 25)

	started := didANewDayOffJustStart(testTime, m)
	require.Nil(t, started)
}

func testDidANewDayOffJustStartDayOff_BalanceLastUpdatedAtAfterStartedAt(t *testing.T) {
	t.Helper()
	m := newUser()

	m.BalanceLastUpdatedAt = testTime
	m.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 49)

	started := didANewDayOffJustStart(testTime, m)
	require.Nil(t, started)
}

func testDidANewDayOffJustStartDayOff_RemainingFreeMiningSessions(t *testing.T) {
	t.Helper()
	m := newUser()

	m.BalanceLastUpdatedAt = time.New(testTime.Add(-1 * stdlibtime.Hour))
	m.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 144)
	m.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 25)

	started := didANewDayOffJustStart(testTime, m)
	require.NotNil(t, started)
	require.Regexp(t, regexp.MustCompile(fmt.Sprintf("%v~[0-9]+", m.UserID)), started.ID)
	require.Equal(t, m.UserID, started.UserID)
	require.EqualValues(t, 6, started.MiningStreak)
	require.EqualValues(t, 1, started.RemainingFreeMiningSessions)
}

func testDidANewDayOffJustStartDayOff_ConcurrentCallsWithTheSameStartedAtTime(t *testing.T) {
	t.Helper()
	m1, m2 := newUser(), newUser()

	m1.BalanceLastUpdatedAt = time.New(testTime.Add(-1 * stdlibtime.Hour))
	m1.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m1.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 144)
	m1.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 25)

	m2.UserID = "test_user_id2"
	m2.BalanceLastUpdatedAt = time.New(testTime.Add(-1 * stdlibtime.Hour))
	m2.MiningSessionSoloLastStartedAt = timeDelta(-stdlibtime.Hour * 49)
	m2.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 144)
	m2.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 25)

	started1 := didANewDayOffJustStart(testTime, m1)
	require.NotNil(t, started1)
	require.Regexp(t, regexp.MustCompile(fmt.Sprintf("%v~[0-9]+", m1.UserID)), started1.ID)
	require.Equal(t, m1.UserID, started1.UserID)
	require.EqualValues(t, 6, started1.MiningStreak)
	require.EqualValues(t, 1, started1.RemainingFreeMiningSessions)

	started2 := didANewDayOffJustStart(testTime, m2)
	require.NotNil(t, started2)
	require.Regexp(t, regexp.MustCompile(fmt.Sprintf("%v~[0-9]+", m2.UserID)), started2.ID)
	require.Equal(t, m2.UserID, started2.UserID)
	require.EqualValues(t, 6, started2.MiningStreak)
	require.EqualValues(t, 1, started2.RemainingFreeMiningSessions)
	require.NotEqual(t, started1.ID, started2.ID)
}

func Test_didANewDayOffJustStart(t *testing.T) {
	t.Parallel()

	t.Run("Empty", testDidANewDayOffJustStartEmpty)
	t.Run("DayOff", testDidANewDayOffJustStartDayOff)
	t.Run("MiningSessionSoloStartedAt is nil", testDidANewDayOffJustStartDayOff_MiningSessionSoloStartedAtIsNil)
	t.Run("MiningSessionSoloLastStartedAt is nil", testDidANewDayOffJustStartDayOff_MiningSessionSoloLastStartedAtIsNil)
	t.Run("BalanceLastUpdatedAt is nil", testDidANewDayOffJustStartDayOff_BalanceLastUpdatedAtIsNil)
	t.Run("MiningSessionSoloEndedAt is nil", testDidANewDayOffJustStartDayOff_MiningSessionSoloEndedAtIsNil)
	t.Run("MiningSessionLastStarted after now", testDidANewDayOffJustStartDayOff_MiningSessionLastStartedAfterNow)
	t.Run("BalanceLastUpdatedAt after startedAt", testDidANewDayOffJustStartDayOff_BalanceLastUpdatedAtAfterStartedAt)
	t.Run("RemainingFreeMiningSessions check", testDidANewDayOffJustStartDayOff_RemainingFreeMiningSessions)
	t.Run("Concurrent calls with the same startedAt time", testDidANewDayOffJustStartDayOff_ConcurrentCallsWithTheSameStartedAtTime)
}
