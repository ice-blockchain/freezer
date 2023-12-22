-- SPDX-License-Identifier: ice License 1.0
SELECT u.created_at as created_at,
       sum((u.balance_solo + (if(t0.id != 0, u.balance_t0, 0)) + verified_balance_t1.balance + verified_balance_t2.balance)*(100.0-u.pre_staking_allocation)/100.0)  AS balance_total_standard,
       sum((u.balance_solo + (if(t0.id != 0, u.balance_t0, 0)) + verified_balance_t1.balance + verified_balance_t2.balance) * (100 + u.pre_staking_bonus) * u.pre_staking_allocation / 10000) AS balance_total_pre_staking,
       sum(u.balance_solo_ethereum + (if(t0.id != 0, u.balance_t0_ethereum, 0)) + verified_balance_t1.ethereum + verified_balance_t2.ethereum) AS balance_total_ethereum
FROM %[1]v u
    GLOBAL LEFT JOIN
    (select DISTINCT ON (id, created_at) id, created_at
    from freezer_user_history
    where created_at IN ['%[2]v']
    AND kyc_step_passed >= %[3]v AND (kyc_step_blocked = 0 OR kyc_step_blocked >= (%[3]v+1))
    group by id, created_at) as t0
ON t0.id = u.id_t0 AND t0.created_at = u.created_at
    GLOBAL LEFT JOIN
    (SELECT DISTINCT ON (id_t0, created_at) id_t0, created_at, sum(balance_for_t0) AS balance, sum(balance_for_t0_ethereum) AS ethereum
    FROM %[1]v
    WHERE created_at IN ['%[2]v']
    AND kyc_step_passed >= %[3]v AND (kyc_step_blocked = 0 OR kyc_step_blocked >= (%[3]v+1))
    GROUP BY id_t0, created_at) AS verified_balance_t1
    ON verified_balance_t1.id_t0 = u.id AND verified_balance_t1.created_at = u.created_at
    GLOBAL LEFT JOIN
    (SELECT DISTINCT ON (id_tminus1, created_at) id_tminus1, created_at, sum(balance_for_tminus1) AS balance, sum(balance_for_tminus1_ethereum) AS ethereum
    FROM %[1]v
    WHERE created_at IN ['%[2]v']
    AND kyc_step_passed >= %[3]v AND ( kyc_step_blocked = 0 OR kyc_step_blocked >= (%[3]v+1))
    GROUP BY id_tminus1, created_at) AS verified_balance_t2 ON verified_balance_t2.id_tminus1 = u.id AND verified_balance_t2.created_at = u.created_at
WHERE u.created_at IN ['%[2]v']
  AND u.kyc_step_passed >= %[3]v AND ( u.kyc_step_blocked = 0 OR u.kyc_step_blocked >= (%[3]v+1))
GROUP BY u.created_at