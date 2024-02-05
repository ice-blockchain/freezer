-- SPDX-License-Identifier: ice License 1.0

CREATE DATABASE IF NOT EXISTS light;
CREATE DATABASE IF NOT EXISTS dark;
CREATE TABLE IF NOT EXISTS light.freezer_user_history
(
       mining_session_solo_last_started_at DateTime64(9,'UTC')  DEFAULT 0,
       mining_session_solo_started_at DateTime64(9,'UTC')  DEFAULT 0,
       mining_session_solo_ended_at DateTime64(9,'UTC')  DEFAULT 0,
       mining_session_solo_previously_ended_at DateTime64(9,'UTC')  DEFAULT 0,
       extra_bonus_started_at DateTime64(9,'UTC')  DEFAULT 0,
       resurrect_solo_used_at DateTime64(9,'UTC')  DEFAULT 0,
       resurrect_t0_used_at DateTime64(9,'UTC')  DEFAULT 0,
       resurrect_tminus1_used_at DateTime64(9,'UTC')  DEFAULT 0,
       mining_session_solo_day_off_last_awarded_at DateTime64(9,'UTC')  DEFAULT 0,
       extra_bonus_last_claim_available_at DateTime64(9,'UTC')  DEFAULT 0,
       solo_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
       for_t0_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
       for_tminus1_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
       created_at DateTime('UTC')  DEFAULT 0,
       balance_total_standard Float64  DEFAULT 0,
       balance_total_pre_staking Float64  DEFAULT 0,
       balance_total_minted Float64  DEFAULT 0,
       balance_total_slashed Float64  DEFAULT 0,
       balance_solo_pending Float64  DEFAULT 0,
       balance_t1_pending Float64  DEFAULT 0,
       balance_t2_pending Float64  DEFAULT 0,
       balance_solo_pending_applied Float64  DEFAULT 0,
       balance_t1_pending_applied Float64  DEFAULT 0,
       balance_t2_pending_applied Float64  DEFAULT 0,
       balance_solo Float64  DEFAULT 0,
       balance_t0 Float64  DEFAULT 0,
       balance_t1 Float64  DEFAULT 0,
       balance_t2 Float64  DEFAULT 0,
       balance_for_t0 Float64  DEFAULT 0,
       balance_for_tminus1 Float64  DEFAULT 0,
       balance_solo_ethereum Float64  DEFAULT 0,
       balance_t0_ethereum Float64  DEFAULT 0,
       balance_t1_ethereum Float64  DEFAULT 0,
       balance_t2_ethereum Float64  DEFAULT 0,
       balance_for_t0_ethereum Float64  DEFAULT 0,
       balance_for_tminus1_ethereum Float64  DEFAULT 0,
       balance_solo_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
       balance_t0_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
       balance_t1_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
       balance_t2_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
       balance_for_t0_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
       balance_for_tminus1_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
       slashing_rate_solo Float64  DEFAULT 0,
       slashing_rate_t0 Float64  DEFAULT 0,
       slashing_rate_t1 Float64  DEFAULT 0,
       slashing_rate_t2 Float64  DEFAULT 0,
       slashing_rate_for_t0 Float64  DEFAULT 0,
       slashing_rate_for_tminus1 Float64  DEFAULT 0,
       id Int64  DEFAULT 0,
       id_t0 Int64  DEFAULT 0,
       id_tminus1 Int64  DEFAULT 0,
       active_t1_referrals Int32  DEFAULT 0,
       active_t2_referrals Int32  DEFAULT 0,
       pre_staking_bonus UInt16  DEFAULT 0,
       pre_staking_allocation UInt16  DEFAULT 0,
       extra_bonus UInt16  DEFAULT 0,
       news_seen UInt16  DEFAULT 0,
       extra_bonus_days_claim_not_available UInt16  DEFAULT 0,
       utc_offset Int16  DEFAULT 0,
       kyc_step_passed UInt8  DEFAULT 0,
       kyc_step_blocked UInt8  DEFAULT 0,
       hide_ranking Bool  DEFAULT FALSE,
       kyc_steps_created_at Array(DateTime64(9,'UTC')) DEFAULT [],
       kyc_steps_last_updated_at Array(DateTime64(9,'UTC')) DEFAULT [],
       country String  DEFAULT '',
       profile_picture_name String  DEFAULT '',
       username String  DEFAULT '',
       mining_blockchain_account_address String  DEFAULT '',
       blockchain_account_address String  DEFAULT '',
       user_id String  DEFAULT ''
) ENGINE=ReplicatedMergeTree('/clickhouse/tables/{cluster}/{shard_light}/freezer_user_history', '{replica_light}')
  PARTITION BY toDate(created_at)
  PRIMARY KEY (id, created_at);
 
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS solo_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER extra_bonus_last_claim_available_at;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS for_t0_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER solo_last_ethereum_coin_distribution_processed_at;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS for_tminus1_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER for_t0_last_ethereum_coin_distribution_processed_at;

ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_solo_ethereum Float64 DEFAULT 0 AFTER balance_for_tminus1;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t0_ethereum Float64 DEFAULT 0 AFTER balance_solo_ethereum;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t1_ethereum Float64 DEFAULT 0 AFTER balance_t0_ethereum;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t2_ethereum Float64 DEFAULT 0 AFTER balance_t1_ethereum;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_t0_ethereum Float64 DEFAULT 0 AFTER balance_t2_ethereum;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_tminus1_ethereum Float64 DEFAULT 0 AFTER balance_for_t0_ethereum;

ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_solo_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_for_tminus1_ethereum;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t0_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_solo_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t1_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t0_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t2_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t1_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_t0_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t2_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_tminus1_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_for_t0_ethereum_mainnet_reward_pool_contribution;

ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_step_passed UInt8 DEFAULT 0 AFTER utc_offset;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_step_blocked UInt8 DEFAULT 0 AFTER kyc_step_passed;

ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_steps_created_at Array(DateTime64(9,'UTC')) DEFAULT [] AFTER hide_ranking;
ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_steps_last_updated_at Array(DateTime64(9,'UTC')) DEFAULT [] AFTER kyc_steps_created_at;

