// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetAdoptionSummary(ctx context.Context) (as *AdoptionSummary, err error) {
	if as = new(AdoptionSummary); ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	if as.TotalActiveUsers, err = r.db.Get(ctx, r.totalActiveUsersKey(*time.Now().Time)).Uint64(); err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrap(err, "failed to get current totalActiveUsers")
	}
	if as.Milestones, err = getAllAdoptions[string](ctx, r.db); err != nil {
		return nil, errors.Wrap(err, "failed to get all adoption milestones")
	}

	return
}

func (r *repository) totalActiveUsersKey(date stdlibtime.Time) string {
	return fmt.Sprintf("%v:%v", totalActiveUsersGlobalKey, date.Format(r.cfg.totalActiveUsersAggregationIntervalDateFormat()))
}

func (r *repository) extractTimeFromTotalActiveUsersKey(key string) *time.Time {
	parseTime, err := stdlibtime.Parse(r.cfg.totalActiveUsersAggregationIntervalDateFormat(), strings.ReplaceAll(key, totalActiveUsersGlobalKey+":", ""))
	log.Panic(err)

	return time.New(parseTime)
}

func (r *repository) incrementTotalActiveUsers(ctx context.Context, ms *MiningSession) (err error) { //nolint:funlen // .
	duplGuardKey := ms.duplGuardKey(r, "incr_total_active_users")
	if set, dErr := r.db.SetNX(ctx, duplGuardKey, "", r.cfg.MiningSessionDuration.Min).Result(); dErr != nil || !set {
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
				errors.Wrapf(r.db.Del(undoCtx, duplGuardKey).Err(), "failed to del mining_session_dupl_guard key"),
			).ErrorOrNil()
		}
	}()
	keys := ms.detectIncrTotalActiveUsersKeys(r)
	responses, err := r.db.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, key := range keys {
			if err = pipeliner.Incr(ctx, key).Err(); err != nil {
				return err
			}
		}

		return nil
	})
	if err == nil {
		errs := make([]error, 0, len(responses))
		for _, response := range responses {
			errs = append(errs, errors.Wrapf(response.Err(), "failed to `%v`", response.FullName()))
		}
		err = multierror.Append(nil, errs...).ErrorOrNil()
	}

	return errors.Wrapf(err, "failed to incr total active users for keys:%#v", keys)
}

func (ms *MiningSession) detectIncrTotalActiveUsersKeys(repo *repository) []string {
	keys := make([]string, 0, int(repo.cfg.MiningSessionDuration.Max/repo.cfg.AdoptionMilestoneSwitch.Duration))
	start, end := ms.EndedAt.Add(-ms.Extension), *ms.EndedAt.Time
	if !ms.LastNaturalMiningStartedAt.Equal(*ms.StartedAt.Time) ||
		(!ms.PreviouslyEndedAt.IsNil() &&
			repo.totalActiveUsersKey(*ms.StartedAt.Time) == repo.totalActiveUsersKey(*ms.PreviouslyEndedAt.Time)) {
		start = start.Add(repo.cfg.AdoptionMilestoneSwitch.Duration)
	}
	start = start.Truncate(repo.cfg.AdoptionMilestoneSwitch.Duration)
	end = end.Truncate(repo.cfg.AdoptionMilestoneSwitch.Duration)
	for start.Before(end) {
		keys = append(keys, repo.totalActiveUsersKey(start))
		start = start.Add(repo.cfg.AdoptionMilestoneSwitch.Duration)
	}
	if ms.PreviouslyEndedAt.IsNil() || repo.totalActiveUsersKey(end) != repo.totalActiveUsersKey(*ms.PreviouslyEndedAt.Time) {
		keys = append(keys, repo.totalActiveUsersKey(end))
	}

	return keys
}

var (
	timeToCheckForAdoptionSwitch   *time.Time
	timeToCheckForAdoptionSwitchMx = new(sync.Mutex)
	timeToCheckShift               = uint64(0)
)

