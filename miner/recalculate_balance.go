// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"sort"
	stdlibtime "time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	storagePG "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

const (
	maxLimit int64 = 10000
)

type (
	pgUser struct {
		Active         *users.NotExpired
		ID, ReferredBy string
		ReferralType   string
	}
	pgUserCreated struct {
		ID        string
		CreatedAt *time.Time
	}

	splittedAdoptionByRange struct {
		TimePoint      *time.Time
		BaseMiningRate float64
	}

	historyRangeTime struct {
		MiningSessionSoloStartedAt         *time.Time
		MiningSessionSoloEndedAt           *time.Time
		MiningSessionSoloLastStartedAt     *time.Time
		MiningSessionSoloPreviouslyEndedAt *time.Time
		CreatedAt                          *time.Time
		ResurrectSoloUsedAt                *time.Time
		SlashingRateSolo                   float64
		BalanceSolo                        float64
		BalanceT1Pending                   float64
		BalanceT1PendingApplied            float64
		BalanceT2Pending                   float64
		BalanceT2PendingApplied            float64
	}

	historyData struct {
		NeedToBeRecalculatedUsers      map[string]struct{}
		HistoryTimeRanges              map[string][]*historyRangeTime
		T1Referrals, T2Referrals       map[string][]string
		T1ActiveCounts, T2ActiveCounts map[string]int32
	}

	balanceRecalculationMetrics struct {
		StartedAt              *time.Time
		EndedAt                *time.Time
		T1BalancePositive      float64
		T1BalanceNegative      float64
		T2BalancePositive      float64
		T2BalanceNegative      float64
		T1ActiveCountsPositive int64
		T1ActiveCountsNegative int64
		T2ActiveCountsPositive int64
		T2ActiveCountsNegative int64
		IterationsNum          int64
		AffectedUsers          int64
		Worker                 int64
	}
)

func (m *miner) getUsers(ctx context.Context, users []*user) (map[string]*pgUserCreated, error) {
	var (
		userIDs []string
		offset  int64 = 0
		result        = make(map[string]*pgUserCreated, len(users))
	)
	for _, val := range users {
		userIDs = append(userIDs, val.UserID)
	}
	for {
		sql := `SELECT
					id,
					created_at
				FROM users
				WHERE id = ANY($1)
				LIMIT $2 OFFSET $3`
		rows, err := storagePG.Select[pgUserCreated](ctx, m.dbPG, sql, userIDs, maxLimit, offset)
		if err != nil {
			return nil, errors.Wrapf(err, "can't get users from pg for: %#v", userIDs)
		}
		if len(rows) == 0 {
			break
		}
		offset += maxLimit
		for _, row := range rows {
			result[row.ID] = row
		}
	}

	return result, nil
}

func (m *miner) collectTiers(ctx context.Context, needToBeRecalculatedUsers map[string]struct{}) (
	t1Referrals, t2Referrals map[string][]string, t1ActiveCounts, t2ActiveCounts map[string]int32, err error,
) {
	var (
		userIDs       = make([]string, 0, len(needToBeRecalculatedUsers))
		offset  int64 = 0
		now           = time.Now()
	)
	for key := range needToBeRecalculatedUsers {
		userIDs = append(userIDs, key)
	}
	t1ActiveCounts, t2ActiveCounts = make(map[string]int32, len(needToBeRecalculatedUsers)), make(map[string]int32, len(needToBeRecalculatedUsers))
	t1Referrals, t2Referrals = make(map[string][]string), make(map[string][]string)
	for {
		sql := `SELECT * FROM(
					SELECT
						id,
						referred_by,
						'T1' AS referral_type,
						(CASE 
							WHEN COALESCE(last_mining_ended_at, to_timestamp(1)) > $1
								THEN COALESCE(last_mining_ended_at, to_timestamp(1))
								ELSE NULL
						END) 														  AS active
					FROM users
					WHERE referred_by = ANY($2)
						AND referred_by != id
						AND username != id
					UNION ALL
					SELECT
						t2.id AS id,
						t0.id AS referred_by,
						'T2'  AS referral_type,
						(CASE 
							WHEN COALESCE(t2.last_mining_ended_at, to_timestamp(1)) > $1
								THEN COALESCE(t2.last_mining_ended_at, to_timestamp(1))
								ELSE NULL
						END) 														  AS active
					FROM users t0
						JOIN users t1
							ON t1.referred_by = t0.id
						JOIN users t2
							ON t2.referred_by = t1.id
					WHERE t0.id = ANY($2)
						AND t2.referred_by != t2.id
						AND t2.username != t2.id
				) X
				LIMIT $3 OFFSET $4`
		rows, err := storagePG.Select[pgUser](ctx, m.dbPG, sql, now.Time, userIDs, maxLimit, offset)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "can't get referrals from pg for showing actual data")
		}
		if len(rows) == 0 {
			break
		}
		offset += maxLimit
		for _, row := range rows {
			if row.ReferredBy != "bogus" && row.ReferredBy != "icenetwork" && row.ID != "bogus" && row.ID != "icenetwork" {
				if row.ReferralType == "T1" {
					t1Referrals[row.ReferredBy] = append(t1Referrals[row.ReferredBy], row.ID)
					if row.Active != nil && *row.Active {
						t1ActiveCounts[row.ReferredBy]++
					}
				} else if row.ReferralType == "T2" {
					t2Referrals[row.ReferredBy] = append(t2Referrals[row.ReferredBy], row.ID)
					if row.Active != nil && *row.Active {
						t2ActiveCounts[row.ReferredBy]++
					}
				} else {
					log.Panic("wrong tier type")
				}
			}
		}
	}

	return t1Referrals, t2Referrals, t1ActiveCounts, t2ActiveCounts, nil
}

