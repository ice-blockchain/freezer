-- SPDX-License-Identifier: ice License 1.0
CREATE TABLE IF NOT EXISTS pending_coin_distributions  (
                    created_at                timestamp NOT NULL,
                    internal_id               bigint    NOT NULL,
                    iceflakes                 bigint    NOT NULL,
                    user_id                   text      NOT NULL PRIMARY KEY,
                    eth_address               text      NOT NULL);

CREATE INDEX IF NOT EXISTS pending_coin_distributions_worker_number_ix ON pending_coin_distributions ((internal_id % 10), created_at ASC);

--- Flow:
--infinite loop: -- with 30 sec sleep between iterations if 0 rows returned
--do in transaction:
--1.      SELECT *
--        FROM pending_coin_distributions
--        WHERE internal_id % 10 = $1
--        ORDER BY created_at ASC
--        LIMIT $2
--        FOR UPDATE
--2.      call ERC-20 smart contract method to airdrop coins
--3.      delete from pending_coin_distributions WHERE user_id = ANY($1)
-- if transaction fails you retry infinitely
-- log.info every successful transaction, log every error
