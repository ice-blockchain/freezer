// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/eskimo/kyc/quiz"
	"github.com/ice-blockchain/eskimo/users"
	balancesynchronizer "github.com/ice-blockchain/freezer/balance-synchronizer"
	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	coindistribution "github.com/ice-blockchain/freezer/coin-distribution"
	extrabonusnotifier "github.com/ice-blockchain/freezer/extra-bonus-notifier"
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
	if cfg.SlashingDaysCount == 0 {
		log.Panic(errors.Errorf("slashingDaysCount is zero"))
	}
	cfg.disableAdvancedTeam = new(atomic.Pointer[[]string])
	cfg.coinDistributionCollectorSettings = new(atomic.Pointer[coindistribution.CollectorSettings])
	cfg.coinDistributionCollectorStartedAt = new(atomic.Pointer[time.Time])
	cfg.miningBoostLevels = new(atomic.Pointer[[]*tokenomics.MiningBoostLevel])

	levels := make([]*tokenomics.MiningBoostLevel, 0, len(cfg.MiningBoost.Levels))
	for dollars, level := range cfg.MiningBoost.Levels {
		level.ICEPrice = strconv.FormatFloat(dollars, 'f', 15, 64)
		levels = append(levels, level)
	}
	sort.SliceStable(levels, func(ii, jj int) bool {
		iiPrice, err := strconv.ParseFloat(levels[ii].ICEPrice, 64)
		log.Panic(err)
		jjPrice, err := strconv.ParseFloat(levels[jj].ICEPrice, 64)
		log.Panic(err)

		return iiPrice < jjPrice
	})
	cfg.miningBoostLevels.Store(&levels)
}

func MustStartMining(ctx context.Context, cancel context.CancelFunc) Client {
	mi := &miner{
		coinDistributionRepository: coindistribution.NewRepository(context.Background(), func() {}),
		mb:                         messagebroker.MustConnect(context.Background(), parentApplicationYamlKey),
		db:                         storage.MustConnect(context.Background(), parentApplicationYamlKey, int(cfg.Workers)),
		dwhClient:                  dwh.MustConnect(context.Background(), applicationYamlKey),
		wg:                         new(sync.WaitGroup),
		telemetry:                  new(telemetry).mustInit(cfg),
		//quizRepository:             quiz.NewReadRepository(context.Background()),
	}
	go mi.startDisableAdvancedTeamCfgSyncer(ctx)
	mi.wg.Add(int(cfg.Workers))
	mi.cancel = cancel
	mi.extraBonusStartDate = extrabonusnotifier.MustGetExtraBonusStartDate(ctx, mi.db)
	mi.mustInitCoinDistributionCollector(ctx)
	if isTenantInDistributionMode() {
		if false {
			mi.usersRepository = users.New(context.Background(), nil)
		}
		go mi.coinDistributionRepository.StartPrepareCoinDistributionsForReviewMonitor(ctx)
	}

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
	<-m.stopCoinDistributionCollectionWorkerManager

	errs := multierror.Append(
		errors.Wrap(m.mb.Close(), "failed to close mb"),
		errors.Wrap(m.db.Close(), "failed to close db"),
		errors.Wrap(m.dwhClient.Close(), "failed to close dwh"),
		errors.Wrap(m.coinDistributionRepository.Close(), "failed to close coinDistributionRepository"),
		//errors.Wrap(m.quizRepository.Close(), "failed to close quizClient"),
	)
	if false && isTenantInDistributionMode() {
		errs = multierror.Append(errs, errors.Wrap(m.usersRepository.Close(), "failed to close usersRepository"))
	}

	return errs.ErrorOrNil()
}

