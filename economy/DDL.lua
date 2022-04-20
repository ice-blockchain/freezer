-- SPDX-License-Identifier: BUSL-1.1

box.execute([[CREATE TABLE IF NOT EXISTS total_users  (
                    key STRING primary key,
                    value UNSIGNED NOT NULL
                    ) WITH ENGINE = 'vinyl';]])

box.execute([[CREATE TABLE IF NOT EXISTS users  (
                    id STRING primary key,
                    hash_code UNSIGNED NOT NULL UNIQUE,
                    referred_by STRING REFERENCES users(id) ON DELETE SET NULL,
                    username STRING NOT NULL UNIQUE,
                    email STRING,
                    full_name STRING,
                    phone_number STRING,
                    profile_picture STRING NOT NULL,
                    country STRING NOT NULL,
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    deleted_at UNSIGNED
                    ) WITH ENGINE = 'vinyl';]])
                    
-- Not sure if this is needed.
box.execute([[CREATE TABLE IF NOT EXISTS adoption  (
                    total_users UNSIGNED primary key,
                    base_hourly_mining_rate DOUBLE NOT NULL
                    ) WITH ENGINE = 'vinyl';]])

box.execute([[CREATE TABLE IF NOT EXISTS user_economy  (
                    user_id STRING primary key,
                    profile_picture_url STRING,
                    balance DOUBLE NOT NULL,
                    staking_percentage DOUBLE,
                    hash_code UNSIGNED NOT NULL UNIQUE,
                    last_mining_started_at UNSIGNED,
                    staking_years UNSIGNED,
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    balance_updated_at UNSIGNED NOT NULL
                    ) WITH ENGINE = 'vinyl';]])

box.execute([[CREATE TABLE IF NOT EXISTS t1_referral_earnings  (
                    user_id STRING NOT NULL,
                    referral_user_id STRING NOT NULL,
                    earnings DOUBLE NOT NULL,
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    primary key(user_id, referral_user_id)
                    ) WITH ENGINE = 'vinyl';]])

box.execute([[CREATE TABLE IF NOT EXISTS t2_referral_earnings  (
                    user_id STRING NOT NULL,
                    referral_user_id STRING NOT NULL,
                    earnings DOUBLE NOT NULL,
                    created_at UNSIGNED NOT NULL,
                    updated_at UNSIGNED NOT NULL,
                    primary key(user_id, referral_user_id)
                    ) WITH ENGINE = 'vinyl';]])
