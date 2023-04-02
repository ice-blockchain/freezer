// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/terror"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) StartNewMiningSession( //nolint:funlen,gocognit // A lot of handling.
	ctx context.Context, ms *MiningSummary, rollbackNegativeMiningProgress *bool,
) error {
	userID := *ms.MiningSession.UserID
	old, err := r.getInternalMiningSummary(ctx, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to getMiningSummary for userID:%v", userID)
	}
	now := time.Now()
	nowValue := *(*now).Time
	if old.LastMiningEndedAt != nil &&
		old.LastNaturalMiningStartedAt != nil &&
		old.LastMiningEndedAt.After(*now.Time) &&
		(now.Sub(*old.LastNaturalMiningStartedAt.Time)/r.cfg.MiningSessionDuration.Min)%2 == 0 {
		return ErrDuplicate
	}
	shouldRollback, err := r.validateRollbackNegativeMiningProgress(old, now, rollbackNegativeMiningProgress)
	if err != nil {
		return err
	}
	newMS := r.newMiningSummary(old, now)
	if err = r.insertNewMiningSession(ctx, userID, old, newMS, shouldRollback); err != nil {
		return errors.Wrapf(err, "failed to insertNewMiningSession:%#v,userID:%v,rollbackNegativeMiningProgress:%v", newMS, userID, shouldRollback)
	}

	return errors.Wrapf(retry(ctx, func() error {
		summary, gErr := r.GetMiningSummary(ctx, userID)
		if gErr == nil {
			if summary.MiningSession == nil || summary.MiningSession.StartedAt == nil || summary.MiningSession.StartedAt.UnixMicro() != nowValue.UnixMicro() {
				gErr = ErrNotFound
			} else {
				*ms = *summary
			}
		}

		return gErr
	}), "permanently failed to GetMiningSummary for userID:%v", userID)
}

func (r *repository) getInternalMiningSummary(ctx context.Context, userID string) (*miningSummary, error) { //nolint:funlen // Big SQL.
	sql := fmt.Sprintf(`SELECT u.last_natural_mining_started_at,
							   u.last_mining_started_at,
							   u.last_mining_ended_at,
							   u.previous_mining_started_at,
							   u.previous_mining_ended_at,
							   u.last_free_mining_session_awarded_at,
							   MAX(negative_balance.amount) as negative_total_no_pre_staking_bonus_balance_amount,
							   MAX(negative_t0_balance.amount) as negative_total_t0_no_pre_staking_bonus_balance_amount,
							   MAX(negative_t1_balance.amount) as negative_total_t1_no_pre_staking_bonus_balance_amount,
							   MAX(negative_t2_balance.amount) as negative_total_t2_no_pre_staking_bonus_balance_amount,
							   0 AS mining_streak,
							   COALESCE(MAX(st.years),0) AS pre_staking_years,
							   COALESCE(MAX(st.allocation),0) AS pre_staking_allocation,
							   COALESCE(MAX(st_b.bonus),0) as pre_staking_bonus
						FROM users u 
							LEFT JOIN pre_stakings st
								   ON st.worker_index = $1
								  AND st.user_id = u.user_id
							LEFT JOIN pre_staking_bonuses st_b
								   ON st.worker_index = $1
								  AND st.years = st_b.years
							LEFT JOIN balances_worker negative_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_balance.worker_index = $1
								  AND negative_balance.user_id = u.user_id
								  AND negative_balance.negative = TRUE
								  AND negative_balance.type = %[1]v
								  AND negative_balance.type_detail = ''
							LEFT JOIN balances_worker negative_t0_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t0_balance.worker_index = $1
								  AND negative_t0_balance.user_id = u.user_id
								  AND negative_t0_balance.negative = TRUE
								  AND negative_t0_balance.type = %[1]v
								  AND negative_t0_balance.type_detail = '%[2]v_' || u.referred_by
							LEFT JOIN balances_worker negative_t1_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t1_balance.worker_index = $1
								  AND negative_t1_balance.user_id = u.user_id
								  AND negative_t1_balance.negative = TRUE
								  AND negative_t1_balance.type = %[1]v
								  AND negative_t1_balance.type_detail = '%[3]v'
							LEFT JOIN balances_worker negative_t2_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t2_balance.worker_index = $1
								  AND negative_t2_balance.user_id = u.user_id
								  AND negative_t2_balance.negative = TRUE
								  AND negative_t2_balance.type = %[1]v
								  AND negative_t2_balance.type_detail = '%[4]v'
						WHERE u.user_id = $2
						GROUP BY u.user_id`, totalNoPreStakingBonusBalanceType, t0BalanceTypeDetail, t1BalanceTypeDetail, t2BalanceTypeDetail)
	resp, err := storage.Get[miningSummary](ctx, r.db, sql, r.workerIndex(ctx), userID)
	if err != nil && storage.IsErr(err, storage.ErrNotFound) {
		return nil, ErrRelationNotFound
	}

	return resp, errors.Wrapf(err, "failed to get the current mining summary for userID:%v", userID)
}