func splitByAdoptionTimeRanges(adoptions []*tokenomics.Adoption[float64], startedAt, endedAt *time.Time) []splittedAdoptionByRange {
	var result []splittedAdoptionByRange

	currentMBR := adoptions[0].BaseMiningRate
	lastAchievedAt := adoptions[0].AchievedAt
	currentAchievedAtIdx := 0

	for idx, adptn := range adoptions {
		if adptn.AchievedAt.IsNil() {
			continue
		}
		if adptn.AchievedAt.Before(*startedAt.Time) {
			currentMBR = adptn.BaseMiningRate
		}
		if (adptn.AchievedAt.After(*startedAt.Time) || adptn.AchievedAt.Equal(*startedAt.Time)) &&
			adptn.AchievedAt.Before(*endedAt.Time) {
			result = append(result, splittedAdoptionByRange{
				TimePoint:      adptn.AchievedAt,
				BaseMiningRate: adptn.BaseMiningRate,
			})
		}
		if adptn.AchievedAt.After(*lastAchievedAt.Time) {
			currentAchievedAtIdx = idx
			lastAchievedAt = adptn.AchievedAt
		}
	}
	result = append(result,
		splittedAdoptionByRange{
			TimePoint:      startedAt,
			BaseMiningRate: currentMBR,
		},
	)
	if endedAt.After(*adoptions[currentAchievedAtIdx].AchievedAt.Time) {
		result = append(result,
			splittedAdoptionByRange{
				TimePoint:      endedAt,
				BaseMiningRate: adoptions[currentAchievedAtIdx].BaseMiningRate,
			})
	} else {
		result = append(result,
			splittedAdoptionByRange{
				TimePoint:      endedAt,
				BaseMiningRate: currentMBR,
			})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TimePoint.Before(*result[j].TimePoint.Time)
	})

	return result
}

func calculateTimeBounds(refTimeRange, usrRange *historyRangeTime) (*time.Time, *time.Time) {
	if refTimeRange.MiningSessionSoloStartedAt.After(*usrRange.MiningSessionSoloEndedAt.Time) || refTimeRange.MiningSessionSoloEndedAt.Before(*usrRange.MiningSessionSoloStartedAt.Time) || refTimeRange.SlashingRateSolo > 0 {
		return nil, nil
	}
	var startedAt, endedAt *time.Time
	if refTimeRange.MiningSessionSoloStartedAt.After(*usrRange.MiningSessionSoloStartedAt.Time) || refTimeRange.MiningSessionSoloStartedAt.Equal(*usrRange.MiningSessionSoloStartedAt.Time) {
		startedAt = refTimeRange.MiningSessionSoloStartedAt
	} else {
		startedAt = usrRange.MiningSessionSoloStartedAt
	}
	if refTimeRange.MiningSessionSoloEndedAt.Before(*usrRange.MiningSessionSoloEndedAt.Time) || refTimeRange.MiningSessionSoloEndedAt.Equal(*usrRange.MiningSessionSoloEndedAt.Time) {
		endedAt = refTimeRange.MiningSessionSoloEndedAt
	} else {
		endedAt = usrRange.MiningSessionSoloEndedAt
	}

	return startedAt, endedAt
}

