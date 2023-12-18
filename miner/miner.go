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
		dwhClient: dwh.MustConnect(context.Background(), applicationYamlKey),
		wg:        new(sync.WaitGroup),
		telemetry: new(telemetry).mustInit(cfg),
		dbPG:      storagePG.MustConnect(context.Background(), eskimoDDL, applicationYamlKey),
	}
	go mi.startDisableAdvancedTeamCfgSyncer(ctx)
	mi.wg.Add(int(cfg.Workers))
	mi.cancel = cancel
	mi.extraBonusStartDate = extrabonusnotifier.MustGetExtraBonusStartDate(ctx, mi.db)
	mi.extraBonusIndicesDistribution = extrabonusnotifier.MustGetExtraBonusIndicesDistribution(ctx, mi.db)

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
		batchNumber                                                                       int64
		totalBatches                                                                      uint64
		iteration                                                                         uint64
		now, lastIterationStartedAt                                                       = time.Now(), time.Now()
		currentAdoption                                                                   = m.getAdoption(ctx, m.db, workerNumber)
		workers                                                                           = cfg.Workers
		batchSize                                                                         = cfg.BatchSize
		userKeys, userHistoryKeys, userInitialBalanceKeys, referralKeys, recalculatedKeys = make([]string, 0, batchSize), make([]string, 0, batchSize), make([]string, 0, 2*batchSize), make([]string, 0, batchSize), make([]string, 0, batchSize)
		dryRunKeys                                                                        = make([]string, 0, batchSize)
		userResults, referralResults, userRecalculatedResults, dryrunResults              = make([]*user, 0, batchSize), make([]*referral, 0, 2*batchSize), make([]*recalculated, 0, batchSize), make([]*dryrunUser, 0, batchSize)
		t0Referrals, tMinus1Referrals                                                     = make(map[int64]*referral, batchSize), make(map[int64]*referral, batchSize)
		t1ReferralsToIncrementActiveValue, t2ReferralsToIncrementActiveValue              = make(map[int64]int32, batchSize), make(map[int64]int32, batchSize)
		t1ReferralsThatStoppedMining, t2ReferralsThatStoppedMining                        = make(map[int64]uint32, batchSize), make(map[int64]uint32, batchSize)
		history                                                                           = make(map[int64][]*dwh.AdjustUserInfo, 2*batchSize)
		recalculatedUsers, dryRunUsers                                                    = make(map[int64]*recalculated, batchSize), make(map[int64]*dryrunUser, batchSize)
		recalculatedUsersUpdated, dryRunUsersUpdated                                      = make([]*recalculated, 0, batchSize), make([]*dryrunUser, 0, batchSize)
		referralsThatStoppedMining                                                        = make([]*referralThatStoppedMining, 0, batchSize)
		msgResponder                                                                      = make(chan error, 3*batchSize)
		msgs                                                                              = make([]*messagebroker.Message, 0, 3*batchSize)
		errs                                                                              = make([]error, 0, 3*batchSize)
		updatedUsers                                                                      = make([]*UpdatedUser, 0, batchSize)
		extraBonusOnlyUpdatedUsers                                                        = make([]*extrabonusnotifier.UpdatedUser, 0, batchSize)
		referralsCountGuardOnlyUpdatedUsers                                               = make([]*referralCountGuardUpdatedUser, 0, batchSize)
		referralsUpdated                                                                  = make([]*referralUpdated, 0, batchSize)
		histories                                                                         = make([]*model.User, 0, batchSize)
		userGlobalRanks                                                                   = make([]redis.Z, 0, batchSize)
		balanceTMinus1RecalculationDryRunItems                                            = make([]*balanceTMinus1RecalculationDryRun, 0, batchSize)
		balanceT2RecalculationDryRunItems                                                 = make([]*balanceT2RecalculationDryRun, 0, batchSize)
		balances                                                                          = make(map[int64]float64, 0)
		historyColumns, historyInsertMetadata                                             = dwh.InsertDDL(int(batchSize))
		shouldSynchronizeBalanceFunc                                                      = func(batchNumberArg uint64) bool { return false }
		allAdoptions                                                                      []*tokenomics.Adoption[float64]
		referralsCollection                                                               = make(map[string]*recalculateReferral, 0)
		t2Referrals                                                                       = make(map[string][]string, 0)
		clickhouseKeysMap                                                                 = make(map[string]struct{}, batchSize)
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
		userKeys, userHistoryKeys, referralKeys, recalculatedKeys, dryRunKeys = userKeys[:0], userHistoryKeys[:0], referralKeys[:0], recalculatedKeys[:0], dryRunKeys[:0]
		userResults, referralResults, userRecalculatedResults, recalculatedUsersUpdated, dryrunResults, dryRunUsersUpdated = userResults[:0], referralResults[:0], userRecalculatedResults[:0], recalculatedUsersUpdated[:0], dryrunResults[:0], dryRunUsersUpdated[:0]
		msgs, errs = msgs[:0], errs[:0]
		updatedUsers = updatedUsers[:0]
		extraBonusOnlyUpdatedUsers = extraBonusOnlyUpdatedUsers[:0]
		referralsCountGuardOnlyUpdatedUsers = referralsCountGuardOnlyUpdatedUsers[:0]
		referralsUpdated = referralsUpdated[:0]
		histories = histories[:0]
		userGlobalRanks = userGlobalRanks[:0]
		referralsThatStoppedMining = referralsThatStoppedMining[:0]
		allAdoptions = allAdoptions[:0]
		balanceTMinus1RecalculationDryRunItems = balanceTMinus1RecalculationDryRunItems[:0]
		balanceT2RecalculationDryRunItems = balanceT2RecalculationDryRunItems[:0]
		userInitialBalanceKeys = userInitialBalanceKeys[:0]
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
		for k := range history {
			delete(history, k)
		}
		for k := range balances {
			delete(balances, k)
		}
		for k := range referralsCollection {
			delete(referralsCollection, k)
		}
		for k := range t2Referrals {
			delete(t2Referrals, k)
		}
		for k := range recalculatedUsers {
			delete(recalculatedUsers, k)
		}
		for k := range clickhouseKeysMap {
			delete(clickhouseKeysMap, k)
		}
	}
	for ctx.Err() == nil {
		/******************************************************************************************************************************************************
			1. Fetching a new batch of users.
		******************************************************************************************************************************************************/
		if len(userKeys) == 0 {
			for ix := batchNumber * batchSize; ix < (batchNumber+1)*batchSize; ix++ {
				userKeys = append(userKeys, model.SerializedUsersKey((workers*ix)+workerNumber))
				if balanceForTMinusBugfixEnabled {
					if balanceForTMinusBugfixDryRunEnabled {
						dryRunKeys = append(dryRunKeys, model.SerializedDryRunUsersKey((workers*ix)+workerNumber))
					} else {
						recalculatedKeys = append(recalculatedKeys, model.SerializedRecalculatedUsersKey((workers*ix)+workerNumber))
					}
				}
				if clearBugfixDebugInfoEnabled {
					dryRunKeys = append(dryRunKeys, model.SerializedDryRunUsersKey((workers*ix)+workerNumber))
					recalculatedKeys = append(recalculatedKeys, model.SerializedRecalculatedUsersKey((workers*ix)+workerNumber))
				}
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
		if clearBugfixDebugInfoEnabled {
			_, err := m.db.Del(reqCtx, recalculatedKeys...).Result()
			if err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf("can't remove recalculated keys:%#v", recalculatedKeys)))
			}
			_, err = m.db.Del(reqCtx, dryRunKeys...).Result()
			if err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf("can't remove dry run keys:%#v", dryRunKeys)))
			}
		}
		reqCancel()
		if len(userKeys) > 0 {
			go m.telemetry.collectElapsed(2, *before.Time)
		}

		/******************************************************************************************************************************************************
			2. Fetching T0 & T-1 referrals of the fetched users.
		******************************************************************************************************************************************************/

		if balanceForTMinusBugfixEnabled {
			if balanceForTMinusBugfixDryRunEnabled {
				if len(dryRunKeys) > 0 {
					reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
					if err := storage.Bind[dryrunUser](reqCtx, m.db, dryRunKeys, &dryrunResults); err != nil {
						log.Error(errors.Wrapf(err, "[miner] failed to get dry run users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
					}
					reqCancel()
					for _, usr := range dryrunResults {
						dryRunUsers[usr.ID] = usr
					}
					for _, usr := range userResults {
						if _, ok := dryRunUsers[usr.ID]; !ok && usr.IDTMinus1 != 0 {
							userInitialBalanceKeys = append(userInitialBalanceKeys, fmt.Sprint(usr.ID))
							clickhouseKeysMap[fmt.Sprint(usr.ID)] = struct{}{}
						}
					}
				}
			} else {
				if len(recalculatedKeys) > 0 {
					reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
					if err := storage.Bind[recalculated](reqCtx, m.db, recalculatedKeys, &userRecalculatedResults); err != nil {
						log.Error(errors.Wrapf(err, "[miner] failed to get user recalculated users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
					}
					reqCancel()
					for _, usr := range userRecalculatedResults {
						recalculatedUsers[usr.ID] = usr
					}
					for _, usr := range userResults {
						if _, ok := recalculatedUsers[usr.ID]; !ok && usr.IDTMinus1 != 0 {
							userInitialBalanceKeys = append(userInitialBalanceKeys, fmt.Sprint(usr.ID))
							clickhouseKeysMap[fmt.Sprint(usr.ID)] = struct{}{}
						}
					}
				}
			}

			if len(userInitialBalanceKeys) > 0 {
				var err error
				reqCtx, reqCancel := context.WithTimeout(context.Background(), requestDeadline)
				balances, err = dwhClient.GetBaseBalanceForTMinus1(reqCtx, userInitialBalanceKeys, int64(len(userKeys)), 0)
				if err != nil {
					log.Error(errors.Wrapf(err, "[miner] failed to fetch base balances for tminus1 batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
				}
				reqCancel()
			}
		}

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
			idTMinus1Key := usr.IDTMinus1
			if usr.IDTMinus1 > 0 {
				tMinus1Referrals[usr.IDTMinus1] = nil
			}
			if usr.IDTMinus1 < 0 {
				tMinus1Referrals[-usr.IDTMinus1] = nil
				idTMinus1Key *= -1
			}
			if balanceForTMinusBugfixEnabled {
				if balanceForTMinusBugfixDryRunEnabled {
					if _, ok := dryRunUsers[usr.ID]; !ok {
						clickhouseKeysMap[fmt.Sprint(idTMinus1Key)] = struct{}{}
					}
				} else {
					if _, ok := recalculatedUsers[usr.ID]; !ok {
						clickhouseKeysMap[fmt.Sprint(idTMinus1Key)] = struct{}{}
					}
				}
			}
		}
		for idT0 := range t0Referrals {
			referralKeys = append(referralKeys, model.SerializedUsersKey(idT0))
		}
		if balanceForTMinusBugfixEnabled || balanceForTMinusBugfixDryRunEnabled {
			var err error
			if len(clickhouseKeysMap) > 0 {
				var clickhouseKeys []string
				for key, _ := range clickhouseKeysMap {
					clickhouseKeys = append(clickhouseKeys, key)
				}
				reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
				history, err = gatherHistory(reqCtx, dwhClient, clickhouseKeys)
				if err != nil {
					log.Error(err, fmt.Sprintf("can't gather history for: %#v", userResults))
				}
				reqCancel()
			}
			reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
			allAdoptions, err = tokenomics.GetAllAdoptions[float64](reqCtx, m.db)
			if err != nil {
				log.Error(err, fmt.Sprintf("can't gather adoptions for workerNumber: %v", workerNumber))
			}
			reqCancel()
		}
		if (balanceT2BugfixEnabled || balanceT2BugfixDryRunEnabled) && !balanceForTMinusBugfixEnabled {
			var err error
			reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
			referralsCollection, t2Referrals, err = m.gatherReferralsInformation(reqCtx, userResults)
			if err != nil {
				log.Error(errors.New("gather referrals info error"), workerNumber, err)
			}
			reqCancel()
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
		shouldSynchronizeBalance := shouldSynchronizeBalanceFunc(uint64(batchNumber))
		for _, usr := range userResults {
			if usr.UserID == "" {
				continue
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
			if balanceForTMinusBugfixEnabled && history != nil && allAdoptions != nil && balances != nil {
				if balanceForTMinusBugfixDryRunEnabled {
					if _, ok := dryRunUsers[usr.ID]; !ok {
						if recalculatedUsr := m.recalculateBalanceTMinus1(usr, allAdoptions, history, balances); recalculatedUsr != nil {
							tMinus1ID := ""
							if tMinus1Ref != nil {
								tMinus1ID = tMinus1Ref.UserID
							}
							balanceTMinus1RecalculationDryRunItems = append(balanceTMinus1RecalculationDryRunItems, &balanceTMinus1RecalculationDryRun{
								OldTMinus1Balance: usr.BalanceForTMinus1,
								NewTMinus1Balance: recalculatedUsr.BalanceForTMinus1,
								UserID:            usr.UserID,
								TMinus1ID:         tMinus1ID,
							})

							dryRunUsersUpdated = append(dryRunUsersUpdated, &dryrunUser{
								DeserializedDryRunUsersKey:           model.DeserializedDryRunUsersKey{ID: usr.ID},
								RecalculatedBalanceForTMinus1AtField: model.RecalculatedBalanceForTMinus1AtField{RecalculatedBalanceForTMinus1At: now},
							})
						}
					}
				} else {
					if _, ok := recalculatedUsers[usr.ID]; !ok {
						if recalculatedUsr := m.recalculateBalanceTMinus1(usr, allAdoptions, history, balances); recalculatedUsr != nil {
							tMinus1ID := ""
							if tMinus1Ref != nil {
								tMinus1ID = tMinus1Ref.UserID
							}
							balanceTMinus1RecalculationDryRunItems = append(balanceTMinus1RecalculationDryRunItems, &balanceTMinus1RecalculationDryRun{
								OldTMinus1Balance: usr.BalanceForTMinus1,
								NewTMinus1Balance: recalculatedUsr.BalanceForTMinus1,
								UserID:            usr.UserID,
								TMinus1ID:         tMinus1ID,
							})

							usr.BalanceForTMinus1 = recalculatedUsr.BalanceForTMinus1

							recalculatedUsersUpdated = append(recalculatedUsersUpdated, &recalculated{
								DeserializedRecalculatedUsersKey:     model.DeserializedRecalculatedUsersKey{ID: usr.ID},
								RecalculatedBalanceForTMinus1AtField: model.RecalculatedBalanceForTMinus1AtField{RecalculatedBalanceForTMinus1At: now},
							})
						}
					}
				}
			}
			if (balanceT2BugfixEnabled || balanceT2BugfixDryRunEnabled) && !balanceForTMinusBugfixEnabled {
				balanceT2 := 0.0
				if referralsT2, ok := t2Referrals[usr.UserID]; ok {
					for _, ref := range referralsT2 {
						if _, ok := referralsCollection[ref]; ok {
							balanceT2 += referralsCollection[ref].BalanceForTMinus1
						}
					}
				}
				balanceT2RecalculationDryRunItems = append(balanceT2RecalculationDryRunItems, &balanceT2RecalculationDryRun{
					OldT2Balance: usr.BalanceT2,
					NewT2Balance: balanceT2,
					UserID:       usr.UserID,
				})

				if !balanceT2BugfixDryRunEnabled {
					usr.BalanceT2 = balanceT2
				}
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
				if userStoppedMining := didUserStoppedMining(now, usr); userStoppedMining != nil {
					referralsCountGuardOnlyUpdatedUsers = append(referralsCountGuardOnlyUpdatedUsers, userStoppedMining)
				}
				if userStoppedMining := didReferralJustStopMining(now, usr, t0Ref, tMinus1Ref); userStoppedMining != nil {
					referralsThatStoppedMining = append(referralsThatStoppedMining, userStoppedMining)
				}
				if dayOffStarted := didANewDayOffJustStart(now, usr); dayOffStarted != nil {
					msgs = append(msgs, dayOffStartedMessage(reqCtx, dayOffStarted))
				}
				if t0Ref != nil {
					if IDT0Changed {
						if !usr.BalanceLastUpdatedAt.IsNil() {
							t1ReferralsToIncrementActiveValue[t0Ref.ID]++
							if t0Ref.IDT0 != 0 {
								t2ReferralsToIncrementActiveValue[t0Ref.IDT0]++
							}
						}
						if usr.ActiveT1Referrals > 0 && t0Ref.ID != 0 {
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
		if len(dryRunUsersUpdated)+len(recalculatedUsersUpdated)+len(t1ReferralsToIncrementActiveValue)+len(t2ReferralsToIncrementActiveValue)+len(referralsCountGuardOnlyUpdatedUsers)+len(t1ReferralsThatStoppedMining)+len(t2ReferralsThatStoppedMining)+len(extraBonusOnlyUpdatedUsers)+len(referralsUpdated)+len(userGlobalRanks) > 0 {
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
			for _, value := range referralsCountGuardOnlyUpdatedUsers {
				if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
					return err
				}
			}
			for _, value := range updatedUsers {
				if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
					return err
				}
			}
			for _, value := range recalculatedUsersUpdated {
				if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
					return err
				}
			}
			for _, value := range dryRunUsersUpdated {
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

		if len(t1ReferralsThatStoppedMining)+len(t2ReferralsThatStoppedMining)+len(updatedUsers)+len(extraBonusOnlyUpdatedUsers)+len(referralsUpdated)+len(userGlobalRanks) > 0 {
			go m.telemetry.collectElapsed(7, *before.Time)
		}
		if balanceForTMinusBugfixEnabled {
			if err := m.insertBalanceTMinus1RecalculationDryRunBatch(ctx, balanceTMinus1RecalculationDryRunItems); err != nil {
				log.Error(err, fmt.Sprintf("can't insert balance tminus1 recalculation dry run information for users:%#v, workerNumber:%v", userResults, workerNumber))
			}
		}
		if balanceT2BugfixDryRunEnabled {
			if err := m.insertBalanceT2RecalculationDryRunBatch(ctx, balanceT2RecalculationDryRunItems); err != nil {
				log.Error(err, fmt.Sprintf("can't insert balance t2 recalculation dry run information for users:%#v, workerNumber:%v", userResults, workerNumber))
			}
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

func didUserStoppedMining(now *time.Time, before *user) *referralCountGuardUpdatedUser {
	if !before.ReferralsCountChangeGuardUpdatedAt.IsNil() &&
		!before.MiningSessionSoloStartedAt.IsNil() &&
		!before.MiningSessionSoloEndedAt.IsNil() &&
		before.ReferralsCountChangeGuardUpdatedAt.Equal(*before.MiningSessionSoloStartedAt.Time) &&
		before.MiningSessionSoloEndedAt.Before(*now.Time) {
		return &referralCountGuardUpdatedUser{
			DeserializedUsersKey:                    before.DeserializedUsersKey,
			ReferralsCountChangeGuardUpdatedAtField: model.ReferralsCountChangeGuardUpdatedAtField{ReferralsCountChangeGuardUpdatedAt: time.Now()},
		}
	}

	return nil
}
