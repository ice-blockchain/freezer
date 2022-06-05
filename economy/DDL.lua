-- SPDX-License-Identifier: BUSL-1.1

box.execute([[CREATE TABLE IF NOT EXISTS global  (
                    key STRING primary key,
                    value SCALAR NOT NULL
                    ) WITH ENGINE = 'vinyl';]])
-- (key,value) : ('TOTAL_USERS', 10000) -----> ++ or -- when user registers/deletes account
-- (key,value) : ('TOTAL_ACTIVE_USERS', 10) -----> you reset this with `select count(1) from user_economy where last_mining_started_at < 24h` and you do that in the same logic you populate adoption_history in

box.execute([[CREATE TABLE IF NOT EXISTS total_users_history  (
                    minute_timestamp UNSIGNED primary key,
                    hour_timestamp UNSIGNED NOT NULL CHECK (minute_timestamp >= hour_timestamp*60 and minute_timestamp < (hour_timestamp+1)*60),
                    day_timestamp UNSIGNED NOT NULL CHECK (hour_timestamp >= day_timestamp*24 and hour_timestamp < (day_timestamp+1)*24),
                    date_ STRING NOT NULL,
                    total_users UNSIGNED NOT NULL DEFAULT 0
                    ) WITH ENGINE = 'vinyl';]])
box.execute([[CREATE INDEX IF NOT EXISTS total_users_history_day_timestamp_ix ON total_users_history (day_timestamp);]])
box.execute([[CREATE INDEX IF NOT EXISTS total_users_history_date_ix ON total_users_history (date_);]])
-- every minute, total_users_history.total_users = global.value where global.key = 'TOTAL_USERS'

box.execute([[CREATE TABLE IF NOT EXISTS adoption  (
                    total_active_users UNSIGNED primary key,
                    base_hourly_mining_rate STRING NOT NULL,
                    active BOOLEAN NOT NULL DEFAULT false
                    ) WITH ENGINE = 'vinyl';]])
box.execute([[INSERT INTO adoption (total_active_users, base_hourly_mining_rate, active)
                          VALUES (0, '16000000000', true),
                                 (50000, '8000000000', false),
                                 (250000, '4000000000', false),
                                 (1250000, '2000000000', false),
                                 (6250000, '1000000000', false),
                                 (31250000, '500000000', false)
          ]])
-- base_hourly_mining_rate is in ice flakes
-- IF the last 168 consecutive hours from adoption_history.hour_timestamp have ALL been >= ANY adoption.total_active_users,
-- then adoption.active of that entry becomes true and the previous active adoption entry becomes false.

box.execute([[CREATE TABLE IF NOT EXISTS adoption_history  (
                    minute_timestamp UNSIGNED primary key,
                    hour_timestamp UNSIGNED NOT NULL CHECK (minute_timestamp >= hour_timestamp*60 and minute_timestamp < (hour_timestamp+1)*60),
                    total_active_users UNSIGNED NOT NULL DEFAULT 0
                    ) WITH ENGINE = 'vinyl';]])
-- minute_timestamp = time.Now().UTC().Unix()/60
-- hour_timestamp = minute_timestamp/60

box.execute([[CREATE TABLE IF NOT EXISTS user_economy  (
                    last_mining_started_at UNSIGNED DEFAULT 0,
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    user_id STRING primary key,
                    username STRING NOT NULL UNIQUE,
                    profile_picture_url STRING,
                    hash_code UNSIGNED NOT NULL UNIQUE
                    ) WITH ENGINE = 'vinyl';]])
-- balance is in ice flakes
-- if staking is enabled for the user, and the percentage is 100%, balances.amount{type=standard} is gonna always be 0.
box.execute([[CREATE INDEX IF NOT EXISTS user_economy_last_mining_started_at_ix ON user_economy (last_mining_started_at);]])
box.execute([[CREATE INDEX IF NOT EXISTS user_economy_username_ix ON user_economy (username);]])

