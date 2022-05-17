// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
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
		return ErrStakingEnabled
	}

	nowUtc := time.Now().UTC()
	nowNano := uint64(nowUtc.UnixNano())

	err = e.enableStaking(userID, staking, nowNano)
	if err != nil {
		return errors.Wrap(err, "unable to enable staking")
	}

	return errors.Wrap(e.notifyStartStaking(ctx, userID, staking, nowUtc), "failed to notify that the user enable staking")
}

func (e *economy) isStakingEnabled(userID UserID) (bool, error) {
	space := userEconomySpace()

	params := map[string]interface{}{
		"userID": userID,
	}

	sql := fmt.Sprintf(`SELECT staking_years, staking_percentage 
		FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1" 
		WHERE user_id = :userID `, space)

	var res []*Staking
	if err := e.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return false, errors.Wrapf(err, "failed to get %q record with userID %v", space, userID)
	}

	if len(res) == 0 {
		return false, ErrNotFound
	}

	return res[0].IsValid(), nil
}

func (s Staking) IsValid() bool {
	return (s.Years >= 1 && s.Years <= 5) &&
		(s.Percentage > 0.0 && s.Percentage <= 100)
}

func (e *economy) enableStaking(userID string, staking Staking, startTime uint64) error {
	params := map[string]interface{}{
		"userId":     userID,
		"years":      staking.Years,
		"percentage": staking.Percentage,
		"updatedAt":  startTime,
	}

	sql := fmt.Sprintf(`
		UPDATE %[1]v 
		SET staking_years = :years, 
			staking_percentage = :percentage, 
			updated_at = :updatedAt 
		WHERE user_id = :userId`, userEconomySpace())

	if err := storage.CheckSQLDMLErr(e.db.PrepareExecute(sql, params)); err != nil {
		return errors.Wrapf(err, "failed set staking_years:%v, staking_persentage:%v for userID:%v", staking.Years, staking.Percentage, userID)
	}

	return nil
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
		Topic:   cfg.MessageBroker.Topics[0].Name,
		Value:   b,
	}, responder)

	return errors.Wrapf(<-responder, "[start-staking] failed to send message to broker")
}
