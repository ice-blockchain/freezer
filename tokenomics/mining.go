// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetRankingSummary(ctx context.Context, userID string) (*RankingSummary, error) { //nolint:funlen // A lot of SQL.
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	if true {
		return &RankingSummary{}, nil
	}
	sql := fmt.Sprintf(`
SELECT count(others.user_id) + 1 AS global_rank
FROM (SELECT x.amount_w0,
			 x.amount_w1,
			 x.amount_w2,
			 x.amount_w3
	  FROM (SELECT amount_w0,
				   amount_w1,
				   amount_w2,
				   amount_w3
			FROM balances
			WHERE user_id = :user_id
			UNION ALL
			SELECT %[1]v AS amount_w0,
				   0 AS amount_w1,
				   0 AS amount_w2,
				   0 AS amount_w3
		   ) AS x
	  LIMIT 1) AS this
	JOIN balances AS others
		ON	(
			 CASE
			   WHEN others.amount_w3 == this.amount_w3
				   THEN (
						 CASE
							WHEN others.amount_w2 == this.amount_w2
								THEN (
									  CASE
										 WHEN others.amount_w1 == this.amount_w1
											 THEN (others.amount_w0 >= this.amount_w0)
										 ELSE others.amount_w1 > this.amount_w1
									  END
									 )
							ELSE others.amount_w2 > this.amount_w2
						 END
						)
			   ELSE others.amount_w3 > this.amount_w3
			 END
			)
		AND others.user_id != :user_id
UNION ALL
SELECT (CASE WHEN hide_ranking == TRUE THEN 1 ELSE 2 END)
FROM users 
WHERE user_id = :user_id`, registrationICEFlakeBonusAmount)
	params := make(map[string]any, 1)
	params["user_id"] = userID
	resp := make([]*RankingSummary, 0, 1+1)
	if err := r.db.PrepareExecuteTyped(sql, params, &resp); err != nil {
		return nil, errors.Wrapf(err, "failed to select miner global rank for userID:%v", userID)
	}
	if len(resp) == 1 {
		return nil, storage.ErrRelationNotFound
	}
	if resp[1].GlobalRank == 1 && userID != requestingUserID(ctx) {
		return nil, ErrGlobalRankHidden
	}

	return resp[0], nil
}

func (r *repository) GetTopMiners(ctx context.Context, keyword string, limit, offset uint64) ([]*Miner, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	if true {
		return make([]*Miner, 0), nil
	}
	if keyword == "" {
		return r.getTopMiners(ctx, limit, offset)
	} else { //nolint:revive // Nope.
		return r.getTopMinersByKeyword(ctx, keyword, limit, offset)
	}
}

func (r *repository) getTopMinersByKeyword(ctx context.Context, keyword string, limit, offset uint64) ([]*Miner, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`SELECT b.amount,
							   u.user_id,
							   u.username,
							   %[2]v AS profile_picture_url
						FROM users u
							JOIN balances b
								ON u.user_id = b.user_id
						WHERE (
								( u.username IS NOT NULL AND u.username LIKE :keyword ESCAPE '\' )
								OR
								( u.first_name IS NOT NULL AND u.first_name != '' AND LOWER(u.first_name) LIKE :keyword ESCAPE '\' )
								OR
								( u.last_name IS NOT NULL AND u.last_name != '' AND LOWER(u.last_name) LIKE :keyword ESCAPE '\' )
							  )
							  AND u.hide_ranking = FALSE
						ORDER BY b.amount_w3 DESC,
								 b.amount_w2 DESC,
								 b.amount_w1 DESC,
								 b.amount_w0 DESC
						LIMIT %[1]v OFFSET :offset`, limit, r.pictureClient.SQLAliasDownloadURL("u.profile_picture_name"))
	params := make(map[string]any, 1+1)
	params["offset"] = offset
	params["keyword"] = fmt.Sprintf("%v%%", strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(keyword), "_", "\\_"), "%", "\\%"))
	resp := make([]*Miner, 0, limit)
	err := errors.Wrapf(r.db.PrepareExecuteTyped(sql, params, &resp), "failed to select for top miners for params:%#v", params)

	return resp, err
}

