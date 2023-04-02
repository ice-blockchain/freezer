// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	stdlibtime "time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/wintr/coin"
	"github.com/ice-blockchain/wintr/time"
)

func TestRepository_CalculateDegradation(t *testing.T) { //nolint:funlen // .
	t.Parallel()
	rep := &repository{cfg: &config{
		RollbackNegativeMining: struct {
			Available struct {
				After stdlibtime.Duration `yaml:"after"`
				Until stdlibtime.Duration `yaml:"until"`
			} `yaml:"available"`
			LastXMiningSessionsCollectingInterval stdlibtime.Duration `yaml:"lastXMiningSessionsCollectingInterval" mapstructure:"lastXMiningSessionsCollectingInterval"` //nolint:lll // .
			AggressiveDegradationStartsAfter      stdlibtime.Duration `yaml:"aggressiveDegradationStartsAfter"`
		}{
			Available: struct {
				After stdlibtime.Duration `yaml:"after"`
				Until stdlibtime.Duration `yaml:"until"`
			}{
				Until: 180 * stdlibtime.Minute,
			},
			AggressiveDegradationStartsAfter: 90 * stdlibtime.Minute,
		},
		MiningSessionDuration: struct {
			Min                      stdlibtime.Duration `yaml:"min"`
			Max                      stdlibtime.Duration `yaml:"max"`
			WarnAboutExpirationAfter stdlibtime.Duration `yaml:"warnAboutExpirationAfter"`
		}{
			Min: 90 * stdlibtime.Second,
			Max: 180 * stdlibtime.Second,
		},
	}}
	now := time.Now()
	const initialAmount = 1_000_000_000_000
	pace := 10 * stdlibtime.Millisecond
	reference := coin.NewAmountUint64(initialAmount)
	total := coin.NewAmountUint64(initialAmount)
	expectedIterations := (rep.cfg.RollbackNegativeMining.Available.Until - rep.cfg.RollbackNegativeMining.AggressiveDegradationStartsAfter) / pace
	expectedSlashedChunk := initialAmount / expectedIterations.Nanoseconds()

	iterations := 0
	for !total.IsZero() {
		now = time.New(now.Add(pace))
		slashedAmount := rep.calculateDegradation(pace, reference, rand.Intn(2) == 1) //nolint:gosec // .
		require.InDelta(t, slashedAmount.Uint64(), expectedSlashedChunk, float64(expectedSlashedChunk/100))
		total = total.Subtract(slashedAmount)
		iterations++
	}
	require.InDelta(t, expectedIterations.Nanoseconds(), iterations, float64(expectedIterations.Nanoseconds()/100))
	require.Less(t, iterations, int(expectedIterations.Nanoseconds()))
}

