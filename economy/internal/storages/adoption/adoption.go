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
		return errors.Wrap(err, "adoption/adoptionSource: failed to update active users")
	}
	if err := a.switchActiveAdoption(ctx); err != nil {
		return errors.Wrap(err, "adoption/adoptionSource: failed to switch active adoption/mining rate")
	}

	return errors.Wrap(a.r.updateTotalUsersHistory(ctx), "adoption/adoptionSource: failed to update total users history")
}

func (a *adoptionSource) switchActiveAdoption(ctx context.Context) error {
	// If the last 168 consecutive hours from adoption_history.hour_timestamp have ALL been >= ANY adoption.total_active_users,
	// then adoption.active of that entry becomes true and the previous active adoption entry becomes false.
	adoptionChangedByHour, err := a.r.getAdoptionsChangedByHoursLastWeek(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to switch active adoption, reading adoptions failed")
	}

	newActiveAdoption := a.calculateNextAdoption(adoptionChangedByHour)
	if newActiveAdoption != nil {
		return errors.Wrapf(a.r.setActiveAdoption(ctx, newActiveAdoption),
			"failed to set active adoption to %#v", newActiveAdoption)
	}

	return nil
}

func (a *adoptionSource) calculateNextAdoption(adoptionsHistory []*adoption) *adoption {
	// We're looking for next  adoption - it will be adoption with minimal count of TotalActiveUsers but greater than current.
	var newActiveAdoption *adoption
	currentActiveAdoption := adoptionsHistory[0] // We'll get current at first because of UNION... And ORDER BY adoptions.active DESC.
	for _, ad := range adoptionsHistory {
		if ad.Active {
			currentActiveAdoption = ad

			continue
		}
		// If in our history by hours if presented any adoption less than current
		// it means that hour has not enough users to switch adoption - no need to switch.
		if ad.TotalActiveUsers < currentActiveAdoption.TotalActiveUsers {
			return nil
		}
		if newActiveAdoption == nil {
			newActiveAdoption = ad
		}
		// Less than minimal but more than currently active.
		if a.checkForMinimalAdoption(ad, newActiveAdoption, currentActiveAdoption) {
			newActiveAdoption = ad
		}
	}
	// If new is current or less - no need to switch anything.
	if newActiveAdoption.TotalActiveUsers <= currentActiveAdoption.TotalActiveUsers {
		return nil
	}

	return newActiveAdoption
}

func (a *adoptionSource) checkForMinimalAdoption(ad, minimal, currentlyActive *adoption) bool {
	lessOrEqToMinimal := ad.TotalActiveUsers <= minimal.TotalActiveUsers
	moreThanCurrent := ad.TotalActiveUsers > currentlyActive.TotalActiveUsers

	return lessOrEqToMinimal && moreThanCurrent
}

func (r *repository) getAdoptionsChangedByHoursLastWeek(ctx context.Context) ([]*adoption, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "failed to get adoptions changed during last week because of context failed")
	}
	sql := `
	SELECT adoption.base_hourly_mining_rate, adoption.total_active_users, adoption.active FROM (
      	SELECT MIN(total_active_users) AS users_count, -- They are stored for minutes, we take minimum value for each hour --
            adoption_history.hour_timestamp
      	FROM adoption_history
      	WHERE :nowHourTimestamp - hour_timestamp < :oneWeek
      	GROUP BY adoption_history.hour_timestamp
  	) AS history
    	JOIN ADOPTION
           ON history.users_count >= adoption.total_active_users
	GROUP BY hour_timestamp
	UNION SELECT ADOPTION.* from ADOPTION WHERE ACTIVE = true -- We need current to compare with --
	ORDER BY adoption.active DESC`

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
		WHERE 1 = 1`
	params := map[string]interface{}{
		"newAdoptionTotalUsers": newAdoption.TotalActiveUsers,
	}

	return errors.Wrapf(storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)), "failed to set current active adoption %#v", newAdoption)
}
