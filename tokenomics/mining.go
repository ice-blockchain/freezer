// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetRankingSummary(ctx context.Context, userID string) (*RankingSummary, error) { //nolint:funlen // .
	id, err := GetOrInitInternalID(ctx, r.db, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	rank, err := r.db.Get(ctx, fmt.Sprintf("global_rank:%v", id)).Uint64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, "failed to get cached global_rank for id:%v", id)
	}
	if rank == 0 {
		if rank, err = r.db.ZRevRank(ctx, "top_miners", model.SerializedUsersKey(id)).Uint64(); err != nil {
			if errors.Is(err, redis.Nil) {
				return &RankingSummary{GlobalRank: rank}, nil
			}

			return nil, errors.Wrapf(err, "failed to ZRevRank top_miners for userID:%v", userID)
		}
		var expiration stdlibtime.Duration
		if r.cfg.MiningSessionDuration.Max == 24*stdlibtime.Hour {
			expiration = 5 * stdlibtime.Minute
		} else {
			expiration = 5 * stdlibtime.Second
		}
		rank++
		if err = r.db.SetEx(ctx, fmt.Sprintf("global_rank:%v", id), rank, expiration).Err(); err != nil {
			return nil, errors.Wrapf(err, "failed to set cached global_rank for id:%v", id)
		}
	}
	if userID != requestingUserID(ctx) {
		if usr, gErr := storage.Get[struct{ model.HideRankingField }](ctx, r.db, model.SerializedUsersKey(id)); gErr != nil || (len(usr) == 1 && usr[0].HideRanking) {
			if gErr == nil {
				gErr = ErrGlobalRankHidden
			}

			return nil, errors.Wrapf(gErr, "failed to get hide_ranking for id:%v", id)
		}
	}

	return &RankingSummary{GlobalRank: rank}, nil
}

const (
	everythingNotAllowedInUsernameRegex = `[^.a-zA-Z0-9]+`
)

var (
	everythingNotAllowedInUsernamePattern = regexp.MustCompile(everythingNotAllowedInUsernameRegex)
)

//nolint:funlen // .
func (r *repository) GetTopMiners(ctx context.Context, keyword string, limit, offset uint64) (topMiners []*Miner, nextOffset uint64, err error) {
	var (
		ids           []string
		sortTopMiners func(int, int) bool
	)
	nextOffset = 1
	topMiners = make([]*Miner, 0)
	for len(topMiners) < int(limit) && nextOffset != 0 {
		if keyword == "" {
			sortTopMiners = func(ii, jj int) bool { return topMiners[ii].balance > topMiners[jj].balance }
			rangeBy := &redis.ZRangeBy{Min: "0", Max: "+inf", Offset: int64(offset), Count: int64(limit)}
			if ids, err = r.db.ZRevRangeByScore(ctx, "top_miners", rangeBy).Result(); err != nil {
				return nil, 0, errors.Wrapf(err, "failed to ZRevRangeByScore for miners for offset:%v,limit:%v", offset, limit)
			}
			if len(ids) > 0 {
				nextOffset = offset + limit
			} else {
				nextOffset = 0
			}
		} else { //nolint:revive // Nope.
			sortTopMiners = func(ii, jj int) bool { return topMiners[ii].Username < topMiners[jj].Username }
			key := string(everythingNotAllowedInUsernamePattern.ReplaceAll([]byte(strings.ToLower(keyword)), []byte("")))
			if key == "" || !strings.EqualFold(key, keyword) {
				return make([]*Miner, 0, 0), 0, nil
			}
			if ids, nextOffset, err = r.db.SScan(ctx, "lookup:"+key, offset, "", int64(limit)).Result(); err != nil {
				return nil, 0, errors.Wrapf(err, "failed to SScan for miners for keyword:%v,offset:%v,limit:%v", key, offset, limit)
			}
		}
		dedupl := make(map[string]struct{}, len(ids))
		for _, id := range ids {
			dedupl[id] = struct{}{}
		}
		ids = ids[:0]
		for id := range dedupl {
			ids = append(ids, id)
		}
		if len(ids) == 0 {
			break
		}
		resp, err := storage.Get[struct {
			model.UserIDField
			model.UsernameField
			model.ProfilePictureNameField
			model.BalanceTotalStandardField
			model.BalanceTotalPreStakingField
			model.HideRankingField
		}](ctx, r.db, ids...)
		if err != nil {
			return nil, 0, errors.Wrapf(err, "failed to get miners for ids:%#v", ids)
		}
		for _, topMiner := range resp {
			if topMiner.HideRanking {
				continue
			}
			topMiners = append(topMiners, &Miner{
				Balance:           fmt.Sprintf(floatToStringFormatter, topMiner.BalanceTotalStandard+topMiner.BalanceTotalPreStaking),
				balance:           topMiner.BalanceTotalStandard + topMiner.BalanceTotalPreStaking,
				UserID:            topMiner.UserID,
				Username:          topMiner.Username,
				ProfilePictureURL: r.pictureClient.DownloadURL(topMiner.ProfilePictureName),
			})
		}
		offset = nextOffset
	}
	sort.SliceStable(topMiners, sortTopMiners)
	// This is the hack for preview GetTopMiners call with limit = 5 (we know this case for sure). We just remove all extra rows from result.
	if limit == 5 {
		bound := limit
		if len(topMiners) < int(limit) {
			bound = uint64(len(topMiners))
		}
		topMiners = topMiners[:bound]
	}

	return topMiners, nextOffset, nil
}