func TestRepository_CalculateMiningRateSummaries(t *testing.T) { //nolint:funlen,maintidx // .
	t.Parallel()
	rep := &repository{cfg: &config{
		GlobalAggregationInterval: struct {
			Parent stdlibtime.Duration `yaml:"parent"`
			Child  stdlibtime.Duration `yaml:"child"`
		}{
			Parent: 24 * stdlibtime.Hour,
			Child:  stdlibtime.Hour,
		},
		ReferralBonusMiningRates: struct {
			T0 uint64 `yaml:"t0"`
			T1 uint64 `yaml:"t1"`
			T2 uint64 `yaml:"t2"`
		}{
			T0: 25,
			T1: 25,
			T2: 5,
		},
	}}
	baseMiningRate := coin.UnsafeParseAmount("16000000000")
	actual := rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      100,
		PreStakingAllocation: 50,
	}, nil, PositiveMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Intended.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("42000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         137,
				T2:         25,
				PreStaking: 0,
				Extra:      50,
				Total:      162,
			},
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         137,
				T2:         25,
				PreStaking: 50,
				Extra:      50,
				Total:      425,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("126000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      687,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      500,
		PreStakingAllocation: 10,
	}, nil, PositiveMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Intended.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("75600000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         247,
				T2:         45,
				PreStaking: 0,
				Extra:      90,
				Total:      372,
			},
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("50400000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         27,
				T2:         5,
				PreStaking: 50,
				Extra:      10,
				Total:      215,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("126000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      687,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      100,
		PreStakingAllocation: 100,
	}, nil, PositiveMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Wrong.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("168000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 100,
				Extra:      100,
				Total:      950,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("168000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 100,
				Extra:      100,
				Total:      950,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:         1,
		T1:         10,
		T2:         10,
		ExtraBonus: 100,
	}, nil, PositiveMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Wrong.
		Type: PositiveMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      500,
		PreStakingAllocation: 10,
	}, coin.UnsafeParseAmount("1111"), NoneMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Wrong.
		Type: NoneMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("0").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         247,
				T2:         45,
				PreStaking: 0,
				Extra:      90,
				Total:      0,
			},
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("0").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         27,
				T2:         5,
				PreStaking: 50,
				Extra:      10,
				Total:      0,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("0").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      0,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("0").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      0,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
	actual = rep.calculateMiningRateSummaries(baseMiningRate, &userMiningRateRecalculationParameters{
		T0:                   1,
		T1:                   10,
		T2:                   10,
		ExtraBonus:           100,
		PreStakingBonus:      500,
		PreStakingAllocation: 10,
	}, coin.UnsafeParseAmount("1000"), NegativeMiningRateType)
	assert.EqualValues(t, &MiningRates[MiningRateSummary[coin.ICE]]{ //nolint:dupl // Wrong.
		Type: NegativeMiningRateType,
		Base: &MiningRateSummary[coin.ICE]{
			Amount: baseMiningRate.UnsafeICE(),
		},
		Standard: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("900").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         247,
				T2:         45,
				PreStaking: 0,
				Extra:      90,
				Total:      0,
			},
		},
		PreStaking: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("600").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         27,
				T2:         5,
				PreStaking: 50,
				Extra:      10,
				Total:      0,
			},
		},
		Total: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("1500").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 50,
				Extra:      100,
				Total:      0,
			},
		},
		TotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("1000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      0,
			},
		},
		PositiveTotalNoPreStakingBonus: &MiningRateSummary[coin.ICE]{
			Amount: coin.UnsafeParseAmount("84000000000").UnsafeICE(),
			Bonuses: &MiningRateBonuses{
				T1:         275,
				T2:         50,
				PreStaking: 0,
				Extra:      100,
				Total:      425,
			},
		},
	}, actual)
}

