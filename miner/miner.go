// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	balancesynchronizer "github.com/ice-blockchain/freezer/balance-synchronizer"
	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	extrabonusnotifier "github.com/ice-blockchain/freezer/extra-bonus-notifier"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	storagePG "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func init() {
	appCfg.MustLoadFromKey(parentApplicationYamlKey, &cfg.Config)
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	cfg.disableAdvancedTeam = new(atomic.Pointer[[]string])
}

func MustStartMining(ctx context.Context, cancel context.CancelFunc) Client {
	mi := &miner{
		mb:        messagebroker.MustConnect(context.Background(), parentApplicationYamlKey),
		db:        storage.MustConnect(context.Background(), parentApplicationYamlKey, int(cfg.Workers)),
		dbPG:      storagePG.MustConnect(ctx, eskimoDDL, applicationYamlKey),
		dwhClient: dwh.MustConnect(context.Background(), applicationYamlKey),
		wg:        new(sync.WaitGroup),
		telemetry: new(telemetry).mustInit(cfg),
	}
	go mi.startDisableAdvancedTeamCfgSyncer(ctx)
	mi.wg.Add(int(cfg.Workers))
	mi.cancel = cancel
	mi.extraBonusStartDate = extrabonusnotifier.MustGetExtraBonusStartDate(ctx, mi.db)
	mi.extraBonusIndicesDistribution = extrabonusnotifier.MustGetExtraBonusIndicesDistribution(ctx, mi.db)
	mi.recalculationBalanceStartDate = mustGetRecalculationBalancesStartDate(ctx, mi.db)

	for workerNumber := int64(0); workerNumber < cfg.Workers; workerNumber++ {
		go func(wn int64) {
			defer mi.wg.Done()
			mi.mine(ctx, wn)
		}(workerNumber)
	}

	return mi
}

func (m *miner) Close() error {
	m.cancel()
	m.wg.Wait()

	return multierror.Append(
		errors.Wrap(m.mb.Close(), "failed to close mb"),
		errors.Wrap(m.db.Close(), "failed to close db"),
		errors.Wrap(m.dbPG.Close(), "failed to close db pg"),
		errors.Wrap(m.dwhClient.Close(), "failed to close dwh"),
	).ErrorOrNil()
}

func (m *miner) CheckHealth(ctx context.Context) error {
	if err := m.dwhClient.Ping(ctx); err != nil {
		return err
	}
	if err := m.checkDBHealth(ctx); err != nil {
		return err
	}
	type ts struct {
		TS *time.Time `json:"ts"`
	}
	now := ts{TS: time.Now()}
	bytes, err := json.MarshalContext(ctx, now)
	if err != nil {
		return errors.Wrapf(err, "[health-check] failed to marshal %#v", now)
	}
	responder := make(chan error, 1)
	m.mb.SendMessage(ctx, &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     cfg.MessageBroker.Topics[0].Name,
		Topic:   cfg.MessageBroker.Topics[0].Name,
		Value:   bytes,
	}, responder)

	return errors.Wrapf(<-responder, "[health-check] failed to send health check message to broker")
}

func (m *miner) checkDBHealth(ctx context.Context) error {
	if resp := m.db.Ping(ctx); resp.Err() != nil || resp.Val() != "PONG" {
		if resp.Err() == nil {
			resp.SetErr(errors.Errorf("response `%v` is not `PONG`", resp.Val()))
		}

		return errors.Wrap(resp.Err(), "[health-check] failed to ping DB")
	}
	if !m.db.IsRW(ctx) {
		return errors.New("db is not writeable")
	}

	return nil
}

