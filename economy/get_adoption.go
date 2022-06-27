// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"

	"github.com/pkg/errors"
)

func (e *economy) GetAdoption(ctx context.Context) (*Adoption, error) {
	milestones, currentTotal, err := e.getAdoptionMilestones(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get adoption milestones")
	}

	return &Adoption{
		Adoption: milestones,
		Users:    *currentTotal,
	}, nil
}

func (e *economy) getAdoptionMilestones(ctx context.Context) ([]*AdoptionMilestone, *UserCounter, error) {
	if ctx.Err() != nil {
		return nil, nil, errors.Wrap(ctx.Err(), "get adoption milestones failed because context failed")
	}
	sql := `
		SELECT
			adoption.base_hourly_mining_rate,
			adoption.total_active_users,
			0 as total_users, -- Milestones are not linked to total users, they are about active. --
			adoption.active
		FROM adoption
		UNION SELECT '' as base_hourly_mining_rate,
					 COALESCE((SELECT value FROM GLOBAL WHERE KEY = 'TOTAL_USERS'),0) as total_active_users,
					 COALESCE((SELECT value FROM GLOBAL WHERE KEY = 'TOTAL_ACTIVE_USERS'),0) as total_users,
					 false as active
		ORDER BY adoption.total_active_users;`
	var queryResult []*adoptionMilestone
	if err := e.db.PrepareExecuteTyped(sql, map[string]interface{}{}, &queryResult); err != nil {
		return nil, nil, errors.Wrap(err, "failed to get adoption milestones")
	}
	milestones, currentTotal := e.handleAdoptionMilestonesResult(queryResult)

	return milestones, currentTotal, nil
}

func (e *economy) handleAdoptionMilestonesResult(queryResult []*adoptionMilestone) ([]*AdoptionMilestone, *UserCounter) {
	result := make([]*AdoptionMilestone, len(queryResult)-1)
	activePassed := false
	var total *UserCounter
	i := uint(0)
	for _, q := range queryResult {
		// If mining rate is 0 = row contains current total users and active users
		// from global instead of adoption milestone.
		if q.HourlyMiningRate.IsZero() {
			total = &UserCounter{
				Total:  q.TotalUsers,
				Active: q.ActiveUsers,
			}
		} else {
			// Adoptions are sorted by total users and when active one is reached - all next are not achieved yet.
			result[i] = q.AdoptionMilestone(!activePassed)
			if q.Active {
				activePassed = true
			}
			i++
		}
	}

	return result, total
}

func (a *adoptionMilestone) AdoptionMilestone(achieved bool) *AdoptionMilestone {
	return &AdoptionMilestone{
		HourlyMiningRate: a.HourlyMiningRate,
		Users: UserCounter{
			Total:  0, // Milestones are not linked to total users, they are about active.
			Active: a.ActiveUsers,
		},
		Achieved: achieved,
	}
}
