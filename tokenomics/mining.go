// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

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
		if rank, err = r.db.ZRevRank(ctx, "top_miners", SerializedUsersKey(id)).Uint64(); err != nil {
			if errors.Is(err, redis.Nil) {
				return &RankingSummary{GlobalRank: 0}, nil
			}

			return nil, errors.Wrapf(err, "failed to ZRevRank top_miners for userID:%v", userID)
		}
		var expiration stdlibtime.Duration
		if r.cfg.MiningSessionDuration.Max == 24*stdlibtime.Hour {
			expiration = 5 * stdlibtime.Minute
		} else {
			expiration = 5 * stdlibtime.Second
		}
		if err = r.db.SetEx(ctx, fmt.Sprintf("global_rank:%v", id), rank, expiration).Err(); err != nil {
			return nil, errors.Wrapf(err, "failed to set cached global_rank for id:%v", id)
		}
	}
	if userID != requestingUserID(ctx) {
		if usr, gErr := storage.Get[struct {
			HideRanking bool `redis:"hide_ranking"`
		}](ctx, r.db, SerializedUsersKey(id)); gErr != nil || (len(usr) == 1 && usr[0].HideRanking) {
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
		ids[ix] = SerializedUsersKey(id)
	}
	topMiners, err = storage.Get[Miner](ctx, r.db, ids...)
	sort.SliceStable(topMiners, func(ii, jj int) bool { return topMiners[ii].Username < topMiners[jj].Username })
	for _, topMiner := range topMiners {
		topMiner.ProfilePictureURL = r.pictureClient.DownloadURL(topMiner.ProfilePictureURL)
	}

	return topMiners, errors.Wrapf(err, "failed to get miners for ids:%#v", ids)
}

type (
	getMiningSummary struct {
		MiningSessionSoloLastStartedAt  *time.Time `redis:"mining_session_solo_last_started_at"`
		MiningSessionSoloStartedAt      *time.Time `redis:"mining_session_solo_started_at"`
		MiningSessionSoloEndedAt        *time.Time `redis:"mining_session_solo_ended_at"`
		MiningSessionT0EndedAt          *time.Time `redis:"mining_session_t0_ended_at"`
		ExtraBonusStartedAt             *time.Time `redis:"extra_bonus_started_at"`
		ExtraBonusLastClaimAvailableAt  *time.Time `redis:"extra_bonus_last_claim_available_at"`
		ExtraBonusDaysClaimNotAvailable uint16     `redis:"extra_bonus_days_claim_not_available"`
		ExtraBonus                      uint16     `redis:"extra_bonus"`
		NewsSeen                        uint16     `redis:"news_seen"`
		PreStakingBonus                 uint16     `redis:"pre_staking_bonus"`
		PreStakingAllocation            uint16     `redis:"pre_staking_allocation"`
		UTCOffset                       int16      `redis:"utc_offset"`
		ActiveT1Referrals               uint32     `redis:"active_t1_referrals"`
		ActiveT2Referrals               uint32     `redis:"active_t2_referrals"`
		IDT0                            int64      `redis:"id_t0"`
		BalanceTotal                    float64    `redis:"balance_total"`
		SlashingRateSolo                float64    `redis:"slashing_rate_solo"`
		SlashingRateT0                  float64    `redis:"slashing_rate_t0"`
		SlashingRateT1                  float64    `redis:"slashing_rate_t1"`
		SlashingRateT2                  float64    `redis:"slashing_rate_t2"`
	}
)