func mustGetRecalculationBalancesStartDate(ctx context.Context, db storage.DB) (recalculationBalancesStartDate *time.Time) {
	recalculationBalancesStartDateString, err := db.Get(ctx, "recalculation_balances_start_date").Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}
	log.Panic(errors.Wrap(err, "failed to get recalculation_balances_start_date"))
	if recalculationBalancesStartDateString != "" {
		recalculationBalancesStartDate = new(time.Time)
		log.Panic(errors.Wrapf(recalculationBalancesStartDate.UnmarshalText([]byte(recalculationBalancesStartDateString)), "failed to parse recalculation_balances_start_date `%v`", recalculationBalancesStartDateString)) //nolint:lll // .
		recalculationBalancesStartDate = time.New(recalculationBalancesStartDate.UTC())

		return
	}
	recalculationBalancesStartDate = time.Now()
	set, sErr := db.SetNX(ctx, "recalculation_balances_start_date", recalculationBalancesStartDate, 0).Result()
	log.Panic(errors.Wrap(sErr, "failed to set recalculation_balances_start_date"))
	if !set {
		return mustGetRecalculationBalancesStartDate(ctx, db)
	}

	return recalculationBalancesStartDate
}

func mustGetBalancesBackupMode(ctx context.Context, db storage.DB) (result bool, err error) {
	balancesBackupModeString, err := db.Get(ctx, "balances_backup_mode").Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}

	return balancesBackupModeString == "true", err
}