func initializeEmptyUser(updatedUser, usr *user) *user {
	var newUser user
	newUser.ID = usr.ID
	newUser.UserID = usr.UserID
	newUser.IDT0 = usr.IDT0
	newUser.IDTMinus1 = usr.IDTMinus1
	newUser.BalanceLastUpdatedAt = nil

	return &newUser
}

func getInternalIDs(ctx context.Context, db storage.DB, keys ...string) ([]string, error) {
	if cmdResults, err := db.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		if err := pipeliner.MGet(ctx, keys...).Err(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	} else {
		results := make([]string, 0, len(cmdResults))
		for _, cmdResult := range cmdResults {
			sliceResult := cmdResult.(*redis.SliceCmd)
			for _, val := range sliceResult.Val() {
				if val == nil {
					continue
				}
				results = append(results, val.(string))
			}
		}

		return results, nil
	}
}

func (m *miner) gatherHistoryAndReferralsInformation(ctx context.Context, users []*user) (history *historyData, err error) {
	if len(users) == 0 {
		return nil, nil
	}
	var (
		needToBeRecalculatedUsers = make(map[string]struct{}, len(users))
		historyTimeRanges         = make(map[string][]*historyRangeTime)
	)
	usrs, err := m.getUsers(ctx, users)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get CreatedAt information for users:%#v", usrs)
	}
	for _, usr := range users {
		if usr.UserID == "" || usr.Username == "" {
			continue
		}
		if _, ok := usrs[usr.UserID]; ok {
			if usrs[usr.UserID].CreatedAt == nil || usrs[usr.UserID].CreatedAt.After(*m.recalculationBalanceStartDate.Time) {
				continue
			}
		}
		needToBeRecalculatedUsers[usr.UserID] = struct{}{}
	}
	if len(needToBeRecalculatedUsers) == 0 {
		return nil, nil
	}
	t1Referrals, t2Referrals, t1ActiveCounts, t2ActiveCounts, err := m.collectTiers(ctx, needToBeRecalculatedUsers)
	if err != nil {
		return nil, errors.Wrap(err, "can't get active users for users")
	}
	if len(t1Referrals) == 0 && len(t2Referrals) == 0 {
		return nil, nil
	}
	userKeys := make([]string, 0, len(t1Referrals)+len(t2Referrals)+len(needToBeRecalculatedUsers))
	for _, values := range t1Referrals {
		for _, val := range values {
			userKeys = append(userKeys, model.SerializedUsersKey(val))
		}
	}
	for _, values := range t2Referrals {
		for _, val := range values {
			userKeys = append(userKeys, model.SerializedUsersKey(val))
		}
	}
	for key, _ := range needToBeRecalculatedUsers {
		userKeys = append(userKeys, model.SerializedUsersKey(key))
	}
	userResults, err := getInternalIDs(ctx, m.db, userKeys...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get internal ids for:%v", userKeys)
	}
	offset := int64(0)
	for {
		historyInformation, err := m.dwhClient.GetAdjustUserInformation(ctx, userResults, maxLimit, offset)
		if err != nil {
			return nil, errors.Wrapf(err, "can't get adjust user information for ids:#%v", userResults)
		}
		if len(historyInformation) == 0 {
			break
		}
		offset += maxLimit
		for _, info := range historyInformation {
			historyTimeRanges[info.UserID] = append(historyTimeRanges[info.UserID], &historyRangeTime{
				MiningSessionSoloPreviouslyEndedAt: info.MiningSessionSoloPreviouslyEndedAt,
				MiningSessionSoloStartedAt:         info.MiningSessionSoloStartedAt,
				MiningSessionSoloEndedAt:           info.MiningSessionSoloEndedAt,
				ResurrectSoloUsedAt:                info.ResurrectSoloUsedAt,
				CreatedAt:                          info.CreatedAt,
				SlashingRateSolo:                   info.SlashingRateSolo,
				BalanceT1Pending:                   info.BalanceT1Pending,
				BalanceT1PendingApplied:            info.BalanceT1PendingApplied,
				BalanceT2Pending:                   info.BalanceT2Pending,
				BalanceT2PendingApplied:            info.BalanceT2PendingApplied,
			})
		}
	}
	if len(historyTimeRanges) == 0 {
		return nil, nil
	}

	return &historyData{
		NeedToBeRecalculatedUsers: needToBeRecalculatedUsers,
		HistoryTimeRanges:         historyTimeRanges,
		T1Referrals:               t1Referrals,
		T2Referrals:               t2Referrals,
		T1ActiveCounts:            t1ActiveCounts,
		T2ActiveCounts:            t2ActiveCounts,
	}, nil
}

