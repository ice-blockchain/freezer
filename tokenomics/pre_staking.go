// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
)

type (
	preStaking struct {
		model.DeserializedUsersKey
		model.PreStakingBonusField
		model.PreStakingAllocationField
	}
)

func (r *repository) GetPreStakingSummary(ctx context.Context, userID string) (*PreStakingSummary, error) {
	ps, _, err := r.getPreStaking(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getPreStaking for userID:%v", userID)
	}

	return &PreStakingSummary{
		PreStaking: &PreStaking{
			Years:      uint64(PreStakingYearsByPreStakingBonuses[ps.PreStakingBonus]),
			Allocation: ps.PreStakingAllocation,
		},
		Bonus: ps.PreStakingBonus,
	}, nil
}

func (r *repository) getPreStaking(ctx context.Context, userID string) (*preStaking, int64, error) {
	id, err := GetOrInitInternalID(ctx, r.db, userID)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	usr, err := storage.Get[preStaking](ctx, r.db, model.SerializedUsersKey(id))
	if err != nil || len(usr) == 0 || usr[0].PreStakingAllocation == 0 {
		if err == nil && (len(usr) == 0 || usr[0].PreStakingAllocation == 0) {
			err = ErrNotFound
		}

		return nil, id, errors.Wrapf(err, "failed to get pre-staking summary for id:%v", id)
	}

	return usr[0], id, nil
}

func (r *repository) StartOrUpdatePreStaking(ctx context.Context, st *PreStakingSummary) error {
	existing, id, err := r.getPreStaking(ctx, st.UserID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return errors.Wrapf(err, "failed to getPreStaking for userID:%v", st.UserID)
	}
	if existing != nil {
		existingYears := uint64(PreStakingYearsByPreStakingBonuses[existing.PreStakingBonus])
		if (existing.PreStakingAllocation == 100 || existing.PreStakingAllocation == st.Allocation) &&
			(existingYears == MaxPreStakingYears || existingYears == st.Years) {
			st.Allocation = existing.PreStakingAllocation
			st.Years = existingYears
			st.Bonus = existing.PreStakingBonus

			return nil
		}
		if existing.PreStakingAllocation > st.Allocation || existingYears > st.Years {
			return ErrDecreasingPreStakingAllocationOrYearsNotAllowed
		}
	}
	st.Bonus = PreStakingBonusesPerYear[uint8(st.Years)]
	existing = &preStaking{
		DeserializedUsersKey:      model.DeserializedUsersKey{ID: id},
		PreStakingBonusField:      model.PreStakingBonusField{PreStakingBonus: st.Bonus},
		PreStakingAllocationField: model.PreStakingAllocationField{PreStakingAllocation: st.Allocation},
	}

	return errors.Wrapf(storage.Set(ctx, r.db, existing), "failed to replace preStaking for %#v", st)
}
