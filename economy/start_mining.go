// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/time"
)

func (e *economy) StartMining(ctx context.Context, userID UserID) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "start mining failed because context failed")
	}

	miningInProgress, err := e.isMiningInProgress(userID)
	if err != nil {
		return errors.Wrap(err, "unable to check is mining in porgress")
	}

	if miningInProgress {
		return ErrMiningInProgress
	}

	nowUtc := time.Now()

	err = e.startMining(userID, nowUtc)
	if err != nil {
		return errors.Wrap(err, "unable to start mining")
	}

	return errors.Wrap(e.notifyStartMining(ctx, userID, nowUtc), "failed to notify that the user started mining")
}

func (e *economy) notifyStartMining(ctx context.Context, userID UserID, startedAt *time.Time) error {
	m := MiningStarted{
		LastMiningStartedAt: startedAt,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(err, "[start-mining] failed to marshal %#v", m)
	}

	responder := make(chan error, 1)
	e.mb.SendMessage(ctx, &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     userID,
		Topic:   cfg.MessageBroker.Topics[0].Name,
		Value:   b,
	}, responder)

	return errors.Wrapf(<-responder, "[start-mining] failed to send message to broker")
}

func (e *economy) startMining(userID string, startTime *time.Time) error {
	params := map[string]interface{}{
		"userId": userID,
		"now":    startTime,
	}

	sql := fmt.Sprintf(`UPDATE %[1]v SET last_mining_started_at = :now, updated_at = :now WHERE user_id = :userId`, userEconomySpace())

	if err := storage.CheckSQLDMLErr(e.db.PrepareExecute(sql, params)); err != nil {
		return errors.Wrapf(err, "failed set last_mining_started_at for userID:%v", userID)
	}

	return nil
}

func (e *economy) isMiningInProgress(userID UserID) (bool, error) {
	params := map[string]interface{}{
		"userId": userID,
	}

	sql := fmt.Sprintf(`SELECT last_mining_started_at FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1" WHERE user_id = :userId`, userEconomySpace())

	var res []*MiningStarted
	if err := e.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return false, errors.Wrapf(err, "failed to get last_mining_started_at for userID:%v", userID)
	}

	if len(res) == 0 {
		return false, errors.Wrapf(storage.ErrNotFound, "unable to find record for UserID:%v", userID)
	}

	inProgress := miningDuration > time.Now().Sub(*res[0].LastMiningStartedAt.Time)

	return inProgress, nil
}

func userEconomySpace() string {
	return "USER_ECONOMY"
}
