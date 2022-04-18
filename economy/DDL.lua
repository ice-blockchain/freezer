-- SPDX-License-Identifier: BUSL-1.1

box.execute([[CREATE TABLE IF NOT EXISTS total_users  (
                    key STRING primary key,
                    value UNSIGNED NOT NULL
                    ) WITH ENGINE = 'vinyl';]])

-- Not sure if this is needed.
box.execute([[CREATE TABLE IF NOT EXISTS adoption  (
                    total_users UNSIGNED primary key,
                    base_hourly_mining_rate DOUBLE NOT NULL
                    ) WITH ENGINE = 'vinyl';]])

box.execute([[CREATE TABLE IF NOT EXISTS user_economy  (
                    user_id STRING primary key,
                    hash_code UNSIGNED NOT NULL UNIQUE,
                    balance DOUBLE NOT NULL,
                    last_mining_started_at UNSIGNED,
                    profile_picture_url STRING,
                    staking_years UNSIGNED,
                    staking_percentage DOUBLE,
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    balance_updated_at UNSIGNED NOT NULL
                    ) WITH ENGINE = 'vinyl';]])

box.execute([[CREATE TABLE IF NOT EXISTS t1_referral_earnings  (
                    user_id STRING STRING NOT NULL,
                    referral_user_id STRING NOT NULL,
                    earnings DOUBLE NOT NULL,
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    primary (user_id, referral_user_id)
                    ) WITH ENGINE = 'vinyl';]])

box.execute([[CREATE TABLE IF NOT EXISTS t2_referral_earnings  (
                    user_id STRING STRING NOT NULL,
                    referral_user_id STRING NOT NULL,
                    earnings DOUBLE NOT NULL,
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    primary (user_id, referral_user_id)
                    ) WITH ENGINE = 'vinyl';]])

-- TODO will add indexes later on