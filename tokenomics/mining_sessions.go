// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/terror"
	"github.com/ice-blockchain/wintr/time"
)

type (
	startOrExtendMiningSession struct {
		ResurrectSoloUsedAt                  *time.Time `redis:"resurrect_solo_used_at"`
		MiningSessionSoloLastStartedAt       *time.Time `redis:"mining_session_solo_last_started_at"`
		MiningSessionSoloStartedAt           *time.Time `redis:"mining_session_solo_started_at"`
		MiningSessionSoloEndedAt             *time.Time `redis:"mining_session_solo_ended_at"`
		MiningSessionSoloDayOffLastAwardedAt *time.Time `redis:"mining_session_solo_day_off_last_awarded_at"`
		MiningSessionSoloPreviouslyEndedAt   *time.Time `redis:"mining_session_solo_previously_ended_at"`
		deserializedUsersKey
	}
	getCurrentMiningSession struct {
		ResurrectSoloUsedAt                  *time.Time `redis:"resurrect_solo_used_at"`
		MiningSessionSoloLastStartedAt       *time.Time `redis:"mining_session_solo_last_started_at"`
		MiningSessionSoloStartedAt           *time.Time `redis:"mining_session_solo_started_at"`
		MiningSessionSoloEndedAt             *time.Time `redis:"mining_session_solo_ended_at"`
		MiningSessionSoloDayOffLastAwardedAt *time.Time `redis:"mining_session_solo_day_off_last_awarded_at"`
		MiningSessionSoloPreviouslyEndedAt   *time.Time `redis:"mining_session_solo_previously_ended_at"`
		PreStakingAllocation                 uint16     `redis:"pre_staking_allocation"`
		PreStakingBonus                      uint16     `redis:"pre_staking_bonus"`
		BalanceTotal                         float64    `redis:"balance_total"`
		SlashingRateSolo                     float64    `redis:"slashing_rate_solo"`
		SlashingRateT0                       float64    `redis:"slashing_rate_t0"`
		SlashingRateT1                       float64    `redis:"slashing_rate_t1"`
		SlashingRateT2                       float64    `redis:"slashing_rate_t2"`
		IDT0                                 int64      `redis:"id_t0"`
		IDTMinus1                            int64      `redis:"id_tminus1"`
	}
)

