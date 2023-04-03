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
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) initializeExtraBonusProcessingWorker(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}

	return errors.Wrapf(retry(ctx, func() error {
		if err := r.initializeWorker(ctx, "extra_bonus_processing_worker", usr.ID, usr.HashCode); err != nil {
			if errors.Is(err, storage.ErrRelationNotFound) {
				return err
			}

			return errors.Wrapf(backoff.Permanent(err),
				"failed to initializeExtraBonusProcessingWorker for userID:%v", usr.ID)
		}

		return nil
	}), "permanently failed to initializeExtraBonusProcessingWorker for userID:%v", usr.ID)
}

func (s *extraBonusProcessingTriggerStreamSource) start(ctx context.Context) {
	log.Info("extraBonusProcessingTriggerStreamSource started")
	defer log.Info("extraBonusProcessingTriggerStreamSource stopped")
	workerIndexes := make([]int16, s.cfg.WorkerCount) //nolint:makezero // Intended.
	for i := 0; i < int(s.cfg.WorkerCount); i++ {
		workerIndexes[i] = int16(i)
	}
	for ctx.Err() == nil {
		stdlibtime.Sleep(s.cfg.Workers.ExtraBonusProcessingSeedingStreamEmitFrequency)
		before := time.Now()
		log.Error(errors.Wrap(executeBatchConcurrently(ctx, s.process, workerIndexes), "failed to executeBatchConcurrently[extraBonusProcessingTriggerStreamSource.process]")) //nolint:lll // .
		log.Error(fmt.Errorf("extraBonusProcessingTriggerStreamSource.process took: %v", stdlibtime.Since(*before.Time)))
	}
}

func (s *extraBonusProcessingTriggerStreamSource) process(ignoredCtx context.Context, workerIndex int16) (err error) {
	if ignoredCtx.Err() != nil {
		return errors.Wrap(ignoredCtx.Err(), "unexpected deadline while processing message")
	}
	const deadline = 5 * stdlibtime.Minute
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	extraBonusIndex, availableExtraBonuses, err := s.getAvailableExtraBonuses(ctx, workerIndex) //nolint:contextcheck // Not needed here.
	if err != nil {
		return errors.Wrapf(err, "failed to getAvailableExtraBonuses for workerIndex:%v", workerIndex)
	}
	if err = executeBatchConcurrently(ctx, s.sendAvailableDailyBonusMessage, availableExtraBonuses); err != nil { //nolint:contextcheck // Not needed here.
		return errors.Wrapf(err, "failed to executeBatchConcurrently[sendAvailableDailyBonusMessage] for availableExtraBonuses:%#v", availableExtraBonuses)
	}
	const table = "extra_bonus_processing_worker"
	params := make(map[string]any, 1)
	params["last_extra_bonus_index_notified"] = extraBonusIndex
	userIDs := make([]string, 0, len(availableExtraBonuses))
	for _, bonus := range availableExtraBonuses {
		userIDs = append(userIDs, bonus.UserID)
	}

	return errors.Wrapf(s.updateWorkerFields(ctx, workerIndex, table, params, userIDs...), //nolint:contextcheck // Not needed here.
		"failed to updateWorkerTimeField for workerIndex:%v,table:%q,params:%#v,userIDs:%#v", workerIndex, table, params, userIDs)
}

