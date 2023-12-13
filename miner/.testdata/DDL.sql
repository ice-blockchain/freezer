-- SPDX-License-Identifier: ice License 1.0
--************************************************************************************************************************************
-- balance_recalculation_metrics
DROP TABLE IF EXISTS balance_recalculation_metrics;
CREATE TABLE IF NOT EXISTS balance_recalculation_metrics (
                        started_at                  timestamp NOT NULL,
                        ended_at                    timestamp NOT NULL,
                        t1_balance_positive         DOUBLE PRECISION NOT NULL,
                        t1_balance_negative         DOUBLE PRECISION NOT NULL,
                        t2_balance_positive         DOUBLE PRECISION NOT NULL,
                        t2_balance_negative         DOUBLE PRECISION NOT NULL,
                        t1_active_counts_positive   BIGINT NOT NULL,
                        t1_active_counts_negative   BIGINT NOT NULL,
                        t2_active_counts_positive   BIGINT NOT NULL,
                        t2_active_counts_negative   BIGINT NOT NULL,
                        iterations_num              BIGINT NOT NULL,
                        affected_users              BIGINT NOT NULL,
                        worker                      BIGINT NOT NULL PRIMARY KEY
                    ) WITH (fillfactor = 70);
--************************************************************************************************************************************
-- balance_recalculation_dry_run
DROP TABLE IF EXISTS balance_recalculation_dry_run;
CREATE TABLE IF NOT EXISTS balance_recalculation_dry_run (
                        diff_t1_balance             DOUBLE PRECISION NOT NULL,
                        diff_t2_balance             DOUBLE PRECISION NOT NULL,
                        diff_t1_active_counts       DOUBLE PRECISION NOT NULL,
                        diff_t2_active_counts       DOUBLE PRECISION NOT NULL,
                        user_id                     text PRIMARY KEY
                    ) WITH (fillfactor = 70);
--************************************************************************************************************************************
