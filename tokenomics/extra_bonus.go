// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"strings"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
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
	bonus, err := r.getAvailableExtraBonus(ctx, now, id)
	if err != nil {
		return errors.Wrapf(err, "failed to getAvailableExtraBonus for userID:%v", ebs.UserID)
	}
	*ebs = *bonus

	return errors.Wrapf(storage.Set(ctx, r.db, &struct {
		ExtraBonusStartedAt *time.Time `redis:"extra_bonus_started_at"`
		deserializedUsersKey
		ExtraBonus uint16 `redis:"extra_bonus"`
		NewsSeen   uint16 `redis:"news_seen"`
	}{
		ExtraBonusStartedAt:  now,
		deserializedUsersKey: deserializedUsersKey{ID: id},
		ExtraBonus:           bonus.AvailableExtraBonus,
	}), "failed to claim extra bonus:%#v", ebs)
}

//nolint:funlen,lll // .
func (r *repository) getAvailableExtraBonus(ctx context.Context, now *time.Time, id int64) (*ExtraBonusSummary, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}

	return &ExtraBonusSummary{}, nil
}

func (r *repository) calculateExtraBonus(
	id int64,
	newsSeen, extraBonusDaysClaimNotAvailable uint16,
	utcOffset int16,
	now, miningSessionSoloStartedAt, miningSessionSoloEndedAt, extraBonusStartedAt, extraBonusLastClaimAvailableAt *time.Time,
) (extraBonus uint16) {
	flatBonus, bonusPercentageRemaining := uint16(0), uint16(0)
	if true {
		return 0
	}
	miningStreak := r.calculateMiningStreak(now, miningSessionSoloStartedAt, miningSessionSoloEndedAt)
	if miningStreak >= uint64(len(r.cfg.ExtraBonuses.MiningStreakValues)) {
		extraBonus += uint16(r.cfg.ExtraBonuses.MiningStreakValues[len(r.cfg.ExtraBonuses.MiningStreakValues)-1])
	} else {
		extraBonus += uint16(r.cfg.ExtraBonuses.MiningStreakValues[miningStreak])
	}
	if newsSeenBonusValues := r.cfg.ExtraBonuses.NewsSeenValues; newsSeen >= uint16(len(newsSeenBonusValues)) {
		extraBonus += uint16(newsSeenBonusValues[len(newsSeenBonusValues)-1])
	} else {
		extraBonus += uint16(newsSeenBonusValues[newsSeen])
	}

	return ((extraBonus + flatBonus) * bonusPercentageRemaining) / 100
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
		deserializedUsersKey
		UTCOffset int16 `redis:"utc_offset"`
	}{
		deserializedUsersKey: deserializedUsersKey{ID: id},
		UTCOffset:            int16(duration / stdlibtime.Minute),
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

	return errors.Wrapf(s.db.HIncrBy(ctx, serializedUsersKey(id), "news_seen", 1).Err(),
		"failed to increment news_seen for userID:%v,id:%v", vn.UserID, id)
}
