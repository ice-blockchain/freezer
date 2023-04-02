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
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) initializeBlockchainBalanceSynchronizationWorker(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}

	return errors.Wrapf(retry(ctx, func() error {
		if err := r.initializeWorker(ctx, "blockchain_balance_synchronization_worker", usr.ID, usr.HashCode); err != nil {
			if errors.Is(err, storage.ErrRelationNotFound) {
				return err
			}

			return errors.Wrapf(backoff.Permanent(err),
				"failed to initializeBlockchainBalanceSynchronizationWorker for userID:%v", usr.ID)
		}

		return nil
	}), "permanently failed to initializeBlockchainBalanceSynchronizationWorker for userID:%v", usr.ID)
}

func (s *blockchainBalanceSynchronizationTriggerStreamSource) start(ctx context.Context) {
	log.Info("blockchainBalanceSynchronizationTriggerStreamSource started")
	defer log.Info("blockchainBalanceSynchronizationTriggerStreamSource stopped")
	workerIndexes := make([]int16, s.cfg.WorkerCount) //nolint:makezero // Intended.
	for i := 0; i < int(s.cfg.WorkerCount); i++ {
		workerIndexes[i] = int16(i)
	}
	for ctx.Err() == nil {
		stdlibtime.Sleep(blockchainBalanceSynchronizationSeedingStreamEmitFrequency)
		before := time.Now()
		log.Error(errors.Wrap(executeBatchConcurrently(ctx, s.process, workerIndexes), "failed to executeBatchConcurrently[blockchainBalanceSynchronizationTriggerStreamSource.process]")) //nolint:lll // .
		log.Info(fmt.Sprintf("blockchainBalanceSynchronizationTriggerStreamSource.process took: %v", stdlibtime.Since(*before.Time)))
	}
}

func (s *blockchainBalanceSynchronizationTriggerStreamSource) process(ignoredCtx context.Context, workerIndex int16) (err error) {
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
	ctx context.Context, workerIndex int16,
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
		TotalNoPreStakingBonusBalanceAmount    *coin.ICEFlake
		MiningBlockchainAccountAddress, UserID string
		PreStakingAllocation, PreStakingBonus  uint64
	}
)

