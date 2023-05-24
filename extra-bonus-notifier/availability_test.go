// SPDX-License-Identifier: ice License 1.0

package extrabonusnotifier

import (
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/wintr/time"
)

var (
	testTime = time.New(stdlibtime.Date(2023, 1, 2, 3, 4, 5, 6, stdlibtime.UTC))
)

func newUser() *user {
	u := new(user)
	u.UserID = "test_user_id"
	u.ID = 111_111

	return u
}

func generateExtraBonusIndicesDistributionFromMap(totalChunkNumber uint16, m map[uint16]bool) map[uint16]map[uint16]uint16 {
	days := make(map[uint16]uint16, len(m))
	for day, value := range m {
		if value {
			days[day] = 42
		}
	}

	result := make(map[uint16]map[uint16]uint16, totalChunkNumber)
	for i := uint16(0); i < totalChunkNumber; i++ {
		result[i] = days
	}

	return result
}

func Test_isExtraBonusAvailable(t *testing.T) {
	t.Parallel()

	weekAgo := testTime.Add(-stdlibtime.Hour * 24 * 7)
	ExtraBonusStartedAt := time.New(stdlibtime.Date(weekAgo.Year(), weekAgo.Month(), weekAgo.Day(), 00, 00, 00, 00, weekAgo.Location()))

	t.Run("Time is before 10am and there is a bonus for that day", func(t *testing.T) {
		now := time.New(stdlibtime.Date(testTime.Year(), testTime.Month(), testTime.Day(), 6, 00, 00, 00, testTime.Location()))
		d := generateExtraBonusIndicesDistributionFromMap(1000, map[uint16]bool{8: true})

		m := newUser()
		m.ExtraBonusStartedAt = ExtraBonusStartedAt
		m.UTCOffset = 180
		m.extraBonusIndex = 7.

		b := isExtraBonusAvailable(now, ExtraBonusStartedAt, d, m)
		require.False(t, b)
		require.EqualValues(t, 7, m.extraBonusIndex)
	})

	t.Run("Time is after 8pm and there is a bonus for that day", func(t *testing.T) {
		now := time.New(stdlibtime.Date(testTime.Year(), testTime.Month(), testTime.Day(), 18, 00, 00, 00, testTime.Location()))
		d := generateExtraBonusIndicesDistributionFromMap(1000, map[uint16]bool{8: true})

		m := newUser()
		m.ExtraBonusStartedAt = ExtraBonusStartedAt
		m.UTCOffset = 180
		m.extraBonusIndex = 7.

		b := isExtraBonusAvailable(now, ExtraBonusStartedAt, d, m)
		require.False(t, b)
		require.EqualValues(t, 7, m.extraBonusIndex)
	})

	t.Run("Time is after 10am but before user's interval, and there is a bonus for that day", func(t *testing.T) {
		now := time.New(stdlibtime.Date(testTime.Year(), testTime.Month(), testTime.Day(), 7, 00, 00, 00, testTime.Location()))
		d := generateExtraBonusIndicesDistributionFromMap(1000, map[uint16]bool{8: true})

		m := newUser()
		m.ExtraBonusStartedAt = ExtraBonusStartedAt
		m.UTCOffset = 180
		m.extraBonusIndex = 7.

		b := isExtraBonusAvailable(now, ExtraBonusStartedAt, d, m)
		require.False(t, b)
		require.EqualValues(t, 7, m.extraBonusIndex)
	})

	t.Run("Time is after 10am and after user's interval, and there is a bonus for that day", func(t *testing.T) {
		now := time.New(stdlibtime.Date(testTime.Year(), testTime.Month(), testTime.Day(), 8, 00, 00, 00, testTime.Location()))
		d := generateExtraBonusIndicesDistributionFromMap(1000, map[uint16]bool{8: true})

		m := newUser()
		m.ExtraBonusStartedAt = ExtraBonusStartedAt
		m.UTCOffset = 180
		m.extraBonusIndex = 7.

		b := isExtraBonusAvailable(now, ExtraBonusStartedAt, d, m)
		require.True(t, b)

		require.EqualValues(t, 8, m.extraBonusIndex)
		require.EqualValues(t, now, m.ExtraBonusLastClaimAvailableAt)
	})
	t.Run("Time is after 10am and after user's interval, but there's no bonus that day", func(t *testing.T) {
		now := time.New(stdlibtime.Date(testTime.Year(), testTime.Month(), testTime.Day(), 8, 00, 00, 00, testTime.Location()))
		d := generateExtraBonusIndicesDistributionFromMap(1000, map[uint16]bool{8: false})

		m := newUser()
		m.ExtraBonusStartedAt = ExtraBonusStartedAt
		m.UTCOffset = 180
		m.extraBonusIndex = 7.

		b := isExtraBonusAvailable(now, ExtraBonusStartedAt, d, m)
		require.False(t, b)
		require.EqualValues(t, 7, m.extraBonusIndex)
	})

	t.Run("Time is after 10am and after user's interval, after user missed bonus for 2 days, and there is a bonus for that day, and there was bonus for day 8 and day 9", func(t *testing.T) {
		ts := testTime.Add(stdlibtime.Hour * 24 * 2)
		now := time.New(stdlibtime.Date(ts.Year(), ts.Month(), ts.Day(), 8, 00, 00, 00, ts.Location()))
		d := generateExtraBonusIndicesDistributionFromMap(1000, map[uint16]bool{8: true, 9: true, 10: true})

		m := newUser()
		m.ExtraBonusStartedAt = ExtraBonusStartedAt
		m.UTCOffset = 180
		m.extraBonusIndex = 7.

		b := isExtraBonusAvailable(now, ExtraBonusStartedAt, d, m)
		require.True(t, b)

		require.EqualValues(t, 10, m.extraBonusIndex)
		require.EqualValues(t, now, m.ExtraBonusLastClaimAvailableAt)
	})

	t.Run("Time is after 10am and after user's interval, after user missed bonus for 2 days, and there is a bonus for that day, and there was bonus for day 9, but not for 8", func(t *testing.T) {
		ts := testTime.Add(stdlibtime.Hour * 24 * 2)
		now := time.New(stdlibtime.Date(ts.Year(), ts.Month(), ts.Day(), 8, 00, 00, 00, ts.Location()))
		d := generateExtraBonusIndicesDistributionFromMap(1000, map[uint16]bool{8: false, 9: true, 10: true})

		m := newUser()
		m.ExtraBonusStartedAt = ExtraBonusStartedAt
		m.UTCOffset = 180
		m.extraBonusIndex = 7.

		b := isExtraBonusAvailable(now, ExtraBonusStartedAt, d, m)
		require.True(t, b)
		require.EqualValues(t, 10, m.extraBonusIndex)
		require.EqualValues(t, now, m.ExtraBonusLastClaimAvailableAt)
	})

	t.Run("With additional ExtraBonusLastClaimAvailableAt", func(t *testing.T) {
		ts := testTime.Add(stdlibtime.Hour * 24 * 2)
		now := time.New(stdlibtime.Date(ts.Year(), ts.Month(), ts.Day(), 8, 00, 00, 00, ts.Location()))
		d := generateExtraBonusIndicesDistributionFromMap(1000, map[uint16]bool{8: false, 9: true, 10: true})

		m := newUser()
		m.ExtraBonusStartedAt = ExtraBonusStartedAt
		m.ExtraBonusLastClaimAvailableAt = time.New(now.Add(-stdlibtime.Hour * 96))
		m.UTCOffset = 180
		m.extraBonusIndex = 7.

		b := isExtraBonusAvailable(now, ExtraBonusStartedAt, d, m)
		require.True(t, b)
		require.EqualValues(t, 10, m.extraBonusIndex)
		require.EqualValues(t, now, m.ExtraBonusLastClaimAvailableAt)
	})
}
