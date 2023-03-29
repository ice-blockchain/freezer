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
			if summary.MiningSession == nil || summary.MiningSession.StartedAt == nil || !summary.MiningSession.StartedAt.Equal(*now.Time) {
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
							   negative_balance.amount,
							   negative_t0_balance.amount,
							   negative_t1_balance.amount,
							   negative_t2_balance.amount,
							   0 AS mining_streak,
							   MAX(st.years) AS years,
							   MAX(st.allocation) AS allocation,
							   st_b.bonus
						FROM users u 
							LEFT JOIN pre_stakings_%[1]v st
								   ON st.user_id = u.user_id
							LEFT JOIN pre_staking_bonuses st_b
								   ON st.years = st_b.years
							LEFT JOIN balances_%[1]v negative_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_balance.user_id = u.user_id
								  AND negative_balance.negative = TRUE
								  AND negative_balance.type = %[2]v
								  AND negative_balance.type_detail = ''
							LEFT JOIN balances_%[1]v negative_t0_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t0_balance.user_id = u.user_id
								  AND negative_t0_balance.negative = TRUE
								  AND negative_t0_balance.type = %[2]v
								  AND negative_t0_balance.type_detail = '%[3]v_' || u.referred_by
							LEFT JOIN balances_%[1]v negative_t1_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t1_balance.user_id = u.user_id
								  AND negative_t1_balance.negative = TRUE
								  AND negative_t1_balance.type = %[2]v
								  AND negative_t1_balance.type_detail = '%[4]v'
							LEFT JOIN balances_%[1]v negative_t2_balance
								   ON u.rollback_used_at IS NULL
								  AND negative_t2_balance.user_id = u.user_id
								  AND negative_t2_balance.negative = TRUE
								  AND negative_t2_balance.type = %[2]v
								  AND negative_t2_balance.type_detail = '%[5]v'
						WHERE u.user_id = :user_id
						GROUP BY u.user_id`, r.workerIndex(ctx), totalNoPreStakingBonusBalanceType, t0BalanceTypeDetail, t1BalanceTypeDetail, t2BalanceTypeDetail)
	params := make(map[string]any, 1)
	params["user_id"] = userID
	resp := make([]*miningSummary, 0, 1)
	if err := r.db.PrepareExecuteTyped(sql, params, &resp); err != nil {
		return nil, errors.Wrapf(err, "failed to get the current mining summary for userID:%v", userID)
	}
	if len(resp) == 0 {
		return nil, ErrRelationNotFound
	}

	return resp[0], nil
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
		rollbackUsedAt = fmt.Sprintf("rollback_used_at = %v,", ms.LastNaturalMiningStartedAt.UnixNano())
		rollbackUsedAtCondition = "AND rollback_used_at IS NULL"
	}
	const null = "null"
	previousMiningEndedAtVal := null
	if old.LastMiningEndedAt != nil {
		previousMiningEndedAtVal = fmt.Sprint(old.LastMiningEndedAt.UnixNano())
	}
	lastFreeMiningSessionAwardedAtVal := null
	if ms.LastFreeMiningSessionAwardedAt != nil {
		lastFreeMiningSessionAwardedAtVal = fmt.Sprint(ms.LastFreeMiningSessionAwardedAt.UnixNano())
	}
	//nolint:dupword // Nope.
	script := fmt.Sprintf(`resp, err = box.execute([[START TRANSACTION;]]) 
if err ~= nil then
	return err
end 
resp, err = box.execute([[ UPDATE users
						   SET updated_at = %[1]v,
							   last_natural_mining_started_at = %[1]v,
							   last_mining_started_at = %[2]v,
							   last_mining_ended_at = %[3]v,
							   previous_mining_started_at = (CASE WHEN last_mining_started_at = %[2]v THEN previous_mining_started_at ELSE last_mining_started_at END),
							   previous_mining_ended_at = (CASE WHEN last_mining_started_at = %[2]v THEN previous_mining_ended_at ELSE last_mining_ended_at END),
 							   %[4]v							   
							   last_free_mining_session_awarded_at = %[5]v
						   WHERE user_id = '%[6]v'
						     AND IFNULL(last_mining_ended_at,0) = IFNULL(%[7]v,0)
							 %[8]v;]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]]) 
	return err
end 
if resp.row_count ~= 1 then
	box.execute([[ROLLBACK;]]) 
	return "race condition"
end 
resp, err = box.execute([[ UPDATE balance_recalculation_worker_%[9]v 
						   SET enabled = TRUE,
							   last_mining_started_at = %[2]v,
							   last_mining_ended_at = %[3]v
						   WHERE user_id = '%[6]v';]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]]) 
	return err
end
if resp.row_count ~= 1 then
	resp, err = box.execute([[ INSERT INTO balance_recalculation_worker_%[9]v(user_id,enabled,last_mining_started_at,last_mining_ended_at) 
																	   VALUES('%[6]v',TRUE,%[2]v,%[3]v);]])
	if err ~= nil then
		box.execute([[ROLLBACK;]]) 
		return "race condition 2"
	end
