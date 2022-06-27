// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/time"
)

func (e *economy) GetUserStats(ctx context.Context, days Days) (*UserStats, error) {
	growth, err := e.getDailyUserGrowth(ctx, days)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get users growth by day for %v days", days)
	}
	result := make([]*DailyUserGrowth, 0, len(growth)-1) // Minus current.
	for ts, grow := range growth {
		if ts > 0 {
			result = append(result, grow.DailyUserGrowth())
		}
	}
	currentTotal := growth[0] // Zero timestamp == current values, see SQL below.

	return &UserStats{
		UserGrowth: result,
		Users: UserCounter{
			Total:  currentTotal.Total,
			Active: currentTotal.Active,
		},
	}, nil
}

//nolint:funlen // Long SQL
func (e *economy) getDailyUserGrowth(ctx context.Context, days Days) (map[uint64]*dailyUserGrowth, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get daily user growth failed because context failed")
	}
	sql := `
	SELECT 
		(:now / 86400000000000 - days.COLUMN_31) as day_ts, -- days.COLUMN_31 = count of values below --
		COALESCE(MIN(total.total_users),0) as total,
		COALESCE(MIN(active.total_active_users),0) as active 
	FROM (VALUES (0),(1),(2),(3),(4),(5),(6),(7),(8),(9),(10),(11),(12),
				 (13),(14),(15),(16),(17),(18),(19),(20),(21),(22),(23),(24),(25),(26),(27),(28),(29),(30)) days
	LEFT JOIN total_users_history AS total ON
		total.DAY_TIMESTAMP = :now / 86400000000000 - days.COLUMN_31
	LEFT JOIN adoption_history AS active ON 
		total.hour_timestamp = active.hour_timestamp
		AND total.MINUTE_TIMESTAMP = active.MINUTE_TIMESTAMP
		AND total.DAY_TIMESTAMP = active.HOUR_TIMESTAMP / 24
	WHERE days.column_31 <= :daysCount
	GROUP BY days.COLUMN_31 
    UNION SELECT 0 as day_ts, -- zero day timestamp == current total / active --
           COALESCE((SELECT value FROM GLOBAL WHERE KEY = 'TOTAL_USERS'),0) as total,
           COALESCE((SELECT value FROM GLOBAL WHERE KEY = 'TOTAL_ACTIVE_USERS'),0) as active
	ORDER BY day_ts;
`
	params := map[string]interface{}{
		"now":       time.Now(),
		"daysCount": days,
	}
	var queryResult []*dailyUserGrowth
	if err := e.db.PrepareExecuteTyped(sql, params, &queryResult); err != nil {
		return nil, errors.Wrap(err, "failed to get daily user growth")
	}
	result := map[uint64]*dailyUserGrowth{}
	for _, q := range queryResult {
		result[q.DayTimestamp] = q
	}

	return result, nil
}

func (d *dailyUserGrowth) DailyUserGrowth() *DailyUserGrowth {
	t := stdlibtime.Unix(int64(d.DayTimestamp*secondsInDay), 0)

	return &DailyUserGrowth{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
		Users: UserCounter{
			Total:  d.Total,
			Active: d.Active,
		},
	}
}
