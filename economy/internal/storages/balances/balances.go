// SPDX-License-Identifier: BUSL-1.1

package balances

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	tm "time"

	"cosmossdk.io/math"
	"github.com/framey-io/go-tarantool"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func New(db tarantool.Connector, mb messagebroker.Client) messagebroker.Processor {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	return &balanceDistributedBatchProcessingStreamSource{db: db, mb: mb}
}

func (b *balanceDistributedBatchProcessingStreamSource) Process(ctx context.Context, m *messagebroker.Message) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}

	economies, err := b.getUserEconomies(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get user economies")
	}

	wg := new(sync.WaitGroup)
	for _, u := range economies {
		wg.Add(1)
		go b.updateBalances(ctx, u, wg)
	}
	wg.Wait()

	return nil
}

//nolint:funlen // Because this is SQL, no sense to split it to different functions. It will reduce readability.
func (b *balanceDistributedBatchProcessingStreamSource) getUserEconomies(ctx context.Context) ([]*userEconomy, error) {
	params := map[string]interface{}{
		"chunks":             cfg.MessageBroker.Topics[3].Partitions,
		"chunkIndex":         cfg.MessageBroker.Topics[3].Partition,
		"inactivityDeadline": (tm.Duration(cfg.InactivityHoursDeadline) * tm.Hour).Nanoseconds(),
		"now":                time.Now(),
	}
	sql := `SELECT
			ue.user_id,
			(SELECT base_hourly_mining_rate
				FROM adoption
				WHERE active = true
			) AS base_hourly_mining_rate,
			(SELECT GROUP_CONCAT(b.type || ';' || b.amount || ';' || CAST(b.updated_at AS string))
				FROM balances b
				LEFT JOIN user_economy bue
					ON bue.user_id = SUBSTR(b.type, POSITION('~', b.type) + 1)
				WHERE b.user_id = ue.user_id AND (:now - ue.last_mining_started_at < :inactivityDeadline OR ue.last_mining_started_at IS NULL) 
			) AS balances,
			sb.percentage AS bonus,
			s.percentage AS allocation,
			(SELECT count(1)
				FROM balances t0 
				JOIN user_economy u
					ON t0.user_id = ue.user_id AND t0.type = 't0_referral_standard_earnings~' || u.user_id
					AND :now - u.last_mining_started_at < :inactivityDeadline
			) AS t0_referrals,
			(SELECT count(1)
				FROM balances t1
				JOIN user_economy u
					ON t1.user_id = ue.user_id AND t1.type = 't1_referral_standard_earnings~' || u.user_id
					AND :now - u.last_mining_started_at < :inactivityDeadline
			) AS t1_referrals,
			(SELECT count(1)
				FROM balances t2
				JOIN user_economy u
					ON t2.user_id = ue.user_id AND t2.type = 't2_referral_standard_earnings~' || u.user_id
					AND :now - u.last_mining_started_at < :inactivityDeadline
			) AS t2_referrals
		FROM user_economy ue
			LEFT JOIN staking s
				ON s.user_id = ue.user_id
			LEFT JOIN staking_bonus sb
				ON sb.years = s.years
		WHERE ue.hash_code % :chunks == :chunkIndex AND :now - ue.last_mining_started_at < :inactivityDeadline
		LIMIT 1000`

	var res []*userEconomy
	err := b.db.PrepareExecuteTyped(sql, params, &res)
	if err != nil {
		return nil, errors.Wrapf(err, "db execute getting user economies for (chunk:%v, chunkIndex:%v) failed",
			cfg.MessageBroker.Topics[3].Partitions, cfg.MessageBroker.Topics[3].Partition)
	}

	return res, nil
}

