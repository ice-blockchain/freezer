-- SPDX-License-Identifier: ice License 1.0
WITH valid_users AS (
                    SELECT DISTINCT ON (id, created_at)
                            created_at, id, id_t0, id_tminus1, pre_staking_allocation, pre_staking_bonus, balance_solo, balance_solo_ethereum, balance_t0, balance_t0_ethereum, balance_for_t0, balance_t1_ethereum, balance_for_tminus1, balance_t2_ethereum
                    FROM %[1]v
                    WHERE created_at IN ['%[2]v']
                      AND kyc_step_passed >= %[3]v
                      AND (kyc_step_blocked = 0 OR kyc_step_blocked >= %[3]v + 1)
                    ),
  valid_t1_users AS (
                    SELECT created_at, id_t0, SUM(balance_for_t0) AS balance_t1
                    FROM valid_users
                    GROUP BY created_at, id_t0
                    ),
  valid_t2_users AS (
                    SELECT created_at, id_tminus1, SUM(balance_for_tminus1) AS balance_t2
                    FROM valid_users
                    GROUP BY created_at, id_tminus1
                    )
SELECT  u.created_at                                                                                                                                                AS created_at,
        SUM(((u.balance_solo+IF(t0.id != 0, u.balance_t0, 0)+t1.balance_t1+t2.balance_t2) * (100.0 - u.pre_staking_allocation)) / 100.0)                            AS balance_total_standard,
        SUM(((u.balance_solo+IF(t0.id != 0, u.balance_t0, 0)+t1.balance_t1+t2.balance_t2) * (100.0 + u.pre_staking_bonus) * u.pre_staking_allocation) / 10000.0)    AS balance_total_pre_staking,
        SUM(u.balance_solo_ethereum+u.balance_t0_ethereum+u.balance_t1_ethereum+u.balance_t2_ethereum)                                                              AS balance_total_ethereum
FROM valid_users u
       GLOBAL LEFT JOIN valid_users t0
                     ON t0.id = u.id_t0
                    AND t0.created_at = u.created_at
       GLOBAL LEFT JOIN valid_t1_users t1
                     ON t1.id_t0 = u.id
                    AND t1.created_at = u.created_at
       GLOBAL LEFT JOIN valid_t2_users t2
                     ON t2.id_tminus1 = u.id
                    AND t2.created_at = u.created_at
GROUP BY u.created_at