ALTER TABLE light.freezer_user_history
    ADD COLUMN IF NOT EXISTS country String  DEFAULT '' AFTER kyc_steps_last_updated_at;
  
CREATE TABLE IF NOT EXISTS dark.freezer_user_history
(
      mining_session_solo_last_started_at DateTime64(9,'UTC')  DEFAULT 0,
      mining_session_solo_started_at DateTime64(9,'UTC')  DEFAULT 0,
      mining_session_solo_ended_at DateTime64(9,'UTC')  DEFAULT 0,
      mining_session_solo_previously_ended_at DateTime64(9,'UTC')  DEFAULT 0,
      extra_bonus_started_at DateTime64(9,'UTC')  DEFAULT 0,
      resurrect_solo_used_at DateTime64(9,'UTC')  DEFAULT 0,
      resurrect_t0_used_at DateTime64(9,'UTC')  DEFAULT 0,
      resurrect_tminus1_used_at DateTime64(9,'UTC')  DEFAULT 0,
      mining_session_solo_day_off_last_awarded_at DateTime64(9,'UTC')  DEFAULT 0,
      extra_bonus_last_claim_available_at DateTime64(9,'UTC')  DEFAULT 0,
      solo_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
      for_t0_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
      for_tminus1_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
      created_at DateTime('UTC')  DEFAULT 0,
      balance_total_standard Float64  DEFAULT 0,
      balance_total_pre_staking Float64  DEFAULT 0,
      balance_total_minted Float64  DEFAULT 0,
      balance_total_slashed Float64  DEFAULT 0,
      balance_solo_pending Float64  DEFAULT 0,
      balance_t1_pending Float64  DEFAULT 0,
      balance_t2_pending Float64  DEFAULT 0,
      balance_solo_pending_applied Float64  DEFAULT 0,
      balance_t1_pending_applied Float64  DEFAULT 0,
      balance_t2_pending_applied Float64  DEFAULT 0,
      balance_solo Float64  DEFAULT 0,
      balance_t0 Float64  DEFAULT 0,
      balance_t1 Float64  DEFAULT 0,
      balance_t2 Float64  DEFAULT 0,
      balance_for_t0 Float64  DEFAULT 0,
      balance_for_tminus1 Float64  DEFAULT 0,
      balance_solo_ethereum Float64  DEFAULT 0,
      balance_t0_ethereum Float64  DEFAULT 0,
      balance_t1_ethereum Float64  DEFAULT 0,
      balance_t2_ethereum Float64  DEFAULT 0,
      balance_for_t0_ethereum Float64  DEFAULT 0,
      balance_for_tminus1_ethereum Float64  DEFAULT 0,
      balance_solo_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
      balance_t0_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
      balance_t1_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
      balance_t2_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
      balance_for_t0_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
      balance_for_tminus1_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
      slashing_rate_solo Float64  DEFAULT 0,
      slashing_rate_t0 Float64  DEFAULT 0,
      slashing_rate_t1 Float64  DEFAULT 0,
      slashing_rate_t2 Float64  DEFAULT 0,
      slashing_rate_for_t0 Float64  DEFAULT 0,
      slashing_rate_for_tminus1 Float64  DEFAULT 0,
      id Int64  DEFAULT 0,
      id_t0 Int64  DEFAULT 0,
      id_tminus1 Int64  DEFAULT 0,
      active_t1_referrals Int32  DEFAULT 0,
      active_t2_referrals Int32  DEFAULT 0,
      pre_staking_bonus UInt16  DEFAULT 0,
      pre_staking_allocation UInt16  DEFAULT 0,
      extra_bonus UInt16  DEFAULT 0,
      news_seen UInt16  DEFAULT 0,
      extra_bonus_days_claim_not_available UInt16  DEFAULT 0,
      utc_offset Int16  DEFAULT 0,
      kyc_step_passed UInt8  DEFAULT 0,
      kyc_step_blocked UInt8  DEFAULT 0,
      hide_ranking Bool  DEFAULT FALSE,
      kyc_steps_created_at Array(DateTime64(9,'UTC')) DEFAULT [],
      kyc_steps_last_updated_at Array(DateTime64(9,'UTC')) DEFAULT [],
      country String  DEFAULT '',
      profile_picture_name String  DEFAULT '',
      username String  DEFAULT '',
      mining_blockchain_account_address String  DEFAULT '',
      blockchain_account_address String  DEFAULT '',
      user_id String  DEFAULT ''
) ENGINE=ReplicatedMergeTree('/clickhouse/tables/{cluster}/{shard_dark}/freezer_user_history', '{replica_dark}')
  PARTITION BY toDate(created_at)
  PRIMARY KEY (id, created_at);
  
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS solo_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER extra_bonus_last_claim_available_at;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS for_t0_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER solo_last_ethereum_coin_distribution_processed_at;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS for_tminus1_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER for_t0_last_ethereum_coin_distribution_processed_at;

ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_solo_ethereum Float64 DEFAULT 0 AFTER balance_for_tminus1;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t0_ethereum Float64 DEFAULT 0 AFTER balance_solo_ethereum;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t1_ethereum Float64 DEFAULT 0 AFTER balance_t0_ethereum;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t2_ethereum Float64 DEFAULT 0 AFTER balance_t1_ethereum;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_t0_ethereum Float64 DEFAULT 0 AFTER balance_t2_ethereum;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_tminus1_ethereum Float64 DEFAULT 0 AFTER balance_for_t0_ethereum;

ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_solo_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_for_tminus1_ethereum;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t0_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_solo_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t1_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t0_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t2_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t1_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_t0_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t2_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_tminus1_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_for_t0_ethereum_mainnet_reward_pool_contribution;

ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_step_passed UInt8 DEFAULT 0 AFTER utc_offset;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_step_blocked UInt8 DEFAULT 0 AFTER kyc_step_passed;

ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_steps_created_at Array(DateTime64(9,'UTC')) DEFAULT [] AFTER hide_ranking;
ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_steps_last_updated_at Array(DateTime64(9,'UTC')) DEFAULT [] AFTER kyc_steps_created_at;