//nolint:funlen,gocognit,revive,tparallel // A lot of assertions.
func TestRepository_CalculateCoins(t *testing.T) {
	t.Parallel()
	const someElapsedNanos = 7_777_777
	repo := &repository{
		cfg: &config{
			GlobalAggregationInterval: struct {
				Parent stdlibtime.Duration `yaml:"parent"`
				Child  stdlibtime.Duration `yaml:"child"`
			}{
				Parent: 24 * stdlibtime.Hour,
				Child:  stdlibtime.Hour,
			},
			ReferralBonusMiningRates: struct {
				T0 uint64 `yaml:"t0"`
				T1 uint64 `yaml:"t1"`
				T2 uint64 `yaml:"t2"`
			}{
				T0: 25,
				T1: 25,
				T2: 5,
			},
		},
	}
	baseMiningRate := coin.UnsafeParseAmount("16000000000")
	for _, test := range calculateCoinsTests() { //nolint:paralleltest // It's not working.
		t.Run(test.name, func(t *testing.T) {
			interval := repo.cfg.GlobalAggregationInterval.Child
			if test.expectedStandardMiningRate != "" {
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStandardMiningRate), repo.calculateMintedStandardCoins(baseMiningRate, test.userMiningRateRecalculationParameters, interval, false)) //nolint:lll // .
			} else {
				assert.Nil(t, repo.calculateMintedStandardCoins(baseMiningRate, test.userMiningRateRecalculationParameters, interval, false))
			}
			if test.expectedStandardMintedCoins != "" {
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStandardMintedCoins), repo.calculateMintedStandardCoins(baseMiningRate, test.userMiningRateRecalculationParameters, someElapsedNanos, false)) //nolint:lll // .
			} else {
				assert.Nil(t, repo.calculateMintedStandardCoins(baseMiningRate, test.userMiningRateRecalculationParameters, someElapsedNanos, false))
			}
			if test.expectedStandardT0ReferralEarnings != "" {
				params := *test.userMiningRateRecalculationParameters
				params.T1 = 0
				params.T2 = 0
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStandardT0ReferralEarnings), repo.calculateMintedStandardCoins(baseMiningRate, &params, someElapsedNanos, true)) //nolint:lll // .
			} else {
				params := *test.userMiningRateRecalculationParameters
				params.T1 = 0
				params.T2 = 0
				assert.Nil(t, repo.calculateMintedStandardCoins(baseMiningRate, &params, someElapsedNanos, true))
			}
			if test.expectedStandardT1ReferralEarnings != "" {
				params := *test.userMiningRateRecalculationParameters
				params.T0 = 0
				params.T2 = 0
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStandardT1ReferralEarnings), repo.calculateMintedStandardCoins(baseMiningRate, &params, someElapsedNanos, true)) //nolint:lll // .
			} else {
				params := *test.userMiningRateRecalculationParameters
				params.T0 = 0
				params.T2 = 0
				assert.Nil(t, repo.calculateMintedStandardCoins(baseMiningRate, &params, someElapsedNanos, true))
			}
			if test.expectedStandardT2ReferralEarnings != "" {
				params := *test.userMiningRateRecalculationParameters
				params.T1 = 0
				params.T0 = 0
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStandardT2ReferralEarnings), repo.calculateMintedStandardCoins(baseMiningRate, &params, someElapsedNanos, true)) //nolint:lll // .
			} else {
				params := *test.userMiningRateRecalculationParameters
				params.T1 = 0
				params.T0 = 0
				assert.Nil(t, repo.calculateMintedStandardCoins(baseMiningRate, &params, someElapsedNanos, true))
			}
			if test.expectedStakingMiningRate != "" {
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStakingMiningRate), repo.calculateMintedPreStakingCoins(baseMiningRate, test.userMiningRateRecalculationParameters, interval, false)) //nolint:lll // .
			} else {
				assert.Nil(t, repo.calculateMintedPreStakingCoins(baseMiningRate, test.userMiningRateRecalculationParameters, interval, false))
			}
			if test.expectedStakingMintedCoins != "" {
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStakingMintedCoins), repo.calculateMintedPreStakingCoins(baseMiningRate, test.userMiningRateRecalculationParameters, someElapsedNanos, false)) //nolint:lll // .
			} else {
				assert.Nil(t, repo.calculateMintedPreStakingCoins(baseMiningRate, test.userMiningRateRecalculationParameters, someElapsedNanos, false))
			}
			if test.expectedStakingT0ReferralEarnings != "" {
				params := *test.userMiningRateRecalculationParameters
				params.T1 = 0
				params.T2 = 0
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStakingT0ReferralEarnings), repo.calculateMintedPreStakingCoins(baseMiningRate, &params, someElapsedNanos, true)) //nolint:lll // .
			} else {
				params := *test.userMiningRateRecalculationParameters
				params.T1 = 0
				params.T2 = 0
				assert.Nil(t, repo.calculateMintedPreStakingCoins(baseMiningRate, &params, someElapsedNanos, true))
			}
			if test.expectedStakingT1ReferralEarnings != "" {
				params := *test.userMiningRateRecalculationParameters
				params.T0 = 0
				params.T2 = 0
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStakingT1ReferralEarnings), repo.calculateMintedPreStakingCoins(baseMiningRate, &params, someElapsedNanos, true)) //nolint:lll // .
			} else {
				params := *test.userMiningRateRecalculationParameters
				params.T0 = 0
				params.T2 = 0
				assert.Nil(t, repo.calculateMintedPreStakingCoins(baseMiningRate, &params, someElapsedNanos, true))
			}
			if test.expectedStakingT2ReferralEarnings != "" {
				params := *test.userMiningRateRecalculationParameters
				params.T1 = 0
				params.T0 = 0
				assert.EqualValues(t, coin.UnsafeParseAmount(test.expectedStakingT2ReferralEarnings), repo.calculateMintedPreStakingCoins(baseMiningRate, &params, someElapsedNanos, true)) //nolint:lll // .
			} else {
				params := *test.userMiningRateRecalculationParameters
				params.T1 = 0
				params.T0 = 0
				assert.Nil(t, repo.calculateMintedPreStakingCoins(baseMiningRate, &params, someElapsedNanos, true))
			}
		})
	}
}