func (m *miner) recalculateUser(usr *user, adoptions []*tokenomics.Adoption[float64], history *historyData) *user {
	if history == nil || history.HistoryTimeRanges == nil || (history.T1Referrals == nil && history.T2Referrals == nil) || adoptions == nil {
		return nil
	}
	if _, ok := history.NeedToBeRecalculatedUsers[usr.UserID]; !ok {
		return nil
	}
	if _, ok := history.HistoryTimeRanges[usr.UserID]; ok {
		var (
			isResurrected                              bool
			slashingLastEndedAt                        *time.Time
			lastMiningSessionSoloEndedAt               *time.Time
			previousUserStartedAt, previousUserEndedAt *time.Time
			now                                        = time.Now()
		)
		clonedUser1 := *usr
		updatedUser := &clonedUser1
		updatedUser.BalanceT1 = 0
		updatedUser.BalanceT2 = 0
		updatedUser.BalanceLastUpdatedAt = nil

		for _, usrRange := range history.HistoryTimeRanges[usr.UserID] {
			if updatedUser == nil {
				updatedUser = initializeEmptyUser(updatedUser, usr)
			}
			lastMiningSessionSoloEndedAt = usrRange.MiningSessionSoloEndedAt

			updatedUser.BalanceT1Pending = usrRange.BalanceT1Pending
			updatedUser.BalanceT1PendingApplied = usrRange.BalanceT1PendingApplied
			updatedUser.BalanceT2Pending = usrRange.BalanceT2Pending
			updatedUser.BalanceT2PendingApplied = usrRange.BalanceT2PendingApplied
			/******************************************************************************************************************************************************
				1. Resurrection check & handling.
			******************************************************************************************************************************************************/
			if !usrRange.ResurrectSoloUsedAt.IsNil() && usrRange.ResurrectSoloUsedAt.Unix() > 0 && !isResurrected {
				var resurrectDelta float64
				if timeSpent := usrRange.MiningSessionSoloStartedAt.Sub(*usrRange.MiningSessionSoloPreviouslyEndedAt.Time); cfg.Development {
					resurrectDelta = timeSpent.Minutes()
				} else {
					resurrectDelta = timeSpent.Hours()
				}
				updatedUser.BalanceT1 += updatedUser.SlashingRateT1 * resurrectDelta
				updatedUser.BalanceT2 += updatedUser.SlashingRateT2 * resurrectDelta
				updatedUser.SlashingRateT1 = 0
				updatedUser.SlashingRateT2 = 0

				isResurrected = true
			}
			/******************************************************************************************************************************************************
				2. Slashing calculations.
			******************************************************************************************************************************************************/
			if usrRange.SlashingRateSolo > 0 {
				if slashingLastEndedAt.IsNil() {
					slashingLastEndedAt = usrRange.MiningSessionSoloEndedAt
				}
				updatedUser.BalanceLastUpdatedAt = slashingLastEndedAt
				updatedUser.ResurrectSoloUsedAt = nil
				updatedUser, _, _ = mine(0., usrRange.CreatedAt, updatedUser, nil, nil)
				slashingLastEndedAt = usrRange.CreatedAt

				continue
			}
			if !slashingLastEndedAt.IsNil() && usrRange.MiningSessionSoloStartedAt.Sub(*slashingLastEndedAt.Time).Nanoseconds() > 0 {
				updatedUser.BalanceLastUpdatedAt = slashingLastEndedAt
				updatedUser.ResurrectSoloUsedAt = nil
				now := usrRange.MiningSessionSoloStartedAt
				updatedUser, _, _ = mine(0., now, updatedUser, nil, nil)
				slashingLastEndedAt = nil
			}
			/******************************************************************************************************************************************************
				3. Saving time range state for the next range for streaks case.
			******************************************************************************************************************************************************/
			if previousUserStartedAt != nil && previousUserStartedAt.Equal(*usrRange.MiningSessionSoloStartedAt.Time) &&
				previousUserEndedAt != nil && (usrRange.MiningSessionSoloEndedAt.After(*previousUserEndedAt.Time) ||
				usrRange.MiningSessionSoloEndedAt.Equal(*previousUserEndedAt.Time)) {

				previousUserStartedAt = usrRange.MiningSessionSoloStartedAt

				usrRange.MiningSessionSoloStartedAt = previousUserEndedAt
				previousUserEndedAt = usrRange.MiningSessionSoloEndedAt
			} else {
				previousUserStartedAt = usrRange.MiningSessionSoloStartedAt
				previousUserEndedAt = usrRange.MiningSessionSoloEndedAt
			}
			/******************************************************************************************************************************************************
				4. T1 Balance calculation for the current user time range.
			******************************************************************************************************************************************************/
			if _, ok := history.T1Referrals[usr.UserID]; ok {
				for _, refID := range history.T1Referrals[usr.UserID] {
					if _, ok := history.HistoryTimeRanges[refID]; ok {
						var previousT1MiningSessionStartedAt, previousT1MiningSessionEndedAt *time.Time
						for _, timeRange := range history.HistoryTimeRanges[refID] {
							if timeRange.SlashingRateSolo > 0 {
								continue
							}
							if previousT1MiningSessionStartedAt != nil && previousT1MiningSessionStartedAt.Equal(*timeRange.MiningSessionSoloStartedAt.Time) &&
								previousT1MiningSessionEndedAt != nil && (timeRange.MiningSessionSoloEndedAt.After(*previousT1MiningSessionEndedAt.Time) ||
								timeRange.MiningSessionSoloEndedAt.Equal(*previousT1MiningSessionEndedAt.Time)) {

								previousT1MiningSessionStartedAt = timeRange.MiningSessionSoloStartedAt
								timeRange.MiningSessionSoloStartedAt = previousT1MiningSessionEndedAt
								previousT1MiningSessionEndedAt = timeRange.MiningSessionSoloEndedAt
							} else {
								previousT1MiningSessionStartedAt = timeRange.MiningSessionSoloStartedAt
								previousT1MiningSessionEndedAt = timeRange.MiningSessionSoloEndedAt
							}
							startedAt, endedAt := calculateTimeBounds(timeRange, usrRange)
							if startedAt == nil && endedAt == nil {
								continue
							}

							adoptionRanges := splitByAdoptionTimeRanges(adoptions, startedAt, endedAt)

							var previousTimePoint *time.Time
							for _, adoptionRange := range adoptionRanges {
								if previousTimePoint == nil {
									previousTimePoint = adoptionRange.TimePoint

									continue
								}
								if previousTimePoint.Equal(*adoptionRange.TimePoint.Time) {
									continue
								}
								updatedUser.ActiveT1Referrals = 1
								updatedUser.ActiveT2Referrals = 0
								updatedUser.MiningSessionSoloStartedAt = previousTimePoint
								updatedUser.MiningSessionSoloEndedAt = time.New(adoptionRange.TimePoint.Add(1 * stdlibtime.Nanosecond))
								updatedUser.BalanceLastUpdatedAt = nil
								updatedUser.ResurrectSoloUsedAt = nil
								now := adoptionRange.TimePoint

								updatedUser, _, _ = mine(adoptionRange.BaseMiningRate, now, updatedUser, nil, nil)

								previousTimePoint = adoptionRange.TimePoint
							}
						}
					}
				}
			}
			/******************************************************************************************************************************************************
				5. T2 Balance calculation for the current user time range.
			******************************************************************************************************************************************************/
			if _, ok := history.T2Referrals[usr.UserID]; ok {
				for _, refID := range history.T2Referrals[usr.UserID] {
					if _, ok := history.HistoryTimeRanges[refID]; ok {
						var previousT2MiningSessionStartedAt, previousT2MiningSessionEndedAt *time.Time
						for _, timeRange := range history.HistoryTimeRanges[refID] {
							if timeRange.SlashingRateSolo > 0 {
								continue
							}
							if previousT2MiningSessionStartedAt != nil && previousT2MiningSessionStartedAt.Equal(*timeRange.MiningSessionSoloStartedAt.Time) &&
								previousT2MiningSessionEndedAt != nil && (timeRange.MiningSessionSoloEndedAt.After(*previousT2MiningSessionEndedAt.Time) ||
								timeRange.MiningSessionSoloEndedAt.Equal(*previousT2MiningSessionEndedAt.Time)) {

								previousT2MiningSessionStartedAt = timeRange.MiningSessionSoloStartedAt
								timeRange.MiningSessionSoloStartedAt = previousT2MiningSessionEndedAt
								previousT2MiningSessionEndedAt = timeRange.MiningSessionSoloEndedAt
							} else {
								previousT2MiningSessionEndedAt = timeRange.MiningSessionSoloEndedAt
								previousT2MiningSessionStartedAt = timeRange.MiningSessionSoloStartedAt
							}
							startedAt, endedAt := calculateTimeBounds(timeRange, usrRange)
							if startedAt == nil && endedAt == nil {
								continue
							}

							adoptionRanges := splitByAdoptionTimeRanges(adoptions, startedAt, endedAt)

							var previousTimePoint *time.Time
							for _, adoptionRange := range adoptionRanges {
								if previousTimePoint == nil {
									previousTimePoint = adoptionRange.TimePoint

									continue
								}
								if previousTimePoint.Equal(*adoptionRange.TimePoint.Time) {
									continue
								}
								updatedUser.ActiveT1Referrals = 0
								updatedUser.ActiveT2Referrals = 1
								updatedUser.MiningSessionSoloPreviouslyEndedAt = usr.MiningSessionSoloPreviouslyEndedAt
								updatedUser.MiningSessionSoloStartedAt = previousTimePoint
								updatedUser.MiningSessionSoloEndedAt = time.New(adoptionRange.TimePoint.Add(1 * stdlibtime.Nanosecond))
								updatedUser.BalanceLastUpdatedAt = nil
								updatedUser.ResurrectSoloUsedAt = nil
								now := adoptionRange.TimePoint

								updatedUser, _, _ = mine(adoptionRange.BaseMiningRate, now, updatedUser, nil, nil)

								previousTimePoint = adoptionRange.TimePoint
							}
						}
					}
				}
			}
		}
		if !lastMiningSessionSoloEndedAt.IsNil() {
			if timeDiff := now.Sub(*lastMiningSessionSoloEndedAt.Time); cfg.Development {
				if timeDiff >= 60*stdlibtime.Minute {
					updatedUser = nil
				}
			} else {
				if timeDiff >= 60*stdlibtime.Hour*24 {
					updatedUser = nil
				}
			}
		}
		if updatedUser == nil {
			updatedUser = initializeEmptyUser(updatedUser, usr)
		}

		return updatedUser
	}

	return nil
}

