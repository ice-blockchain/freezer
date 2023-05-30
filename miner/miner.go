// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func init() {
	appCfg.MustLoadFromKey(parentApplicationYamlKey, &cfg.Config)
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
}

func MustStartMining(ctx context.Context) {
	mi := &miner{
		mb: messagebroker.MustConnect(context.Background(), parentApplicationYamlKey),
	}
	defer func() { log.Panic(errors.Wrap(mi.Close(), "failed to stop miner")) }()

	wg := new(sync.WaitGroup)
	wg.Add(int(cfg.Workers))
	defer wg.Wait()

	for workerNumber := int64(0); workerNumber < cfg.Workers; workerNumber++ {
		go func(wn int64) {
			defer wg.Done()
			mi.mine(ctx, wn)
		}(workerNumber)
	}
}

func (m *miner) Close() error {
	return multierror.Append(
		errors.Wrap(m.mb.Close(), "failed to close mb"),
	).ErrorOrNil()
}

func (m *miner) mine(ctx context.Context, workerNumber int64) {
	db := storage.MustConnect(context.Background(), parentApplicationYamlKey, 1)
	defer func() {
		if err := recover(); err != nil {
			log.Error(db.Close())
			panic(err)
		}
		log.Error(db.Close())
	}()
	dwhClient := dwh.MustConnect(context.Background(), applicationYamlKey)
	defer func() {
		if err := recover(); err != nil {
			log.Error(dwhClient.Close())
			panic(err)
		}
		log.Error(dwhClient.Close())
	}()
	var (
		batchNumber                                                int64
		now                                                        = time.Now()
		currentAdoption                                            *tokenomics.Adoption[float64]
		workers                                                    = cfg.Workers
		batchSize                                                  = cfg.BatchSize
		userKeys, userHistoryKeys, referralKeys                    = make([]string, 0, batchSize), make([]string, 0, batchSize), make([]string, 0, 2*batchSize)
		userResults, referralResults                               = make([]*user, 0, batchSize), make([]*referral, 0, 2*batchSize)
		t0Referrals, tMinus1Referrals                              = make(map[int64]*referral, batchSize), make(map[int64]*referral, batchSize)
		t1ReferralsThatStoppedMining, t2ReferralsThatStoppedMining = make(map[int64]uint32, batchSize), make(map[int64]uint32, batchSize)
		referralsThatStoppedMining                                 = make([]*referralThatStoppedMining, 0, batchSize)
		msgResponder                                               = make(chan error, batchSize)
		msgs                                                       = make([]*messagebroker.Message, 0, batchSize)
		errs                                                       = make([]error, 0, batchSize)
		updatedUsers                                               = make([]*UpdatedUser, 0, batchSize)
		histories                                                  = make([]*model.User, 0, batchSize)
		historyColumns, historyInsertMetadata                      = dwh.InsertDDL(int(batchSize))
	)
	resetVars := func(success bool) {
		now = time.Now()
		if success && len(userResults) < int(batchSize) {
			batchNumber = 0
		}
		if batchNumber == 0 || currentAdoption == nil {
			for err := errors.New("init"); ctx.Err() == nil && err != nil; {
				reqCtx, reqCancel := context.WithTimeout(context.Background(), requestDeadline)
				currentAdoption, err = tokenomics.GetCurrentAdoption(reqCtx, db)
				reqCancel()
				log.Error(errors.Wrapf(err, "[miner] failed to GetCurrentAdoption for workerNumber:%v", workerNumber))
			}
		}
		userKeys, userHistoryKeys, referralKeys = userKeys[:0], userHistoryKeys[:0], referralKeys[:0]
		userResults, referralResults = userResults[:0], referralResults[:0]
		msgs, errs = msgs[:0], errs[:0]
		updatedUsers = updatedUsers[:0]
		histories = histories[:0]
		referralsThatStoppedMining = referralsThatStoppedMining[:0]
		for k := range t0Referrals {
			delete(t0Referrals, k)
		}
		for k := range tMinus1Referrals {
			delete(tMinus1Referrals, k)
		}
		for k := range t1ReferralsThatStoppedMining {
			delete(t1ReferralsThatStoppedMining, k)
		}
		for k := range t2ReferralsThatStoppedMining {
			delete(t2ReferralsThatStoppedMining, k)
		}
	}
	resetVars(true)
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
			log.Error(errors.Wrapf(err, "[miner] failed to get users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			2. Fetching T0 & T-1 referrals of the fetched users.
		******************************************************************************************************************************************************/

		for _, usr := range userResults {
			if usr.IDT0 > 0 {
				t0Referrals[usr.IDT0] = nil
			}
			if usr.IDT0 < 0 {
				t0Referrals[-usr.IDT0] = nil
			}
			if usr.IDTMinus1 > 0 {
				tMinus1Referrals[usr.IDTMinus1] = nil
			}
			if usr.IDTMinus1 < 0 {
				tMinus1Referrals[-usr.IDTMinus1] = nil
			}
		}
		for idT0 := range t0Referrals {
			referralKeys = append(referralKeys, model.SerializedUsersKey(idT0))
		}
		for idTMinus1 := range tMinus1Referrals {
			referralKeys = append(referralKeys, model.SerializedUsersKey(idTMinus1))
		}

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[referral](reqCtx, db, referralKeys, &referralResults); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to get referrees for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			3. Mining for the users.
		******************************************************************************************************************************************************/

		for _, ref := range referralResults {
			if _, found := tMinus1Referrals[ref.ID]; found {
				tMinus1Referrals[ref.ID] = ref
			}
			if _, found := t0Referrals[ref.ID]; found {
				t0Referrals[ref.ID] = ref
			}
		}
		for _, usr := range userResults {
			var t0Ref, tMinus1Ref *referral
			if usr.IDT0 > 0 {
				t0Ref = t0Referrals[usr.IDT0]
			}
			if usr.IDT0 < 0 {
				t0Ref = t0Referrals[-usr.IDT0]
			}
			if usr.IDTMinus1 > 0 {
				tMinus1Ref = tMinus1Referrals[usr.IDTMinus1]
			}
			if usr.IDTMinus1 < 0 {
				tMinus1Ref = tMinus1Referrals[-usr.IDTMinus1]
			}
			updatedUser, shouldGenerateHistory := mine(currentAdoption.BaseMiningRate, now, usr, t0Ref, tMinus1Ref)
			if updatedUser != nil {
				updatedUsers = append(updatedUsers, &updatedUser.UpdatedUser)
			}
			if shouldGenerateHistory {
				userHistoryKeys = append(userHistoryKeys, usr.Key())
			}
			if userStoppedMining := didReferralJustStopMining(now, usr, updatedUser); userStoppedMining != nil {
				referralsThatStoppedMining = append(referralsThatStoppedMining, userStoppedMining)
			}
			if dayOffStarted := didANewDayOffJustStart(now, usr); dayOffStarted != nil {
				msgs = append(msgs, dayOffStartedMessage(reqCtx, dayOffStarted))
			}
		}

		/******************************************************************************************************************************************************
			4. Sending messages to the broker.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		for _, message := range msgs {
			m.mb.SendMessage(reqCtx, message, msgResponder)
		}
		for (len(msgs) > 0 && len(errs) < len(msgs)) || len(msgResponder) > 0 {
			errs = append(errs, <-msgResponder)
		}
		if err := multierror.Append(reqCtx.Err(), errs...).ErrorOrNil(); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to send messages to broker for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			5. Fetching all relevant fields that will be added to the history/bookkeeping.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[model.User](reqCtx, db, userHistoryKeys, &histories); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to get histories for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			6. Inserting history/bookkeeping data.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := dwhClient.Insert(reqCtx, historyColumns, historyInsertMetadata, histories); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to insert histories for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			7. Persisting the mining progress for the users.
		******************************************************************************************************************************************************/

		for _, usr := range referralsThatStoppedMining {
			if usr.IDT0 > 0 {
				t1ReferralsThatStoppedMining[usr.IDT0]++
			}
			if usr.IDTMinus1 > 0 {
				t2ReferralsThatStoppedMining[usr.IDTMinus1]++
			}
		}

		var pipeliner redis.Pipeliner
		if len(t1ReferralsThatStoppedMining)+len(t2ReferralsThatStoppedMining) > 0 {
			pipeliner = db.TxPipeline()
		} else {
			pipeliner = db.Pipeline()
		}

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if responses, err := pipeliner.Pipelined(reqCtx, func(pipeliner redis.Pipeliner) error {
			for id, value := range t1ReferralsThatStoppedMining {
				if err := pipeliner.HIncrBy(reqCtx, model.SerializedUsersKey(id), "active_t1_referrals", -int64(value)).Err(); err != nil {
					return err
				}
			}
			for id, value := range t2ReferralsThatStoppedMining {
				if err := pipeliner.HIncrBy(reqCtx, model.SerializedUsersKey(id), "active_t2_referrals", -int64(value)).Err(); err != nil {
					return err
				}
			}
			for _, value := range updatedUsers {
				if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			log.Error(errors.Wrapf(err, "[miner] [1]failed to persist mining process for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		} else {
			if len(errs) != 0 {
				errs = errs[:0]
			}
			for _, response := range responses {
				if err = response.Err(); err != nil {
					errs = append(errs, errors.Wrapf(err, "failed to `%v`", response.FullName()))
				}
			}
			if err = multierror.Append(nil, errs...).ErrorOrNil(); err != nil {
				log.Error(errors.Wrapf(err, "[miner] [2]failed to persist mining progress for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
				reqCancel()
				resetVars(false)

				continue
			}
		}

		batchNumber++
		reqCancel()
		resetVars(true)
	}

}
