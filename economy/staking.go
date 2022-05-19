// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"time"

	"github.com/framey-io/go-tarantool"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
)

func (e *economy) StartStaking(ctx context.Context, userID UserID, staking Staking) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "start mining failed because context failed")
	}

	stakingEnabled, err := e.isStakingEnabled(userID)
	if err != nil {
		return errors.Wrap(err, "unable to check is staking enabled")
	}

	if stakingEnabled {
		return ErrStakingAlreadyEnabled
	}

	nowUtc := time.Now().UTC()

	err = e.enableStaking(userID, staking, nowUtc)
	if err != nil {
		return errors.Wrap(err, "unable to enable staking")
	}

	return errors.Wrap(e.notifyStartStaking(ctx, userID, staking, nowUtc), "failed to notify that the user enable staking")
}

func (e *economy) isStakingEnabled(userID UserID) (bool, error) {
	params := map[string]interface{}{
		"userID": userID,
	}

	sql := `SELECT staking_years > 0 AND staking_percentage > 0 
		FROM USER_ECONOMY INDEXED BY "pk_unnamed_USER_ECONOMY_1" 
		WHERE user_id = :userID`

	var res []stakingAlreadyEnabled
	if err := e.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return false, errors.Wrapf(err, "failed to get user_economy record with userID %v", userID)
	}

	if len(res) == 0 {
		return false, ErrNotFound
	}

	return res[0].Value, nil
}

func (e *economy) enableStaking(userID string, staking Staking, updatedAt time.Time) error {
	space := "USER_ECONOMY"
	index := "pk_unnamed_USER_ECONOMY_1"
	key := tarantool.StringKey{S: userID}

	//nolint:gomnd // Those are not magic numbers, those are the indexes of the fields.
	ops := []tarantool.Op{
		{Op: "=", Field: 4, Arg: staking.Percentage},
		{Op: "=", Field: 7, Arg: staking.Years},
		{Op: "=", Field: 9, Arg: updatedAt.UnixNano()},
	}

	return errors.Wrapf(e.db.UpdateTyped(space, index, key, ops, &[]*userEconomy{}),
		"failed set staking_years:%v, staking_persentage:%v for userID:%v", staking.Years, staking.Percentage, userID)
}

func (e *economy) notifyStartStaking(ctx context.Context, userID UserID, staking Staking, startedAt time.Time) error {
	m := stakingEnabled{
		TS:      startedAt,
		Staking: staking,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(err, "[start-staking] failed to marshal %#v", m)
	}

	responder := make(chan error, 1)
	e.mb.SendMessage(ctx, &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     userID,
		Topic:   cfg.MessageBroker.Topics[1].Name,
		Value:   b,
	}, responder)

	return errors.Wrapf(<-responder, "[start-staking] failed to send message to broker")
}
