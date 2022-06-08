// SPDX-License-Identifier: BUSL-1.1

package adoption

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) updateTotalUsersHistory(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "failed to get count of active users because of context failed")
	}
	sql := `REPLACE INTO total_users_history (
					minute_timestamp,
                    hour_timestamp,
                    day_timestamp,
                    date_,
                    total_users
				) VALUES(
					:minuteTS,
                    :hourTS,
                    :dayTS,
                    :date_,
					(SELECT CAST(value AS unsigned) FROM global WHERE key = 'TOTAL_USERS')
				)`
	now := time.Now()
	params := map[string]interface{}{
		"minuteTS": now.Unix() / secsInMinute,
		"hourTS":   now.Unix() / secsInHour,
		"dayTS":    now.Unix() / (hoursInDay * secsInHour),
		"date_":    now.Format(dateFormat),
	}

	return errors.Wrapf(storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)), "failed to update total_users_history %v", params)
}
