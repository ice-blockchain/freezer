package balances

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/framey-io/go-tarantool"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
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

	economies, err := b.getEconomy()
	if err != nil {
		return errors.Wrap(err, "can't get user economies")
	}

	b.calculate(economies)

	return nil
}

func (b *balanceSource) getEconomy() ([]*userEconomy, error) {
	params := map[string]interface{}{
		"chunks":             b.cfg.MessageBroker.Topics[3].Partitions,
		"chunkIndex":         b.cfg.MessageBroker.Topics[3].Partition,
		"inactivityDeadline": time.Duration(b.cfg.InactivityHoursDeadline) * time.Hour,
		"now":                uint64(time.Now().UnixNano()),
	}

	// TODO: batching?
	sql := `SELECT
				ue.user_id,
				ue.last_mining_started_at,
				(SELECT count(1)
					FROM balances t0 
					JOIN user_economy u
						ON t0.user_id = ue.user_id AND t0.type = 't0_referral_standard_earnings~' || u.user_id AND u.last_mining_started_at < :inactivityDeadline) AS t0_referrals,
				(SELECT count(1)
					FROM balances t1
					JOIN user_economy u
						ON t1.user_id = ue.user_id AND t1.type = 't1_referral_standard_earnings~' || u.user_id AND u.last_mining_started_at < :inactivityDeadline) AS t1_referrals,
				(SELECT count(1)
					FROM balances t2
					JOIN user_economy u
						ON t2.user_id = ue.user_id AND t2.type = 't2_referral_standard_earnings~' || u.user_id AND u.last_mining_started_at < :inactivityDeadline) AS t2_referrals,
				(SELECT base_hourly_mining_rate
					FROM adoption
					WHERE active = true) AS base_hourly_mining_rate,
				(SELECT GROUP_CONCAT(CAST(b.percentage AS string) || ':' || CAST(s.percentage AS string))
					FROM staking_bonus b
						JOIN staking s
							ON b.years = s.years AND s.user_id = ue.user_id) AS staking_info,
				(SELECT GROUP_CONCAT(b.type || ':' || b.amount || ',' || CAST(b.updated_at AS string))
					FROM balances b
					LEFT JOIN user_economy ue -- as we must earn only from active users, so need to use LEFT JOIN user_economy.
						ON ue.user_id = SUBSTR(b.type, POSITION('~', b.type) + 1)
					WHERE b.user_id = :userId AND (ue.last_mining_started_at < :inactivityDeadline OR ue.last_mining_started_at IS NULL) -- 'is null' is not fully right solution maybe. Need to be checked.
				) AS balances
			FROM user_economy ue	
			WHERE ue.hash_code % :chunks == :chunkIndex
			LIMIT 1000` // TODO: handle limit.

	var res []*userEconomy
	err := b.db.PrepareExecuteTyped(sql, params, &res)
	if err != nil {
		return nil, errors.Wrapf(err, "db execute getting user economy for (chunk:%v, chunkIndex:%v) failed",
			b.cfg.MessageBroker.Topics[3].Partitions, b.cfg.MessageBroker.Topics[3].Partition)
	}

	return res, nil
}

