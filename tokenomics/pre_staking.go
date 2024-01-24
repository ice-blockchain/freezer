// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/time"
)

type (
	preStaking struct {
		model.DeserializedUsersKey
		model.PreStakingBonusResettableField
		model.PreStakingAllocationResettableField
	}
	preStakingWithKYC struct {
		model.DeserializedUsersKey
		model.PreStakingBonusField
		model.PreStakingAllocationField
		model.KYCState
	}
)

func (r *repository) GetPreStakingSummary(ctx context.Context, userID string) (*PreStakingSummary, error) {
	ps, _, err := r.getPreStaking(ctx, userID)
	if err != nil || (ps != nil && (ps.PreStakingAllocation == 0 || ps.QuizWasReset(time.Now()))) {
		if err == nil && ps.QuizWasReset(time.Now()) {
			err = ErrPrestakingDisabled
		} else if err == nil && ps.PreStakingAllocation == 0 {
			err = ErrNotFound
		}

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

func (r *repository) getPreStaking(ctx context.Context, userID string) (*preStakingWithKYC, int64, error) {
	id, err := GetOrInitInternalID(ctx, r.db, userID)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	usr, err := storage.Get[preStakingWithKYC](ctx, r.db, model.SerializedUsersKey(id))
	if err != nil || len(usr) == 0 {
		if err == nil && (len(usr) == 0) {
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

	isPrestakingDisabled := false
	if existing != nil {
		isPrestakingDisabled = existing.QuizWasReset(time.Now())
		if !isPrestakingDisabled {
			existingYears := uint64(PreStakingYearsByPreStakingBonuses[existing.PreStakingBonus])
			if existing.PreStakingAllocation == st.Allocation && existingYears == st.Years {
				st.Allocation = existing.PreStakingAllocation
				st.Years = existingYears
				st.Bonus = existing.PreStakingBonus

				return nil
			}
		} else {
			st.Allocation = 0
			st.Years = 0
			st.Bonus = 0
		}
	}
	if st.Allocation == 0 || st.Years == 0 {
		st.Allocation = 0
		st.Years = 0
		st.Bonus = 0
	} else {
		st.Bonus = PreStakingBonusesPerYear[uint8(st.Years)]
	}
	bonus := model.FlexibleFloat64(st.Bonus)
	alloc := model.FlexibleFloat64(st.Allocation)
	prestaking := &preStaking{
		DeserializedUsersKey:                model.DeserializedUsersKey{ID: id},
		PreStakingBonusResettableField:      model.PreStakingBonusResettableField{PreStakingBonus: &bonus},
		PreStakingAllocationResettableField: model.PreStakingAllocationResettableField{PreStakingAllocation: &alloc},
	}
	if !isPrestakingDisabled || (existing != nil && existing.PreStakingAllocation != 0) {
		err = storage.Set(ctx, r.db, prestaking)
	}
	if isPrestakingDisabled && err == nil {
		err = ErrPrestakingDisabled
	}

	return errors.Wrapf(err, "failed to replace preStaking for %#v", st)
}
