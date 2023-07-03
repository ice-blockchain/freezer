// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"testing"
	stdlibtime "time"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	"github.com/ice-blockchain/wintr/time"
	"github.com/stretchr/testify/assert"
)

func TestCalculateDates(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}

	/******************************************************************************************************************************************************
		1. limit = 24, offset = 0, factor = 1
	******************************************************************************************************************************************************/
	limit := uint64(24)
	offset := uint64(0)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	dates, notBeforeTime := repo.calculateDates(limit, offset, start, 1)
	assert.Len(t, dates, 30)

	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 4, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
		*expectedStart.Time,
		expectedStart.Add(1 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(5 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(6 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(7 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(8 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(9 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(10 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(11 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(12 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(13 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(14 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(15 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(16 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(17 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(18 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(19 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(20 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(21 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(22 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(23 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(24 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(25 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(26 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(27 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(28 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(29 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)

	/******************************************************************************************************************************************************
		2. limit = 12, offset = 0, factor = 1
	******************************************************************************************************************************************************/
	limit = uint64(12)
	offset = uint64(0)
	start = time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	dates, notBeforeTime = repo.calculateDates(limit, offset, start, 1)
	assert.Len(t, dates, 6)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)

	expectedStart = time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected = []stdlibtime.Time{
		*expectedStart.Time,
		expectedStart.Add(1 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(5 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)

	/******************************************************************************************************************************************************
		3. limit = 36, offset = 0, factor = 1
	******************************************************************************************************************************************************/
	limit = uint64(36)
	offset = uint64(0)
	start = time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	dates, notBeforeTime = repo.calculateDates(limit, offset, start, 1)
	assert.Len(t, dates, 30)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 4, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	expectedStart = time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected = []stdlibtime.Time{
		*expectedStart.Time,
		expectedStart.Add(1 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(5 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(6 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(7 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(8 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(9 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(10 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(11 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(12 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(13 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(14 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(15 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(16 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(17 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(18 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(19 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(20 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(21 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(22 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(23 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(24 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(25 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(26 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(27 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(28 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(29 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)

	/******************************************************************************************************************************************************
		4. limit = 48, offset = 0, factor = 1
	******************************************************************************************************************************************************/
	limit = uint64(48)
	offset = uint64(0)
	start = time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	dates, notBeforeTime = repo.calculateDates(limit, offset, start, 1)
	assert.Len(t, dates, 54)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 3, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	expectedStart = time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected = []stdlibtime.Time{
		*expectedStart.Time,
		expectedStart.Add(1 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(5 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(6 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(7 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(8 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(9 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(10 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(11 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(12 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(13 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(14 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(15 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(16 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(17 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(18 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(19 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(20 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(21 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(22 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(23 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(24 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(25 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(26 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(27 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(28 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(29 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(30 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(31 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(32 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(33 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(34 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(35 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(36 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(37 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(38 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(39 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(40 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(41 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(42 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(43 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(44 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(45 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(46 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(47 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(48 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(49 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(50 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(51 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(52 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(53 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)

	/******************************************************************************************************************************************************
		5. limit = 24, offset = 0, factor = -1
	******************************************************************************************************************************************************/
	limit = uint64(24)
	offset = uint64(0)
	start = time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	dates, notBeforeTime = repo.calculateDates(limit, offset, start, -1)
	assert.Len(t, dates, 30)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 4, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)

	expectedStart = time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected = []stdlibtime.Time{
		*expectedStart.Time,
		expectedStart.Add(-1 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-5 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-6 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-7 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-8 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-9 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-10 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-11 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-12 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-13 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-14 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-15 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-16 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-17 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-18 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-19 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-20 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-21 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-22 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-23 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-24 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-25 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-26 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-27 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-28 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-29 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)

	/******************************************************************************************************************************************************
		6. limit = 24, offset = 24, factor = -1
	******************************************************************************************************************************************************/
	limit = uint64(24)
	offset = uint64(24)
	start = time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	dates, notBeforeTime = repo.calculateDates(limit, offset, start, -1)
	assert.Len(t, dates, 30)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 3, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	expectedStart = time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected = []stdlibtime.Time{
		expectedStart.Add(-24 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-25 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-26 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-27 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-28 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-29 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-30 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-31 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-32 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-33 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-34 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-35 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-36 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-37 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-38 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-39 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-40 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-41 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-42 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-43 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-44 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-45 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-46 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-47 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-48 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-49 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-50 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-51 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-52 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-53 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)

	/******************************************************************************************************************************************************
		7. limit = 36, offset = 36, factor = -1
	******************************************************************************************************************************************************/
	limit = uint64(24)
	offset = uint64(24)
	start = time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	dates, notBeforeTime = repo.calculateDates(limit, offset, start, -1)
	assert.Len(t, dates, 30)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 3, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	expectedStart = time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected = []stdlibtime.Time{
		expectedStart.Add(-24 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-25 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-26 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-27 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-28 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-29 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-30 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-31 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-32 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-33 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-34 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-35 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-36 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-37 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-38 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-39 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-40 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-41 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-42 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-43 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-44 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-45 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-46 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-47 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-48 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-49 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-50 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-51 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-52 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-53 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)

	/******************************************************************************************************************************************************
		7. limit = 48, offset = 48, factor = -1
	******************************************************************************************************************************************************/
	limit = uint64(24)
	offset = uint64(24)
	start = time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	dates, notBeforeTime = repo.calculateDates(limit, offset, start, -1)
	assert.Len(t, dates, 30)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 3, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	expectedStart = time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected = []stdlibtime.Time{
		expectedStart.Add(-24 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-25 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-26 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-27 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-28 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-29 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-30 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-31 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-32 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-33 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-34 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-35 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-36 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-37 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-38 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-39 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-40 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-41 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-42 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-43 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-44 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-45 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-46 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-47 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-48 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-49 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-50 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-51 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-52 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-53 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)
}

func TestProcessBalanceHistory_ChildIsHour(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	now := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))

	/******************************************************************************************************************************************************
		1. History - data from clickhouse.
	******************************************************************************************************************************************************/
	history := []*dwh.BalanceHistory{
		{
			CreatedAt:           time.New(now.Add(-1 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  25.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-2 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  28.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-3 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  32.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-4 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  31.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-5 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  25.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-6 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 17.,
		},
		{
			CreatedAt:           time.New(now.Add(-7 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  20.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-8 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 15.,
		},
		{
			CreatedAt:           time.New(now.Add(-9 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 15.,
		},
		{
			CreatedAt:           time.New(now.Add(-10 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  30.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-11 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 29.,
		},
		{
			CreatedAt:           time.New(now.Add(-12 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  27.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-13 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  30.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-14 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  31.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-15 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 27.,
		},
		{
			CreatedAt:           time.New(now.Add(-16 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  10.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-17 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 8.,
		},
		{
			CreatedAt:           time.New(now.Add(-18 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  15.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-19 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 10.,
		},
		{
			CreatedAt:           time.New(now.Add(-20 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  30.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-21 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  28.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-22 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 25.,
		},
		{
			CreatedAt:           time.New(now.Add(-23 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  22.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-24 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 21.,
		},
		{
			CreatedAt:           time.New(now.Add(-25 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  15.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-26 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  30.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-27 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 32.,
		},
		{
			CreatedAt:           time.New(now.Add(-28 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  29.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-29 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  27.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-30 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  1.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-31 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 25.,
		},
		{
			CreatedAt:           time.New(now.Add(-32 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  20.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-33 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  32.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(-34 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 10.,
		},
		{
			CreatedAt:           time.New(now.Add(-35 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  32.,
			BalanceTotalSlashed: 0.,
		},
	}
	/******************************************************************************************************************************************************
		1. Not before time is -10 hours. startDateIsBeforeEndDate = true.
	******************************************************************************************************************************************************/
	notBeforeTime := time.New(now.Add(-10 * repo.cfg.GlobalAggregationInterval.Child))
	startDateIsBeforeEndDate := true

	entries := repo.processBalanceHistory(history, startDateIsBeforeEndDate, notBeforeTime)

	assert.Len(t, entries, 3)
	assert.Len(t, entries[0].TimeSeries, 0)
	assert.Len(t, entries[1].TimeSeries, 5)
	assert.Len(t, entries[2].TimeSeries, 5)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 3, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[0].Time.UnixNano())
	assert.Equal(t, 50., entries[0].Balance.amount)
	assert.Equal(t, "50.00", entries[0].Balance.Amount)
	assert.Equal(t, false, entries[0].Balance.Negative)
	assert.Equal(t, int64(0), entries[0].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[1].Time.UnixNano())
	assert.Equal(t, "145.00", entries[1].Balance.Amount)
	assert.Equal(t, 145., entries[1].Balance.amount)
	assert.Equal(t, false, entries[1].Balance.Negative)
	assert.Equal(t, int64(190), entries[1].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[2].Time.UnixNano())
	assert.Equal(t, "141.00", entries[2].Balance.Amount)
	assert.Equal(t, 141., entries[2].Balance.amount)
	assert.Equal(t, false, entries[2].Balance.Negative)
	assert.Equal(t, int64(-2), entries[2].Balance.Bonus)

	timeSeries0 := entries[0].TimeSeries
	assert.Len(t, timeSeries0, 0)
	assert.Empty(t, timeSeries0)

	timeSeries1 := entries[1].TimeSeries
	assert.Len(t, timeSeries1, 5)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 19, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[0].Time.UnixNano())
	assert.Equal(t, "30.00", timeSeries1[0].Balance.Amount)
	assert.Equal(t, 30., timeSeries1[0].Balance.amount)
	assert.Equal(t, false, timeSeries1[0].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries1[0].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 20, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[1].Time.UnixNano())
	assert.Equal(t, "15.00", timeSeries1[1].Balance.Amount)
	assert.Equal(t, -15., timeSeries1[1].Balance.amount)
	assert.Equal(t, true, timeSeries1[1].Balance.Negative)
	assert.Equal(t, int64(-150), timeSeries1[1].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 21, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[2].Time.UnixNano())
	assert.Equal(t, "15.00", timeSeries1[2].Balance.Amount)
	assert.Equal(t, -15., timeSeries1[2].Balance.amount)
	assert.Equal(t, true, timeSeries1[2].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries1[2].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 22, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[3].Time.UnixNano())
	assert.Equal(t, "20.00", timeSeries1[3].Balance.Amount)
	assert.Equal(t, 20., timeSeries1[3].Balance.amount)
	assert.Equal(t, false, timeSeries1[3].Balance.Negative)
	assert.Equal(t, int64(233), timeSeries1[3].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 23, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[4].Time.UnixNano())
	assert.Equal(t, "17.00", timeSeries1[4].Balance.Amount)
	assert.Equal(t, -17., timeSeries1[4].Balance.amount)
	assert.Equal(t, true, timeSeries1[4].Balance.Negative)
	assert.Equal(t, int64(-185), timeSeries1[4].Balance.Bonus)

	timeSeries2 := entries[2].TimeSeries
	assert.Len(t, timeSeries2, 5)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 0, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[0].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries2[0].Balance.Amount)
	assert.Equal(t, 25., timeSeries2[0].Balance.amount)
	assert.Equal(t, false, timeSeries2[0].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries2[0].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 1, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[1].Time.UnixNano())
	assert.Equal(t, "31.00", timeSeries2[1].Balance.Amount)
	assert.Equal(t, 31., timeSeries2[1].Balance.amount)
	assert.Equal(t, false, timeSeries2[1].Balance.Negative)
	assert.Equal(t, int64(24), timeSeries2[1].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 2, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[2].Time.UnixNano())
	assert.Equal(t, "32.00", timeSeries2[2].Balance.Amount)
	assert.Equal(t, 32., timeSeries2[2].Balance.amount)
	assert.Equal(t, false, timeSeries2[2].Balance.Negative)
	assert.Equal(t, int64(3), timeSeries2[2].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 3, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[3].Time.UnixNano())
	assert.Equal(t, "28.00", timeSeries2[3].Balance.Amount)
	assert.Equal(t, 28., timeSeries2[3].Balance.amount)
	assert.Equal(t, false, timeSeries2[3].Balance.Negative)
	assert.Equal(t, int64(-12), timeSeries2[3].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 4, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[4].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries2[4].Balance.Amount)
	assert.Equal(t, 25., timeSeries2[4].Balance.amount)
	assert.Equal(t, false, timeSeries2[4].Balance.Negative)
	assert.Equal(t, int64(-10), timeSeries2[4].Balance.Bonus)

	/******************************************************************************************************************************************************
		2. Not before time is -5 hours. startDateIsBeforeEndDate = true.
	******************************************************************************************************************************************************/
	notBeforeTime = time.New(now.Add(-5 * repo.cfg.GlobalAggregationInterval.Child))
	startDateIsBeforeEndDate = true

	entries = repo.processBalanceHistory(history, startDateIsBeforeEndDate, notBeforeTime)

	assert.Len(t, entries, 3)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 3, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[0].Time.UnixNano())
	assert.Equal(t, "50.00", entries[0].Balance.Amount)
	assert.Equal(t, 50., entries[0].Balance.amount)
	assert.Equal(t, false, entries[0].Balance.Negative)
	assert.Equal(t, int64(0), entries[0].Balance.Bonus)

	assert.Len(t, entries[1].TimeSeries, 0)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 3, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[0].Time.UnixNano())
	assert.Equal(t, "145.00", entries[1].Balance.Amount)
	assert.Equal(t, 145., entries[1].Balance.amount)
	assert.Equal(t, false, entries[1].Balance.Negative)
	assert.Equal(t, int64(190), entries[1].Balance.Bonus)

	assert.Len(t, entries[2].TimeSeries, 5)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 3, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[0].Time.UnixNano())
	assert.Equal(t, "141.00", entries[2].Balance.Amount)
	assert.Equal(t, 141., entries[2].Balance.amount)
	assert.Equal(t, false, entries[2].Balance.Negative)
	assert.Equal(t, int64(-2), entries[2].Balance.Bonus)

	timeSeries0 = entries[0].TimeSeries
	assert.Len(t, timeSeries0, 0)

	timeSeries1 = entries[1].TimeSeries
	assert.Len(t, timeSeries1, 0)

	timeSeries2 = entries[2].TimeSeries
	assert.Len(t, timeSeries2, 5)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 0, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[0].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries2[0].Balance.Amount)
	assert.Equal(t, 25., timeSeries2[0].Balance.amount)
	assert.Equal(t, false, timeSeries2[0].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries2[0].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 1, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[1].Time.UnixNano())
	assert.Equal(t, "31.00", timeSeries2[1].Balance.Amount)
	assert.Equal(t, 31., timeSeries2[1].Balance.amount)
	assert.Equal(t, false, timeSeries2[1].Balance.Negative)
	assert.Equal(t, int64(24), timeSeries2[1].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 2, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[2].Time.UnixNano())
	assert.Equal(t, "32.00", timeSeries2[2].Balance.Amount)
	assert.Equal(t, 32., timeSeries2[2].Balance.amount)
	assert.Equal(t, false, timeSeries2[2].Balance.Negative)
	assert.Equal(t, int64(3), timeSeries2[2].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 3, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[3].Time.UnixNano())
	assert.Equal(t, "28.00", timeSeries2[3].Balance.Amount)
	assert.Equal(t, 28., timeSeries2[3].Balance.amount)
	assert.Equal(t, false, timeSeries2[3].Balance.Negative)
	assert.Equal(t, int64(-12), timeSeries2[3].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 4, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[4].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries2[4].Balance.Amount)
	assert.Equal(t, 25., timeSeries2[4].Balance.amount)
	assert.Equal(t, false, timeSeries2[4].Balance.Negative)
	assert.Equal(t, int64(-10), timeSeries2[4].Balance.Bonus)

	/******************************************************************************************************************************************************
		3. Not before time is -5 hours. startDateIsBeforeEndDate = false.
	******************************************************************************************************************************************************/
	notBeforeTime = time.New(now.Add(-5 * repo.cfg.GlobalAggregationInterval.Child))
	startDateIsBeforeEndDate = false

	entries = repo.processBalanceHistory(history, startDateIsBeforeEndDate, notBeforeTime)
	assert.Len(t, entries, 3)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[0].Time.UnixNano())
	assert.Equal(t, "141.00", entries[0].Balance.Amount)
	assert.Equal(t, 141., entries[0].Balance.amount)
	assert.Equal(t, false, entries[0].Balance.Negative)
	assert.Equal(t, int64(-2), entries[0].Balance.Bonus)

	assert.Len(t, entries[1].TimeSeries, 0)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[1].Time.UnixNano())
	assert.Equal(t, "145.00", entries[1].Balance.Amount)
	assert.Equal(t, 145., entries[1].Balance.amount)
	assert.Equal(t, false, entries[1].Balance.Negative)
	assert.Equal(t, int64(190), entries[1].Balance.Bonus)

	assert.Len(t, entries[2].TimeSeries, 0)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 3, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[2].Time.UnixNano())
	assert.Equal(t, "50.00", entries[2].Balance.Amount)
	assert.Equal(t, 50., entries[2].Balance.amount)
	assert.Equal(t, false, entries[2].Balance.Negative)
	assert.Equal(t, int64(0), entries[2].Balance.Bonus)

	timeSeries0 = entries[0].TimeSeries
	assert.Len(t, timeSeries0, 5)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 4, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries0[0].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries0[0].Balance.Amount)
	assert.Equal(t, 25., timeSeries0[0].Balance.amount)
	assert.Equal(t, false, timeSeries0[0].Balance.Negative)
	assert.Equal(t, int64(-10), timeSeries0[0].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 3, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries0[1].Time.UnixNano())
	assert.Equal(t, "28.00", timeSeries0[1].Balance.Amount)
	assert.Equal(t, 28., timeSeries0[1].Balance.amount)
	assert.Equal(t, false, timeSeries0[1].Balance.Negative)
	assert.Equal(t, int64(-12), timeSeries0[1].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 2, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries0[2].Time.UnixNano())
	assert.Equal(t, "32.00", timeSeries0[2].Balance.Amount)
	assert.Equal(t, 32., timeSeries0[2].Balance.amount)
	assert.Equal(t, false, timeSeries0[2].Balance.Negative)
	assert.Equal(t, int64(3), timeSeries0[2].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 1, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries0[3].Time.UnixNano())
	assert.Equal(t, "31.00", timeSeries0[3].Balance.Amount)
	assert.Equal(t, 31., timeSeries0[3].Balance.amount)
	assert.Equal(t, false, timeSeries0[3].Balance.Negative)
	assert.Equal(t, int64(24), timeSeries0[3].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 0, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries0[4].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries0[4].Balance.Amount)
	assert.Equal(t, 25., timeSeries0[4].Balance.amount)
	assert.Equal(t, false, timeSeries0[4].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries0[4].Balance.Bonus)

	timeSeries1 = entries[1].TimeSeries
	assert.Len(t, timeSeries1, 0)

	timeSeries2 = entries[2].TimeSeries
	assert.Len(t, timeSeries2, 0)

	/******************************************************************************************************************************************************
		4. Not before time is -30 hours. startDateIsBeforeEndDate = true.
	******************************************************************************************************************************************************/
	notBeforeTime = time.New(now.Add(-30 * repo.cfg.GlobalAggregationInterval.Child))
	startDateIsBeforeEndDate = true

	entries = repo.processBalanceHistory(history, startDateIsBeforeEndDate, notBeforeTime)
	assert.Len(t, entries, 3)
	assert.Len(t, entries[0].TimeSeries, 1)
	assert.Len(t, entries[1].TimeSeries, 24)
	assert.Len(t, entries[2].TimeSeries, 5)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 3, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[0].Time.UnixNano())
	assert.Equal(t, "50.00", entries[0].Balance.Amount)
	assert.Equal(t, 50., entries[0].Balance.amount)
	assert.Equal(t, false, entries[0].Balance.Negative)
	assert.Equal(t, int64(0), entries[0].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[1].Time.UnixNano())
	assert.Equal(t, "145.00", entries[1].Balance.Amount)
	assert.Equal(t, 145., entries[1].Balance.amount)
	assert.Equal(t, false, entries[1].Balance.Negative)
	assert.Equal(t, int64(190), entries[1].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 0, 0, 0, 0, stdlibtime.UTC)).UnixNano(), entries[2].Time.UnixNano())
	assert.Equal(t, "141.00", entries[2].Balance.Amount)
	assert.Equal(t, 141., entries[2].Balance.amount)
	assert.Equal(t, false, entries[2].Balance.Negative)
	assert.Equal(t, int64(-2), entries[2].Balance.Bonus)

	timeSeries0 = entries[0].TimeSeries
	assert.Len(t, timeSeries0, 1)
	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 3, 23, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries0[0].Time.UnixNano())
	assert.Equal(t, "1.00", timeSeries0[0].Balance.Amount)
	assert.Equal(t, 1., timeSeries0[0].Balance.amount)
	assert.Equal(t, false, timeSeries0[0].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries0[0].Balance.Bonus)

	timeSeries1 = entries[1].TimeSeries
	assert.Len(t, timeSeries1, 24)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 0, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[0].Time.UnixNano())
	assert.Equal(t, "27.00", timeSeries1[0].Balance.Amount)
	assert.Equal(t, 27., timeSeries1[0].Balance.amount)
	assert.Equal(t, false, timeSeries1[0].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries1[0].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 1, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[1].Time.UnixNano())
	assert.Equal(t, "29.00", timeSeries1[1].Balance.Amount)
	assert.Equal(t, 29., timeSeries1[1].Balance.amount)
	assert.Equal(t, false, timeSeries1[1].Balance.Negative)
	assert.Equal(t, int64(7), timeSeries1[1].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 2, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[2].Time.UnixNano())
	assert.Equal(t, "32.00", timeSeries1[2].Balance.Amount)
	assert.Equal(t, -32., timeSeries1[2].Balance.amount)
	assert.Equal(t, true, timeSeries1[2].Balance.Negative)
	assert.Equal(t, int64(-210), timeSeries1[2].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 3, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[3].Time.UnixNano())
	assert.Equal(t, "30.00", timeSeries1[3].Balance.Amount)
	assert.Equal(t, 30., timeSeries1[3].Balance.amount)
	assert.Equal(t, false, timeSeries1[3].Balance.Negative)
	assert.Equal(t, int64(193), timeSeries1[3].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 4, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[4].Time.UnixNano())
	assert.Equal(t, "15.00", timeSeries1[4].Balance.Amount)
	assert.Equal(t, 15., timeSeries1[4].Balance.amount)
	assert.Equal(t, false, timeSeries1[4].Balance.Negative)
	assert.Equal(t, int64(-50), timeSeries1[4].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 5, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[5].Time.UnixNano())
	assert.Equal(t, "21.00", timeSeries1[5].Balance.Amount)
	assert.Equal(t, -21., timeSeries1[5].Balance.amount)
	assert.Equal(t, true, timeSeries1[5].Balance.Negative)
	assert.Equal(t, int64(-240), timeSeries1[5].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 6, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[6].Time.UnixNano())
	assert.Equal(t, "22.00", timeSeries1[6].Balance.Amount)
	assert.Equal(t, 22., timeSeries1[6].Balance.amount)
	assert.Equal(t, false, timeSeries1[6].Balance.Negative)
	assert.Equal(t, int64(204), timeSeries1[6].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 7, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[7].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries1[7].Balance.Amount)
	assert.Equal(t, -25., timeSeries1[7].Balance.amount)
	assert.Equal(t, true, timeSeries1[7].Balance.Negative)
	assert.Equal(t, int64(-213), timeSeries1[7].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 8, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[8].Time.UnixNano())
	assert.Equal(t, "28.00", timeSeries1[8].Balance.Amount)
	assert.Equal(t, 28., timeSeries1[8].Balance.amount)
	assert.Equal(t, false, timeSeries1[8].Balance.Negative)
	assert.Equal(t, int64(212), timeSeries1[8].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 9, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[9].Time.UnixNano())
	assert.Equal(t, "30.00", timeSeries1[9].Balance.Amount)
	assert.Equal(t, 30., timeSeries1[9].Balance.amount)
	assert.Equal(t, false, timeSeries1[9].Balance.Negative)
	assert.Equal(t, int64(7), timeSeries1[9].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 10, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[10].Time.UnixNano())
	assert.Equal(t, "10.00", timeSeries1[10].Balance.Amount)
	assert.Equal(t, -10., timeSeries1[10].Balance.amount)
	assert.Equal(t, true, timeSeries1[10].Balance.Negative)
	assert.Equal(t, int64(-133), timeSeries1[10].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 11, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[11].Time.UnixNano())
	assert.Equal(t, "15.00", timeSeries1[11].Balance.Amount)
	assert.Equal(t, 15., timeSeries1[11].Balance.amount)
	assert.Equal(t, false, timeSeries1[11].Balance.Negative)
	assert.Equal(t, int64(250), timeSeries1[11].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 12, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[12].Time.UnixNano())
	assert.Equal(t, "8.00", timeSeries1[12].Balance.Amount)
	assert.Equal(t, -8., timeSeries1[12].Balance.amount)
	assert.Equal(t, true, timeSeries1[12].Balance.Negative)
	assert.Equal(t, int64(-153), timeSeries1[12].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 13, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[13].Time.UnixNano())
	assert.Equal(t, "10.00", timeSeries1[13].Balance.Amount)
	assert.Equal(t, 10., timeSeries1[13].Balance.amount)
	assert.Equal(t, false, timeSeries1[13].Balance.Negative)
	assert.Equal(t, int64(225), timeSeries1[13].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 14, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[14].Time.UnixNano())
	assert.Equal(t, "27.00", timeSeries1[14].Balance.Amount)
	assert.Equal(t, -27., timeSeries1[14].Balance.amount)
	assert.Equal(t, true, timeSeries1[14].Balance.Negative)
	assert.Equal(t, int64(-370), timeSeries1[14].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 15, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[15].Time.UnixNano())
	assert.Equal(t, "31.00", timeSeries1[15].Balance.Amount)
	assert.Equal(t, 31., timeSeries1[15].Balance.amount)
	assert.Equal(t, false, timeSeries1[15].Balance.Negative)
	assert.Equal(t, int64(214), timeSeries1[15].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 16, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[16].Time.UnixNano())
	assert.Equal(t, "30.00", timeSeries1[16].Balance.Amount)
	assert.Equal(t, 30., timeSeries1[16].Balance.amount)
	assert.Equal(t, false, timeSeries1[16].Balance.Negative)
	assert.Equal(t, int64(-3), timeSeries1[16].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 17, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[17].Time.UnixNano())
	assert.Equal(t, "27.00", timeSeries1[17].Balance.Amount)
	assert.Equal(t, 27., timeSeries1[17].Balance.amount)
	assert.Equal(t, false, timeSeries1[17].Balance.Negative)
	assert.Equal(t, int64(-10), timeSeries1[17].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 18, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[18].Time.UnixNano())
	assert.Equal(t, "29.00", timeSeries1[18].Balance.Amount)
	assert.Equal(t, -29., timeSeries1[18].Balance.amount)
	assert.Equal(t, true, timeSeries1[18].Balance.Negative)
	assert.Equal(t, int64(-207), timeSeries1[18].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 19, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[19].Time.UnixNano())
	assert.Equal(t, "30.00", timeSeries1[19].Balance.Amount)
	assert.Equal(t, 30., timeSeries1[19].Balance.amount)
	assert.Equal(t, false, timeSeries1[19].Balance.Negative)
	assert.Equal(t, int64(203), timeSeries1[19].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 20, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[20].Time.UnixNano())
	assert.Equal(t, "15.00", timeSeries1[20].Balance.Amount)
	assert.Equal(t, -15., timeSeries1[20].Balance.amount)
	assert.Equal(t, true, timeSeries1[20].Balance.Negative)
	assert.Equal(t, int64(-150), timeSeries1[20].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 21, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[21].Time.UnixNano())
	assert.Equal(t, "15.00", timeSeries1[21].Balance.Amount)
	assert.Equal(t, -15., timeSeries1[21].Balance.amount)
	assert.Equal(t, true, timeSeries1[21].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries1[21].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 22, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[22].Time.UnixNano())
	assert.Equal(t, "20.00", timeSeries1[22].Balance.Amount)
	assert.Equal(t, 20., timeSeries1[22].Balance.amount)
	assert.Equal(t, false, timeSeries1[22].Balance.Negative)
	assert.Equal(t, int64(233), timeSeries1[22].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 4, 23, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries1[23].Time.UnixNano())
	assert.Equal(t, "17.00", timeSeries1[23].Balance.Amount)
	assert.Equal(t, -17., timeSeries1[23].Balance.amount)
	assert.Equal(t, true, timeSeries1[23].Balance.Negative)
	assert.Equal(t, int64(-185), timeSeries1[23].Balance.Bonus)

	timeSeries2 = entries[2].TimeSeries
	assert.Len(t, timeSeries2, 5)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 0, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[0].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries2[0].Balance.Amount)
	assert.Equal(t, 25., timeSeries2[0].Balance.amount)
	assert.Equal(t, false, timeSeries2[0].Balance.Negative)
	assert.Equal(t, int64(0), timeSeries2[0].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 1, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[1].Time.UnixNano())
	assert.Equal(t, "31.00", timeSeries2[1].Balance.Amount)
	assert.Equal(t, 31., timeSeries2[1].Balance.amount)
	assert.Equal(t, false, timeSeries2[1].Balance.Negative)
	assert.Equal(t, int64(24), timeSeries2[1].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 2, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[2].Time.UnixNano())
	assert.Equal(t, "32.00", timeSeries2[2].Balance.Amount)
	assert.Equal(t, 32., timeSeries2[2].Balance.amount)
	assert.Equal(t, false, timeSeries2[2].Balance.Negative)
	assert.Equal(t, int64(3), timeSeries2[2].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 3, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[3].Time.UnixNano())
	assert.Equal(t, "28.00", timeSeries2[3].Balance.Amount)
	assert.Equal(t, 28., timeSeries2[3].Balance.amount)
	assert.Equal(t, false, timeSeries2[3].Balance.Negative)
	assert.Equal(t, int64(-12), timeSeries2[3].Balance.Bonus)

	assert.EqualValues(t, time.New(stdlibtime.Date(2023, 6, 5, 4, 15, 10, 1, stdlibtime.UTC)).UnixNano(), timeSeries2[4].Time.UnixNano())
	assert.Equal(t, "25.00", timeSeries2[4].Balance.Amount)
	assert.Equal(t, 25., timeSeries2[4].Balance.amount)
	assert.Equal(t, false, timeSeries2[4].Balance.Negative)
	assert.Equal(t, int64(-10), timeSeries2[4].Balance.Bonus)
}