func (r *repository) validateRollbackNegativeMiningProgress(
	currentMiningSummary *miningSummary, now *time.Time, rollbackNegativeMiningProgress *bool,
) (*bool, error) {
	if currentMiningSummary.LastMiningEndedAt == nil {
		return nil, nil //nolint:nilnil // Nope.
	}
	amountLost := currentMiningSummary.calculateAmountLost()
	if !amountLost.IsZero() &&
		(now.Sub(*currentMiningSummary.LastMiningEndedAt.Time) < r.cfg.RollbackNegativeMining.Available.After ||
			now.Sub(*currentMiningSummary.LastMiningEndedAt.Time) > r.cfg.RollbackNegativeMining.Available.Until) {
		amountLost = nil
	}
	if rollbackNegativeMiningProgress == nil && !amountLost.IsZero() {
		return nil, terror.New(ErrNegativeMiningProgressDecisionRequired, map[string]any{
			"amount":                amountLost.UnsafeICE(),
			"duringTheLastXSeconds": now.Sub(*currentMiningSummary.LastMiningEndedAt.Time).Milliseconds() / 1e3, //nolint:gomnd // To get to seconds.
		})
	} else if rollbackNegativeMiningProgress != nil && amountLost.IsZero() {
		return nil, nil //nolint:nilnil // Nope.
	}

	return rollbackNegativeMiningProgress, nil
}

func (m *miningSummary) calculateAmountLost() *coin.ICEFlake {
	standardAmount := m.NegativeTotalNoPreStakingBonusBalanceAmount.
		MultiplyUint64(percentage100 - m.PreStakingAllocation).
		DivideUint64(percentage100)
	preStakingAmount := m.NegativeTotalNoPreStakingBonusBalanceAmount.
		MultiplyUint64(m.PreStakingAllocation * (m.PreStakingBonus + percentage100)).
		DivideUint64(percentage100 * percentage100)
	standardT0Amount := m.NegativeTotalT0NoPreStakingBonusBalanceAmount.
		MultiplyUint64(percentage100 - m.PreStakingAllocation).
		DivideUint64(percentage100)
	preStakingT0Amount := m.NegativeTotalT0NoPreStakingBonusBalanceAmount.
		MultiplyUint64(m.PreStakingAllocation * (m.PreStakingBonus + percentage100)).
		DivideUint64(percentage100 * percentage100)
	standardT1Amount := m.NegativeTotalT1NoPreStakingBonusBalanceAmount.
		MultiplyUint64(percentage100 - m.PreStakingAllocation).
		DivideUint64(percentage100)
	preStakingT1Amount := m.NegativeTotalT1NoPreStakingBonusBalanceAmount.
		MultiplyUint64(m.PreStakingAllocation * (m.PreStakingBonus + percentage100)).
		DivideUint64(percentage100 * percentage100)
	standardT2Amount := m.NegativeTotalT2NoPreStakingBonusBalanceAmount.
		MultiplyUint64(percentage100 - m.PreStakingAllocation).
		DivideUint64(percentage100)
	preStakingT2Amount := m.NegativeTotalT2NoPreStakingBonusBalanceAmount.
		MultiplyUint64(m.PreStakingAllocation * (m.PreStakingBonus + percentage100)).
		DivideUint64(percentage100 * percentage100)

	return standardAmount.Add(preStakingAmount).
		Add(standardT0Amount).Add(preStakingT0Amount).
		Add(standardT1Amount).Add(preStakingT1Amount).
		Add(standardT2Amount).Add(preStakingT2Amount)
}

