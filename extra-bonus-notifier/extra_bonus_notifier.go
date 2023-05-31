// SPDX-License-Identifier: ice License 1.0

package extrabonusnotifier

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/freezer/model"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func init() {
	appCfg.MustLoadFromKey(parentApplicationYamlKey, &cfg.messagebrokerConfig)
	appCfg.MustLoadFromKey(parentApplicationYamlKey, &cfg.ExtraBonusConfig)
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
}

func MustStartNotifyingExtraBonusAvailability(ctx context.Context) {
	ebs := &extraBonusNotifier{
		mb: messagebroker.MustConnect(context.Background(), parentApplicationYamlKey),
	}
	tmpDb := storage.MustConnect(context.Background(), parentApplicationYamlKey, 1)
	ebs.extraBonusStartDate = MustGetExtraBonusStartDate(ctx, tmpDb)
	ebs.extraBonusIndicesDistribution = MustGetExtraBonusIndicesDistribution(ctx, tmpDb)
	log.Panic(tmpDb.Close())

	defer func() { log.Panic(errors.Wrap(ebs.Close(), "failed to stop extraBonusNotifier")) }()

	wg := new(sync.WaitGroup)
	wg.Add(int(cfg.Workers))
	defer wg.Wait()

	for workerNumber := int64(0); workerNumber < cfg.Workers; workerNumber++ {
		go func(wn int64) {
			defer wg.Done()
			ebs.notifyingExtraBonusAvailability(ctx, wn)
		}(workerNumber)
	}
}

func MustGetExtraBonusStartDate(ctx context.Context, db storage.DB) (extraBonusStartDate *time.Time) {
	extraBonusStartDateString, err := db.Get(ctx, "extra_bonus_start_date").Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}
	log.Panic(errors.Wrap(err, "failed to get extra_bonus_start_date"))
	if extraBonusStartDateString != "" {
		extraBonusStartDate = new(time.Time)
		log.Panic(errors.Wrapf(extraBonusStartDate.UnmarshalText([]byte(extraBonusStartDateString)), "failed to parse extra_bonus_start_date `%v`", extraBonusStartDateString)) //nolint:lll // .

		return
	}
	extraBonusStartDate = time.New(stdlibtime.Now().Truncate(24 * stdlibtime.Hour))
	set, sErr := db.SetNX(ctx, "extra_bonus_start_date", extraBonusStartDate, 0).Result()
	log.Panic(errors.Wrap(sErr, "failed to set extra_bonus_start_date"))
	if !set {
		return MustGetExtraBonusStartDate(ctx, db)
	}

	return extraBonusStartDate
}

func MustGetExtraBonusIndicesDistribution(ctx context.Context, db storage.DB) map[uint16]map[uint16]uint16 {
	totalChunkNumber, totalExtraBonusDays := cfg.Chunks, uint16(len(cfg.ExtraBonuses.FlatValues))
	extraBonusIndicesDistribution := make(map[uint16]map[uint16]uint16, totalChunkNumber)
	flatResult, err := db.Get(ctx, "extra_bonus_distribution").Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}
	log.Panic(errors.Wrap(err, "failed to get extra_bonus_distribution"))
	if flatResult != "" {
		for _, elem := range strings.Split(flatResult, ",") {
			parts := strings.Split(elem, ":")
			i, cErr := strconv.Atoi(parts[0])
			log.Panic(cErr)
			j, cErr := strconv.Atoi(parts[1])
			log.Panic(cErr)
			k, cErr := strconv.Atoi(parts[2])
			log.Panic(cErr)
			if _, found := extraBonusIndicesDistribution[uint16(i)]; !found {
				extraBonusIndicesDistribution[uint16(i)] = make(map[uint16]uint16, totalExtraBonusDays)
			}
			extraBonusIndicesDistribution[uint16(i)][uint16(j)] = uint16(k)
		}

		return extraBonusIndicesDistribution
	}
	value := make([]string, 0, totalChunkNumber)
	for j := uint16(1); j <= totalExtraBonusDays; j++ {
		offsets := make([]uint16, totalChunkNumber)
		for i := uint16(0); i < totalChunkNumber; i++ {
			offsets[i] = i
		}
		rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(offsets), func(i, jj int) {
			offsets[i], offsets[jj] = offsets[jj], offsets[i]
		})
		for i := uint16(0); i < totalChunkNumber; i++ {
			if _, found := extraBonusIndicesDistribution[i]; !found {
				extraBonusIndicesDistribution[i] = make(map[uint16]uint16, totalExtraBonusDays)
			}
			extraBonusIndicesDistribution[i][j] = offsets[i]
			value = append(value, fmt.Sprintf("%v:%v:%v", i, j, offsets[i]))
		}
	}
	set, err := db.SetNX(ctx, "extra_bonus_distribution", strings.Join(value, ","), 0).Result()
	log.Panic(errors.Wrap(err, "failed to set extra_bonus_distribution"))
	if !set {
		return MustGetExtraBonusIndicesDistribution(ctx, db)
	}

	return extraBonusIndicesDistribution
}