func (r *repository) trySwitchToNextAdoption(ctx context.Context) error {
	if now := *time.Now().Time; !timeToCheckForAdoptionSwitch.IsNil() && timeToCheckForAdoptionSwitch.After(now) {
		return nil
	}
	timeToCheckForAdoptionSwitchMx.Lock()
	defer timeToCheckForAdoptionSwitchMx.Unlock()
	if now := *time.Now().Time; !timeToCheckForAdoptionSwitch.IsNil() && timeToCheckForAdoptionSwitch.After(now) {
		return nil
	}
	currentAdoption, err := GetCurrentAdoption(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to getCurrentAdoption")
	}
	if timeToCheckForAdoptionSwitch.IsNil() {
		timeToCheckForAdoptionSwitch = time.New(currentAdoption.AchievedAt.Add(r.cfg.AdoptionMilestoneSwitch.Duration)) //nolint:lll // .
		if now := *time.Now().Time; !timeToCheckForAdoptionSwitch.IsNil() && timeToCheckForAdoptionSwitch.After(now) {
			return nil
		}
	}
	nextAdoption, err := r.getNextAdoption(ctx, currentAdoption)
	if err != nil || nextAdoption == nil {
		return errors.Wrap(err, "failed to try to get next adoption")
	}
	if err = r.switchToNextAdoption(ctx, nextAdoption); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil // This is a concurrency check. Multiple goroutines will try to update it, but only the 1st will succeed.
		}

		return errors.Wrap(err, "failed to try to get next adoption")
	}

	if err = r.notifyAdoptionChange(ctx, nextAdoption); err != nil {
		revertCtx, revertCancel := context.WithTimeout(context.Background(), requestDeadline)
		defer revertCancel()

		return multierror.Append( //nolint:wrapcheck // Not needed.
			errors.Wrapf(err, "failed notifyAdoptionChange for:%#v", nextAdoption),
			errors.Wrapf(r.revertSwitchToNextAdoption(revertCtx, nextAdoption), //nolint:contextcheck // It might be cancelled.
				"failed to revertSwitchToNextAdoption for:%#v", nextAdoption),
		).ErrorOrNil()
	}

	return nil
}

func (r *repository) notifyAdoptionChange(ctx context.Context, nextAdoption *Adoption[float64]) error {
	oldAdoption, err := getAdoption(ctx, r.db, nextAdoption.Milestone-1)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return errors.Wrapf(err, "failed to get adoption for milestone:%v", nextAdoption.Milestone-1)
	}
	snapshot := &AdoptionSnapshot{Adoption: nextAdoption, Before: oldAdoption}
	if err = r.sendAdoptionSnapshotMessage(ctx, snapshot); err != nil {
		return errors.Wrapf(err, "failed to sendAdoptionSnapshotMessage: %#v", snapshot)
	}

	return nil
}

func (r *repository) mustInitAdoptions(ctx context.Context) {
	responses, err := r.db.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		errs := make([]error, 0, len(r.cfg.AdoptionMilestoneSwitch.ActiveUserMilestones)*3+2) //nolint:gomnd // .
		for ix, milestone := range r.cfg.AdoptionMilestoneSwitch.ActiveUserMilestones {
			id := ix + 1
			key := adoptionsKey(uint64(id))
			if ix == 0 {
				errs = append(errs, pipeliner.HSetNX(ctx, key, "achieved_at", time.Now()).Err())
			}
			errs = append(errs,
				pipeliner.HSetNX(ctx, key, "base_mining_rate", milestone.BaseMiningRate).Err(),
				pipeliner.HSetNX(ctx, key, "milestone", id).Err(),
				pipeliner.HSetNX(ctx, key, "total_active_users", milestone.Users).Err())
		}
		errs = append(errs, pipeliner.SetNX(ctx, "current_adoption_milestone", 1, 0).Err())

		return multierror.Append(nil, errs...).ErrorOrNil() //nolint:wrapcheck // .
	})
	log.Panic(err)
	for _, response := range responses {
		log.Panic(errors.Wrapf(response.Err(), "failed to `%v`", response.FullName()))
	}
}

func getAdoption(ctx context.Context, db storage.DB, milestone uint64) (*Adoption[float64], error) {
	resp, err := storage.Get[Adoption[float64]](ctx, db, adoptionsKey(milestone))
	if err != nil || len(resp) == 0 {
		if err == nil || errors.Is(err, redis.Nil) {
			err = ErrNotFound
		}

		return nil, errors.Wrapf(err, "failed to get the adoption by milestone:%v", milestone)
	}

	return resp[0], nil
}