box.execute([[CREATE TABLE IF NOT EXISTS staking  (
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    user_id STRING primary key REFERENCES user_economy(user_id) ON DELETE CASCADE,
                    percentage UNSIGNED NOT NULL,
                    years UNSIGNED NOT NULL
                    ) WITH ENGINE = 'vinyl';]])
-- When staking happens, you move staking.percentage*balances.amount{type=standard}/100 to balances.amount{type=staking}, for that user_id

box.execute([[CREATE TABLE IF NOT EXISTS balances (
                    updated_at UNSIGNED NOT NULL,
                    amount STRING NOT NULL DEFAULT '0',
                    user_id STRING NOT NULL REFERENCES user_economy(user_id) ON DELETE CASCADE,
                    type STRING NOT NULL
                        CHECK (
                            lower(type) == 'standard' or
                            lower(type) == 'staking' or
                            lower(type) == 'total' or
                            lower(type) == 't0_referral_standard_earnings' or
                            POSITION('t0_referral_standard_earnings~', lower(type)) != 0  or
                            lower(type) == 't0_referral_staking_earnings' or
                            POSITION('t0_referral_staking_earnings~', lower(type)) != 0  or
                            lower(type) == 't1_referral_standard_earnings' or
                            POSITION('t1_referral_standard_earnings~', lower(type)) != 0  or
                            lower(type) == 't1_referral_staking_earnings' or
                            POSITION('t1_referral_staking_earnings~', lower(type)) != 0  or
                            lower(type) == 't2_referral_standard_earnings' or
                            POSITION('t2_referral_standard_earnings~', lower(type)) != 0  or
                            lower(type) == 't2_referral_staking_earnings' or
                            POSITION('t2_referral_staking_earnings~', lower(type)) != 0
                        ),
                    amount_w0 UNSIGNED NOT NULL DEFAULT 0,
                    amount_w1 UNSIGNED NOT NULL DEFAULT 0,
                    amount_w2 UNSIGNED NOT NULL DEFAULT 0,
                    amount_w3 UNSIGNED NOT NULL DEFAULT 0,
                    
                    primary key (user_id, type)) WITH ENGINE = 'vinyl';]])
box.execute([[CREATE INDEX IF NOT EXISTS balances_amount_words_ix ON balances (amount_w3, amount_w2, amount_w1, amount_w0);]])
-- amount is in ice flakes
-- t0 is the user that referred/invited the user to the app, so T0 -invited> user_id -invited> T1 -invited> T2
-- tX_referral_standard_earnings~{userId} => I.E. type = 't1_referral_standard_earnings~did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2'

box.execute([[CREATE TABLE IF NOT EXISTS staking_bonus  (
                    years UNSIGNED primary key,
                    percentage UNSIGNED NOT NULL
                    ) WITH ENGINE = 'vinyl';]])
box.execute([[INSERT INTO staking_bonus (years, percentage)
                          VALUES (1, 100),
                                 (2, 200),
                                 (3, 300),
                                 (4, 400),
                                 (5, 500)
          ]])

-- BALANCE calculation (for user_id = '1'):
-- t0Referrals = select count(1) from balances t0 join user_economy u on t0.user_id = '1' and t0.type = 't0_referral_standard_earnings~' || u.user_id and u.last_mining_started_at < 24h ago
-- ## t0Referrals is <= 1, always
-- t1Referrals = select count(1) from balances t1 join user_economy u on t1.user_id = '1' and t1.type = 't1_referral_standard_earnings~' || u.user_id and u.last_mining_started_at < 24h ago
-- t2Referrals = select count(1) from balances t2 join user_economy u on t2.user_id = '1' and t2.type = 't2_referral_standard_earnings~' || u.user_id and u.last_mining_started_at < 24h ago
-- baseHourlyMiningRate = select base_hourly_mining_rate from adoption where active = true
-- (stakingPercentageBonus, stakingPercentageAllocation) = select b.percentage as bonus, s.percentage as allocation from staking_bonus b join staking s on b.years = s.years and s.user_id = '1'
-- elapsedNanoseconds = now - balances.updated_at