//nolint:funlen // .
func (r *repository) GetMiningSummary(ctx context.Context, userID string) (*MiningSummary, error) {
	id, err := r.getOrInitInternalID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	now := time.Now()
	ms, err := storage.Get[getMiningSummary](ctx, r.db, SerializedUsersKey(id))
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
	var t0, extraBonus uint16
	if ms[0].IDT0 > 0 && !ms[0].MiningSessionT0EndedAt.IsNil() && ms[0].MiningSessionT0EndedAt.After(*now.Time) {
		t0 = 1
	}
	if !ms[0].ExtraBonusStartedAt.IsNil() && ms[0].ExtraBonusStartedAt.Add(r.cfg.ExtraBonuses.Duration).After(*now.Time) {
		extraBonus = ms[0].ExtraBonus
	}
	negativeMiningRate := ms[0].SlashingRateSolo + ms[0].SlashingRateT0 + ms[0].SlashingRateT1 + ms[0].SlashingRateT2

	return &MiningSummary{
		MiningStreak:                r.calculateMiningStreak(now, ms[0].MiningSessionSoloStartedAt, ms[0].MiningSessionSoloEndedAt),
		MiningSession:               r.calculateMiningSession(now, ms[0].MiningSessionSoloLastStartedAt, ms[0].MiningSessionSoloEndedAt),
		RemainingFreeMiningSessions: r.calculateRemainingFreeMiningSessions(now, ms[0].MiningSessionSoloEndedAt),
		MiningRates:                 r.calculateMiningRateSummaries(extraBonus, t0, ms[0].PreStakingAllocation, ms[0].PreStakingBonus, ms[0].ActiveT1Referrals, ms[0].ActiveT2Referrals, currentAdoption.BaseMiningRate, negativeMiningRate, ms[0].BalanceTotal, now, ms[0].MiningSessionSoloEndedAt), //nolint:lll // .
		ExtraBonusSummary: ExtraBonusSummary{
			AvailableExtraBonus: r.calculateExtraBonus(id, ms[0].NewsSeen, ms[0].ExtraBonusDaysClaimNotAvailable, ms[0].UTCOffset, now, ms[0].MiningSessionSoloStartedAt, ms[0].MiningSessionSoloEndedAt, ms[0].ExtraBonusStartedAt, ms[0].ExtraBonusLastClaimAvailableAt), //nolint:lll // .
		},
	}, nil
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
	extraBonus, t0, preStakingAllocation, preStakingBonus uint16,
	t1, t2 uint32,
	baseMiningRate, negativeMiningRate, totalBalance float64,
	now, miningSessionSoloEndedAt *time.Time,
) (miningRates *MiningRates[*MiningRateSummary[string]]) {
	miningRates = new(MiningRates[*MiningRateSummary[string]])
	var (
		standardMiningRate         float64
		preStakingMiningRate       float64
		totalNoPreStakingBonusRate float64

		totalBonusVal                     uint64
		totalNoPreStakingBonusVal         uint64
		positiveTotalNoPreStakingBonusVal uint64
		preStakingBonusVal                uint64
	)
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
		Amount: fmt.Sprint(baseMiningRate),
	}
	if preStakingAllocation != 100 {
		var localTotalBonus uint64
		switch miningRates.Type {
		case PositiveMiningRateType:
			standardMiningRate = r.calculateMintedStandardCoins(t0, extraBonus, preStakingAllocation, t1, t2, baseMiningRate, r.cfg.GlobalAggregationInterval.Child, false)
			if standardMiningRate > baseMiningRate {
				localTotalBonus = uint64(((standardMiningRate - baseMiningRate) * 100) / baseMiningRate)
			}
		case NegativeMiningRateType:
			standardMiningRate = (negativeMiningRate * float64(100-preStakingAllocation)) / 100
		case NoneMiningRateType:
		}
		miningRates.Standard = &MiningRateSummary[string]{
			Amount: fmt.Sprint(standardMiningRate),
			Bonuses: &MiningRateBonuses{
				T1:    uint64(float64((uint64(t0*r.cfg.ReferralBonusMiningRates.T0)+uint64(t1*r.cfg.ReferralBonusMiningRates.T1))*uint64(100-preStakingAllocation)) / 100),
				T2:    uint64(float64(t2*r.cfg.ReferralBonusMiningRates.T2*uint32(100-preStakingAllocation)) / 100),
				Extra: uint64(float64(extraBonus*(100-preStakingAllocation)) / 100),
				Total: localTotalBonus,
			},
		}
	}
	if preStakingAllocation != 0 {
		var localTotalBonus uint64
		switch miningRates.Type {
		case PositiveMiningRateType:
			preStakingMiningRate = r.calculateMintedPreStakingCoins(t0, extraBonus, preStakingAllocation, preStakingBonus, t1, t2, baseMiningRate, r.cfg.GlobalAggregationInterval.Child, false)
			if preStakingMiningRate > baseMiningRate {
				localTotalBonus = uint64(((preStakingMiningRate - baseMiningRate) * 100) / baseMiningRate)
			}
		case NegativeMiningRateType:
			preStakingMiningRate = (negativeMiningRate * float64(preStakingBonus+100) * float64(preStakingAllocation)) / (100 * 100)
		case NoneMiningRateType:
		}
		t1Bonus := float64((uint64(t0*r.cfg.ReferralBonusMiningRates.T0)+uint64(t1*r.cfg.ReferralBonusMiningRates.T1))*uint64(preStakingAllocation)) / 100
		t2Bonus := float64(t2*r.cfg.ReferralBonusMiningRates.T2*uint32(preStakingAllocation)) / 100
		extraBonusVal := float64(extraBonus*preStakingAllocation) / 100
		preStakingBonusVal = uint64((float64(preStakingAllocation) * float64(preStakingBonus)) / 100)
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
	positiveTotalNoPreStakingBonus := r.calculateMintedStandardCoins(t0, extraBonus, 0, t1, t2, baseMiningRate, r.cfg.GlobalAggregationInterval.Child, false)
	if positiveTotalNoPreStakingBonus > baseMiningRate {
		positiveTotalNoPreStakingBonusVal = uint64(((positiveTotalNoPreStakingBonus - baseMiningRate) * 100) / baseMiningRate)
	}
	switch miningRates.Type {
	case PositiveMiningRateType:
		if standardMiningRate+preStakingMiningRate > baseMiningRate {
			totalBonusVal = uint64(((standardMiningRate + preStakingMiningRate - baseMiningRate) * 100) / baseMiningRate)
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
			T1:         uint64(t0*r.cfg.ReferralBonusMiningRates.T0) + uint64(t1*r.cfg.ReferralBonusMiningRates.T1),
			T2:         uint64(t2 * r.cfg.ReferralBonusMiningRates.T2),
			Extra:      uint64(extraBonus),
			PreStaking: preStakingBonusVal,
			Total:      totalBonusVal,
		},
	}
	miningRates.TotalNoPreStakingBonus = &MiningRateSummary[string]{
		Amount: fmt.Sprint(totalNoPreStakingBonusRate),
		Bonuses: &MiningRateBonuses{
			T1:         uint64(t0*r.cfg.ReferralBonusMiningRates.T0) + uint64(t1*r.cfg.ReferralBonusMiningRates.T1),
			T2:         uint64(t2 * r.cfg.ReferralBonusMiningRates.T2),
			Extra:      uint64(extraBonus),
			PreStaking: 0,
			Total:      totalNoPreStakingBonusVal,
		},
	}
	miningRates.PositiveTotalNoPreStakingBonus = &MiningRateSummary[string]{
		Amount: fmt.Sprint(positiveTotalNoPreStakingBonus),
		Bonuses: &MiningRateBonuses{
			T1:         uint64(t0*r.cfg.ReferralBonusMiningRates.T0) + uint64(t1*r.cfg.ReferralBonusMiningRates.T1),
			T2:         uint64(t2 * r.cfg.ReferralBonusMiningRates.T2),
			Extra:      uint64(extraBonus),
			PreStaking: 0,
			Total:      positiveTotalNoPreStakingBonusVal,
		},
	}

	return miningRates
}

func (r *repository) calculateMintedStandardCoins(
	t0, extraBonus, preStakingAllocation uint16,
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
		return (baseMiningRate * mintedBase * float64(100-preStakingAllocation)) / (100 * 100)
	}

	return (baseMiningRate * float64(elapsedNanos.Nanoseconds()) * mintedBase * float64(100-preStakingAllocation)) /
		(float64(r.cfg.GlobalAggregationInterval.Child.Nanoseconds()) * 100 * 100)
}

func (r *repository) calculateMintedPreStakingCoins(
	t0, extraBonus, preStakingAllocation, preStakingBonus uint16,
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
