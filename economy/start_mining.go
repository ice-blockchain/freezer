// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
)

func (r *userEconomyRepository) StartMining(ctx context.Context, userID UserID) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "start mining failed because context failed")
	}

	miningInProgress, err := r.isMiningInProgress(userID)
	if err != nil {
		return errors.Wrap(err, "unable to check is mining in porgress")
	}

	if miningInProgress {
		return ErrMiningInProgress
	}

	miningStarted, err := r.startMining(userID)
	if err != nil {
		return errors.Wrap(err, "unable to start mining")
	}

	if err = r.notifyStartMining(ctx, userID, miningStarted); err != nil {
		return errors.Wrap(err, "unable to notify start mining")
	}

	return nil
}

func (r *userEconomyRepository) notifyStartMining(ctx context.Context, userID UserID, startedAt uint64) error {
	m := miningStarted{
		TS:     startedAt,
		UserID: userID,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(err, "[start-mining] failed to marshal %#v", m)
	}

	responder := make(chan error, 1)
	r.mb.SendMessage(ctx, &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     userID,
		Topic:   cfg.MessageBroker.Topics[0].Name,
		Value:   b,
	}, responder)

	return errors.Wrapf(<-responder, "[start-mining] failed to send message to broker")
}

func (r *userEconomyRepository) startMining(userID string) (uint64, error) {
	nowNano := uint64(time.Now().UnixNano())

	params := map[string]interface{}{
		"userId":        userID,
		"miningStarted": nowNano,
		"updatedAt":     nowNano,
	}

	sql := fmt.Sprintf(`UPDATE %[1]v SET last_mining_started_at = :miningStarted, updated_at = :updatedAt WHERE user_id = :userId`, userEconomySpace())

	if err := storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)); err != nil {
		return 0, errors.Wrapf(err, "failed set last_mining_started_at for userID:%v", userID)
	}

	return nowNano, nil
}

func (r *userEconomyRepository) isMiningInProgress(userID UserID) (bool, error) {
	params := map[string]interface{}{
		"userId": userID,
	}

	sql := fmt.Sprintf(`SELECT * FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1" WHERE user_id = :userId`, userEconomySpace())

	var res []*userEconomy
	if err := r.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return false, errors.Wrapf(err, "failed to get last_mining_started_at for userID:%v", userID)
	}

	if len(res) == 0 {
		return false, errors.Wrapf(nil, "unable to find record for UserID:%v", userID)
	}

	miningStared := time.Unix(0, int64(res[0].LastMiningStartedAt))
	inProgress := miningDuration > time.Now().Sub(miningStared)

	return inProgress, nil
}
