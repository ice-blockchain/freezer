// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/cenkalti/backoff/v4"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) initializeBlockchainBalanceSynchronizationWorker(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	workerIndex := usr.HashCode % r.cfg.WorkerCount
	err := retry(ctx, func() error {
		if err := r.initializeWorker(ctx, "blockchain_balance_synchronization_worker_", usr.ID, workerIndex); err != nil {
			if errors.Is(err, storage.ErrRelationNotFound) {
				return err
			}

			return errors.Wrapf(backoff.Permanent(err),
				"failed to initializeBlockchainBalanceSynchronizationWorker for userID:%v,workerIndex:%v", usr.ID, workerIndex)
		}

		return nil
	})

	return errors.Wrapf(err, "permanently failed to initializeBlockchainBalanceSynchronizationWorker for userID:%v,workerIndex:%v", usr.ID, workerIndex)
}

func (s *blockchainBalanceSynchronizationTriggerStreamSource) start(ctx context.Context) {
	log.Info("blockchainBalanceSynchronizationTriggerStreamSource started")
	defer log.Info("blockchainBalanceSynchronizationTriggerStreamSource stopped")
	workerIndexes := make([]uint64, s.cfg.WorkerCount) //nolint:makezero // Intended.
	for i := 0; i < int(s.cfg.WorkerCount); i++ {
		workerIndexes[i] = uint64(i)
	}
	for ctx.Err() == nil {
		stdlibtime.Sleep(blockchainBalanceSynchronizationSeedingStreamEmitFrequency)
		before := time.Now()
		log.Error(errors.Wrap(executeBatchConcurrently(ctx, s.process, workerIndexes), "failed to executeBatchConcurrently[blockchainBalanceSynchronizationTriggerStreamSource.process]")) //nolint:lll // .
		log.Info(fmt.Sprintf("blockchainBalanceSynchronizationTriggerStreamSource.process took: %v", stdlibtime.Since(*before.Time)))
	}
}

func (s *blockchainBalanceSynchronizationTriggerStreamSource) process(ignoredCtx context.Context, workerIndex uint64) (err error) {
	if ignoredCtx.Err() != nil {
		return errors.Wrap(ignoredCtx.Err(), "unexpected deadline while processing message")
	}
	const deadline = 5 * stdlibtime.Minute
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	rows, err := s.getLatestBalances(ctx, workerIndex) //nolint:contextcheck // We use context with longer deadline.
	if err != nil || len(rows) == 0 {
		return errors.Wrapf(err, "failed to getLatestBalances for workerIndex:%v", workerIndex)
	}
	if err = s.updateBalances(ctx, rows); err != nil { //nolint:contextcheck // Intended.
		return errors.Wrapf(err, "failed to updateBalances:%#v", rows)
	}
	if err = executeBatchConcurrently(ctx, s.sendBalancesMessage, rows); err != nil { //nolint:contextcheck // We use context with longer deadline.
		return errors.Wrapf(err, "failed to sendBalancesMessages for:%#v", rows)
	}

	return errors.Wrapf(s.updateLastIterationFinishedAt(ctx, workerIndex, rows), //nolint:contextcheck // We use context with longer deadline.
		"failed to updateLastIterationFinishedAt for workerIndex:%v,rows:%#v", workerIndex, rows)
}

func (s *blockchainBalanceSynchronizationTriggerStreamSource) getLatestBalances( //nolint:funlen // .
	ctx context.Context, workerIndex uint64,
) ([]*Balances[coin.ICEFlake], error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	batch, err := s.getLatestBalancesNewBatch(ctx, workerIndex)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getLatestBalancesNewBatch for workerIndex:%v", workerIndex)
	}
	if len(batch) == 0 {
		return nil, nil
	}
	res := make([]*Balances[coin.ICEFlake], 0, len(batch))
	for _, row := range batch {
		var standard, preStaking *coin.ICEFlake
		switch row.PreStakingAllocation {
		case 0:
			standard = row.TotalNoPreStakingBonusBalanceAmount
		case percentage100:
			preStaking = row.TotalNoPreStakingBonusBalanceAmount.
				MultiplyUint64(row.PreStakingBonus + percentage100).
				DivideUint64(percentage100)
		default:
			standard = row.TotalNoPreStakingBonusBalanceAmount.
				MultiplyUint64(percentage100 - row.PreStakingAllocation).
				DivideUint64(percentage100)
			preStaking = row.TotalNoPreStakingBonusBalanceAmount.
				MultiplyUint64(row.PreStakingAllocation * (row.PreStakingBonus + percentage100)).
				DivideUint64(percentage100 * percentage100)
		}
		res = append(res, &Balances[coin.ICEFlake]{
			Standard:                       standard,
			PreStaking:                     preStaking,
			UserID:                         row.UserID,
			miningBlockchainAccountAddress: row.MiningBlockchainAccountAddress,
		})
	}

	return res, nil
}

