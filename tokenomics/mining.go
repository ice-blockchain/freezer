// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"regexp"
	"sort"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetRankingSummary(ctx context.Context, userID string) (*RankingSummary, error) { //nolint:funlen // .
	id, err := r.getOrInitInternalID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	rank, err := r.db.Get(ctx, fmt.Sprintf("global_rank:%v", id)).Uint64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, "failed to get cached global_rank for id:%v", id)
	}
	if rank == 0 {
		if rank, err = r.db.ZRevRank(ctx, "top_miners", serializedUsersKey(id)).Uint64(); err != nil {
			if errors.Is(err, redis.Nil) {
				return &RankingSummary{GlobalRank: 0}, nil
			}

			return nil, errors.Wrapf(err, "failed to ZRevRank top_miners for userID:%v", userID)
		}
		if err = r.db.SetEx(ctx, fmt.Sprintf("global_rank:%v", id), rank, stdlibtime.Hour).Err(); err != nil {
			return nil, errors.Wrapf(err, "failed to set cached global_rank for id:%v", id)
		}
	}
	if userID != requestingUserID(ctx) {
		if usr, gErr := storage.Get[struct {
			HideRanking bool `redis:"hide_ranking"`
		}](ctx, r.db, serializedUsersKey(id)); gErr != nil || (len(usr) == 1 && usr[0].HideRanking) {
			if gErr == nil {
				gErr = ErrGlobalRankHidden
			}

			return nil, errors.Wrapf(gErr, "failed to get hide_ranking for id:%v", id)
		}
	}

	return &RankingSummary{GlobalRank: rank + 1}, nil
}

const (
	everythingNotAllowedInUsernameRegex = `[^.a-zA-Z0-9]+`
)

var (
	everythingNotAllowedInUsernamePattern = regexp.MustCompile(everythingNotAllowedInUsernameRegex)
)

func (r *repository) GetTopMiners(ctx context.Context, keyword string, limit, offset uint64) (topMiners []*Miner, err error) {
	var ids []string
	if keyword == "" {
		rangeBy := &redis.ZRangeBy{Min: "0", Max: "+inf", Offset: int64(offset), Count: int64(limit)}
		if ids, err = r.db.ZRevRangeByScore(ctx, "top_miners", rangeBy).Result(); err != nil {
			return nil, errors.Wrapf(err, "failed to ZRevRangeByScore for miners for offset:%v,limit:%v", offset, limit)
		}
	} else { //nolint:revive // Nope.
		key := string(everythingNotAllowedInUsernamePattern.ReplaceAll([]byte(strings.ToLower(keyword)), []byte("")))
		if key == "" || !strings.EqualFold(key, keyword) {
			return nil, nil
		}
		if ids, _, err = r.db.SScan(ctx, key, offset, "", int64(limit)).Result(); err != nil {
			return nil, errors.Wrapf(err, "failed to SScan for miners for keyword:%v,offset:%v,limit:%v", key, offset, limit)
		}
	}
	for ix, id := range ids {
		ids[ix] = serializedUsersKey(id)
	}
	topMiners, err = storage.Get[Miner](ctx, r.db, ids...)
	sort.SliceStable(topMiners, func(ii, jj int) bool { return topMiners[ii].Username < topMiners[jj].Username })
	for _, topMiner := range topMiners {
		topMiner.ProfilePictureURL = r.pictureClient.DownloadURL(topMiner.ProfilePictureURL)
	}

	return topMiners, errors.Wrapf(err, "failed to get miners for ids:%#v", ids)
}

//nolint:funlen,lll // .
func (r *repository) GetMiningSummary(ctx context.Context, userID string) (*MiningSummary, error) {
	id, err := r.getOrInitInternalID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	ms, err := storage.Get[miningSummary2](ctx, r.db, serializedUsersKey(id))
	if err != nil || len(ms) == 0 {
		if err == nil {
			err = errors.Wrapf(ErrRelationNotFound, "missing state for id:%v", id)
		}

		return nil, errors.Wrapf(err, "failed to get miningSummary for id:%v", id)
	}
	currentAdoption, err := r.getCurrentAdoption(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to getCurrentAdoption")
	}
	now := time.Now()

	miningStreak := r.calculateMiningStreak(now, ms[0].MiningSessionSoloStartedAt, ms[0].MiningSessionSoloEndedAt)

	return &MiningSummary{
		MiningStreak:                r.calculateMiningStreak(now, ms[0].MiningSessionSoloStartedAt, ms[0].MiningSessionSoloEndedAt),
		MiningSession:               r.calculateMiningSession(now, ms[0].MiningSessionSoloLastStartedAt, ms[0].MiningSessionSoloEndedAt),
		RemainingFreeMiningSessions: r.calculateRemainingFreeMiningSessions(now, ms[0].MiningSessionSoloEndedAt),
		MiningRates:                 r.calculateMiningRateSummaries(currentAdoption.BaseMiningRate, ms[0], now),
		ExtraBonusSummary: ExtraBonusSummary{
			AvailableExtraBonus: r.calculateExtraBonus(ms[0].NewsSeen, now, ms[0].MiningSessionSoloStartedAt, ms[0].MiningSessionSoloEndedAt),
		},
	}, nil
}