func (r *repository) getTopMiners(ctx context.Context, limit, offset uint64) ([]*Miner, error) { //nolint:revive // .
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`SELECT b.amount,
							   u.user_id,
							   u.username,
							   %[2]v AS profile_picture_url
							FROM balances b
								JOIN users u
								    ON u.user_id = b.user_id
									AND u.hide_ranking = FALSE
							ORDER BY b.amount_w3 DESC,
									 b.amount_w2 DESC,
									 b.amount_w1 DESC,
									 b.amount_w0 DESC
              				LIMIT %[1]v OFFSET :offset`, limit, r.pictureClient.SQLAliasDownloadURL("u.profile_picture_name"))
	params := make(map[string]any, 1)
	params["offset"] = offset
	resp := make([]*Miner, 0, limit)
	if err := r.db.PrepareExecuteTyped(sql, params, &resp); err != nil {
		return nil, errors.Wrapf(err, "failed to select for top miners for limit:%v,offset:%v", limit, offset)
	}

	return resp, nil
}

//nolint:funlen,lll // .
func (r *repository) GetMiningSummary(ctx context.Context, userID string) (*MiningSummary, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	now := time.Now()
	sql := fmt.Sprintf(`
SELECT u.last_natural_mining_started_at,
	   u.last_mining_started_at,
	   u.last_mining_ended_at,
	   (select current_adoption.base_mining_rate from (%[2]v) current_adoption),
	   COALESCE(btotal.amount,'0') AS total_amount,
	   COALESCE(bt0.amount,'0') AS t0_amount,
	   COALESCE(bt1.amount,'0') AS t1_amount,
	   COALESCE(bt2.amount,'0') AS t2_amount,
	   COALESCE(aggressive_degradation_btotal.amount,'0') AS aggressive_degradation_reference_total_amount,
	   COALESCE(aggressive_degradation_bt0.amount,'0') AS aggressive_degradation_reference_t0_amount,
	   COALESCE(aggressive_degradation_bt1.amount,'0') AS aggressive_degradation_reference_t1_amount,
	   COALESCE(aggressive_degradation_bt2.amount,'0') AS aggressive_degradation_reference_t2_amount,
	   COALESCE(degradation_btotalt0t1t2.amount,'0') AS degradation_reference_total_t0_t1_t2_amount,
	   u.user_id,
	   (CASE WHEN t0.user_id IS NULL THEN 0 ELSE 1 END) AS t0,
	   COALESCE(ar_worker.t1,0) as t1,
	   COALESCE(ar_worker.t2,0) as t2,
	   (CASE WHEN EXTRACT(epoch from COALESCE(eb_worker.extra_bonus_ended_at, to_timestamp(0)))*%[11]v > $2::bigint THEN eb_worker.extra_bonus ELSE 0 END) AS extra_bonus,
	   COALESCE(x.pre_staking_allocation,0) as pre_staking_allocation,
	   COALESCE(st_b.bonus, 0) as pre_staking_bonus,
	   COALESCE(eb.bonus,0) AS flat_bonus,
	   (CASE WHEN (eb_worker.user_id IS NOT NULL AND ($2::bigint - EXTRACT(epoch from COALESCE(eb_worker.extra_bonus_started_at, to_timestamp(0)))*%[11]v > $9::bigint) AND ebw.extra_bonus_index IS NOT NULL)
			 	THEN (100 - (25 *  ((CASE WHEN ($2::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (ebw.extra_bonus_index * $3::bigint)) - $8::bigint - ((ebw.offset_value * $5::bigint) / $10)) < $7::bigint THEN 0 ELSE ($2::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (ebw.extra_bonus_index * $3::bigint)) - $8::bigint - ((ebw.offset_value * $5::bigint) / $10)) END)/$6::bigint)))
	   		 ELSE 0
	    END) AS bonus_percentage_remaining,
	   eb_worker.news_seen
FROM (SELECT MAX(st.years) AS pre_staking_years,
		     MAX(st.allocation) AS pre_staking_allocation,
			 x.user_id
			 FROM ( SELECT CAST($1 AS VARCHAR) AS user_id ) x
				 LEFT JOIN pre_stakings_%[1]v st
						ON st.user_id = x.user_id
			 GROUP BY x.user_id 
	 ) x
	    JOIN users u
		  ON u.user_id = x.user_id
   		JOIN extra_bonus_start_date sd 
		  ON sd.key = 0
   LEFT	JOIN extra_bonus_processing_worker_%[1]v eb_worker
		  ON eb_worker.user_id = x.user_id
   LEFT JOIN active_referrals_%[1]v ar_worker
		  ON ar_worker.user_id = x.user_id
   LEFT JOIN extra_bonuses eb 
          ON eb.ix = ($2::bigint + (eb_worker.utc_offset * $4::bigint) - sd.value) / $3::bigint
		 AND $2::bigint + (eb_worker.utc_offset * $4::bigint) > sd.value
		 AND eb.bonus > 0
   LEFT JOIN extra_bonuses_worker_%[1]v ebw
		  ON eb.ix = ebw.extra_bonus_index
		 AND $2::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (ebw.extra_bonus_index * $3::bigint)) - $8::bigint - ((ebw.offset_value * $5::bigint) / $10::bigint) < $9::bigint
		 AND $2::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (ebw.extra_bonus_index * $3::bigint)) - $8::bigint - ((ebw.offset_value * $5::bigint) / $10::bigint) > 0
   LEFT JOIN pre_staking_bonuses st_b
		  ON st_b.years = x.pre_staking_years
   LEFT JOIN balances_worker_%[1]v btotal
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND btotal.user_id = u.user_id
	     AND btotal.negative = FALSE
	     AND btotal.type = %[3]v
	     AND btotal.type_detail = ''
   LEFT JOIN balances_worker_%[1]v bt0
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND bt0.user_id = u.user_id
	     AND bt0.negative = FALSE
	     AND bt0.type = %[3]v
	     AND bt0.type_detail = '%[4]v_' || u.referred_by
   LEFT JOIN balances_worker_%[1]v bt1
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND bt1.user_id = u.user_id
	     AND bt1.negative = FALSE
	     AND bt1.type = %[3]v
	     AND bt1.type_detail = '%[5]v'
   LEFT JOIN balances_worker_%[1]v bt2
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND bt2.user_id = u.user_id
	     AND bt2.negative = FALSE
	     AND bt2.type = %[3]v
	     AND bt2.type_detail = '%[6]v'
   LEFT JOIN balances_worker_%[1]v aggressive_degradation_btotal
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND aggressive_degradation_btotal.user_id = u.user_id
	     AND aggressive_degradation_btotal.negative = FALSE
	     AND aggressive_degradation_btotal.type = %[3]v
	     AND aggressive_degradation_btotal.type_detail = '%[7]v'
   LEFT JOIN balances_worker_%[1]v aggressive_degradation_bt0
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND aggressive_degradation_bt0.user_id = u.user_id
	     AND aggressive_degradation_bt0.negative = FALSE
	     AND aggressive_degradation_bt0.type = %[3]v
	     AND aggressive_degradation_bt0.type_detail = '%[4]v_' || u.referred_by || '_'
   LEFT JOIN balances_worker_%[1]v aggressive_degradation_bt1
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND aggressive_degradation_bt1.user_id = u.user_id
	     AND aggressive_degradation_bt1.negative = FALSE
	     AND aggressive_degradation_bt1.type = %[3]v
	     AND aggressive_degradation_bt1.type_detail = '%[8]v'
   LEFT JOIN balances_worker_%[1]v aggressive_degradation_bt2
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND aggressive_degradation_bt2.user_id = u.user_id
	     AND aggressive_degradation_bt2.negative = FALSE
	     AND aggressive_degradation_bt2.type = %[3]v
	     AND aggressive_degradation_bt2.type_detail = '%[9]v'
   LEFT JOIN balances_worker_%[1]v degradation_btotalt0t1t2
		  ON (u.last_mining_ended_at IS NOT NULL AND EXTRACT(epoch from u.last_mining_ended_at)*%[11]v < $2::bigint )
	     AND degradation_btotalt0t1t2.user_id = u.user_id
	     AND degradation_btotalt0t1t2.negative = FALSE
	     AND degradation_btotalt0t1t2.type = %[3]v
	     AND degradation_btotalt0t1t2.type_detail = '%[10]v'
   LEFT JOIN users t0
	  	  ON t0.user_id = u.referred_by
	     AND t0.user_id != x.user_id
	  	 AND t0.last_mining_ended_at IS NOT NULL
	  	 AND EXTRACT(epoch from t0.last_mining_ended_at)*%[11]v  > $2::bigint`,
		r.workerIndex(ctx),
		currentAdoptionSQL(),
		totalNoPreStakingBonusBalanceType,
		t0BalanceTypeDetail,
		t1BalanceTypeDetail,
		t2BalanceTypeDetail,
		aggressiveDegradationTotalReferenceBalanceTypeDetail,
		aggressiveDegradationT1ReferenceBalanceTypeDetail,
		aggressiveDegradationT2ReferenceBalanceTypeDetail,
		degradationT0T1T2TotalReferenceBalanceTypeDetail,
		pgMicrosecordsPrecision,
	)
	const networkLagDelta = 1.3
	params := []any{
		userID,
		now.UnixMicro(),
		r.cfg.ExtraBonuses.Duration,
		r.cfg.ExtraBonuses.UTCOffsetDuration,
		r.cfg.ExtraBonuses.AvailabilityWindow,
		r.cfg.ExtraBonuses.DelayedClaimPenaltyWindow,
		stdlibtime.Duration(float64(r.cfg.ExtraBonuses.DelayedClaimPenaltyWindow.Nanoseconds()) * networkLagDelta),
		r.cfg.ExtraBonuses.TimeToAvailabilityWindow,
		r.cfg.ExtraBonuses.ClaimWindow,
		r.cfg.WorkerCount,
	}
	type respStruct struct {
		LastNaturalMiningStartedAt                *time.Time
		LastMiningStartedAt                       *time.Time
		LastMiningEndedAt                         *time.Time
		BaseMiningRate                            *coin.ICEFlake
		TotalAmount                               *coin.ICEFlake
		T0Amount                                  *coin.ICEFlake
		T1Amount                                  *coin.ICEFlake
		T2Amount                                  *coin.ICEFlake
		AggressiveDegradationReferenceTotalAmount *coin.ICEFlake
		AggressiveDegradationReferenceT0Amount    *coin.ICEFlake
		AggressiveDegradationReferenceT1Amount    *coin.ICEFlake
		AggressiveDegradationReferenceT2Amount    *coin.ICEFlake
		DegradationReferenceTotalT0T1T2Amount     *coin.ICEFlake
		userMiningRateRecalculationParameters
		FlatBonus                uint64
		BonusPercentageRemaining uint64
		NewsSeen                 uint64
	}
	resp, err := storage.Select[respStruct](ctx, r.dbV2, sql, params...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select for mining summary for userID:%v", userID)
	}
	if len(resp) == 0 {
		return nil, storage.ErrRelationNotFound
	}
	var mrt MiningRateType
	var negativeMiningRate *coin.ICEFlake
	if resp[0].LastMiningEndedAt == nil { //nolint:gocritic,nestif // Wrong.
		mrt = NoneMiningRateType
	} else if resp[0].LastMiningEndedAt.After(*now.Time) {
		mrt = PositiveMiningRateType
	} else if resp[0].TotalAmount.Add(resp[0].T0Amount).Add(resp[0].T1Amount).Add(resp[0].T2Amount).IsZero() {
		mrt = NoneMiningRateType
	} else {
		mrt = NegativeMiningRateType
		if aggressive := resp[0].LastMiningEndedAt.Add(r.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time); aggressive {
			referenceAmount := resp[0].AggressiveDegradationReferenceTotalAmount.
				Add(resp[0].AggressiveDegradationReferenceT0Amount).
				Add(resp[0].AggressiveDegradationReferenceT1Amount).
				Add(resp[0].AggressiveDegradationReferenceT2Amount)
			negativeMiningRate = r.calculateDegradation(r.cfg.GlobalAggregationInterval.Child, referenceAmount, true)
		} else {
			negativeMiningRate = r.calculateDegradation(r.cfg.GlobalAggregationInterval.Child, resp[0].DegradationReferenceTotalT0T1T2Amount, false)
		}
	}
	miningStreak := r.calculateMiningStreak(now, resp[0].LastMiningStartedAt, resp[0].LastMiningEndedAt)

	return &MiningSummary{
		MiningStreak:                miningStreak,
		MiningSession:               r.calculateMiningSession(now, resp[0].LastNaturalMiningStartedAt, resp[0].LastMiningEndedAt),
		RemainingFreeMiningSessions: r.calculateRemainingFreeMiningSessions(now, resp[0].LastMiningEndedAt),
		MiningRates:                 r.calculateMiningRateSummaries(resp[0].BaseMiningRate, &resp[0].userMiningRateRecalculationParameters, negativeMiningRate, mrt), //nolint:lll // .
		ExtraBonusSummary: ExtraBonusSummary{
			AvailableExtraBonus: r.calculateExtraBonus(resp[0].FlatBonus, resp[0].BonusPercentageRemaining, resp[0].NewsSeen, miningStreak),
		},
	}, nil
}

