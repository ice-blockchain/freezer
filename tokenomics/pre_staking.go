// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v3"
)

type (
	preStaking struct {
		deserializedUsersKey
		PreStakingBonus      uint16 `redis:"pre_staking_bonus"`
		PreStakingAllocation uint16 `redis:"pre_staking_allocation"`
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
			Allocation: uint64(ps.PreStakingAllocation),
		},
		Bonus: uint64(ps.PreStakingBonus),
	}, nil
}

func (r *repository) getPreStaking(ctx context.Context, userID string) (*preStaking, int64, error) {
	id, err := r.getOrInitInternalID(ctx, userID)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	usr, err := storage.Get[preStaking](ctx, r.db, serializedUsersKey(id))
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
		if (existing.PreStakingAllocation == uint16(percentage100) || existing.PreStakingAllocation == uint16(st.Allocation)) &&
			(existingYears == MaxPreStakingYears || existingYears == st.Years) {
			st.Allocation = uint64(existing.PreStakingAllocation)
			st.Years = existingYears
			st.Bonus = uint64(existing.PreStakingBonus)

			return nil
		}
		if existing.PreStakingAllocation > uint16(st.Allocation) || existingYears > st.Years {
			return ErrDecreasingPreStakingAllocationOrYearsNotAllowed
		}
	}
	st.Bonus = uint64(PreStakingBonusesPerYear[uint8(st.Years)])
	existing = &preStaking{
		deserializedUsersKey: deserializedUsersKey{ID: id},
		PreStakingBonus:      uint16(st.Bonus),
		PreStakingAllocation: uint16(st.Allocation),
	}

	return errors.Wrapf(storage.Set(ctx, r.db, existing), "failed to replace preStaking for %#v", st)
}
