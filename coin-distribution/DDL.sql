-- SPDX-License-Identifier: ice License 1.0
DO $$ BEGIN
    CREATE DOMAIN uint256 AS NUMERIC(78,0) NOT NULL DEFAULT 0
    CHECK (VALUE >= 0 AND VALUE <= 115792089237316195423570985008687907853269984665640564039457584007913129639935)
    CHECK (SCALE(VALUE) = 0);
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS pending_coin_distributions  (
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    day                       date      NOT NULL,
                    iceflakes                 uint256,
                    user_id                   text      NOT NULL,
                    eth_address               text      NOT NULL,
                    PRIMARY KEY(day, user_id));

CREATE INDEX IF NOT EXISTS pending_coin_distributions_worker_number_ix ON pending_coin_distributions ((internal_id % 10), created_at ASC);

CREATE TABLE IF NOT EXISTS global (
                    key       text NOT NULL primary key,
                    value     text NOT NULL );
INSERT INTO global (key,value)
            VALUES ('coin_distributer_enabled','true')
         ON CONFLICT(key) DO NOTHING;

CREATE TABLE IF NOT EXISTS coin_distributions_by_earner (
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    balance                   bigint    NOT NULL,
                    day                       date      NOT NULL,
                    username                  text      NOT NULL,
                    referred_by_username      text      NOT NULL,
                    user_id                   text      NOT NULL,
                    earner_user_id            text      NOT NULL,
                    eth_address               text      NOT NULL,
                    PRIMARY KEY(day, user_id, earner_user_id));

CREATE TABLE IF NOT EXISTS coin_distributions_pending_review  (
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    ice                       bigint    NOT NULL,
                    day                       date      NOT NULL,
                    iceflakes                 uint256           ,
                    username                  text      NOT NULL,
                    referred_by_username      text      NOT NULL,
                    user_id                   text      NOT NULL,
                    eth_address               text      NOT NULL,
                    PRIMARY KEY(day, user_id));

CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_internal_id_ix ON coin_distributions_pending_review (internal_id);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_created_at_ix ON coin_distributions_pending_review (created_at);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_ice_ix ON coin_distributions_pending_review (ice);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_username_ix ON coin_distributions_pending_review (username);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_referred_by_username_ix ON coin_distributions_pending_review (referred_by_username);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_lookup1_ix ON coin_distributions_pending_review (ice,internal_id);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_lookup2_ix ON coin_distributions_pending_review (created_at,internal_id);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_lookup3_ix ON coin_distributions_pending_review (username,internal_id);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_lookup4_ix ON coin_distributions_pending_review (ice,username,internal_id);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_lookup5_ix ON coin_distributions_pending_review (referred_by_username,internal_id);
CREATE INDEX IF NOT EXISTS coin_distributions_pending_review_lookup6_ix ON coin_distributions_pending_review (ice,referred_by_username,internal_id);

CREATE TABLE IF NOT EXISTS reviewed_coin_distributions  (
                    reviewed_at               timestamp NOT NULL,
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    ice                       bigint    NOT NULL,
                    day                       date      NOT NULL,
                    review_day                date      NOT NULL,
                    iceflakes                 uint256           ,
                    username                  text      NOT NULL,
                    referred_by_username      text      NOT NULL,
                    user_id                   text      NOT NULL ,
                    eth_address               text      NOT NULL,
                    reviewer_user_id          text      NOT NULL,
                    decision                  text      NOT NULL,
                    PRIMARY KEY(user_id, day, review_day));