func (r *repository) newMiningSummary(old *miningSummary, now *time.Time) *miningSummary {
	resp := &miningSummary{
		LastMiningStartedAt:        now,
		LastNaturalMiningStartedAt: now,
		LastMiningEndedAt:          time.New(now.Add(r.cfg.MiningSessionDuration.Max)),
	}
	if old.LastMiningEndedAt == nil || old.LastMiningStartedAt == nil || old.LastMiningEndedAt.Before(*now.Time) {
		return resp
	}
	resp.LastMiningStartedAt = old.LastMiningStartedAt
	resp.LastFreeMiningSessionAwardedAt = old.LastFreeMiningSessionAwardedAt
	resp.MiningStreak = r.calculateMiningStreak(now, resp.LastMiningStartedAt, resp.LastMiningEndedAt)
	var durationSinceLastFreeMiningSessionAwarded stdlibtime.Duration
	if resp.LastFreeMiningSessionAwardedAt == nil {
		durationSinceLastFreeMiningSessionAwarded = now.Sub(*resp.LastMiningStartedAt.Time)
	} else {
		durationSinceLastFreeMiningSessionAwarded = now.Sub(*resp.LastFreeMiningSessionAwardedAt.Time)
	}
	freeMiningSession := uint64(0)
	minimumDurationForAwardingFreeMiningSession := stdlibtime.Duration(r.cfg.ConsecutiveNaturalMiningSessionsRequiredFor1ExtraFreeArtificialMiningSession.Max) * r.cfg.MiningSessionDuration.Max //nolint:lll // .
	if durationSinceLastFreeMiningSessionAwarded >= minimumDurationForAwardingFreeMiningSession {
		resp.LastFreeMiningSessionAwardedAt = now
		freeMiningSession++
	}
	if freeSessions := stdlibtime.Duration(r.calculateRemainingFreeMiningSessions(now, old.LastMiningEndedAt) + freeMiningSession); freeSessions > 0 {
		resp.LastMiningEndedAt = time.New(resp.LastMiningEndedAt.Add(freeSessions * r.cfg.MiningSessionDuration.Max))
	}

	return resp
}

