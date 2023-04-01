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

	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetAdoptionSummary(ctx context.Context) (as *AdoptionSummary, err error) {
	if as = new(AdoptionSummary); ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	key := r.totalActiveUsersGlobalParentKey(time.Now().Time)
	if as.TotalActiveUsers, err = r.getGlobalUnsignedValue(ctx, key); err != nil && !errors.Is(err, storagev2.ErrNotFound) {
		return nil, errors.Wrapf(err, "failed to get totalActiveUsers getGlobalUnsignedValue for key:%v", key)
	}
	if as.Milestones, err = getAllAdoptions[coin.ICE](ctx, r.dbV2); err != nil {
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
		if errors.Is(err, storagev2.ErrNotFound) {
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
	if err != nil && !errors.Is(err, storagev2.ErrNotFound) {
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
	resp, err := storagev2.Get[Adoption[coin.ICEFlake]](ctx, r.dbV2, `SELECT * FROM adoption WHERE milestone = $1`, milestone)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get the adoption by milestone:%v", milestone)
	}
	if resp.Milestone == 0 {
		return nil, storagev2.ErrNotFound
	}

	return resp, nil
}

func getAllAdoptions[DENOM coin.ICEFlake | coin.ICE](ctx context.Context, dbV2 *storagev2.DB) ([]*Adoption[DENOM], error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	sql := `SELECT  achieved_at,
					base_mining_rate,
					milestone,
					total_active_users
		FROM adoption LIMIT $1`
	resp, err := storagev2.Select[Adoption[DENOM]](ctx, dbV2, sql, lastAdoptionMilestone)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select for all adoptions")
	}

	return resp, nil
}

func (r *repository) getCurrentAdoption(ctx context.Context) (*Adoption[coin.ICEFlake], error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	resp, err := storagev2.Select[Adoption[coin.ICEFlake]](ctx, r.dbV2, currentAdoptionSQL())
	if err != nil {
		return nil, errors.Wrap(err, "failed to select for the current adoption")
	}
	if len(resp) == 0 || resp[0] == nil || resp[0].Milestone == 0 { //nolint:revive // Nope.
		return nil, storagev2.ErrNotFound // Should never happen.
	}

	return resp[0], nil
}

func currentAdoptionSQL() string {
	return `SELECT achieved_at,
			       base_mining_rate,
				   MAX(milestone) AS milestone,
				   total_active_users
		    FROM adoption
		    WHERE achieved_at IS NOT NULL
		    GROUP BY achieved_at, base_mining_rate, total_active_users`
}

func (r *repository) getNextAdoption(ctx context.Context) (*Adoption[coin.ICEFlake], error) { //nolint:funlen // Alot of SQL & mappings.
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "context failed")
	}
	var (
		consecutiveDurationsRequired = stdlibtime.Duration(r.cfg.AdoptionMilestoneSwitch.ConsecutiveDurationsRequired)
		keyParams                    = make([]string, 0, consecutiveDurationsRequired)
		args                         = make([]any, 0)
		now                          = time.Now()
	)
	paramCount := 1
	for duration := stdlibtime.Duration(0); duration < consecutiveDurationsRequired; duration++ {
		relativeTime := now.Add(-duration * r.cfg.AdoptionMilestoneSwitch.Duration)
		args = append(args, r.totalActiveUsersGlobalChildKey(&relativeTime))
		keyParams = append(keyParams, fmt.Sprintf("$%d", paramCount))
		paramCount++
	}
	args = append(args, time.New(now.Add(-consecutiveDurationsRequired*r.cfg.AdoptionMilestoneSwitch.Duration)).Time)
	args = append(args, consecutiveDurationsRequired)
	paramCount += 1
	sql := fmt.Sprintf(`SELECT x.achieved_at,
							   x.base_mining_rate,
							   x.milestone,
							   x.total_active_users
						FROM (SELECT next_adoption.*,
									 COUNT(g.key) AS consecutive_durations
							  FROM global g
									   JOIN (%[2]v) current_adoption
									   CROSS JOIN adoption next_adoption
											ON g.key IN (%[1]v)
												AND next_adoption.milestone = current_adoption.milestone + 1
												AND current_adoption.achieved_at < $%[3]v
												AND g.value >= next_adoption.total_active_users
							  GROUP BY next_adoption.achieved_at, next_adoption.base_mining_rate, next_adoption.milestone
							) x
							  WHERE x.consecutive_durations = $%[4]v::bigint
							    AND x.achieved_at IS NULL`, strings.Join(keyParams, ","), currentAdoptionSQL(), paramCount-1, paramCount)
	resp, err := storagev2.Select[Adoption[coin.ICEFlake]](ctx, r.dbV2, sql, args...)
	if err != nil {
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
			SET achieved_at = $1
			WHERE milestone = $2
			AND achieved_at IS NULL`
	affectedRows, err := storagev2.Exec(ctx, r.dbV2, sql, nextAdoption.AchievedAt.Time, nextAdoption.Milestone)
	if err != nil || affectedRows == 0 {
		return errors.Wrapf(err,
			"failed to update the next adoption to switch to it, achievedAt:%v, milestone: %v", nextAdoption.AchievedAt, nextAdoption.Milestone)
	}

	return nil

}

func (r *repository) revertSwitchToNextAdoption(ctx context.Context, nextAdoption *Adoption[coin.ICEFlake]) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	sql := `UPDATE adoption
			SET achieved_at = NULL
			WHERE milestone = $1
			AND achieved_at IS NOT NULL AND achieved_at = $2`
	resp, err := storagev2.Exec(ctx, r.dbV2, sql, nextAdoption.Milestone, nextAdoption.AchievedAt)
	if err != nil || resp == 0 {
		return errors.Wrapf(err, "failed to revert to update the next adoption to switch to it, milestone:%v, achievedAt:%v",
			nextAdoption.Milestone, nextAdoption.AchievedAt)
	}

	return nil
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
