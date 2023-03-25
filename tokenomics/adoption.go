// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/go-tarantool-client"
	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetAdoptionSummary(ctx context.Context) (as *AdoptionSummary, err error) {
	if as = new(AdoptionSummary); ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	before := time.Now()
	defer func() {
		if elapsed := stdlibtime.Since(*before.Time); elapsed > 100*stdlibtime.Millisecond {
			log.Info("[response]GetAdoptionSummary took: %v", elapsed)
		}
	}()
	key := r.totalActiveUsersGlobalParentKey(time.Now().Time)
	if as.TotalActiveUsers, err = r.getGlobalUnsignedValue(ctx, key); err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, errors.Wrapf(err, "failed to get totalActiveUsers getGlobalUnsignedValue for key:%v", key)
	}
	if as.Milestones, err = getAllAdoptions[coin.ICE](ctx, r.db); err != nil {
		return nil, errors.Wrapf(err, "failed to get all adoption milestones")
	}

	return
}

func (r *repository) totalActiveUsersGlobalParentKey(date *stdlibtime.Time) string {
	return fmt.Sprintf("%v_%v", totalActiveUsersGlobalKey, date.Format(r.cfg.globalAggregationIntervalParentDateFormat()))
}

func (r *repository) totalActiveUsersGlobalChildKey(date *stdlibtime.Time) string {
	return fmt.Sprintf("%v_%v", totalActiveUsersGlobalKey, date.Format(r.cfg.globalAggregationIntervalChildDateFormat()))
}

func (r *repository) trySwitchToNextAdoption(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	nextAdoption, err := r.getNextAdoption(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to try to get next adoption")
	}
	if nextAdoption == nil {
		return nil
	}
	if err = r.switchToNextAdoption(ctx, nextAdoption); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil // This is a concurrency check. Multiple goroutines will try to update it, but only the 1st will succeed.
		}

		return errors.Wrap(err, "failed to try to get next adoption")
	}

	if err = r.notifyAdoptionChange(ctx, nextAdoption); err != nil {
		revertCtx, revertCancel := context.WithTimeout(context.Background(), requestDeadline)
		defer revertCancel()

		return multierror.Append(
			errors.Wrapf(err, "failed notifyAdoptionChange for:%#v", nextAdoption),
			errors.Wrapf(r.revertSwitchToNextAdoption(revertCtx, nextAdoption), //nolint:contextcheck // It might be cancelled.
				"failed to revertSwitchToNextAdoption for:%#v", nextAdoption))
	}

	return nil
}

func (r *repository) notifyAdoptionChange(ctx context.Context, nextAdoption *Adoption[coin.ICEFlake]) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	adoption, err := r.getAdoption(ctx, nextAdoption.Milestone-1)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return errors.Wrapf(err, "failed to get adoption for milestone:%v", nextAdoption.Milestone-1)
	}
	snapshot := &AdoptionSnapshot{Adoption: nextAdoption, Before: adoption}
	if err = r.sendAdoptionSnapshotMessage(ctx, snapshot); err != nil {
		return errors.Wrapf(err, "failed to sendAdoptionSnapshotMessage: %#v", snapshot)
	}

	return nil
}

func (r *repository) getAdoption(ctx context.Context, milestone uint64) (*Adoption[coin.ICEFlake], error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	resp := new(Adoption[coin.ICEFlake])
	if err := r.db.GetTyped("ADOPTION", "pk_unnamed_ADOPTION_1", tarantool.UintKey{I: uint(milestone)}, resp); err != nil {
		return nil, errors.Wrapf(err, "failed to get the adoption by milestone:%v", milestone)
	}
	if resp.Milestone == 0 {
		return nil, storage.ErrNotFound
	}

	return resp, nil
}

func getAllAdoptions[DENOM coin.ICEFlake | coin.ICE](ctx context.Context, db tarantool.Connector) ([]*Adoption[DENOM], error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	resp := make([]*Adoption[DENOM], 0, lastAdoptionMilestone)
	if err := db.SelectTyped("ADOPTION", "pk_unnamed_ADOPTION_1", 0, lastAdoptionMilestone, tarantool.IterAll, []any{}, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to select for all adoptions")
	}

	return resp, nil
}

func (r *repository) getCurrentAdoption(ctx context.Context) (*Adoption[coin.ICEFlake], error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	resp := make([]*Adoption[coin.ICEFlake], 0, 1)
	if err := r.db.PrepareExecuteTyped(currentAdoptionSQL(), map[string]any{}, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to select for the current adoption")
	}
	if len(resp) == 0 || resp[0] == nil || resp[0].Milestone == 0 { //nolint:revive // Nope.
		return nil, storage.ErrNotFound // Should never happen.
	}

	return resp[0], nil
}

