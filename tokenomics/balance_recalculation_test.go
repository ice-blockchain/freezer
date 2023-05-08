// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"testing"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CalculateMiningRateSummaries(t *testing.T) { //nolint:funlen,maintidx // .
	t.Parallel()
	rep := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
		ReferralBonusMiningRates: struct {
			T0 uint64 `yaml:"t0"`
			T1 uint64 `yaml:"t1"`
			T2 uint64 `yaml:"t2"`
		}{
			T0: 25,
			T1: 25,
			T2: 5,
		},
	}}
	baseMiningRate := coin.UnsafeParseAmount("16000000000")
	actual := rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      100,
		PreStakingAllocation: 50,
	}, nil, PositiveMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Intended.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("42000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         137,
				T2:         25,
				PreStaking: 0,
				Extra:      50,
				Total:      162,
			},
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         137,
				T2:         25,
				PreStaking: 50,
				Extra:      50,
				Total:      425,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("126000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      687,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      500,
		PreStakingAllocation: 10,
	}, nil, PositiveMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Intended.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("75600000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         247,
				T2:         45,
				PreStaking: 0,
				Extra:      90,
				Total:      372,
			},
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("50400000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         27,
				T2:         5,
				PreStaking: 50,
				Extra:      10,
				Total:      215,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("126000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      687,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      100,
		PreStakingAllocation: 100,
	}, nil, PositiveMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Wrong.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("168000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 100,
				Extra:      100,
				Total:      950,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("168000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 100,
				Extra:      100,
				Total:      950,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:         1,
		T1:         10,
		T2:         10,
		ExtraBonus: 100,
	}, nil, PositiveMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Wrong.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      500,
		PreStakingAllocation: 10,
	}, coin.UnsafeParseAmount("1111"), NoneMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Wrong.
		Type: NoneMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("0").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         247,
				T2:         45,
				PreStaking: 0,
				Extra:      90,
				Total:      0,
			},
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("0").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         27,
				T2:         5,
				PreStaking: 50,
				Extra:      10,
				Total:      0,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("0").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      0,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("0").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      0,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      500,
		PreStakingAllocation: 10,
	}, coin.UnsafeParseAmount("1000"), NegativeMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Wrong.
		Type: NegativeMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("900").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         247,
				T2:         45,
				PreStaking: 0,
				Extra:      90,
				Total:      0,
			},
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("600").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         27,
				T2:         5,
				PreStaking: 50,
				Extra:      10,
				Total:      0,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("1500").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      0,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("1000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      0,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
}