type (
	latestBalanceSQLRow struct {
		_msgpack                               struct{} `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
		TotalNoPreStakingBonusBalanceAmount    *coin.ICEFlake
		MiningBlockchainAccountAddress, UserID string
		PreStakingAllocation, PreStakingBonus  uint64
	}
)

func (s *blockchainBalanceSynchronizationTriggerStreamSource) getLatestBalancesNewBatch( //nolint:funlen // .
	ctx context.Context, workerIndex uint64,
) ([]*latestBalanceSQLRow, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	var (
		now         = *time.Now().Time
		limit       = maxICEBlockchainConcurrentOperations / s.cfg.WorkerCount
		typeDetails = make([]string, 0, 1+1)
		params      = make(map[string]any, 1+1)
	)
	for i := stdlibtime.Duration(0); i <= 1; i++ {
		dateFormat := now.Add(-1 * i * s.cfg.GlobalAggregationInterval.Child).Format(s.cfg.globalAggregationIntervalChildDateFormat())
		params[fmt.Sprintf("type_detail%v", i)] = fmt.Sprintf("@%v", dateFormat)
		typeDetails = append(typeDetails, fmt.Sprintf(":type_detail%v", i))
	}
	sql := fmt.Sprintf(`
SELECT IFNULL(IFNULL(x.amount,b.amount),'0'),
	   x.mining_blockchain_account_address,
	   x.user_id,
	   x.pre_staking_allocation,
	   st_b.bonus AS pre_staking_bonus
FROM (SELECT MAX(st.years) AS pre_staking_years,
		     MAX(st.allocation) AS pre_staking_allocation,
		     MAX(b.updated_at),
		     b.amount AS amount,
			 x.mining_blockchain_account_address,
			 x.user_id
	  FROM ( SELECT user_id,
					mining_blockchain_account_address
			 FROM blockchain_balance_synchronization_worker_%[2]v
			 ORDER BY last_iteration_finished_at
			 LIMIT %[1]v ) x
		 LEFT JOIN pre_stakings_%[2]v st
		        ON st.user_id = x.user_id
		 LEFT JOIN balances_%[2]v b	
			    ON b.user_id = x.user_id
			   AND b.negative = FALSE
			   AND b.type = %[3]v
			   AND b.type_detail IN (%[5]v)
	  GROUP BY x.user_id
	 ) x
   LEFT JOIN pre_staking_bonuses st_b
		  ON st_b.years = x.pre_staking_years
   LEFT JOIN balance_recalculation_worker_%[2]v not_started_yet_bal_worker
		  ON not_started_yet_bal_worker.user_id = x.user_id
         AND (not_started_yet_bal_worker.last_iteration_finished_at IS NULL OR not_started_yet_bal_worker.last_mining_ended_at IS NULL)
   LEFT JOIN balances_%[2]v b
		  ON b.user_id = not_started_yet_bal_worker.user_id
	     AND b.negative = FALSE
	     AND b.type = %[4]v
	     AND b.type_detail = ''`, limit, workerIndex, totalNoPreStakingBonusBalanceType, pendingXBalanceType, strings.Join(typeDetails, ","))
	res := make([]*latestBalanceSQLRow, 0, limit)
	if err := s.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return nil, errors.Wrapf(err,
			"failed to select a batch of latest information about latest calculating balances for workerIndex:%v,params:%#v", workerIndex, params)
	}

	return res, nil
}

func (s *blockchainBalanceSynchronizationTriggerStreamSource) updateBalances( //nolint:funlen // Mostly mappings.
	ctx context.Context, bs []*Balances[coin.ICEFlake],
) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	//nolint:godox // .
	type blockchainMessage struct { // TODO: delete this and use the actual one.
		AccountAddress     string
		ICEFlake           string
		PreStakingICEFlake string
	}
	values := make([]string, 0, len(bs))
	blockchainMessages := make([]*blockchainMessage, 0, len(bs))
	for _, bal := range bs {
		if bal.Standard.IsNil() {
			bal.Standard = coin.ZeroICEFlakes()
		}
		if bal.PreStaking.IsNil() {
			bal.PreStaking = coin.ZeroICEFlakes()
		}
		total := coin.New(bal.Standard.Add(bal.PreStaking))
		totalAmount, err := total.Amount.Uint.Marshal()
		log.Panic(err) //nolint:revive // Intended.
		values = append(values, fmt.Sprintf("('%[1]v',%[2]v,%[3]v,%[4]v,%[5]v,'%[6]v')",
			string(totalAmount), total.AmountWord0, total.AmountWord1, total.AmountWord2, total.AmountWord3, bal.UserID))
		if bal.miningBlockchainAccountAddress != "" {
			blockchainMessages = append(blockchainMessages, &blockchainMessage{
				AccountAddress:     bal.miningBlockchainAccountAddress,
				ICEFlake:           bal.Standard.String(),
				PreStakingICEFlake: bal.PreStaking.String(),
			})
		}
	}
	sql := fmt.Sprintf(`REPLACE INTO balances (amount,amount_w0,amount_w1,amount_w2,amount_w3,user_id) VALUES %v`, strings.Join(values, ","))
	if _, err := storage.CheckSQLDMLResponse(s.db.Execute(sql)); err != nil {
		return errors.Wrapf(err, "failed to replace into balances, values:%#v", values)
	}
	if len(blockchainMessages) != 0 { //nolint:revive,staticcheck // .
		//nolint:godox // .
		// TODO use blockchainMessages.
	}

	return nil
}

func (s *blockchainBalanceSynchronizationTriggerStreamSource) sendBalancesMessage(ctx context.Context, bs *Balances[coin.ICEFlake]) error {
	valueBytes, err := json.MarshalContext(ctx, bs)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %#v", bs)
	}
	msg := &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     bs.UserID,
		Topic:   s.cfg.MessageBroker.Topics[4].Name,
		Value:   valueBytes,
	}
	responder := make(chan error, 1)
	defer close(responder)
	s.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send `%v` message to broker", msg.Topic)
}

func (s *blockchainBalanceSynchronizationTriggerStreamSource) updateLastIterationFinishedAt(
	ctx context.Context, workerIndex uint64, rows []*Balances[coin.ICEFlake],
) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	userIDs := make([]string, 0, len(rows))
	for i := range rows {
		userIDs = append(userIDs, rows[i].UserID)
	}
	const table = "blockchain_balance_synchronization_worker_"
	params := make(map[string]any, 1)
	params["last_iteration_finished_at"] = time.Now()
	err := s.updateWorkerFields(ctx, workerIndex, table, params, userIDs...)

	return errors.Wrapf(err, "failed to updateWorkerTimeField for workerIndex:%v,table:%q,params:%#v,userIDs:%#v", workerIndex, table, params, userIDs)
}

func (r *repository) updateBlockchainBalanceSynchronizationWorkerBlockchainAccountAddress(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	const table = "blockchain_balance_synchronization_worker_"
	workerIndex := usr.HashCode % r.cfg.WorkerCount
	params := make(map[string]any, 1)
	params["mining_blockchain_account_address"] = usr.MiningBlockchainAccountAddress
	err := r.updateWorkerFields(ctx, workerIndex, table, params, usr.ID)

	return errors.Wrapf(err, "failed to updateWorkerFields for workerIndex:%v,table:%q,params:%#v,userIDs:%#v", workerIndex, table, params, usr.ID)
}
