// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	"github.com/ice-blockchain/wintr/time"
)

func TestCalculateDates_Limit24_Offset0_Factor1(t *testing.T) {
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

	limit := uint64(24)
	offset := uint64(0)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	end := time.New(stdlibtime.Date(2023, 6, 7, 5, 15, 10, 1, stdlibtime.UTC))
	factor := stdlibtime.Duration(1)

	dates, notBeforeTime, notAfterTime := repo.calculateDates(limit, offset, start, end, factor)
	assert.Len(t, dates, 48)

	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 6, 5, 15, 10, 1, stdlibtime.UTC)), notAfterTime)

	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
		expectedStart.Add(-5 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-1 * repo.cfg.GlobalAggregationInterval.Child),
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
	}
	assert.EqualValues(t, expected, dates)

}

func TestCalculateDates_Limit12_Offset0_Factor1(t *testing.T) {
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
	limit := uint64(12)
	offset := uint64(0)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	end := time.New(stdlibtime.Date(2023, 6, 7, 5, 15, 10, 1, stdlibtime.UTC))
	factor := stdlibtime.Duration(1)

	dates, notBeforeTime, notAfterTime := repo.calculateDates(limit, offset, start, end, factor)
	assert.Len(t, dates, 24)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC)), notAfterTime) // Cuz calculated limit is 0.

	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
		expectedStart.Add(-5 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-1 * repo.cfg.GlobalAggregationInterval.Child),
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
	}
	assert.EqualValues(t, expected, dates)
}

func TestCalculateDates_Limit36_Offset0_Factor1(t *testing.T) {
	repo := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}

	limit := uint64(36)
	offset := uint64(0)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	end := time.New(stdlibtime.Date(2023, 6, 7, 5, 15, 10, 1, stdlibtime.UTC))
	factor := stdlibtime.Duration(1)

	dates, notBeforeTime, notAfterTime := repo.calculateDates(limit, offset, start, end, factor)
	assert.Len(t, dates, 48)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 6, 5, 15, 10, 1, stdlibtime.UTC)), notAfterTime)

	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
		expectedStart.Add(-5 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-1 * repo.cfg.GlobalAggregationInterval.Child),
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
	}
	assert.EqualValues(t, expected, dates)
}