//nolint:funlen // A lot of cases.
func calculateCoinsTests() []*struct {
	*userMiningRateRecalculationParameters
	name                               string
	expectedStandardMiningRate         string
	expectedStandardMintedCoins        string
	expectedStandardT0ReferralEarnings string
	expectedStandardT1ReferralEarnings string
	expectedStandardT2ReferralEarnings string
	expectedStakingMiningRate          string
	expectedStakingMintedCoins         string
	expectedStakingT0ReferralEarnings  string
	expectedStakingT1ReferralEarnings  string
	expectedStakingT2ReferralEarnings  string
} {
	return []*struct {
		*userMiningRateRecalculationParameters
		name                               string
		expectedStandardMiningRate         string
		expectedStandardMintedCoins        string
		expectedStandardT0ReferralEarnings string
		expectedStandardT1ReferralEarnings string
		expectedStandardT2ReferralEarnings string
		expectedStakingMiningRate          string
		expectedStakingMintedCoins         string
		expectedStakingT0ReferralEarnings  string
		expectedStakingT1ReferralEarnings  string
		expectedStakingT2ReferralEarnings  string
	}{
		// NO staking.
		{
			name:                        "no referrals active, no staking",
			expectedStandardMiningRate:  "16000000000",
			expectedStandardMintedCoins: "34567",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				PreStakingBonus: 100, // Has no effect cuz there's 0 allocation.
			},
		},
		{
			name:                               "1 T0 referrals active, no staking",
			expectedStandardMiningRate:         "20000000000",
			expectedStandardMintedCoins:        "43209",
			expectedStandardT0ReferralEarnings: "8641",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T0:              1,
				PreStakingBonus: 100, // Has no effect cuz there's 0 allocation.
			},
		},
		{
			name:                               "1 T1 referrals active, no staking",
			expectedStandardMiningRate:         "20000000000",
			expectedStandardMintedCoins:        "43209",
			expectedStandardT1ReferralEarnings: "8641",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T1:              1,
				PreStakingBonus: 100, // Has no effect cuz there's 0 allocation.
			},
		},
		{
			name:                               "1 T2 referrals active, no staking",
			expectedStandardMiningRate:         "16800000000",
			expectedStandardMintedCoins:        "36296",
			expectedStandardT2ReferralEarnings: "1728",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T2:              1,
				PreStakingBonus: 100, // Has no effect cuz there's 0 allocation.
			},
		},
		{
			name:                               "a lot of T0,T1 & T2 referrals active, no staking",
			expectedStandardMiningRate:         "840000020000000000",
			expectedStandardMintedCoins:        "1814814676543",
			expectedStandardT0ReferralEarnings: "8641",
			expectedStandardT1ReferralEarnings: "86419744444",
			expectedStandardT2ReferralEarnings: "1728394888888",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T0:              1,
				T1:              10_000_000,
				T2:              1_000_000_000,
				PreStakingBonus: 100, // Has no effect cuz there's 0 allocation.
			},
		},
		// 100% Staking.
		{
			name:                       "no referrals active, 100% staking",
			expectedStakingMiningRate:  "32000000000",
			expectedStakingMintedCoins: "69135",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				PreStakingBonus:      100,
				PreStakingAllocation: 100,
			},
		},
		{
			name:                              "1 T0 referrals active, 100% staking",
			expectedStakingMiningRate:         "40000000000",
			expectedStakingMintedCoins:        "86419",
			expectedStakingT0ReferralEarnings: "17283",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T0:                   1,
				PreStakingBonus:      100,
				PreStakingAllocation: 100,
			},
		},
		{
			name:                              "1 T1 referrals active, 100% staking",
			expectedStakingMiningRate:         "40000000000",
			expectedStakingMintedCoins:        "86419",
			expectedStakingT1ReferralEarnings: "17283",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T1:                   1,
				PreStakingBonus:      100,
				PreStakingAllocation: 100,
			},
		},
		{
			name:                              "1 T2 referrals active, 100% staking",
			expectedStakingMiningRate:         "33600000000",
			expectedStakingMintedCoins:        "72592",
			expectedStakingT2ReferralEarnings: "3456",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T2:                   1,
				PreStakingBonus:      100,
				PreStakingAllocation: 100,
			},
		},
		{
			name:                              "a lot of T0,T1 & T2 referrals active, 100% staking",
			expectedStakingMiningRate:         "1680000040000000000",
			expectedStakingMintedCoins:        "3629629353086",
			expectedStakingT0ReferralEarnings: "17283",
			expectedStakingT1ReferralEarnings: "172839488888",
			expectedStakingT2ReferralEarnings: "3456789777777",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T0:                   1,
				T1:                   10_000_000,
				T2:                   1_000_000_000,
				PreStakingBonus:      100,
				PreStakingAllocation: 100,
			},
		},
		// 50% Staking.
		{
			name:                        "no referrals active, 50% staking",
			expectedStandardMiningRate:  "8000000000",
			expectedStakingMiningRate:   "16000000000",
			expectedStakingMintedCoins:  "34567",
			expectedStandardMintedCoins: "17283",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				PreStakingBonus:      100,
				PreStakingAllocation: 50,
			},
		},
		{
			name:                               "1 T0 referrals active, 50% staking",
			expectedStandardMiningRate:         "10000000000",
			expectedStakingMiningRate:          "20000000000",
			expectedStandardMintedCoins:        "21604",
			expectedStakingMintedCoins:         "43209",
			expectedStandardT0ReferralEarnings: "4320",
			expectedStakingT0ReferralEarnings:  "8641",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T0:                   1,
				PreStakingBonus:      100,
				PreStakingAllocation: 50,
			},
		},
		{
			name:                               "1 T1 referrals active, 50% staking",
			expectedStandardMiningRate:         "10000000000",
			expectedStakingMiningRate:          "20000000000",
			expectedStandardMintedCoins:        "21604",
			expectedStakingMintedCoins:         "43209",
			expectedStandardT1ReferralEarnings: "4320",
			expectedStakingT1ReferralEarnings:  "8641",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T1:                   1,
				PreStakingBonus:      100,
				PreStakingAllocation: 50,
			},
		},
		{
			name:                               "1 T2 referrals active, 50% staking",
			expectedStandardMiningRate:         "8400000000",
			expectedStakingMiningRate:          "16800000000",
			expectedStandardMintedCoins:        "18148",
			expectedStakingMintedCoins:         "36296",
			expectedStandardT2ReferralEarnings: "864",
			expectedStakingT2ReferralEarnings:  "1728",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T2:                   1,
				PreStakingBonus:      100,
				PreStakingAllocation: 50,
			},
		},
		{
			name:                               "a lot of T0,T1 & T2 referrals active, 50% staking",
			expectedStandardMiningRate:         "420000018000000000",
			expectedStakingMiningRate:          "840000036000000000",
			expectedStandardMintedCoins:        "907407355555",
			expectedStakingMintedCoins:         "1814814711111",
			expectedStandardT0ReferralEarnings: "4320",
			expectedStakingT0ReferralEarnings:  "8641",
			expectedStandardT1ReferralEarnings: "43209872222",
			expectedStakingT1ReferralEarnings:  "86419744444",
			expectedStandardT2ReferralEarnings: "864197444444",
			expectedStakingT2ReferralEarnings:  "1728394888888",
			userMiningRateRecalculationParameters: &userMiningRateRecalculationParameters{
				T0:                   1,
				T1:                   10_000_000,
				T2:                   1_000_000_000,
				ExtraBonus:           100,
				PreStakingBonus:      100,
				PreStakingAllocation: 50,
			},
		},
	}
}