//nolint:funlen // .
func (r *repository) GetMiningSummary(ctx context.Context, userID string) (*MiningSummary, error) {
	id, err := GetOrInitInternalID(ctx, r.db, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	now := time.Now()
	ms, err := storage.Get[struct {
		model.MiningSessionSoloLastStartedAtField
		model.MiningSessionSoloStartedAtField
		model.MiningSessionSoloEndedAtField
		model.ExtraBonusStartedAtField
		model.ExtraBonusLastClaimAvailableAtField
		model.UserIDField
		model.BalanceTotalStandardField
		model.BalanceTotalPreStakingField
		model.SlashingRateSoloField
		model.SlashingRateT0Field
		model.SlashingRateT1Field
		model.SlashingRateT2Field
		model.PreStakingBonusField
		model.ExtraBonusField
		model.IDT0Field
		model.ActiveT1ReferralsField
		model.ActiveT2ReferralsField
		model.ExtraBonusDaysClaimNotAvailableResettableField
		model.NewsSeenField
		model.PreStakingAllocationField
		model.UTCOffsetField
	}](ctx, r.db, model.SerializedUsersKey(id))
	if err != nil || len(ms) == 0 {
		if err == nil {
			err = errors.Wrapf(ErrRelationNotFound, "missing state for id:%v", id)
		}

		return nil, errors.Wrapf(err, "failed to get miningSummary for id:%v", id)
	}
	currentAdoption, err := GetCurrentAdoption(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to getCurrentAdoption")
	}
	t0, err := r.isT0Online(ctx, ms[0].IDT0, now)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to check if t0 is online for idT0:%v", ms[0].IDT0)
	}
	var extraBonus float64
	if !ms[0].ExtraBonusStartedAt.IsNil() && ms[0].ExtraBonusStartedAt.Add(r.cfg.ExtraBonuses.Duration).After(*now.Time) {
		extraBonus = ms[0].ExtraBonus
	}
	negativeMiningRate := ms[0].SlashingRateSolo + ms[0].SlashingRateT0 + ms[0].SlashingRateT1 + ms[0].SlashingRateT2
	var availableExtraBonusVal float64
	if avb, gErr := r.getAvailableExtraBonus(now, id, ms[0].ExtraBonusStartedAtField, ms[0].ExtraBonusLastClaimAvailableAtField, ms[0].MiningSessionSoloStartedAtField, ms[0].MiningSessionSoloEndedAtField, ms[0].ExtraBonusDaysClaimNotAvailableResettableField, ms[0].UTCOffsetField, ms[0].NewsSeenField); gErr == nil { //nolint:lll // .
		availableExtraBonusVal = avb.ExtraBonus
	}

	return &MiningSummary{
		MiningStreak:                r.calculateMiningStreak(now, ms[0].MiningSessionSoloStartedAt, ms[0].MiningSessionSoloEndedAt),
		MiningSession:               r.calculateMiningSession(now, ms[0].MiningSessionSoloLastStartedAt, ms[0].MiningSessionSoloEndedAt),
		RemainingFreeMiningSessions: r.calculateRemainingFreeMiningSessions(now, ms[0].MiningSessionSoloEndedAt),
		MiningRates:                 r.calculateMiningRateSummaries(t0, extraBonus, ms[0].PreStakingAllocation, ms[0].PreStakingBonus, ms[0].ActiveT1Referrals, ms[0].ActiveT2Referrals, currentAdoption.BaseMiningRate, negativeMiningRate, ms[0].BalanceTotalStandard+ms[0].BalanceTotalPreStaking, now, ms[0].MiningSessionSoloEndedAt), //nolint:lll // .
		ExtraBonusSummary:           ExtraBonusSummary{AvailableExtraBonus: availableExtraBonusVal},
	}, nil
}