end
box.execute([[ INSERT INTO blockchain_balance_synchronization_worker_%[9]v(user_id) VALUES ('%[6]v');]])
box.execute([[ INSERT INTO extra_bonus_processing_worker_%[9]v(user_id) VALUES ('%[6]v');]])
box.execute([[ INSERT INTO mining_rates_recalculation_worker_%[9]v(user_id) VALUES ('%[6]v');]])
resp,err = box.execute([[COMMIT;]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]])
	return err
end 
return ''`,
		ms.LastNaturalMiningStartedAt.UnixNano(),
		ms.LastMiningStartedAt.UnixNano(),
		ms.LastMiningEndedAt.UnixNano(),
		rollbackUsedAt,
		lastFreeMiningSessionAwardedAtVal,
		userID,
		previousMiningEndedAtVal,
		rollbackUsedAtCondition,
		r.workerIndex(ctx))
	resp := make([]string, 0, 1)
	if err := r.db.EvalTyped(script, []any{}, &resp); err != nil {
		return errors.Wrapf(err, "failed to eval script to insert mining session for %#v", ms)
	} else if errMessage := resp[0]; errMessage != "" {
		if strings.Contains(errMessage, `race condition`) {
			return ErrRaceCondition
		}

		return errors.Errorf("insert mining session script returned unexpected error message:`%v`, for %#v", errMessage, ms)
	}

	return nil
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

		return multierror.Append( //nolint:wrapcheck // Not needed.
			errors.Wrapf(err, "failed to send a new mining session message: %#v", sess),
			errors.Wrapf(r.db.InsertTyped(fmt.Sprintf("MINING_SESSIONS_DLQ_%v", r.workerIndex(ctx)), dlq, &[]*MiningSessionDLQ{}),
				"failed to dlqMiningSessionMessage:%#v because sendMiningSessionMessage failed", sess),
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
	sql := `SELECT t0.user_id AS t0_user_id,
				   tMinus1.user_id AS tMinus1_user_id,
				   t0.hash_code AS t0HashCode,
				   tMinus1.hash_code AS tMinus1HashCode
			  FROM users u
			  		LEFT JOIN users t0
			  			   ON t0.user_id = u.referred_by
						  AND t0.user_id != u.user_id
			  		LEFT JOIN users tMinus1
			  			   ON tMinus1.user_id = t0.referred_by
						  AND tMinus1.user_id != t0.user_id
						  AND tMinus1.user_id != u.user_id
			  WHERE u.user_id = :user_id`
	params := make(map[string]any, 1)
	params["user_id"] = *ms.UserID
	rows := make([]*struct {
		_msgpack                    struct{} `msgpack:",asArray"` //nolint:tagliatelle,revive,nosnakecase // To insert we need asArray
		T0UserID, TMinus1UserID     string
		T0HashCode, TMinus1HashCode uint64
	}, 0, 1)
	//nolint:revive // Nope.
	if err := s.db.PrepareExecuteTyped(sql, params, &rows); err != nil || len(rows) == 0 || (rows[0].T0UserID == "" && rows[0].TMinus1UserID == "") {
		return errors.Wrapf(err, "failed to select for t0/t-1 information for userID:%v", *ms.UserID)
	}
	//nolint:dupword // Nope.
	script := fmt.Sprintf(`resp, err = box.execute([[START TRANSACTION;]]) 
if err ~= nil then
	return err
end 
resp, err = box.execute([[ INSERT INTO processed_mining_sessions(session_number, user_id) VALUES (%[1]v, '%[2]v'); ]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]]) 
	return err
end 
resp, err = box.execute([[ UPDATE active_referrals_%[3]v 
						   SET t1 = t1 + 1
						   WHERE user_id = '%[4]v';]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]]) 
	return err
end
if resp.row_count ~= 1 and '%[4]v' ~= '' then
	resp, err = box.execute([[ INSERT INTO active_referrals_%[3]v(t1, user_id) VALUES (1, '%[4]v'); ]]) 
	if err ~= nil then
		box.execute([[ROLLBACK;]]) 
		return err
	end 
end
resp, err = box.execute([[ UPDATE active_referrals_%[5]v 
						   SET t2 = t2 + 1
						   WHERE user_id = '%[6]v';]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]]) 
	return err
end
if resp.row_count ~= 1 and '%[6]v' ~= '' then
	resp, err = box.execute([[ INSERT INTO active_referrals_%[5]v(t2, user_id) VALUES (1, '%[6]v'); ]]) 
	if err ~= nil then
		box.execute([[ROLLBACK;]]) 
		return err
	end 
end
resp, err = box.execute([[COMMIT;]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]])
	return err