func getAllAdoptions[DENOM ~string | ~float64](ctx context.Context, db storage.DB) ([]*Adoption[DENOM], error) {
	const max = 20 // We try to get 20 just to be sure we get all of them. We're never going to have more than 20 milestones.
	keys := make([]string, max)
	for ix := uint64(1); ix <= max; ix++ {
		keys[ix-1] = adoptionsKey(ix)
	}
	allAdoptions, err := storage.Get[Adoption[DENOM]](ctx, db, keys...)
	sort.SliceStable(allAdoptions, func(ii, jj int) bool { return allAdoptions[ii].Milestone < allAdoptions[jj].Milestone })

	return allAdoptions, errors.Wrap(err, "failed to get all adoptions")
}

func GetCurrentAdoption(ctx context.Context, db storage.DB) (*Adoption[float64], error) {
	if milestone, err := db.Get(ctx, "current_adoption_milestone").Uint64(); err != nil || milestone == 0 {
		if (err == nil && milestone == 0) || (err != nil && errors.Is(err, redis.Nil)) {
			err = ErrNotFound
		}

		return nil, errors.Wrap(err, "failed to get current_adoption_milestone")
	} else {
		return getAdoption(ctx, db, milestone)
	}
}

func (r *repository) getNextAdoption(ctx context.Context, currentAdoption *Adoption[float64]) (*Adoption[float64], error) { //nolint:funlen // .
	now := time.Now()
	timeToSwitchBasedOnPreviousAdoption := time.New(currentAdoption.AchievedAt.Add(stdlibtime.Duration(r.cfg.AdoptionMilestoneSwitch.ConsecutiveDurationsRequired) * r.cfg.AdoptionMilestoneSwitch.Duration)) //nolint:lll // .
	timeToCheckForAdoptionSwitchBasedOnRemainingDurations := time.New(now.Add(stdlibtime.Duration(atomic.AddUint64(&timeToCheckShift, 1)) * r.cfg.AdoptionMilestoneSwitch.Duration))                          //nolint:lll
	timeToCheckForAdoptionSwitchBasedOnPreviousAdoption := time.New(currentAdoption.AchievedAt.Add(r.cfg.AdoptionMilestoneSwitch.Duration))
	timeToCheckForAdoptionSwitch = maxTime(timeToCheckForAdoptionSwitchBasedOnRemainingDurations, timeToCheckForAdoptionSwitchBasedOnPreviousAdoption) //nolint:lll // .
	if timeToCheckForAdoptionSwitch.After(*timeToSwitchBasedOnPreviousAdoption.Time) {
		atomic.StoreUint64(&timeToCheckShift, 0)
		timeToCheckForAdoptionSwitch = time.New(now.Add(r.cfg.AdoptionMilestoneSwitch.Duration))
	}

	if timeToSwitchBasedOnPreviousAdoption.After(*now.Time) {
		return nil, nil
	}

	nextAdoption, err := getAdoption(ctx, r.db, currentAdoption.Milestone+1)
	if err != nil || !nextAdoption.AchievedAt.IsNil() {
		if err != nil && errors.Is(err, ErrNotFound) {
			timeToCheckForAdoptionSwitch = time.New(time.Now().Add(stdlibtime.Duration(atomic.AddUint64(&timeToCheckShift, 1)) * r.cfg.AdoptionMilestoneSwitch.Duration)) //nolint:lll // .

			return nil, nil
		}

		return nil, errors.Wrapf(err, "failed to get next adoption `%v`", currentAdoption.Milestone+1)
	}
	globalKeys := make([]string, 0, stdlibtime.Duration(r.cfg.AdoptionMilestoneSwitch.ConsecutiveDurationsRequired))
	for duration := stdlibtime.Duration(0); duration < stdlibtime.Duration(r.cfg.AdoptionMilestoneSwitch.ConsecutiveDurationsRequired); duration++ {
		globalKeys = append(globalKeys, r.totalActiveUsersKey(now.Add(-duration*r.cfg.AdoptionMilestoneSwitch.Duration)))
	}
	activeUsersCounters, err := r.db.MGet(ctx, globalKeys...).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get global active users count")
	}
	atLeastOneInvalid := false
	for _, activeUsersCounterAny := range activeUsersCounters {
		if val, ok := activeUsersCounterAny.(string); !ok || val == "" {
			atLeastOneInvalid = true
		} else if activeUsersCounter, pErr := strconv.ParseUint(val, 10, 64); pErr != nil {
			return nil, errors.Wrapf(pErr, "failed to ParseUint: %#v", activeUsersCounterAny)
		} else if activeUsersCounter < nextAdoption.TotalActiveUsers {
			atLeastOneInvalid = true
		}
	}

	if atLeastOneInvalid || len(activeUsersCounters) != len(globalKeys) {
		return nil, nil
	}

	nextAdoption.AchievedAt = now

	return nextAdoption, nil
}

