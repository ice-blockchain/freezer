-- SPDX-License-Identifier: ice License 1.0
--************************************************************************************************************************************
-- pre_staking_bonuses
box.execute([[CREATE TABLE IF NOT EXISTS pre_staking_bonuses (
                                                          years UNSIGNED PRIMARY KEY,
                                                          bonus UNSIGNED NOT NULL CHECK (bonus > 0)
                                                         )
                                                         WITH ENGINE = 'memtx';]])
box.execute([[INSERT INTO pre_staking_bonuses (years, bonus)
                                   VALUES (1,     35),
                                          (2,     70),
                                          (3,     115),
                                          (4,     170),
                                          (5,     250);]])
--************************************************************************************************************************************
-- extra_bonuses
box.execute([[CREATE TABLE IF NOT EXISTS extra_bonuses (
                                                          ix    UNSIGNED PRIMARY KEY,
                                                          bonus UNSIGNED NOT NULL DEFAULT 0
                                                         )
                                                         WITH ENGINE = 'memtx';]])
box.execute([[INSERT INTO extra_bonuses (ix, bonus) VALUES %[3]v;]])
--************************************************************************************************************************************
-- global
box.execute([[CREATE TABLE IF NOT EXISTS global (
                                                 key   STRING PRIMARY KEY,
                                                 value SCALAR NOT NULL
                                                )
                                                WITH ENGINE = 'memtx';]])
--************************************************************************************************************************************
-- extra_bonus_start_date
box.execute([[CREATE TABLE IF NOT EXISTS extra_bonus_start_date (
                                                 key UNSIGNED NOT NULL PRIMARY KEY,
                                                 value UNSIGNED NOT NULL
                                                )
                                                WITH ENGINE = 'memtx';]])
box.execute([[INSERT INTO extra_bonus_start_date (key, value) VALUES (0,%[4]v);]])
--************************************************************************************************************************************                                                
-- adoption                                                
box.execute([[CREATE TABLE IF NOT EXISTS adoption (
                                                   achieved_at             UNSIGNED,
                                                   base_mining_rate        STRING NOT NULL,
                                                   milestone               UNSIGNED PRIMARY KEY,
                                                   total_active_users      UNSIGNED NOT NULL
                                                  )
                                                  WITH ENGINE = 'memtx';]])
box.execute([[INSERT INTO adoption (milestone, total_active_users, base_mining_rate, achieved_at)
                            VALUES (1,         0,                  '16000000000',           %[1]v),
                                   (2,         %[5]v,              '8000000000',            null),
                                   (3,         %[6]v,              '4000000000',            null),
                                   (4,         %[7]v,              '2000000000',            null),
                                   (5,         %[8]v,              '1000000000',            null),
                                   (6,         %[9]v,              '500000000',             null);]])
--************************************************************************************************************************************
-- users
box.execute([[CREATE TABLE IF NOT EXISTS users (
                        created_at                                              UNSIGNED NOT NULL,
                        updated_at                                              UNSIGNED NOT NULL,
                        rollback_used_at                                        UNSIGNED,
                        last_natural_mining_started_at                          UNSIGNED,
                        last_mining_started_at                                  UNSIGNED,
                        last_mining_ended_at                                    UNSIGNED,
                        previous_mining_started_at                              UNSIGNED,
                        previous_mining_ended_at                                UNSIGNED,
                        last_free_mining_session_awarded_at                     UNSIGNED,
                        user_id                                                 STRING PRIMARY KEY,
                        referred_by                                             STRING,
                        username                                                STRING,
                        first_name                                              STRING,
                        last_name                                               STRING,
                        profile_picture_name                                    STRING,
                        mining_blockchain_account_address                       STRING,
                        blockchain_account_address                              STRING,
                        hash_code                                               UNSIGNED NOT NULL,
                        hide_ranking                                            BOOLEAN NOT NULL DEFAULT FALSE,
                        verified                                                BOOLEAN NOT NULL DEFAULT FALSE
                    )
                     WITH ENGINE = 'memtx';]])