func TestCalculateDates_Limit48_Offset0_Factor1(t *testing.T) {
	repo := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	limit := uint64(48)
	offset := uint64(0)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	end := time.New(stdlibtime.Date(2023, 6, 6, 5, 15, 10, 1, stdlibtime.UTC))
	factor := stdlibtime.Duration(1)

	dates, notBeforeTime, notAfterTime := repo.calculateDates(limit, offset, start, end, factor)
	assert.Len(t, dates, 72)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 6, 5, 15, 10, 1, stdlibtime.UTC)), notAfterTime)

	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
		expectedStart.Add(-5 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-4 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-3 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-2 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-1 * repo.cfg.GlobalAggregationInterval.Child),
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
		expectedStart.Add(54 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(55 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(56 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(57 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(58 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(59 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(60 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(61 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(62 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(63 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(64 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(65 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(66 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)
}

func TestCalculateDates_Limit24_Offset0_FactorMinus1(t *testing.T) {
	repo := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	limit := uint64(24)
	offset := uint64(0)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	var end *time.Time
	factor := stdlibtime.Duration(-1)

	dates, notBeforeTime, notAfterTime := repo.calculateDates(limit, offset, start, end, factor)
	assert.Len(t, dates, 30)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 4, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	assert.Equal(t, start, notAfterTime)

	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
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
}

func TestCalculateDates_Limit24_Offset24_FactorMinus1(t *testing.T) {
	repo := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	limit := uint64(24)
	offset := uint64(24)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	var end *time.Time
	factor := stdlibtime.Duration(-1)

	dates, notBeforeTime, notAfterTime := repo.calculateDates(limit, offset, start, end, factor)
	assert.Len(t, dates, 30)

	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 3, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 4, 5, 15, 10, 1, stdlibtime.UTC)), notAfterTime)

	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
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

func TestCalculateDates_Limit24_Offset24_Factor1(t *testing.T) {
	repo := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	limit := uint64(24)
	offset := uint64(24)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	end := time.New(stdlibtime.Date(2023, 6, 7, 5, 15, 10, 1, stdlibtime.UTC))
	factor := stdlibtime.Duration(1)

	dates, notBeforeTime, notAfterTime := repo.calculateDates(limit, offset, start, end, factor)
	assert.Len(t, dates, 48)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 6, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 7, 5, 15, 10, 1, stdlibtime.UTC)), notAfterTime)
	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
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
		expectedStart.Add(54 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(55 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(56 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(57 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(58 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(59 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(60 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(61 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(62 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(63 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(64 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(65 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(66 * repo.cfg.GlobalAggregationInterval.Child),
	}
	assert.EqualValues(t, expected, dates)
}

func TestCalculateDates_Limit48_Offset48_FactorMinus1(t *testing.T) {
	repo := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}

	limit := uint64(48)
	offset := uint64(48)
	start := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	end := time.New(stdlibtime.Date(2023, 6, 5, 5, 15, 10, 1, stdlibtime.UTC))
	factor := stdlibtime.Duration(-1)

	dates, notBeforeTime, notAfterTime := repo.calculateDates(limit, offset, start, end, factor)
	assert.Len(t, dates, 54)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 1, 5, 15, 10, 1, stdlibtime.UTC)), notBeforeTime)
	assert.Equal(t, time.New(stdlibtime.Date(2023, 6, 3, 5, 15, 10, 1, stdlibtime.UTC)), notAfterTime)
	expectedStart := time.New(stdlibtime.Date(2023, 6, 5, 5, 0, 0, 0, stdlibtime.UTC))
	expected := []stdlibtime.Time{
		expectedStart.Add(-48 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-49 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-50 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-51 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-52 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-53 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-54 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-55 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-56 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-57 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-58 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-59 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-60 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-61 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-62 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-63 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-64 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-65 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-66 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-67 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-68 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-69 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-70 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-71 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-72 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-73 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-74 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-75 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-76 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-77 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-78 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-79 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-80 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-81 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-82 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-83 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-84 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-85 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-86 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-87 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-88 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-89 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-90 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-91 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-92 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-93 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-94 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-95 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-96 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-97 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-98 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-99 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-100 * repo.cfg.GlobalAggregationInterval.Child),
		expectedStart.Add(-101 * repo.cfg.GlobalAggregationInterval.Child),
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
		2. Not before time is -10 hours. Not after time = now. startDateIsBeforeEndDate = true.
	******************************************************************************************************************************************************/
	notBeforeTime := time.New(now.Add(-10 * repo.cfg.GlobalAggregationInterval.Child))
	notAfterTime := now
	startDateIsBeforeEndDate := true

	entries := repo.processBalanceHistory(history, startDateIsBeforeEndDate, notBeforeTime, notAfterTime)
	expected := []*BalanceHistoryEntry{
		{
			Time: stdlibtime.Date(2023, 6, 4, 0, 0, 0, 0, stdlibtime.UTC),
			Balance: &BalanceHistoryBalanceDiff{
				amount:   145.,
				Amount:   "145.00",
				Bonus:    int64(190),
				Negative: false,
			},
			TimeSeries: []*BalanceHistoryEntry{
				{
					Time: stdlibtime.Date(2023, 6, 4, 19, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(0),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 4, 20, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -15.,
						Amount:   "15.00",
						Negative: true,
						Bonus:    int64(-150),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 4, 21, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -15.,
						Amount:   "15.00",
						Negative: true,
						Bonus:    int64(0),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 4, 22, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   20.,
						Amount:   "20.00",
						Negative: false,
						Bonus:    int64(233),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 4, 23, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -17.,
						Amount:   "17.00",
						Negative: true,
						Bonus:    int64(-185),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
			},
		},
		{
			Time: *time.New(stdlibtime.Date(2023, 6, 5, 0, 0, 0, 0, stdlibtime.UTC)).Time,
			Balance: &BalanceHistoryBalanceDiff{
				amount:   141.,
				Amount:   "141.00",
				Bonus:    int64(-2),
				Negative: false,
			},
			TimeSeries: []*BalanceHistoryEntry{
				{
					Time: stdlibtime.Date(2023, 6, 5, 0, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   25.,
						Amount:   "25.00",
						Negative: false,
						Bonus:    int64(247),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 1, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   31.,
						Amount:   "31.00",
						Negative: false,
						Bonus:    int64(24),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 2, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   32.,
						Amount:   "32.00",
						Negative: false,
						Bonus:    int64(3),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 3, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   28.,
						Amount:   "28.00",
						Negative: false,
						Bonus:    int64(-12),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 4, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   25.,
						Amount:   "25.00",
						Negative: false,
						Bonus:    int64(-10),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
			},
		},
	}
	assert.EqualValues(t, expected, entries)

	/******************************************************************************************************************************************************
		3. Not before time is -5 hours. Not after time = now. startDateIsBeforeEndDate = true.
	******************************************************************************************************************************************************/
	notBeforeTime = time.New(now.Add(-5 * repo.cfg.GlobalAggregationInterval.Child))
	notAfterTime = now
	startDateIsBeforeEndDate = true

	entries = repo.processBalanceHistory(history, startDateIsBeforeEndDate, notBeforeTime, notAfterTime)

	expected = []*BalanceHistoryEntry{
		{
			Time: *time.New(stdlibtime.Date(2023, 6, 5, 0, 0, 0, 0, stdlibtime.UTC)).Time,
			Balance: &BalanceHistoryBalanceDiff{
				amount:   141.,
				Amount:   "141.00",
				Bonus:    int64(-2),
				Negative: false,
			},
			TimeSeries: []*BalanceHistoryEntry{
				{
					Time: stdlibtime.Date(2023, 6, 5, 0, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   25.,
						Amount:   "25.00",
						Negative: false,
						Bonus:    int64(0),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 1, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   31.,
						Amount:   "31.00",
						Negative: false,
						Bonus:    int64(24),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 2, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   32.,
						Amount:   "32.00",
						Negative: false,
						Bonus:    int64(3),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 3, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   28.,
						Amount:   "28.00",
						Negative: false,
						Bonus:    int64(-12),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 4, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   25.,
						Amount:   "25.00",
						Negative: false,
						Bonus:    int64(-10),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
			},
		},
	}
	assert.EqualValues(t, expected, entries)
}

func TestProcessBalanceHistory_ChildIsHour_TimeGrow(t *testing.T) {
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
			CreatedAt:           time.New(now.Add(1 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  25.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(2 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  28.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  32.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(4 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  31.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  25.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(6 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 17.,
		},
		{
			CreatedAt:           time.New(now.Add(7 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  20.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(8 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 15.,
		},
		{
			CreatedAt:           time.New(now.Add(9 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 15.,
		},
		{
			CreatedAt:           time.New(now.Add(10 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  30.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(11 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 29.,
		},
		{
			CreatedAt:           time.New(now.Add(12 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  27.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(13 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  30.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(14 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  31.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(15 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 27.,
		},
		{
			CreatedAt:           time.New(now.Add(16 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  10.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(17 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 8.,
		},
		{
			CreatedAt:           time.New(now.Add(18 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  15.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(19 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 10.,
		},
		{
			CreatedAt:           time.New(now.Add(20 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  30.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(21 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  28.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(22 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 25.,
		},
		{
			CreatedAt:           time.New(now.Add(23 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  22.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(24 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 21.,
		},
		{
			CreatedAt:           time.New(now.Add(25 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  15.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(26 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  30.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(27 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 32.,
		},
		{
			CreatedAt:           time.New(now.Add(28 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  29.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(29 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  27.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(30 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  1.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(31 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 25.,
		},
		{
			CreatedAt:           time.New(now.Add(32 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  20.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(33 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  32.,
			BalanceTotalSlashed: 0.,
		},
		{
			CreatedAt:           time.New(now.Add(34 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  0.,
			BalanceTotalSlashed: 10.,
		},
		{
			CreatedAt:           time.New(now.Add(35 * repo.cfg.GlobalAggregationInterval.Child)),
			BalanceTotalMinted:  32.,
			BalanceTotalSlashed: 0.,
		},
	}

	/******************************************************************************************************************************************************
		2. Not before time is now. Not after time = +24 hours. startDateIsBeforeEndDate = false.
	******************************************************************************************************************************************************/
	notBeforeTime := now
	notAfterTime := time.New(now.Add(30 * repo.cfg.GlobalAggregationInterval.Child))
	startDateIsBeforeEndDate := true

	entries := repo.processBalanceHistory(history, startDateIsBeforeEndDate, notBeforeTime, notAfterTime)
	expected := []*BalanceHistoryEntry{
		{
			Time: stdlibtime.Date(2023, 6, 5, 0, 0, 0, 0, stdlibtime.UTC),
			Balance: &BalanceHistoryBalanceDiff{
				amount:   193.,
				Amount:   "193.00",
				Bonus:    int64(0),
				Negative: false,
			},
			TimeSeries: []*BalanceHistoryEntry{
				{
					Time: stdlibtime.Date(2023, 6, 5, 6, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   25.,
						Amount:   "25.00",
						Negative: false,
						Bonus:    int64(0),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 7, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   28.,
						Amount:   "28.00",
						Negative: false,
						Bonus:    int64(12),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 8, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   32.,
						Amount:   "32.00",
						Negative: false,
						Bonus:    int64(14),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 9, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   31.,
						Amount:   "31.00",
						Negative: false,
						Bonus:    int64(-3),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 10, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   25.,
						Amount:   "25.00",
						Negative: false,
						Bonus:    int64(-19),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 11, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -17.,
						Amount:   "17.00",
						Negative: true,
						Bonus:    int64(-168),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 12, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   20.,
						Amount:   "20.00",
						Negative: false,
						Bonus:    int64(217),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 13, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -15.,
						Amount:   "15.00",
						Negative: true,
						Bonus:    int64(-175),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 14, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -15.,
						Amount:   "15.00",
						Negative: true,
						Bonus:    int64(0),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 15, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(300),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 16, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -29.,
						Amount:   "29.00",
						Negative: true,
						Bonus:    int64(-196),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 17, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   27.,
						Amount:   "27.00",
						Negative: false,
						Bonus:    int64(193),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 18, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(11),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 19, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   31.,
						Amount:   "31.00",
						Negative: false,
						Bonus:    int64(3),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 20, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -27.,
						Amount:   "27.00",
						Negative: true,
						Bonus:    int64(-187),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 21, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   10.,
						Amount:   "10.00",
						Negative: false,
						Bonus:    int64(137),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 22, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -8.,
						Amount:   "8.00",
						Negative: true,
						Bonus:    int64(-180),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 23, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   15.,
						Amount:   "15.00",
						Negative: false,
						Bonus:    int64(287),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
			},
		},
		{
			Time: *time.New(stdlibtime.Date(2023, 6, 6, 0, 0, 0, 0, stdlibtime.UTC)).Time,
			Balance: &BalanceHistoryBalanceDiff{
				amount:   143.,
				Amount:   "143.00",
				Bonus:    int64(-25),
				Negative: false,
			},
			TimeSeries: []*BalanceHistoryEntry{
				{
					Time: stdlibtime.Date(2023, 6, 6, 0, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -10.,
						Amount:   "10.00",
						Negative: true,
						Bonus:    int64(-166),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 1, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(400),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 2, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   28.,
						Amount:   "28.00",
						Negative: false,
						Bonus:    int64(-6),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 3, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -25.,
						Amount:   "25.00",
						Negative: true,
						Bonus:    int64(-189),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 4, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   22.,
						Amount:   "22.00",
						Negative: false,
						Bonus:    int64(188),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 5, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -21.,
						Amount:   "21.00",
						Negative: true,
						Bonus:    int64(-195),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 6, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   15.,
						Amount:   "15.00",
						Negative: false,
						Bonus:    int64(171),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 7, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(100),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 8, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -32.,
						Amount:   "32.00",
						Negative: true,
						Bonus:    int64(-206),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 9, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   29.,
						Amount:   "29.00",
						Negative: false,
						Bonus:    int64(190),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 10, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   27.,
						Amount:   "27.00",
						Negative: false,
						Bonus:    int64(-6),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 11, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   1.,
						Amount:   "1.00",
						Negative: false,
						Bonus:    int64(-96),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
			},
		},
	}
	assert.EqualValues(t, expected, entries)

	startDateIsBeforeEndDate = false
	entries = repo.processBalanceHistory(history, startDateIsBeforeEndDate, notBeforeTime, notAfterTime)
	expected = []*BalanceHistoryEntry{
		{
			Time: *time.New(stdlibtime.Date(2023, 6, 6, 0, 0, 0, 0, stdlibtime.UTC)).Time,
			Balance: &BalanceHistoryBalanceDiff{
				amount:   143.,
				Amount:   "143.00",
				Bonus:    int64(-25),
				Negative: false,
			},
			TimeSeries: []*BalanceHistoryEntry{
				{
					Time: stdlibtime.Date(2023, 6, 6, 11, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   1.,
						Amount:   "1.00",
						Negative: false,
						Bonus:    int64(-96),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 10, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   27.,
						Amount:   "27.00",
						Negative: false,
						Bonus:    int64(-6),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 9, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   29.,
						Amount:   "29.00",
						Negative: false,
						Bonus:    int64(190),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 8, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -32.,
						Amount:   "32.00",
						Negative: true,
						Bonus:    int64(-206),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 7, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(100),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 6, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   15.,
						Amount:   "15.00",
						Negative: false,
						Bonus:    int64(171),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 5, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -21.,
						Amount:   "21.00",
						Negative: true,
						Bonus:    int64(-195),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 4, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   22.,
						Amount:   "22.00",
						Negative: false,
						Bonus:    int64(188),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 3, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -25.,
						Amount:   "25.00",
						Negative: true,
						Bonus:    int64(-189),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 2, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   28.,
						Amount:   "28.00",
						Negative: false,
						Bonus:    int64(-6),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 1, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(400),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 6, 0, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -10.,
						Amount:   "10.00",
						Negative: true,
						Bonus:    int64(-166),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
			},
		},
		{
			Time: stdlibtime.Date(2023, 6, 5, 0, 0, 0, 0, stdlibtime.UTC),
			Balance: &BalanceHistoryBalanceDiff{
				amount:   193.,
				Amount:   "193.00",
				Bonus:    int64(0),
				Negative: false,
			},
			TimeSeries: []*BalanceHistoryEntry{
				{
					Time: stdlibtime.Date(2023, 6, 5, 23, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   15.,
						Amount:   "15.00",
						Negative: false,
						Bonus:    int64(287),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 22, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -8.,
						Amount:   "8.00",
						Negative: true,
						Bonus:    int64(-180),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 21, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   10.,
						Amount:   "10.00",
						Negative: false,
						Bonus:    int64(137),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 20, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -27.,
						Amount:   "27.00",
						Negative: true,
						Bonus:    int64(-187),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 19, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   31.,
						Amount:   "31.00",
						Negative: false,
						Bonus:    int64(3),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 18, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(11),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 17, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   27.,
						Amount:   "27.00",
						Negative: false,
						Bonus:    int64(193),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 16, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -29.,
						Amount:   "29.00",
						Negative: true,
						Bonus:    int64(-196),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 15, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   30.,
						Amount:   "30.00",
						Negative: false,
						Bonus:    int64(300),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 14, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -15.,
						Amount:   "15.00",
						Negative: true,
						Bonus:    int64(0),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 13, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -15.,
						Amount:   "15.00",
						Negative: true,
						Bonus:    int64(-175),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 12, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   20.,
						Amount:   "20.00",
						Negative: false,
						Bonus:    int64(217),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 11, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   -17.,
						Amount:   "17.00",
						Negative: true,
						Bonus:    int64(-168),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 10, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   25.,
						Amount:   "25.00",
						Negative: false,
						Bonus:    int64(-19),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 9, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   31.,
						Amount:   "31.00",
						Negative: false,
						Bonus:    int64(-3),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 8, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   32.,
						Amount:   "32.00",
						Negative: false,
						Bonus:    int64(14),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 7, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   28.,
						Amount:   "28.00",
						Negative: false,
						Bonus:    int64(12),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
				{
					Time: stdlibtime.Date(2023, 6, 5, 6, 15, 10, 1, stdlibtime.UTC),
					Balance: &BalanceHistoryBalanceDiff{
						amount:   25.,
						Amount:   "25.00",
						Negative: false,
						Bonus:    int64(0),
					},
					TimeSeries: []*BalanceHistoryEntry{},
				},
			},
		},
	}
	assert.EqualValues(t, expected, entries)
}