func (r *repository) insertNewMiningSession( //nolint:funlen // Big script.
	ctx context.Context, userID string, old, ms *miningSummary, rollbackNegativeMiningSession *bool,
) error {
	var rollbackUsedAt, rollbackUsedAtCondition string
	if rollbackNegativeMiningSession != nil && *rollbackNegativeMiningSession {
		rollbackUsedAt = "rollback_used_at = $1,"
		rollbackUsedAtCondition = "AND rollback_used_at IS NULL"
	}
	var previousMiningEndedAtVal *stdlibtime.Time
	if old.LastMiningEndedAt != nil {
		previousMiningEndedAtVal = old.LastMiningEndedAt.Time
	}
	var lastFreeMiningSessionAwardedAtVal *stdlibtime.Time
	if ms.LastFreeMiningSessionAwardedAt != nil {
		lastFreeMiningSessionAwardedAtVal = ms.LastFreeMiningSessionAwardedAt.Time
	}

	return storage.DoInTransaction(ctx, r.db, func(conn storage.QueryExecer) error {
		sql := fmt.Sprintf(`UPDATE users
							   SET updated_at = $1,
								   last_natural_mining_started_at = $1,
								   last_mining_started_at = $2,
								   last_mining_ended_at = $3,
								   previous_mining_started_at = (CASE WHEN last_mining_started_at = $2 THEN previous_mining_started_at ELSE last_mining_started_at END),
								   previous_mining_ended_at = (CASE WHEN last_mining_started_at = $2 THEN previous_mining_ended_at ELSE last_mining_ended_at END),
								   %[1]v							   
								   last_free_mining_session_awarded_at = $4
							   WHERE user_id = $5
								 AND COALESCE(last_mining_ended_at,'1999-01-08 04:05:06'::timestamp) = COALESCE($6,'1999-01-08 04:05:06'::timestamp)
								 %[2]v`, rollbackUsedAt, rollbackUsedAtCondition)
		if rowsAffected, err := conn.Exec(ctx, sql,
			ms.LastNaturalMiningStartedAt.Time,
			ms.LastMiningStartedAt.Time,
			ms.LastMiningEndedAt.Time,
			lastFreeMiningSessionAwardedAtVal,
			userID,
			previousMiningEndedAtVal); err != nil {
			return errors.Wrapf(err, "failed to update users for starting a new mining session for: %#v", ms)
		} else if rowsAffected.RowsAffected() != 1 {
			return ErrRaceCondition
		}
		sql = `INSERT INTO balance_recalculation_worker(user_id,enabled,last_mining_started_at,last_mining_ended_at,hash_code,worker_index) 
												 VALUES($1     ,TRUE   ,$2                    ,$3                  ,$4       ,$5)
			   ON CONFLICT(worker_index, user_id)
					DO UPDATE
						  SET enabled = EXCLUDED.enabled,
							  last_mining_started_at = EXCLUDED.last_mining_started_at,
							  last_mining_ended_at = EXCLUDED.last_mining_ended_at
						WHERE balance_recalculation_worker.enabled != EXCLUDED.enabled
						   OR coalesce(balance_recalculation_worker.last_mining_started_at,'1999-01-08 04:05:06'::timestamp) != EXCLUDED.last_mining_started_at
						   OR coalesce(balance_recalculation_worker.last_mining_ended_at,'1999-01-08 04:05:06'::timestamp) != EXCLUDED.last_mining_ended_at`
		if _, err := conn.Exec(ctx, sql,
			userID,
			ms.LastMiningStartedAt.Time,
			ms.LastMiningEndedAt.Time,
			r.hashCode(ctx),
			r.workerIndex(ctx)); err != nil {
			return errors.Wrapf(err, "failed to update balance_recalculation_worker for starting a new mining session for: %#v", ms)
		}
		sql = `INSERT INTO blockchain_balance_synchronization_worker(user_id, hash_code, worker_index) 
 															 VALUES ($1     , $2        , $3) 
 			   ON CONFLICT (worker_index, user_id) DO NOTHING`
		if _, err := conn.Exec(ctx, sql, userID, r.hashCode(ctx), r.workerIndex(ctx)); err != nil {
			return errors.Wrapf(err, "failed to insert blockchain_balance_synchronization_worker_%v for userId:%v", r.workerIndex(ctx), userID)
		}
		sql = `INSERT INTO extra_bonus_processing_worker(user_id, hash_code, worker_index) 
 										         VALUES ($1     , $2        , $3) 
 			   ON CONFLICT (worker_index, user_id) DO NOTHING`
		if _, err := conn.Exec(ctx, sql, userID, r.hashCode(ctx), r.workerIndex(ctx)); err != nil {
			return errors.Wrapf(err, "failed to insert extra_bonus_processing_worker_%v for userId:%v", r.workerIndex(ctx), userID)
		}
		sql = `INSERT INTO mining_rates_recalculation_worker(user_id, hash_code, worker_index) 
 										             VALUES ($1     , $2        , $3) 
 			   ON CONFLICT (worker_index, user_id) DO NOTHING`
		if _, err := conn.Exec(ctx, sql, userID, r.hashCode(ctx), r.workerIndex(ctx)); err != nil {
			return errors.Wrapf(err, "failed to insert mining_rates_recalculation_worker_%v for userId:%v", r.workerIndex(ctx), userID)
		}
		sess := &MiningSession{
			LastNaturalMiningStartedAt: ms.LastNaturalMiningStartedAt,
			StartedAt:                  ms.LastMiningStartedAt,
			EndedAt:                    ms.LastMiningEndedAt,
			MiningStreak:               ms.MiningStreak,
			UserID:                     &userID,
		}

		return errors.Wrapf(r.sendMiningSessionMessage(ctx, sess), "failed to sendMiningSessionMessage:%#v", sess)
	})
}