box.execute([[CREATE INDEX IF NOT EXISTS users_referred_by_idx ON users (referred_by);]])
box.execute([[CREATE INDEX IF NOT EXISTS top_miners_lookup_idx ON users (username,first_name,last_name);]])
--************************************************************************************************************************************
-- balances
box.execute([[CREATE TABLE IF NOT EXISTS balances (
                                                   amount     STRING NOT NULL DEFAULT '0',
                                                   amount_w0  UNSIGNED NOT NULL DEFAULT 0,
                                                   amount_w1  UNSIGNED NOT NULL DEFAULT 0,
                                                   amount_w2  UNSIGNED NOT NULL DEFAULT 0,
                                                   amount_w3  UNSIGNED NOT NULL DEFAULT 0,
                                                   user_id    STRING NOT NULL PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE
                                                  )
                                                  WITH ENGINE = 'memtx';]])
box.execute([[CREATE INDEX IF NOT EXISTS balances_amount_words_ix ON balances (amount_w3, amount_w2, amount_w1, amount_w0);]])
--************************************************************************************************************************************
-- processed_add_balance_commands
box.execute([[CREATE TABLE IF NOT EXISTS processed_add_balance_commands (
                                                 user_id STRING NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                                                 key     STRING NOT NULL,
                                                 PRIMARY KEY (user_id, key)
                                                )
                                                 WITH ENGINE = 'vinyl';]])
--************************************************************************************************************************************
-- processed_seen_news
box.execute([[CREATE TABLE IF NOT EXISTS processed_seen_news (
                                                 user_id STRING NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                                                 news_id STRING NOT NULL,
                                                 PRIMARY KEY (user_id, news_id)
                                                )
                                                 WITH ENGINE = 'vinyl';]])
--************************************************************************************************************************************
-- mining_sessions_dlq
for worker_index=0,%[2]v do
        box.execute([[CREATE TABLE IF NOT EXISTS mining_sessions_dlq_]] .. worker_index .. [[ (
                               id              STRING NOT NULL PRIMARY KEY,
                               user_id         STRING NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                               message         STRING NOT NULL
                              )
                            WITH ENGINE = 'vinyl';]])
end
--************************************************************************************************************************************
-- extra_bonuses
for worker_index=0,%[2]v do
        box.execute([[CREATE TABLE IF NOT EXISTS extra_bonuses_]] .. worker_index .. [[ (
                              extra_bonus_index UNSIGNED NOT NULL PRIMARY KEY REFERENCES extra_bonuses(ix) ON DELETE CASCADE,
                              offset            UNSIGNED NOT NULL DEFAULT 0
                             )
                             WITH ENGINE = 'memtx';]])
end
--************************************************************************************************************************************
-- pre_stakings
for worker_index=0,%[2]v do
        box.execute([[CREATE TABLE IF NOT EXISTS pre_stakings_]] .. worker_index .. [[ (
                                                           created_at   UNSIGNED NOT NULL,
                                                           user_id      STRING NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                                                           years        UNSIGNED NOT NULL REFERENCES pre_staking_bonuses(years),
                                                           allocation   UNSIGNED NOT NULL CHECK (allocation > 0 AND allocation <= 100),
                                                           PRIMARY KEY (user_id, years, allocation)
                                                          )
                                                        WITH ENGINE = 'memtx';]])
        box.execute([[CREATE INDEX IF NOT EXISTS pre_stakings_]] .. worker_index .. [[_years_idx ON pre_stakings_]] .. worker_index .. [[ (years);]])
