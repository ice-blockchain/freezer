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
		close:                 db.Close,
		UserEconomyRepository: &economy{db: db},
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

func (p *economy) GetUserEconomy(ctx context.Context, userID UserID, ownEconomy bool) (*UserEconomy, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get user economy failed because context failed")
	}
	adoption, err := p.getAdoptions(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select adoption information")
	}
	if ownEconomy {
		return p.getOwnUserEconomy(ctx, userID, adoption)
	}

	return p.getAnotherUserEconomy(ctx, userID, adoption)
}

func (p *economy) getOwnUserEconomy(ctx context.Context, userID UserID, adoptions map[TotalUsers]BaseHourlyMiningRate) (*UserEconomy, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get user economy failed because context failed")
	}
	var result []*userEconomy
	params := map[string]interface{}{
		"userId":             userID,
		"now":                time.Now().UTC().UnixNano(),
		"inactivityDeadline": time.Duration(cfg.InactivityHoursDeadline) * time.Hour,
	}
	if err := p.db.PrepareExecuteTyped(getUserEconomySQL(), params, &result); err != nil {
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
							   (SELECT CASE WHEN count(1) = 0 THEN 1 ELSE count(1) END
										FROM %[3]v
										WHERE balance > (SELECT ue.balance 
																FROM %[3]v ue INDEXED BY "pk_unnamed_%[3]v_1"
																WHERE ue.user_id = :userId)) as global_rank,
							   (%[4]v) as t1_earnings_sum,
							   (%[5]v) as t2_earnings_sum,
							   (SELECT count(1) FROM users) as current_total_users
							FROM %[3]v ue INDEXED BY "pk_unnamed_%[3]v_1"
							WHERE ue.user_id = :userId`, t1ActiveUsersSQL, t2ActiveUsersSQL, userEconomySpace(), t1EarningsSumSQL, t2EarningsSumSQL)
}

func (p *economy) getAnotherUserEconomy(ctx context.Context, userID UserID, adoptions map[TotalUsers]BaseHourlyMiningRate) (*UserEconomy, error) {
	// For now we return the same as for own user. It will be replaced later.
	return p.getOwnUserEconomy(ctx, userID, adoptions)
}

func getActiveUsersSQL(sqlCondition string) string {
	return fmt.Sprintf(`SELECT count(1)
				FROM %[1]v u INDEXED BY "pk_unnamed_%[1]v_1"
					JOIN %[2]v ue INDEXED BY "pk_unnamed_%[2]v_1"
						ON u.ID = ue.user_id
				WHERE :now - ue.last_mining_started_at < :inactivityDeadline %[3]v`, usersSpace(), userEconomySpace(), sqlCondition)
}

func getActiveUsersT1Condition() string {
	return "AND u.referred_by = :userId"
}

func getActiveUsersT2Condition() string {
	return fmt.Sprintf(`AND u.id IN (SELECT id
										FROM %[1]v
										WHERE referred_by = :userId)`, usersSpace())
}

func (p *economy) getAdoptions(ctx context.Context) (map[TotalUsers]BaseHourlyMiningRate, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get adoptions failed because context failed")
	}
	var res []*adoption
	sql := fmt.Sprintf(`SELECT * FROM %v`, adoptionSpace())
	if err := p.db.PrepareExecuteTyped(sql, nil, &res); err != nil {
		return nil, errors.Wrapf(err, "failed to get adoptions")
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

func getAdoptionBaseRate(adoption map[TotalUsers]BaseHourlyMiningRate, currentTotalUsers uint64) float64 {
	baseHourlyMiningRate := cfg.DefaultRates.BaseHourlyMiningRate
	for totalUsers, rate := range adoption {
		if currentTotalUsers <= totalUsers {
			baseHourlyMiningRate = rate
		}
	}

	return baseHourlyMiningRate
}

func (u *userEconomy) toUserEconomy(adoption map[TotalUsers]BaseHourlyMiningRate) *UserEconomy {
	adoptionBaseRate := getAdoptionBaseRate(adoption, u.CurrentTotalUsers)

	return &UserEconomy{
		Balance: Balance{
			Total: u.Balance,
			Referrals: ReferralBalance{
				T1: u.T1EarningsSum,
				T2: u.T2EarningsSum,
			},
		},
		HourlyMiningRate:    adoptionBaseRate * (float64(u.T1Count)*cfg.DefaultRates.Tier1Rate + float64(u.T2Count)*cfg.DefaultRates.Tier2Rate),
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

func usersSpace() string {
	return "USERS"
}

func userEconomySpace() string {
	return "USER_ECONOMY"
}

func adoptionSpace() string {
	return "ADOPTION"
}
