// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/model"
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
		CreatedAt *time.Time
		ID        string
	}

	historyData struct {
		Referrals                      map[string]*recalculateReferral
		NeedToBeRecalculatedUsers      map[string]struct{}
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

	balanceRecalculationDryRun struct {
		T1BalanceDiff      float64
		T2BalanceDiff      float64
		T1ActiveCountsDiff int32
		T2ActiveCountsDiff int32
		UserID             string
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
	referralsUserKeys []string, t1Referrals, t2Referrals map[string][]string, t1ActiveCounts, t2ActiveCounts map[string]int32, err error,
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
			return nil, nil, nil, nil, nil, errors.Wrap(err, "can't get referrals from pg for showing actual data")
		}
		if len(rows) == 0 {
			break
		}
		offset += maxLimit
		for _, row := range rows {
			if row.ReferredBy != "bogus" && row.ReferredBy != "icenetwork" && row.ID != "bogus" && row.ID != "icenetwork" {
				if row.ReferralType == "T1" {
					t1Referrals[row.ReferredBy] = append(t1Referrals[row.ReferredBy], row.ID)
					referralsUserKeys = append(referralsUserKeys, model.SerializedUsersKey(row.ID))
					if row.Active != nil && *row.Active {
						t1ActiveCounts[row.ReferredBy]++
					}
				} else if row.ReferralType == "T2" {
					t2Referrals[row.ReferredBy] = append(t2Referrals[row.ReferredBy], row.ID)
					referralsUserKeys = append(referralsUserKeys, model.SerializedUsersKey(row.ID))
					if row.Active != nil && *row.Active {
						t2ActiveCounts[row.ReferredBy]++
					}
				} else {
					log.Panic("wrong tier type")
				}
			}
		}
	}

	return referralsUserKeys, t1Referrals, t2Referrals, t1ActiveCounts, t2ActiveCounts, nil
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
				results = append(results, model.SerializedUsersKey(val.(string)))
			}
		}

		return results, nil
	}
}

func (m *miner) gatherReferralsInformation(ctx context.Context, users []*user) (history *historyData, err error) {
	if len(users) == 0 {
		return nil, nil
	}
	needToBeRecalculatedUsers := make(map[string]struct{}, len(users))
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
	referralsUserKeys, t1Referrals, t2Referrals, t1ActiveCounts, t2ActiveCounts, err := m.collectTiers(ctx, needToBeRecalculatedUsers)
	if err != nil {
		return nil, errors.Wrap(err, "can't get active users for users")
	}
	if len(t1Referrals) == 0 && len(t2Referrals) == 0 {
		return nil, nil
	}
	internalIDKeys, err := getInternalIDs(ctx, m.db, referralsUserKeys...)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get internal ids for:%#v", referralsUserKeys)
	}
	referrals := make([]*recalculateReferral, 0)
	if err := storage.Bind[recalculateReferral](ctx, m.db, internalIDKeys, &referrals); err != nil {
		return nil, errors.Wrapf(err, "failed to get referrals for:%v", users)
	}
	referralsCollection := make(map[string]*recalculateReferral, len(referrals))
	for _, ref := range referrals {
		referralsCollection[ref.UserID] = ref
	}

	return &historyData{
		Referrals:                 referralsCollection,
		NeedToBeRecalculatedUsers: needToBeRecalculatedUsers,
		T1Referrals:               t1Referrals,
		T2Referrals:               t2Referrals,
		T1ActiveCounts:            t1ActiveCounts,
		T2ActiveCounts:            t2ActiveCounts,
	}, nil
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

func (m *miner) insertBalanceRecalculationDryRunBatch(ctx context.Context, infos []*balanceRecalculationDryRun) error {
	if len(infos) == 0 {
		return nil
	}
	paramCounter := 1
	params := []any{}
	sqlParams := []string{}
	for _, info := range infos {
		sqlParams = append(sqlParams, fmt.Sprintf("($%v, $%v, $%v, $%v, $%v)", paramCounter, paramCounter+1, paramCounter+2, paramCounter+3, paramCounter+4))
		paramCounter += 5
		params = append(params, info.T1BalanceDiff, info.T2BalanceDiff, info.T1ActiveCountsDiff, info.T2ActiveCountsDiff, info.UserID)
	}

	sql := fmt.Sprintf(`INSERT INTO balance_recalculation_dry_run(
							diff_t1_balance,
							diff_t2_balance,
							diff_t1_active_counts,
							diff_t2_active_counts,
							user_id
						)
					VALUES %v
					ON CONFLICT(user_id) DO NOTHING`, strings.Join(sqlParams, ","))
	_, err := storagePG.Exec(ctx, m.dbPG, sql, params...)

	return errors.Wrap(err, "failed to insert dry run info")
}
