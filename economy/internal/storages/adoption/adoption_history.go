// SPDX-License-Identifier: BUSL-1.1

package adoption

import (
	"context"

	"github.com/framey-io/go-tarantool"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/time"
)

func newRepository(db tarantool.Connector) Repository {
	return &repository{db: db}
}

func (r *repository) updateActiveUsers(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "failed to update count of active users because of context failed")
	}
	// Get it with another query because we're reusing its value in adoption_history and global (key = TOTAL_ACTIVE_USERS).
	activeUsersCount, timestamp, err := r.getActiveUsersCount(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to update active users count because of count failed")
	}
	mins := uint64(timestamp.Unix()) / secsInMinute
	hours := mins / minsInHour
	ah := &adoptionHistory{
		MinuteTimestamp:  mins,
		HoursTimestamp:   hours,
		TotalActiveUsers: activeUsersCount,
	}
	if err = r.db.ReplaceTyped(spaceAdoptionHistory, ah, &[]*adoptionHistory{}); err != nil {
		return errors.Wrapf(err, "failed to update adoption history for moment %v:%v", hours, mins)
	}

	return errors.Wrapf(r.updateGlobalActiveUsersCount(ctx, activeUsersCount), "failed to update global total users count")
}

func (r *repository) getActiveUsersCount(ctx context.Context) (uint64, *time.Time, error) {
	now := time.Now()
	if ctx.Err() != nil {
		return 0, now, errors.Wrap(ctx.Err(), "failed to get count of active users because of context failed")
	}
	var queryResult []*withCount
	sql := `SELECT count(1) AS c FROM user_economy WHERE :now - last_mining_started_at < :inactivityDeadline`
	params := map[string]interface{}{
		"now":                time.Now(),
		"inactivityDeadline": inactivityDeadline,
	}
	if err := r.db.PrepareExecuteTyped(sql, params, &queryResult); err != nil {
		return 0, now, errors.Wrap(err, "failed to get count of active users")
	}
	if len(queryResult) == 0 {
		return uint64(0), now, nil
	}

	return queryResult[0].Count, now, nil
}

func (r *repository) updateGlobalActiveUsersCount(ctx context.Context, count uint64) error {
	newValue := &global{
		Key:   keyTotalActiveUsers,
		Value: count,
	}
	updateOp := []tarantool.Op{
		{Op: "=", Field: fieldGlobalValue, Arg: count},
	}

	return errors.Wrapf(r.db.UpsertAsync(spaceGlobal, newValue, updateOp).GetTyped(&[]*global{}), "failed to update %v key %v", spaceGlobal, keyTotalActiveUsers)
}
