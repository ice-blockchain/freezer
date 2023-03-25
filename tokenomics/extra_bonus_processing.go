// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	stdlibtime "time"

	"github.com/cenkalti/backoff/v4"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) initializeExtraBonusProcessingWorker(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	workerIndex := usr.HashCode % r.cfg.WorkerCount
	err := retry(ctx, func() error {
		if err := r.initializeWorker(ctx, "extra_bonus_processing_worker_", usr.ID, workerIndex); err != nil {
			if errors.Is(err, storage.ErrRelationNotFound) {
				return err
			}

			return errors.Wrapf(backoff.Permanent(err),
				"failed to initializeExtraBonusProcessingWorker for userID:%v,workerIndex:%v", usr.ID, workerIndex)
		}

		return nil
	})

	return errors.Wrapf(err, "permanently failed to initializeExtraBonusProcessingWorker for userID:%v,workerIndex:%v", usr.ID, workerIndex)
}

func (p *processor) startExtraBonusProcessingTriggerSeedingStream(ctx context.Context) {
	nilBodyForEachWorker := make([]any, p.cfg.WorkerCount) //nolint:makezero // Intended.
	ticker := stdlibtime.NewTicker(extraBonusProcessingSeedingStreamEmitFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Error(errors.Wrap(sendMessagesConcurrently[any](ctx, p.sendExtraBonusProcessingTriggerMessage, nilBodyForEachWorker),
				"failed to sendMessagesConcurrently[sendExtraBonusProcessingTriggerMessage]"))
		case <-ctx.Done():
			return
		}
	}
}

func (p *processor) sendExtraBonusProcessingTriggerMessage(ctx context.Context, _ any) error {
	msg := &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     uuid.NewString(),
		Topic:   p.cfg.MessageBroker.Topics[12].Name,
	}
	responder := make(chan error, 1)
	defer close(responder)
	p.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send `%v` message to broker", msg.Topic)
}

func (s *extraBonusProcessingTriggerStreamSource) Process(ignoredCtx context.Context, msg *messagebroker.Message) (err error) {
	if ignoredCtx.Err() != nil {
		return errors.Wrap(ignoredCtx.Err(), "unexpected deadline while processing message")
	}
	if true {
		return nil
	}
	const deadline = 5 * stdlibtime.Minute
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	extraBonusIndex, availableExtraBonuses, err := s.getAvailableExtraBonuses(ctx, uint64(msg.Partition)) //nolint:contextcheck // Not needed here.
	if err != nil {
		return errors.Wrapf(err, "failed to getAvailableExtraBonuses for workerIndex:%v", uint64(msg.Partition))
	}
	if err = sendMessagesConcurrently(ctx, s.sendAvailableDailyBonusMessage, availableExtraBonuses); err != nil { //nolint:contextcheck // Not needed here.
		return errors.Wrapf(err, "failed to sendMessagesConcurrently[sendAvailableDailyBonusMessage] for availableExtraBonuses:%#v", availableExtraBonuses)
	}
	const table = "extra_bonus_processing_worker_"
	params := make(map[string]any, 1)
	params["last_extra_bonus_index_notified"] = extraBonusIndex
	userIDs := make([]string, 0, len(availableExtraBonuses))
	for _, bonus := range availableExtraBonuses {
		userIDs = append(userIDs, bonus.UserID)
	}

	return errors.Wrapf(s.updateWorkerFields(ctx, uint64(msg.Partition), table, params, userIDs...), //nolint:contextcheck // Not needed here.
		"failed to updateWorkerTimeField for workerIndex:%v,table:%q,params:%#v,userIDs:%#v", uint64(msg.Partition), table, params, userIDs)
}

//nolint:funlen,lll // .
func (s *extraBonusProcessingTriggerStreamSource) getAvailableExtraBonuses(
	ctx context.Context, workerIndex uint64,
) (extraBonusIndex uint64, availableExtraBonuses []*ExtraBonusSummary, err error) {
	if ctx.Err() != nil {
		return 0, nil, errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	sql := fmt.Sprintf(`WITH sd AS (SELECT value FROM extra_bonus_start_date WHERE key = 0)
						SELECT bal_worker.last_mining_started_at,
							   bal_worker.last_mining_ended_at,
							   eb_worker.user_id,
							   eb_worker.news_seen,
							   b.bonus,
							   (100 - (25 *  ((CASE WHEN (:now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - (sd.value + (b.ix * :duration)) - :time_to_availability_window - ((e.offset * :availability_window) / :worker_count)) < :first_delayed_claim_penalty_window THEN 0 ELSE (:now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - (sd.value + (b.ix * :duration)) - :time_to_availability_window - ((e.offset * :availability_window) / :worker_count)) END)/:delayed_claim_penalty_window))) AS bonus_percentage_remaining,
							   b.ix
						FROM extra_bonus_processing_worker_%[1]v eb_worker
							JOIN balance_recalculation_worker_%[1]v bal_worker
							  ON bal_worker.user_id = eb_worker.user_id
							JOIN sd 
							JOIN extra_bonuses b 
							  ON b.ix = (:now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - sd.value) / :duration
		 					 AND :now_nanos + (eb_worker.utc_offset * :utc_offset_duration) > sd.value
							 AND b.bonus > 0
							JOIN extra_bonuses_%[1]v e
							  ON e.extra_bonus_index = b.ix
							 AND :now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - (sd.value + (e.extra_bonus_index * :duration)) - :time_to_availability_window - ((e.offset * :availability_window) / :worker_count) < :claim_window
							 AND :now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - (sd.value + (e.extra_bonus_index * :duration)) - :time_to_availability_window - ((e.offset * :availability_window) / :worker_count) > 0
						WHERE (eb_worker.last_extra_bonus_index_notified IS NULL OR eb_worker.last_extra_bonus_index_notified < b.ix)
					      AND :now_nanos - IFNULL(eb_worker.extra_bonus_started_at, 0) > :claim_window 
						LIMIT %[2]v`, workerIndex, extraBonusProcessingBatchSize)
	now := time.Now()
	const networkLagDelta = 1.3
	params := make(map[string]any, 9) //nolint:gomnd // .
	params["now_nanos"] = now
	params["duration"] = s.cfg.ExtraBonuses.Duration
	params["utc_offset_duration"] = s.cfg.ExtraBonuses.UTCOffsetDuration
	params["availability_window"] = s.cfg.ExtraBonuses.AvailabilityWindow
	params["time_to_availability_window"] = s.cfg.ExtraBonuses.TimeToAvailabilityWindow
	params["claim_window"] = s.cfg.ExtraBonuses.ClaimWindow
	params["worker_count"] = s.cfg.WorkerCount
	params["delayed_claim_penalty_window"] = s.cfg.ExtraBonuses.DelayedClaimPenaltyWindow
	params["first_delayed_claim_penalty_window"] = stdlibtime.Duration(float64(s.cfg.ExtraBonuses.DelayedClaimPenaltyWindow.Nanoseconds()) * networkLagDelta)
	resp := make([]*struct {
		_msgpack                                                       struct{} `msgpack:",asArray"` //nolint:tagliatelle,revive,nosnakecase // .
		LastMiningStartedAt, LastMiningEndedAt                         *time.Time
		UserID                                                         string
		NewsSeen, FlatBonus, BonusPercentageRemaining, ExtraBonusIndex uint64
	}, 0, extraBonusProcessingBatchSize)
	if err = s.db.PrepareExecuteTyped(sql, params, &resp); err != nil {
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