func (r *repository) StartNewMiningSession( //nolint:funlen,gocognit // A lot of handling.
	ctx context.Context, ms *MiningSummary, rollbackNegativeMiningProgress *bool,
) error {
	userID := *ms.MiningSession.UserID
	id, err := r.getOrInitInternalID(ctx, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	now := time.Now()
	old, err := storage.Get[getCurrentMiningSession](ctx, r.db, serializedUsersKey(id))
	if err != nil || len(old) == 0 {
		if err == nil {
			err = errors.Wrapf(ErrRelationNotFound, "missing state for id:%v", id)
		}

		return errors.Wrapf(err, "failed to get miningSummary for id:%v", id)
	}
	if !old[0].MiningSessionSoloEndedAt.IsNil() &&
		!old[0].MiningSessionSoloLastStartedAt.IsNil() &&
		old[0].MiningSessionSoloEndedAt.After(*now.Time) &&
		(now.Sub(*old[0].MiningSessionSoloLastStartedAt.Time)/r.cfg.MiningSessionDuration.Min)%2 == 0 {
		return ErrDuplicate
	}
	shouldRollback, err := r.validateRollbackNegativeMiningProgress(old[0].PreStakingAllocation, old[0].PreStakingBonus, old[0].SlashingRateSolo, old[0].SlashingRateT0, old[0].SlashingRateT1, old[0].SlashingRateT2, old[0].MiningSessionSoloEndedAt, old[0].ResurrectSoloUsedAt, now, rollbackNegativeMiningProgress) //nolint:lll // .
	if err != nil {
		return err
	}
	if err = r.updateTMinus1(ctx, id, old[0].IDT0, old[0].IDTMinus1); err != nil {
		return errors.Wrapf(err, "failed to updateTMinus1 for id:%v", id)
	}
	oldMS := &startOrExtendMiningSession{
		ResurrectSoloUsedAt:                  old[0].ResurrectSoloUsedAt,
		MiningSessionSoloLastStartedAt:       old[0].MiningSessionSoloLastStartedAt,
		MiningSessionSoloStartedAt:           old[0].MiningSessionSoloStartedAt,
		MiningSessionSoloEndedAt:             old[0].MiningSessionSoloEndedAt,
		MiningSessionSoloDayOffLastAwardedAt: old[0].MiningSessionSoloDayOffLastAwardedAt,
		MiningSessionSoloPreviouslyEndedAt:   old[0].MiningSessionSoloPreviouslyEndedAt,
	}
	newMS, extension := r.newStartOrExtendMiningSession(oldMS, now)
	newMS.ID = id
	if shouldRollback != nil && *shouldRollback && oldMS.ResurrectSoloUsedAt.IsNil() {
		newMS.ResurrectSoloUsedAt = time.New(stdlibtime.Date(3000, 0, 0, 0, 0, 0, 0, nil)) //nolint:gomnd // .
	}
	sess := &MiningSession{
		LastNaturalMiningStartedAt: newMS.MiningSessionSoloLastStartedAt,
		StartedAt:                  newMS.MiningSessionSoloStartedAt,
		EndedAt:                    newMS.MiningSessionSoloEndedAt,
		PreviouslyEndedAt:          newMS.MiningSessionSoloPreviouslyEndedAt,
		Extension:                  extension,
		MiningStreak:               r.calculateMiningStreak(now, newMS.MiningSessionSoloStartedAt, newMS.MiningSessionSoloEndedAt),
		UserID:                     &userID,
	}
	if err = r.sendMiningSessionMessage(ctx, sess); err != nil {
		return errors.Wrapf(err, "failed to sendMiningSessionMessage:%#v", sess)
	}
	if err = storage.Set(ctx, r.db, newMS); err != nil {
		return errors.Wrapf(err, "failed to insertNewMiningSession:%#v", newMS)
	}

	return errors.Wrapf(retry(ctx, func() error {
		summary, gErr := r.GetMiningSummary(ctx, userID)
		if gErr == nil {
			if summary.MiningSession == nil || summary.MiningSession.StartedAt.IsNil() || !summary.MiningSession.StartedAt.Equal(*now.Time) {
				gErr = ErrNotFound
			} else {
				*ms = *summary
			}
		}

		return gErr
	}), "permanently failed to GetMiningSummary for userID:%v", userID)
}

//nolint:funlen // .
func (r *repository) updateTMinus1(ctx context.Context, id, idT0, idTMinus1 int64) error {
	if idTMinus1 < 1 || idT0 < 1 {
		return nil
	}
	if oldTminus1Data, err := storage.Get[struct {
		UserID string `redis:"user_id"`
	}](ctx, r.db, serializedUsersKey(idTMinus1)); err != nil || len(oldTminus1Data) != 0 {
		return errors.Wrapf(err, "failed to get state for t-1:%v", idTMinus1)
	}
	idTMinus1 = 0
	if t0Data, err := storage.Get[struct {
		IDT0 int64 `redis:"id_t0"`
	}](ctx, r.db, serializedUsersKey(idT0)); err != nil {
		return errors.Wrapf(err, "failed to get state for t0:%v", idT0)
	} else if len(t0Data) != 0 {
		idTMinus1 = t0Data[0].IDT0
	}
	type (
		replaceIDTMinus1 struct {
			deserializedUsersKey
			IDTMinus1              int64   `redis:"id_tminus1"`
			BalanceForTMinus1      float64 `redis:"balance_for_tminus1"`
			SlashingRateForTMinus1 float64 `redis:"slashing_rate_for_tminus1"`
		}
	)
	if err := storage.Set(ctx, r.db, &replaceIDTMinus1{deserializedUsersKey: deserializedUsersKey{ID: id}, IDTMinus1: idTMinus1}); err != nil {
		return errors.Wrapf(err, "failed to replaceIDTMinus1, id:%v, newIDTMinus1:%v", id, idTMinus1)
	}
	stdlibtime.Sleep(stdlibtime.Second)
	afterReplaceIDTMinus1, err := storage.Get[struct {
		BalanceForTMinus1      float64 `redis:"balance_for_tminus1"`
		SlashingRateForTMinus1 float64 `redis:"slashing_rate_for_tminus1"`
	}](ctx, r.db, serializedUsersKey(id))
	if err != nil || len(afterReplaceIDTMinus1) == 0 || (afterReplaceIDTMinus1[0].BalanceForTMinus1 == 0.0 && afterReplaceIDTMinus1[0].SlashingRateForTMinus1 == 0.0) { //nolint:lll // .
		if err == nil && len(afterReplaceIDTMinus1) == 0 {
			err = errors.Wrapf(ErrRelationNotFound, "missing state[2] for id:%v", id)
		}

		return errors.Wrapf(err, "failed to get state for id:%v, after t-1 id was updated", id)
	}

	return errors.Wrapf(storage.Set(ctx, r.db, &replaceIDTMinus1{deserializedUsersKey: deserializedUsersKey{ID: id}, IDTMinus1: idTMinus1}),
		"failed[2] to replaceIDTMinus1, id:%v, newIDTMinus1:%v", id, idTMinus1)
}

func (r *repository) validateRollbackNegativeMiningProgress(
	preStakingAllocation, preStakingBonus uint16,
	slashingRateSolo, slashingRateT0, slashingRateT1, slashingRateT2 float64,
	miningSessionSoloEndedAt, resurrectSoloUsedAt, now *time.Time,
	rollbackNegativeMiningProgress *bool,
) (*bool, error) {
	if !resurrectSoloUsedAt.IsNil() || miningSessionSoloEndedAt.IsNil() ||
		(now.Sub(*miningSessionSoloEndedAt.Time) < r.cfg.RollbackNegativeMining.Available.After ||
			now.Sub(*miningSessionSoloEndedAt.Time) > r.cfg.RollbackNegativeMining.Available.Until) {
		return nil, nil //nolint:nilnil // Nope.
	}
	amountLost := (slashingRateSolo + slashingRateT0 + slashingRateT1 + slashingRateT2) * now.Sub(*miningSessionSoloEndedAt.Time).Seconds()
	amountLost = ((amountLost * float64(100-preStakingAllocation)) / 100) + ((amountLost * float64(preStakingAllocation*(preStakingBonus+100))) / (100 * 100))
	if amountLost == 0.0 {
		return nil, nil //nolint:nilnil // Nope.
	}
	if rollbackNegativeMiningProgress == nil {
		return nil, terror.New(ErrNegativeMiningProgressDecisionRequired, map[string]any{
			"amount":                fmt.Sprint(amountLost),
			"duringTheLastXSeconds": uint64(now.Sub(*miningSessionSoloEndedAt.Time).Seconds()),
		})
	}

	return rollbackNegativeMiningProgress, nil
}

func (r *repository) newStartOrExtendMiningSession(old *startOrExtendMiningSession, now *time.Time) (*startOrExtendMiningSession, stdlibtime.Duration) {
	resp := &startOrExtendMiningSession{
		ResurrectSoloUsedAt:                old.ResurrectSoloUsedAt,
		MiningSessionSoloStartedAt:         now,
		MiningSessionSoloLastStartedAt:     now,
		MiningSessionSoloEndedAt:           time.New(now.Add(r.cfg.MiningSessionDuration.Max)),
		MiningSessionSoloPreviouslyEndedAt: old.MiningSessionSoloEndedAt,
	}
	if old.MiningSessionSoloEndedAt.IsNil() || old.MiningSessionSoloStartedAt.IsNil() || old.MiningSessionSoloEndedAt.Before(*now.Time) {
		return resp, r.cfg.MiningSessionDuration.Max
	}
	resp.MiningSessionSoloPreviouslyEndedAt = old.MiningSessionSoloPreviouslyEndedAt
	resp.MiningSessionSoloStartedAt = old.MiningSessionSoloStartedAt
	resp.MiningSessionSoloDayOffLastAwardedAt = old.MiningSessionSoloDayOffLastAwardedAt
	var durationSinceLastFreeMiningSessionAwarded stdlibtime.Duration
	if resp.MiningSessionSoloDayOffLastAwardedAt.IsNil() {
		durationSinceLastFreeMiningSessionAwarded = now.Sub(*resp.MiningSessionSoloStartedAt.Time)
	} else {
		durationSinceLastFreeMiningSessionAwarded = now.Sub(*resp.MiningSessionSoloDayOffLastAwardedAt.Time)
	}
	freeMiningSession := uint64(0)
	minimumDurationForAwardingFreeMiningSession := stdlibtime.Duration(r.cfg.ConsecutiveNaturalMiningSessionsRequiredFor1ExtraFreeArtificialMiningSession.Max) * r.cfg.MiningSessionDuration.Max //nolint:lll // .
	if durationSinceLastFreeMiningSessionAwarded >= minimumDurationForAwardingFreeMiningSession {
		resp.MiningSessionSoloDayOffLastAwardedAt = now
		freeMiningSession++
	}
	if freeSessions := stdlibtime.Duration(r.calculateRemainingFreeMiningSessions(now, old.MiningSessionSoloEndedAt) + freeMiningSession); freeSessions > 0 {
		resp.MiningSessionSoloEndedAt = time.New(resp.MiningSessionSoloEndedAt.Add(freeSessions * r.cfg.MiningSessionDuration.Max))
	}

	return resp, resp.MiningSessionSoloEndedAt.Sub(*old.MiningSessionSoloEndedAt.Time)
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
	if ctx.Err() != nil || len(msg.Value) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	ms := new(MiningSession)
	if err := json.UnmarshalContext(ctx, msg.Value, ms); err != nil || ms.UserID == nil {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), ms)
	}

	return multierror.Append( //nolint:wrapcheck // Not needed.
		errors.Wrapf(s.incrementTotalActiveUsers(ctx, ms), "failed to incrementTotalActiveUsers for %#v", ms),
		errors.Wrapf(s.incrementActiveReferralCountForT0AndTMinus1(ctx, ms), "failed to incrementActiveReferralCountForT0AndTMinus1 for %#v", ms),
		errors.Wrapf(s.trySwitchToNextAdoption(ctx), "failed to trySwitchToNextAdoption"),
	).ErrorOrNil()
}

