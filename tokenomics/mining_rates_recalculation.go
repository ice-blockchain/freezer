// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	stdlibtime "time"

	"github.com/cenkalti/backoff/v4"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) initializeMiningRatesRecalculationWorker(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}

	return errors.Wrapf(retry(ctx, func() error {
		if err := r.initializeWorker(ctx, "mining_rates_recalculation_worker", usr.ID, usr.HashCode); err != nil {
			if errors.Is(err, storage.ErrRelationNotFound) {
				return err
			}

			return errors.Wrapf(backoff.Permanent(err), "failed to initializeMiningRatesRecalculationWorker for userID:%v", usr.ID)
		}

		return nil
	}), "permanently failed to initializeMiningRatesRecalculationWorker for userID:%v", usr.ID)
}

func (s *miningRatesRecalculationTriggerStreamSource) start(ctx context.Context) {
	log.Info("miningRatesRecalculationTriggerStreamSource started")
	defer log.Info("miningRatesRecalculationTriggerStreamSource stopped")
	workerIndexes := make([]int16, s.cfg.WorkerCount) //nolint:makezero // Intended.
	for i := 0; i < int(s.cfg.WorkerCount); i++ {
		workerIndexes[i] = int16(i)
	}
	for ctx.Err() == nil {
		stdlibtime.Sleep(s.cfg.Workers.RefreshMiningRatesProcessingSeedingStreamEmitFrequency)
		before := time.Now()
		log.Error(errors.Wrap(executeBatchConcurrently(ctx, s.process, workerIndexes), "failed to executeBatchConcurrently[miningRatesRecalculationTriggerStreamSource.process]")) //nolint:lll // .
		log.Error(errors.Errorf("miningRatesRecalculationTriggerStreamSource.process took: %v", stdlibtime.Since(*before.Time)))
	}
}

func (s *miningRatesRecalculationTriggerStreamSource) process(ignoredCtx context.Context, workerIndex int16) (err error) {
	if ignoredCtx.Err() != nil {
		return errors.Wrap(ignoredCtx.Err(), "unexpected deadline while processing message")
	}
	const deadline = 5 * stdlibtime.Minute
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	now := time.Now()
	rows, err := s.getLatestMiningRates(ctx, workerIndex, now) //nolint:contextcheck // We use context with longer deadline.
	if err != nil || len(rows) == 0 {
		return errors.Wrapf(err, "failed to getLatestMiningRates for workerIndex:%v", workerIndex)
	}
	if err = executeBatchConcurrently(ctx, s.sendMiningRatesMessage, rows); err != nil { //nolint:contextcheck // We use context with longer deadline.
		return errors.Wrapf(err, "failed to sendMiningRatesMessages for:%#v", rows)
	}

	return errors.Wrapf(s.updateLastIterationFinishedAt(ctx, workerIndex, rows, now), //nolint:contextcheck // We use context with longer deadline.
		"failed to updateLastIterationFinishedAt for workerIndex:%v,rows:%#v", workerIndex, rows)
}

func (s *miningRatesRecalculationTriggerStreamSource) getLatestMiningRates( //nolint:funlen,gocognit // .
	ctx context.Context, workerIndex int16, now *time.Time,
) ([]*MiningRates[coin.ICEFlake], error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	currentAdoption, err := s.getCurrentAdoption(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to getCurrentAdoption")
	}
	batch, err := s.getUserMiningRateCalculationParametersNewBatch(ctx, workerIndex, now)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getUserMiningRateCalculationParametersBatch for workerIndex:%v", workerIndex)
	}
	if len(batch) == 0 {
		return nil, nil
	}
	res := make([]*MiningRates[coin.ICEFlake], 0, len(batch))
	for _, row := range batch {
		var negativeMiningRate *coin.ICEFlake
		if row.LastMiningEndedAt != nil && row.LastMiningEndedAt.Before(*now.Time) {
			if aggressive := row.LastMiningEndedAt.Add(s.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter).Before(*now.Time); aggressive {
				referenceAmount := row.AggressiveDegradationReferenceTotalAmount.
					Add(row.AggressiveDegradationReferenceT0Amount).
					Add(row.AggressiveDegradationReferenceT1Amount).
					Add(row.AggressiveDegradationReferenceT2Amount)
				negativeMiningRate = s.calculateDegradation(s.cfg.GlobalAggregationInterval.Child, referenceAmount, true)
			} else {
				negativeMiningRate = s.calculateDegradation(s.cfg.GlobalAggregationInterval.Child, row.DegradationReferenceTotalT0T1T2Amount, false)
			}
			if negativeMiningRate.IsNil() {
				negativeMiningRate = coin.ZeroICEFlakes()
			}
		}
		res = append(res, s.calculateICEFlakeMiningRates(currentAdoption.BaseMiningRate, row, negativeMiningRate))
	}

	return res, nil
}