func (r *repository) isT0Online(ctx context.Context, idT0 int64, now *time.Time) (uint16, error) {
	if idT0 == 0 {
		return 0, nil
	}
	if idT0 < 0 {
		idT0 *= -1
	}
	t0Ref, err := storage.Get[struct {
		model.MiningSessionSoloEndedAtField
	}](ctx, r.db, model.SerializedUsersKey(idT0))
	if err == nil && len(t0Ref) == 1 && !t0Ref[0].MiningSessionSoloEndedAt.IsNil() && t0Ref[0].MiningSessionSoloEndedAt.After(*now.Time) {
		return 1, nil
	}

	return 0, errors.Wrapf(err, "failed to get MiningSessionSoloEndedAtField for idT0:%v", idT0)
}

func (r *repository) calculateMiningSession(now, start, end *time.Time) (ms *MiningSession) {
	if ms = CalculateMiningSession(now, start, end, r.cfg.MiningSessionDuration.Max); ms != nil {
		ms.ResettableStartingAt = time.New(ms.StartedAt.Add(r.cfg.MiningSessionDuration.Min))
		ms.WarnAboutExpirationStartingAt = time.New(ms.StartedAt.Add(r.cfg.MiningSessionDuration.WarnAboutExpirationAfter))
	}

	return ms
}

func CalculateMiningSession(now, start, end *time.Time, miningSessionDuration stdlibtime.Duration) (ms *MiningSession) {
	if start.IsNil() || end.IsNil() || end.Before(*now.Time) {
		return nil
	}
	lastMiningStartedAt := time.New(start.Add((now.Sub(*start.Time) / miningSessionDuration) * miningSessionDuration))
	free := start.Add(miningSessionDuration).Before(*now.Time)

	return &MiningSession{
		StartedAt: lastMiningStartedAt,
		EndedAt:   time.New(lastMiningStartedAt.Add(miningSessionDuration)),
		Free:      &free,
	}
}

