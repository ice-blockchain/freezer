// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/model"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
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
		model.KYCQuizResetAtField
		model.KYCQuizResetAtAppliedField
	}
)

func (p *preStakingWithKYC) PreStakingAlreadyDisabled() bool {
	return p.KYCQuizResetAt != nil && len(*p.KYCQuizResetAt) > 0 && p.KYCQuizResetAt.Equals(p.KYCQuizResetAtApplied)
}

func (r *repository) GetPreStakingSummary(ctx context.Context, userID string) (*PreStakingSummary, error) {
	ps, _, err := r.getPreStaking(ctx, userID)
	if err != nil || (ps != nil && (ps.PreStakingAllocation == 0 || ps.PreStakingAlreadyDisabled())) {
		if err == nil && ps.PreStakingAlreadyDisabled() {
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
		isPrestakingDisabled = existing.PreStakingAlreadyDisabled()
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
		if err == nil {
			exisingBonus, existingAllocation := 0.0, 0.0
			if existing != nil {
				exisingBonus, existingAllocation = existing.PreStakingBonus, existing.PreStakingAllocation
			}
			if sErr := r.sendPreStakingSnapshotMessage(ctx, exisingBonus, existingAllocation, st); sErr != nil {
				if existing == nil {
					existing = &preStakingWithKYC{
						DeserializedUsersKey:      model.DeserializedUsersKey{ID: id},
						PreStakingBonusField:      model.PreStakingBonusField{0},
						PreStakingAllocationField: model.PreStakingAllocationField{0},
					}
				}
				bonus = model.FlexibleFloat64(existing.PreStakingBonus)
				alloc = model.FlexibleFloat64(existing.PreStakingAllocation)
				prestaking = &preStaking{
					DeserializedUsersKey:                model.DeserializedUsersKey{ID: id},
					PreStakingBonusResettableField:      model.PreStakingBonusResettableField{PreStakingBonus: &bonus},
					PreStakingAllocationResettableField: model.PreStakingAllocationResettableField{PreStakingAllocation: &alloc},
				}
				rollbackCtx, cancel := context.WithTimeout(context.Background(), requestDeadline)
				defer cancel()
				rErr := storage.Set(rollbackCtx, r.db, prestaking)

				return multierror.Append(
					sErr,
					rErr,
				).ErrorOrNil()
			}
		}
	}
	if isPrestakingDisabled && err == nil {
		err = ErrPrestakingDisabled
	}

	return errors.Wrapf(err, "failed to replace preStaking for %#v", st)
}

func PreStakingMessage(ctx context.Context, producer, topic, userID string, existingBonus, existingAllocation float64, newPrestaking *PreStakingSummary) *messagebroker.Message {
	if newPrestaking == nil {
		newPrestaking = &PreStakingSummary{
			Bonus: 0,
			PreStaking: &PreStaking{
				UserID:     userID,
				Years:      0,
				Allocation: 0,
			},
		}
	}
	if newPrestaking.Years == 0 {
		newPrestaking.Years = uint64(PreStakingYearsByPreStakingBonuses[newPrestaking.Bonus])
	}
	snapshot := &PreStakingSnapshot{
		PreStakingSummary: newPrestaking,
		Before: &PreStakingSummary{
			PreStaking: &PreStaking{
				UserID:     newPrestaking.UserID,
				Years:      uint64(PreStakingYearsByPreStakingBonuses[existingBonus]),
				Allocation: existingAllocation,
			},
			Bonus: existingBonus,
		},
	}

	valueBytes, err := json.MarshalContext(ctx, snapshot)
	log.Panic(errors.Wrapf(err, "failed to marshal %#v", newPrestaking))

	return &messagebroker.Message{
		Headers: map[string]string{"producer": producer},
		Key:     newPrestaking.UserID,
		Topic:   topic,
		Value:   valueBytes,
	}
}

func (r *repository) sendPreStakingSnapshotMessage(ctx context.Context, existingBonus, existingAllocation float64, st *PreStakingSummary) error {
	msg := PreStakingMessage(ctx, freezerRefrigerantProducer, r.cfg.MessageBroker.Topics[6].Name, st.UserID, existingBonus, existingAllocation, st)
	responder := make(chan error, 1)
	defer close(responder)
	r.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send `%v` message to broker", msg.Topic)
}