func (r *repository) sendMiningSessionMessage(ctx context.Context, ms *MiningSession) error {
	valueBytes, err := json.MarshalContext(ctx, ms)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %#v", ms)
	}
	msg := &messagebroker.Message{
		Timestamp: *ms.LastNaturalMiningStartedAt.Time,
		Headers:   map[string]string{"producer": "freezer"},
		Key:       *ms.UserID,
		Topic:     r.cfg.MessageBroker.Topics[2].Name,
		Value:     valueBytes,
	}
	responder := make(chan error, 1)
	defer close(responder)
	r.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send `%v` message to broker", msg.Topic)
}

func (s *miningSessionsTableSource) Process(ctx context.Context, msg *messagebroker.Message) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	if len(msg.Value) == 0 {
		return nil
	}
	var ms MiningSession
	if err := json.UnmarshalContext(ctx, msg.Value, &ms); err != nil {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &ms)
	}
	if ms.UserID == nil {
		return nil
	}

	return multierror.Append( //nolint:wrapcheck // Not needed.
		errors.Wrapf(s.trySwitchToNextAdoption(ctx), "failed to trySwitchToNextAdoption"),
		errors.Wrapf(s.incrementActiveReferralCountForT0AndTMinus1(ctx, &ms), "failed to incrementActiveReferralCountForT0AndTMinus1 for %#v", ms),
	).ErrorOrNil()
}