func currentAdoptionSQL() string {
	return `SELECT achieved_at,
			       base_mining_rate,
				   MAX(milestone) AS milestone,
				   total_active_users
		    FROM adoption
		    WHERE achieved_at IS NOT NULL`
}

func (r *repository) getNextAdoption(ctx context.Context) (*Adoption[coin.ICEFlake], error) { //nolint:funlen // Alot of SQL & mappings.
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	var (
		consecutiveDurationsRequired = stdlibtime.Duration(r.cfg.AdoptionMilestoneSwitch.ConsecutiveDurationsRequired)
		keyParams                    = make([]string, 0, consecutiveDurationsRequired)
		params                       = make(map[string]any, cap(keyParams)+1+1)
		now                          = time.Now()
	)
	params["expected_consecutive_durations"] = consecutiveDurationsRequired
	params["minimum_time_for_the_previous_adoption_to_be_achieved"] = time.New(now.Add(-consecutiveDurationsRequired * r.cfg.AdoptionMilestoneSwitch.Duration))
	for duration := stdlibtime.Duration(0); duration < consecutiveDurationsRequired; duration++ {
		relativeTime := now.Add(-duration * r.cfg.AdoptionMilestoneSwitch.Duration)
		params[fmt.Sprintf("total_active_per_duration%v_key", duration)] = r.totalActiveUsersGlobalChildKey(&relativeTime)
		keyParams = append(keyParams, fmt.Sprintf(":total_active_per_duration%v_key", duration))
	}
	sql := fmt.Sprintf(`SELECT x.achieved_at,
							   x.base_mining_rate,
							   x.milestone,
							   x.total_active_users
						FROM (SELECT next_adoption.*,
									 COUNT(g.key) AS consecutive_durations
							  FROM global g
									   JOIN (%[2]v) current_adoption
									   JOIN adoption next_adoption
											ON g.key IN (%[1]v)
												AND next_adoption.milestone = current_adoption.milestone + 1
												AND current_adoption.achieved_at < :minimum_time_for_the_previous_adoption_to_be_achieved
												AND CAST(g.value AS UNSIGNED) >= next_adoption.total_active_users) x
							  WHERE x.consecutive_durations == :expected_consecutive_durations
							    AND x.achieved_at IS NULL`, strings.Join(keyParams, ","), currentAdoptionSQL())
	resp := make([]*Adoption[coin.ICEFlake], 0, 1)
	before2 := time.Now()
	defer func() {
		if elapsed := stdlibtime.Since(*before2.Time); elapsed > 100*stdlibtime.Millisecond {
			log.Info("[response]getNextAdoption SQL took: %v", elapsed)
		}
	}()
	if err := r.db.PrepareExecuteTyped(sql, params, &resp); err != nil {
		return nil, errors.Wrap(err, "failed to select if the next adoption is achieved")
	}
	if len(resp) == 0 || resp[0] == nil || resp[0].Milestone == 0 { //nolint:revive // Nope.
		return nil, nil //nolint:nilnil // Nope.
	}
	resp[0].AchievedAt = now

	return resp[0], nil
}

func (r *repository) switchToNextAdoption(ctx context.Context, nextAdoption *Adoption[coin.ICEFlake]) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	sql := `UPDATE adoption
			SET achieved_at = :achieved_at
			WHERE milestone = :milestone
			AND achieved_at IS NULL`
	params := make(map[string]any, 1+1)
	params["milestone"] = nextAdoption.Milestone
	params["achieved_at"] = nextAdoption.AchievedAt

	return errors.Wrapf(storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)),
		"failed to update the next adoption to switch to it, params:%#v", params)
}

func (r *repository) revertSwitchToNextAdoption(ctx context.Context, nextAdoption *Adoption[coin.ICEFlake]) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	sql := `UPDATE adoption
			SET achieved_at = NULL
			WHERE milestone = :milestone
			AND achieved_at IS NOT NULL AND achieved_at = :achieved_at`
	params := make(map[string]any, 1+1)
	params["milestone"] = nextAdoption.Milestone
	params["achieved_at"] = nextAdoption.AchievedAt

	return errors.Wrapf(storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)),
		"failed to revert to update the next adoption to switch to it, params:%#v", params)
}

func (r *repository) mustNotifyCurrentAdoption(ctx context.Context) {
	adoption, err := r.getCurrentAdoption(ctx)
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