end
--************************************************************************************************************************************
-- balance_recalculation_workers
for worker_index=0,%[2]v do
    	box.execute([[CREATE TABLE IF NOT EXISTS balance_recalculation_worker_]] .. worker_index .. [[
                           (
                            last_iteration_finished_at UNSIGNED,
                            last_mining_started_at     UNSIGNED,
                            last_mining_ended_at       UNSIGNED,
                            enabled                    BOOLEAN NOT NULL DEFAULT FALSE,
                            user_id                    STRING PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE
                           )
                            WITH ENGINE = 'memtx';]])
        box.execute([[CREATE INDEX IF NOT EXISTS balance_recalculation_worker_]] .. worker_index .. [[_iterator_ix ON balance_recalculation_worker_]] .. worker_index .. [[ (enabled,last_iteration_finished_at);]])
        box.execute([[CREATE INDEX IF NOT EXISTS balance_recalculation_worker_]] .. worker_index .. [[_iterator2_ix ON balance_recalculation_worker_]] .. worker_index .. [[ (user_id,last_iteration_finished_at);]])
end
--************************************************************************************************************************************
-- mining_rates_recalculation_workers
for worker_index=0,%[2]v do
        box.execute([[CREATE TABLE IF NOT EXISTS mining_rates_recalculation_worker_]] .. worker_index .. [[
                                   (
                                    last_iteration_finished_at UNSIGNED,
                                    user_id                    STRING PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE
                                   )
                                    WITH ENGINE = 'memtx';]])
        box.execute([[CREATE INDEX IF NOT EXISTS mining_rates_recalculation_worker_]] .. worker_index .. [[_last_iteration_finished_at_ix ON mining_rates_recalculation_worker_]] .. worker_index .. [[ (last_iteration_finished_at);]])
end
--************************************************************************************************************************************
-- blockchain_balance_synchronization_workers
for worker_index=0,%[2]v do
        box.execute([[CREATE TABLE IF NOT EXISTS blockchain_balance_synchronization_worker_]] .. worker_index .. [[
                                   (
                                    last_iteration_finished_at        UNSIGNED,
                                    mining_blockchain_account_address STRING,
                                    user_id                           STRING PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE
                                   )
                                    WITH ENGINE = 'memtx';]])
        box.execute([[CREATE INDEX IF NOT EXISTS blockchain_balance_synchronization_worker_]] .. worker_index .. [[_last_iteration_finished_at_ix ON blockchain_balance_synchronization_worker_]] .. worker_index .. [[ (last_iteration_finished_at);]])
end
--************************************************************************************************************************************
-- extra_bonus_processing_workers
for worker_index=0,%[2]v do
    	box.execute([[CREATE TABLE IF NOT EXISTS extra_bonus_processing_worker_]] .. worker_index .. [[
                           (
                            extra_bonus_started_at          UNSIGNED,
                            extra_bonus_ended_at            UNSIGNED,
                            user_id                         STRING NOT NULL PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
                            utc_offset                      INT NOT NULL DEFAULT 0,
                            news_seen                       UNSIGNED NOT NULL DEFAULT 0,
                            extra_bonus                     UNSIGNED NOT NULL DEFAULT 0,
                            last_extra_bonus_index_notified UNSIGNED REFERENCES extra_bonuses(ix) ON DELETE SET NULL
                           )
                            WITH ENGINE = 'memtx';]])
        box.execute([[CREATE INDEX IF NOT EXISTS extra_bonus_processing_worker_]] .. worker_index .. [[_iterator_ix ON extra_bonus_processing_worker_]] .. worker_index .. [[ (last_extra_bonus_index_notified);]])
end
--************************************************************************************************************************************
-- balances
for worker_index=0,%[2]v do
        box.execute([[CREATE TABLE IF NOT EXISTS balances_]] .. worker_index .. [[
                        (
                           updated_at  UNSIGNED NOT NULL,
                           amount      STRING NOT NULL DEFAULT '0',
                           user_id     STRING NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                           type_detail STRING NOT NULL DEFAULT '',
                           type        UNSIGNED NOT NULL,
                           negative    BOOLEAN NOT NULL DEFAULT FALSE,
                           PRIMARY KEY (user_id, negative, type, type_detail)
                        )
                         WITH ENGINE = 'memtx';]])
end