ALTER TABLE dark.freezer_user_history
    ADD COLUMN IF NOT EXISTS country String  DEFAULT '' AFTER kyc_steps_last_updated_at;
  
CREATE TABLE IF NOT EXISTS freezer_user_history
(
     mining_session_solo_last_started_at DateTime64(9,'UTC')  DEFAULT 0,
     mining_session_solo_started_at DateTime64(9,'UTC')  DEFAULT 0,
     mining_session_solo_ended_at DateTime64(9,'UTC')  DEFAULT 0,
     mining_session_solo_previously_ended_at DateTime64(9,'UTC')  DEFAULT 0,
     extra_bonus_started_at DateTime64(9,'UTC')  DEFAULT 0,
     resurrect_solo_used_at DateTime64(9,'UTC')  DEFAULT 0,
     resurrect_t0_used_at DateTime64(9,'UTC')  DEFAULT 0,
     resurrect_tminus1_used_at DateTime64(9,'UTC')  DEFAULT 0,
     mining_session_solo_day_off_last_awarded_at DateTime64(9,'UTC')  DEFAULT 0,
     extra_bonus_last_claim_available_at DateTime64(9,'UTC')  DEFAULT 0,
     solo_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
     for_t0_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
     for_tminus1_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC')  DEFAULT 0,
     created_at DateTime('UTC')  DEFAULT 0,
     balance_total_standard Float64  DEFAULT 0,
     balance_total_pre_staking Float64  DEFAULT 0,
     balance_total_minted Float64  DEFAULT 0,
     balance_total_slashed Float64  DEFAULT 0,
     balance_solo_pending Float64  DEFAULT 0,
     balance_t1_pending Float64  DEFAULT 0,
     balance_t2_pending Float64  DEFAULT 0,
     balance_solo_pending_applied Float64  DEFAULT 0,
     balance_t1_pending_applied Float64  DEFAULT 0,
     balance_t2_pending_applied Float64  DEFAULT 0,
     balance_solo Float64  DEFAULT 0,
     balance_t0 Float64  DEFAULT 0,
     balance_t1 Float64  DEFAULT 0,
     balance_t2 Float64  DEFAULT 0,
     balance_for_t0 Float64  DEFAULT 0,
     balance_for_tminus1 Float64  DEFAULT 0,
     balance_solo_ethereum Float64  DEFAULT 0,
     balance_t0_ethereum Float64  DEFAULT 0,
     balance_t1_ethereum Float64  DEFAULT 0,
     balance_t2_ethereum Float64  DEFAULT 0,
     balance_for_t0_ethereum Float64  DEFAULT 0,
     balance_for_tminus1_ethereum Float64  DEFAULT 0,
     balance_solo_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
     balance_t0_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
     balance_t1_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
     balance_t2_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
     balance_for_t0_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
     balance_for_tminus1_ethereum_mainnet_reward_pool_contribution Float64  DEFAULT 0,
     slashing_rate_solo Float64  DEFAULT 0,
     slashing_rate_t0 Float64  DEFAULT 0,
     slashing_rate_t1 Float64  DEFAULT 0,
     slashing_rate_t2 Float64  DEFAULT 0,
     slashing_rate_for_t0 Float64  DEFAULT 0,
     slashing_rate_for_tminus1 Float64  DEFAULT 0,
     id Int64  DEFAULT 0,
     id_t0 Int64  DEFAULT 0,
     id_tminus1 Int64  DEFAULT 0,
     active_t1_referrals Int32  DEFAULT 0,
     active_t2_referrals Int32  DEFAULT 0,
     pre_staking_bonus UInt16  DEFAULT 0,
     pre_staking_allocation UInt16  DEFAULT 0,
     extra_bonus UInt16  DEFAULT 0,
     news_seen UInt16  DEFAULT 0,
     extra_bonus_days_claim_not_available UInt16  DEFAULT 0,
     utc_offset Int16  DEFAULT 0,
     kyc_step_passed UInt8  DEFAULT 0,
     kyc_step_blocked UInt8  DEFAULT 0,
     hide_ranking Bool  DEFAULT FALSE,
     kyc_steps_created_at Array(DateTime64(9,'UTC')) DEFAULT [],
     kyc_steps_last_updated_at Array(DateTime64(9,'UTC')) DEFAULT [],
     country String  DEFAULT '',
     profile_picture_name String  DEFAULT '',
     username String  DEFAULT '',
     mining_blockchain_account_address String  DEFAULT '',
     blockchain_account_address String  DEFAULT '',
     user_id String  DEFAULT ''
) ENGINE = Distributed('{cluster}', '', 'freezer_user_history', toUInt64(toDate(created_at)));

ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS solo_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER extra_bonus_last_claim_available_at;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS for_t0_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER solo_last_ethereum_coin_distribution_processed_at;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS for_tminus1_last_ethereum_coin_distribution_processed_at DateTime64(9,'UTC') DEFAULT 0 AFTER for_t0_last_ethereum_coin_distribution_processed_at;

ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_solo_ethereum Float64 DEFAULT 0 AFTER balance_for_tminus1;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t0_ethereum Float64 DEFAULT 0 AFTER balance_solo_ethereum;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t1_ethereum Float64 DEFAULT 0 AFTER balance_t0_ethereum;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t2_ethereum Float64 DEFAULT 0 AFTER balance_t1_ethereum;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_t0_ethereum Float64 DEFAULT 0 AFTER balance_t2_ethereum;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_tminus1_ethereum Float64 DEFAULT 0 AFTER balance_for_t0_ethereum;

ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_solo_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_for_tminus1_ethereum;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t0_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_solo_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t1_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t0_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_t2_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t1_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_t0_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_t2_ethereum_mainnet_reward_pool_contribution;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS balance_for_tminus1_ethereum_mainnet_reward_pool_contribution Float64 DEFAULT 0 AFTER balance_for_t0_ethereum_mainnet_reward_pool_contribution;

ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_step_passed UInt8 DEFAULT 0 AFTER utc_offset;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_step_blocked UInt8 DEFAULT 0 AFTER kyc_step_passed;

ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_steps_created_at Array(DateTime64(9,'UTC')) DEFAULT [] AFTER hide_ranking;
ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS kyc_steps_last_updated_at Array(DateTime64(9,'UTC')) DEFAULT [] AFTER kyc_steps_created_at;

ALTER TABLE freezer_user_history
    ADD COLUMN IF NOT EXISTS country String  DEFAULT '' AFTER kyc_steps_last_updated_at;