//nolint:funlen // All calculation in one place to be more clear.
func (b *balanceDistributedBatchProcessingStreamSource) updateBalances(ctx context.Context, u *userEconomy, wg *sync.WaitGroup) {
	balances := u.parseUserBalances()
	totals := initializeTotals(balances)

	rateMultiplier := u.T0Referrals*cfg.Rates.Tier0 + u.T1Referrals*cfg.Rates.Tier1 + u.T2Referrals*cfg.Rates.Tier2 + percentage100
	hourlyMiningRate := u.BaseHourlyMiningRate.MulUint64(rateMultiplier).QuoUint64(percentage100)
	normalHourlyMiningRate := math.NewUint(percentage100 - u.Allocation).Mul(hourlyMiningRate).QuoUint64(percentage100)
	stakedHourlyMiningRate := math.NewUint(u.Bonus).Mul(hourlyMiningRate).Mul(math.NewUint(u.Allocation)).QuoUint64(stakedHourlyMiningRateDivider)

	standard := normalHourlyMiningRate.MulUint64(uint64(time.Now().UnixNano()) - balances["standard"].UpdatedAt).QuoUint64(generalBalanceDivider)
	staking := stakedHourlyMiningRate.MulUint64(uint64(time.Now().UnixNano()) - balances["staking"].UpdatedAt).QuoUint64(generalBalanceDivider)

	totals["standard"] = balances["standard"].Amount.Add(standard)
	totals["staking"] = balances["staking"].Amount.Add(staking)
	totals["total"] = standard.Add(staking)

	standardGeneral := math.NewUint(percentage100 - u.Allocation).Mul(u.BaseHourlyMiningRate.Uint)
	stakingGeneral := math.NewUint(u.Allocation * u.Bonus).Mul(u.BaseHourlyMiningRate.Uint)

	for balanceType, value := range balances {
		var result, earnings math.Uint
		if strings.Contains(balanceType, "~") {
			elapsedNanoseconds := uint64(time.Now().UnixNano()) - value.UpdatedAt
			parts := strings.Split(balanceType, "~")

			switch parts[0] {
			case "t0_referral_standard_earnings":
				earnings = standardGeneral.MulUint64(u.T0Referrals).MulUint64(elapsedNanoseconds).QuoUint64(t0StandardDivider)
				result = value.Amount.Add(earnings)
			case "t1_referral_standard_earnings":
				earnings = standardGeneral.MulUint64(u.T1Referrals).MulUint64(elapsedNanoseconds).QuoUint64(t1StandardDivider)
				result = value.Amount.Add(earnings)
			case "t2_referral_standard_earnings":
				earnings = standardGeneral.MulUint64(u.T2Referrals).MulUint64(elapsedNanoseconds).QuoUint64(t2StandardDivider)
				result = value.Amount.Add(earnings)
			case "t0_referral_staking_earnings":
				earnings = stakingGeneral.MulUint64(u.T0Referrals).MulUint64(elapsedNanoseconds).QuoUint64(t0StakingDivider)
				result = value.Amount.Add(earnings)
			case "t1_referral_staking_earnings":
				earnings = stakingGeneral.MulUint64(u.T1Referrals).MulUint64(elapsedNanoseconds).QuoUint64(t1StakingDivider)
				result = value.Amount.Add(earnings)
			case "t2_referral_staking_earnings":
				earnings = stakingGeneral.MulUint64(u.T2Referrals).MulUint64(elapsedNanoseconds).QuoUint64(t2StakingDivider)
				result = value.Amount.Add(earnings)
			default:
				continue
			}

			totals[parts[0]] = totals[parts[0]].Add(earnings)
			log.Error(errors.Wrapf(b.updateBalance(ctx, u.UserID, balanceType, result),
				"can't update %v balance for userID:%v by amount:%v", balanceType, u.UserID, result))
		}
	}

	for balanceType, amount := range totals {
		log.Error(errors.Wrapf(b.updateBalance(ctx, u.UserID, balanceType, amount),
			"can't update %v balance for userID:%v by amount:%v", balanceType, u.UserID, amount))
	}

	log.Error(errors.Wrapf(b.produceBalancesEventMessage(totals), "can't call produceBalancesEventMessage for userID:%v", u.UserID))
	wg.Done()
}

func (u *userEconomy) parseUserBalances() map[string]balance {
	balances := strings.Split(u.Balances, ",")
	res := make(map[string]balance, len(balances))

	for _, val := range balances {
		parts := strings.Split(val, ";")
		updatedAt, err := strconv.ParseUint(parts[2], base10, bitSize64)
		log.Panic(errors.Wrapf(err, "can't parse uint64 value for balance updatedAt:%v", parts[2]))

		res[parts[0]] = balance{Amount: coin.UnsafeNewAmount(parts[1]), UpdatedAt: updatedAt}
	}

	return res
}

func (b *balanceDistributedBatchProcessingStreamSource) updateBalance(ctx context.Context, userID UserID, balanceType BalanceType, balance math.Uint) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	if !balance.IsZero() {
		var c coin.Coin
		c.SetAmount(&coin.ICEFlake{Uint: balance})
		params := map[string]interface{}{
			"amount":    c.Amount,
			"amountW0":  c.AmountWord0,
			"amountW1":  c.AmountWord1,
			"amountW2":  c.AmountWord2,
			"amountW3":  c.AmountWord3,
			"updatedAt": time.Now(),
			"userId":    userID,
			"type":      balanceType,
		}

		sql := `UPDATE balances SET
						amount = :amount,
						amount_w0 = :amountW0, amount_w1 = :amountW1, amount_w2 = :amountW2, amount_w3 = :amountW3,
						updated_at = :updatedAt
					WHERE user_id = :userId AND type = :type`

		query, err := b.db.PrepareExecute(sql, params)
		if err = storage.CheckSQLDMLErr(query, err); err != nil {
			return errors.Wrapf(err, "failed to update balances with userID:%v and type:%v", userID, balanceType)
		}
	}

	return nil
}

func (b *balanceDistributedBatchProcessingStreamSource) produceBalancesEventMessage(balances map[string]math.Uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), userBalanceMessageDeadline)
	defer cancel()

	message, err := json.Marshal(balances)
	if err != nil {
		return errors.Wrapf(err, "can't marshal %#v balance message", balances)
	}

	responder := make(chan error, 1)
	m := &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     uuid.NewString(),
		Topic:   cfg.MessageBroker.Topics[5].Name,
		Value:   message,
	}

	defer close(responder)
	b.mb.SendMessage(ctx, m, responder)

	return errors.Wrapf(<-responder, "failed to send user balances message: %#v", m)
}

func initializeTotals(balances map[string]balance) map[string]math.Uint {
	return map[string]math.Uint{
		"standard":                      balances["standard"].Amount.Uint,
		"staking":                       balances["staking"].Amount.Uint,
		"total":                         balances["total"].Amount.Uint,
		"t0_referral_standard_earnings": balances["t0_referral_standard_earnings"].Amount.Uint,
		"t1_referral_standard_earnings": balances["t1_referral_standard_earnings"].Amount.Uint,
		"t2_referral_standard_earnings": balances["t2_referral_standard_earnings"].Amount.Uint,
		"t0_referral_staking_earnings":  balances["t0_referral_staking_earnings"].Amount.Uint,
		"t1_referral_staking_earnings":  balances["t1_referral_staking_earnings"].Amount.Uint,
		"t2_referral_staking_earnings":  balances["t2_referral_staking_earnings"].Amount.Uint,
	}
}
