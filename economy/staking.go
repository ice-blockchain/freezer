// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"

	"cosmossdk.io/math"
	"github.com/framey-io/go-tarantool"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/time"
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

	nowUtc := time.Now()

	err = e.enableStaking(userID, staking, nowUtc)
	if err != nil {
		return errors.Wrap(err, "unable to enable staking")
	}

	return errors.Wrap(e.notifyStartStaking(ctx, userID, staking, nowUtc), "failed to notify that the user enable staking")
}

func (e *economy) isStakingEnabled(userID UserID) (bool, error) {
	params := map[string]interface{}{
		"userID": userID,
		"type":   balanceTypeStaking,
	}

	sql := `SELECT s.years > 0 AND s.percentage > 0 AND b.balance IS NOT NULL
		FROM staking s INDEXED BY "pk_unnamed_STAKING_1" 
			JOIN balances b
				ON b.user_id = s.user_id and b.type = 'staking'
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

func (e *economy) enableStaking(userID string, staking Staking, now *time.Time) error {
	b, err := e.getUserBalance(userID)
	if err != nil {
		return errors.Wrapf(err, "failed to get user balance for userID %v", userID)
	}

	allocation := math.NewUint(staking.Percentage)

	var stakingBalance coin.Coin
	var remainings coin.Coin
	stakingBalance.SetAmount(&coin.ICEFlake{Uint: b.Mul(allocation).QuoUint64(100)})
	remainings.SetAmount(&coin.ICEFlake{Uint: b.Sub(stakingBalance.Amount.Uint)})

	var errs error
	if err := e.updateBalance(userID, &stakingBalance, balanceTypeStaking); err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "failed to update staking balance for userID %v", userID))
	}
	if err := e.updateBalance(userID, &remainings, balanceTypeStandard); err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "failed to update standard balance for userID %v", userID))
	}
	if err := e.updateStaking(userID, staking); err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "failed to update staking information for userID %v", userID))
	}

	return errors.Wrapf(errs, "unable to enable staking")
}

func (e *economy) getUserBalance(userID UserID) (*coin.ICEFlake, error) {
	params := map[string]interface{}{
		"userId": userID,
		"type":   balanceTypeStandard,
	}

	sql := `SELECT amount
		FROM balances INDEXED BY "pk_unnamed_BALANCES_1" 
		WHERE user_id = :userId AND type = :type`

	var res []userBalance
	if err := e.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return nil, errors.Wrapf(err, "failed to get user_economy record with userID %v", userID)
	}

	if len(res) == 0 {
		return nil, ErrNotFound
	}

	return res[0].Balance, nil
}

func (e *economy) updateStaking(userID UserID, s Staking) error {
	space := "STAKING"
	index := "pk_unnamed_STAKING_1"
	key := tarantool.StringKey{S: userID}

	ops := []tarantool.Op{
		{Op: "=", Field: 2, Arg: s.Percentage},
		{Op: "=", Field: 3, Arg: s.Years},
		{Op: "=", Field: 5, Arg: time.Now()},
	}

	return errors.Wrapf(e.db.UpdateTyped(space, index, key, ops, &[]*staking{}),
		"failed update staking info:%v with staking percentage:%v, staking years:%v, for userID:%v",
		s.Percentage, s.Years, userID)
}

func (e *economy) updateBalance(userID UserID, balance *coin.Coin, balanceType string) error {
	params := map[string]interface{}{
		"amount":    balance.Amount,
		"amountW0":  balance.AmountWord0,
		"amountW1":  balance.AmountWord1,
		"amountW2":  balance.AmountWord2,
		"amountW3":  balance.AmountWord3,
		"userId":    userID,
		"type":      balanceType,
		"updatedAt": time.Now(),
	}

	sql := `UPDATE balances SET
				   amount = :balance,
				   amount_w0 = :amountW0,
				   amount_w1 = :amountW1,
				   amount_w2 = :amountW2,
				   amount_w3 = :amountW3,
				   updated_at = :updatedAt
				WHERE user_id = :userId AND type = :type`

	query, err := e.db.PrepareExecute(sql, params)

	if err = storage.CheckSQLDMLErr(query, err); err != nil {
		return errors.Wrapf(err, "failed to update balances with userID:%v and type:%v", userID, balanceType)
	}

	return nil
}

func (e *economy) notifyStartStaking(ctx context.Context, userID UserID, staking Staking, startedAt *time.Time) error {
	m := StakingEnabled{
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