func (r *repository) calculateMiningSession(now, start, end *time.Time) (ms *MiningSession) {
	if end == nil || end.Before(*now.Time) {
		return nil
	}
	lastMiningStartedAt := time.New(start.Add((now.Sub(*start.Time) / r.cfg.MiningSessionDuration.Max) * r.cfg.MiningSessionDuration.Max))
	free := start.Add(r.cfg.MiningSessionDuration.Max).Before(*now.Time)

	return &MiningSession{
		StartedAt:                     lastMiningStartedAt,
		EndedAt:                       time.New(lastMiningStartedAt.Add(r.cfg.MiningSessionDuration.Max)),
		Free:                          &free,
		ResettableStartingAt:          time.New(lastMiningStartedAt.Add(r.cfg.MiningSessionDuration.Min)),
		WarnAboutExpirationStartingAt: time.New(lastMiningStartedAt.Add(r.cfg.MiningSessionDuration.WarnAboutExpirationAfter)),
	}
}

//nolint:funlen // A lot of calculations.
func (r *repository) calculateMiningRateSummaries(
	baseMiningRate float64, ms *miningSummary2, now *time.Time,
) (miningRates *MiningRates[MiningRateSummary[string]]) {
	miningRates = new(MiningRates[MiningRateSummary[string]])
	var (
		standardMiningRate         float64
		preStakingMiningRate       float64
		negativeMiningRate         float64
		totalNoPreStakingBonusRate float64

		totalBonusVal                     uint64
		totalNoPreStakingBonusVal         uint64
		positiveTotalNoPreStakingBonusVal uint64
		preStakingBonusVal                uint64

		extraBonus           = uint64(ms.ExtraBonus)
		t0                   = uint64(1)
		t1                   = uint64(ms.ActiveT1Referrals)
		t2                   = uint64(ms.ActiveT2Referrals)
		preStakingAllocation = uint64(ms.PreStakingAllocation)
		preStakingBonus      = uint64(ms.PreStakingBonus)
	)
	if ms.MiningSessionSoloEndedAt.IsNil() { //nolint:gocritic,nestif // Wrong.
		miningRates.Type = NoneMiningRateType
	} else if ms.MiningSessionSoloEndedAt.After(*now.Time) {
		miningRates.Type = PositiveMiningRateType
	} else if ms.BalanceTotal <= 0.0 {
		miningRates.Type = NoneMiningRateType
	} else {
		extraBonus, t0, t1, t2 = 0, 0, 0, 0
		miningRates.Type = NegativeMiningRateType
		negativeMiningRate = ms.SlashingRateSolo + ms.SlashingRateT0 + ms.SlashingRateT1 + ms.SlashingRateT2
	}
	if !ms.ExtraBonusStartedAt.IsNil() &&
		!ms.MiningSessionSoloEndedAt.IsNil() &&
		ms.ExtraBonusStartedAt.Add(r.cfg.ExtraBonuses.Duration).After(*now.Time) &&
		ms.MiningSessionSoloEndedAt.After(*now.Time) {

	}
	if ms.IDT0 <= 0 || ms.MiningSessionT0EndedAt.IsNil() || ms.MiningSessionT0EndedAt.Before(*now.Time) {
		t0 = 0
	}
	miningRates.Base = &MiningRateSummary[string]{
		Amount: fmt.Sprint(baseMiningRate),
	}
	if ms.PreStakingAllocation != 100 {
		var localTotalBonus uint64
		switch miningRates.Type {
		case PositiveMiningRateType:
			standardMiningRate = r.calculateMintedStandardCoins(baseMiningRate, t0, t1, t2, extraBonus, preStakingAllocation, r.cfg.GlobalAggregationInterval.Child, false)
			if standardMiningRate > baseMiningRate {
				localTotalBonus = uint64(((standardMiningRate - baseMiningRate) * 100) / baseMiningRate)
			}
		case NegativeMiningRateType:
			standardMiningRate = (negativeMiningRate * float64(100-ms.PreStakingAllocation)) / float64(100)
		case NoneMiningRateType:
		}
		miningRates.Standard = &MiningRateSummary[string]{
			Amount: fmt.Sprint(baseMiningRate),
			Bonuses: &MiningRateBonuses{
				T1:    uint64(float64((t0*r.cfg.ReferralBonusMiningRates.T0+t1*r.cfg.ReferralBonusMiningRates.T1)*(100-preStakingAllocation)) / float64(100)),
				T2:    uint64(float64(t2*r.cfg.ReferralBonusMiningRates.T2*(100-preStakingAllocation)) / float64(100)),
				Extra: uint64(float64(extraBonus*(100-preStakingAllocation)) / float64(100)),
				Total: localTotalBonus,
			},
		}
	}
	if preStakingAllocation != 0 {
		var localTotalBonus uint64
		switch miningRates.Type {
		case PositiveMiningRateType:
			preStakingMiningRate = r.calculateMintedPreStakingCoins(baseMiningRate, t0, t1, t2, extraBonus, preStakingAllocation, preStakingBonus, r.cfg.GlobalAggregationInterval.Child, false) //nolint:lll // .
			if preStakingMiningRate > baseMiningRate {
				localTotalBonus = uint64(((preStakingMiningRate - baseMiningRate) * float64(100)) / baseMiningRate)
			}
		case NegativeMiningRateType:
			preStakingMiningRate = (negativeMiningRate * float64(preStakingBonus+100) * float64(preStakingAllocation)) / float64(100*100)
		case NoneMiningRateType:
		}
		t1Bonus := float64((t0*r.cfg.ReferralBonusMiningRates.T0+t1*r.cfg.ReferralBonusMiningRates.T1)*preStakingAllocation) / float64(100)
		t2Bonus := float64(t2*r.cfg.ReferralBonusMiningRates.T2*preStakingAllocation) / float64(100)
		extraBonusVal := float64(extraBonus*preStakingAllocation) / float64(100)
		preStakingBonusVal = uint64((float64(preStakingAllocation) * float64(preStakingBonus)) / float64(100))
		miningRates.PreStaking = &MiningRateSummary[string]{
			Amount: fmt.Sprint(preStakingMiningRate),
			Bonuses: &MiningRateBonuses{
				T1:         uint64(t1Bonus),
				T2:         uint64(t2Bonus),
				Extra:      uint64(extraBonusVal),
				PreStaking: preStakingBonusVal,
				Total:      localTotalBonus,
			},
		}
	}
	positiveTotalNoPreStakingBonus := r.calculateMintedStandardCoins(baseMiningRate, t0, t1, t2, extraBonus, 0, r.cfg.GlobalAggregationInterval.Child, false)
	if positiveTotalNoPreStakingBonus > baseMiningRate {
		positiveTotalNoPreStakingBonusVal = uint64(((positiveTotalNoPreStakingBonus - baseMiningRate) * float64(100)) / baseMiningRate)
	}
	switch miningRates.Type {
	case PositiveMiningRateType:
		if standardMiningRate+preStakingMiningRate > baseMiningRate {
			totalBonusVal = uint64(((standardMiningRate + preStakingMiningRate - baseMiningRate) * float64(100)) / baseMiningRate)
		}
		totalNoPreStakingBonusRate = positiveTotalNoPreStakingBonus
		totalNoPreStakingBonusVal = positiveTotalNoPreStakingBonusVal
	case NegativeMiningRateType:
		totalNoPreStakingBonusRate = negativeMiningRate
	case NoneMiningRateType:
	}
	miningRates.Total = &MiningRateSummary[string]{
		Amount: fmt.Sprint(standardMiningRate + preStakingMiningRate),
		Bonuses: &MiningRateBonuses{
			T1:         t0*r.cfg.ReferralBonusMiningRates.T0 + t1*r.cfg.ReferralBonusMiningRates.T1,
			T2:         t2 * r.cfg.ReferralBonusMiningRates.T2,
			Extra:      extraBonus,
			PreStaking: preStakingBonusVal,
			Total:      totalBonusVal,
		},
	}
	miningRates.TotalNoPreStakingBonus = &MiningRateSummary[string]{
		Amount: fmt.Sprint(totalNoPreStakingBonusRate),
		Bonuses: &MiningRateBonuses{
			T1:         t0*r.cfg.ReferralBonusMiningRates.T0 + t1*r.cfg.ReferralBonusMiningRates.T1,
			T2:         t2 * r.cfg.ReferralBonusMiningRates.T2,
			Extra:      extraBonus,
			PreStaking: 0,
			Total:      totalNoPreStakingBonusVal,
		},
	}
	miningRates.PositiveTotalNoPreStakingBonus = &MiningRateSummary[string]{
		Amount: fmt.Sprint(positiveTotalNoPreStakingBonus),
		Bonuses: &MiningRateBonuses{
			T1:         t0*r.cfg.ReferralBonusMiningRates.T0 + t1*r.cfg.ReferralBonusMiningRates.T1,
			T2:         t2 * r.cfg.ReferralBonusMiningRates.T2,
			Extra:      extraBonus,
			PreStaking: 0,
			Total:      positiveTotalNoPreStakingBonusVal,
		},
	}

	return miningRates
}