func (m *miner) getBalanceRecalculationMetrics(ctx context.Context, workerNumber int64) (brm *balanceRecalculationMetrics, err error) {
	sql := `SELECT * FROM balance_recalculation_metrics WHERE worker = $1`
	res, err := storagePG.Get[balanceRecalculationMetrics](ctx, m.dbPG, sql, workerNumber)
	if err != nil {
		if err == storagePG.ErrNotFound {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "failed to get balance recalculation metrics:%v", res)
	}

	return res, nil
}

func (m *miner) insertBalanceRecalculationMetrics(ctx context.Context, brm *balanceRecalculationMetrics) error {
	sql := `INSERT INTO balance_recalculation_metrics(
						worker,
						started_at,
						ended_at,
						t1_balance_positive,
						t1_balance_negative,
						t2_balance_positive,
						t2_balance_negative,
						t1_active_counts_positive,
						t1_active_counts_negative,
						t2_active_counts_positive,
						t2_active_counts_negative,
						iterations_num,
						affected_users
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := storagePG.Exec(ctx, m.dbPG, sql, brm.Worker, brm.StartedAt.Time, brm.EndedAt.Time, brm.T1BalancePositive, brm.T1BalanceNegative, brm.T2BalancePositive, brm.T2BalanceNegative,
		brm.T1ActiveCountsPositive, brm.T1ActiveCountsNegative, brm.T2ActiveCountsPositive, brm.T2ActiveCountsNegative, brm.IterationsNum, brm.AffectedUsers)

	return errors.Wrapf(err, "failed to insert metrics for worker:%v, params:%#v", brm.Worker, brm)
}

func (b *balanceRecalculationMetrics) reset() {
	b.EndedAt = nil
	b.AffectedUsers = 0
	b.IterationsNum = 0
	b.T1BalancePositive = 0
	b.T1BalanceNegative = 0
	b.T2BalancePositive = 0
	b.T2BalanceNegative = 0
	b.T1ActiveCountsPositive = 0
	b.T1ActiveCountsNegative = 0
	b.T2ActiveCountsPositive = 0
	b.T2ActiveCountsNegative = 0
	b.StartedAt = time.Now()
}
