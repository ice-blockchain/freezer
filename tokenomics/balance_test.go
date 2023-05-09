// SPDX-License-Identifier: ice License 1.0
//go:build xxx

package tokenomics

import (
	"fmt"
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/time"
)

//nolint:funlen,dupl,maintidx // .
func TestProcessBalanceHistory_ChildIsHour_ParentIsDay_Minus30MinutesTimezone(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	utcOffset := -870 * stdlibtime.Minute // -14:30.
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	now := time.Now()
	now = time.New(now.Add(-users.NanosSinceMidnight(now)))
	adoptions := []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     nil,
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	preStakingSummaries := []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 10000,
	}}
	childFormat := repo.cfg.globalAggregationIntervalChildDateFormat()
	balances := []*balance{
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(4*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(5*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(40000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(15*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
	}
	expected := []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}},
	}}

	actual := repo.processBalanceHistory(balances, true, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}, {
		AchievedAt:     time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
		BaseMiningRate: coin.UnsafeParseAmount("4000000000"),
		Milestone:      3,
	}}
	preStakingSummaries = []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
			Years:      3,
			Allocation: 50,
		},
		Bonus: 100,
	}, {
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 200,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(192000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(192000000000),
			Negative: true,
			Bonus:    -4900,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(120000000000),
				Negative: true,
				Bonus:    -3100,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(96000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(96000000000),
				Negative: true,
				Bonus:    -2500,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    500,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: false,
				Bonus:    1700,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: true,
				Bonus:    -1900,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: true,
				Bonus:    -1900,
			},
		}, {
			Time: now.In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: false,
				Bonus:    1700,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)
}

//nolint:funlen,dupl,maintidx // .
func TestProcessBalanceHistory_ChildIsHour_ParentIsDay_Plus30MinutesTimezone(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	utcOffset := 870 * stdlibtime.Minute // +14:30.
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	now := time.Now()
	now = time.New(now.Add(-users.NanosSinceMidnight(now)))
	adoptions := []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     nil,
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	preStakingSummaries := []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 10000,
	}}
	childFormat := repo.cfg.globalAggregationIntervalChildDateFormat()
	balances := []*balance{
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(4*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(5*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(40000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(15*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
	}
	expected := []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}},
	}}
	actual := repo.processBalanceHistory(balances, true, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}, {
		AchievedAt:     time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
		BaseMiningRate: coin.UnsafeParseAmount("4000000000"),
		Milestone:      3,
	}}
	preStakingSummaries = []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
			Years:      3,
			Allocation: 50,
		},
		Bonus: 100,
	}, {
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 200,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location).Add(-30 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)
}

//nolint:funlen,dupl,maintidx // .
func TestProcessBalanceHistory_ChildIsHour_ParentIsDay_Minus45MinutesTimezone(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	utcOffset := -765 * stdlibtime.Minute // -12:45.
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	now := time.Now()
	now = time.New(now.Add(-users.NanosSinceMidnight(now)))
	adoptions := []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     nil,
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	preStakingSummaries := []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 10000,
	}}
	childFormat := repo.cfg.globalAggregationIntervalChildDateFormat()
	balances := []*balance{
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(4*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(5*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(40000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(15*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
	}
	expected := []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}},
	}}
	actual := repo.processBalanceHistory(balances, true, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}, {
		AchievedAt:     time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
		BaseMiningRate: coin.UnsafeParseAmount("4000000000"),
		Milestone:      3,
	}}
	preStakingSummaries = []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
			Years:      3,
			Allocation: 50,
		},
		Bonus: 100,
	}, {
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 200,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(192000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(192000000000),
			Negative: true,
			Bonus:    -4900,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(120000000000),
				Negative: true,
				Bonus:    -3100,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(96000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(96000000000),
				Negative: true,
				Bonus:    -2500,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    500,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: false,
				Bonus:    1700,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: true,
				Bonus:    -1900,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: true,
				Bonus:    -1900,
			},
		}, {
			Time: now.In(location).Add(-15 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: false,
				Bonus:    1700,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)
}

//nolint:funlen,dupl,maintidx // .
func TestProcessBalanceHistory_ChildIsHour_ParentIsDay_Plus45MinutesTimezone(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
	}}
	utcOffset := 765 * stdlibtime.Minute // +12:45.
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	now := time.Now()
	now = time.New(now.Add(-users.NanosSinceMidnight(now)))
	adoptions := []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     nil,
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	preStakingSummaries := []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 10000,
	}}
	childFormat := repo.cfg.globalAggregationIntervalChildDateFormat()
	balances := []*balance{
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(4*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(5*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(40000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(15*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
	}
	expected := []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}},
	}}
	actual := repo.processBalanceHistory(balances, true, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}, {
		AchievedAt:     time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
		BaseMiningRate: coin.UnsafeParseAmount("4000000000"),
		Milestone:      3,
	}}
	preStakingSummaries = []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
			Years:      3,
			Allocation: 50,
		},
		Bonus: 100,
	}, {
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Minute)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 200,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(64000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(64000000000),
			Negative: true,
			Bonus:    -500,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location).Add(-45 * stdlibtime.Minute),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)
}

