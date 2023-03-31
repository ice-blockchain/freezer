-- SPDX-License-Identifier: ice License 1.0
--************************************************************************************************************************************
-- pre_staking_bonuses
create table if not exists pre_staking_bonuses
(
    years smallint not null primary key check (years > 0),
    bonus smallint not null check (bonus > 0)
);
----
insert into pre_staking_bonuses (years, bonus)
                         values (1,     35),
                                (2,     70),
                                (3,     115),
                                (4,     170),
                                (5,     250)
on conflict(years) do nothing;
----
--************************************************************************************************************************************
-- extra_bonuses
create table if not exists extra_bonuses
(
    ix    smallint not null primary key check (ix >= 0),
    bonus smallint not null default 0 check (bonus >= 0)
);
----
insert into extra_bonuses (ix, bonus)
                    values %[3]v
on conflict(ix) do nothing;
----
--************************************************************************************************************************************
-- global
create table if not exists global
(
    key   text not null primary key,
    value bigint not null
);
----
--************************************************************************************************************************************
-- extra_bonus_start_date
create table if not exists extra_bonus_start_date
(
    key smallint not null primary key check (key >= 0),
    value bigint not null check (value > 0)
);
----
insert into extra_bonus_start_date (key, value)
                            values (0,   %[4]v)
on conflict(key) do nothing;
----
--************************************************************************************************************************************
-- adoption
create table if not exists adoption
(
    achieved_at             timestamp,
    base_mining_rate        text not null default '0',
    milestone               smallint not null primary key check (milestone > 0),
    total_active_users      bigint not null check (total_active_users >= 0)
);
----
insert into adoption (milestone, total_active_users, base_mining_rate, achieved_at)
              values (1,         0,                  '16000000000',          '%[1]v'),
                     (2,         %[5]v,              '8000000000',            null),
                     (3,         %[6]v,              '4000000000',            null),
                     (4,         %[7]v,              '2000000000',            null),
                     (5,         %[8]v,              '1000000000',            null),
                     (6,         %[9]v,              '500000000',             null)
on conflict(milestone) do nothing;
----
--************************************************************************************************************************************
-- users
create table if not exists users
(
    created_at                                              timestamp not null,
    updated_at                                              timestamp not null,
    rollback_used_at                                        timestamp,
    last_natural_mining_started_at                          timestamp,
    last_mining_started_at                                  timestamp,
    last_mining_ended_at                                    timestamp,
    previous_mining_started_at                              timestamp,
    previous_mining_ended_at                                timestamp,
    last_free_mining_session_awarded_at                     timestamp,
    user_id                                                 text not null primary key,
    referred_by                                             text,
    username                                                text,
    first_name                                              text,
    last_name                                               text,
    profile_picture_name                                    text,
    mining_blockchain_account_address                       text,
    blockchain_account_address                              text,
    hash_code                                               bigint not null,
    hide_ranking                                            boolean not null default false,
    verified                                                boolean not null default false
);
----
create index if not exists users_referred_by_idx ON users (referred_by);
----
create index if not exists top_miners_lookup_idx ON users (username,first_name,last_name);
----
--************************************************************************************************************************************
-- balances
create table if not exists balances
(
    amount     text not null default '0',
    amount_w0  bigint not null default 0,
    amount_w1  bigint not null default 0,
    amount_w2  bigint not null default 0,
    amount_w3  bigint not null default 0,
    user_id    text not null primary key references users(user_id) on delete cascade
);
----
create index if not exists balances_amount_words_ix ON balances (amount_w3, amount_w2, amount_w1, amount_w0);
----
--************************************************************************************************************************************
-- processed_add_balance_commands
create table if not exists processed_add_balance_commands
(
    user_id text not null references users(user_id) on delete cascade,
    key     text not null,
    primary key (user_id, key)
);
----
--************************************************************************************************************************************
-- processed_seen_news
create table if not exists processed_seen_news
(
    user_id text not null references users(user_id) on delete cascade,
    news_id text not null,
    primary key (user_id, news_id)
);
----
--************************************************************************************************************************************
-- processed_mining_sessions
create table if not exists processed_mining_sessions
(
    user_id         text not null references users(user_id) on delete cascade,
    session_number  integer not null check (session_number >= 0),
    negative        boolean not null default false,
    primary key(session_number, negative, user_id)
);
----
--************************************************************************************************************************************
-- functions
CREATE OR REPLACE FUNCTION createListWorkerPartition(tableName text, count smallint)
  RETURNS VOID AS
$$
BEGIN
    FOR worker_index IN 0 .. count-1 BY 1
    LOOP
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %%s_%%s PARTITION OF %%s FOR VALUES WITH (MODULUS %%s,REMAINDER %%s);',
           tableName,
           worker_index,
           tableName,
           count,
           worker_index
        );
    END LOOP;
END
$$ LANGUAGE plpgsql;
----
CREATE OR REPLACE FUNCTION createListWorkerPartition(tableName text, count smallint)
  RETURNS VOID AS
$$
BEGIN
    FOR worker_index IN 0 .. count-1 BY 1
    LOOP
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %%s_%%s PARTITION OF %%s FOR VALUES IN (%%s);',
           tableName,
           worker_index,
           tableName,
           worker_index
        );
    END LOOP;