func (ebn *extraBonusNotifier) Close() error {
	return multierror.Append(
		errors.Wrap(ebn.mb.Close(), "failed to close mb"),
	).ErrorOrNil()
}

func (ebn *extraBonusNotifier) notifyingExtraBonusAvailability(ctx context.Context, workerNumber int64) {
	db := storage.MustConnect(context.Background(), parentApplicationYamlKey, 1)
	defer func() {
		if err := recover(); err != nil {
			log.Error(db.Close())
			panic(err)
		}
		log.Error(db.Close())
	}()
	var (
		batchNumber  int64
		now          = time.Now()
		workers      = cfg.Workers
		batchSize    = cfg.BatchSize
		userKeys     = make([]string, 0, batchSize)
		userResults  = make([]*User, 0, batchSize)
		msgResponder = make(chan error, batchSize)
		msgs         = make([]*messagebroker.Message, 0, batchSize)
		errs         = make([]error, 0, batchSize)
		updatedUsers = make([]interface{ Key() string }, 0, batchSize)
	)
	resetVars := func(success bool) {
		now = time.Now()
		if success && len(userResults) < int(batchSize) {
			batchNumber = 0
		}
		userKeys = userKeys[:0]
		userResults = userResults[:0]
		msgs, errs = msgs[:0], errs[:0]
		updatedUsers = updatedUsers[:0]
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
		if err := storage.Bind[User](reqCtx, db, userKeys, &userResults); err != nil {
			log.Error(errors.Wrapf(err, "[extraBonusNotifier] failed to get users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			2. Processing batch.
		******************************************************************************************************************************************************/

		for _, usr := range userResults {
			if isAvailable, _ := IsExtraBonusAvailable(now, ebn.extraBonusStartDate, ebn.extraBonusIndicesDistribution, usr); isAvailable {
				eba := &ExtraBonusAvailable{UserID: usr.UserID, ExtraBonusIndex: usr.ExtraBonusIndex}
				updatedUsers = append(updatedUsers, &usr.UpdatedUser)
				msgs = append(msgs, extraBonusAvailableMessage(ctx, eba))
			}
		}

		/******************************************************************************************************************************************************
			3. Sending messages to the broker.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		for _, message := range msgs {
			ebn.mb.SendMessage(reqCtx, message, msgResponder)
		}
		for (len(msgs) > 0 && len(errs) < len(msgs)) || len(msgResponder) > 0 {
			errs = append(errs, <-msgResponder)
		}
		if err := multierror.Append(reqCtx.Err(), errs...).ErrorOrNil(); err != nil {
			log.Error(errors.Wrapf(err, "[extraBonusNotifier] failed to send messages to broker for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			4. Persisting the extra bonus availability progress for the users.
		******************************************************************************************************************************************************/

		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Set(reqCtx, db, updatedUsers...); err != nil {
			log.Error(errors.Wrapf(err, "[extraBonusNotifier] failed to persist the extra bonus availability progress for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber)) //nolint:lll // .
			reqCancel()
			resetVars(false)

			continue
		}

		batchNumber++
		reqCancel()
		resetVars(true)
	}
}

func extraBonusAvailableMessage(ctx context.Context, event *ExtraBonusAvailable) *messagebroker.Message {
	valueBytes, err := json.MarshalContext(ctx, event)
	log.Panic(errors.Wrapf(err, "failed to marshal %#v", event))

	return &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     event.UserID,
		Topic:   cfg.MessageBroker.Topics[4].Name,
		Value:   valueBytes,
	}
}