//nolint:funlen,gomnd,lll // A lot of calculations.
func (r *repository) calculateMiningRateSummaries(
	t0 uint16, extraBonus, preStakingAllocation, preStakingBonus float64,
	t1, t2 int32,
	baseMiningRate, negativeMiningRate, totalBalance float64,
	now, miningSessionSoloEndedAt *time.Time,
) (miningRates *MiningRates[*MiningRateSummary[string]]) {
	miningRates = new(MiningRates[*MiningRateSummary[string]])
	var (
		standardMiningRate         float64
		preStakingMiningRate       float64
		totalNoPreStakingBonusRate float64

		totalBonusVal                     float64
		totalNoPreStakingBonusVal         float64
		positiveTotalNoPreStakingBonusVal float64
		preStakingBonusVal                float64
	)
	if t1 < 0 {
		t1 = 0
	}
	if t2 < 0 {
		t2 = 0
	}
	positiveTotalNoPreStakingBonus := r.calculateMintedStandardCoins(t0, extraBonus, 0, uint32(t1), uint32(t2), baseMiningRate, r.cfg.GlobalAggregationInterval.Child, false)
	if positiveTotalNoPreStakingBonus > baseMiningRate {
		positiveTotalNoPreStakingBonusVal = ((positiveTotalNoPreStakingBonus - baseMiningRate) * 100) / baseMiningRate
	}
	miningRates.PositiveTotalNoPreStakingBonus = &MiningRateSummary[string]{
		Amount: fmt.Sprintf(floatToStringFormatter, roundFloat64(positiveTotalNoPreStakingBonus)),
		Bonuses: &MiningRateBonuses{
			T1:         float64(t0*r.cfg.ReferralBonusMiningRates.T0) + float64(uint32(t1)*r.cfg.ReferralBonusMiningRates.T1),
			T2:         float64(uint32(t2) * r.cfg.ReferralBonusMiningRates.T2),
			Extra:      float64(extraBonus),
			PreStaking: 0,
			Total:      positiveTotalNoPreStakingBonusVal,
		},
	}
	if miningSessionSoloEndedAt.IsNil() { //nolint:gocritic,nestif // Wrong.
		miningRates.Type = NoneMiningRateType
	} else if miningSessionSoloEndedAt.After(*now.Time) {
		miningRates.Type = PositiveMiningRateType
	} else if totalBalance <= 0.0 {
		miningRates.Type = NoneMiningRateType
	} else {
		extraBonus, t0, t1, t2 = 0, 0, 0, 0
		miningRates.Type = NegativeMiningRateType
	}
	miningRates.Base = &MiningRateSummary[string]{
		Amount: fmt.Sprintf(floatToStringFormatter, baseMiningRate),
	}
	if preStakingAllocation != 100 {
		var localTotalBonus float64
		switch miningRates.Type {
		case PositiveMiningRateType:
			standardMiningRate = r.calculateMintedStandardCoins(t0, extraBonus, preStakingAllocation, uint32(t1), uint32(t2), baseMiningRate, r.cfg.GlobalAggregationInterval.Child, false)
			if standardMiningRate > baseMiningRate {
				localTotalBonus = ((standardMiningRate - baseMiningRate) * 100) / baseMiningRate
			}
		case NegativeMiningRateType:
			standardMiningRate = (negativeMiningRate * (100 - preStakingAllocation)) / 100
		case NoneMiningRateType:
		}
		miningRates.Standard = &MiningRateSummary[string]{
			Amount: fmt.Sprintf(floatToStringFormatter, roundFloat64(standardMiningRate)),
			Bonuses: &MiningRateBonuses{
				T1:    ((float64(t0*r.cfg.ReferralBonusMiningRates.T0) + float64(uint32(t1)*r.cfg.ReferralBonusMiningRates.T1)) * (100 - preStakingAllocation)) / 100,
				T2:    float64(uint32(t2)*r.cfg.ReferralBonusMiningRates.T2) * (100 - preStakingAllocation) / 100,
				Extra: extraBonus * (100 - preStakingAllocation) / 100,
				Total: localTotalBonus,
			},
		}
	}
	if preStakingAllocation != 0 {
		var localTotalBonus float64
		switch miningRates.Type {
		case PositiveMiningRateType:
			preStakingMiningRate = r.calculateMintedPreStakingCoins(t0, extraBonus, preStakingAllocation, preStakingBonus, uint32(t1), uint32(t2), baseMiningRate, r.cfg.GlobalAggregationInterval.Child, false)
			if preStakingMiningRate > baseMiningRate {
				localTotalBonus = ((preStakingMiningRate - baseMiningRate) * 100) / baseMiningRate
			}
		case NegativeMiningRateType:
			preStakingMiningRate = (negativeMiningRate * (preStakingBonus + 100) * preStakingAllocation) / (100 * 100)
		case NoneMiningRateType:
		}
		t1Bonus := float64((uint64(t0*r.cfg.ReferralBonusMiningRates.T0)+uint64(uint32(t1)*r.cfg.ReferralBonusMiningRates.T1))*uint64(preStakingAllocation)) / 100
		t2Bonus := float64(uint32(t2)*r.cfg.ReferralBonusMiningRates.T2*uint32(preStakingAllocation)) / 100
		extraBonusVal := extraBonus * preStakingAllocation / 100
		preStakingBonusVal = (preStakingAllocation * preStakingBonus) / 100
		miningRates.PreStaking = &MiningRateSummary[string]{
			Amount: fmt.Sprintf(floatToStringFormatter, roundFloat64(preStakingMiningRate)),
			Bonuses: &MiningRateBonuses{
				T1:         t1Bonus,
				T2:         t2Bonus,
				Extra:      extraBonusVal,
				PreStaking: preStakingBonusVal,
				Total:      localTotalBonus,
			},
		}
	}
	switch miningRates.Type {
	case PositiveMiningRateType:
		if standardMiningRate+preStakingMiningRate > baseMiningRate {
			totalBonusVal = ((standardMiningRate + preStakingMiningRate - baseMiningRate) * 100) / baseMiningRate
		}
		totalNoPreStakingBonusRate = positiveTotalNoPreStakingBonus
		totalNoPreStakingBonusVal = positiveTotalNoPreStakingBonusVal
	case NegativeMiningRateType:
		totalNoPreStakingBonusRate = negativeMiningRate
	case NoneMiningRateType:
	}
	miningRates.Total = &MiningRateSummary[string]{
		Amount: fmt.Sprintf(floatToStringFormatter, roundFloat64(standardMiningRate+preStakingMiningRate)),
		Bonuses: &MiningRateBonuses{
			T1:         float64(t0*r.cfg.ReferralBonusMiningRates.T0) + float64(uint32(t1)*r.cfg.ReferralBonusMiningRates.T1),
			T2:         float64(uint32(t2) * r.cfg.ReferralBonusMiningRates.T2),
			Extra:      extraBonus,
			PreStaking: preStakingBonusVal,
			Total:      totalBonusVal,
		},
	}
	miningRates.TotalNoPreStakingBonus = &MiningRateSummary[string]{
		Amount: fmt.Sprintf(floatToStringFormatter, roundFloat64(totalNoPreStakingBonusRate)),
		Bonuses: &MiningRateBonuses{
			T1:         float64(t0*r.cfg.ReferralBonusMiningRates.T0) + float64(uint32(t1)*r.cfg.ReferralBonusMiningRates.T1),
			T2:         float64(uint32(t2) * r.cfg.ReferralBonusMiningRates.T2),
			Extra:      extraBonus,
			PreStaking: 0,
			Total:      totalNoPreStakingBonusVal,
		},
	}

	return miningRates
}