func (r *repository) calculateMiningSession(now, start, end *time.Time) (ms *MiningSession) {
	if end == nil || end.Before(*now.Time) {
		return nil
	}
	lastMiningStartedAt := time.New(start.Add((now.Sub(*start.Time) / r.cfg.MiningSessionDuration.Max) * r.cfg.MiningSessionDuration.Max))
	free := start.Add(r.cfg.MiningSessionDuration.Max).Before(*now.Time)

	return &MiningSession{
		StartedAt:                     lastMiningStartedAt,
		EndedAt:                       time.New(lastMiningStartedAt.Add(r.cfg.MiningSessionDuration.Max)),
		Free:                          &free,
		ResettableStartingAt:          time.New(lastMiningStartedAt.Add(r.cfg.MiningSessionDuration.Min)),
		WarnAboutExpirationStartingAt: time.New(lastMiningStartedAt.Add(r.cfg.MiningSessionDuration.WarnAboutExpirationAfter)),
	}
}

func (r *repository) calculateMiningRateSummaries( //nolint:funlen // A lot of calculations.
	baseMiningRate *coin.ICEFlake, params *userMiningRateRecalculationParameters, negativeMiningRate *coin.ICEFlake, miningRateType MiningRateType,
) (miningRates *MiningRates[MiningRateSummary[coin.ICE]]) {
	miningRates = new(MiningRates[MiningRateSummary[coin.ICE]])
	miningRates.Type = miningRateType
	var (
		standardMiningRate, preStakingMiningRate *coin.ICEFlake
		preStakingBonusVal                       uint64
	)
	miningRates.Base = &MiningRateSummary[coin.ICE]{
		Amount: baseMiningRate.UnsafeICE(),
	}
	if params.PreStakingAllocation != percentage100 {
		var totalBonus uint64
		switch miningRates.Type {
		case PositiveMiningRateType:
			standardMiningRate = r.calculateMintedStandardCoins(baseMiningRate, params, r.cfg.GlobalAggregationInterval.Child, false)
			totalBonus = coin.ZeroICEFlakes().
				Add(standardMiningRate).
				Subtract(baseMiningRate).
				MultiplyUint64(percentage100).
				Divide(baseMiningRate).Uint64()
		case NegativeMiningRateType:
			standardMiningRate = coin.ZeroICEFlakes().
				Add(negativeMiningRate.
					MultiplyUint64(percentage100 - params.PreStakingAllocation).
					DivideUint64(percentage100))
		case NoneMiningRateType:
			standardMiningRate = coin.ZeroICEFlakes()
		}
		miningRates.Standard = &MiningRateSummary[coin.ICE]{
			Amount: standardMiningRate.UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:    uint64(float64((params.T0*r.cfg.ReferralBonusMiningRates.T0+params.T1*r.cfg.ReferralBonusMiningRates.T1)*(percentage100-params.PreStakingAllocation)) / float64(percentage100)), //nolint:lll // .
				T2:    uint64(float64(params.T2*r.cfg.ReferralBonusMiningRates.T2*(percentage100-params.PreStakingAllocation)) / float64(percentage100)),
				Extra: uint64(float64(params.ExtraBonus*(percentage100-params.PreStakingAllocation)) / float64(percentage100)),
				Total: totalBonus,
			},
		}
	}
	if params.PreStakingAllocation != 0 {
		var totalBonus uint64
		switch miningRates.Type {
		case PositiveMiningRateType:
			preStakingMiningRate = r.calculateMintedPreStakingCoins(baseMiningRate, params, r.cfg.GlobalAggregationInterval.Child, false)
			totalBonus = coin.ZeroICEFlakes().
				Add(preStakingMiningRate).
				Subtract(baseMiningRate).
				MultiplyUint64(percentage100).
				Divide(baseMiningRate).Uint64()
		case NegativeMiningRateType:
			preStakingMiningRate = coin.ZeroICEFlakes().
				Add(negativeMiningRate.
					MultiplyUint64((params.PreStakingBonus + percentage100) * params.PreStakingAllocation).
					DivideUint64(percentage100 * percentage100))
		case NoneMiningRateType:
			preStakingMiningRate = coin.ZeroICEFlakes()
		}
		t1Bonus := float64((params.T0*r.cfg.ReferralBonusMiningRates.T0+params.T1*r.cfg.ReferralBonusMiningRates.T1)*params.PreStakingAllocation) / float64(percentage100) //nolint:lll // .
		t2Bonus := float64(params.T2*r.cfg.ReferralBonusMiningRates.T2*params.PreStakingAllocation) / float64(percentage100)
		extraBonus := float64(params.ExtraBonus*params.PreStakingAllocation) / float64(percentage100)
		preStakingBonusVal = uint64((float64(params.PreStakingAllocation) * float64(params.PreStakingBonus)) / float64(percentage100))
		miningRates.PreStaking = &MiningRateSummary[coin.ICE]{
			Amount: preStakingMiningRate.UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         uint64(t1Bonus),
				T2:         uint64(t2Bonus),
				Extra:      uint64(extraBonus),
				PreStaking: preStakingBonusVal,
				Total:      totalBonus,
			},
		}
	}
	totalNoStakingBonusParams := *params
	totalNoStakingBonusParams.PreStakingAllocation, totalNoStakingBonusParams.PreStakingBonus = 0, 0
	positiveTotalNoPreStakingBonus := r.calculateMintedStandardCoins(baseMiningRate, &totalNoStakingBonusParams, r.cfg.GlobalAggregationInterval.Child, false)
	positiveTotalNoPreStakingBonusVal := positiveTotalNoPreStakingBonus.
		Subtract(baseMiningRate).
		MultiplyUint64(percentage100).
		Divide(baseMiningRate).Uint64()
	var totalBonus, totalNoPreStakingBonusVal uint64
	var totalNoPreStakingBonus *coin.ICEFlake
	switch miningRates.Type {
	case PositiveMiningRateType:
		totalBonus = coin.ZeroICEFlakes().
			Add(standardMiningRate).
			Add(preStakingMiningRate).
			Subtract(baseMiningRate).
			MultiplyUint64(percentage100).
			Divide(baseMiningRate).Uint64()
		totalNoPreStakingBonus = positiveTotalNoPreStakingBonus
		totalNoPreStakingBonusVal = positiveTotalNoPreStakingBonusVal
	case NegativeMiningRateType:
		totalNoPreStakingBonus = coin.ZeroICEFlakes().Add(negativeMiningRate)
	case NoneMiningRateType:
		totalNoPreStakingBonus = coin.ZeroICEFlakes()
	}
	miningRates.Total = &MiningRateSummary[coin.ICE]{
		Amount: standardMiningRate.Add(preStakingMiningRate).UnsafeICE(),
		Bonuses: &MiningRateBonuses{
			T1:         params.T0*r.cfg.ReferralBonusMiningRates.T0 + params.T1*r.cfg.ReferralBonusMiningRates.T1,
			T2:         params.T2 * r.cfg.ReferralBonusMiningRates.T2,
			Extra:      params.ExtraBonus,
			PreStaking: preStakingBonusVal,
			Total:      totalBonus,
		},
	}
	miningRates.TotalNoPreStakingBonus = &MiningRateSummary[coin.ICE]{
		Amount: totalNoPreStakingBonus.UnsafeICE(),
		Bonuses: &MiningRateBonuses{
			T1:         params.T0*r.cfg.ReferralBonusMiningRates.T0 + params.T1*r.cfg.ReferralBonusMiningRates.T1,
			T2:         params.T2 * r.cfg.ReferralBonusMiningRates.T2,
			Extra:      params.ExtraBonus,
			PreStaking: 0,
			Total:      totalNoPreStakingBonusVal,
		},
	}
	miningRates.PositiveTotalNoPreStakingBonus = &MiningRateSummary[coin.ICE]{
		Amount: positiveTotalNoPreStakingBonus.UnsafeICE(),
		Bonuses: &MiningRateBonuses{
			T1:         params.T0*r.cfg.ReferralBonusMiningRates.T0 + params.T1*r.cfg.ReferralBonusMiningRates.T1,
			T2:         params.T2 * r.cfg.ReferralBonusMiningRates.T2,
			Extra:      params.ExtraBonus,
			PreStaking: 0,
			Total:      positiveTotalNoPreStakingBonusVal,
		},
	}

	return miningRates
}

