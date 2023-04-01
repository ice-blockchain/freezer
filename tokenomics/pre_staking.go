// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetPreStakingSummary(ctx context.Context, userID string) (resp *PreStakingSummary, err error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := `SELECT st.created_at,
				   st.user_id,
				   MAX(st.years) AS years,
				   MAX(st.allocation) AS allocation,
				   st.hash_code,	
				   st.worker_index,	
				   coalesce(st_b.bonus,0) AS bonus
			FROM pre_stakings st
				LEFT JOIN pre_staking_bonuses st_b
					   ON st.worker_index = $1 
					  AND st.years = st_b.years
			WHERE st.worker_index = $1 
			  AND st.user_id = $2
			GROUP BY st.worker_index,st.user_id,st.hash_code,st.created_at,st_b.bonus`
	if resp, err = storagev2.Get[PreStakingSummary](ctx, r.dbV2, sql, r.workerIndex(ctx), userID); err != nil {
		return nil, errors.Wrapf(err, "failed to select for pre-staking summary for userID:%v", userID)
	}

	return resp, nil
}

func (r *repository) getAllPreStakingSummaries(ctx context.Context, userID string) (resp []*PreStakingSummary, err error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := `SELECT st.*,
				   st_b.bonus	
			FROM pre_stakings st
				JOIN pre_staking_bonuses st_b
			      ON st.worker_index = $1 
			     AND st.years = st_b.years
			WHERE st.worker_index = $1 
			  AND st.user_id = $2
			ORDER BY st.created_at`
	if resp, err = storagev2.Select[PreStakingSummary](ctx, r.dbV2, sql, r.workerIndex(ctx), userID); err != nil {
		return nil, errors.Wrapf(err, "failed to select all pre-staking summaries for userID:%v", userID)
	}

	return
}

func (r *repository) StartOrUpdatePreStaking(ctx context.Context, st *PreStakingSummary) error { //nolint:funlen,gocognit // Can't properly split it further.
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	existing, err := r.GetPreStakingSummary(ctx, st.UserID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return errors.Wrapf(err, "failed to GetPreStakingSummary for userID:%v", st.UserID)
	}
	if existing != nil {
		if (existing.Allocation == percentage100 || existing.Allocation == st.Allocation) &&
			(existing.Years == MaxPreStakingYears || existing.Years == st.Years) {
			*st = *existing

			return nil
		}
		if existing.Allocation > st.Allocation || existing.Years > st.Years {
			return ErrDecreasingPreStakingAllocationOrYearsNotAllowed
		}
	}
	res, err := storagev2.Get[struct{ Bonus uint64 }](ctx, r.dbV2, `SELECT bonus FROM pre_staking_bonuses WHERE years = $1`, st.Years)
	if err != nil {
		return errors.Wrapf(err, "failed to get pre-staking bonus for years:%v", st.Years)
	}
	st.CreatedAt, st.Bonus = time.Now(), res.Bonus

	return errors.Wrap(storagev2.DoInTransaction(ctx, r.dbV2, func(conn storagev2.QueryExecer) error {
		sql := `INSERT INTO pre_stakings (created_at, user_id, years, allocation, hash_code , worker_index) 
							      VALUES ($1        , $2     , $3   , $4        , $5::bigint, $6)`
		if _, err = storagev2.Exec(ctx, conn, sql, *st.CreatedAt.Time, st.UserID, st.Years, st.Allocation, r.hashCode(ctx), r.workerIndex(ctx)); err != nil {
			return errors.Wrapf(err, "failed to insertNewPreStaking:%#v", st)
		}
		ss := &PreStakingSnapshot{PreStakingSummary: st, Before: existing}

		return errors.Wrapf(r.sendPreStakingSnapshotMessage(ctx, ss), "failed to send pre-staking snapshot message:%#v", ss)
	}), "DoInTransaction failed")
}

func (r *repository) sendPreStakingSnapshotMessage(ctx context.Context, st *PreStakingSnapshot) error {
	valueBytes, err := json.MarshalContext(ctx, st)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %#v", st)
	}
	msg := &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     st.UserID,
		Topic:   r.cfg.MessageBroker.Topics[6].Name,
		Value:   valueBytes,
	}
	responder := make(chan error, 1)
	defer close(responder)
	r.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send `%v` message to broker", msg.Topic)
}
