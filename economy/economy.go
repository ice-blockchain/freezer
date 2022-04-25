// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	appCfg "github.com/ICE-Blockchain/wintr/config"
	"github.com/ICE-Blockchain/wintr/connectors/storage"
	"github.com/ICE-Blockchain/wintr/log"
)

func New(ctx context.Context, cancel context.CancelFunc) Repository {
	db := storage.MustConnect(ctx, cancel, ddl, applicationYamlKey)
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	return &repository{
		close: db.Close,
		db:    db,
	}
}

func (r *repository) Close() error {
	log.Info("closing economy repository...")

	return errors.Wrap(r.close(), "closing economy repository failed")
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor {
	//nolint:nolintlint // TODO implement me
	return nil
}

func (p *processor) Close() error {
	//nolint:nolintlint // TODO implement me.

	return nil
}

func (p *processor) CheckHealth(ctx context.Context) error {
	//nolint:nolintlint // TODO implement me.

	return nil
}

func (r *repository) GetUserEconomy(ctx context.Context, userID UserID, ownEconomy bool) (*UserEconomy, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get user economy failed because context failed")
	}
	adoption, err := r.getAdoptions(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select adoption information")
	}
	if ownEconomy {
		return r.getOwnUserEconomy(ctx, userID, adoption)
	}

	return r.getAnotherUserEconomy(ctx, userID, adoption)
}

func (r *repository) getOwnUserEconomy(ctx context.Context, userID UserID, adoptions map[TotalUsers]BaseHourlyMiningRate) (*UserEconomy, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get user economy failed because context failed")
	}
	var result []*userEconomy
	params := map[string]interface{}{
		"userId":             userID,
		"now":                time.Now().UTC().UnixNano(),
		"inactivityDeadline": time.Duration(cfg.InactivityHoursDeadline) * time.Hour,
	}
	if err := r.db.PrepareExecuteTyped(getUserEconomySQL(), params, &result); err != nil {
		return nil, errors.Wrapf(err, "failed to get user economy for userID:%v", userID)
	}
	if len(result) == 0 {
		return nil, errors.Wrapf(ErrNotFound, "no user economy found for id:%v", userID)
	}

	return result[0].toUserEconomy(adoptions), nil
}

func getUserEconomySQL() string {
	t1ActiveUsersSQL := getActiveUsersSQL(getActiveUsersT1Condition())
	t2ActiveUsersSQL := getActiveUsersSQL(getActiveUsersT2Condition())
	t1EarningsSumSQL := getTiersEarningsSumSQL(t1ReferralsSpace())
	t2EarningsSumSQL := getTiersEarningsSumSQL(t2ReferralsSpace())

	return fmt.Sprintf(`SELECT ue.user_id,
			ue.profile_picture_url,
			ue.balance,
			ue.staking_percentage,
			ue.hash_code,
			ue.last_mining_started_at,
			ue.staking_years,
			ue.created_at,
			ue.updated_at,
			ue.balance_updated_at,
			(%[1]v) as t1_count,
			(%[2]v) as t2_count,
			(SELECT count(1) + 1
					FROM %[3]v
					WHERE balance >= (SELECT ue.balance 
											FROM %[3]v ue INDEXED BY "pk_unnamed_%[3]v_1"
											WHERE ue.user_id = :userId) AND user_id != :userId) as global_rank,
			(%[4]v) as t1_earnings_sum,
			(%[5]v) as t2_earnings_sum,
			(SELECT count(1) FROM %[3]v) as current_total_users,
			(%[6]v) as base_hourly_mining_rate
		FROM %[3]v ue INDEXED BY "pk_unnamed_%[3]v_1"
		WHERE ue.user_id = :userId`,
		t1ActiveUsersSQL, t2ActiveUsersSQL, userEconomySpace(), t1EarningsSumSQL, t2EarningsSumSQL, getBaseHourlyMiningRateSQL())
}

func (r *repository) getAnotherUserEconomy(ctx context.Context, userID UserID, adoptions map[TotalUsers]BaseHourlyMiningRate) (*UserEconomy, error) {
	// For now we return the same as for own user. It will be replaced later.
	return r.getOwnUserEconomy(ctx, userID, adoptions)
}

func getActiveUsersSQL(sqlCondition string) string {
	return fmt.Sprintf(`SELECT count(1)
				FROM %[1]v t1 INDEXED BY "pk_unnamed_%[1]v_1"
					JOIN %[2]v ue INDEXED BY "pk_unnamed_%[2]v_1"
						ON t1.user_id = ue.user_id
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

func getBaseHourlyMiningRateSQL() string {
	return fmt.Sprintf(`SELECT base_hourly_mining_rate
						FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1"
						WHERE total_users >= (SELECT count(1) FROM %[2]v)
						ORDER BY total_users ASC LIMIT 1`, adoptionSpace(), userEconomySpace())
}

func (r *repository) getAdoptions(ctx context.Context) (map[TotalUsers]BaseHourlyMiningRate, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get adoptions failed because context failed")
	}
	var res []*adoption
	sql := fmt.Sprintf(`SELECT * FROM %v`, adoptionSpace())
	if err := r.db.PrepareExecuteTyped(sql, nil, &res); err != nil {
		return nil, errors.Wrapf(err, "failed to get adoptions")
	}
	if len(res) == 0 {
		return nil, errors.New("no adoptions in the database")
	}

	return toAdoptionsMap(res), nil
}

func toAdoptionsMap(a []*adoption) map[TotalUsers]BaseHourlyMiningRate {
	res := make(map[TotalUsers]BaseHourlyMiningRate, len(a))
	for _, r := range a {
		res[r.TotalUsers] = r.BaseHourlyMiningRate
	}

	return res
}

func getTiersEarningsSumSQL(table string) string {
	return fmt.Sprintf(`SELECT
							SUM(earnings)
						FROM %v INDEXED BY "pk_unnamed_%[1]v_1"
						WHERE user_id = :userId`, table)
}

func (u *userEconomy) toUserEconomy(adoption map[TotalUsers]BaseHourlyMiningRate) *UserEconomy {
	return &UserEconomy{
		Balance: Balance{
			Total: u.Balance,
			Referrals: ReferralBalance{
				T1: u.T1EarningsSum,
				T2: u.T2EarningsSum,
			},
		},
		HourlyMiningRate:    u.BaseHourlyMiningRate * (float64(u.T1Count)*cfg.Rates.Tier1 + float64(u.T2Count)*cfg.Rates.Tier2), // * 100?
		GlobalRank:          u.GlobalRank,
		CurrentTotalUsers:   u.CurrentTotalUsers,
		Adoption:            adoption,
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