-- hourlyMiningRate = baseHourlyMiningRate * (t0Referrals*25 + t1Referrals*25 + t2Referrals*5 + 100) / 100
-- normalHourlyMiningRate = (100-stakingPercentageAllocation) * hourlyMiningRate / 100
-- balances{type=standard}.amount += normalHourlyMiningRate * elapsedNanoseconds / 3600000000000

-- stakedHourlyMiningRate = stakingPercentageBonus * hourlyMiningRate * stakingPercentageAllocation / 10000
-- balances{type=staking}.amount += stakedHourlyMiningRate * elapsedNanoseconds / 3600000000000

-- balances{type=total}.amount = balances{type=standard}.amount + balances{type=staking}.amount

-- Referral EARNINGS
-- balances{type=t0_referral_standard_earnings~{userId}}.amount += t0Referrals * (100-stakingPercentageAllocation) * baseHourlyMiningRate * elapsedNanoseconds / 1440000000000000
-- balances{type=t1_referral_standard_earnings~{userId}}.amount += t1Referrals * (100-stakingPercentageAllocation) * baseHourlyMiningRate * elapsedNanoseconds / 1440000000000000
-- balances{type=t2_referral_standard_earnings~{userId}}.amount += t2Referrals * (100-stakingPercentageAllocation) * baseHourlyMiningRate * elapsedNanoseconds / 7200000000000000
-- balances{type=t0_referral_standard_earnings}.amount += t0_referral_standard_earnings~*.amount
-- balances{type=t1_referral_standard_earnings}.amount += sum(balances{type=t1_referral_standard_earnings~*}.amount)
-- balances{type=t2_referral_standard_earnings}.amount += sum(balances{type=t2_referral_standard_earnings~*}.amount)
-- staking
-- balances{type=t0_referral_staking_earnings~{userId}}.amount += t0Referrals * stakingPercentageBonus * stakingPercentageAllocation * baseHourlyMiningRate * elapsedNanoseconds / 144000000000000000
-- balances{type=t1_referral_staking_earnings~{userId}}.amount += t1Referrals * stakingPercentageBonus * stakingPercentageAllocation * baseHourlyMiningRate * elapsedNanoseconds / 144000000000000000
-- balances{type=t2_referral_staking_earnings~{userId}}.amount += t2Referrals * stakingPercentageBonus * stakingPercentageAllocation * baseHourlyMiningRate * elapsedNanoseconds / 720000000000000000
-- balances{type=t0_referral_staking_earnings}.amount += balances{type=t0_referral_staking_earnings}.amount
-- balances{type=t1_referral_staking_earnings}.amount += sum(balances{type=t1_referral_staking_earnings~*}.amount)
-- balances{type=t2_referral_staking_earnings}.amount += sum(balances{type=t2_referral_staking_earnings~*}.amount)

-- Balance related SQL example
--SELECT ...
--FROM ... x1
--WHERE   ...
---------- ‼️This is how you compare balances. I.E. In this case we look for balances that are >= than a specific one provided as an arg (via its words)‼️
--        AND (CASE
--                WHEN x1.amount_w3 == :amount_w3
--                    THEN (CASE
--                             WHEN x1.amount_w2 == :amount_w2
--                                 THEN (CASE
--                                          WHEN x1.amount_w1 == :amount_w1
--                                              THEN (x1.amount_w0 >= :amount_w0)
--                                          ELSE x1.amount_w1 > :amount_w1
--                                       END)
--                             ELSE x1.amount_w2 > :amount_w2
--                          END)
--                ELSE x1.amount_w3 > :amount_w3
--             END)
-- ‼️This is how you sort balances‼️
--ORDER BY x1.amount_w3 DESC,
--         x1.amount_w2 DESC,
--         x1.amount_w1 DESC,
--         x1.amount_w0 DESC;