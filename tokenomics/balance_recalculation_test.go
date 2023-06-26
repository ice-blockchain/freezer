// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"fmt"
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"

	"github.com/ice-blockchain/wintr/time"
)

func TestRepository_CalculateMiningRateSummaries(t *testing.T) { //nolint:funlen,maintidx // .
	t.Parallel()
	rep := &repository{cfg: &Config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
		ReferralBonusMiningRates: struct {
			T0 uint16 `yaml:"t0"`
			T1 uint32 `yaml:"t1"`
			T2 uint32 `yaml:"t2"`
		}{
			T0: 25,
			T1: 25,
			T2: 5,
		},
	}}
	var (
		baseMiningRate       = 16.0
		negativeMiningRate   = 1000.0
		totalBalance         = 0.0
		t0                   = uint16(1)
		t1                   = int32(10)
		t2                   = int32(10)
		extraBonus           = uint16(100)
		preStakingBonus      = uint16(100)
		preStakingAllocation = uint16(50)
		now                  = time.Now()
		endedAt              = time.New(now.Add(stdlibtime.Second))
	)
	actual := rep.calculateMiningRateSummaries(extraBonus, t0, preStakingAllocation, preStakingBonus, t1, t2, baseMiningRate, negativeMiningRate, totalBalance, now, endedAt)
	assert.EqualValues(t, &MiningRates[*MiningRateSummary[string]]{ //nolint:dupl // Intended.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[string]{
			Amount: fmt.Sprintf(floatToStringFormatter, baseMiningRate),
		},
		Standard: &MiningRateSummary[string]{
			Amount: "42.00",
			Bonuses: &MiningRateBonuses{
				T1:         137,
				T2:         25,
				PreStaking: 0,
				Extra:      50,
				Total:      162,
			},
		},
		PreStaking: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         137,
				T2:         25,
				PreStaking: 50,
				Extra:      50,
				Total:      425,
			},
		},
		Total: &MiningRateSummary[string]{
			Amount: "126.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      687,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	preStakingBonus = uint16(500)
	preStakingAllocation = uint16(10)
	actual = rep.calculateMiningRateSummaries(extraBonus, t0, preStakingAllocation, preStakingBonus, t1, t2, baseMiningRate, negativeMiningRate, totalBalance, now, endedAt)
	assert.EqualValues(t, &MiningRates[*MiningRateSummary[string]]{ //nolint:dupl // Intended.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[string]{
			Amount: fmt.Sprintf(floatToStringFormatter, baseMiningRate),
		},
		Standard: &MiningRateSummary[string]{
			Amount: "75.60",
			Bonuses: &MiningRateBonuses{
				T1:         247,
				T2:         45,
				PreStaking: 0,
				Extra:      90,
				Total:      372,
			},
		},
		PreStaking: &MiningRateSummary[string]{
			Amount: "50.40",
			Bonuses: &MiningRateBonuses{
				T1:         27,
				T2:         5,
				PreStaking: 50,
				Extra:      10,
				Total:      215,
			},
		},
		Total: &MiningRateSummary[string]{
			Amount: "126.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      687,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	preStakingBonus = uint16(100)
	preStakingAllocation = uint16(100)
	actual = rep.calculateMiningRateSummaries(extraBonus, t0, preStakingAllocation, preStakingBonus, t1, t2, baseMiningRate, negativeMiningRate, totalBalance, now, endedAt)
	assert.EqualValues(t, &MiningRates[*MiningRateSummary[string]]{ //nolint:dupl // Wrong.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[string]{
			Amount: fmt.Sprintf(floatToStringFormatter, baseMiningRate),
		},
		PreStaking: &MiningRateSummary[string]{
			Amount: "168.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 100,
				Extra:      100,
				Total:      950,
			},
		},
		Total: &MiningRateSummary[string]{
			Amount: "168.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 100,
				Extra:      100,
				Total:      950,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	preStakingBonus = uint16(0)
	preStakingAllocation = uint16(0)
	actual = rep.calculateMiningRateSummaries(extraBonus, t0, preStakingAllocation, preStakingBonus, t1, t2, baseMiningRate, negativeMiningRate, totalBalance, now, endedAt)
	assert.EqualValues(t, &MiningRates[*MiningRateSummary[string]]{ //nolint:dupl // Wrong.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[string]{
			Amount: fmt.Sprintf(floatToStringFormatter, baseMiningRate),
		},
		Standard: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		Total: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)

	preStakingBonus = uint16(500)
	preStakingAllocation = uint16(10)
	endedAt = now
	actual = rep.calculateMiningRateSummaries(extraBonus, t0, preStakingAllocation, preStakingBonus, t1, t2, baseMiningRate, negativeMiningRate, totalBalance, now, endedAt)
	assert.EqualValues(t, &MiningRates[*MiningRateSummary[string]]{ //nolint:dupl // Wrong.
		Type: NoneMiningRateType,
		Base: &MiningRateSummary[string]{
			Amount: fmt.Sprintf(floatToStringFormatter, baseMiningRate),
		},
		Standard: &MiningRateSummary[string]{
			Amount: "0.00",
			Bonuses: &MiningRateBonuses{
				T1:         247,
				T2:         45,
				PreStaking: 0,
				Extra:      90,
				Total:      0,
			},
		},
		PreStaking: &MiningRateSummary[string]{
			Amount: "0.00",
			Bonuses: &MiningRateBonuses{
				T1:         27,
				T2:         5,
				PreStaking: 50,
				Extra:      10,
				Total:      0,
			},
		},
		Total: &MiningRateSummary[string]{
			Amount: "0.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      0,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "0.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      0,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	totalBalance = 1.0
	actual = rep.calculateMiningRateSummaries(extraBonus, t0, preStakingAllocation, preStakingBonus, t1, t2, baseMiningRate, negativeMiningRate, totalBalance, now, endedAt)
	assert.EqualValues(t, &MiningRates[*MiningRateSummary[string]]{ //nolint:dupl // Wrong.
		Type: NegativeMiningRateType,
		Base: &MiningRateSummary[string]{
			Amount: fmt.Sprintf(floatToStringFormatter, baseMiningRate),
		},
		Standard: &MiningRateSummary[string]{
			Amount: "900.00",
			Bonuses: &MiningRateBonuses{
				T1:         0,
				T2:         0,
				PreStaking: 0,
				Extra:      0,
				Total:      0,
			},
		},
		PreStaking: &MiningRateSummary[string]{
			Amount: "600.00",
			Bonuses: &MiningRateBonuses{
				T1:         0,
				T2:         0,
				PreStaking: 50,
				Extra:      0,
				Total:      0,
			},
		},
		Total: &MiningRateSummary[string]{
			Amount: "1500.00",
			Bonuses: &MiningRateBonuses{
				T1:         0,
				T2:         0,
				PreStaking: 50,
				Extra:      0,
				Total:      0,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "1000.00",
			Bonuses: &MiningRateBonuses{
				T1:         0,
				T2:         0,
				PreStaking: 0,
				Extra:      0,
				Total:      0,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[string]{
			Amount: "84.00",
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