func (b *balanceSource) calculate(ue []*userEconomy) {
	// Here we should update all balances:
	// balances{type=standard}
	// balances{type=staking}
	// balances{type=total}

	// t0_referral_standard_earnings~{userId}
	// t1_referral_standard_earnings~{userId}
	// t2_referral_standard_earnings~{userId}
	// t0_referral_standard_earnings
	// t1_referral_standard_earnings
	// t2_referral_standard_earnings

	// t0_referral_staking_earnings~{userId}
	// t1_referral_staking_earnings~{userId}
	// t2_referral_staking_earnings~{userId}
	// t0_referral_staking_earnings
	// t1_referral_staking_earnings
	// t2_referral_staking_earnings

	for _, u := range ue {
		si := u.parseStakingInfo()
		bc := u.parseBalances()

		hourlyMiningRate, normalHourlyMiningRate := u.calculateMiningRates(b.cfg, si)
		stakedHourlyMiningRate := u.calculateStakingRate(si, hourlyMiningRate)

		val := bc["standard"]
		balanceStandard := u.calculateStandardBalance(&val, normalHourlyMiningRate)

		val = bc["staking"]
		balanceStaking := u.calculateStakingBalance(&val, stakedHourlyMiningRate)
		balanceTotal := u.calculateTotalBalance(balanceStandard, balanceStaking)

		// TODO: elapsedNanoseconds should be different for each balance. It is taken from balance struct.
		// TODO: handle balances. All of them can have different elapsed time, so maybe we need to handle them separately.
		for t, v := range bc {
			if t == "standard" {
				if err := b.updateBalance(u.UserID, "standard", coin.UnsafeNew(balanceStandard.String())); err != nil {
					// TODO: handle error.
				}
			} else if t == "staking" {
				if err := b.updateBalance(u.UserID, "staking", coin.UnsafeNew(balanceStaking.String())); err != nil {
					// TODO: handle error.
				}
			} else if t == "total" {
				if err := b.updateBalance(u.UserID, "staking", coin.UnsafeNew(balanceTotal.String())); err != nil {
					// TODO: handle error.
				}
			} else if t == "t0_referral_standard_earnings" {

			} else if t == "t1_referral_standard_earnings" {

			} else if t == "t2_referral_standard_earnings" {

			} else if strings.Contains(t, "t0_referral_standard_earnings~") {

			} else if strings.Contains(t, "t1_referral_standard_earnings~") {

			} else if strings.Contains(t, "t2_referral_standard_earnings~") {

			} else if t == "t0_referral_staking_earnings" {

			} else if t == "t1_referral_staking_earnings" {

			} else if t == "t2_referral_staking_earnings" {

			} else if strings.Contains(t, "t0_referral_staking_earnings~") {

			} else if strings.Contains(t, "t1_referral_staking_earnings~") {

			} else if strings.Contains(t, "t2_referral_staking_earnings~") {

			}
		}

		// -------------------------------------------------------
		// Referrals standard.
		//referralStandardCommonPart := math.NewUint(100 - si.Allocation).Mul(u.BaseHourlyMiningRate.Uint).Mul(elapsedNanoseconds)
		//t0ReferralStandardEarningsUser := referralStandardCommonPart.MulUint64(u.T0Referrals).QuoUint64(1440000000000000)
		// TODO: +=
		//t1ReferralStandardEarningsUser := referralStandardCommonPart.MulUint64(u.T1Referrals).QuoUint64(1440000000000000)
		// TODO: +=
		//t2ReferralStandardEarningsUser := referralStandardCommonPart.MulUint64(u.T2Referrals).QuoUint64(7200000000000000)

		// TODO: +=
		//t0ReferralStandardEarnings := t0ReferralStandardEarningsUser

		// TODO: summ
		// -------------------------------------------------------
		// Referrals staking.
		//referralStakingCommonPart := math.NewUint(si.Allocation * si.Bonus).Mul(elapsedNanoseconds)
		//t0ReferralStakingEarningsUser := math.NewUint(u.T0Referrals).Mul(referralStakingCommonPart).QuoUint64(144000000000000000)
		//t1ReferralStakingEarningsUser := math.NewUint(u.T1Referrals).Mul(referralStakingCommonPart).QuoUint64(144000000000000000)
		//t2ReferralStakingEarningsUser := math.NewUint(u.T2Referrals).Mul(referralStakingCommonPart).QuoUint64(720000000000000000)

		// TODO: +=
		//t0ReferralStandardEarnings := t0ReferralStandardEarningsUser
	}
}

func (u *userEconomy) calculateStandardBalance(b *balance, normalHourlyMiningRate *coin.ICEFlake) *coin.ICEFlake {
	elapsedNanoseconds := math.NewUint(uint64(time.Now().UnixNano()) - b.UpdatedAt)
	balanceStandardAddition := normalHourlyMiningRate.Mul(elapsedNanoseconds).Quo(math.NewUint(3600000000000))
	balanceStandard := math.Uint(b.Amount.Uint).Add(balanceStandardAddition)

	return &coin.ICEFlake{Uint: balanceStandard}
}

