// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/go-tarantool-client"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetPreStakingSummary(ctx context.Context, userID string) (*PreStakingSummary, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`SELECT st.created_at,
							   st.user_id,
							   MAX(st.years) AS years,
							   MAX(st.allocation) AS allocation,
							   st_b.bonus	
						FROM pre_stakings_%[1]v st
							LEFT JOIN pre_staking_bonuses st_b
								   ON st.years = st_b.years
						WHERE st.user_id = :user_id
						GROUP BY st.user_id`, r.workerIndex(ctx))
	params := make(map[string]any, 1)
	params["user_id"] = userID
	resp := make([]*PreStakingSummary, 0, 1)
	if err := r.db.PrepareExecuteTyped(sql, params, &resp); err != nil {
		return nil, errors.Wrapf(err, "failed to select for pre-staking summary for userID:%v", userID)
	}
	if len(resp) == 0 {
		return nil, storage.ErrNotFound
	}
	resp[0].UserID = ""
	resp[0].CreatedAt = nil

	return resp[0], nil
}

func (r *repository) getAllPreStakingSummaries(ctx context.Context, userID string) (resp []*PreStakingSummary, err error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`SELECT st.*,
							   st_b.bonus	
						FROM pre_stakings_%[1]v st
							JOIN pre_staking_bonuses st_b
							  ON st.years = st_b.years
						WHERE st.user_id = :user_id
						ORDER BY st.created_at`, r.workerIndex(ctx))
	params := make(map[string]any, 1)
	params["user_id"] = userID
	err = errors.Wrapf(r.db.PrepareExecuteTyped(sql, params, &resp), "failed to select all pre-staking summaries for userID:%v", userID)

	return
}

func (r *repository) StartOrUpdatePreStaking(ctx context.Context, st *PreStakingSummary) error { //nolint:funlen,gocognit // Can't properly split it further.
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	existing, err := r.GetPreStakingSummary(ctx, st.UserID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
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
	st.CreatedAt = time.Now()
	sql := fmt.Sprintf(`INSERT INTO pre_stakings_%[1]v (created_at,user_id,years,allocation) 
												VALUES (:created_at,:user_id,:years,:allocation)`, r.workerIndex(ctx))
	params := make(map[string]any, 1+1+1+1)
	params["created_at"] = st.CreatedAt
	params["user_id"] = st.UserID
	params["years"] = st.Years
	params["allocation"] = st.Allocation
	if err = storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)); err != nil {
		return errors.Wrapf(err, "failed to insertNewPreStaking:%#v", st)
	}
	var res struct {
		_msgpack     struct{} `msgpack:",asArray"` //nolint:tagliatelle,revive,nosnakecase // To insert we need asArray
		Years, Bonus uint64
	}
	if err = r.db.GetTyped("PRE_STAKING_BONUSES", "pk_unnamed_PRE_STAKING_BONUSES_1", tarantool.UintKey{I: uint(st.Years)}, &res); err != nil {
		return errors.Wrapf(err, "failed to get pre-staking bonus for years:%v", st.Years)
	}
	st.Bonus = res.Bonus
	ss := &PreStakingSnapshot{PreStakingSummary: st, Before: existing}
	if err = r.sendPreStakingSnapshotMessage(ctx, ss); err != nil {
		pkIndex := fmt.Sprintf("pk_unnamed_PRE_STAKINGS_%v_1", r.workerIndex(ctx))
		key := []any{st.UserID, st.Years, st.Allocation}

		return multierror.Append( //nolint:wrapcheck // Not needed.
			errors.Wrapf(err, "failed to send pre-staking snapshot message:%#v", ss),
			errors.Wrapf(r.db.DeleteTyped(fmt.Sprintf("PRE_STAKINGS_%v", r.workerIndex(ctx)), pkIndex, key, &[]*PreStaking{}),
				"failed to revertInsertNewPreStaking: %#v", st),
		).ErrorOrNil()
	}
	st.UserID = ""
	st.CreatedAt = nil

	return nil
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
