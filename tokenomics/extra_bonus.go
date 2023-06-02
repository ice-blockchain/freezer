// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	extrabonusnotifier "github.com/ice-blockchain/freezer/extra-bonus-notifier"
	"github.com/ice-blockchain/freezer/model"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

type (
	availableExtraBonus struct {
		model.ExtraBonusLastClaimAvailableAtField
		model.ExtraBonusStartedAtField
		model.DeserializedUsersKey
		model.ExtraBonusField
		model.NewsSeenField
		model.ExtraBonusDaysClaimNotAvailableField
	}
)

func (r *repository) ClaimExtraBonus(ctx context.Context, ebs *ExtraBonusSummary) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	id, err := r.getOrInitInternalID(ctx, ebs.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", ebs.UserID)
	}
	now := time.Now()
	stateForUpdate, err := r.detectAvailableExtraBonus(ctx, now, id)
	if err != nil {
		return errors.Wrapf(err, "failed to getAvailableExtraBonus for userID:%v", ebs.UserID)
	}
	ebs.AvailableExtraBonus = stateForUpdate.ExtraBonus

	return errors.Wrapf(storage.Set(ctx, r.db, stateForUpdate), "failed to claim extra bonus:%#v", stateForUpdate)
}

func (r *repository) detectAvailableExtraBonus(ctx context.Context, now *time.Time, id int64) (*availableExtraBonus, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	usr, err := storage.Get[struct {
		model.MiningSessionSoloStartedAtField
		model.MiningSessionSoloEndedAtField
		model.ExtraBonusLastClaimAvailableAtField
		model.ExtraBonusStartedAtField
		model.ExtraBonusDaysClaimNotAvailableField
		model.NewsSeenField
		model.UTCOffsetField
	}](ctx, r.db, model.SerializedUsersKey(id))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get extra bonus state before claiming it for id:%v", id)
	}
	if len(usr) == 0 {
		return nil, ErrNotFound
	}

	return r.getAvailableExtraBonus(now, id, usr[0].ExtraBonusStartedAtField, usr[0].ExtraBonusLastClaimAvailableAtField, usr[0].MiningSessionSoloStartedAtField, usr[0].MiningSessionSoloEndedAtField, usr[0].ExtraBonusDaysClaimNotAvailableField, usr[0].UTCOffsetField, usr[0].NewsSeenField) //nolint:lll // .
}

//nolint:funlen,lll // .
func (r *repository) getAvailableExtraBonus(
	now *time.Time,
	id int64,
	extraBonusStartedAtField model.ExtraBonusStartedAtField,
	extraBonusLastClaimAvailableAtField model.ExtraBonusLastClaimAvailableAtField,
	miningSessionSoloStartedAtField model.MiningSessionSoloStartedAtField,
	miningSessionSoloEndedAtField model.MiningSessionSoloEndedAtField,
	extraBonusDaysClaimNotAvailableField model.ExtraBonusDaysClaimNotAvailableField,
	utcOffsetField model.UTCOffsetField,
	newsSeenField model.NewsSeenField,
) (*availableExtraBonus, error) {
	var (
		extraBonus uint16
		ebUsr      = &extrabonusnotifier.User{
			ExtraBonusStartedAtField: extraBonusStartedAtField,
			UTCOffsetField:           utcOffsetField,
			UpdatedUser: extrabonusnotifier.UpdatedUser{
				DeserializedUsersKey:                 model.DeserializedUsersKey{ID: id},
				ExtraBonusLastClaimAvailableAtField:  extraBonusLastClaimAvailableAtField,
				ExtraBonusDaysClaimNotAvailableField: extraBonusDaysClaimNotAvailableField,
			},
		}
		calculateExtraBonus = func() uint16 {
			return extrabonusnotifier.CalculateExtraBonus(newsSeenField.NewsSeen, ebUsr.ExtraBonusDaysClaimNotAvailable, ebUsr.ExtraBonusIndex-1, now, ebUsr.ExtraBonusLastClaimAvailableAt, miningSessionSoloStartedAtField.MiningSessionSoloStartedAt, miningSessionSoloEndedAtField.MiningSessionSoloEndedAt) //nolint:lll // .
		}
	)
	if !ebUsr.ExtraBonusStartedAt.IsNil() &&
		now.After(*ebUsr.ExtraBonusLastClaimAvailableAt.Time) &&
		ebUsr.ExtraBonusStartedAt.After(*ebUsr.ExtraBonusLastClaimAvailableAt.Time) &&
		ebUsr.ExtraBonusStartedAt.Before(ebUsr.ExtraBonusLastClaimAvailableAt.Add(r.cfg.ExtraBonuses.ClaimWindow)) &&
		now.Before(ebUsr.ExtraBonusLastClaimAvailableAt.Add(r.cfg.ExtraBonuses.ClaimWindow)) {
		return nil, ErrDuplicate
	}
	log.Info(fmt.Sprintf("getAvailableExtraBonus:before:%#v,newsSeen:%v", ebUsr, newsSeenField.NewsSeen))
	defer func() {
		log.Info(fmt.Sprintf("getAvailableExtraBonus:after:%#v,extraBonus:%v", ebUsr, calculateExtraBonus()))
	}()
	if bonusAvailable, bonusClaimable := extrabonusnotifier.IsExtraBonusAvailable(now, r.extraBonusStartDate, r.extraBonusIndicesDistribution, ebUsr); bonusAvailable {
		if extraBonus = calculateExtraBonus(); extraBonus == 0 {
			return nil, ErrNotFound
		} else {
			return &availableExtraBonus{
				ExtraBonusLastClaimAvailableAtField: ebUsr.ExtraBonusLastClaimAvailableAtField,
				ExtraBonusStartedAtField:            model.ExtraBonusStartedAtField{ExtraBonusStartedAt: now},
				DeserializedUsersKey:                ebUsr.DeserializedUsersKey,
				ExtraBonusField:                     model.ExtraBonusField{ExtraBonus: extraBonus},
			}, nil
		}
	} else if !bonusClaimable {
		return nil, ErrNotFound
	} else {
		if extraBonus = calculateExtraBonus(); extraBonus == 0 {
			return nil, ErrNotFound
		} else {
			ebUsr.ExtraBonusLastClaimAvailableAt = nil
		}
	}

	return &availableExtraBonus{
		ExtraBonusLastClaimAvailableAtField: ebUsr.ExtraBonusLastClaimAvailableAtField,
		ExtraBonusStartedAtField:            model.ExtraBonusStartedAtField{ExtraBonusStartedAt: now},
		DeserializedUsersKey:                ebUsr.DeserializedUsersKey,
		ExtraBonusField:                     model.ExtraBonusField{ExtraBonus: extraBonus},
	}, nil
}

