// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"fmt"
	"sort"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	storagePG "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

const (
	maxLimit int64 = 100000
)

type (
	pgUser struct {
		ID, ReferredBy string
	}

	splittedAdoptionByRange struct {
		TimePoint      *time.Time
		BaseMiningRate float64
	}

	balanceTMinus1RecalculationDryRun struct {
		UpdatedAt                            *time.Time
		OldTMinus1Balance, NewTMinus1Balance float64
		UserID                               string
		TMinus1ID                            string
	}

	balanceT2RecalculationDryRun struct {
		UpdatedAt                  *time.Time
		OldT2Balance, NewT2Balance float64
		UserID                     string
	}
)

func (m *miner) collectReferralsTier2(ctx context.Context, usersKeys []string) (
	t2Referrals map[string][]string, err error,
) {
	var offset int64 = 0
	t2Referrals = make(map[string][]string)
	for {
		sql := `SELECT
					t2.id AS id,
					t0.id AS referred_by
				FROM users t0
					JOIN users t1
						ON t1.referred_by = t0.id
					JOIN users t2
						ON t2.referred_by = t1.id
				WHERE t0.id = ANY($1)
					AND t2.referred_by != t2.id
					AND t2.username != t2.id
				LIMIT $2 OFFSET $3`
		rows, err := storagePG.Select[pgUser](ctx, m.dbPG, sql, usersKeys, maxLimit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "can't get referrals from pg for showing actual data")
		}
		if len(rows) == 0 {
			break
		}
		offset += maxLimit
		for _, row := range rows {
			if row.ReferredBy != "bogus" && row.ReferredBy != "icenetwork" && row.ID != "bogus" && row.ID != "icenetwork" {
				t2Referrals[row.ReferredBy] = append(t2Referrals[row.ReferredBy], row.ID)
			}
		}
	}

	return t2Referrals, nil
}

func getInternalIDsBatch(ctx context.Context, db storage.DB, keys ...string) ([]string, error) {
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
				results = append(results, model.SerializedUsersKey(val.(string)))
			}
		}

		return results, nil
	}
}

func getInternalIDs(ctx context.Context, db storage.DB, keys ...string) (result []string, err error) {
	var batchKeys []string
	for _, key := range keys {
		batchKeys = append(batchKeys, key)
		if len(batchKeys) >= int(cfg.BatchSize) {
			res, err := getInternalIDsBatch(ctx, db, batchKeys...)
			if err != nil {
				return nil, err
			}
			result = append(result, res...)
			batchKeys = batchKeys[:0]
		}
	}
	if len(batchKeys) > 0 {
		res, err := getInternalIDsBatch(ctx, db, batchKeys...)
		if err != nil {
			return nil, err
		}
		result = append(result, res...)
	}

	return result, nil
}

func getReferralsBatch(ctx context.Context, db storage.DB, keys ...string) ([]*recalculateReferral, error) {
	referrals := make([]*recalculateReferral, 0)
	if err := storage.Bind[recalculateReferral](ctx, db, keys, &referrals); err != nil {
		return nil, errors.Wrapf(err, "failed to get referrals for:%v", keys)
	}

	return referrals, nil
}

func getReferrals(ctx context.Context, db storage.DB, keys ...string) (result []*recalculateReferral, err error) {
	var batchKeys []string
	for _, key := range keys {
		batchKeys = append(batchKeys, key)
		if len(batchKeys) >= int(cfg.BatchSize) {
			referrals, err := getReferralsBatch(ctx, db, batchKeys...)
			if err != nil {
				return nil, err
			}
			result = append(result, referrals...)
			batchKeys = batchKeys[:0]
		}
	}
	if len(batchKeys) > 0 {
		referrals, err := getReferralsBatch(ctx, db, batchKeys...)
		if err != nil {
			return nil, err
		}
		result = append(result, referrals...)
	}

	return result, nil
}