//nolint:lll,funlen,maintidx // .
func TestBalanceRecalculationTriggerStreamSource_recalculateBalances(t *testing.T) {
	t.Parallel()
	source := &balanceRecalculationTriggerStreamSource{
		processor: &processor{
			repository: &repository{
				cfg: &config{
					RollbackNegativeMining: struct {
						Available struct {
							After stdlibtime.Duration `yaml:"after"`
							Until stdlibtime.Duration `yaml:"until"`
						} `yaml:"available"`
						LastXMiningSessionsCollectingInterval stdlibtime.Duration `yaml:"lastXMiningSessionsCollectingInterval" mapstructure:"lastXMiningSessionsCollectingInterval"`
						AggressiveDegradationStartsAfter      stdlibtime.Duration `yaml:"aggressiveDegradationStartsAfter"`
					}{
						Available: struct {
							After stdlibtime.Duration `yaml:"after"`
							Until stdlibtime.Duration `yaml:"until"`
						}{
							After: 7 * 24 * stdlibtime.Hour,
							Until: 60 * 24 * stdlibtime.Hour,
						},
						LastXMiningSessionsCollectingInterval: 24 * stdlibtime.Hour,
						AggressiveDegradationStartsAfter:      30 * 24 * stdlibtime.Hour,
					},
					MiningSessionDuration: struct {
						Min                      stdlibtime.Duration `yaml:"min"`
						Max                      stdlibtime.Duration `yaml:"max"`
						WarnAboutExpirationAfter stdlibtime.Duration `yaml:"warnAboutExpirationAfter"`
					}{
						Min: 12 * stdlibtime.Hour,
						Max: 24 * stdlibtime.Hour,
					},
					ReferralBonusMiningRates: struct {
						T0 uint64 `yaml:"t0"`
						T1 uint64 `yaml:"t1"`
						T2 uint64 `yaml:"t2"`
					}{
						T0: 25,
						T1: 25,
						T2: 5,
					},
					GlobalAggregationInterval: struct {
						Parent stdlibtime.Duration `yaml:"parent"`
						Child  stdlibtime.Duration `yaml:"child"`
					}{
						Parent: 24 * stdlibtime.Hour,
						Child:  stdlibtime.Hour,
					},
					WorkerCount: 10,
				},
			},
		},
	}
	now := time.Now()
	baseMiningRate := coin.NewAmountUint64(16 * uint64(coin.Denomination))
	userID, t0UserID, tMinus1UserID := uuid.NewString(), uuid.NewString(), uuid.NewString()
	details := &BalanceRecalculationDetails{
		LastNaturalMiningStartedAt:   now,
		LastMiningStartedAt:          now,
		T0LastMiningStartedAt:        now,
		TMinus1LastMiningStartedAt:   now,
		LastMiningEndedAt:            time.New(now.Add(24 * stdlibtime.Hour)),
		T0LastMiningEndedAt:          time.New(now.Add(24 * stdlibtime.Hour)),
		TMinus1LastMiningEndedAt:     time.New(now.Add(24 * stdlibtime.Hour)),
		PreviousMiningEndedAt:        nil,
		T0PreviousMiningEndedAt:      nil,
		TMinus1PreviousMiningEndedAt: nil,
		RollbackUsedAt:               nil,
		T0RollbackUsedAt:             nil,
		TMinus1RollbackUsedAt:        nil,
		BaseMiningRate:               baseMiningRate,
		UUserID:                      userID,
		T0UserID:                     t0UserID,
		TMinus1UserID:                tMinus1UserID,
		T0:                           1,
		T1:                           1,
		T2:                           1,
		ExtraBonus:                   100,
	}
	balances := []*balanceRecalculationRow{{
		BalanceRecalculationDetails: details,
	}}
	now = time.New(now.Add(1 * stdlibtime.Hour))
	balancesForReplace, balancesForDelete, processingStoppedForUserIDs, dayOffStartedEvents, userIDs, _ := source.recalculateBalances(now, 0, balances)
	expectedBalancesForReplace := []*balance{{
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(32_000_000_000),
		UserID:     userID,
		TypeDetail: source.thisDurationDegradationReferenceTypeDetail(now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(4_000_000_000),
		UserID:     userID,
		TypeDetail: balances[0].t0TypeDetail(),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(4_000_000_000),
		UserID:     userID,
		TypeDetail: source.t0ThisDurationDegradationReferenceTypeDetail(balances[0].BalanceRecalculationDetails, now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(4_000_000_000),
		UserID:     userID,
		TypeDetail: balances[0].reverseT0TypeDetail(),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(4_000_000_000),
		UserID:     userID,
		TypeDetail: source.reverseT0ThisDurationDegradationReferenceTypeDetail(balances[0].BalanceRecalculationDetails, now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(800_000_000),
		UserID:     userID,
		TypeDetail: balances[0].reverseTMinus1TypeDetail(),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(800_000_000),
		UserID:     userID,
		TypeDetail: source.reverseTMinus1ThisDurationDegradationReferenceTypeDetail(balances[0].BalanceRecalculationDetails, now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(4_000_000_000),
		UserID:     userID,
		TypeDetail: t1BalanceTypeDetail,
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(800_000_000),
		UserID:     userID,
		TypeDetail: t2BalanceTypeDetail,
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(4_000_000_000),
		UserID:     userID,
		TypeDetail: source.t1ThisDurationDegradationReferenceTypeDetail(now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(800_000_000),
		UserID:     userID,
		TypeDetail: source.t2ThisDurationDegradationReferenceTypeDetail(now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(32_000_000_000),
		UserID:     userID,
		TypeDetail: "",
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(40_800_000_000),
		UserID:     userID,
		TypeDetail: fmt.Sprintf("/%v", now.Format(source.cfg.globalAggregationIntervalChildDateFormat())),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(40_800_000_000),
		UserID:     userID,
		TypeDetail: fmt.Sprintf("@%v", now.Format(source.cfg.globalAggregationIntervalChildDateFormat())),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}}
	sort.SliceStable(balancesForReplace, func(i, j int) bool {
		return fmt.Sprint(balancesForReplace[i].Negative, balancesForReplace[i].Type, balancesForReplace[i].TypeDetail, balancesForReplace[i].UserID) < fmt.Sprint(balancesForReplace[j].Negative, balancesForReplace[j].Type, balancesForReplace[j].TypeDetail, balancesForReplace[j].UserID)
	})
	sort.SliceStable(expectedBalancesForReplace, func(i, j int) bool {
		return fmt.Sprint(expectedBalancesForReplace[i].Negative, expectedBalancesForReplace[i].Type, expectedBalancesForReplace[i].TypeDetail, expectedBalancesForReplace[i].UserID) < fmt.Sprint(expectedBalancesForReplace[j].Negative, expectedBalancesForReplace[j].Type, expectedBalancesForReplace[j].TypeDetail, expectedBalancesForReplace[j].UserID)
	})
	assert.EqualValues(t, expectedBalancesForReplace, balancesForReplace)
	assert.EqualValues(t, []*balance{}, balancesForDelete)
	assert.EqualValues(t, map[string]*time.Time{}, processingStoppedForUserIDs)
	assert.EqualValues(t, []*FreeMiningSessionStarted{}, dayOffStartedEvents)
	assert.EqualValues(t, []string{userID}, userIDs)

	balances = balances[:0]
	for i := range expectedBalancesForReplace {
		if expectedBalancesForReplace[i].TypeDetail == fmt.Sprintf("@%v", now.Format(source.cfg.globalAggregationIntervalChildDateFormat())) {
			continue
		}
		balances = append(balances, &balanceRecalculationRow{BalanceRecalculationDetails: details, B: expectedBalancesForReplace[i]})
	}
	now = time.New(now.Add(1 * stdlibtime.Hour))
	balancesForReplace, balancesForDelete, processingStoppedForUserIDs, dayOffStartedEvents, userIDs, _ = source.recalculateBalances(now, 0, balances)
	expectedBalancesForReplace = []*balance{{
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(64_000_000_000),
		UserID:     userID,
		TypeDetail: source.thisDurationDegradationReferenceTypeDetail(now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(8_000_000_000),
		UserID:     userID,
		TypeDetail: balances[0].t0TypeDetail(),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(8_000_000_000),
		UserID:     userID,
		TypeDetail: source.t0ThisDurationDegradationReferenceTypeDetail(balances[0].BalanceRecalculationDetails, now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(8_000_000_000),
		UserID:     userID,
		TypeDetail: balances[0].reverseT0TypeDetail(),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(8_000_000_000),
		UserID:     userID,
		TypeDetail: source.reverseT0ThisDurationDegradationReferenceTypeDetail(balances[0].BalanceRecalculationDetails, now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(1_600_000_000),
		UserID:     userID,
		TypeDetail: balances[0].reverseTMinus1TypeDetail(),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(1_600_000_000),
		UserID:     userID,
		TypeDetail: source.reverseTMinus1ThisDurationDegradationReferenceTypeDetail(balances[0].BalanceRecalculationDetails, now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(8_000_000_000),
		UserID:     userID,
		TypeDetail: t1BalanceTypeDetail,
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(1_600_000_000),
		UserID:     userID,
		TypeDetail: t2BalanceTypeDetail,
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(8_000_000_000),
		UserID:     userID,
		TypeDetail: source.t1ThisDurationDegradationReferenceTypeDetail(now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(1_600_000_000),
		UserID:     userID,
		TypeDetail: source.t2ThisDurationDegradationReferenceTypeDetail(now),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(64_000_000_000),
		UserID:     userID,
		TypeDetail: "",
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(40_800_000_000),
		UserID:     userID,
		TypeDetail: fmt.Sprintf("/%v", now.Format(source.cfg.globalAggregationIntervalChildDateFormat())),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(40_800_000_000),
		UserID:     userID,
		TypeDetail: fmt.Sprintf("/%v", now.Add(-1*stdlibtime.Hour).Format(source.cfg.globalAggregationIntervalChildDateFormat())),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(40_800_000_000),
		UserID:     userID,
		TypeDetail: degradationT0T1T2TotalReferenceBalanceTypeDetail,
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}, {
		UpdatedAt:  now,
		Amount:     coin.NewAmountUint64(81_600_000_000),
		UserID:     userID,
		TypeDetail: fmt.Sprintf("@%v", now.Format(source.cfg.globalAggregationIntervalChildDateFormat())),
		Type:       totalNoPreStakingBonusBalanceType,
		Negative:   false,
	}}
	sort.SliceStable(balancesForReplace, func(i, j int) bool {
		return fmt.Sprint(balancesForReplace[i].Negative, balancesForReplace[i].Type, balancesForReplace[i].TypeDetail, balancesForReplace[i].UserID) < fmt.Sprint(balancesForReplace[j].Negative, balancesForReplace[j].Type, balancesForReplace[j].TypeDetail, balancesForReplace[j].UserID)
	})
	sort.SliceStable(expectedBalancesForReplace, func(i, j int) bool {
		return fmt.Sprint(expectedBalancesForReplace[i].Negative, expectedBalancesForReplace[i].Type, expectedBalancesForReplace[i].TypeDetail, expectedBalancesForReplace[i].UserID) < fmt.Sprint(expectedBalancesForReplace[j].Negative, expectedBalancesForReplace[j].Type, expectedBalancesForReplace[j].TypeDetail, expectedBalancesForReplace[j].UserID)
	})
	//nolint:gocritic,godot,revive // No.
	//assert.EqualValues(t, expectedBalancesForReplace, balancesForReplace)
	assert.EqualValues(t, []*balance{}, balancesForDelete)
	assert.EqualValues(t, map[string]*time.Time{}, processingStoppedForUserIDs)
	assert.EqualValues(t, []*FreeMiningSessionStarted{}, dayOffStartedEvents)
	assert.EqualValues(t, []string{userID}, userIDs)
}
