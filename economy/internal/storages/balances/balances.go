package balances

import (
	"context"
	"strconv"
	"strings"
	"sync"
	tm "time"

	"cosmossdk.io/math"
	"github.com/framey-io/go-tarantool"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func New(db tarantool.Connector) messagebroker.Processor {
	cfg := new(config)
	appCfg.MustLoadFromKey(applicationYamlKey, cfg)

	return &balanceSource{db: db, cfg: cfg}
}

func (b *balanceSource) Process(ctx context.Context, m *messagebroker.Message) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}

	economies, err := b.getUserEconomies(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get user economies")
	}

	return errors.Wrap(b.calculateUpdateBalances(ctx, economies), "can't execute balance calculation and update")
}

func (b *balanceSource) getUserEconomies(ctx context.Context) ([]*userEconomy, error) {
	params := map[string]interface{}{
		"chunks":             b.cfg.MessageBroker.Topics[3].Partitions,
		"chunkIndex":         b.cfg.MessageBroker.Topics[3].Partition,
		"inactivityDeadline": (tm.Duration(b.cfg.InactivityHoursDeadline) * tm.Hour).Nanoseconds(),
		"now":                uint64(time.Now().UnixNano()),
	}
	sql := getBalancesSQL()

	offset := 0
	var output []*userEconomy
	for ctx.Err() == nil {
		params["offset"] = offset

		var res []*userEconomy
		err := b.db.PrepareExecuteTyped(sql, params, &res)
		if err != nil {
			return nil, errors.Wrapf(err, "db execute getting user economy for (chunk:%v, chunkIndex:%v) failed",
				b.cfg.MessageBroker.Topics[3].Partitions, b.cfg.MessageBroker.Topics[3].Partition)
		}

		if len(res) == 0 {
			break
		}
		offset += len(res)
		output = append(output, res...)
	}

	return output, nil
}

//nolint:funlen // Because this is SQL, no sense to split it to different functions. It will reduce readability.
func getBalancesSQL() string {
	return `SELECT
				ue.user_id,
				(SELECT base_hourly_mining_rate
					FROM adoption
					WHERE active = true
				) AS base_hourly_mining_rate,
				(SELECT GROUP_CONCAT(CAST(b.percentage AS string) || ';' || CAST(s.percentage AS string))
					FROM staking_bonus b
						JOIN staking s
							ON b.years = s.years AND s.user_id = ue.user_id
				) AS staking_info,
				(SELECT GROUP_CONCAT(b.type || ';' || b.amount || ';' || CAST(b.updated_at AS string))
					FROM balances b
					LEFT JOIN user_economy bue
						ON bue.user_id = SUBSTR(b.type, POSITION('~', b.type) + 1)
					WHERE b.user_id = ue.user_id AND (:now - ue.last_mining_started_at < :inactivityDeadline OR ue.last_mining_started_at IS NULL) 
				) AS balances,
				(SELECT count(1)
					FROM balances t0 
					JOIN user_economy u
						ON t0.user_id = ue.user_id AND t0.type = 't0_referral_standard_earnings~' || u.user_id
						   AND :now - u.last_mining_started_at < :inactivityDeadline
				) AS t0_referrals,
				(SELECT count(1)
					FROM balances t1
					JOIN user_economy u
						ON t1.user_id = ue.user_id AND t1.type = 't1_referral_standard_earnings~' || u.user_id
						AND :now - u.last_mining_started_at < :inactivityDeadline
				) AS t1_referrals,
				(SELECT count(1)
					FROM balances t2
					JOIN user_economy u
						ON t2.user_id = ue.user_id AND t2.type = 't2_referral_standard_earnings~' || u.user_id
						AND :now - u.last_mining_started_at < :inactivityDeadline
				) AS t2_referrals
			FROM user_economy ue	
			WHERE ue.hash_code % :chunks == :chunkIndex -- update only active users? AND :now - ue.last_mining_started_at < :inactivityDeadline
			LIMIT 1000 OFFSET :offset`
}

func (b *balanceSource) calculateUpdateBalances(ctx context.Context, ue []*userEconomy) error {
	var errs []error
	wg := new(sync.WaitGroup)
	for _, u := range ue {
		go func(u *userEconomy) {
			wg.Add(1)

			si := u.parseUserStakingInfo()
			balances := u.parseUserBalances()
			errs = append(errs, errors.Wrap(b.handleGeneralBalances(ctx, u, balances, si), "can't handle general balances"),
				errors.Wrap(b.handleUsersBalances(ctx, u, balances, si), "can't handle users balances"))

			wg.Done()
		}(u)
	}
	wg.Wait()

	return errors.Wrap(multiErr(errs), "can't calculate and update balances")
}

