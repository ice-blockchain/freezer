package adoption

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/framey-io/go-tarantool"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
)

func New(db tarantool.Connector) messagebroker.Processor {
	return &adoptionSource{r: newRepository(db).(*repository)}
}

func (a *adoptionSource) Process(ctx context.Context, _ *messagebroker.Message) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	if err := a.r.updateActiveUsers(ctx); err != nil {
		return errors.Wrapf(err, "adoption/adoptionSource: failed to update active users")
	}
	if err := a.switchActiveAdoption(ctx); err != nil {
		return errors.Wrapf(err, "adoption/adoptionSource: failed to switch active adoption/mining rate")
	}

	return errors.Wrapf(a.r.updateTotalUsersHistory(ctx), "adoption/adoptionSource: failed to update total users history")
}

func (a *adoptionSource) switchActiveAdoption(ctx context.Context) error {
	// If the last 168 consecutive hours from adoption_history.hour_timestamp have ALL been >= ANY adoption.total_active_users,
	// then adoption.active of that entry becomes true and the previous active adoption entry becomes false.
	adoptionsWithHistory, err := a.r.getAdoptionsWithHistoryForLastWeek(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to switch active adoption, reading adoptions failed")
	}

	newActiveAdoption := a.calculateNextAdoption(adoptionsWithHistory)
	if newActiveAdoption != nil {
		return errors.Wrapf(a.r.setActiveAdoption(ctx, newActiveAdoption),
			"failed to set active adoption to %#v", newActiveAdoption)
	}

	return nil
}

func (a *adoptionSource) calculateNextAdoption(adoptionsWithHistory []*adoptionWithHistory) *adoption {
	var newActiveAdoption *adoption
	for _, adoptionValue := range adoptionsWithHistory {
		historyValues := parseHistory(adoptionValue.HistoryByHour)
		allHistoryAboveAdoptionRequirements := isAllHistoryAboveAdoptionRequirements(historyValues, adoptionValue)
		if allHistoryAboveAdoptionRequirements {
			if !adoptionValue.Active {
				newActiveAdoption = &adoptionValue.adoption
			} else {
				continue
			}
			// We switch incrementally only to next adoption (cannot skip them).
			break
		}
	}

	return newActiveAdoption
}

func (r *repository) getAdoptionsWithHistoryForLastWeek(ctx context.Context) ([]*adoptionWithHistory, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "failed to get last week adoption history because of context failed")
	}
	sql := fmt.Sprintf(`
	SELECT 
	(
		SELECT GROUP_CONCAT(history.users_count) FROM (
			SELECT MIN(total_active_users) AS users_count FROM adoption_history -- They are stored for minutes, we take minimum value for each hour --
			GROUP BY adoption_history.hour_timestamp
			ORDER BY adoption_history.hour_timestamp DESC
			LIMIT %v
		) history
	) AS users_by_hour,
	adoption.base_hourly_mining_rate, adoption.total_active_users, adoption.active
	FROM adoption 
	WHERE adoption.total_active_users >= (SELECT value FROM global WHERE key = 'TOTAL_ACTIVE_USERS')	
	ORDER BY adoption.total_active_users;
`, adoptionSwitchRequirementsDuration)
	var queryResult []*adoptionWithHistory
	if err := r.db.PrepareExecuteTyped(sql, map[string]interface{}{}, &queryResult); err != nil {
		return nil, errors.Wrap(err, "failed to get last week adoption history")
	}

	return queryResult, nil
}

func (r *repository) setActiveAdoption(ctx context.Context, newAdoption *adoption) error {
	if ctx.Err() != nil {
		return errors.Wrapf(ctx.Err(), "failed to set active adoption to %#v because of context failed", newAdoption)
	}
	sql := `
	UPDATE adoption SET active = CASE
		WHEN adoption.total_active_users = :newAdoptionTotalUsers AND adoption.base_hourly_mining_rate = :newAdoptionRate THEN TRUE
		ELSE FALSE END
	WHERE 1 = 1;`
	params := map[string]interface{}{
		"newAdoptionTotalUsers": newAdoption.TotalActiveUsers,
		"newAdoptionRate":       newAdoption.BaseHourlyMiningRate,
	}

	return errors.Wrapf(storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)), "failed to set current active adoption %#v", newAdoption)
}

func parseHistory(historyValue string) []uint64 {
	values := strings.Split(historyValue, ",")
	uints := make([]uint64, len(values))
	for i, val := range values {
		//nolint:errcheck // There is UNSIGNED field in db, no chance if we can get error here.
		uints[i], _ = strconv.ParseUint(val, base10, bitSize64)
	}

	return uints
}

func isAllHistoryAboveAdoptionRequirements(historyValues []uint64, adoptionValue *adoptionWithHistory) bool {
	allHistoryAboveAdoptionLimit := len(historyValues) > 0 // Default = true.
	for _, history := range historyValues {
		if history < adoptionValue.TotalActiveUsers { // And we switch it to false on first non-fitting history entry.
			allHistoryAboveAdoptionLimit = false

			break
		}
	}

	return allHistoryAboveAdoptionLimit
}