end 
return ''`,
		s.sessionNumber(ms.LastNaturalMiningStartedAt),
		*ms.UserID,
		rows[0].T0HashCode%s.cfg.WorkerCount,
		rows[0].T0UserID,
		rows[0].TMinus1HashCode%s.cfg.WorkerCount,
		rows[0].TMinus1UserID)
	resp := make([]string, 0, 1)
	if err := s.db.EvalTyped(script, []any{}, &resp); err != nil {
		return errors.Wrapf(err, "failed to eval script to increment active_referrals for t0&t-1, for %#v", ms)
	} else if errMessage := resp[0]; errMessage != "" {
		return errors.Errorf("increment active_referrals for t0&t-1 script returned unexpected error message:`%v`, for %#v", errMessage, ms)
	}

	return nil
}

type (
	userThatStoppedMining struct {
		LastMiningEndedAt               *time.Time
		UserID, T0UserID, TMinus1UserID string
		T0HashCode, TMinus1HashCode     uint64
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
		if _, found := t0UserIDsPerWorkerIndex[usr.T0HashCode%r.cfg.WorkerCount]; !found {
			t0UserIDsPerWorkerIndex[usr.T0HashCode%r.cfg.WorkerCount] = make([]string, 0, len(usersThatStoppedMining))
		}
		if _, found := tMinus1UserIDsPerWorkerIndex[usr.TMinus1HashCode%r.cfg.WorkerCount]; !found {
			tMinus1UserIDsPerWorkerIndex[usr.TMinus1HashCode%r.cfg.WorkerCount] = make([]string, 0, len(usersThatStoppedMining))
		}
		t0UserIDsPerWorkerIndex[usr.T0HashCode%r.cfg.WorkerCount] = append(t0UserIDsPerWorkerIndex[usr.T0HashCode%r.cfg.WorkerCount], usr.T0UserID)
		tMinus1UserIDsPerWorkerIndex[usr.TMinus1HashCode%r.cfg.WorkerCount] = append(tMinus1UserIDsPerWorkerIndex[usr.TMinus1HashCode%r.cfg.WorkerCount], usr.TMinus1UserID) //nolint:lll // .
		processedMiningSessionValues = append(processedMiningSessionValues, fmt.Sprintf(`(%v,true,'%v')`, r.sessionNumber(usr.LastMiningEndedAt), usr.UserID))
	}
	t0Values := make([]string, 0, len(t0UserIDsPerWorkerIndex))
	tMinus1Values := make([]string, 0, len(tMinus1UserIDsPerWorkerIndex))
	for t0WorkerIndex, t0UserIDs := range t0UserIDsPerWorkerIndex {
		for i := range t0UserIDs {
			t0UserIDs[i] = fmt.Sprintf(`'%v'`, t0UserIDs[i])
		}
		t0Values = append(t0Values, fmt.Sprintf(`
resp, err = box.execute([[ UPDATE active_referrals_%[1]v 
						   SET t1 = t1 - 1
						   WHERE user_id in (%[2]v);]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]]) 
	return err
end`, t0WorkerIndex, strings.Join(t0UserIDs, ",")))
	}
	for tMinus1WorkerIndex, tMinus1UserIDs := range tMinus1UserIDsPerWorkerIndex {
		for i := range tMinus1UserIDs {
			tMinus1UserIDs[i] = fmt.Sprintf(`'%v'`, tMinus1UserIDs[i])
		}
		tMinus1Values = append(tMinus1Values, fmt.Sprintf(`
resp, err = box.execute([[ UPDATE active_referrals_%[1]v 
						   SET t2 = t2 - 1
						   WHERE user_id in (%[2]v);]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]]) 
	return err
end`, tMinus1WorkerIndex, strings.Join(tMinus1UserIDs, ",")))
	}
	script := fmt.Sprintf(`resp, err = box.execute([[START TRANSACTION;]]) 
if err ~= nil then
	return err
end 
resp, err = box.execute([[ INSERT INTO processed_mining_sessions(session_number, negative, user_id) VALUES %[1]v; ]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]]) 
	return err
end 
%[2]v
%[3]v
resp, err = box.execute([[COMMIT;]]) 
if err ~= nil then
	box.execute([[ROLLBACK;]])
	return err
end 
return ''`,
		strings.Join(processedMiningSessionValues, ","),
		strings.Join(t0Values, "\n"),
		strings.Join(tMinus1Values, "\n"))
	resp := make([]string, 0, 1)
	if err := r.db.EvalTyped(script, []any{}, &resp); err != nil {
		return errors.Wrapf(err, "failed to eval script to decrement active_referrals for t0&t-1, for %#v", usersThatStoppedMining)
	} else if errMessage := resp[0]; errMessage != "" {
		if strings.Contains(errMessage, "Duplicate key exists in unique index \"pk_unnamed_PROCESSED_MINING_SESSIONS_1\"") {
			return nil
		}

		return errors.Errorf("decrement active_referrals for t0&t-1 script returned unexpected error message:`%v`,for %#v", errMessage, usersThatStoppedMining)
	}

	return nil
}

func (r *repository) sessionNumber(date *time.Time) uint64 {
	return uint64(date.Unix()) / uint64(r.cfg.MiningSessionDuration.Max/stdlibtime.Second)
}
