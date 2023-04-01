// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	storage "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/terror"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) StartNewMiningSession( //nolint:funlen,gocognit // A lot of handling.
	ctx context.Context, ms *MiningSummary, rollbackNegativeMiningProgress *bool,
) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
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
		return errors.Wrapf(err,
			"failed to insertNewMiningSession:%#v,userID:%v,rollbackNegativeMiningProgress:%v", newMS, userID, shouldRollback)
	}
	if err = retry(ctx, func() error {
		summary, gErr := r.GetMiningSummary(ctx, userID)
		if gErr == nil {
			if summary.MiningSession == nil || summary.MiningSession.StartedAt == nil || !(summary.MiningSession.StartedAt.UnixMicro() == nowValue.UnixMicro()) {
				gErr = ErrNotFound
			} else {
				*ms = *summary
			}
		}

		return gErr
	}); err != nil {
		return errors.Wrapf(err, "permanently failed to GetMiningSummary for userID:%v", userID)
	}

	return errors.Wrapf(r.trySendMiningSessionMessage(ctx, userID, newMS),
		"failed to trySendMiningSessionMessage:%#v,userID:%v", ms, userID)
}

func (r *repository) getInternalMiningSummary(ctx context.Context, userID string) (*miningSummary, error) { //nolint:funlen // Big SQL.
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`SELECT u.last_natural_mining_started_at,
							   u.last_mining_started_at,
							   u.last_mining_ended_at,
							   u.previous_mining_started_at,
							   u.previous_mining_ended_at,
							   u.last_free_mining_session_awarded_at,
							   COALESCE(MAX(negative_balance.amount),'0') as negative_total_no_pre_staking_bonus_balance_amount,
							   COALESCE(MAX(negative_t0_balance.amount),'0') as negative_total_t0_no_pre_staking_bonus_balance_amount,
							   COALESCE(MAX(negative_t1_balance.amount),'0') as negative_total_t1_no_pre_staking_bonus_balance_amount,
							   COALESCE(MAX(negative_t2_balance.amount),'0') negative_total_t2_no_pre_staking_bonus_balance_amount,
							   0 AS mining_streak,
							   COALESCE(MAX(st.years),'0') AS pre_staking_years,
							   COALESCE(MAX(st.allocation),'0') AS pre_staking_allocation,
							   COALESCE(MAX(st_b.bonus),'0') as pre_staking_bonus
						FROM users u 
							LEFT JOIN pre_stakings st
								   ON st.user_id = u.user_id
						    	   AND st.worker_index = $1
							LEFT JOIN pre_staking_bonuses st_b
								   ON st.years = st_b.years
							LEFT JOIN balances_worker negative_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_balance.user_id = u.user_id
								  AND negative_balance.worker_index = $1
								  AND negative_balance.negative = TRUE
								  AND negative_balance.type = %[1]v
								  AND negative_balance.type_detail = ''
							LEFT JOIN balances_worker negative_t0_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t0_balance.user_id = u.user_id
								  AND negative_t0_balance.worker_index = $1
								  AND negative_t0_balance.negative = TRUE
								  AND negative_t0_balance.type = %[1]v
								  AND negative_t0_balance.type_detail = '%[2]v_' || u.referred_by
							LEFT JOIN balances_worker negative_t1_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t1_balance.user_id = u.user_id
								  AND negative_t1_balance.worker_index = $1
								  AND negative_t1_balance.negative = TRUE
								  AND negative_t1_balance.type = %[1]v
								  AND negative_t1_balance.type_detail = '%[3]v'
							LEFT JOIN balances_worker negative_t2_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t2_balance.user_id = u.user_id
								  AND negative_t2_balance.worker_index = $1
								  AND negative_t2_balance.negative = TRUE
								  AND negative_t2_balance.type = %[1]v
								  AND negative_t2_balance.type_detail = '%[4]v'
						WHERE u.user_id = $2
						GROUP BY u.user_id`, totalNoPreStakingBonusBalanceType, t0BalanceTypeDetail, t1BalanceTypeDetail, t2BalanceTypeDetail)
	resp, err := storage.Get[miningSummary](ctx, r.dbV2, sql, r.workerIndex(ctx), userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get the current mining summary for userID:%v", userID)
	}

	return resp, nil
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
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	var rollbackUsedAt, rollbackUsedAtCondition string
	if rollbackNegativeMiningSession != nil && *rollbackNegativeMiningSession {
		rollbackUsedAt = fmt.Sprintf("rollback_used_at = '%v',", ms.LastNaturalMiningStartedAt.Format(pgTimeFormat))
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

	return storage.DoInTransaction(ctx, r.dbV2, func(conn storage.QueryExecer) error {
		userHashCode, _ := ctx.Value(userHashCodeCtxValueKey).(uint64)
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
								 AND COALESCE(last_mining_ended_at,'1999-01-08 04:05:06')::timestamp = COALESCE($6,'1999-01-08 04:05:06')::timestamp
								 %[2]v;`, rollbackUsedAt, rollbackUsedAtCondition)
		rowsAffected, err := conn.Exec(ctx, sql,
			ms.LastNaturalMiningStartedAt.Time,
			ms.LastMiningStartedAt.Time,
			ms.LastMiningEndedAt.Time,
			lastFreeMiningSessionAwardedAtVal,
			userID,
			previousMiningEndedAtVal.Format(pgTimeFormat))
		if err != nil {
			return errors.Wrapf(err, "failed to update users mining time")
		}
		if rowsAffected.RowsAffected() != 1 {
			return ErrRaceCondition
		}
		rowsAffected, err = conn.Exec(ctx, fmt.Sprintf(`
						UPDATE balance_recalculation_worker_%[1]v 
							   SET enabled = TRUE,
								   last_mining_started_at = $1,
								   last_mining_ended_at = $2
							   WHERE user_id = $3;`, r.workerIndex(ctx)),
			ms.LastMiningStartedAt.Time,
			ms.LastMiningEndedAt.Time,
			userID,
		)
		if rowsAffected.RowsAffected() != 1 {
			if _, err := conn.Exec(ctx, fmt.Sprintf(`
						INSERT INTO balance_recalculation_worker(user_id,enabled,last_mining_started_at,last_mining_ended_at, hash_code, worker_index) 
																		   VALUES($3,TRUE,$1,$2, $4, %[1]v);`,
				r.workerIndex(ctx)),
				ms.LastMiningStartedAt.Time,
				ms.LastMiningEndedAt.Time,
				userID,
				userHashCode,
			); err != nil {
				return ErrRaceCondition
			}
		}
		if _, err := conn.Exec(ctx, fmt.Sprintf(
			`INSERT INTO blockchain_balance_synchronization_worker(user_id, hash_code, worker_index) VALUES ($1,$2,%[1]v) ON CONFLICT (worker_index, user_id) DO NOTHING;`, r.workerIndex(ctx)),
			userID, userHashCode,
		); err != nil {
			return errors.Wrapf(err, "failed to insert blockchain_balance_synchronization_worker_%v for userId:%v", r.workerIndex(ctx), userID)
		}
		if _, err := conn.Exec(ctx, fmt.Sprintf(
			`INSERT INTO extra_bonus_processing_worker(user_id, hash_code, worker_index) VALUES ($1,$2, %[1]v) ON CONFLICT (worker_index, user_id) DO NOTHING;`,
			r.workerIndex(ctx)),
			userID, userHashCode,
		); err != nil {
			return errors.Wrapf(err, "failed to insert extra_bonus_processing_worker_%v for userId:%v", r.workerIndex(ctx), userID)
		}
		if _, err := conn.Exec(ctx, fmt.Sprintf(
			`INSERT INTO mining_rates_recalculation_worker(user_id,  hash_code, worker_index) VALUES ($1, $2, %[1]v) ON CONFLICT (worker_index, user_id) DO NOTHING;`, r.workerIndex(ctx)),
			userID, userHashCode,
		); err != nil {
			return errors.Wrapf(err, "failed to insert mining_rates_recalculation_worker_%v for userId:%v", r.workerIndex(ctx), userID)
		}

		return nil
	})
}

func (r *repository) trySendMiningSessionMessage(ctx context.Context, userID string, newMS *miningSummary) error { //nolint:funlen // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sess := &MiningSession{
		LastNaturalMiningStartedAt: newMS.LastNaturalMiningStartedAt,
		StartedAt:                  newMS.LastMiningStartedAt,
		EndedAt:                    newMS.LastMiningEndedAt,
		UserID:                     &userID,
		MiningStreak:               newMS.MiningStreak,
	}
	if err := r.sendMiningSessionMessage(ctx, sess); err != nil {
		valueBytes, mErr := json.MarshalContext(ctx, sess)
		if mErr != nil {
			return multierror.Append( //nolint:wrapcheck // Not needed.
				errors.Wrapf(err, "failed to send a new mining session message: %#v", sess),
				errors.Wrapf(mErr, "failed to marshal %#v", sess),
			).ErrorOrNil()
		}
		type (
			MiningSessionDLQ struct {
				_msgpack            struct{} `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
				ID, UserID, Message string
			}
		)
		dlq := &MiningSessionDLQ{ID: uuid.NewString(), UserID: userID, Message: string(valueBytes)}
		_, insertErr := storage.Exec(ctx, r.dbV2,
			`INSERT INTO mining_sessions_dlq (id, user_id, message,hash_code) VALUES ($1,$2,$3, $4)`, dlq.ID, dlq.UserID, dlq.Message, 0)
		return multierror.Append( //nolint:wrapcheck // Not needed.
			errors.Wrapf(err, "failed to send a new mining session message: %#v", sess),
			errors.Wrapf(insertErr, "failed to dlqMiningSessionMessage:%#v because sendMiningSessionMessage failed", sess),
		).ErrorOrNil()
	}

	return nil
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
			  WHERE u.user_id = $1;`
	type referrals struct {
		T0UserID, TMinus1UserID     string
		T0HashCode, TMinus1HashCode int64
	}
	//nolint:revive // Nope.
	t0tMinus1, err := storage.Get[referrals](ctx, s.dbV2, sql, *ms.UserID)
	if err != nil || (t0tMinus1.T0UserID == "" && t0tMinus1.TMinus1UserID == "") {
		return errors.Wrapf(err, "failed to select for t0/t-1 information for userID:%v", *ms.UserID)
	}

	return errors.Wrapf(storage.DoInTransaction(ctx, s.dbV2, func(conn storage.QueryExecer) error {
		if _, err := conn.Exec(ctx, `INSERT INTO processed_mining_sessions(session_number, user_id) VALUES ($1, $2);`,
			s.sessionNumber(ms.LastNaturalMiningStartedAt),
			*ms.UserID,
		); err != nil {
			return errors.Wrapf(err, "failed to insert processed_mining_sessions for userId:%v, sessionNumber:%v",
				*ms.UserID, s.sessionNumber(ms.LastNaturalMiningStartedAt))
		}
		if rowsAffected, err := conn.Exec(ctx, `UPDATE active_referrals 
			   SET t1 = t1 + 1 
			   WHERE user_id = $1 and worker_index = $2;`,
			t0tMinus1.T0UserID,
			t0tMinus1.T0HashCode%int64(s.cfg.WorkerCount),
		); err != nil {
			return errors.Wrapf(err, "failed to update active_referrals for userId:%v", t0tMinus1.T0UserID)
		} else if rowsAffected.RowsAffected() != 1 && t0tMinus1.T0UserID != "" {
			if _, err := conn.Exec(ctx,
				`INSERT INTO active_referrals(t1, user_id, hash_code, worker_index) VALUES (1, $1, $2, $3);`,
				t0tMinus1.T0UserID, t0tMinus1.T0HashCode, t0tMinus1.T0HashCode%int64(s.cfg.WorkerCount),
			); err != nil {
				return errors.Wrapf(err, "failed to insert active_referrals for userId:%v", t0tMinus1.T0UserID)
			}
		}
		if rowsAffected, err := conn.Exec(ctx, `UPDATE active_referrals 
			   SET t2 = t2 + 1
			   WHERE user_id = $1 AND worker_index = $2;`,
			t0tMinus1.TMinus1UserID,
			t0tMinus1.TMinus1HashCode%int64(s.cfg.WorkerCount),
		); err != nil {
			return errors.Wrapf(err, "failed to update active_referrals_%v for userId:%v",
				t0tMinus1.TMinus1HashCode%int64(s.cfg.WorkerCount), t0tMinus1.TMinus1UserID)
		} else if rowsAffected.RowsAffected() != 1 && t0tMinus1.TMinus1UserID != "" {
			if _, err := conn.Exec(ctx,
				`INSERT INTO active_referrals(t1, user_id, hash_code, worker_index) VALUES (1, $1, $2, $3);`,
				t0tMinus1.TMinus1UserID,
				t0tMinus1.TMinus1HashCode,
				t0tMinus1.TMinus1HashCode%int64(s.cfg.WorkerCount),
			); err != nil {
				return errors.Wrapf(err, "failed to insert active_referrals for userId:%v", t0tMinus1.TMinus1UserID)
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

//nolint:funlen,gocognit,revive // .
func (r *repository) decrementActiveReferralCountForT0AndTMinus1(ctx context.Context, usersThatStoppedMining ...*userThatStoppedMining) error {
	if ctx.Err() != nil || len(usersThatStoppedMining) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	t0UserIDsPerWorkerIndex := make(map[uint64][]string)
	tMinus1UserIDsPerWorkerIndex := make(map[uint64][]string)
	processedMiningSessionValues := make([]string, 0, len(usersThatStoppedMining))
	for _, usr := range usersThatStoppedMining {
		t0WorkerIndex := uint64(usr.T0HashCode) % uint64(r.cfg.WorkerCount)
		tMinus1WorkerIndex := uint64(usr.TMinus1HashCode) % uint64(r.cfg.WorkerCount)
		if _, found := t0UserIDsPerWorkerIndex[t0WorkerIndex]; !found {
			t0UserIDsPerWorkerIndex[t0WorkerIndex] = make([]string, 0, len(usersThatStoppedMining))
		}
		if _, found := tMinus1UserIDsPerWorkerIndex[tMinus1WorkerIndex]; !found {
			tMinus1UserIDsPerWorkerIndex[tMinus1WorkerIndex] = make([]string, 0, len(usersThatStoppedMining))
		}
		t0UserIDsPerWorkerIndex[t0WorkerIndex] = append(t0UserIDsPerWorkerIndex[t0WorkerIndex], usr.T0UserID)
		tMinus1UserIDsPerWorkerIndex[tMinus1WorkerIndex] = append(tMinus1UserIDsPerWorkerIndex[tMinus1WorkerIndex], usr.TMinus1UserID)
		processedMiningSessionValues = append(processedMiningSessionValues, fmt.Sprintf(`(%v,true,'%v')`, r.sessionNumber(usr.LastMiningEndedAt), usr.UserID))
	}
	t0Values := make([]string, 0, len(t0UserIDsPerWorkerIndex))
	tMinus1Values := make([]string, 0, len(tMinus1UserIDsPerWorkerIndex))
	for workerIndex, t0UserIDs := range t0UserIDsPerWorkerIndex {
		for i := range t0UserIDs {
			t0UserIDs[i] = fmt.Sprintf(`('%v', %v)`, t0UserIDs[i], workerIndex)
		}
		t0Values = append(t0Values, fmt.Sprintf(`
			UPDATE active_referrals
			SET t1 = GREATEST(t1 - 1, 0)
			WHERE (user_id, worker_index) in (%[1]v);
		`, strings.Join(t0UserIDs, ",")))
	}
	for workerIndex, tMinus1UserIDs := range tMinus1UserIDsPerWorkerIndex {
		for i := range tMinus1UserIDs {
			tMinus1UserIDs[i] = fmt.Sprintf(`('%v', %v)`, tMinus1UserIDs[i], workerIndex)
		}
		tMinus1Values = append(tMinus1Values, fmt.Sprintf(`
		   UPDATE active_referrals
		   SET t2 = GREATEST(t2 - 1, 0)
		   WHERE (user_id, worker_index) in (%[1]v); 
		`, strings.Join(tMinus1UserIDs, ",")))
	}

	err := storage.DoInTransaction(ctx, r.dbV2, func(conn storage.QueryExecer) error {
		if _, err := conn.Exec(ctx, fmt.Sprintf(`INSERT INTO processed_mining_sessions(session_number, negative, user_id) VALUES %[1]v;`,
			strings.Join(processedMiningSessionValues, ",")),
		); err != nil {
			return errors.Wrap(err, "failed to insert processed_mining_sessions")
		}
		for _, t0Query := range t0Values {
			if _, err := conn.Exec(ctx, t0Query); err != nil {
				return errors.Wrapf(err, "failed to update active_referrals: %v", t0Query)
			}
		}
		for _, tMinus1Query := range tMinus1Values {
			if _, err := conn.Exec(ctx, tMinus1Query); err != nil {
				return errors.Wrapf(err, "failed to update active_referrals: %v", tMinus1Query)
			}
		}

		return nil
	})
	if err != nil {
		if storage.IsErr(err, storage.ErrDuplicate, "pkey") {
			return nil
		}
		return errors.Wrapf(err, "failed to eval script to decrement active_referrals for t0&t-1, for %#v", usersThatStoppedMining)
	}

	return nil
}

func (r *repository) sessionNumber(date *time.Time) uint64 {
	return uint64(date.Unix()) / uint64(r.cfg.MiningSessionDuration.Max/stdlibtime.Second)
}