//nolint:funlen,lll // .
func (s *extraBonusProcessingTriggerStreamSource) getAvailableExtraBonuses(
	ctx context.Context, workerIndex int16,
) (extraBonusIndex uint64, availableExtraBonuses []*ExtraBonusSummary, err error) {
	if ctx.Err() != nil {
		return 0, nil, errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	sql := `WITH sd AS (SELECT value FROM extra_bonus_start_date WHERE key = 0)
			SELECT bal_worker.last_mining_started_at,
				   bal_worker.last_mining_ended_at,
				   eb_worker.user_id,
				   eb_worker.news_seen,
				   b.bonus AS flat_bonus,
				   (100 - (25 *  ((CASE WHEN ($3::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (e.extra_bonus_index * $5::bigint)) - $6::bigint - ((e.offset_value * $7::bigint) / $9)) < $11::bigint THEN 0 ELSE ($3::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (e.extra_bonus_index * $5::bigint)) - $6::bigint - ((e.offset_value * $7::bigint) / $9)) END)/$10::bigint))) AS bonus_percentage_remaining,
				   b.ix AS extra_bonus_index
			FROM extra_bonus_processing_worker eb_worker
				JOIN balance_recalculation_worker bal_worker
				  ON bal_worker.worker_index = $1
				 AND bal_worker.user_id = eb_worker.user_id
				JOIN sd 
				  ON 1=1
				JOIN extra_bonuses b 
				  ON b.ix = ($3::bigint + (eb_worker.utc_offset * $4::bigint) - sd.value) / $5::bigint
				 AND $3::bigint + (eb_worker.utc_offset * $4::bigint) > sd.value
				 AND b.bonus > 0
				JOIN extra_bonuses_worker e
				  ON e.worker_index = $1
				 AND e.extra_bonus_index = b.ix
				 AND $3::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (e.extra_bonus_index * $5::bigint)) - $6::bigint - ((e.offset_value * $7::bigint) / $9) < $8::bigint
				 AND $3::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (e.extra_bonus_index * $5::bigint)) - $6::bigint - ((e.offset_value * $7::bigint) / $9) > 0
			WHERE eb_worker.worker_index = $1
			  AND (eb_worker.last_extra_bonus_index_notified IS NULL OR eb_worker.last_extra_bonus_index_notified < b.ix)
			  AND $12 > coalesce(eb_worker.extra_bonus_started_at, '1999-01-08 04:05:06'::timestamp) 
			LIMIT $2`
	now := time.Now()
	const networkLagDelta, argCount = 1.3, 12
	args := append(make([]any, 0, argCount),
		workerIndex,
		s.cfg.Workers.ExtraBonusProcessingBatchSize,
		now.UnixNano(),
		s.cfg.ExtraBonuses.UTCOffsetDuration,
		s.cfg.ExtraBonuses.Duration,
		s.cfg.ExtraBonuses.TimeToAvailabilityWindow,
		s.cfg.ExtraBonuses.AvailabilityWindow,
		s.cfg.ExtraBonuses.ClaimWindow,
		s.cfg.WorkerCount,
		s.cfg.ExtraBonuses.DelayedClaimPenaltyWindow,
		stdlibtime.Duration(float64(s.cfg.ExtraBonuses.DelayedClaimPenaltyWindow.Nanoseconds())*networkLagDelta),
		now.Add(-s.cfg.ExtraBonuses.ClaimWindow))
	resp, err := storage.Select[struct {
		LastMiningStartedAt, LastMiningEndedAt                         *time.Time
		UserID                                                         string
		NewsSeen, FlatBonus, BonusPercentageRemaining, ExtraBonusIndex uint64
	}](ctx, s.db, sql, args...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "failed to select for availableExtraBonuses for workerIndex:%v", workerIndex)
	}
	if len(resp) != 0 {
		extraBonusIndex = resp[0].ExtraBonusIndex
	}
	availableExtraBonuses = make([]*ExtraBonusSummary, 0, len(resp))
	for _, row := range resp {
		availableExtraBonuses = append(availableExtraBonuses, &ExtraBonusSummary{
			UserID:              row.UserID,
			AvailableExtraBonus: s.calculateExtraBonus(row.FlatBonus, row.BonusPercentageRemaining, row.NewsSeen, s.calculateMiningStreak(now, row.LastMiningStartedAt, row.LastMiningEndedAt)), //nolint:lll // .
			ExtraBonusIndex:     row.ExtraBonusIndex,
		})
	}

	return extraBonusIndex, availableExtraBonuses, nil
}

func (s *extraBonusProcessingTriggerStreamSource) sendAvailableDailyBonusMessage(ctx context.Context, ebs *ExtraBonusSummary) error {
	valueBytes, err := json.MarshalContext(ctx, ebs)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %#v", ebs)
	}

	msg := &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     ebs.UserID,
		Topic:   s.cfg.MessageBroker.Topics[7].Name,
		Value:   valueBytes,
	}

	responder := make(chan error, 1)
	defer close(responder)
	s.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send %v message to broker, msg:%#v", msg.Topic, ebs)
}