type (
	latestMiningRateCalculationSQLRow struct {
		LastMiningEndedAt                         *time.Time
		AggressiveDegradationReferenceTotalAmount *coin.ICEFlake
		AggressiveDegradationReferenceT0Amount    *coin.ICEFlake
		AggressiveDegradationReferenceT1Amount    *coin.ICEFlake
		AggressiveDegradationReferenceT2Amount    *coin.ICEFlake
		DegradationReferenceTotalT0T1T2Amount     *coin.ICEFlake
		userMiningRateRecalculationParameters
	}
)

func (s *miningRatesRecalculationTriggerStreamSource) getUserMiningRateCalculationParametersNewBatch( //nolint:funlen,gocritic// .
	ctx context.Context, workerIndex int16, now *time.Time,
) ([]*latestMiningRateCalculationSQLRow, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`
SELECT u.last_mining_ended_at,
	   aggressive_degradation_btotal.amount AS aggressive_degradation_reference_total_amount,
	   aggressive_degradation_bt0.amount AS aggressive_degradation_reference_t0_amount,
	   aggressive_degradation_bt1.amount AS aggressive_degradation_reference_t1_amount,
	   aggressive_degradation_bt2.amount AS aggressive_degradation_reference_t2_amount,
	   degradation_btotalt0t1t2.amount AS degradation_reference_total_t0_t1_t2_amount,
	   u.user_id,
	   (CASE WHEN t0.user_id IS NULL THEN 0 ELSE 1 END) AS t0,
	   coalesce(ar_worker.t1,0) AS t1,
	   coalesce(ar_worker.t2,0) AS t2,
	   (CASE WHEN coalesce(eb_worker.extra_bonus_ended_at, '1999-01-08 04:05:06'::timestamp) > $3 THEN eb_worker.extra_bonus ELSE 0 END) AS extra_bonus, 
	   coalesce(x.pre_staking_allocation,0) AS pre_staking_allocation,
	   coalesce(st_b.bonus,0) AS pre_staking_bonus
FROM (SELECT MAX(st.years) AS pre_staking_years,
		     MAX(st.allocation) AS pre_staking_allocation,
			 x.user_id
	  FROM ( SELECT user_id
		     FROM mining_rates_recalculation_worker
			 WHERE worker_index = $1
		     ORDER BY last_iteration_finished_at
		     LIMIT $2 ) x
			 LEFT JOIN pre_stakings st
					ON st.worker_index = $1
				   AND st.user_id = x.user_id
	  GROUP BY x.user_id
	 ) x
	    JOIN users u
		  ON u.user_id = x.user_id
   LEFT JOIN extra_bonus_processing_worker eb_worker
		  ON eb_worker.worker_index = $1
		 AND eb_worker.user_id = x.user_id
   LEFT JOIN active_referrals ar_worker
		  ON ar_worker.worker_index = $1
		 AND ar_worker.user_id = x.user_id
   LEFT JOIN pre_staking_bonuses st_b
		  ON st_b.years = x.pre_staking_years
   LEFT JOIN balances_worker aggressive_degradation_btotal
		  ON (u.last_mining_ended_at IS NOT NULL AND u.last_mining_ended_at < $3 )
		 AND aggressive_degradation_btotal.user_id = u.user_id
		 AND aggressive_degradation_btotal.negative = FALSE
		 AND aggressive_degradation_btotal.type = %[1]v
		 AND aggressive_degradation_btotal.type_detail = '%[2]v'
   LEFT JOIN balances_worker aggressive_degradation_bt0
		  ON (u.last_mining_ended_at IS NOT NULL AND u.last_mining_ended_at < $3 )
		 AND aggressive_degradation_bt0.worker_index = $1
		 AND aggressive_degradation_bt0.user_id = u.user_id
		 AND aggressive_degradation_bt0.negative = FALSE
		 AND aggressive_degradation_bt0.type = %[1]v
		 AND aggressive_degradation_bt0.type_detail = '%[3]v_' || u.referred_by || '_'
   LEFT JOIN balances_worker aggressive_degradation_bt1
		  ON (u.last_mining_ended_at IS NOT NULL AND u.last_mining_ended_at < $3 )
		 AND aggressive_degradation_bt1.worker_index = $1
		 AND aggressive_degradation_bt1.user_id = u.user_id
		 AND aggressive_degradation_bt1.negative = FALSE
		 AND aggressive_degradation_bt1.type = %[1]v
		 AND aggressive_degradation_bt1.type_detail = '%[4]v'
   LEFT JOIN balances_worker aggressive_degradation_bt2
		  ON (u.last_mining_ended_at IS NOT NULL AND u.last_mining_ended_at < $3 )
		 AND aggressive_degradation_bt2.worker_index = $1
		 AND aggressive_degradation_bt2.user_id = u.user_id
		 AND aggressive_degradation_bt2.negative = FALSE
		 AND aggressive_degradation_bt2.type = %[1]v
		 AND aggressive_degradation_bt2.type_detail = '%[5]v'
   LEFT JOIN balances_worker degradation_btotalt0t1t2
		  ON (u.last_mining_ended_at IS NOT NULL AND u.last_mining_ended_at < $3 )
		 AND degradation_btotalt0t1t2.worker_index = $1
		 AND degradation_btotalt0t1t2.user_id = u.user_id
		 AND degradation_btotalt0t1t2.negative = FALSE
		 AND degradation_btotalt0t1t2.type = %[1]v
		 AND degradation_btotalt0t1t2.type_detail = '%[6]v'
   LEFT JOIN users t0
	  	  ON t0.user_id = u.referred_by
	     AND t0.user_id != x.user_id
	  	 AND t0.last_mining_ended_at IS NOT NULL
	  	 AND t0.last_mining_ended_at  > $3`,
		totalNoPreStakingBonusBalanceType,
		aggressiveDegradationTotalReferenceBalanceTypeDetail,
		t0BalanceTypeDetail,
		aggressiveDegradationT1ReferenceBalanceTypeDetail,
		aggressiveDegradationT2ReferenceBalanceTypeDetail,
		degradationT0T1T2TotalReferenceBalanceTypeDetail)
	res, err := storage.Select[latestMiningRateCalculationSQLRow](ctx, s.db, sql, workerIndex, s.cfg.Workers.MiningRatesRecalculationBatchSize, *now.Time)

	return res, errors.Wrapf(err, "failed to select a batch of latest user mining rate calculation parameters for workerIndex:%v", workerIndex)
}