//nolint:funlen,dupl,maintidx // .
func TestProcessBalanceHistory_ChildIsMinute_ParentIsHour_Minus30MinuteTimezone(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: stdlibtime.Hour,
			Child:  stdlibtime.Minute,
		},
	}}
	utcOffset := -870 * stdlibtime.Minute // -14:30.
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	now := time.Now()
	now = time.New(now.Add(-users.NanosSinceMidnight(now)))
	adoptions := []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     nil,
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	preStakingSummaries := []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 10000,
	}}
	childFormat := repo.cfg.globalAggregationIntervalChildDateFormat()
	balances := []*balance{
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(4*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(5*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(40000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(15*repo.cfg.GlobalAggregationInterval.Parent).Format(childFormat)),
			Negative:   true,
		},
	}
	expected := []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}},
	}, {
		Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(4040000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(4040000000000),
			Negative: true,
			Bonus:    -25350,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(4040000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(4040000000000),
				Negative: true,
				Bonus:    -25350,
			},
		}},
	}}
	actual := repo.processBalanceHistory(balances, true, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(4040000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(4040000000000),
			Negative: true,
			Bonus:    -50600,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(4040000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(4040000000000),
				Negative: true,
				Bonus:    -50600,
			},
		}},
	}, {
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}, {
		AchievedAt:     time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
		BaseMiningRate: coin.UnsafeParseAmount("4000000000"),
		Milestone:      3,
	}}
	preStakingSummaries = []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
			Years:      3,
			Allocation: 50,
		},
		Bonus: 100,
	}, {
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 200,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(120000000000),
			Negative: true,
			Bonus:    -3100,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(120000000000),
				Negative: true,
				Bonus:    -3100,
			},
		}},
	}, {
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(72000000000),
			Negative: true,
			Bonus:    -1900,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(96000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(96000000000),
				Negative: true,
				Bonus:    -2500,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    500,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: false,
				Bonus:    1700,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: true,
				Bonus:    -1900,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: true,
				Bonus:    -1900,
			},
		}, {
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: false,
				Bonus:    1700,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)
}