func (r *repository) switchToNextAdoption(ctx context.Context, nextAdoption *Adoption[float64]) error {
	if responses, err := r.db.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		return multierror.Append( //nolint:wrapcheck // .
			pipeliner.HSetNX(ctx, adoptionsKey(nextAdoption.Milestone), "achieved_at", nextAdoption.AchievedAt).Err(),
			pipeliner.Set(ctx, "current_adoption_milestone", fmt.Sprint(nextAdoption.Milestone), 0).Err(),
		).ErrorOrNil()
	}); err != nil {
		return errors.Wrapf(err, "failed to set current_adoption_milestone to milestone: %v", nextAdoption.Milestone)
	} else {
		errs := make([]error, 0, 1+1)
		for _, response := range responses {
			if err = response.Err(); err != nil {
				errs = append(errs, errors.Wrapf(err, "failed to `%v`", response.FullName()))
			} else if boolCmd, ok := response.(*redis.BoolCmd); ok && !boolCmd.Val() {
				errs = append(errs, errors.Wrapf(ErrNotFound, "failed to `%v`", response.FullName()))
			}
		}

		return multierror.Append(nil, errs...).ErrorOrNil() //nolint:wrapcheck // .
	}
}

func (r *repository) revertSwitchToNextAdoption(ctx context.Context, nextAdoption *Adoption[float64]) error {
	if responses, err := r.db.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		return multierror.Append( //nolint:wrapcheck // .
			pipeliner.HDel(ctx, adoptionsKey(nextAdoption.Milestone), "achieved_at").Err(),
			pipeliner.Set(ctx, "current_adoption_milestone", fmt.Sprint(nextAdoption.Milestone-1), 0).Err(),
		).ErrorOrNil()
	}); err != nil {
		return errors.Wrapf(err, "failed to revert set current_adoption_milestone to milestone: %v", nextAdoption.Milestone)
	} else {
		errs := make([]error, 0, 1+1)
		for _, response := range responses {
			if err = response.Err(); err != nil {
				errs = append(errs, errors.Wrapf(err, "failed to `%v`", response.FullName()))
			} else if intCmd, ok := response.(*redis.IntCmd); ok && intCmd.Val() == 0 {
				errs = append(errs, errors.Wrapf(ErrNotFound, "failed to `%v`", response.FullName()))
			}
		}

		return multierror.Append(nil, errs...).ErrorOrNil() //nolint:wrapcheck // .
	}
}

func (r *repository) mustNotifyCurrentAdoption(ctx context.Context) {
	adoption, err := GetCurrentAdoption(ctx, r.db)
	log.Panic(errors.Wrapf(err, "failed to get getCurrentAdoption")) //nolint:revive // Intended.
	snapshot := &AdoptionSnapshot{Adoption: adoption, Before: adoption}
	log.Panic(errors.Wrapf(r.sendAdoptionSnapshotMessage(ctx, snapshot), "failed to sendAdoptionSnapshotMessage: %#v", snapshot))
}

func (r *repository) sendAdoptionSnapshotMessage(ctx context.Context, snapshot *AdoptionSnapshot) error {
	valueBytes, err := json.MarshalContext(ctx, snapshot)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %#v", snapshot)
	}

	msg := &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     strconv.FormatUint(snapshot.Milestone, 10),
		Topic:   r.cfg.MessageBroker.Topics[1].Name,
		Value:   valueBytes,
	}

	responder := make(chan error, 1)
	defer close(responder)
	r.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send %v message to broker, msg:%#v", msg.Topic, snapshot)
}

func adoptionsKey(milestone uint64) string {
	return "adoptions:" + strconv.FormatUint(milestone, 10)
}

func maxTime(first, second *time.Time) *time.Time {
	if first.After(*second.Time) {
		return first
	}

	return second
}