func (m *miner) mine(ctx context.Context, workerNumber int64) {
	dwhClient := dwh.MustConnect(context.Background(), applicationYamlKey)
	defer func() {
		if err := recover(); err != nil {
			log.Error(dwhClient.Close())
			panic(err)
		}
		log.Error(dwhClient.Close())
	}()
	var (
		batchNumber                                                          int64
		totalBatches                                                         uint64
		iteration                                                            uint64
		now, lastIterationStartedAt                                          = time.Now(), time.Now()
		currentAdoption                                                      = m.getAdoption(ctx, m.db, workerNumber)
		workers                                                              = cfg.Workers
		batchSize                                                            = cfg.BatchSize
		metrics                                                              = new(balanceRecalculationMetrics)
		userKeys, userBackupKeys, userHistoryKeys, referralKeys              = make([]string, 0, batchSize), make([]string, 0, batchSize), make([]string, 0, batchSize), make([]string, 0, 2*batchSize)
		userResults, backupUserResults, referralResults                      = make([]*user, 0, batchSize), make([]*backupUserUpdated, 0, batchSize), make([]*referral, 0, 2*batchSize)
		t0Referrals, tMinus1Referrals                                        = make(map[int64]*referral, batchSize), make(map[int64]*referral, batchSize)
		t1ReferralsToIncrementActiveValue, t2ReferralsToIncrementActiveValue = make(map[int64]int32, batchSize), make(map[int64]int32, batchSize)
		t1ReferralsThatStoppedMining, t2ReferralsThatStoppedMining           = make(map[int64]uint32, batchSize), make(map[int64]uint32, batchSize)
		referralsThatStoppedMining                                           = make([]*referralThatStoppedMining, 0, batchSize)
		msgResponder                                                         = make(chan error, 3*batchSize)
		msgs                                                                 = make([]*messagebroker.Message, 0, 3*batchSize)
		errs                                                                 = make([]error, 0, 3*batchSize)
		updatedUsers                                                         = make([]*UpdatedUser, 0, batchSize)
		extraBonusOnlyUpdatedUsers                                           = make([]*extrabonusnotifier.UpdatedUser, 0, batchSize)
		referralsUpdated                                                     = make([]*referralUpdated, 0, batchSize)
		histories                                                            = make([]*model.User, 0, batchSize)
		userGlobalRanks                                                      = make([]redis.Z, 0, batchSize)
		backupedUsers                                                        = make(map[int64]*backupUserUpdated, batchSize)
		backupUsersUpdated                                                   = make([]*backupUserUpdated, 0, batchSize)
		recalculatedTiersBalancesUsers                                       = make(map[int64]*user, batchSize)
		historyColumns, historyInsertMetadata                                = dwh.InsertDDL(int(batchSize))
		shouldSynchronizeBalanceFunc                                         = func(batchNumberArg uint64) bool { return false }
		recalculationHistory                                                 *historyData
		allAdoptions                                                         []*tokenomics.Adoption[float64]
	)
	resetVars := func(success bool) {
		if success && len(userKeys) == int(batchSize) && len(userResults) == 0 {
			go m.telemetry.collectElapsed(0, *lastIterationStartedAt.Time)
			lastIterationStartedAt = time.Now()
			iteration++
			if batchNumber < 1 {
				panic("unexpected batch number: " + fmt.Sprint(batchNumber))
			}
			totalBatches = uint64(batchNumber - 1)
			metricsExists, err := m.getBalanceRecalculationMetrics(ctx, workerNumber)
			if err != nil {
				log.Error(err, "can't get balance recalculation metrics for worker:", workerNumber)
			}
			if metricsExists == nil && err == nil {
				metrics.IterationsNum = int64(totalBatches)
				metrics.EndedAt = time.Now()
				metrics.Worker = workerNumber
				if err := m.insertBalanceRecalculationMetrics(ctx, metrics); err != nil {
					log.Error(err, "can't insert balance recalculation metrics for worker:", workerNumber)
				}
				metrics.reset()
			}
			if totalBatches != 0 && iteration > 2 {
				shouldSynchronizeBalanceFunc = m.telemetry.shouldSynchronizeBalanceFunc(uint64(workerNumber), totalBatches, iteration)
			}
			batchNumber = 0
		} else if success {
			go m.telemetry.collectElapsed(1, *now.Time)
		}
		now = time.Now()
		if batchNumber == 0 || currentAdoption == nil {
			currentAdoption = m.getAdoption(ctx, m.db, workerNumber)
		}
		userKeys, userBackupKeys, userHistoryKeys, referralKeys = userKeys[:0], userBackupKeys[:0], userHistoryKeys[:0], referralKeys[:0]
		userResults, referralResults = userResults[:0], referralResults[:0]
		msgs, errs = msgs[:0], errs[:0]
		updatedUsers = updatedUsers[:0]
		extraBonusOnlyUpdatedUsers = extraBonusOnlyUpdatedUsers[:0]
		referralsUpdated = referralsUpdated[:0]
		histories = histories[:0]
		userGlobalRanks = userGlobalRanks[:0]
		referralsThatStoppedMining = referralsThatStoppedMining[:0]
		allAdoptions = allAdoptions[:0]
		backupUsersUpdated = backupUsersUpdated[:0]
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
		for k := range t1ReferralsToIncrementActiveValue {
			delete(t1ReferralsToIncrementActiveValue, k)
		}
		for k := range t2ReferralsToIncrementActiveValue {
			delete(t2ReferralsToIncrementActiveValue, k)
		}
		for k := range backupedUsers {
			delete(backupedUsers, k)
		}
		for k := range recalculatedTiersBalancesUsers {
			delete(recalculatedTiersBalancesUsers, k)
		}
	}
	metrics.StartedAt = time.Now()
	for ctx.Err() == nil {
		/******************************************************************************************************************************************************
			1. Fetching a new batch of users.
		******************************************************************************************************************************************************/
		if len(userKeys) == 0 {
			for ix := batchNumber * batchSize; ix < (batchNumber+1)*batchSize; ix++ {
				userKeys = append(userKeys, model.SerializedUsersKey((workers*ix)+workerNumber))
			}
		}
		before := time.Now()
		reqCtx, reqCancel := context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[user](reqCtx, m.db, userKeys, &userResults); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to get users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			now = time.Now()

			continue
		}
		reqCancel()
		if len(userKeys) > 0 {
			go m.telemetry.collectElapsed(2, *before.Time)
		}
		for ix := batchNumber * batchSize; ix < (batchNumber+1)*batchSize; ix++ {
			userBackupKeys = append(userBackupKeys, model.SerializedBackupUsersKey((workers*ix)+workerNumber))
		}
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[backupUserUpdated](reqCtx, m.db, userBackupKeys, &backupUserResults); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to get backuped users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			now = time.Now()

			continue
		}
		balanceBackupMode, err := mustGetBalancesBackupMode(reqCtx, m.db)
		if err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to get backup flag for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()

			continue
		}
		reqCancel()
		for _, usr := range backupUserResults {
			backupedUsers[usr.ID] = usr
		}

		/******************************************************************************************************************************************************
			2. Fetching T0 & T-1 referrals of the fetched users.
		******************************************************************************************************************************************************/

		for _, usr := range userResults {
			if usr.UserID == "" {
				continue
			}
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

		before = time.Now()
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[referral](reqCtx, m.db, referralKeys, &referralResults); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to get referrees for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()
		if len(referralKeys) > 0 {
			go m.telemetry.collectElapsed(3, *before.Time)
		}

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
		if !balanceBackupMode {
			reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
			recalculationHistory, err = m.gatherHistoryAndReferralsInformation(reqCtx, userResults)
			if err != nil {
				log.Error(errors.New("tiers diff balances error"), workerNumber, err)
			}
			reqCancel()
			reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
			allAdoptions, err = tokenomics.GetAllAdoptions[float64](reqCtx, m.db)
			if err != nil {
				log.Error(errors.New("can't get all adoptions"), workerNumber, err)
			}
			reqCancel()
		}
		shouldSynchronizeBalance := shouldSynchronizeBalanceFunc(uint64(batchNumber))
		for _, usr := range userResults {
			if usr.UserID == "" {
				continue
			}
			backupedUsr, backupExists := backupedUsers[usr.ID]
			if balanceBackupMode {
				if backupExists {
					diffT1ActiveValue := backupedUsr.ActiveT1Referrals - usr.ActiveT1Referrals
					diffT2ActiveValue := backupedUsr.ActiveT2Referrals - usr.ActiveT2Referrals
					if diffT1ActiveValue < 0 && diffT1ActiveValue*-1 > usr.ActiveT1Referrals {
						diffT1ActiveValue = -usr.ActiveT1Referrals
					}
					if diffT2ActiveValue < 0 && diffT2ActiveValue*-1 > usr.ActiveT2Referrals {
						diffT2ActiveValue = -usr.ActiveT2Referrals
					}
					if false {
						t1ReferralsToIncrementActiveValue[usr.ID] += diffT1ActiveValue
						t2ReferralsToIncrementActiveValue[usr.ID] += diffT2ActiveValue

						usr.BalanceT1 = backupedUsr.BalanceT1
						usr.BalanceT2 = backupedUsr.BalanceT2

						usr.SlashingRateT1 = backupedUsr.SlashingRateT1
						usr.SlashingRateT2 = backupedUsr.SlashingRateT2

						backupedUsr.BalancesBackupUsedAt = time.Now()
						backupUsersUpdated = append(backupUsersUpdated, backupedUsr)
					}
				}
			} else {
				if recalculatedUsr := m.recalculateUser(usr, allAdoptions, recalculationHistory); recalculatedUsr != nil {
					diffT1ActiveValue := recalculationHistory.T1ActiveCounts[usr.UserID] - usr.ActiveT1Referrals
					diffT2ActiveValue := recalculationHistory.T2ActiveCounts[usr.UserID] - usr.ActiveT2Referrals

					oldBalanceT1 := usr.BalanceT1
					oldBalanceT2 := usr.BalanceT2

					if diffT1ActiveValue < 0 && diffT1ActiveValue*-1 > usr.ActiveT1Referrals {
						diffT1ActiveValue = -usr.ActiveT1Referrals
					}
					if diffT2ActiveValue < 0 && diffT2ActiveValue*-1 > usr.ActiveT2Referrals {
						diffT2ActiveValue = -usr.ActiveT2Referrals
					}

					oldSlashingT1Rate := usr.SlashingRateT1
					oldSlashingT2Rate := usr.SlashingRateT2

					if false {
						t1ReferralsToIncrementActiveValue[usr.ID] += diffT1ActiveValue
						t2ReferralsToIncrementActiveValue[usr.ID] += diffT2ActiveValue

						usr.BalanceT1 = recalculatedUsr.BalanceT1
						usr.BalanceT2 = recalculatedUsr.BalanceT2

						usr.SlashingRateT1 = recalculatedUsr.SlashingRateT1
						usr.SlashingRateT2 = recalculatedUsr.SlashingRateT2
					}
					if !backupExists {
						metrics.AffectedUsers += 1
						if recalculatedUsr.BalanceT1-oldBalanceT1 >= 0 {
							metrics.T1BalancePositive += recalculatedUsr.BalanceT1 - oldBalanceT1
						} else {
							metrics.T1BalanceNegative += recalculatedUsr.BalanceT1 - oldBalanceT1
						}
						if recalculatedUsr.BalanceT2-oldBalanceT2 >= 0 {
							metrics.T2BalancePositive += recalculatedUsr.BalanceT2 - oldBalanceT2
						} else {
							metrics.T2BalanceNegative += recalculatedUsr.BalanceT2 - oldBalanceT2
						}
						if diffT1ActiveValue < 0 {
							metrics.T1ActiveCountsNegative += int64(diffT1ActiveValue)
						} else {
							metrics.T1ActiveCountsPositive += int64(diffT1ActiveValue)
						}
						if diffT2ActiveValue < 0 {
							metrics.T2ActiveCountsNegative += int64(diffT2ActiveValue)
						} else {
							metrics.T2ActiveCountsPositive += int64(diffT2ActiveValue)
						}

						backupUsersUpdated = append(backupUsersUpdated, &backupUserUpdated{
							DeserializedBackupUsersKey:              model.DeserializedBackupUsersKey{ID: usr.ID},
							UserIDField:                             usr.UserIDField,
							BalanceT1Field:                          model.BalanceT1Field{BalanceT1: oldBalanceT1},
							BalanceT2Field:                          model.BalanceT2Field{BalanceT2: oldBalanceT2},
							SlashingRateT1Field:                     model.SlashingRateT1Field{SlashingRateT1: oldSlashingT1Rate},
							SlashingRateT2Field:                     model.SlashingRateT2Field{SlashingRateT2: oldSlashingT2Rate},
							ActiveT1ReferralsField:                  model.ActiveT1ReferralsField{ActiveT1Referrals: usr.ActiveT1Referrals},
							ActiveT2ReferralsField:                  model.ActiveT2ReferralsField{ActiveT2Referrals: usr.ActiveT2Referrals},
							FirstRecalculatedBalanceT1Field:         model.FirstRecalculatedBalanceT1Field{FirstRecalculatedBalanceT1: recalculatedUsr.BalanceT1},
							FirstRecalculatedBalanceT2Field:         model.FirstRecalculatedBalanceT2Field{FirstRecalculatedBalanceT2: recalculatedUsr.BalanceT2},
							FirstRecalculatedSlashingRateT1Field:    model.FirstRecalculatedSlashingRateT1Field{FirstRecalculatedSlashingRateT1: recalculatedUsr.SlashingRateT1},
							FirstRecalculatedSlashingRateT2Field:    model.FirstRecalculatedSlashingRateT2Field{FirstRecalculatedSlashingRateT2: recalculatedUsr.SlashingRateT2},
							FirstRecalculatedActiveT1ReferralsField: model.FirstRecalculatedActiveT1ReferralsField{FirstRecalculatedActiveT1Referrals: usr.ActiveT1Referrals + diffT1ActiveValue},
							FirstRecalculatedActiveT2ReferralsField: model.FirstRecalculatedActiveT2ReferralsField{FirstRecalculatedActiveT2Referrals: usr.ActiveT2Referrals + diffT2ActiveValue},
						})
					}
				}
			}

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
			if isAdvancedTeamDisabled(usr.LatestDevice) {
				usr.ActiveT2Referrals = 0
			}
			updatedUser, shouldGenerateHistory, IDT0Changed := mine(currentAdoption.BaseMiningRate, now, usr, t0Ref, tMinus1Ref)
			if shouldGenerateHistory {
				userHistoryKeys = append(userHistoryKeys, usr.Key())
			}
			if updatedUser != nil {
				var extraBonusIndex uint16
				if isAvailable, _ := extrabonusnotifier.IsExtraBonusAvailable(now, m.extraBonusStartDate, updatedUser.ExtraBonusStartedAt, m.extraBonusIndicesDistribution, updatedUser.ID, int16(updatedUser.UTCOffset), &extraBonusIndex, &updatedUser.ExtraBonusDaysClaimNotAvailable, &updatedUser.ExtraBonusLastClaimAvailableAt); isAvailable {
					eba := &extrabonusnotifier.ExtraBonusAvailable{UserID: updatedUser.UserID, ExtraBonusIndex: extraBonusIndex}
					msgs = append(msgs, extrabonusnotifier.ExtraBonusAvailableMessage(reqCtx, eba))
				} else {
					updatedUser.ExtraBonusDaysClaimNotAvailable = 0
					updatedUser.ExtraBonusLastClaimAvailableAt = nil
				}
				if true || balanceBackupMode || !backupExists {
					if userStoppedMining := didReferralJustStopMining(now, usr, t0Ref, tMinus1Ref); userStoppedMining != nil {
						referralsThatStoppedMining = append(referralsThatStoppedMining, userStoppedMining)
					}
				}

				if dayOffStarted := didANewDayOffJustStart(now, usr); dayOffStarted != nil {
					msgs = append(msgs, dayOffStartedMessage(reqCtx, dayOffStarted))
				}
				if t0Ref != nil {
					if IDT0Changed {
						if !usr.BalanceLastUpdatedAt.IsNil() {
							log.Info(fmt.Sprintf("idT0 changed for:%v from:%v to:%v, t1 referrals for:%v were incremented by 1", usr.ID, usr.IDT0, updatedUser.IDT0, t0Ref.ID))

							t1ReferralsToIncrementActiveValue[t0Ref.ID]++
							if t0Ref.IDT0 != 0 {
								t2ReferralsToIncrementActiveValue[t0Ref.IDT0]++
							}
						}
						if usr.ActiveT1Referrals > 0 && t0Ref.ID != 0 {
							log.Info(fmt.Sprintf("idT0 changed for:%v from:%v to:%v, t2 referrals for:%v were incremented by: %v", usr.ID, usr.IDT0, updatedUser.IDT0, t0Ref.ID, usr.ActiveT1Referrals))

							t2ReferralsToIncrementActiveValue[t0Ref.ID] += usr.ActiveT1Referrals
						}
					}
					if usr.IDTMinus1 != t0Ref.IDT0 {
						updatedUser.IDTMinus1 = t0Ref.IDT0
					}
				}
				updatedUsers = append(updatedUsers, &updatedUser.UpdatedUser)
			} else {
				extraBonusOnlyUpdatedUsr := extrabonusnotifier.UpdatedUser{
					ExtraBonusLastClaimAvailableAtField:            usr.ExtraBonusLastClaimAvailableAtField,
					DeserializedUsersKey:                           usr.DeserializedUsersKey,
					ExtraBonusDaysClaimNotAvailableResettableField: model.ExtraBonusDaysClaimNotAvailableResettableField{ExtraBonusDaysClaimNotAvailable: usr.ExtraBonusDaysClaimNotAvailable},
				}
				if isAvailable, _ := extrabonusnotifier.IsExtraBonusAvailable(now, m.extraBonusStartDate, usr.ExtraBonusStartedAt, m.extraBonusIndicesDistribution, usr.ID, int16(usr.UTCOffset), &extraBonusOnlyUpdatedUsr.ExtraBonusIndex, &extraBonusOnlyUpdatedUsr.ExtraBonusDaysClaimNotAvailable, &extraBonusOnlyUpdatedUsr.ExtraBonusLastClaimAvailableAt); isAvailable {
					eba := &extrabonusnotifier.ExtraBonusAvailable{UserID: usr.UserID, ExtraBonusIndex: extraBonusOnlyUpdatedUsr.ExtraBonusIndex}
					msgs = append(msgs, extrabonusnotifier.ExtraBonusAvailableMessage(reqCtx, eba))
					extraBonusOnlyUpdatedUsers = append(extraBonusOnlyUpdatedUsers, &extraBonusOnlyUpdatedUsr)
				}
				if updUsr := updateT0AndTMinus1ReferralsForUserHasNeverMined(usr); updUsr != nil {
					referralsUpdated = append(referralsUpdated, updUsr)
					if t0Ref != nil && t0Ref.ID != 0 && usr.ActiveT1Referrals > 0 {
						t2ReferralsToIncrementActiveValue[t0Ref.ID] += usr.ActiveT1Referrals
					}
				}
			}
			totalStandardBalance, totalPreStakingBalance := usr.BalanceTotalStandard, usr.BalanceTotalPreStaking
			if updatedUser != nil {
				totalStandardBalance, totalPreStakingBalance = updatedUser.BalanceTotalStandard, updatedUser.BalanceTotalPreStaking
			}
			totalBalance := totalStandardBalance + totalPreStakingBalance
			if shouldSynchronizeBalance {
				userGlobalRanks = append(userGlobalRanks, balancesynchronizer.GlobalRank(usr.ID, totalBalance))
				msgs = append(msgs, balancesynchronizer.BalanceUpdatedMessage(reqCtx, usr.UserID, totalStandardBalance, totalPreStakingBalance))
			}
		}

		/******************************************************************************************************************************************************
			4. Sending messages to the broker.
		******************************************************************************************************************************************************/

		before = time.Now()
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
		if len(msgs) > 0 {
			go m.telemetry.collectElapsed(4, *before.Time)
		}

		/******************************************************************************************************************************************************
			5. Fetching all relevant fields that will be added to the history/bookkeeping.
		******************************************************************************************************************************************************/

		before = time.Now()
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[model.User](reqCtx, m.db, userHistoryKeys, &histories); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to get histories for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()
		if len(userHistoryKeys) > 0 {
			go m.telemetry.collectElapsed(5, *before.Time)
		}

		/******************************************************************************************************************************************************
			6. Inserting history/bookkeeping data.
		******************************************************************************************************************************************************/

		before = time.Now()
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := dwhClient.Insert(reqCtx, historyColumns, historyInsertMetadata, histories); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to insert histories for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()
		if len(histories) > 0 {
			go m.telemetry.collectElapsed(6, *before.Time)
		}

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
		if len(t1ReferralsThatStoppedMining)+len(t2ReferralsThatStoppedMining)+len(extraBonusOnlyUpdatedUsers)+len(referralsUpdated)+len(userGlobalRanks)+len(backupUsersUpdated) > 0 {
			pipeliner = m.db.TxPipeline()
		} else {
			pipeliner = m.db.Pipeline()
		}
		before = time.Now()
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if responses, err := pipeliner.Pipelined(reqCtx, func(pipeliner redis.Pipeliner) error {
			for id, value := range t1ReferralsToIncrementActiveValue {
				if err := pipeliner.HIncrBy(reqCtx, model.SerializedUsersKey(id), "active_t1_referrals", int64(value)).Err(); err != nil {
					return err
				}
			}
			for id, value := range t2ReferralsToIncrementActiveValue {
				if err := pipeliner.HIncrBy(reqCtx, model.SerializedUsersKey(id), "active_t2_referrals", int64(value)).Err(); err != nil {
					return err
				}
			}
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
			for _, value := range extraBonusOnlyUpdatedUsers {
				if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
					return err
				}
			}
			for _, value := range referralsUpdated {
				if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
					return err
				}
			}
			for _, value := range backupUsersUpdated {
				if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
					return err
				}
			}
			if len(userGlobalRanks) > 0 {
				if err := pipeliner.ZAdd(reqCtx, "top_miners", userGlobalRanks...).Err(); err != nil {
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
		if len(t1ReferralsToIncrementActiveValue)+len(t2ReferralsToIncrementActiveValue)+len(t1ReferralsThatStoppedMining)+len(t2ReferralsThatStoppedMining)+len(updatedUsers)+len(extraBonusOnlyUpdatedUsers)+len(referralsUpdated)+len(backupUsersUpdated)+len(userGlobalRanks) > 0 {
			go m.telemetry.collectElapsed(7, *before.Time)
		}

		batchNumber++
		reqCancel()
		resetVars(true)
	}
}

