// SPDX-License-Identifier: ice License 1.0

package balancesynchronizer

import (
	"context"
	"sync"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/freezer/model"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
)

func init() {
	appCfg.MustLoadFromKey(parentApplicationYamlKey, &cfg.Config)
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
}

func MustStartSynchronizingBalance(ctx context.Context) {
	bs := &balanceSynchronizer{
		mb: messagebroker.MustConnect(context.Background(), parentApplicationYamlKey),
	}
	defer log.Panic(errors.Wrap(bs.Close(), "failed to stop balanceSynchronizer"))

	wg := new(sync.WaitGroup)
	wg.Add(int(cfg.Workers))
	defer wg.Wait()

	for workerNumber := int64(0); workerNumber < cfg.Workers; workerNumber++ {
		go func(wn int64) {
			defer wg.Done()
			bs.synchronize(ctx, wn)
		}(workerNumber)
	}
}

func (bs *balanceSynchronizer) Close() error {
	return multierror.Append(
		errors.Wrap(bs.mb.Close(), "failed to close mb"),
	).ErrorOrNil()
}

func (bs *balanceSynchronizer) synchronize(ctx context.Context, workerNumber int64) {
	db := storage.MustConnect(context.Background(), parentApplicationYamlKey, 1)
	defer func() {
		if err := recover(); err != nil {
			log.Error(db.Close())
			panic(err)
		}
		log.Error(db.Close())
	}()
	var (
		batchNumber        int64
		iteration          uint64
		workers            = cfg.Workers
		batchSize          = cfg.BatchSize
		userKeys           = make([]string, 0, batchSize)
		userResults        = make([]*user, 0, batchSize)
		msgResponder       = make(chan error, batchSize)
		msgs               = make([]*messagebroker.Message, 0, batchSize)
		errs               = make([]error, 0, batchSize)
		updatedUsers       = make([]redis.Z, 0, batchSize)
		blockchainMessages = make([]*blockchainMessage, 0, batchSize)
	)
	resetVars := func(success bool) {
		if success && len(userKeys) < int(batchSize) {
			batchNumber = 0
			iteration++
		}
		userKeys = userKeys[:0]
		userResults = userResults[:0]
		msgs, errs = msgs[:0], errs[:0]
		updatedUsers = updatedUsers[:0]
		blockchainMessages = blockchainMessages[:0]
	}
	for ctx.Err() == nil {
		/******************************************************************************************************************************************************
			1. Fetching a new batch of users.
		******************************************************************************************************************************************************/
		if len(userKeys) == 0 {
			for ix := batchNumber * batchSize; ix < (batchNumber+1)*batchSize; ix++ {
				userKeys = append(userKeys, model.SerializedUsersKey((workers*ix)+workerNumber))
			}
		}
		reqCtx, reqCancel := context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[user](reqCtx, db, userKeys, &userResults); err != nil {
			log.Error(errors.Wrapf(err, "[balanceSynchronizer] failed to get users for batchNumer:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			2. Processing batch.
		******************************************************************************************************************************************************/

		for _, usr := range userResults {
			updatedUsers = append(updatedUsers, redis.Z{
				Score:  usr.BalanceTotalStandard + usr.BalanceTotalPreStaking,
				Member: model.SerializedUsersKey(usr.ID),
			})
			if msg := shouldSendBalanceUpdatedMessage(ctx, iteration, usr); msg != nil {
				msgs = append(msgs, msg)
			}
			if msg := shouldSynchronizeBlockchainAccount(iteration, usr); msg != nil {
				blockchainMessages = append(blockchainMessages, msg)
			}
		}

		/******************************************************************************************************************************************************
			3. Sending messages to the broker.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		for _, message := range msgs {
			bs.mb.SendMessage(reqCtx, message, msgResponder)
		}
		for len(errs) < cap(errs) || len(msgResponder) > 0 {
			errs = append(errs, <-msgResponder)
		}
		if err := multierror.Append(reqCtx.Err(), errs...).ErrorOrNil(); err != nil {
			log.Error(errors.Wrapf(err, "[balanceSynchronizer] failed to send messages to broker for batchNumer:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			4. Updating user scores in `top_miners` sorted set.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := db.ZAdd(reqCtx, "top_miners", updatedUsers...).Err(); err != nil {
			log.Error(errors.Wrapf(err, "[balanceSynchronizer] failed to ZAdd top_miners for batchNumer:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			5. Updating balances in the blockchain for that batch of users.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := bs.synchronizeBlockchainAccounts(reqCtx, blockchainMessages); err != nil {
			log.Error(errors.Wrapf(err, "[balanceSynchronizer] failed to synchronizeBlockchainAccount for batchNumer:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}

		batchNumber++
		reqCancel()
		resetVars(true)
	}

}

func shouldSendBalanceUpdatedMessage(ctx context.Context, iteration uint64, usr *user) *messagebroker.Message {
	if iteration%(uint64(usr.ID)%10) != 0 {
		return nil
	}
	event := &BalanceUpdated{
		UserID:     usr.UserID,
		Standard:   usr.BalanceTotalStandard,
		PreStaking: usr.BalanceTotalPreStaking,
	}
	valueBytes, err := json.MarshalContext(ctx, event)
	log.Panic(errors.Wrapf(err, "failed to marshal %#v", event))

	return &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     event.UserID,
		Topic:   cfg.MessageBroker.Topics[3].Name,
		Value:   valueBytes,
	}
}
