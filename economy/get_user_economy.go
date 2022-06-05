// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	tm "time"

	"cosmossdk.io/math"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (e *economy) GetUserEconomy(ctx context.Context, userID UserID, ownEconomy bool) (*UserEconomy, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get user economy failed because context failed")
	}
	if ownEconomy {
		return e.getOwnUserEconomy(ctx, userID)
	}

	return e.getAnotherUserEconomy(ctx, userID)
}

func (e *economy) getOwnUserEconomy(ctx context.Context, userID UserID) (*UserEconomy, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get user economy failed because context failed")
	}
	var result []*userEconomySummary
	params := map[string]interface{}{
		"userId":             userID,
		"now":                time.Now().UnixNano(),
		"inactivityDeadline": tm.Duration(cfg.InactivityHoursDeadline) * tm.Hour,
	}
	if err := e.db.PrepareExecuteTyped(getUserEconomySQL(), params, &result); err != nil {
		return nil, errors.Wrapf(err, "failed to get user economy for userID:%v", userID)
	}
	if len(result) == 0 {
		return nil, errors.Wrapf(ErrNotFound, "no user economy found for id:%v", userID)
	}

	return result[0].toUserEconomySummary(), nil
}

//nolint:funlen // Because this is SQL, no sense to split it to different functions.
func getUserEconomySQL() string {
	return fmt.Sprintf(`SELECT
		ue.last_mining_started_at,
		s.updated_at AS staking_balance_updated_at,
		b.amount AS balance,
		(SELECT amount FROM balances WHERE user_id = :userId AND type = 'staking') AS staking_balance,
		(SELECT base_hourly_mining_rate
			FROM adoption
			WHERE active = true
		) AS base_hourly_mining_rate,
		(SELECT amount
			FROM balances b INDEXED BY "pk_unnamed_BALANCES_1"
				WHERE type = 't0_referral_standard_earnings' AND b.user_id = :userId
		) AS t0_amount,
		(SELECT amount
			FROM balances b INDEXED BY "pk_unnamed_BALANCES_1"
				WHERE type = 't1_referral_standard_earnings' AND b.user_id = :userId
		) AS t1_amount,
		(SELECT amount
			FROM balances b INDEXED BY "pk_unnamed_BALANCES_1"
				WHERE type = 't2_referral_standard_earnings' AND b.user_id = :userId
		) AS t2_amount,
		ue.user_id,
		ue.username,
		ue.profile_picture_url,
		(SELECT
			GROUP_CONCAT(CAST(total_active_users AS string) || ':' || base_hourly_mining_rate || CAST(active AS string))
			FROM adoption
			ORDER BY total_active_users ASC
		) AS adoptions,
		ue.hash_code,
		(SELECT count(1)
			FROM balances b INDEXED BY "pk_unnamed_BALANCES_1"
				JOIN user_economy ue INDEXED BY "pk_unnamed_USER_ECONOMY_1"
					ON b.user_id = ue.user_id
			WHERE b.user_id = :userId AND :now - ue.last_mining_started_at < :inactivityDeadline
				  AND POSITION('t0_referral_standard_earnings~', lower(b.type)) || ue.user_id
		) AS t0_count,
		(SELECT count(1)
			FROM balances b INDEXED BY "pk_unnamed_BALANCES_1"
				JOIN user_economy ue INDEXED BY "pk_unnamed_USER_ECONOMY_1"
					ON b.user_id = ue.user_id
			WHERE b.user_id = :userId AND :now - ue.last_mining_started_at < :inactivityDeadline
				   AND POSITION('t1_referral_standard_earnings~', lower(b.type)) != 0
		) AS t1_count,
		(SELECT count(1)
			FROM balances b INDEXED BY "pk_unnamed_BALANCES_1"
				JOIN user_economy ue INDEXED BY "pk_unnamed_USER_ECONOMY_1"
					ON b.user_id = ue.user_id
			WHERE b.user_id = :userId AND :now - ue.last_mining_started_at < :inactivityDeadline
				   AND POSITION('t2_referral_standard_earnings~', lower(b.type))
		) AS t2_count,
		(%[1]v) AS global_rank,
		s.percentage AS staking_percentage_allocation,
		s.years AS staking_years,
		(SELECT value
				FROM global
				WHERE key = 'TOTAL_USERS'
		) AS current_total_users,
		(SELECT sb.percentage
				FROM staking_bonus sb
					JOIN staking s ON sb.years = s.years AND s.user_id = :userId
		) AS staking_percentage_bonus
	FROM user_economy ue INDEXED BY "pk_unnamed_USER_ECONOMY_1"
		INNER JOIN balances b
			ON b.user_id = ue.user_id AND type = 'standard'
		LEFT JOIN staking s
			ON s.user_id = ue.user_id
	WHERE ue.user_id = :userId`,
		getGlobalRankSQL(),
	)
}

func (e *economy) getAnotherUserEconomy(ctx context.Context, userID UserID) (*UserEconomy, error) {
	// For now we return the same as for own user. It will be replaced later.
	return e.getOwnUserEconomy(ctx, userID)
}

func getGlobalRankSQL() string {
	return `
		SELECT count(1) - 1
		FROM balances b_cmp
		WHERE CASE
			WHEN b_cmp.amount_w3 == b.amount_w3
			THEN (CASE
					WHEN b_cmp.amount_w2 == b.amount_w2
						THEN (CASE
								WHEN b_cmp.amount_w1 == b.amount_w1
								THEN (b_cmp.amount_w0 >= b.amount_w0)
								ELSE b_cmp.amount_w1 > b.amount_w1
							END)
					ELSE b_cmp.amount_w2 > b.amount_w2
					END)
				ELSE b_cmp.amount_w3 > b.amount_w3
			END
	`
}

func parseAdoptions(adoptions string) map[uint64]*coin.ICEFlake {
	a := strings.Split(adoptions, ",")
	res := make(map[uint64]*coin.ICEFlake, len(a))

	for _, adoption := range a {
		parts := strings.Split(adoption, ":")
		totalUsers, err := strconv.ParseUint(parts[0], base10, bitSize64)
		log.Panic(errors.Wrapf(err, "can't parse rate uint for adoption:%v", parts[0]))

		res[totalUsers] = coin.UnsafeNewAmount(parts[1])
	}

	return res
}

func (u *userEconomySummary) toUserEconomySummary() *UserEconomy {
	adoptions := parseAdoptions(u.Adoptions)
	hmr := u.calculateHourlyMiningRate()

	return &UserEconomy{
		LastMiningStartedAt: u.LastMiningStartedAt,
		HourlyMiningRate:    hmr,
		Adoption:            adoptions,
		Balance: Balance{
			Total: u.Balance,
			Referrals: ReferralBalance{
				T1: u.T1Amount,
				T2: u.T2Amount,
			},
		},
		CurrentTotalUsers: u.CurrentTotalUsers,
		Staking: Staking{
			Years:      u.StakingYears,
			Percentage: u.StakingPercentageAllocation,
		},
		GlobalRank: u.GlobalRank,
	}
}

func (u *userEconomySummary) calculateHourlyMiningRate() *coin.ICEFlake {
	tierCountPart := (u.T0Count*cfg.Rates.Tier0 + u.T1Count*cfg.Rates.Tier1 + u.T2Count*cfg.Rates.Tier2 + percentage100)
	hmr := u.BaseHourlyMiningRate.Mul(math.NewUint(tierCountPart)).Quo(math.NewUint(percentage100))

	return &coin.ICEFlake{Uint: hmr}
}