func (r *repository) calculateMintedStandardCoins( //nolint:revive // Not an issue here.
	baseMiningRate *coin.ICEFlake, params *userMiningRateRecalculationParameters, elapsedNanos stdlibtime.Duration, excludeBaseRate bool,
) *coin.ICEFlake {
	if params.PreStakingAllocation == percentage100 || elapsedNanos <= 0 {
		return nil
	}
	var includeBaseMiningRate uint64
	if !excludeBaseRate {
		includeBaseMiningRate = percentage100 + params.ExtraBonus
	}
	mintedBase := includeBaseMiningRate +
		params.T0*r.cfg.ReferralBonusMiningRates.T0 +
		params.T1*r.cfg.ReferralBonusMiningRates.T1 +
		params.T2*r.cfg.ReferralBonusMiningRates.T2
	if mintedBase == 0 {
		return nil
	}
	if elapsedNanos == r.cfg.GlobalAggregationInterval.Child {
		return baseMiningRate.
			MultiplyUint64(mintedBase).
			MultiplyUint64(percentage100 - params.PreStakingAllocation).
			DivideUint64(percentage100 * percentage100)
	}

	return baseMiningRate.
		MultiplyUint64(uint64(elapsedNanos.Nanoseconds())).
		MultiplyUint64(mintedBase).
		MultiplyUint64(percentage100 - params.PreStakingAllocation).
		DivideUint64(uint64(r.cfg.GlobalAggregationInterval.Child.Nanoseconds()) * percentage100 * percentage100)
}