func (m *miner) getAdoption(ctx context.Context, db storage.DB, workerNumber int64) (currentAdoption *tokenomics.Adoption[float64]) {
	for err := errors.New("init"); ctx.Err() == nil && err != nil; {
		reqCtx, reqCancel := context.WithTimeout(context.Background(), requestDeadline)
		currentAdoption, err = tokenomics.GetCurrentAdoption(reqCtx, db)
		reqCancel()
		log.Error(errors.Wrapf(err, "[miner] failed to GetCurrentAdoption for workerNumber:%v", workerNumber))
	}

	return currentAdoption
}

func (m *miner) startDisableAdvancedTeamCfgSyncer(ctx context.Context) {
	ticker := stdlibtime.NewTicker(5 * stdlibtime.Minute) //nolint:gosec,gomnd // Not an  issue.
	defer ticker.Stop()
	log.Panic(errors.Wrap(m.syncDisableAdvancedTeamCfg(ctx), "failed to syncDisableAdvancedTeamCfg"))

	for {
		select {
		case <-ticker.C:
			reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
			log.Error(errors.Wrap(m.syncDisableAdvancedTeamCfg(reqCtx), "failed to syncDisableAdvancedTeamCfg"))
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

func (m *miner) syncDisableAdvancedTeamCfg(ctx context.Context) error {
	result, err := m.db.Get(ctx, "disable_advanced_team_cfg").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errors.Wrap(err, "could not get `disable_advanced_team_cfg`")
	}
	var (
		oldCfg []string
		newCfg = strings.Split(strings.ReplaceAll(strings.ToLower(result), " ", ""), ",")
	)
	sort.SliceStable(newCfg, func(ii, jj int) bool { return newCfg[ii] < newCfg[jj] })
	if old := cfg.disableAdvancedTeam.Swap(&newCfg); old != nil {
		oldCfg = *old
	}
	if strings.Join(oldCfg, "") != strings.Join(newCfg, "") {
		log.Info(fmt.Sprintf("`disable_advanced_team_cfg` changed from: %#v, to: %#v", oldCfg, newCfg))
	}

	return nil
}

func isAdvancedTeamEnabled(device string) bool {
	if device == "" {
		return true
	}
	var disableAdvancedTeamFor []string
	if cfgVal := cfg.disableAdvancedTeam.Load(); cfgVal != nil {
		disableAdvancedTeamFor = *cfgVal
	}
	if len(disableAdvancedTeamFor) == 0 {
		return true
	}
	for _, disabled := range disableAdvancedTeamFor {
		if strings.EqualFold(device, disabled) {
			return false
		}
	}

	return true
}

func isAdvancedTeamDisabled(device string) bool {
	return !isAdvancedTeamEnabled(device)
}
