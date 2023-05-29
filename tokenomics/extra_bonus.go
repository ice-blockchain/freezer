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
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/freezer/model"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
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
	ebs.AvailableExtraBonus, err = r.getAvailableExtraBonus(ctx, now, id)
	if err != nil {
		return errors.Wrapf(err, "failed to getAvailableExtraBonus for userID:%v", ebs.UserID)
	}

	return errors.Wrapf(storage.Set(ctx, r.db, &struct {
		model.ExtraBonusStartedAtField
		model.DeserializedUsersKey
		model.ExtraBonusField
		model.NewsSeenField
	}{
		ExtraBonusStartedAtField: model.ExtraBonusStartedAtField{ExtraBonusStartedAt: now},
		DeserializedUsersKey:     model.DeserializedUsersKey{ID: id},
		ExtraBonusField:          model.ExtraBonusField{ExtraBonus: ebs.AvailableExtraBonus},
	}), "failed to claim extra bonus:%#v", ebs)
}

//nolint:funlen,lll // .
func (r *repository) getAvailableExtraBonus(ctx context.Context, now *time.Time, id int64) (uint16, error) {
	if ctx.Err() != nil {
		return 0, errors.Wrap(ctx.Err(), "unexpected deadline")
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
		return 0, errors.Wrapf(err, "failed to get extra bonus state before claiming it for id:%v", id)
	}
	if len(usr) == 0 {
		return 0, ErrNotFound
	}
	var (
		extraBonusLastClaimAvailableAt  = usr[0].ExtraBonusLastClaimAvailableAt
		extraBonusStartedAt             = usr[0].ExtraBonusStartedAt
		extraBonusDaysClaimNotAvailable = usr[0].ExtraBonusDaysClaimNotAvailable
	)
	if extraBonusLastClaimAvailableAt.IsNil() ||
		extraBonusLastClaimAvailableAt.After(*now.Time) ||
		extraBonusLastClaimAvailableAt.Add(r.cfg.ExtraBonuses.ClaimWindow).Before(*now.Time) {
		return 0, ErrNotFound
	}
	if !extraBonusStartedAt.IsNil() &&
		extraBonusStartedAt.After(*extraBonusLastClaimAvailableAt.Time) &&
		extraBonusStartedAt.Before(extraBonusLastClaimAvailableAt.Add(r.cfg.ExtraBonuses.ClaimWindow)) {
		return 0, ErrDuplicate
	}
	extraBonus := r.calculateExtraBonus(usr[0].NewsSeen, extraBonusDaysClaimNotAvailable, usr[0].UTCOffset, now, extraBonusLastClaimAvailableAt, usr[0].MiningSessionSoloStartedAt, usr[0].MiningSessionSoloEndedAt) //nolint:lll // .
	if extraBonus == 0 {
		return 0, ErrNotFound
	}

	return extraBonus, nil
}

func (r *repository) calculateExtraBonus(
	newsSeen, extraBonusDaysClaimNotAvailable uint16,
	utcOffset int16,
	now, extraBonusLastClaimAvailableAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt *time.Time,
) uint16 {
	const networkDelayDelta = 1.333
	var (
		firstDelayedClaimPenaltyWindow = int64(float64(r.cfg.ExtraBonuses.DelayedClaimPenaltyWindow.Nanoseconds()) * networkDelayDelta)
		newsSeenBonusValues            = r.cfg.ExtraBonuses.NewsSeenValues
		miningStreakValues             = r.cfg.ExtraBonuses.MiningStreakValues
		utcOffsetDuration              = stdlibtime.Duration(utcOffset) * r.cfg.ExtraBonuses.UTCOffsetDuration
		location                       = stdlibtime.FixedZone(utcOffsetDuration.String(), int(utcOffsetDuration.Seconds()))
		extraBonusIndex                = uint16(extraBonusLastClaimAvailableAt.In(location).Sub(r.extraBonusStartDate.In(location)) / r.cfg.ExtraBonuses.Duration)
		bonusPercentageRemaining       = 100 + extraBonusDaysClaimNotAvailable
		miningStreak                   = r.calculateMiningStreak(now, miningSessionSoloStartedAt, miningSessionSoloEndedAt)
		flatBonusValue                 = r.cfg.ExtraBonuses.FlatValues[extraBonusIndex]
	)
	if flatBonusValue == 0 {
		return 0
	}
	if delay := now.Sub(*extraBonusLastClaimAvailableAt.Time); delay.Nanoseconds() > firstDelayedClaimPenaltyWindow {
		bonusPercentageRemaining -= 25 * uint16(delay/r.cfg.ExtraBonuses.DelayedClaimPenaltyWindow)
	}
	if miningStreak >= uint64(len(miningStreakValues)) {
		miningStreak = uint64(len(miningStreakValues) - 1)
	}
	if newsSeen >= uint16(len(newsSeenBonusValues)) {
		newsSeen = uint16(len(newsSeenBonusValues) - 1)
	}

	return ((flatBonusValue + miningStreakValues[miningStreak] + newsSeenBonusValues[newsSeen]) * bonusPercentageRemaining) / 100
}

func MustGetExtraBonusStartDate(ctx context.Context, db storage.DB) (extraBonusStartDate *time.Time) {
	extraBonusStartDateString, err := db.Get(ctx, "extra_bonus_start_date").Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}
	log.Panic(errors.Wrap(err, "failed to get extra_bonus_start_date"))
	if extraBonusStartDateString != "" {
		extraBonusStartDate = new(time.Time)
		log.Panic(errors.Wrapf(extraBonusStartDate.UnmarshalText([]byte(extraBonusStartDateString)), "failed to parse extra_bonus_start_date `%v`", extraBonusStartDateString)) //nolint:lll // .

		return
	}
	extraBonusStartDate = time.New(stdlibtime.Now().Truncate(24 * stdlibtime.Hour))
	set, sErr := db.SetNX(ctx, "extra_bonus_start_date", extraBonusStartDate, 0).Result()
	log.Panic(errors.Wrap(sErr, "failed to set extra_bonus_start_date"))
	if !set {
		return MustGetExtraBonusStartDate(ctx, db)
	}

	return extraBonusStartDate
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
