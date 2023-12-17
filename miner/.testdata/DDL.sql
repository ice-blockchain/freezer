-- SPDX-License-Identifier: ice License 1.0

CREATE TABLE IF NOT EXISTS balance_tminus1_recalculation_dry_run (
                        updated_at                      timestamp NOT NULL,
                        old_tminus1_balance             DOUBLE PRECISION NOT NULL,
                        new_tminus1_balance             DOUBLE PRECISION NOT NULL,
                        user_id                         text PRIMARY KEY
                    ) WITH (fillfactor = 70);

CREATE TABLE IF NOT EXISTS balance_t2_recalculation_dry_run (
                        updated_at                      timestamp NOT NULL,
                        old_t2_balance                  DOUBLE PRECISION NOT NULL,
                        new_t1_balance                  DOUBLE PRECISION NOT NULL,
                        user_id                         text PRIMARY KEY
                    ) WITH (fillfactor = 70);