//nolint:funlen,dupl,maintidx // .
func TestProcessBalanceHistory_ChildIsMinute_ParentIsHour_Plus30MinuteTimezone(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: stdlibtime.Hour,
			Child:  stdlibtime.Minute,
		},
	}}
	utcOffset := 690 * stdlibtime.Minute // +11:30.
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	now := time.Now()
	now = time.New(now.Add(-users.NanosSinceMidnight(now)))
	adoptions := []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     nil,
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	preStakingSummaries := []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 10000,
	}}
	childFormat := repo.cfg.globalAggregationIntervalChildDateFormat()
	balances := []*balance{
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(4*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(5*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(40000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(15*repo.cfg.GlobalAggregationInterval.Parent).Format(childFormat)),
			Negative:   true,
		},
	}
	expected := []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}},
	}, {
		Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(40000000000),
			Negative: true,
			Bonus:    -350,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}},
	}}
	actual := repo.processBalanceHistory(balances, true, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.Add(15 * stdlibtime.Hour).In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(40000000000),
			Negative: true,
			Bonus:    -350,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}},
	}, {
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}, {
		AchievedAt:     time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
		BaseMiningRate: coin.UnsafeParseAmount("4000000000"),
		Milestone:      3,
	}}
	preStakingSummaries = []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
			Years:      3,
			Allocation: 50,
		},
		Bonus: 100,
	}, {
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 200,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(120000000000),
			Negative: true,
			Bonus:    -3100,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(120000000000),
				Negative: true,
				Bonus:    -3100,
			},
		}},
	}, {
		Time: now.In(location).Add(-30 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)
}

