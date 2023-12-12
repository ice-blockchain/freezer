-- SPDX-License-Identifier: ice License 1.0
CREATE TABLE IF NOT EXISTS pending_coin_distributions  (
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    iceflakes                 bigint    NOT NULL,
                    user_id                   text      NOT NULL PRIMARY KEY,
                    eth_address               text      NOT NULL);

CREATE INDEX IF NOT EXISTS pending_coin_distributions_worker_number_ix ON pending_coin_distributions ((internal_id % 10), created_at ASC);

CREATE TABLE IF NOT EXISTS pending_coin_distribution_configurations (
                    key       text NOT NULL primary key,
                    value     text NOT NULL );
INSERT INTO pending_coin_distribution_configurations(key,value) VALUES ('enabled','true') ON CONFLICT(key) DO NOTHING;

--- Flow:
--infinite loop: -- with 30 sec sleep between iterations if 0 rows returned
--do in transaction:
--1.      SELECT *
--        FROM pending_coin_distributions
--        WHERE internal_id % 10 = $1
--        ORDER BY created_at ASC
--        LIMIT $2
--        FOR UPDATE
--2.      delete from pending_coin_distributions WHERE user_id = ANY($1)
--3.      call ERC-20 smart contract method to airdrop coins

DO $$ BEGIN
    CREATE DOMAIN uint256 AS NUMERIC(78,0) NOT NULL DEFAULT 0
    CHECK (VALUE >= 0 AND VALUE <= 115792089237316195423570985008687907853269984665640564039457584007913129639935)
    CHECK (SCALE(VALUE) = 0);
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS coin_distributions_by_earner (
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    balance                   bigint    NOT NULL,
                    username                  text      NOT NULL,
                    referred_by_username      text      NOT NULL,
                    user_id                   text      NOT NULL,
                    earner_user_id            text      NOT NULL,
                    eth_address               text      NOT NULL,
                    PRIMARY KEY(user_id, earner_user_id));

CREATE TABLE IF NOT EXISTS coin_distributions_pending_review  (
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    iceflakes                 uint256,
                    username                  text      NOT NULL,
                    referred_by_username      text      NOT NULL,
                    user_id                   text      NOT NULL PRIMARY KEY,
                    eth_address               text      NOT NULL);

CREATE TABLE IF NOT EXISTS reviewed_coin_distributions  (
                    reviewed_at               timestamp NOT NULL,
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    iceflakes                 uint256,
                    username                  text      NOT NULL,
                    referred_by_username      text      NOT NULL,
                    user_id                   text      NOT NULL ,
                    eth_address               text      NOT NULL,
                    reviewer_user_id          text      NOT NULL,
                    decision                  text      NOT NULL,
                    PRIMARY KEY(user_id, reviewed_at));

CREATE TABLE IF NOT EXISTS pending_coin_distribution_statistics  (
                    created_at                          timestamp NOT NULL PRIMARY KEY,
                    total_iceflakes                     uint256,
                    marketing_25percent_iceflakes       uint256,
                    marketing_45percent_iceflakes       uint256,
                    marketing_30percent_iceflakes       uint256,
                    marketing_25percent_eth_address     text NOT NULL,
                    marketing_45percent_eth_address     text NOT NULL,
                    marketing_30percent_eth_address     text NOT NULL);