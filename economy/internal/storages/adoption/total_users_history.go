package adoption

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage"
)

func (r *repository) updateTotalUsersHistory(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "failed to get count of active users because of context failed")
	}
	// Using sql here to get value from global by one query.
	sql := `REPLACE INTO total_users_history (
					minute_timestamp,
                    hour_timestamp,
                    day_timestamp,
                    date_val,
                    total_users
) VALUES(
					:minuteTS,
                    :hourTS,
                    :dayTS,
                    :dateVal,
					(SELECT value FROM global WHERE key = 'TOTAL_USERS')
);`
	now := time.Now().UTC()
	params := map[string]interface{}{
		"minuteTS": now.Unix() / secsInMinute,
		"hourTS":   now.Unix() / secsInHour,
		"dayTS":    now.Unix() / (hoursInDay * secsInHour),
		"dateVal":  now.Format("2006-01-02"),
	}

	return errors.Wrapf(storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)), "failed to add total_users_history %v", params)
}