func (m *miner) gatherReferralsInformation(ctx context.Context, users []*user) (map[string]*recalculateReferral, map[string][]string, error) {
	if len(users) == 0 {
		return nil, nil, nil
	}
	var userIDs []string

	for _, usr := range users {
		userIDs = append(userIDs, usr.UserID)
	}
	t2Referrals, err := m.collectReferralsTier2(ctx, userIDs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't get users")
	}
	if len(t2Referrals) == 0 {
		return nil, nil, nil
	}
	var serializedReferralsUserKeys []string
	for _, val := range t2Referrals {
		for _, id := range val {
			serializedReferralsUserKeys = append(serializedReferralsUserKeys, model.SerializedUsersKey(id))
		}
	}
	internalIDKeys, err := getInternalIDs(ctx, m.db, serializedReferralsUserKeys...)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get referrals internal id keys for:%v", users)
	}
	referrals, err := getReferrals(ctx, m.db, internalIDKeys...)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get referrals for:%v", users)
	}
	referralsCollection := make(map[string]*recalculateReferral, len(referrals))
	for _, ref := range referrals {
		referralsCollection[ref.UserID] = ref
	}

	return referralsCollection, t2Referrals, nil
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

func calculateTimeBounds(refTimeRange, usrRange *dwh.AdjustUserInfo) (*time.Time, *time.Time) {
	if refTimeRange.MiningSessionSoloStartedAt.After(*usrRange.MiningSessionSoloEndedAt.Time) || refTimeRange.MiningSessionSoloEndedAt.Before(*usrRange.MiningSessionSoloStartedAt.Time) {
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

func gatherHistory(ctx context.Context, dwhClient dwh.Client, keys []string) (history map[int64][]*dwh.AdjustUserInfo, err error) {
	if len(keys) == 0 {
		return nil, nil
	}
	offset := int64(0)
	historyTimeRanges := make(map[int64][]*dwh.AdjustUserInfo, 0)
	for {
		historyInformation, err := dwhClient.GetAdjustUserInformation(ctx, keys, startRecalculationsFrom, maxLimit, offset)
		if err != nil {
			return nil, errors.Wrapf(err, "can't get adjust user information for ids:#%v", keys)
		}
		for _, info := range historyInformation {
			historyTimeRanges[info.ID] = append(historyTimeRanges[info.ID], info)
		}
		if len(historyInformation) == 0 || len(historyInformation) < int(maxLimit) {
			break
		}
		offset += maxLimit
	}
	if len(historyTimeRanges) == 0 {
		return nil, nil
	}

	return historyTimeRanges, nil
}

func (m *miner) recalculateBalanceTMinus1(usr *user, adoptions []*tokenomics.Adoption[float64], history map[int64][]*dwh.AdjustUserInfo, baseTMinus1Balances map[int64]float64) *user {
	if adoptions == nil || history == nil || baseTMinus1Balances == nil {
		return nil
	}
	startTime, err := stdlibtime.Parse(timeLayout, startRecalculationsFrom)
	if err != nil {
		log.Panic(err, "can't parse start recalculations from time")
	}
	baseBalanceTMinus1, ok := baseTMinus1Balances[usr.ID]
	if !ok {
		baseBalanceTMinus1 = 0
	}
	idTMinus1 := usr.IDTMinus1
	if idTMinus1 < 0 {
		idTMinus1 *= -1
	}
	var (
		lastMiningSessionSoloEndedAt *time.Time
		now                          = time.Now()
	)
	if _, ok = history[idTMinus1]; ok {
		var isResurrected bool

		updatedUser := new(user)
		updatedUser.BalanceForTMinus1 = baseBalanceTMinus1

		for _, tminus1Range := range history[idTMinus1] {
			if !tminus1Range.MiningSessionSoloEndedAt.IsNil() {
				lastMiningSessionSoloEndedAt = tminus1Range.MiningSessionSoloEndedAt
			}
			if tminus1Range.MiningSessionSoloStartedAt.Before(startTime) {
				tminus1Range.MiningSessionSoloStartedAt = time.New(startTime)
			}

			/******************************************************************************************************************************************************
				1. Resurrection check.
			******************************************************************************************************************************************************/
			if !tminus1Range.ResurrectSoloUsedAt.IsNil() && !tminus1Range.ResurrectSoloUsedAt.IsZero() && tminus1Range.ResurrectSoloUsedAt.After(startTime) && !isResurrected {
				var resurrectDelta float64
				if timeSpent := tminus1Range.MiningSessionSoloStartedAt.Sub(*tminus1Range.MiningSessionSoloPreviouslyEndedAt.Time); cfg.Development {
					resurrectDelta = timeSpent.Minutes()
				} else {
					timeSpent := tminus1Range.MiningSessionSoloStartedAt.Sub(*tminus1Range.MiningSessionSoloPreviouslyEndedAt.Time)
					resurrectDelta = timeSpent.Hours()
				}
				updatedUser.BalanceForTMinus1 += updatedUser.SlashingRateForTMinus1 * resurrectDelta
				isResurrected = true
			} else {

				/******************************************************************************************************************************************************
					2. Slashing.
				******************************************************************************************************************************************************/
				if !tminus1Range.MiningSessionSoloPreviouslyEndedAt.IsNil() && !tminus1Range.MiningSessionSoloPreviouslyEndedAt.IsZero() {
					var (
						elapsedTimeFraction float64
						miningSessionRatio  float64
					)
					if tminus1Range.MiningSessionSoloPreviouslyEndedAt.Before(startTime) {
						tminus1Range.MiningSessionSoloPreviouslyEndedAt = time.New(startTime)
					}
					if timeSpent := tminus1Range.MiningSessionSoloStartedAt.Sub(*tminus1Range.MiningSessionSoloPreviouslyEndedAt.Time); cfg.Development {
						elapsedTimeFraction = timeSpent.Minutes()
						miningSessionRatio = 1
					} else {
						elapsedTimeFraction = timeSpent.Hours()
						miningSessionRatio = 24.
					}

					if elapsedTimeFraction > 0 {
						if updatedUser.SlashingRateForTMinus1 == 0 {
							updatedUser.SlashingRateForTMinus1 = updatedUser.BalanceForTMinus1 / 60. / miningSessionRatio
						}
						updatedUser.BalanceForTMinus1 -= updatedUser.SlashingRateForTMinus1 * elapsedTimeFraction
					}
				}
			}

			/******************************************************************************************************************************************************
				3. TMinus1 balance calculation for the current user time range.
			******************************************************************************************************************************************************/
			for _, timeRange := range history[usr.ID] {
				if timeRange.MiningSessionSoloStartedAt.Before(startTime) {
					timeRange.MiningSessionSoloStartedAt = time.New(startTime)
				}

				startedAt, endedAt := calculateTimeBounds(timeRange, tminus1Range)
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
					var elapsedTimeFraction float64
					if timeSpent := adoptionRange.TimePoint.Sub(*previousTimePoint.Time); cfg.Development {
						elapsedTimeFraction = timeSpent.Minutes()
					} else {
						elapsedTimeFraction = timeSpent.Hours()
					}
					rate := 5 * adoptionRange.BaseMiningRate * elapsedTimeFraction / 100

					updatedUser.BalanceForTMinus1 += rate

					if updatedUser.SlashingRateForTMinus1 > 0 {
						updatedUser.SlashingRateForTMinus1 = 0
					}

					previousTimePoint = adoptionRange.TimePoint
				}
			}
		}
		if !lastMiningSessionSoloEndedAt.IsNil() {
			if timeDiff := now.Sub(*lastMiningSessionSoloEndedAt.Time); cfg.Development {
				if timeDiff >= 60*stdlibtime.Minute {
					updatedUser.BalanceForTMinus1 = 0
				} else if timeDiff > 0 {
					elapsedTimeFraction := now.Truncate(stdlibtime.Minute).Sub(*lastMiningSessionSoloEndedAt.Time).Minutes()
					miningSessionRatio := 1.
					if updatedUser.SlashingRateForTMinus1 == 0 {
						updatedUser.SlashingRateForTMinus1 = updatedUser.BalanceForTMinus1 / 60. / miningSessionRatio
					}
					updatedUser.BalanceForTMinus1 -= updatedUser.SlashingRateForTMinus1 * elapsedTimeFraction
				}
			} else {
				if timeDiff >= 60*stdlibtime.Hour*24 {
					updatedUser.BalanceForTMinus1 = 0
				} else if timeDiff > 0 {
					elapsedTimeFraction := now.Truncate(stdlibtime.Hour).Sub(*lastMiningSessionSoloEndedAt.Time).Hours()
					miningSessionRatio := 24.
					if updatedUser.SlashingRateForTMinus1 == 0 {
						updatedUser.SlashingRateForTMinus1 = updatedUser.BalanceForTMinus1 / 60. / miningSessionRatio
					}
					updatedUser.BalanceForTMinus1 -= updatedUser.SlashingRateForTMinus1 * elapsedTimeFraction
				}
			}
		}
		if updatedUser.BalanceForTMinus1 < 0 {
			updatedUser.BalanceForTMinus1 = 0
		}

		return updatedUser
	}

	return nil
}

func (m *miner) insertBalanceTMinus1RecalculationDryRunBatch(ctx context.Context, infos []*balanceTMinus1RecalculationDryRun) error {
	if len(infos) == 0 {
		return nil
	}
	now := time.Now()
	paramCounter := 1
	params := []any{}
	sqlParams := []string{}
	for _, info := range infos {
		sqlParams = append(sqlParams, fmt.Sprintf("($%v, $%v, $%v, $%v, $%v)", paramCounter, paramCounter+1, paramCounter+2, paramCounter+3, paramCounter+4))
		paramCounter += 5
		params = append(params, now.Time, info.OldTMinus1Balance, info.NewTMinus1Balance, info.UserID, info.TMinus1ID)
	}

	sql := fmt.Sprintf(`INSERT INTO balance_tminus1_recalculation_dry_run(
							updated_at,
							old_tminus1_balance,
							new_tminus1_balance,
							user_id,
                            tminus1_id
						)
					VALUES %v
					ON CONFLICT(user_id) DO NOTHING`, strings.Join(sqlParams, ","))
	_, err := storagePG.Exec(ctx, m.dbPG, sql, params...)

	return errors.Wrap(err, "failed to insert dry run balanceTMinus1 info")
}

func (m *miner) insertBalanceT2RecalculationDryRunBatch(ctx context.Context, infos []*balanceT2RecalculationDryRun) error {
	if len(infos) == 0 {
		return nil
	}
	now := time.Now()
	paramCounter := 1
	params := []any{}
	sqlParams := []string{}
	for _, info := range infos {
		sqlParams = append(sqlParams, fmt.Sprintf("($%v, $%v, $%v, $%v)", paramCounter, paramCounter+1, paramCounter+2, paramCounter+3))
		paramCounter += 4
		params = append(params, now.Time, info.OldT2Balance, info.NewT2Balance, info.UserID)
	}

	sql := fmt.Sprintf(`INSERT INTO balance_t2_recalculation_dry_run(
							updated_at,
							old_t2_balance,
							new_t1_balance,
							user_id
						)
					VALUES %v
					ON CONFLICT(user_id) DO NOTHING`, strings.Join(sqlParams, ","))
	_, err := storagePG.Exec(ctx, m.dbPG, sql, params...)

	return errors.Wrap(err, "failed to insert dry run info")
}