func (r *repository) calculateMintedStandardCoins(
	baseMiningRate float64,
	t0, t1, t2, extraBonus, preStakingAllocation uint64,
	elapsedNanos stdlibtime.Duration,
	excludeBaseRate bool,
) float64 {
	if preStakingAllocation == 100 || elapsedNanos <= 0 {
		return 0
	}
	var includeBaseMiningRate float64
	if !excludeBaseRate {
		includeBaseMiningRate = float64(100 + extraBonus)
	}
	mintedBase := includeBaseMiningRate +
		float64(t0*r.cfg.ReferralBonusMiningRates.T0) +
		float64(t1*r.cfg.ReferralBonusMiningRates.T1) +
		float64(t2*r.cfg.ReferralBonusMiningRates.T2)
	if mintedBase == 0 {
		return 0
	}
	if elapsedNanos == r.cfg.GlobalAggregationInterval.Child {
		return (baseMiningRate * mintedBase * float64(100-preStakingAllocation)) / float64(100*100)
	}

	return (baseMiningRate * float64(elapsedNanos.Nanoseconds()) * mintedBase * float64(100-preStakingAllocation)) /
		(float64(r.cfg.GlobalAggregationInterval.Child.Nanoseconds()) * float64(100*100))
}

func (r *repository) calculateMintedPreStakingCoins(
	baseMiningRate float64,
	t0, t1, t2, extraBonus, preStakingAllocation, preStakingBonus uint64,
	elapsedNanos stdlibtime.Duration,
	excludeBaseRate bool,
) float64 {
	if preStakingAllocation == 0 || elapsedNanos <= 0 {
		return 0
	}
	var includeBaseMiningRate float64
	if !excludeBaseRate {
		includeBaseMiningRate = float64(100 + extraBonus)
	}
	mintedBase := includeBaseMiningRate +
		float64(t0*r.cfg.ReferralBonusMiningRates.T0) +
		float64(t1*r.cfg.ReferralBonusMiningRates.T1) +
		float64(t2*r.cfg.ReferralBonusMiningRates.T2)
	if mintedBase == 0 {
		return 0
	}
	if elapsedNanos == r.cfg.GlobalAggregationInterval.Child {
		return (baseMiningRate * mintedBase * float64((preStakingBonus+100)*preStakingAllocation)) /
			float64(100*100*100)
	}

	return (baseMiningRate * mintedBase * float64(elapsedNanos.Nanoseconds()) * float64((preStakingBonus+100)*preStakingAllocation)) /
		float64(uint64(r.cfg.GlobalAggregationInterval.Child.Nanoseconds())*100*100*100)
}

func (r *repository) calculateMiningStreak(now, start, end *time.Time) uint64 {
	if start == nil || end == nil || now.After(*end.Time) || now.Before(*start.Time) {
		return 0
	}

	return uint64(now.Sub(*start.Time) / r.cfg.MiningSessionDuration.Max)
}

func (r *repository) calculateRemainingFreeMiningSessions(now, end *time.Time) uint64 {
	if end == nil || now.After(*end.Time) {
		return 0
	}

	return uint64(end.Sub(*now.Time) / r.cfg.MiningSessionDuration.Max)
}
