// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strconv"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetAdoptionSummary(ctx context.Context) (as *AdoptionSummary, err error) {
	if as = new(AdoptionSummary); ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
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

		return multierror.Append( //nolint:wrapcheck // Not needed.
			errors.Wrapf(err, "failed notifyAdoptionChange for:%#v", nextAdoption),
			errors.Wrapf(r.revertSwitchToNextAdoption(revertCtx, nextAdoption), //nolint:contextcheck // It might be cancelled.
				"failed to revertSwitchToNextAdoption for:%#v", nextAdoption),
		).ErrorOrNil()
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
	resp, err := storage.Get[Adoption[coin.ICEFlake]](ctx, r.db, `SELECT * FROM adoption WHERE milestone = $1`, milestone)

	return resp, errors.Wrapf(err, "failed to get the adoption by milestone:%v", milestone)
}

func getAllAdoptions[DENOM coin.ICEFlake | coin.ICE](ctx context.Context, db *storage.DB) ([]*Adoption[DENOM], error) {
	resp, err := storage.Select[Adoption[DENOM]](ctx, db, `SELECT * FROM adoption`)

	return resp, errors.Wrap(err, "failed to select for all adoptions")
}

func (r *repository) getCurrentAdoption(ctx context.Context) (*Adoption[coin.ICEFlake], error) {
	resp, err := storage.Get[Adoption[coin.ICEFlake]](ctx, r.db, currentAdoptionSQL())

	return resp, errors.Wrap(err, "failed to select for the current adoption")
}

func currentAdoptionSQL() string {
	return `SELECT achieved_at,
			       base_mining_rate,
				   milestone,
				   total_active_users
		    FROM adoption
		    WHERE achieved_at IS NOT NULL
		    ORDER BY milestone DESC
			LIMIT 1`
}

func (r *repository) getNextAdoption(ctx context.Context) (*Adoption[coin.ICEFlake], error) { //nolint:funlen // Alot of SQL & mappings.
	var (
		now                                      = time.Now()
		consecutiveDurationsRequired             = stdlibtime.Duration(r.cfg.AdoptionMilestoneSwitch.ConsecutiveDurationsRequired)
		minimumTimeBetween2MilestoneAchievements = time.New(now.Add(-consecutiveDurationsRequired * r.cfg.AdoptionMilestoneSwitch.Duration)).Time
		globalKeys                               = make([]string, 0, consecutiveDurationsRequired)
	)
	for duration := stdlibtime.Duration(0); duration < consecutiveDurationsRequired; duration++ {
		relativeTime := now.Add(-duration * r.cfg.AdoptionMilestoneSwitch.Duration)
		globalKeys = append(globalKeys, r.totalActiveUsersGlobalChildKey(&relativeTime))
	}
	sql := fmt.Sprintf(`SELECT x.achieved_at,
							   x.base_mining_rate,
							   x.milestone,
							   x.total_active_users
						FROM (SELECT next_adoption.*,
									 COUNT(g.key) AS consecutive_durations
							  FROM global g
									   JOIN (%[1]v) current_adoption
									   CROSS JOIN adoption next_adoption
											ON g.key = ANY($3)
												AND next_adoption.milestone = current_adoption.milestone + 1
												AND current_adoption.achieved_at < $1
												AND g.value >= next_adoption.total_active_users
							  GROUP BY next_adoption.milestone) x
					    WHERE x.consecutive_durations = $2::bigint
						  AND x.achieved_at IS NULL`, currentAdoptionSQL())
	resp, err := storage.Get[Adoption[coin.ICEFlake]](ctx, r.db, sql, minimumTimeBetween2MilestoneAchievements, consecutiveDurationsRequired, globalKeys)
	if err != nil {
		if storage.IsErr(err, storage.ErrNotFound) {
			return nil, nil //nolint:nilnil // Nope.
		}

		return nil, errors.Wrap(err, "failed to select if the next adoption is achieved")
	}
	resp.AchievedAt = now

	return resp, nil
}

func (r *repository) switchToNextAdoption(ctx context.Context, nextAdoption *Adoption[coin.ICEFlake]) error {
	sql := `UPDATE adoption
			SET achieved_at = $1
			WHERE milestone = $2
			AND achieved_at IS NULL`
	_, err := storage.Exec(ctx, r.db, sql, nextAdoption.AchievedAt.Time, nextAdoption.Milestone)

	return errors.Wrapf(err,
		"failed to update the next adoption to switch to it, achievedAt:%v, milestone: %v", nextAdoption.AchievedAt, nextAdoption.Milestone)
}

func (r *repository) revertSwitchToNextAdoption(ctx context.Context, nextAdoption *Adoption[coin.ICEFlake]) error {
	sql := `UPDATE adoption
			SET achieved_at = NULL
			WHERE milestone = $1
			AND achieved_at IS NOT NULL 
			AND achieved_at = $2`
	_, err := storage.Exec(ctx, r.db, sql, nextAdoption.Milestone, nextAdoption.AchievedAt)

	return errors.Wrapf(err, "failed to revert to update the next adoption to switch to it, milestone:%v, achievedAt:%v",
		nextAdoption.Milestone, nextAdoption.AchievedAt)
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