func (s *blockchainBalanceSynchronizationTriggerStreamSource) getLatestBalancesNewBatch( //nolint:funlen // .
	ctx context.Context, workerIndex int16,
) ([]*latestBalanceSQLRow, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`
SELECT coalesce(coalesce(x.amount,b.amount),'0') AS total_no_pre_staking_bonus_balance_amount,
	   coalesce(x.mining_blockchain_account_address,'') AS mining_blockchain_account_address,
	   x.user_id,
	   coalesce(x.pre_staking_allocation,0) AS pre_staking_allocation,
	   coalesce(st_b.bonus,0) AS pre_staking_bonus
FROM (SELECT MAX(st.years) AS pre_staking_years,
		     MAX(st.allocation) AS pre_staking_allocation,
		     MAX(b.updated_at),
		     b.amount AS amount,
			 x.mining_blockchain_account_address,
			 x.user_id
	  FROM ( SELECT user_id,
					mining_blockchain_account_address
			 FROM blockchain_balance_synchronization_worker
			 WHERE worker_index = $1
			 ORDER BY last_iteration_finished_at
			 LIMIT $2 ) x
		 LEFT JOIN pre_stakings st
			    ON st.worker_index = $1
		       AND st.user_id = x.user_id
		 LEFT JOIN balances_worker b	
			    ON b.worker_index = $1
			   AND b.user_id = x.user_id
			   AND b.negative = FALSE
			   AND b.type = %[1]v
			   AND b.type_detail = ANY($3)
	  GROUP BY x.user_id, x.mining_blockchain_account_address, b.amount
	 ) x
   LEFT JOIN pre_staking_bonuses st_b
		  ON st_b.years = x.pre_staking_years
   LEFT JOIN balance_recalculation_worker not_started_yet_bal_worker
		  ON not_started_yet_bal_worker.worker_index = $1
		 AND not_started_yet_bal_worker.user_id = x.user_id
         AND (not_started_yet_bal_worker.last_iteration_finished_at IS NULL OR not_started_yet_bal_worker.last_mining_ended_at IS NULL)
   LEFT JOIN balances_worker b
		  ON b.worker_index = $1
		 AND b.user_id = not_started_yet_bal_worker.user_id
	     AND b.negative = FALSE
	     AND b.type = %[2]v
	     AND b.type_detail = '%[3]v_%[4]v'`, totalNoPreStakingBonusBalanceType, pendingXBalanceType, rootBalanceTypeDetail, registrationICEBonusEventID)
	var (
		now         = *time.Now().Time
		limit       = maxICEBlockchainConcurrentOperations / int(s.cfg.WorkerCount)
		typeDetails = make([]string, 0, 1+1)
	)
	for i := stdlibtime.Duration(0); i <= 1; i++ {
		dateFormat := now.Add(-1 * i * s.cfg.GlobalAggregationInterval.Child).Format(s.cfg.globalAggregationIntervalChildDateFormat())
		typeDetails = append(typeDetails, fmt.Sprintf("@%v", dateFormat))
	}
	args := append(make([]any, 0, 1+1+1), workerIndex, limit, typeDetails)
	res, err := storage.Select[latestBalanceSQLRow](ctx, s.db, sql, args...)

	return res, errors.Wrapf(err,
		"failed to select a batch of latest information about latest calculating balances for workerIndex:%v,typeDetails:%#v", workerIndex, typeDetails)
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
	const columnNumber = 6
	args := make([]any, 0, len(bs)*columnNumber)
	blockchainMessages := make([]*blockchainMessage, 0, len(bs))
	for ix, bal := range bs {
		if bal.Standard.IsNil() {
			bal.Standard = coin.ZeroICEFlakes()
		}
		if bal.PreStaking.IsNil() {
			bal.PreStaking = coin.ZeroICEFlakes()
		}
		total := coin.New(bal.Standard.Add(bal.PreStaking))
		totalAmount, err := total.Amount.Uint.Marshal()
		log.Panic(err) //nolint:revive // Intended.
		args = append(args, string(totalAmount), total.AmountWord0, total.AmountWord1, total.AmountWord2, total.AmountWord3, bal.UserID)
		values = append(values, fmt.Sprintf("($%[1]v,$%[2]v,$%[3]v,$%[4]v,$%[5]v,$%[6]v)",
			ix*columnNumber+1, ix*columnNumber+2, ix*columnNumber+3, ix*columnNumber+4, ix*columnNumber+5, ix*columnNumber+columnNumber))
		if bal.miningBlockchainAccountAddress != "" {
			blockchainMessages = append(blockchainMessages, &blockchainMessage{
				AccountAddress:     bal.miningBlockchainAccountAddress,
				ICEFlake:           bal.Standard.String(),
				PreStakingICEFlake: bal.PreStaking.String(),
			})
		}
	}
	sql := fmt.Sprintf(`INSERT INTO balances (amount,amount_w0,amount_w1,amount_w2,amount_w3,user_id) 
											 VALUES %v 
						ON CONFLICT (user_id) 
							 DO UPDATE
								   SET amount = EXCLUDED.amount,
									   amount_w0 = EXCLUDED.amount_w0,
									   amount_w1 = EXCLUDED.amount_w1,
									   amount_w2 = EXCLUDED.amount_w2,
									   amount_w3 = EXCLUDED.amount_w3
							 WHERE balances.amount != EXCLUDED.amount`, strings.Join(values, ","))
	if _, err := storage.Exec(ctx, s.db, sql, args...); err != nil {
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
	ctx context.Context, workerIndex int16, rows []*Balances[coin.ICEFlake],
) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	userIDs := make([]string, 0, len(rows))
	for i := range rows {
		userIDs = append(userIDs, rows[i].UserID)
	}
	const table = "blockchain_balance_synchronization_worker"
	params := make(map[string]any, 1)
	params["last_iteration_finished_at"] = *time.Now().Time
	err := s.updateWorkerFields(ctx, workerIndex, table, params, userIDs...)

	return errors.Wrapf(err, "failed to updateWorkerTimeField for workerIndex:%v,table:%q,params:%#v,userIDs:%#v", workerIndex, table, params, userIDs)
}

func (r *repository) updateBlockchainBalanceSynchronizationWorkerBlockchainAccountAddress(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil || usr.MiningBlockchainAccountAddress == "" {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	const table = "blockchain_balance_synchronization_worker"
	workerIndex := int16(usr.HashCode % uint64(r.cfg.WorkerCount))
	params := make(map[string]any, 1)
	params["mining_blockchain_account_address"] = usr.MiningBlockchainAccountAddress
	err := r.updateWorkerFields(ctx, workerIndex, table, params, usr.ID)

	return errors.Wrapf(err, "failed to updateWorkerFields for workerIndex:%v,table:%q,params:%#v,userIDs:%#v", workerIndex, table, params, usr.ID)
}