//nolint:funlen,revive,gocognit // .
func (s *miningSessionsTableSource) incrementActiveReferralCountForT0AndTMinus1(ctx context.Context, ms *MiningSession) (err error) {
	if ctx.Err() != nil || !ms.LastNaturalMiningStartedAt.Equal(*ms.StartedAt.Time) {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	duplGuardKey := ms.duplGuardKey(s.repository, "incr_active_ref")
	if set, dErr := s.db.SetNX(ctx, duplGuardKey, "", s.cfg.MiningSessionDuration.Min).Result(); dErr != nil || !set {
		if dErr == nil {
			dErr = ErrDuplicate
		}

		return errors.Wrapf(dErr, "SetNX failed for mining_session_dupl_guard, miningSession: %#v", ms)
	}
	defer func() {
		if err != nil {
			undoCtx, cancelUndo := context.WithTimeout(context.Background(), requestDeadline)
			defer cancelUndo()
			err = multierror.Append( //nolint:wrapcheck // .
				err,
				errors.Wrapf(s.db.Del(undoCtx, duplGuardKey).Err(), "failed to del mining_session_dupl_guard key"),
			).ErrorOrNil()
		}
	}()
	id, err := s.getOrInitInternalID(ctx, *ms.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", *ms.UserID)
	}
	referees, err := storage.Get[struct {
		deserializedUsersKey
		IDT0      int64 `redis:"id_t0"`
		IDTMinus1 int64 `redis:"id_tminus1"`
	}](ctx, s.db, serializedUsersKey(id))
	if err != nil || len(referees) == 0 || (referees[0].IDT0 < 1 && referees[0].IDTMinus1 < 1) {
		return errors.Wrapf(err, "failed to get referees for id:%v, userID:%v", id, *ms.UserID)
	}
	if referees[0].IDT0 < 1 || referees[0].IDTMinus1 < 1 {
		if referees[0].IDT0 >= 1 {
			err = s.db.HIncrBy(ctx, serializedUsersKey(referees[0].IDT0), "active_t1_referrals", 1).Err()
		}
		if referees[0].IDTMinus1 >= 1 {
			err = s.db.HIncrBy(ctx, serializedUsersKey(referees[0].IDTMinus1), "active_t2_referrals", 1).Err()
		}
	} else {
		responses, txErr := s.db.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			return multierror.Append( //nolint:wrapcheck // .
				pipeliner.HIncrBy(ctx, serializedUsersKey(referees[0].IDT0), "active_t1_referrals", 1).Err(),
				pipeliner.HIncrBy(ctx, serializedUsersKey(referees[0].IDTMinus1), "active_t2_referrals", 1).Err(),
			).ErrorOrNil()
		})
		if txErr == nil {
			errs := make([]error, 0, len(responses))
			for _, response := range responses {
				errs = append(errs, errors.Wrapf(response.Err(), "failed to `%v`", response.FullName()))
			}
			txErr = multierror.Append(nil, errs...).ErrorOrNil()
		}
		err = txErr
	}

	return errors.Wrapf(err, "failed to increment active referrals for t0&t-1, id:%v, userID:%v, ref:%#v", id, *ms.UserID, referees[0])
}

func (r *repository) sessionNumber(date *time.Time) uint64 {
	return uint64(date.Unix()) / uint64(r.cfg.MiningSessionDuration.Min/stdlibtime.Second)
}

func (ms *MiningSession) duplGuardKey(repo *repository, guardType string) string {
	return fmt.Sprintf("mining_session_dupl_guards:%v~%v~%v", guardType, ms.UserID, repo.sessionNumber(ms.StartedAt))
}