func (u *userEconomy) calculateStakingBalance(b *balance, stakedHourlyMiningRate *coin.ICEFlake) *coin.ICEFlake {
	elapsedNanoseconds := math.NewUint(uint64(time.Now().UnixNano()) - b.UpdatedAt)
	balanceStaking := math.Uint(b.Amount.Uint).Mul(stakedHourlyMiningRate.Mul(elapsedNanoseconds).QuoUint64(3600000000000))

	return &coin.ICEFlake{Uint: balanceStaking}
}

func (u *userEconomy) calculateTotalBalance(standard *coin.ICEFlake, staked *coin.ICEFlake) *coin.ICEFlake {
	return &coin.ICEFlake{Uint: standard.Uint.Add(staked.Uint)}
}

func (u *userEconomy) calculateMiningRates(cfg *config, si *stakingInfo) (*coin.ICEFlake, *coin.ICEFlake) {
	multiplier := math.NewUint(u.T0Referrals*cfg.Rates.Tier0 + u.T1Referrals*cfg.Rates.Tier1 + u.T2Referrals*cfg.Rates.Tier2)
	hourlyMiningRate := u.BaseHourlyMiningRate.Mul(multiplier).QuoUint64(100)
	normalHourlyMiningRate := math.NewUint(100 - si.Allocation).Mul(hourlyMiningRate).QuoUint64(100)

	return &coin.ICEFlake{Uint: hourlyMiningRate}, &coin.ICEFlake{Uint: normalHourlyMiningRate}
}

func (u *userEconomy) calculateStakingRate(si *stakingInfo, hourlyMiningRate *coin.ICEFlake) *coin.ICEFlake {
	return &coin.ICEFlake{Uint: math.NewUint(si.Bonus).Mul(hourlyMiningRate.Uint).Mul(math.NewUint(si.Allocation)).QuoUint64(10000)}
}

func (u *userEconomy) parseStakingInfo() *stakingInfo {
	parts := strings.Split(u.StakingInfo, ":")
	bonus, err := strconv.ParseUint(parts[0], base10, bitSize64)
	log.Panic(errors.Wrapf(err, "can't parse rate uint64 for adoption:%v", parts[0]))

	allocation, err := strconv.ParseUint(parts[1], base10, bitSize64)
	log.Panic(errors.Wrapf(err, "can't parse rate uint64 for adoption:%v", parts[0]))

	return &stakingInfo{Bonus: bonus, Allocation: allocation}
}

func (u *userEconomy) parseBalances() map[string]balance {
	bs := strings.Split(u.Balances, ",")
	res := make(map[string]balance, len(bs))

	for _, v := range bs {
		parts := strings.Split(v, ":")
		updatedAt, err := strconv.ParseUint(parts[2], base10, bitSize64)
		log.Panic(errors.Wrapf(err, "can't parse rate uint64 for balance updateAt:%v", parts[2]))

		res[parts[0]] = balance{Amount: coin.UnsafeNewAmount(parts[1]), UpdatedAt: updatedAt}
	}

	return res
}

func (b *balanceSource) updateBalance(userID UserID, balanceType string, balance *coin.Coin) error {
	params := map[string]interface{}{
		"amount":    balance.Amount,
		"amountW0":  balance.AmountWord0,
		"amountW1":  balance.AmountWord1,
		"amountW2":  balance.AmountWord2,
		"amountW3":  balance.AmountWord3,
		"userId":    userID,
		"type":      balanceType,
		"updatedAt": time.Now(),
	}

	sql := `UPDATE balances SET
				   amount = :balance,
				   amount_w0 = :amountW0,
				   amount_w1 = :amountW1,
				   amount_w2 = :amountW2,
				   amount_w3 = :amountW3,
				   updated_at = :updatedAt
				WHERE user_id = :userId AND type = :type`

	query, err := p.db.PrepareExecute(sql, params)

	if err = storage.CheckSQLDMLErr(query, err); err != nil {
		return errors.Wrapf(err, "failed to update balances with userID:%v and type:%v", userID, balanceType)
	}

	return nil
}