func (s *miningRatesRecalculationTriggerStreamSource) sendMiningRatesMessage(ctx context.Context, mnrs *MiningRates[coin.ICEFlake]) error {
	valueBytes, err := json.MarshalContext(ctx, mnrs)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %#v", mnrs)
	}

	msg := &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     mnrs.UserID,
		Topic:   s.cfg.MessageBroker.Topics[3].Name,
		Value:   valueBytes,
	}

	responder := make(chan error, 1)
	defer close(responder)
	s.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send %v message to broker, msg:%#v", msg.Topic, mnrs)
}

func (r *repository) calculateICEFlakeMiningRates(
	baseMiningRate *coin.ICEFlake, row *latestMiningRateCalculationSQLRow, negativeMiningRate *coin.ICEFlake,
) (miningRates *MiningRates[coin.ICEFlake]) {
	miningRates = new(MiningRates[coin.ICEFlake])
	miningRates.UserID = row.UserID
	if !negativeMiningRate.IsNil() {
		miningRates.Type = NegativeMiningRateType
		if row.PreStakingAllocation != percentage100 {
			miningRates.Standard = negativeMiningRate.
				MultiplyUint64(percentage100 - row.PreStakingAllocation).
				DivideUint64(percentage100)
		}
		if row.PreStakingAllocation != 0 {
			miningRates.PreStaking = negativeMiningRate.
				MultiplyUint64((row.PreStakingBonus + percentage100) * row.PreStakingAllocation).
				DivideUint64(percentage100 * percentage100)
		}

		return miningRates
	}

	miningRates.Type = PositiveMiningRateType
	params := &row.userMiningRateRecalculationParameters
	if standard := r.calculateMintedStandardCoins(baseMiningRate, params, r.cfg.GlobalAggregationInterval.Child, false); !standard.IsZero() {
		miningRates.Standard = standard
	}
	if preStaking := r.calculateMintedPreStakingCoins(baseMiningRate, params, r.cfg.GlobalAggregationInterval.Child, false); !preStaking.IsZero() {
		miningRates.PreStaking = preStaking
	}

	return miningRates
}

func (s *miningRatesRecalculationTriggerStreamSource) updateLastIterationFinishedAt(
	ctx context.Context, workerIndex int16, rows []*MiningRates[coin.ICEFlake], now *time.Time,
) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	userIDs := make([]string, 0, len(rows))
	for i := range rows {
		userIDs = append(userIDs, rows[i].UserID)
	}
	const table = "mining_rates_recalculation_worker"
	params := make(map[string]any, 1)
	params["last_iteration_finished_at"] = *now.Time
	err := s.updateWorkerFields(ctx, workerIndex, table, params, userIDs...)

	return errors.Wrapf(err, "failed to updateWorkerTimeField for workerIndex:%v,table:%q,params:%#v,userIDs:%#v", workerIndex, table, params, userIDs)
}