func (r *repository) calculateMintedPreStakingCoins( //nolint:revive // Not an issue here.
	baseMiningRate *coin.ICEFlake, params *userMiningRateRecalculationParameters, elapsedNanos stdlibtime.Duration, excludeBaseRate bool,
) *coin.ICEFlake {
	if params.PreStakingAllocation == 0 || elapsedNanos <= 0 {
		return nil
	}
	var includeBaseMiningRate uint64
	if !excludeBaseRate {
		includeBaseMiningRate = percentage100 + params.ExtraBonus
	}
	mintedBase := includeBaseMiningRate +
		params.T0*r.cfg.ReferralBonusMiningRates.T0 +
		params.T1*r.cfg.ReferralBonusMiningRates.T1 +
		params.T2*r.cfg.ReferralBonusMiningRates.T2
	if mintedBase == 0 {
		return nil
	}
	if elapsedNanos == r.cfg.GlobalAggregationInterval.Child {
		return baseMiningRate.
			MultiplyUint64(mintedBase).
			MultiplyUint64((params.PreStakingBonus + percentage100) * params.PreStakingAllocation).
			DivideUint64(percentage100 * percentage100 * percentage100)
	}

	return baseMiningRate.
		MultiplyUint64(uint64(elapsedNanos.Nanoseconds())).
		MultiplyUint64(mintedBase).
		MultiplyUint64((params.PreStakingBonus + percentage100) * params.PreStakingAllocation).
		DivideUint64(uint64(r.cfg.GlobalAggregationInterval.Child.Nanoseconds()) * percentage100 * percentage100 * percentage100)
}

func (r *repository) calculateMiningStreak(now, start, end *time.Time) uint64 {
	if start == nil || end == nil || now.After(*end.Time) || now.Before(*start.Time) {
		return 0
	}

	return uint64(now.Sub(*start.Time) / r.cfg.MiningSessionDuration.Max)
}

func (r *repository) calculateRemainingFreeMiningSessions(now, end *time.Time) uint64 {
	if end == nil || now.After(*end.Time) {
		return 0
	}

	return uint64(end.Sub(*now.Time) / r.cfg.MiningSessionDuration.Max)
}
