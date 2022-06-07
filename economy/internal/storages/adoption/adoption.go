// SPDX-License-Identifier: BUSL-1.1

package adoption

import (
	"context"

	"github.com/framey-io/go-tarantool"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/time"
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
	adoptionChangedByHour, err := a.r.getAdoptionsChangedByHoursLastWeek(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to switch active adoption, reading adoptions failed")
	}

	newActiveAdoption := a.calculateNextAdoption(adoptionChangedByHour)
	if newActiveAdoption != nil {
		return errors.Wrapf(a.r.setActiveAdoption(ctx, newActiveAdoption),
			"failed to set active adoption to %#v", newActiveAdoption)
	}

	return nil
}

func (a *adoptionSource) calculateNextAdoption(adoptionsHistory []*adoption) *adoption {
	var newActiveAdoption *adoption
	adoptionChanged := false
	if len(adoptionsHistory) > 0 {
		newActiveAdoption = adoptionsHistory[0]
	} else {
		newActiveAdoption = &adoption{TotalActiveUsers: 0}
	}
	// Check if all hours during the week (last 168 hours) was the same next adoption.
	for _, currentAdoption := range adoptionsHistory {
		if currentAdoption.TotalActiveUsers != newActiveAdoption.TotalActiveUsers {
			adoptionChanged = true

			break
		}
	}
	// If it was changed - it means NOT ALL hours had enough users to switch the adoption, so no need to switch.
	if adoptionChanged {
		return nil
	}

	return newActiveAdoption
}

func (r *repository) getAdoptionsChangedByHoursLastWeek(ctx context.Context) ([]*adoption, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "failed to get adoptions changed during last week because of context failed")
	}
	sql := `
	SELECT adoption.base_hourly_mining_rate, adoption.total_active_users, adoption.active FROM (
      SELECT MIN(total_active_users) AS users_count, -- They are stored for minutes, we take minimum value for each hour --Я
             adoption_history.hour_timestamp
      FROM adoption_history
      WHERE :nowHourTimestamp - hour_timestamp < :oneWeek
      GROUP BY adoption_history.hour_timestamp
  	) as history
      JOIN ADOPTION
           ON users_count >= adoption.total_active_users
	GROUP BY HOUR_TIMESTAMP
`
	params := map[string]interface{}{
		"oneWeek":          adoptionSwitchRequirementsDuration,
		"nowHourTimestamp": time.Now().Unix() / secsInHour,
	}
	var queryResult []*adoption
	if err := r.db.PrepareExecuteTyped(sql, params, &queryResult); err != nil {
		return nil, errors.Wrap(err, "failed to get adoptions changed during last week")
	}

	return queryResult, nil
}

func (r *repository) setActiveAdoption(ctx context.Context, newAdoption *adoption) error {
	if ctx.Err() != nil {
		return errors.Wrapf(ctx.Err(), "failed to set active adoption to %#v because of context failed", newAdoption)
	}
	sql := `
	UPDATE adoption 
		SET active = adoption.total_active_users = :newAdoptionTotalUsers
	WHERE 1 = 1;`
	params := map[string]interface{}{
		"newAdoptionTotalUsers": newAdoption.TotalActiveUsers,
	}

	return errors.Wrapf(storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)), "failed to set current active adoption %#v", newAdoption)
}

//func parseHistory(historyValue string) []uint64 {
//	values := strings.Split(historyValue, ",")
//	uints := make([]uint64, len(values))
//	for i, val := range values {
//		//nolint:errcheck // There is UNSIGNED field in db, no chance if we can get error here.
//		uints[i], _ = strconv.ParseUint(val, base10, bitSize64)
//	}
//
//	return uints
//}
//
//func isAllHistoryAboveAdoptionRequirements(historyValues []uint64, adoptionValue *adoptionWithHistory) bool {
//	allHistoryAboveAdoptionLimit := len(historyValues) > 0 // Default = true.
//	for _, history := range historyValues {
//		if history < adoptionValue.TotalActiveUsers { // And we switch it to false on first non-fitting history entry.
//			allHistoryAboveAdoptionLimit = false
//
//			break
//		}
//	}
//
//	return allHistoryAboveAdoptionLimit
//}