END
$$ LANGUAGE plpgsql;
----
--************************************************************************************************************************************
-- mining_sessions_dlq
create table if not exists mining_sessions_dlq
(
    id              text not null,
    user_id         text not null references users(user_id) on delete cascade,
    message         text not null,
    hash_code       bigint not null,
    worker_index    smallint not null check (worker_index >= 0),
    primary key(worker_index, id)
) partition by list (worker_index);
----
select createListWorkerPartition('mining_sessions_dlq'::text,%[2]v::smallint);
----
--************************************************************************************************************************************
-- extra_bonuses_worker
create table if not exists extra_bonuses_worker
(
    extra_bonus_index smallint not null references extra_bonuses(ix) on delete cascade,
    offset_value      smallint not null default 0 check (offset_value >= 0),
    worker_index      smallint not null check (worker_index >= 0),
    primary key(worker_index, extra_bonus_index)
) partition by list (worker_index);
----
select createListWorkerPartition('extra_bonuses_worker'::text,%[2]v::smallint);
----
--************************************************************************************************************************************
-- pre_stakings
create table if not exists pre_stakings
(
    created_at   timestamp not null,
    user_id      text not null references users(user_id) on delete cascade,
    years        smallint not null references pre_staking_bonuses(years),
    allocation   smallint not null check (allocation > 0 AND allocation <= 100),
    hash_code    bigint not null,
    worker_index smallint not null check (worker_index >= 0),
    primary key (worker_index, user_id, years, allocation)
) partition by list (worker_index);
----
create index if not exists pre_stakings_years_idx ON pre_stakings(worker_index,years);
----
select createListWorkerPartition('pre_stakings'::text,%[2]v::smallint);
----
--************************************************************************************************************************************
-- balance_recalculation_workers
create table if not exists balance_recalculation_worker
(
    last_iteration_finished_at timestamp,
    last_mining_started_at     timestamp,
    last_mining_ended_at       timestamp,
    enabled                    boolean not null default false,
    user_id                    text not null references users(user_id) on delete cascade,
    hash_code                  bigint not null,
    worker_index               smallint not null check (worker_index >= 0),
    primary key (worker_index, user_id)
) partition by list (worker_index);
----
create index if not exists balance_recalculation_worker_iterator_ix ON balance_recalculation_worker(worker_index,enabled,last_iteration_finished_at);
----
select createListWorkerPartition('balance_recalculation_worker'::text,%[2]v::smallint);
----
--************************************************************************************************************************************
-- mining_rates_recalculation_workers
create table if not exists mining_rates_recalculation_worker
(
    last_iteration_finished_at timestamp,
    user_id                    text not null references users(user_id) on delete cascade,
    hash_code                  bigint not null,
    worker_index               smallint not null check (worker_index >= 0),
    primary key (worker_index, user_id)
) partition by list (worker_index);
----
create index if not exists mining_rates_recalculation_worker_last_iteration_finished_at_ix ON mining_rates_recalculation_worker(worker_index,last_iteration_finished_at);
----
select createListWorkerPartition('mining_rates_recalculation_worker'::text,%[2]v::smallint);
----
--************************************************************************************************************************************
-- blockchain_balance_synchronization_workers
create table if not exists blockchain_balance_synchronization_worker
(
    last_iteration_finished_at        timestamp,
    mining_blockchain_account_address text,
    user_id                           text not null references users(user_id) on delete cascade,
    hash_code                         bigint not null,
    worker_index                      smallint not null check (worker_index >= 0),
    primary key (worker_index, user_id)
) partition by list (worker_index);
----
create index if not exists blockchain_balance_synchronization_worker_last_iteration_finished_at_ix ON blockchain_balance_synchronization_worker(worker_index,last_iteration_finished_at);
----
select createListWorkerPartition('blockchain_balance_synchronization_worker'::text,%[2]v::smallint);
----
--************************************************************************************************************************************
-- extra_bonus_processing_workers
create table if not exists extra_bonus_processing_worker
(
    extra_bonus_started_at          timestamp,
    extra_bonus_ended_at            timestamp,
    user_id                         text not null references users(user_id) on delete cascade,
    utc_offset                      smallint not null default 0,
    news_seen                       smallint not null default 0 check (news_seen >= 0),
    extra_bonus                     smallint not null default 0 check (extra_bonus >= 0),
    last_extra_bonus_index_notified smallint references extra_bonuses(ix) on delete set null,
    hash_code                       bigint not null,
    worker_index                    smallint not null check (worker_index >= 0),
    primary key (worker_index, user_id)
) partition by list (worker_index);
----
create index if not exists extra_bonus_processing_worker_iterator_ix ON extra_bonus_processing_worker(worker_index,last_extra_bonus_index_notified);
----
select createListWorkerPartition('extra_bonus_processing_worker'::text,%[2]v::smallint);
----
--************************************************************************************************************************************
-- balances_worker
create table if not exists balances_worker
(
    updated_at    timestamp not null,
    amount        text not null default '0',
    user_id       text not null references users(user_id) on delete cascade,
    type_detail   text not null default '',
    type          smallint not null check (type >= 0),
    negative      boolean not null default false,
    hash_code     bigint not null,
    worker_index  smallint not null check (worker_index >= 0),
    primary key (worker_index, user_id, negative, type, type_detail)
) partition by list (worker_index);
----
select createListWorkerPartition('balances_worker'::text,%[2]v::smallint);
----
--************************************************************************************************************************************
-- active_referrals
create table if not exists active_referrals
(
    user_id      text not null references users(user_id) on delete cascade,
    t1           integer not null default 0 check (t1 >= 0),
    t2           integer not null default 0 check (t2 >= 0),
    hash_code    bigint not null,
    worker_index smallint not null check (worker_index >= 0),
    primary key (worker_index, user_id)
) partition by list (worker_index);
----
select createListWorkerPartition('active_referrals'::text,%[2]v::smallint);