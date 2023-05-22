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

	"github.com/ice-blockchain/freezer/tokenomics"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func init() {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
}

func MustStartNotifyingExtraBonusAvailability(ctx context.Context) {
	ebs := &extraBonusNotifier{
		db: storage.MustConnect(context.Background(), applicationYamlKey),
		mb: messagebroker.MustConnect(context.Background(), applicationYamlKey),
	}
	ebs.mustGetExtraBonusStartDate(ctx)
	ebs.mustGetExtraBonusIndicesDistribution(ctx)

	defer log.Panic(errors.Wrap(ebs.Close(), "failed to stop extraBonusNotifier"))

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

func (ebn *extraBonusNotifier) Close() error {
	return multierror.Append(
		errors.Wrap(ebn.db.Close(), "failed to close db"),
		errors.Wrap(ebn.mb.Close(), "failed to close mb"),
	).ErrorOrNil()
}

func (ebn *extraBonusNotifier) mustGetExtraBonusStartDate(ctx context.Context) {
	extraBonusStartDateString, err := ebn.db.Get(ctx, "extra_bonus_start_date").Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}
	log.Panic(errors.Wrap(err, "failed to get extra_bonus_start_date"))
	if extraBonusStartDateString != "" {
		ebn.extraBonusStartDate = new(time.Time)
		log.Panic(errors.Wrapf(ebn.extraBonusStartDate.UnmarshalText([]byte(extraBonusStartDateString)), "failed to parse extra_bonus_start_date `%v`", extraBonusStartDateString)) //nolint:lll // .

		return
	}
	ebn.extraBonusStartDate = time.New(stdlibtime.Now().Truncate(24 * stdlibtime.Hour))
	set, sErr := ebn.db.SetNX(ctx, "extra_bonus_start_date", ebn.extraBonusStartDate, 0).Result()
	log.Panic(errors.Wrap(sErr, "failed to set extra_bonus_start_date"))
	if !set {
		ebn.mustGetExtraBonusStartDate(ctx)
	}
}

func (ebn *extraBonusNotifier) mustGetExtraBonusIndicesDistribution(ctx context.Context) {
	totalChunkNumber, totalExtraBonusDays := cfg.Chunks, uint16(len(cfg.ExtraBonuses.FlatValues))
	ebn.extraBonusIndicesDistribution = make(map[uint16]map[uint16]uint16, totalChunkNumber)
	flatResult, err := ebn.db.Get(ctx, "extra_bonus_distribution").Result()
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
			if _, found := ebn.extraBonusIndicesDistribution[uint16(i)]; !found {
				ebn.extraBonusIndicesDistribution[uint16(i)] = make(map[uint16]uint16, totalExtraBonusDays)
			}
			ebn.extraBonusIndicesDistribution[uint16(i)][uint16(j)] = uint16(k)
		}

		return
	}
	value := make([]string, cfg.Chunks)
	for j := uint16(1); j <= totalExtraBonusDays; j++ {
		offsets := make([]uint16, totalChunkNumber)
		for i := uint16(0); i < totalChunkNumber; i++ {
			offsets[i] = i
		}
		rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(offsets), func(i, jj int) {
			offsets[i], offsets[jj] = offsets[jj], offsets[i]
		})
		for i := uint16(0); i < totalChunkNumber; i++ {
			if _, found := ebn.extraBonusIndicesDistribution[i]; !found {
				ebn.extraBonusIndicesDistribution[i] = make(map[uint16]uint16, totalExtraBonusDays)
			}
			ebn.extraBonusIndicesDistribution[i][j] = offsets[i]
			value = append(value, fmt.Sprintf("%v:%v:%v", i, j, offsets[i]))
		}
	}
	set, err := ebn.db.SetNX(ctx, "extra_bonus_distribution", strings.Join(value, ","), 0).Result()
	log.Panic(errors.Wrap(err, "failed to set extra_bonus_distribution"))
	if !set {
		ebn.mustGetExtraBonusIndicesDistribution(ctx)
	}
}

func (ebn *extraBonusNotifier) notifyingExtraBonusAvailability(ctx context.Context, workerNumber int64) {
	var (
		batchNumber  int64
		now          = time.Now()
		workers      = cfg.Workers
		batchSize    = cfg.BatchSize
		userKeys     = make([]string, 0, batchSize)
		userResults  = make([]*user, 0, batchSize)
		msgResponder = make(chan error, batchSize)
		msgs         = make([]*messagebroker.Message, 0, batchSize)
		errs         = make([]error, 0, batchSize)
		updatedUsers = make([]interface{ Key() string }, 0, batchSize)
	)
	resetVars := func(success bool) {
		now = time.Now()
		if success && len(userKeys) < int(batchSize) {
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
				userKeys = append(userKeys, tokenomics.SerializedUsersKey((workers*ix)+workerNumber))
			}
		}
		reqCtx, reqCancel := context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[user](reqCtx, ebn.db, userKeys, &userResults); err != nil {
			log.Error(errors.Wrapf(err, "[extraBonusNotifier] failed to get users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			2. Processing batch.
		******************************************************************************************************************************************************/

		for _, usr := range userResults {
			if isExtraBonusAvailable(now, ebn.extraBonusStartDate, ebn.extraBonusIndicesDistribution, usr) {
				eba := &ExtraBonusAvailable{UserID: usr.UserID, ExtraBonusIndex: usr.extraBonusIndex}
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
		for len(errs) < cap(errs) || len(msgResponder) > 0 {
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
		if err := storage.Set(reqCtx, ebn.db, updatedUsers...); err != nil {
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