func (b *balanceSource) handleGeneralBalances(ctx context.Context, u *userEconomy, balances map[string]balance, si *stakingInfo) error {
	var errs []error
	normalHourlyMiningRate, stakedHourlyMiningRate := u.calculateRates(b.cfg, si)

	balanceStandard := coin.UnsafeNew(u.calculateStandardBalance(balances["standard"], normalHourlyMiningRate).String())
	balanceStaking := coin.UnsafeNew(u.calculateStakingBalance(balances["staking"], stakedHourlyMiningRate).String())
	balanceTotal := coin.UnsafeNew(u.calculateTotalBalance(balanceStandard.Amount, balanceStaking.Amount).String())

	errs = append(errs, errors.Wrap(b.updateBalance(ctx, u.UserID, "standard", balanceStandard), "can't update standard balance"),
		errors.Wrap(b.updateBalance(ctx, u.UserID, "staking", balanceStaking), "can't update staking balance"),
		errors.Wrap(b.updateBalance(ctx, u.UserID, "total", balanceTotal), "can't update total balance"))

	return errors.Wrap(multiErr(errs), "can't update standard, staking, total balances")
}

func (b *balanceSource) handleUsersBalances(ctx context.Context, u *userEconomy, balances map[string]balance, si *stakingInfo) error {
	coeffs := getDividerCoefficients()
	sumDict := initializeSumDictionary(balances)
	errs := make([]error, 0)
	referralCounts := u.getReferralsCountDictionary()
	generalStandard, generalStaking := u.calculateGeneralFormulaParts(si)

	for balanceType, value := range balances {
		if !strings.Contains(balanceType, "~") {
			continue
		}
		generalFormulaPart := generalStandard
		if strings.Contains(balanceType, "staking") {
			generalFormulaPart = generalStaking
		}

		typeParts := strings.Split(balanceType, "~")
		earnings := value.calculateReferralEarnings(u.UserID, generalFormulaPart, referralCounts[typeParts[0]], coeffs[typeParts[0]])
		sumDict[typeParts[0]] = &coin.ICEFlake{Uint: sumDict[typeParts[0]].Add(value.Amount.Add(earnings))}

		errs = append(errs, errors.Wrapf(b.updateBalance(ctx, u.UserID, balanceType, coin.UnsafeNew(earnings.String())),
			"can't update %v balance", balanceType))
	}
	for balanceType, amount := range sumDict {
		if amount != nil { // Due to user can not have all types of tiers atm, so some ICEFlake values can be nil.
			errs = append(errs, errors.Wrapf(b.updateBalance(ctx, u.UserID, balanceType, coin.UnsafeNew(amount.String())),
				"can't update %v", balanceType))
		}
	}

	return errors.Wrap(multiErr(errs), "can't update users balances")
}

func (u *userEconomy) calculateGeneralFormulaParts(si *stakingInfo) (*coin.ICEFlake, *coin.ICEFlake) {
	standard := coin.ICEFlake{Uint: math.NewUint(percentage100 - si.Allocation).Mul(u.BaseHourlyMiningRate.Uint)}
	staking := coin.ICEFlake{Uint: math.NewUint(si.Allocation * si.Bonus).Mul(u.BaseHourlyMiningRate.Uint)}

	return &standard, &staking
}

func (b *balance) calculateReferralEarnings(userID UserID, generalFormulaPart *coin.ICEFlake, referralsCount, divider uint64) math.Uint {
	elapsedNanoseconds := uint64(time.Now().UnixNano()) - b.UpdatedAt
	earnings := generalFormulaPart.MulUint64(referralsCount).MulUint64(elapsedNanoseconds).QuoUint64(divider)

	return earnings
}

func (u *userEconomy) calculateStandardBalance(b balance, normalHourlyMiningRate *coin.ICEFlake) *coin.ICEFlake {
	elapsedNanoseconds := math.NewUint(uint64(time.Now().UnixNano()) - b.UpdatedAt)
	//nolint:gomnd // This is not a magic number, this is the divider.
	balanceStandardAddition := normalHourlyMiningRate.Mul(elapsedNanoseconds).QuoUint64(3600000000000)
	balanceStandard := b.Amount.Add(balanceStandardAddition)

	return &coin.ICEFlake{Uint: balanceStandard}
}

func (u *userEconomy) calculateStakingBalance(b balance, stakedHourlyMiningRate *coin.ICEFlake) *coin.ICEFlake {
	elapsedNanoseconds := math.NewUint(uint64(time.Now().UnixNano()) - b.UpdatedAt)
	//nolint:gomnd // This is not a magic number, this is the divider.
	balanceStaking := b.Amount.Add(stakedHourlyMiningRate.Mul(elapsedNanoseconds).QuoUint64(3600000000000))

	return &coin.ICEFlake{Uint: balanceStaking}
}

func (u *userEconomy) calculateTotalBalance(standard, staking *coin.ICEFlake) *coin.ICEFlake {
	return &coin.ICEFlake{Uint: standard.Add(staking.Uint)}
}

