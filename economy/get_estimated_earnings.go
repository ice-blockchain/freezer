// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"

	"cosmossdk.io/math"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
)

func (e *economy) GetEstimatedEarnings(ctx context.Context, arg *GetEstimatedEarningsArg) (*EstimatedEarnings, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get estimated earnings failed because context failed")
	}
	ed, err := e.getEstimatedEarningsCalculationData(ctx, arg.StakingYears)
	if err != nil {
		return nil, errors.Wrap(err, "can't call getEstimatedEarningsCalculationData function")
	}
	ee := e.calculateEstimatedEarnings(&calculateEstimatedEarningsArg{
		GetEstimatedEarningsArg:          arg,
		estimatedEarningsCalculationData: ed,
	})

	return ee, nil
}

func (e *economy) getEstimatedEarningsCalculationData(ctx context.Context, stakingYears uint8) (*estimatedEarningsCalculationData, error) {
	sql := `SELECT 
				base_hourly_mining_rate,
				(SELECT percentage FROM staking_bonus WHERE years = :stakingYears) AS bonus
			FROM adoption a
			WHERE a.active = true`

	params := map[string]interface{}{
		"stakingYears": stakingYears,
	}

	var res []*estimatedEarningsCalculationData
	if err := e.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return nil, errors.Wrap(err, "failed to get base hourly mining rate and staking bonus data")
	}
	if len(res) == 0 {
		return nil, errors.New("no base hourly mining rate and staking bonus data")
	}

	return res[0], nil
}

func (e *economy) calculateEstimatedEarnings(arg *calculateEstimatedEarningsArg) *EstimatedEarnings {
	t0Referrals := uint64(0)
	if arg.T0ActiveReferee {
		t0Referrals = 1
	}

	rateMultiplier := t0Referrals*cfg.Rates.Tier0 + arg.T1ActiveReferrals*cfg.Rates.Tier1 + arg.T2ActiveReferrals*cfg.Rates.Tier2 + percentage100
	hmr := arg.BaseHourlyMiningRate.MulUint64(rateMultiplier).QuoUint64(percentage100)

	normalHMR := math.NewUint(percentage100 - uint64(arg.StakingAllocation)).Mul(hmr).QuoUint64(percentage100)
	stakedHMR := math.NewUint(arg.StakingPercentageBonus).Mul(hmr).Mul(math.NewUint(uint64(arg.StakingAllocation))).QuoUint64(stakedHourlyMiningRateDivider)

	return &EstimatedEarnings{
		StandardHourlyMiningRate: &coin.ICEFlake{Uint: normalHMR},
		StakingHourlyMiningRate:  &coin.ICEFlake{Uint: stakedHMR},
	}
}