//nolint:funlen,dupl,maintidx // .
func TestProcessBalanceHistory_ChildIsMinute_ParentIsHour_Minus45MinuteTimezone(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: stdlibtime.Hour,
			Child:  stdlibtime.Minute,
		},
	}}
	utcOffset := -705 * stdlibtime.Minute // -11:45.
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	now := time.Now()
	now = time.New(now.Add(-users.NanosSinceMidnight(now)))
	adoptions := []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     nil,
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	preStakingSummaries := []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 10000,
	}}
	childFormat := repo.cfg.globalAggregationIntervalChildDateFormat()
	balances := []*balance{
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(4*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(5*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(40000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(15*repo.cfg.GlobalAggregationInterval.Parent).Format(childFormat)),
			Negative:   true,
		},
	}
	expected := []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}},
	}, {
		Time: now.In(location).Add(15 * repo.cfg.GlobalAggregationInterval.Parent).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(4040000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(4040000000000),
			Negative: true,
			Bonus:    -25350,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(4040000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(4040000000000),
				Negative: true,
				Bonus:    -25350,
			},
		}},
	}}
	actual := repo.processBalanceHistory(balances, true, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(15 * repo.cfg.GlobalAggregationInterval.Parent).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(4040000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(4040000000000),
			Negative: true,
			Bonus:    -50600,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(4040000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(4040000000000),
				Negative: true,
				Bonus:    -50600,
			},
		}},
	}, {
		Time: now.In(location).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}, {
		AchievedAt:     time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
		BaseMiningRate: coin.UnsafeParseAmount("4000000000"),
		Milestone:      3,
	}}
	preStakingSummaries = []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
			Years:      3,
			Allocation: 50,
		},
		Bonus: 100,
	}, {
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 200,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(15 * repo.cfg.GlobalAggregationInterval.Parent).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(120000000000),
			Negative: true,
			Bonus:    -3100,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(120000000000),
				Negative: true,
				Bonus:    -3100,
			},
		}},
	}, {
		Time: now.In(location).Add(-15 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(72000000000),
			Negative: true,
			Bonus:    -1900,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(96000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(96000000000),
				Negative: true,
				Bonus:    -2500,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    500,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: false,
				Bonus:    1700,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: true,
				Bonus:    -1900,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: true,
				Bonus:    -1900,
			},
		}, {
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(72000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(72000000000),
				Negative: false,
				Bonus:    1700,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)
}

//nolint:funlen,dupl,maintidx // .
func TestProcessBalanceHistory_ChildIsMinute_ParentIsHour_Plus45MinuteTimezone(t *testing.T) {
	t.Parallel()
	repo := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: stdlibtime.Hour,
			Child:  stdlibtime.Minute,
		},
	}}
	utcOffset := 705 * stdlibtime.Minute // +11:45.
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	now := time.Now()
	now = time.New(now.Add(-users.NanosSinceMidnight(now)))
	adoptions := []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     nil,
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	preStakingSummaries := []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 10000,
	}}
	childFormat := repo.cfg.globalAggregationIntervalChildDateFormat()
	balances := []*balance{
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(2*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(3*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(8000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(4*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   false,
		},
		{
			Amount:     coin.NewAmountUint64(32000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(5*repo.cfg.GlobalAggregationInterval.Child).Format(childFormat)),
			Negative:   true,
		},
		{
			Amount:     coin.NewAmountUint64(40000000000),
			TypeDetail: fmt.Sprintf("/%v", now.Add(15*repo.cfg.GlobalAggregationInterval.Parent).Format(childFormat)),
			Negative:   true,
		},
	}
	expected := []*BalanceHistoryEntry{{
		Time: now.In(location).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}},
	}, {
		Time: now.In(location).Add(15 * repo.cfg.GlobalAggregationInterval.Parent).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(40000000000),
			Negative: true,
			Bonus:    -350,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}},
	}}
	actual := repo.processBalanceHistory(balances, true, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(1000 * repo.cfg.GlobalAggregationInterval.Child)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(15 * repo.cfg.GlobalAggregationInterval.Parent).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(40000000000),
			Negative: true,
			Bonus:    -350,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(40000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(40000000000),
				Negative: true,
				Bonus:    -350,
			},
		}},
	}, {
		Time: now.In(location).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)

	adoptions = []*Adoption[coin.ICEFlake]{{
		AchievedAt:     time.New(now.Add(-24 * stdlibtime.Hour)),
		BaseMiningRate: coin.UnsafeParseAmount("16000000000"),
		Milestone:      1,
	}, {
		AchievedAt:     time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
		BaseMiningRate: coin.UnsafeParseAmount("8000000000"),
		Milestone:      2,
	}, {
		AchievedAt:     time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
		BaseMiningRate: coin.UnsafeParseAmount("4000000000"),
		Milestone:      3,
	}}
	preStakingSummaries = []*PreStakingSummary{{
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
			Years:      3,
			Allocation: 50,
		},
		Bonus: 100,
	}, {
		PreStaking: &PreStaking{
			CreatedAt:  time.New(now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).Add(8 * stdlibtime.Second)),
			Years:      5,
			Allocation: 100,
		},
		Bonus: 200,
	}}
	expected = []*BalanceHistoryEntry{{
		Time: now.In(location).Add(15 * repo.cfg.GlobalAggregationInterval.Parent).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(120000000000),
			Negative: true,
			Bonus:    -3100,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(15 * repo.cfg.GlobalAggregationInterval.Parent).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(120000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(120000000000),
				Negative: true,
				Bonus:    -3100,
			},
		}},
	}, {
		Time: now.In(location).Add(-45 * stdlibtime.Minute),
		Balance: &BalanceHistoryBalanceDiff{
			Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
			amount:   coin.NewAmountUint64(24000000000),
			Negative: true,
			Bonus:    -250,
		},
		TimeSeries: []*BalanceHistoryEntry{{
			Time: now.Add(5 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(32000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(32000000000),
				Negative: true,
				Bonus:    -300,
			},
		}, {
			Time: now.Add(4 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(8000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(8000000000),
				Negative: false,
				Bonus:    -50,
			},
		}, {
			Time: now.Add(3 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}, {
			Time: now.Add(2 * repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.Add(repo.cfg.GlobalAggregationInterval.Child).In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: true,
				Bonus:    -250,
			},
		}, {
			Time: now.In(location),
			Balance: &BalanceHistoryBalanceDiff{
				Amount:   coin.NewAmountUint64(24000000000).UnsafeICE(),
				amount:   coin.NewAmountUint64(24000000000),
				Negative: false,
				Bonus:    50,
			},
		}},
	}}
	actual = repo.processBalanceHistory(balances, false, utcOffset, adoptions, preStakingSummaries)
	assert.EqualValues(t, expected, actual)
}
