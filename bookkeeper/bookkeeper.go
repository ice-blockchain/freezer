// SPDX-License-Identifier: ice License 1.0

package bookkeeper

import (
	"context"
	"fmt"
	"sync"
	stdlibtime "time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	"github.com/ice-blockchain/freezer/model"
	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
)

func init() {
	appCfg.MustLoadFromKey(parentApplicationYamlKey, &cfg.Config)
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
}

func MustStartBookkeeping(ctx context.Context) {
	bk := &bookkeeper{}
	defer func() { log.Panic(errors.Wrap(bk.Close(), "failed to stop bookkeeper")) }()

	wg := new(sync.WaitGroup)
	wg.Add(int(cfg.Workers))
	defer wg.Wait()

	for workerNumber := int64(0); workerNumber < cfg.Workers; workerNumber++ {
		go func(wn int64) {
			defer wg.Done()
			bk.bookKeep(ctx, wn)
		}(workerNumber)
	}
}

func (bk *bookkeeper) Close() error {
	return nil
}

func (bk *bookkeeper) bookKeep(ctx context.Context, workerNumber int64) {
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
		historyKey                            = fmt.Sprintf("user_historical_chunks:%v", workerNumber)
		userResults                           = make([]*model.User, 0, cfg.BatchSize)
		errs                                  = make([]error, 0, cfg.BatchSize)
		historyColumns, historyInsertMetadata = dwh.InsertDDL(int(cfg.BatchSize))
	)
	for ctx.Err() == nil {
		/******************************************************************************************************************************************************
			1. Fetching a new batch of users.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel := context.WithTimeout(context.Background(), requestDeadline)
		userKeys, err := db.LRange(reqCtx, historyKey, 0, cfg.BatchSize-1).Result()
		reqCancel()

		if err != nil || len(userKeys) == 0 {
			log.Error(errors.Wrapf(err, "[bookkeeper] failed to LRange for users for workerNumber:%v", workerNumber))
			stdlibtime.Sleep(stdlibtime.Duration(10*workerNumber) * stdlibtime.Millisecond)

			continue
		}

		/******************************************************************************************************************************************************
			2. Getting the historical data for that batch.
		******************************************************************************************************************************************************/

		if len(userResults) != 0 {
			userResults = userResults[:0]
		}
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		err = storage.Bind[model.User](reqCtx, db, userKeys, &userResults)
		reqCancel()

		if err != nil {
			log.Error(errors.Wrapf(err, "[bookkeeper] failed to get users for workerNumber:%v", workerNumber))

			continue
		}

		/******************************************************************************************************************************************************
			3. Sending data to analytics/dwh storage.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		err = dwhClient.Insert(reqCtx, historyColumns, historyInsertMetadata, userResults)
		reqCancel()

		if err != nil {
			log.Error(errors.Wrapf(err, "[bookkeeper] failed to XXXXX for workerNumber:%v", workerNumber))

			continue
		}

		/******************************************************************************************************************************************************
			4. Deleting historical data from originating storage.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		responses, err := db.Pipelined(reqCtx, func(pipeliner redis.Pipeliner) error {
			return multierror.Append( //nolint:wrapcheck // Not needed.
				pipeliner.Del(reqCtx, userKeys...).Err(),
				pipeliner.LPopCount(reqCtx, historyKey, len(userKeys)).Err(),
			).ErrorOrNil()
		})
		reqCancel()

		if len(errs) != 0 {
			errs = errs[:0]
		}
		for _, response := range responses {
			if rErr := response.Err(); rErr != nil {
				errs = append(errs, errors.Wrapf(rErr, "failed to `%v`", response.FullName()))
			}
		}
		if rErr := multierror.Append(err, errs...).ErrorOrNil(); rErr != nil {
			log.Error(errors.Wrapf(rErr, "[bookkeeper] failed to del originating historical data for workerNumber:%v", workerNumber))

			continue
		}
	}
}