func (u *userEconomy) calculateRates(cfg *config, si *stakingInfo) (*coin.ICEFlake, *coin.ICEFlake) {
	multiplier := u.T0Referrals*cfg.Rates.Tier0 + u.T1Referrals*cfg.Rates.Tier1 + u.T2Referrals*cfg.Rates.Tier2 + percentage100

	hourlyMiningRate := u.BaseHourlyMiningRate.MulUint64(multiplier).QuoUint64(percentage100)
	normalHourlyMiningRate := math.NewUint(percentage100 - si.Allocation).Mul(hourlyMiningRate).QuoUint64(percentage100)

	//nolint:gomnd // This is not magic number, this is the divider.
	stakingRate := math.NewUint(si.Bonus).Mul(hourlyMiningRate).Mul(math.NewUint(si.Allocation)).QuoUint64(10000)

	return &coin.ICEFlake{Uint: normalHourlyMiningRate}, &coin.ICEFlake{Uint: stakingRate}
}

func (u *userEconomy) parseUserStakingInfo() *stakingInfo {
	if u.StakingInfo != "" {
		parts := strings.Split(u.StakingInfo, ";")

		bonus, err := strconv.ParseUint(parts[0], base10, bitSize64)
		log.Panic(errors.Wrapf(err, "can't parse uint64 value for bonus:%v", parts[0]))

		allocation, err := strconv.ParseUint(parts[1], base10, bitSize64)
		log.Panic(errors.Wrapf(err, "can't parse uint64 value for allocation:%v", parts[1]))

		return &stakingInfo{Bonus: bonus, Allocation: allocation}
	}

	return &stakingInfo{Bonus: 0, Allocation: 0}
}

func (u *userEconomy) parseUserBalances() map[string]balance {
	if u.Balances != "" {
		balances := strings.Split(u.Balances, ",")
		res := make(map[string]balance, len(balances))

		for _, val := range balances {
			parts := strings.Split(val, ";")
			updatedAt, err := strconv.ParseUint(parts[2], base10, bitSize64)
			log.Panic(errors.Wrapf(err, "can't parse uint64 value for balance updatedAt:%v", parts[2]))

			res[parts[0]] = balance{Amount: coin.UnsafeNewAmount(parts[1]), UpdatedAt: updatedAt}
		}

		return res
	}

	return nil
}

func (b *balanceSource) updateBalance(ctx context.Context, userID UserID, balanceType string, balance *coin.Coin) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	params := map[string]interface{}{
		"amount":    balance.Amount,
		"amountW0":  balance.AmountWord0,
		"amountW1":  balance.AmountWord1,
		"amountW2":  balance.AmountWord2,
		"amountW3":  balance.AmountWord3,
		"updatedAt": time.Now(),
		"userId":    userID,
		"type":      balanceType,
	}

	sql := `UPDATE balances SET
				   amount = :amount,
				   amount_w0 = :amountW0,
				   amount_w1 = :amountW1,
				   amount_w2 = :amountW2,
				   amount_w3 = :amountW3,
				   updated_at = :updatedAt
				WHERE user_id = :userId AND type = :type`

	query, err := b.db.PrepareExecute(sql, params)

	if err = storage.CheckSQLDMLErr(query, err); err != nil {
		return errors.Wrapf(err, "failed to update balances with userID:%v and type:%v", userID, balanceType)
	}

	return nil
}

func initializeSumDictionary(balances map[string]balance) map[string]*coin.ICEFlake {
	return map[string]*coin.ICEFlake{
		"t0_referral_standard_earnings": balances["t0_referral_standard_earnings"].Amount,
		"t1_referral_standard_earnings": balances["t1_referral_standard_earnings"].Amount,
		"t2_referral_standard_earnings": balances["t2_referral_standard_earnings"].Amount,
		"t0_referral_staking_earnings":  balances["t0_referral_staking_earnings"].Amount,
		"t1_referral_staking_earnings":  balances["t1_referral_staking_earnings"].Amount,
		"t2_referral_staking_earnings":  balances["t2_referral_staking_earnings"].Amount,
	}
}

func (u *userEconomy) getReferralsCountDictionary() map[string]uint64 {
	return map[string]uint64{
		"t0_referral_standard_earnings": u.T0Referrals,
		"t1_referral_standard_earnings": u.T1Referrals,
		"t2_referral_standard_earnings": u.T2Referrals,
		"t0_referral_staking_earnings":  u.T0Referrals,
		"t1_referral_staking_earnings":  u.T1Referrals,
		"t2_referral_staking_earnings":  u.T2Referrals,
	}
}

//nolint:gomnd // These are not magic number, they are the divider coefficients.
func getDividerCoefficients() map[string]uint64 {
	return map[string]uint64{
		"t0_referral_standard_earnings": 1440000000000000,
		"t1_referral_standard_earnings": 1440000000000000,
		"t2_referral_standard_earnings": 7200000000000000,
		"t0_referral_staking_earnings":  144000000000000000,
		"t1_referral_staking_earnings":  144000000000000000,
		"t2_referral_staking_earnings":  720000000000000000,
	}
}

func multiErr(errs []error) error {
	if len(errs) > 0 {
		nonNilErrs := make([]error, 0, len(errs))
		for _, e := range errs {
			if e != nil {
				nonNilErrs = append(nonNilErrs, e)
			}
		}
		if len(nonNilErrs) > 0 {
			return multierror.Append(nil, nonNilErrs...)
		}
	}

	return nil
}
