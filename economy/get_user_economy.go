// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ICE-Blockchain/wintr/log"
)

func (u *userEconomyRepository) GetUserEconomy(ctx context.Context, userID UserID, ownEconomy bool) (*UserEconomy, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get user economy failed because context failed")
	}
	if ownEconomy {
		return u.getOwnUserEconomy(ctx, userID)
	}

	return u.getAnotherUserEconomy(ctx, userID)
}

func (u *userEconomyRepository) getOwnUserEconomy(ctx context.Context, userID UserID) (*UserEconomy, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get user economy failed because context failed")
	}
	var result []*userEconomy
	params := map[string]interface{}{
		"userId":             userID,
		"now":                time.Now().UTC().UnixNano(),
		"inactivityDeadline": time.Duration(cfg.InactivityHoursDeadline) * time.Hour,
	}
	if err := u.db.PrepareExecuteTyped(getUserEconomySQL(), params, &result); err != nil {
		return nil, errors.Wrapf(err, "failed to get user economy for userID:%v", userID)
	}
	if len(result) == 0 {
		return nil, errors.Wrapf(ErrNotFound, "no user economy found for id:%v", userID)
	}

	return result[0].toUserEconomy(), nil
}

func getUserEconomySQL() string {
	t1ActiveUsersCountSQL := getActiveUsersSQL(getActiveUsersT1Condition())
	t2ActiveUsersCountSQL := getActiveUsersSQL(getActiveUsersT2Condition())
	t1EarningsSumSQL := getTiersEarningsSumSQL(t1ReferralsSpace())
	t2EarningsSumSQL := getTiersEarningsSumSQL(t2ReferralsSpace())

	return fmt.Sprintf(`SELECT
		ue.user_id,
		ue.profile_picture_url,
		(%[1]v) as adoptions,
		ue.balance,
		ue.staking_percentage,
		ue.hash_code,
		ue.last_mining_started_at,
		ue.staking_years,
		ue.created_at,
		ue.updated_at,
		ue.balance_updated_at,
		(%[2]v) as t1_count,
		(%[3]v) as t2_count,
		(%[4]v) as global_rank,
		(%[5]v) as t1_earnings_sum,
		(%[6]v) as t2_earnings_sum,
		(SELECT value FROM %[7]v WHERE key = 'TOTAL_USERS') as current_total_users
	FROM %[8]v ue INDEXED BY "pk_unnamed_%[8]v_1"
	WHERE ue.user_id = :userId`,
		getAdoptionsSQL(), t1ActiveUsersCountSQL, t2ActiveUsersCountSQL,
		getGlobalRankSQL(), t1EarningsSumSQL, t2EarningsSumSQL, totalUsersSpace(), userEconomySpace())
}

func (u *userEconomyRepository) getAnotherUserEconomy(ctx context.Context, userID UserID) (*UserEconomy, error) {
	// For now we return the same as for own user. It will be replaced later.
	return u.getOwnUserEconomy(ctx, userID)
}

func getGlobalRankSQL() string {
	return fmt.Sprintf(`SELECT count(1) - 1
			FROM %[1]v
			WHERE balance >= (SELECT ue.balance 
					FROM %[1]v ue INDEXED BY "pk_unnamed_%[1]v_1"
					WHERE ue.user_id = :userId) AND user_id != :userId`, userEconomySpace())
}

func getActiveUsersSQL(sqlCondition string) string {
	return fmt.Sprintf(`SELECT count(1)
			FROM %[1]v t INDEXED BY "pk_unnamed_%[1]v_1"
				JOIN %[2]v ue INDEXED BY "pk_unnamed_%[2]v_1"
					ON t.user_id = ue.user_id
			WHERE %[3]v AND :now - ue.last_mining_started_at < :inactivityDeadline`, t1ReferralsSpace(), userEconomySpace(), sqlCondition)
}

func getActiveUsersT1Condition() string {
	return "ue.user_id = :userId"
}

func getActiveUsersT2Condition() string {
	return fmt.Sprintf(`ue.user_id IN (SELECT referral_user_id
			FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1"
			WHERE user_id = :userId)`, t1ReferralsSpace())
}

func getAdoptionsSQL() string {
	return fmt.Sprintf(`SELECT
				GROUP_CONCAT(CAST(total_users as string) || ':' || CAST(base_hourly_mining_rate as string))
			FROM %[1]v
			ORDER BY total_users ASC`, adoptionSpace())
}

func getTiersEarningsSumSQL(table string) string {
	return fmt.Sprintf(`SELECT SUM(earnings)
			FROM %v INDEXED BY "pk_unnamed_%[1]v_1"
			WHERE user_id = :userId`, table)
}

func parseAdoptions(adoptions string, currentTotalUsers uint64) (map[uint64]float64, float64) {
	a := strings.Split(adoptions, ",")
	res := make(map[uint64]float64, len(a))
	var baseHourlyMiningRate float64

	for _, adoption := range a {
		parts := strings.Split(adoption, ":")
		totalUsers, err := strconv.ParseUint(parts[0], digitBase, digitBitSize)
		if err != nil {
			log.Error(err, "can't parse rate uint for adoption:%v", parts[0])

			continue
		}
		rate, err := strconv.ParseFloat(parts[1], digitBitSize)
		if err != nil {
			log.Error(err, "can't parse baseHourlyMiningrate float64 %[1]v for adoption with total users:%[2]v", parts[1], parts[0])

			continue
		}

		res[totalUsers] = rate
		if currentTotalUsers <= totalUsers {
			baseHourlyMiningRate = rate
		}
	}

	return res, baseHourlyMiningRate
}

func (u *userEconomy) toUserEconomy() *UserEconomy {
	adoptions, baseHourlyMiningRate := parseAdoptions(u.Adoptions, u.CurrentTotalUsers)

	return &UserEconomy{
		Balance: Balance{
			Total: u.Balance,
			Referrals: ReferralBalance{
				T1: u.T1EarningsSum,
				T2: u.T2EarningsSum,
			},
		},
		HourlyMiningRate:    baseHourlyMiningRate * (float64(u.T1Count)*cfg.Rates.Tier1 + float64(u.T2Count)*cfg.Rates.Tier2 + 1),
		GlobalRank:          u.GlobalRank,
		CurrentTotalUsers:   u.CurrentTotalUsers,
		Adoption:            adoptions,
		LastMiningStartedAt: time.Unix(int64(u.LastMiningStartedAt), 0),
		Staking: Staking{
			Years:      u.StakingYears,
			Percentage: u.StakingPercentage,
		},
	}
}

func t1ReferralsSpace() string {
	return "T1_REFERRAL_EARNINGS"
}

func t2ReferralsSpace() string {
	return "T2_REFERRAL_EARNINGS"
}

func userEconomySpace() string {
	return "USER_ECONOMY"
}

func adoptionSpace() string {
	return "ADOPTION"
}

func totalUsersSpace() string {
	return "TOTAL_USERS"
}
