// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetPreStakingSummary(ctx context.Context, userID string) (resp *PreStakingSummary, err error) {
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
	resp, err = storage.Get[PreStakingSummary](ctx, r.db, sql, r.workerIndex(ctx), userID)

	return resp, errors.Wrapf(err, "failed to select for pre-staking summary for userID:%v", userID)
}

func (r *repository) getAllPreStakingSummaries(ctx context.Context, userID string) (resp []*PreStakingSummary, err error) {
	sql := `SELECT st.*,
				   st_b.bonus	
			FROM pre_stakings st
				JOIN pre_staking_bonuses st_b
			      ON st.worker_index = $1 
			     AND st.years = st_b.years
			WHERE st.worker_index = $1 
			  AND st.user_id = $2
			ORDER BY st.created_at`
	resp, err = storage.Select[PreStakingSummary](ctx, r.db, sql, r.workerIndex(ctx), userID)

	return resp, errors.Wrapf(err, "failed to select all pre-staking summaries for userID:%v", userID)
}

func (r *repository) StartOrUpdatePreStaking(ctx context.Context, st *PreStakingSummary) error { //nolint:funlen,gocognit // Can't properly split it further.
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
	res, err := storage.Get[struct{ Bonus uint64 }](ctx, r.db, `SELECT bonus FROM pre_staking_bonuses WHERE years = $1`, st.Years)
	if err != nil {
		return errors.Wrapf(err, "failed to get pre-staking bonus for years:%v", st.Years)
	}
	st.CreatedAt, st.Bonus = time.Now(), res.Bonus

	return errors.Wrap(storage.DoInTransaction(ctx, r.db, func(conn storage.QueryExecer) error {
		sql := `INSERT INTO pre_stakings (created_at, user_id, years, allocation, hash_code , worker_index) 
							      VALUES ($1        , $2     , $3   , $4        , $5::bigint, $6)`
		if _, err = storage.Exec(ctx, conn, sql, *st.CreatedAt.Time, st.UserID, st.Years, st.Allocation, r.hashCode(ctx), r.workerIndex(ctx)); err != nil {
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