func (s *deviceMetadataTableSource) Process(ctx context.Context, msg *messagebroker.Message) error { //nolint:funlen // .
	if ctx.Err() != nil || len(msg.Value) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	type (
		deviceMetadata struct {
			Before *deviceMetadata `json:"before,omitempty"`
			UserID string          `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
			TZ     string          `json:"tz,omitempty" example:"+03:00"`
		}
	)
	var dm deviceMetadata
	if err := json.UnmarshalContext(ctx, msg.Value, &dm); err != nil || dm.UserID == "" || dm.TZ == "" || (dm.Before != nil && dm.Before.TZ == dm.TZ) {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &dm)
	}
	duration, err := stdlibtime.ParseDuration(strings.Replace(dm.TZ+"m", ":", "h", 1))
	if err != nil {
		return errors.Wrapf(err, "invalid timezone:%#v", &dm)
	}
	id, err := s.getOrInitInternalID(ctx, dm.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getOrInitInternalID for %#v", &dm)
	}
	val := &struct {
		model.DeserializedUsersKey
		model.UTCOffsetField
	}{
		DeserializedUsersKey: model.DeserializedUsersKey{ID: id},
		UTCOffsetField:       model.UTCOffsetField{UTCOffset: int16(duration / stdlibtime.Minute)},
	}

	return errors.Wrapf(storage.Set(ctx, s.db, val), "failed to update users' timezone for %#v", &dm)
}

func (s *viewedNewsSource) Process(ctx context.Context, msg *messagebroker.Message) (err error) { //nolint:funlen // .
	if ctx.Err() != nil || len(msg.Value) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	var vn struct {
		UserID string `json:"userId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
		NewsID string `json:"newsId,omitempty" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2"`
	}
	if err = json.UnmarshalContext(ctx, msg.Value, &vn); err != nil || vn.UserID == "" {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &vn)
	}
	duplGuardKey := fmt.Sprintf("news_seen_dupl_guards:%v", vn.UserID)
	if set, dErr := s.db.SetNX(ctx, duplGuardKey, "", s.cfg.MiningSessionDuration.Min).Result(); dErr != nil || !set {
		if dErr == nil {
			dErr = ErrDuplicate
		}

		return errors.Wrapf(dErr, "SetNX failed for news_seen_dupl_guard, %#v", vn)
	}
	defer func() {
		if err != nil {
			undoCtx, cancelUndo := context.WithTimeout(context.Background(), requestDeadline)
			defer cancelUndo()
			err = multierror.Append( //nolint:wrapcheck // .
				err,
				errors.Wrapf(s.db.Del(undoCtx, duplGuardKey).Err(), "failed to del news_seen_dupl_guard key"),
			).ErrorOrNil()
		}
	}()
	id, err := s.getOrInitInternalID(ctx, vn.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getOrInitInternalID for %#v", &vn)
	}

	return errors.Wrapf(s.db.HIncrBy(ctx, model.SerializedUsersKey(id), "news_seen", 1).Err(),
		"failed to increment news_seen for userID:%v,id:%v", vn.UserID, id)
}