func roundFloat64(val float64) float64 {
	const precision = 100

	return math.Round(val*precision) / precision
}

func (r *repository) calculateMintedStandardCoins(
	t0 uint16, extraBonus, preStakingAllocation float64,
	t1, t2 uint32,
	baseMiningRate float64,
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
		return (baseMiningRate * mintedBase * (100 - preStakingAllocation)) / (100 * 100)
	}

	return (baseMiningRate * float64(elapsedNanos.Nanoseconds()) * mintedBase * (100 - preStakingAllocation)) /
		(float64(r.cfg.GlobalAggregationInterval.Child.Nanoseconds()) * 100 * 100)
}

func (r *repository) calculateMintedPreStakingCoins(
	t0 uint16, extraBonus, preStakingAllocation, preStakingBonus float64,
	t1, t2 uint32,
	baseMiningRate float64,
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
		return (baseMiningRate * mintedBase * ((preStakingBonus + 100) * preStakingAllocation)) /
			float64(100*100*100)
	}

	return (baseMiningRate * mintedBase * float64(elapsedNanos.Nanoseconds()) * ((preStakingBonus + 100) * preStakingAllocation)) /
		float64(uint64(r.cfg.GlobalAggregationInterval.Child.Nanoseconds())*100*100*100)
}

func (r *repository) calculateMiningStreak(now, start, end *time.Time) uint64 {
	if start.IsNil() || end.IsNil() || now.After(*end.Time) || now.Before(*start.Time) {
		return 0
	}

	return uint64(now.Sub(*start.Time) / r.cfg.MiningSessionDuration.Max)
}

func (r *repository) calculateRemainingFreeMiningSessions(now, end *time.Time) uint64 {
	if end.IsNil() || now.After(*end.Time) {
		return 0
	}

	return uint64(end.Sub(*now.Time) / r.cfg.MiningSessionDuration.Max)
}
