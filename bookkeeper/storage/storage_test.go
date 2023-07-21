// SPDX-License-Identifier: ice License 1.0

package storage

import (
	"context"
	"sort"
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/time"
)

func TestStorage(t *testing.T) {
	cl := MustConnect(context.Background(), "self")
	defer func() {
		if err := recover(); err != nil {
			cl.Close()
			panic(err)
		}
		cl.Close()
	}()
	require.NoError(t, cl.Ping(context.Background()))
	t1, t2 := stdlibtime.Now().UTC().Truncate(stdlibtime.Minute), stdlibtime.Now().UTC().Add(stdlibtime.Hour).Truncate(stdlibtime.Minute)
	id1, id2 := t1.UnixNano(), t2.UnixNano()
	columns, input := InsertDDL(2)

	usrs := []*model.User{
		{
			BalanceLastUpdatedAtField:                 model.BalanceLastUpdatedAtField{BalanceLastUpdatedAt: time.New(t1)},
			MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: time.Now()},
			MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: time.Now()},
			MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.Now()},
			MiningSessionSoloPreviouslyEndedAtField:   model.MiningSessionSoloPreviouslyEndedAtField{MiningSessionSoloPreviouslyEndedAt: time.Now()},
			ExtraBonusStartedAtField:                  model.ExtraBonusStartedAtField{ExtraBonusStartedAt: time.Now()},
			ResurrectSoloUsedAtField:                  model.ResurrectSoloUsedAtField{ResurrectSoloUsedAt: time.Now()},
			ResurrectT0UsedAtField:                    model.ResurrectT0UsedAtField{ResurrectT0UsedAt: time.Now()},
			ResurrectTMinus1UsedAtField:               model.ResurrectTMinus1UsedAtField{ResurrectTMinus1UsedAt: time.Now()},
			MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: time.Now()},
			ExtraBonusLastClaimAvailableAtField:       model.ExtraBonusLastClaimAvailableAtField{ExtraBonusLastClaimAvailableAt: time.Now()},
			ProfilePictureNameField:                   model.ProfilePictureNameField{ProfilePictureName: "ProfilePictureName"},
			UsernameField:                             model.UsernameField{Username: "Username"},
			MiningBlockchainAccountAddressField:       model.MiningBlockchainAccountAddressField{MiningBlockchainAccountAddress: "MiningBlockchainAccountAddress"},
			BlockchainAccountAddressField:             model.BlockchainAccountAddressField{BlockchainAccountAddress: "BlockchainAccountAddress"},
			UserIDField:                               model.UserIDField{UserID: "UserID"},
			DeserializedUsersKey:                      model.DeserializedUsersKey{ID: id1},
			BalanceTotalStandardField:                 model.BalanceTotalStandardField{BalanceTotalStandard: 1},
			BalanceTotalPreStakingField:               model.BalanceTotalPreStakingField{BalanceTotalPreStaking: 2},
			BalanceTotalMintedField:                   model.BalanceTotalMintedField{BalanceTotalMinted: 3},
			BalanceTotalSlashedField:                  model.BalanceTotalSlashedField{BalanceTotalSlashed: 4},
			BalanceSoloPendingField:                   model.BalanceSoloPendingField{BalanceSoloPending: 5},
			BalanceT1PendingField:                     model.BalanceT1PendingField{BalanceT1Pending: 6},
			BalanceT2PendingField:                     model.BalanceT2PendingField{BalanceT2Pending: 7},
			BalanceSoloPendingAppliedField:            model.BalanceSoloPendingAppliedField{BalanceSoloPendingApplied: 8},
			BalanceT1PendingAppliedField:              model.BalanceT1PendingAppliedField{BalanceT1PendingApplied: 9},
			BalanceT2PendingAppliedField:              model.BalanceT2PendingAppliedField{BalanceT2PendingApplied: 10},
			BalanceSoloField:                          model.BalanceSoloField{BalanceSolo: 11},
			BalanceT0Field:                            model.BalanceT0Field{BalanceT0: 12},
			BalanceT1Field:                            model.BalanceT1Field{BalanceT1: 13},
			BalanceT2Field:                            model.BalanceT2Field{BalanceT2: 14},
			BalanceForT0Field:                         model.BalanceForT0Field{BalanceForT0: 15},
			BalanceForTMinus1Field:                    model.BalanceForTMinus1Field{BalanceForTMinus1: 16},
			SlashingRateSoloField:                     model.SlashingRateSoloField{SlashingRateSolo: 17},
			PreStakingBonusField:                      model.PreStakingBonusField{PreStakingBonus: 27.},
			PreStakingAllocationField:                 model.PreStakingAllocationField{PreStakingAllocation: 28.},
			ExtraBonusField:                           model.ExtraBonusField{ExtraBonus: 29.},
			SlashingRateT0Field:                       model.SlashingRateT0Field{SlashingRateT0: 18},
			SlashingRateT1Field:                       model.SlashingRateT1Field{SlashingRateT1: 19},
			SlashingRateT2Field:                       model.SlashingRateT2Field{SlashingRateT2: 20},
			SlashingRateForT0Field:                    model.SlashingRateForT0Field{SlashingRateForT0: 21},
			SlashingRateForTMinus1Field:               model.SlashingRateForTMinus1Field{SlashingRateForTMinus1: 22},
			IDT0Field:                                 model.IDT0Field{IDT0: 23},
			IDTMinus1Field:                            model.IDTMinus1Field{IDTMinus1: 24},
			ActiveT1ReferralsField:                    model.ActiveT1ReferralsField{ActiveT1Referrals: 25},
			ActiveT2ReferralsField:                    model.ActiveT2ReferralsField{ActiveT2Referrals: 26},
			NewsSeenField:                             model.NewsSeenField{NewsSeen: 30},
			ExtraBonusDaysClaimNotAvailableField:      model.ExtraBonusDaysClaimNotAvailableField{ExtraBonusDaysClaimNotAvailable: 31},
			UTCOffsetField:                            model.UTCOffsetField{UTCOffset: -32},
			HideRankingField:                          model.HideRankingField{HideRanking: true},
		}, {
			DeserializedUsersKey:      model.DeserializedUsersKey{ID: id2},
			BalanceLastUpdatedAtField: model.BalanceLastUpdatedAtField{BalanceLastUpdatedAt: time.New(t1)},
			BalanceTotalMintedField:   model.BalanceTotalMintedField{BalanceTotalMinted: 33},
			BalanceTotalSlashedField:  model.BalanceTotalSlashedField{BalanceTotalSlashed: 44},
		}}
	require.NoError(t, cl.Insert(context.Background(), columns, input, usrs))

	usrs = []*model.User{
		{
			DeserializedUsersKey:      model.DeserializedUsersKey{ID: id1},
			BalanceLastUpdatedAtField: model.BalanceLastUpdatedAtField{BalanceLastUpdatedAt: time.New(t2)},
			BalanceTotalMintedField:   model.BalanceTotalMintedField{BalanceTotalMinted: 333},
			BalanceTotalSlashedField:  model.BalanceTotalSlashedField{BalanceTotalSlashed: 444},
		},
		{
			DeserializedUsersKey:      model.DeserializedUsersKey{ID: id2},
			BalanceLastUpdatedAtField: model.BalanceLastUpdatedAtField{BalanceLastUpdatedAt: time.New(t2)},
			BalanceTotalMintedField:   model.BalanceTotalMintedField{BalanceTotalMinted: 3333},
			BalanceTotalSlashedField:  model.BalanceTotalSlashedField{BalanceTotalSlashed: 4444},
		}}
	require.NoError(t, cl.Insert(context.Background(), columns, input, usrs))

	h1, err := cl.SelectBalanceHistory(context.Background(), id1, []stdlibtime.Time{t1, t2})
	require.NoError(t, err)
	sort.SliceStable(h1, func(ii, jj int) bool { return h1[ii].CreatedAt.Before(*h1[jj].CreatedAt.Time) })
	assert.EqualValues(t, []*BalanceHistory{}, h1)
	h2, err := cl.SelectBalanceHistory(context.Background(), id2, []stdlibtime.Time{t1, t2})
	require.NoError(t, err)
	sort.SliceStable(h2, func(ii, jj int) bool { return h2[ii].CreatedAt.Before(*h2[jj].CreatedAt.Time) })
	assert.EqualValues(t, []*BalanceHistory{}, h2)
}