func (m *miner) CheckHealth(ctx context.Context) error {
	if err := m.coinDistributionRepository.CheckHealth(ctx); err != nil {
		return err
	}
	if err := m.dwhClient.Ping(ctx); err != nil {
		return err
	}
	if err := m.checkDBHealth(ctx); err != nil {
		return err
	}
	//if err := m.quizRepository.CheckHealth(ctx); err != nil {
	//	return err
	//}
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
		batchNumber                                                                                         int64
		totalBatches                                                                                        uint64
		iteration                                                                                           uint64
		now, lastIterationStartedAt                                                                         = time.Now(), time.Now()
		workers                                                                                             = cfg.Workers
		batchSize                                                                                           = cfg.BatchSize
		userKeys, userHistoryKeys, referralKeys, syncQuizUserIDs, syncMandatoryUserFieldsForDistributionIDs = make([]string, 0, batchSize), make([]string, 0, batchSize), make([]string, 0, 2*batchSize), make([]string, 0, batchSize), make([]string, 0, batchSize)
		userResults, referralResults                                                                        = make([]*user, 0, batchSize), make([]*referral, 0, 2*batchSize)
		t0Referrals, tMinus1Referrals                                                                       = make(map[int64]*referral, batchSize), make(map[int64]*referral, batchSize)
		t1ReferralsToIncrementActiveValue, t2ReferralsToIncrementActiveValue                                = make(map[int64]int32, batchSize), make(map[int64]int32, batchSize)
		t1ReferralsThatStoppedMining, t2ReferralsThatStoppedMining                                          = make(map[int64]uint32, batchSize), make(map[int64]uint32, batchSize)
		balanceT1EthereumIncr, balanceT2EthereumIncr                                                        = make(map[int64]float64, batchSize), make(map[int64]float64, batchSize)
		balanceT1WelcomeBonusIncr                                                                           = make(map[int64]float64, batchSize)
		pendingBalancesForTMinus1, pendingBalancesForT0                                                     = make(map[int64]float64, batchSize), make(map[int64]float64, batchSize)
		referralsThatStoppedMining                                                                          = make([]*referralThatStoppedMining, 0, batchSize)
		coinDistributions                                                                                   = make([]*coindistribution.ByEarnerForReview, 0, 4*batchSize)
		msgResponder                                                                                        = make(chan error, 3*batchSize)
		msgs                                                                                                = make([]*messagebroker.Message, 0, 3*batchSize)
		errs                                                                                                = make([]error, 0, 3*batchSize)
		updatedUsers                                                                                        = make([]*UpdatedUser, 0, batchSize)
		extraBonusOnlyUpdatedUsers                                                                          = make([]*extrabonusnotifier.UpdatedUser, 0, batchSize)
		referralsCountGuardOnlyUpdatedUsers                                                                 = make([]*referralCountGuardUpdatedUser, 0, batchSize)
		usersThatStoppedMiningForDistribution                                                               = make([]*userThatStoppedMiningForDistribution, 0, batchSize)
		referralsUpdated                                                                                    = make([]*referralUpdated, 0, batchSize)
		histories                                                                                           = make([]*model.User, 0, batchSize)
		quizStatuses                                                                                        = make(map[string]*quiz.QuizStatus, batchSize)
		mandatoryUserFieldsForDistributionProfileList                                                       = make(map[string]*users.MandatoryForDistributionFieldsProfile, batchSize)
		ranks                                                                                               = make(map[string]int, batchSize)
		userGlobalRanks                                                                                     = make([]redis.Z, 0, batchSize)
		historyColumns, historyInsertMetadata                                                               = dwh.InsertDDL(int(batchSize))
		shouldSynchronizeBalanceFunc                                                                        = func(batchNumberArg uint64) bool { return false }
		startedCoinDistributionCollecting                                                                   = isCoinDistributionCollectorEnabled(now)
	)
	if startedCoinDistributionCollecting {
		m.coinDistributionStartedSignaler <- struct{}{}
	}
	resetVars := func(success bool) {
		if success && len(userKeys) == int(batchSize) && len(userResults) == 0 {
			go m.telemetry.collectElapsed(0, *lastIterationStartedAt.Time)
			if !startedCoinDistributionCollecting && iteration%2 == 1 && isCoinDistributionCollectorEnabled(now) {
				m.coinDistributionStartedSignaler <- struct{}{}
				startedCoinDistributionCollecting = true
			}
			if startedCoinDistributionCollecting && iteration%2 == 0 && isCoinDistributionCollectorEnabled(now) {
				m.coinDistributionEndedSignaler <- struct{}{}
				m.coinDistributionWorkerMX.Lock()
				m.coinDistributionWorkerMX.Unlock()
				startedCoinDistributionCollecting = false
			}
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
		userKeys, userHistoryKeys, referralKeys = userKeys[:0], userHistoryKeys[:0], referralKeys[:0]
		userResults, referralResults = userResults[:0], referralResults[:0]
		syncQuizUserIDs = syncQuizUserIDs[:0]
		syncMandatoryUserFieldsForDistributionIDs = syncMandatoryUserFieldsForDistributionIDs[:0]
		msgs, errs = msgs[:0], errs[:0]
		updatedUsers = updatedUsers[:0]
		extraBonusOnlyUpdatedUsers = extraBonusOnlyUpdatedUsers[:0]
		referralsCountGuardOnlyUpdatedUsers = referralsCountGuardOnlyUpdatedUsers[:0]
		referralsUpdated = referralsUpdated[:0]
		histories = histories[:0]
		userGlobalRanks = userGlobalRanks[:0]
		referralsThatStoppedMining = referralsThatStoppedMining[:0]
		coinDistributions = coinDistributions[:0]
		usersThatStoppedMiningForDistribution = usersThatStoppedMiningForDistribution[:0]

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
		for k := range balanceT1WelcomeBonusIncr {
			delete(balanceT1WelcomeBonusIncr, k)
		}
		for k := range balanceT1EthereumIncr {
			delete(balanceT1EthereumIncr, k)
		}
		for k := range balanceT2EthereumIncr {
			delete(balanceT2EthereumIncr, k)
		}
		for k := range pendingBalancesForTMinus1 {
			delete(pendingBalancesForTMinus1, k)
		}
		for k := range pendingBalancesForT0 {
			delete(pendingBalancesForT0, k)
		}
		for k := range quizStatuses {
			delete(quizStatuses, k)
		}
		for k := range mandatoryUserFieldsForDistributionProfileList {
			delete(mandatoryUserFieldsForDistributionProfileList, k)
		}
		for k := range ranks {
			delete(ranks, k)
		}
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
			if !isAdvancedTeamDisabled(ref.LatestDevice) {
				if _, found := tMinus1Referrals[ref.ID]; found {
					tMinus1Referrals[ref.ID] = ref
				}
			}
			if _, found := t0Referrals[ref.ID]; found {
				t0Referrals[ref.ID] = ref
			}
		}
		shouldSynchronizeBalance := shouldSynchronizeBalanceFunc(uint64(batchNumber))

		if isTenantInDistributionMode() {
			reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
			t1Referrals := make([]*coindistribution.T1Referrals, 0, len(userResults))
			for _, usr := range userResults {
				if t0Referrals[usr.IDT0] == nil || !t0Referrals[usr.IDT0].isVerified() || !usr.isVerified() {
					continue
				}
				t1Referrals = append(t1Referrals, &coindistribution.T1Referrals{
					Balance:    usr.BalanceForT0,
					ReferredBy: t0Referrals[usr.IDT0].UserID,
					UserID:     usr.UserID,
				})
			}
			if err := m.coinDistributionRepository.InsertT1Referrals(reqCtx, t1Referrals); err != nil {
				reqCancel()
				log.Error(errors.Wrapf(err, "[miner] failed to insert t1 referrals for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
				resetVars(false)
				continue
			}
			reqCancel()

			if startedCoinDistributionCollecting {
				reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
				var pairs []coindistribution.ReferralPair
				for _, usr := range userResults {
					if t0Referrals[usr.IDT0] == nil || !t0Referrals[usr.IDT0].isVerified() || !usr.isVerified() {
						continue
					}
					pairs = append(pairs, coindistribution.ReferralPair{
						UserID:     usr.UserID,
						ReferredBy: t0Referrals[usr.IDT0].UserID,
					})
				}
				var err error
				ranks, err = m.coinDistributionRepository.CollectT1Ranks(reqCtx, pairs)
				if err != nil {
					reqCancel()
					log.Error(errors.Wrapf(err, "[miner] failed to get t1 ranks for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
					resetVars(false)
					continue
				}
				reqCancel()
			}
		}
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
			beforeWelcomeBonusV2NotApplied := usr.WelcomeBonusV2Applied == nil || !*usr.WelcomeBonusV2Applied
			updatedUser, shouldGenerateHistory, IDT0Changed, pendingAmountForTMinus1, pendingAmountForT0 := mine(now, usr, t0Ref, tMinus1Ref)
			if shouldGenerateHistory {
				syncQuizUserIDs = append(syncQuizUserIDs, usr.UserID)
				userHistoryKeys = append(userHistoryKeys, usr.Key())
			}
			if isTenantInDistributionMode() && !usr.isMandatoryFieldsSetForDistributionValid() && !usr.MiningSessionSoloStartedAt.IsNil() {
				syncMandatoryUserFieldsForDistributionIDs = append(syncMandatoryUserFieldsForDistributionIDs, usr.UserID)
			}

			if updatedUser != nil {
				if !isTenantInDistributionMode() {
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
							tMinus1Ref = tMinus1Referrals[updatedUser.IDTMinus1]
						}
					}
				} else {
					if updatedUser.MiningSessionSoloEndedAt.After(*now.Time) {
						usersThatStoppedMiningForDistribution = append(usersThatStoppedMiningForDistribution, &userThatStoppedMiningForDistribution{
							DeserializedUsersKey:                    usr.DeserializedUsersKey,
							ReferralsCountChangeGuardUpdatedAtField: model.ReferralsCountChangeGuardUpdatedAtField{ReferralsCountChangeGuardUpdatedAt: now},
							MiningSessionSoloEndedAtField:           model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: now},
						})
					}
				}
				userCoinDistributions, balanceDistributedForT0, balanceDistributedForTMinus1 := updatedUser.processEthereumCoinDistribution(startedCoinDistributionCollecting, now, t0Ref, tMinus1Ref, ranks, cfg.miningBoostLevels)
				coinDistributions = append(coinDistributions, userCoinDistributions...)
				if balanceDistributedForT0 > 0 {
					balanceT1EthereumIncr[t0Ref.ID] += balanceDistributedForT0
				}
				if balanceDistributedForTMinus1 > 0 {
					balanceT2EthereumIncr[tMinus1Ref.ID] += balanceDistributedForTMinus1
				}
				if !isTenantInDistributionMode() {
					if tMinus1Ref != nil && tMinus1Ref.ID != 0 && pendingAmountForTMinus1 != 0 {
						pendingBalancesForTMinus1[tMinus1Ref.ID] += pendingAmountForTMinus1
					}
					if t0Ref != nil && t0Ref.ID != 0 && pendingAmountForT0 != 0 {
						pendingBalancesForT0[t0Ref.ID] += pendingAmountForT0
					}
					if afterWelcomeBonusV2Applied := updatedUser.WelcomeBonusV2Applied != nil && *updatedUser.WelcomeBonusV2Applied; t0Ref != nil && t0Ref.ID != 0 && beforeWelcomeBonusV2NotApplied && afterWelcomeBonusV2Applied {
						idT0 := t0Ref.ID
						if idT0 < 0 {
							idT0 *= -1
						}
						balanceT1WelcomeBonusIncr[idT0] += cfg.WelcomeBonusV2Amount
					}
				}
				updatedUsers = append(updatedUsers, &updatedUser.UpdatedUser)
			} else {
				if !isTenantInDistributionMode() {
					if updUsr := updateT0AndTMinus1ReferralsForUserHasNeverMined(usr); updUsr != nil {
						referralsUpdated = append(referralsUpdated, updUsr)
						if t0Ref != nil && t0Ref.ID != 0 && usr.ActiveT1Referrals > 0 {
							t2ReferralsToIncrementActiveValue[t0Ref.ID] += usr.ActiveT1Referrals
						}
					}
				}
			}
			if !isTenantInDistributionMode() {
				totalStandardBalance, totalPreStakingBalance := usr.BalanceTotalStandard, usr.BalanceTotalPreStaking
				if updatedUser != nil {
					totalStandardBalance, totalPreStakingBalance = updatedUser.BalanceTotalStandard, updatedUser.BalanceTotalPreStaking
				}
				totalBalance := totalStandardBalance + totalPreStakingBalance
				if shouldSynchronizeBalance {
					userGlobalRanks = append(userGlobalRanks, balancesynchronizer.GlobalRank(usr.ID, totalBalance))
					if math.IsNaN(totalStandardBalance) || math.IsNaN(totalPreStakingBalance) {
						log.Info(fmt.Sprintf("bmr[%#v],before[%+v], after[%+v]", updatedUser.baseMiningRate(now), usr, updatedUser))
					}
					msgs = append(msgs, balancesynchronizer.BalanceUpdatedMessage(reqCtx, usr.UserID, totalStandardBalance, totalPreStakingBalance))
				}
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
			6. Syncing quiz state/syncing mandatory user fields for distribution
		******************************************************************************************************************************************************/
		before = time.Now()
		if false && (len(syncQuizUserIDs) > 0 && len(histories) > 0) {
			reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
			var err error
			quizStatuses, err = m.quizRepository.GetQuizStatus(reqCtx, syncQuizUserIDs...)
			if err != nil {
				log.Error(errors.Wrapf(err, "[miner] failed to sync quiz status (%v entries) batchNumber:%v,workerNumber:%v", len(syncQuizUserIDs), batchNumber, workerNumber))
				reqCancel()
				resetVars(false)

				continue
			}
			reqCancel()
			if len(quizStatuses) > 0 {
				for i := range histories {
					if quizSync, hasQuizSync := quizStatuses[histories[i].UserID]; hasQuizSync && quizSync != nil {
						histories[i].KYCQuizDisabled = quizSync.KYCQuizDisabled
						histories[i].KYCQuizCompleted = quizSync.KYCQuizCompleted
					}
				}
				go m.telemetry.collectElapsed(6, *before.Time)
			}
		}
		if false && len(syncMandatoryUserFieldsForDistributionIDs) > 0 {
			reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
			var err error
			mandatoryUserFieldsForDistributionProfileList, err = m.usersRepository.GetMandatoryForDistributionUserFieldsByIDList(reqCtx, syncMandatoryUserFieldsForDistributionIDs)
			if err != nil {
				log.Error(errors.Wrapf(err, "[miner] failed to sync mandatory fields for distribution (%v entries) batchNumber:%v,workerNumber:%v", len(syncMandatoryUserFieldsForDistributionIDs), batchNumber, workerNumber))
				reqCancel()
				resetVars(false)

				continue
			}
			reqCancel()
			if len(mandatoryUserFieldsForDistributionProfileList) > 0 {
				go m.telemetry.collectElapsed(6, *before.Time)
			}
		}

		/******************************************************************************************************************************************************
			6b. Insert user balances for distribution (PostgreSQL).
		******************************************************************************************************************************************************/
		if isTenantInDistributionMode() {
			before = time.Now()
			reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
			userBalances := make([]*coindistribution.UserBalance, 0, len(userResults))
			for _, usr := range userResults {
				if usr.UserID == "" {
					continue
				}
				email := usr.Email
				if mandatoryFieldsSync, ok := mandatoryUserFieldsForDistributionProfileList[usr.UserID]; ok && mandatoryFieldsSync != nil && mandatoryFieldsSync.Email != "" {
					email = mandatoryFieldsSync.Email
				}
				userBalances = append(userBalances, &coindistribution.UserBalance{
					UserID:        usr.UserID,
					Username:      usr.Username,
					Email:         email,
					Balance:       usr.BalanceTotalStandard + usr.BalanceTotalPreStaking,
					InternalID:    usr.ID,
					KYCStepPassed: int16(usr.KYCStepPassed),
					Verified:      usr.isVerified(),
				})
			}
			if len(userBalances) > 0 {
				log.Info(fmt.Sprintf("[miner] inserting user balances: count:%v, batchNumber:%v, totalBatches:%v, workerNumber:%v",
					len(userBalances), batchNumber, totalBatches, workerNumber))
			}
			if err := m.coinDistributionRepository.InsertUserBalances(reqCtx, userBalances); err != nil {
				log.Error(errors.Wrapf(err, "[miner] failed to insert user balances for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
				reqCancel()
				resetVars(false)

				continue
			}
			reqCancel()
			if len(userBalances) > 0 {
				go m.telemetry.collectElapsed(6, *before.Time)
			}
		}

		/******************************************************************************************************************************************************
			7. Inserting history/bookkeeping data.
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
			go m.telemetry.collectElapsed(7, *before.Time)
		}

		/******************************************************************************************************************************************************
			8. Processing Ethereum Coin Distributions for eligible users.
		******************************************************************************************************************************************************/

		before = time.Now()
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := m.coinDistributionRepository.CollectCoinDistributionsForReview(reqCtx, coinDistributions); err != nil {
			log.Error(errors.Wrapf(err, "[miner] failed to CollectCoinDistributionsForReview for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()
			resetVars(false)

			continue
		}
		reqCancel()
		if len(coinDistributions) > 0 {
			go m.telemetry.collectElapsed(8, *before.Time)
		}

		/******************************************************************************************************************************************************
			9. Persisting the mining progress for the users.
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
		var transactional bool
		if len(pendingBalancesForTMinus1)+len(pendingBalancesForT0)+len(balanceT1WelcomeBonusIncr)+len(balanceT1EthereumIncr)+len(balanceT2EthereumIncr)+len(t1ReferralsToIncrementActiveValue)+len(t2ReferralsToIncrementActiveValue)+len(referralsCountGuardOnlyUpdatedUsers)+len(t1ReferralsThatStoppedMining)+len(t2ReferralsThatStoppedMining)+len(extraBonusOnlyUpdatedUsers)+len(referralsUpdated)+len(userGlobalRanks)+len(usersThatStoppedMiningForDistribution) > 0 {
			pipeliner = m.db.TxPipeline()
			transactional = true
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
				if quizSync, hasQuizSync := quizStatuses[value.UserID]; hasQuizSync && quizSync != nil {
					disabled := model.FlexibleBool(quizSync.KYCQuizDisabled)
					completed := model.FlexibleBool(quizSync.KYCQuizCompleted)
					value.KYCQuizDisabled = &disabled
					value.KYCQuizCompleted = &completed
				}
				if mandatoryFieldsSync, hasMandatoryFieldsSync := mandatoryUserFieldsForDistributionProfileList[value.UserID]; hasMandatoryFieldsSync && mandatoryFieldsSync != nil {
					value.PhoneNumber = mandatoryFieldsSync.PhoneNumber
					value.Email = mandatoryFieldsSync.Email
					value.TelegramUserID = mandatoryFieldsSync.TelegramUserID
					value.TelegramBotID = mandatoryFieldsSync.TelegramBotID
					if mandatoryFieldsSync.DistributionScenariosVerified != nil {
						value.DistributionScenariosVerified = *mandatoryFieldsSync.DistributionScenariosVerified
					}
				}
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
			for idT0, amount := range balanceT1WelcomeBonusIncr {
				if err := pipeliner.HIncrByFloat(reqCtx, model.SerializedUsersKey(idT0), "balance_t1_welcome_bonus_pending", amount).Err(); err != nil {
					return err
				}
			}
			for idT0, amount := range balanceT1EthereumIncr {
				if amount == 0 {
					continue
				}
				if err := pipeliner.HIncrByFloat(reqCtx, model.SerializedUsersKey(idT0), "balance_t1_ethereum_pending", amount).Err(); err != nil {
					return err
				}
			}
			for idTMinus1, amount := range balanceT2EthereumIncr {
				if amount == 0 {
					continue
				}
				if err := pipeliner.HIncrByFloat(reqCtx, model.SerializedUsersKey(idTMinus1), "balance_t2_ethereum_pending", amount).Err(); err != nil {
					return err
				}
			}
			for idT0, amount := range pendingBalancesForT0 {
				if err := pipeliner.HIncrByFloat(reqCtx, model.SerializedUsersKey(idT0), "balance_t1_pending", amount).Err(); err != nil {
					return err
				}
			}
			for idTMinus1, amount := range pendingBalancesForTMinus1 {
				if err := pipeliner.HIncrByFloat(reqCtx, model.SerializedUsersKey(idTMinus1), "balance_t2_pending", amount).Err(); err != nil {
					return err
				}
			}
			if isTenantInDistributionMode() {
				for _, value := range usersThatStoppedMiningForDistribution {
					if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
						return err
					}
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
		if transactional || len(updatedUsers) > 0 {
			go m.telemetry.collectElapsed(9, *before.Time)
		}

		batchNumber++
		reqCancel()
		resetVars(true)
	}
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
	if true {
		return true
	}

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

func isTenantInDistributionMode() bool {
	return true
}