func (s *miningSessionsTableSource) incrementActiveReferralCountForT0AndTMinus1(ctx context.Context, ms *MiningSession) error { //nolint:funlen // .
	if ctx.Err() != nil || !ms.LastNaturalMiningStartedAt.Equal(*ms.StartedAt.Time) {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := `SELECT COALESCE(t0.user_id,'') AS t0_user_id,
				   COALESCE(tMinus1.user_id, '') AS t_Minus1_user_id,
				   COALESCE(t0.hash_code,0) AS t0_Hash_Code,
				   COALESCE(tMinus1.hash_code,0) AS t_Minus1_Hash_Code
			  FROM users u
			  		LEFT JOIN users t0
			  			   ON t0.user_id = u.referred_by
						  AND t0.user_id != u.user_id
			  		LEFT JOIN users tMinus1
			  			   ON tMinus1.user_id = t0.referred_by
						  AND tMinus1.user_id != t0.user_id
						  AND tMinus1.user_id != u.user_id
			  WHERE u.user_id = $1`
	//nolint:revive // Nope.
	referees, err := storage.Get[struct {
		T0UserID, TMinus1UserID     string
		T0HashCode, TMinus1HashCode int64
	}](ctx, s.db, sql, *ms.UserID)
	if err != nil || (referees.T0UserID == "" && referees.TMinus1UserID == "") {
		return errors.Wrapf(err, "failed to select for t0/t-1 information for userID:%v", *ms.UserID)
	}

	return errors.Wrapf(storage.DoInTransaction(ctx, s.db, func(conn storage.QueryExecer) error {
		sessionNumber := s.sessionNumber(ms.LastNaturalMiningStartedAt)
		sql = `INSERT INTO processed_mining_sessions(session_number, user_id) VALUES ($1, $2);`
		if _, err = conn.Exec(ctx, sql, sessionNumber, *ms.UserID); err != nil {
			return errors.Wrapf(err, "failed to insert processed_mining_sessions for userId:%v, sessionNumber:%v", *ms.UserID, sessionNumber)
		}
		if referees.T0UserID != "" {
			sql = `INSERT INTO active_referrals (t1, user_id, hash_code, worker_index) 
									     VALUES (1 , $1     , $2       , $3)
				   ON CONFLICT (worker_index, user_id)
							DO UPDATE
								  SET t1 = t1 + EXCLUDED.t1`
			t0WorkerIndex := int16(uint64(referees.T0HashCode) % uint64(s.cfg.WorkerCount))
			if _, err = conn.Exec(ctx, sql, referees.T0UserID, referees.T0HashCode, t0WorkerIndex); err != nil {
				return errors.Wrapf(err, "failed to increment t1 active_referrals for userId:%v", referees.T0UserID)
			}
		}
		if referees.TMinus1UserID != "" {
			sql = `INSERT INTO active_referrals (t2, user_id, hash_code, worker_index) 
									     VALUES (1 , $1     , $2       , $3)
				   ON CONFLICT (worker_index, user_id)
							DO UPDATE
								  SET t2 = t2 + EXCLUDED.t2`
			tMinus1WorkerIndex := int16(uint64(referees.TMinus1HashCode) % uint64(s.cfg.WorkerCount))
			if _, err = conn.Exec(ctx, sql, referees.TMinus1UserID, referees.TMinus1HashCode, tMinus1WorkerIndex); err != nil {
				return errors.Wrapf(err, "failed to increment t2 active_referrals for userId:%v", referees.TMinus1UserID)
			}
		}

		return nil
	}), "failed to execute transaction to increment active_referrals for t0&t-1, for %#v", ms)
}

type (
	userThatStoppedMining struct {
		LastMiningEndedAt               *time.Time
		UserID, T0UserID, TMinus1UserID string
		T0HashCode, TMinus1HashCode     int64
	}
)

//nolint:funlen,gocognit,revive,gomnd // .
func (r *repository) decrementActiveReferralCountForT0AndTMinus1(ctx context.Context, usersThatStoppedMining ...*userThatStoppedMining) error {
	if ctx.Err() != nil || len(usersThatStoppedMining) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	t0Referrals := make(map[string]int, len(usersThatStoppedMining))
	tMinus1Referrals := make(map[string]int, len(usersThatStoppedMining))
	t0WorkerIndexes := make(map[string]int16, len(usersThatStoppedMining))
	tMinus1WorkerIndexes := make(map[string]int16, len(usersThatStoppedMining))
	processedMiningSessionValues := make([]string, 0, len(usersThatStoppedMining))
	processedMiningSessionArgs := make([]any, 0, 2*len(usersThatStoppedMining))
	for ix, usr := range usersThatStoppedMining {
		processedMiningSessionValues = append(processedMiningSessionValues, fmt.Sprintf(`($%v,true,$%v)`, 2*ix+1, 2*ix+2))
		processedMiningSessionArgs = append(processedMiningSessionArgs, r.sessionNumber(usr.LastMiningEndedAt), usr.UserID)
		if usr.TMinus1UserID != "" {
			tMinus1Referrals[usr.TMinus1UserID]++
			tMinus1WorkerIndexes[usr.TMinus1UserID] = int16(uint64(usr.TMinus1HashCode) % uint64(r.cfg.WorkerCount))
		}
		if usr.T0UserID != "" {
			t0Referrals[usr.T0UserID]++
			t0WorkerIndexes[usr.T0UserID] = int16(uint64(usr.T0HashCode) % uint64(r.cfg.WorkerCount))
		}
	}
	ix := 0
	referralsValues := make([]string, 0, len(t0Referrals)+len(tMinus1Referrals))
	referralsArgs := make([]any, 0, 3*(len(t0Referrals)+len(tMinus1Referrals)))
	t0Conditions := make([]string, 0, len(t0Referrals))
	tMinus1Conditions := make([]string, 0, len(tMinus1Referrals))
	for userID, referralCount := range t0Referrals {
		pk := fmt.Sprintf(`($%v, $%v)`, 3*ix+1, 3*ix+2)
		referralsValues = append(referralsValues, pk)
		referralsArgs = append(referralsArgs, t0WorkerIndexes[userID], userID, referralCount)
		t0Conditions = append(t0Conditions, fmt.Sprintf(`WHEN (worker_index, user_id) = %v THEN GREATEST(t1 - $%v, 0)`, pk, 3*ix+3))
		ix++
	}
	for userID, referralCount := range tMinus1Referrals {
		pk := fmt.Sprintf(`($%v, $%v)`, 3*ix+1, 3*ix+2)
		referralsValues = append(referralsValues, pk)
		referralsArgs = append(referralsArgs, tMinus1WorkerIndexes[userID], userID, referralCount)
		tMinus1Conditions = append(tMinus1Conditions, fmt.Sprintf(`WHEN (worker_index, user_id) = %v THEN GREATEST(t2 - $%v, 0)`, pk, 3*ix+3))
		ix++
	}
	if len(referralsValues) == 0 {
		sql := fmt.Sprintf(`INSERT INTO processed_mining_sessions(session_number, negative, user_id) 
																  VALUES %[1]v`, strings.Join(processedMiningSessionValues, ","))
		if _, err := storage.Exec(ctx, r.db, sql, processedMiningSessionArgs...); err == nil || storage.IsErr(err, storage.ErrDuplicate, "pk") {
			return nil
		} else {
			return errors.Wrapf(err, "failed to insert processed_mining_sessions for args:%#v", processedMiningSessionArgs)
		}
	}
	if err := storage.DoInTransaction(ctx, r.db, func(conn storage.QueryExecer) error {
		sql := fmt.Sprintf(`INSERT INTO processed_mining_sessions(session_number, negative, user_id) 
																  VALUES %[1]v`, strings.Join(processedMiningSessionValues, ","))
		if _, err := conn.Exec(ctx, sql, processedMiningSessionArgs...); err != nil {
			return errors.Wrapf(err, "failed to insert processed_mining_sessions for args:%#v", processedMiningSessionArgs)
		}
		sql = fmt.Sprintf(`UPDATE active_referrals
						   SET t1 = (CASE %[2]v ELSE t1 END),
							   t2 = (CASE %[3]v ELSE t2 END)
						   WHERE (worker_index, user_id) in (%[1]v)`,
			strings.Join(referralsValues, ","),
			strings.Join(t0Conditions, ","),
			strings.Join(tMinus1Conditions, ","))
		_, err := conn.Exec(ctx, sql, referralsArgs...)

		return errors.Wrapf(err, "failed to decrement t1 and t2 active_referrals for args:%#v", referralsArgs)
	}); err != nil {
		if storage.IsErr(err, storage.ErrDuplicate, "pk") {
			return nil
		}

		return errors.Wrapf(err, "failed to run transaction to decrement active_referrals for t0&t-1, for %#v", usersThatStoppedMining)
	}

	return nil
}

func (r *repository) sessionNumber(date *time.Time) uint64 {
	return uint64(date.Unix()) / uint64(r.cfg.MiningSessionDuration.Max/stdlibtime.Second)
}
