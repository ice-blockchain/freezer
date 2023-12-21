// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"github.com/google/uuid"
	"github.com/ice-blockchain/eskimo/users"
	coindistribution "github.com/ice-blockchain/freezer/coin-distribution"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/time"
	"github.com/stretchr/testify/require"
	"testing"
	stdlibtime "time"
)

type eligibleFlags uint8

const (
	notEligible                        = 0
	eligibleByEthAddress eligibleFlags = 1 << iota
	eligibleByKYC
	eligibleByCountry
	eligibleByStreak
	eligibleByBalance
)
const eligible = eligibleByCountry | eligibleByKYC | eligibleByBalance | eligibleByEthAddress | eligibleByBalance | eligibleByStreak

func userEligibleForDistribution(now *time.Time, u *user, distributionCfg *coindistribution.CollectorSettings, cfg config, flags eligibleFlags, balance float64) *user {
	u.UserID = uuid.NewString()
	u.SoloLastEthereumCoinDistributionProcessedAt = nil
	if flags&eligibleByBalance != 0 && distributionCfg.MinBalanceRequired != 0 {
		u.BalanceSolo = balance
		u.BalanceTotalStandard = u.BalanceSolo
		u.BalanceSoloEthereum = 0
		u.BalanceT0Ethereum = 0
		u.BalanceT1Ethereum = 0
		u.BalanceT2Ethereum = 0
	}
	if flags&eligibleByEthAddress != 0 {
		u.MiningBlockchainAccountAddress = "skip"
	}
	if flags&eligibleByKYC != 0 {
		u.KYCStepPassed = users.QuizKYCStep
		passedTimes := model.TimeSlice([]*time.Time{
			timeDelta(0), timeDelta(0), timeDelta(0), timeDelta(0),
		})
		u.KYCStepsCreatedAt = &passedTimes
		u.KYCStepsLastUpdatedAt = &passedTimes
	}
	if flags&eligibleByCountry != 0 {
		u.Country = "US"
	}
	if flags&eligibleByStreak != 0 && distributionCfg.MinMiningStreaksRequired != 0 {
		u.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 144)
		u.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 25)
	}

	return u
}

func refEligibleForDistribution(now *time.Time, r *referral, distributionCfg *coindistribution.CollectorSettings, cfg config, flags eligibleFlags) *referral {
	r.SoloLastEthereumCoinDistributionProcessedAt = nil
	if flags&eligibleByBalance != 0 && distributionCfg.MinBalanceRequired != 0 {
		delta := distributionCfg.EndDate.Truncate(cfg.EthereumDistributionFrequency.Min).Sub(now.Truncate(cfg.EthereumDistributionFrequency.Min))
		if delta <= cfg.EthereumDistributionFrequency.Max {
			r.BalanceTotalStandard = distributionCfg.MinBalanceRequired + 1
		} else {
			r.BalanceTotalStandard = (distributionCfg.MinBalanceRequired + 1) * float64(int64(delta/cfg.EthereumDistributionFrequency.Max)+1)
		}

		r.BalanceSoloEthereum = 0
		r.BalanceT0Ethereum = 0
		r.BalanceT1Ethereum = 0
		r.BalanceT2Ethereum = 0
	}
	if flags&eligibleByEthAddress != 0 {
		r.MiningBlockchainAccountAddress = "skip"
	}
	if flags&eligibleByKYC != 0 {
		r.KYCStepPassed = users.QuizKYCStep
		passedTimes := model.TimeSlice([]*time.Time{
			timeDelta(0), timeDelta(0), timeDelta(0), timeDelta(0),
		})
		r.KYCStepsCreatedAt = &passedTimes
		r.KYCStepsLastUpdatedAt = &passedTimes
	}
	if flags&eligibleByCountry != 0 {
		r.Country = "US"
	}
	if flags&eligibleByStreak != 0 && distributionCfg.MinMiningStreaksRequired != 0 {
		r.MiningSessionSoloStartedAt = timeDelta(-stdlibtime.Hour * 144)
		r.MiningSessionSoloEndedAt = timeDelta(stdlibtime.Hour * 25)
	}

	return r
}

func testCollectorConfig(minBalance float64, streak uint64) *coindistribution.CollectorSettings {
	cfg.EthereumDistributionFrequency.Min = 24 * stdlibtime.Hour
	cfg.EthereumDistributionFrequency.Max = 7 * 24 * stdlibtime.Hour
	return &coindistribution.CollectorSettings{
		DeniedCountries:          nil,
		LatestDate:               time.New(testTime.Truncate(24 * stdlibtime.Hour).Add(364 * 24 * stdlibtime.Hour)),
		StartDate:                testTime,
		EndDate:                  timeDelta(365 * 24 * stdlibtime.Hour),
		MinBalanceRequired:       minBalance,
		StartHour:                0,
		MinMiningStreaksRequired: streak,
		Enabled:                  true,
		ForcedExecution:          false,
	}
}

func Test_user_isEligibleForSelfForEthereumDistribution(t *testing.T) {
	now := testTime
	cfg.EthereumDistributionFrequency.Min = 24 * stdlibtime.Hour
	cfg.EthereumDistributionFrequency.Max = 7 * 24 * stdlibtime.Hour
	testCfg := testCollectorConfig(1, 3)
	cfg.coinDistributionCollectorSettings.Store(testCfg)
	t.Run("empty user", func(t *testing.T) {
		var u *user
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
		u = new(user)
		u.ID = 0
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})

	t.Run("not eligible at all", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, notEligible, 0)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("only eth address set", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress, 0)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("only eth address set", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress, 0)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("eth and valid country", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry, 0)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("eth, valid country and kyc passed", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC, 0)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("eth, valid country, kyc passed and user have balance for distribution", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance, 5300)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("user have everything", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak, 5300)
		require.True(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("user have everything but country denied", func(t *testing.T) {
		testCfg.DeniedCountries = map[string]struct{}{"us": struct{}{}}
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak, 5300)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
		testCfg.DeniedCountries = nil
	})
	t.Run("user have everything but did not mine recently", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak, 5300)
		u.MiningSessionSoloEndedAt = timeDelta(-1 * stdlibtime.Hour)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("all balance already distributed", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak, 5300)
		u.BalanceSoloEthereum = u.BalanceSolo
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("user had distribution recently", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak, 5300)
		u.SoloLastEthereumCoinDistributionProcessedAt = timeDelta(-1 * stdlibtime.Minute)
		require.False(t, u.isEligibleForSelfForEthereumDistribution(now))
	})
}

func Test_referral_isEligibleForSelfForEthereumDistribution(t *testing.T) {
	now := testTime
	cfg.MiningSessionDuration.Max = 24 * stdlibtime.Hour
	cfg.MiningSessionDuration.Min = 12 * stdlibtime.Hour
	cfg.EthereumDistributionFrequency.Min = 24 * stdlibtime.Hour
	cfg.EthereumDistributionFrequency.Max = 7 * 24 * stdlibtime.Hour
	testCfg := testCollectorConfig(1, 3)
	cfg.coinDistributionCollectorSettings.Store(testCfg)
	t.Run("empty referral", func(t *testing.T) {
		var r *referral
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
		r = new(referral)
		r.ID = 0
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})

	t.Run("not eligible at all", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, notEligible)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("only eth address set", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("only eth address set", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("eth and valid country", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("eth, valid country and kyc passed", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("eth, valid country, kyc passed and referral have balance for distribution", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("referral have everything", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak)
		require.True(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("referral have everything but country denied", func(t *testing.T) {
		testCfg.DeniedCountries = map[string]struct{}{"us": struct{}{}}
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
		testCfg.DeniedCountries = nil
	})
	t.Run("ref have everything but did not mine recently", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak)
		r.MiningSessionSoloEndedAt = timeDelta(-1 * stdlibtime.Hour)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("all balance already distributed", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak)
		r.BalanceSoloEthereum = r.BalanceTotalStandard
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
	t.Run("ref had distribution recently", func(t *testing.T) {
		r := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak)
		r.SoloLastEthereumCoinDistributionProcessedAt = timeDelta(-1 * stdlibtime.Minute)
		require.False(t, r.isEligibleForSelfForEthereumDistribution(now))
	})
}

func Test_user_isEligibleForT0TMinus1ForEthereumDistribution(t *testing.T) {
	now := testTime
	cfg.MiningSessionDuration.Max = 24 * stdlibtime.Hour
	cfg.MiningSessionDuration.Min = 12 * stdlibtime.Hour
	cfg.EthereumDistributionFrequency.Min = 24 * stdlibtime.Hour
	cfg.EthereumDistributionFrequency.Max = 7 * 24 * stdlibtime.Hour
	testCfg := testCollectorConfig(1, 3)
	cfg.coinDistributionCollectorSettings.Store(testCfg)
	t.Run("empty user", func(t *testing.T) {
		var u *user
		require.False(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.False(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
		u = new(user)
		u.ID = 0
		require.False(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.False(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
	})

	t.Run("not eligible at all", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, notEligible, 0)

		require.False(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.False(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
	})
	t.Run("T0 had distribution recently", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak, 28)
		u.ForT0LastEthereumCoinDistributionProcessedAt = timeDelta(-1 * stdlibtime.Minute)
		require.False(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.True(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
	})
	t.Run("TMinus1 had distribution recently", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak, 28)
		u.ForTMinus1LastEthereumCoinDistributionProcessedAt = timeDelta(-1 * stdlibtime.Minute)
		require.True(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.False(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
	})
	t.Run("user have everything and both T0, TMinus1 had no distribution recently", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC|eligibleByBalance|eligibleByStreak, 28)
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
		u.ForTMinus1LastEthereumCoinDistributionProcessedAt = nil
		require.True(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.True(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
	})
	t.Run("streak and balance are not required for referral distribution", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByEthAddress|eligibleByCountry|eligibleByKYC, 28)
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
		require.True(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.True(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
	})
	t.Run("eth address not required for T0 distribution", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByCountry|eligibleByKYC, 28)
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
		require.True(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.True(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
	})
	t.Run("valid country is required for T0 distribution", func(t *testing.T) {
		testCfg.DeniedCountries = map[string]struct{}{"us": struct{}{}}
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByCountry|eligibleByKYC, 28)
		require.False(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.False(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
		testCfg.DeniedCountries = nil
		require.True(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.True(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
	})
	t.Run("kyc is required for T0 distribution", func(t *testing.T) {
		testCfg.DeniedCountries = map[string]struct{}{"us": struct{}{}}
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligibleByCountry|eligibleByKYC, 28)
		require.False(t, u.isEligibleForT0ForEthereumDistribution(now))
		require.False(t, u.isEligibleForTMinus1ForEthereumDistribution(now))
		testCfg.DeniedCountries = nil
	})
}

func Test_processEthereumCoinDistribution(t *testing.T) {
	ethereumDistributionDryRunModeEnabled = false
	t.Run("coin distribution disabled", testProcessEthereumCoinDistributionDisabled)
	t.Run("only solo eligible", testProcessEthereumCoinDistributionSolo)
	t.Run("solo and T0 eligible", testProcessEthereumCoinDistributionT0)
	t.Run("solo, T0 and TMinus1 are eligible", testProcessEthereumCoinDistributionT0TMinus1)
}

func testProcessEthereumCoinDistributionT0TMinus1(t *testing.T) {
	testCfg := testCollectorConfig(1, 3)
	cfg.coinDistributionCollectorSettings.Store(testCfg)
	cfg.MiningSessionDuration.Max = 24 * stdlibtime.Hour
	cfg.MiningSessionDuration.Min = 12 * stdlibtime.Hour
	now := timeDelta(1 * stdlibtime.Hour)
	t.Run("TMinus1 is empty - no record for TMinus1", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 5300)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		u.IDT0 = t0.ID
		u.BalanceForT0 = 5300
		u.BalanceT0 = 10600
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, nil)
		require.EqualValues(t, 100, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 3)
		compareEarnerForUser(t, u, earners[0], now, 100, t0)
		require.EqualValues(t, 100, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 200, earners[1].Balance)
		require.EqualValues(t, 200, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, 100)
		require.EqualValues(t, 100, u.BalanceForT0Ethereum)
		tMinus1 := newRef()
		tMinus1.UserID = t0.UserID
		earners, distributionsForT0, distributionsForTMinus1 = u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 0, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 1)
	})
	t.Run("TMinus1 is not eligible for distribution", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 5300)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		tMinus1 := refEligibleForDistribution(now, newRef(), testCfg, cfg, notEligible)
		u.BalanceForT0 = 5300
		u.BalanceForTMinus1 = 5300
		u.BalanceT0 = 10600
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 100, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 100, t0)
		require.EqualValues(t, 100, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 200, earners[1].Balance)
		require.EqualValues(t, 200, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, 100)
		require.EqualValues(t, 100, u.BalanceForT0Ethereum)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		require.EqualValues(t, 0, u.BalanceForTMinus1Ethereum)
		require.Equal(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Equal(t, now, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("TMinus1 eligible for distribution, balance is distributed to TMinus1", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 5300)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		tMinus1 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		u.BalanceForT0 = 5300
		u.BalanceForTMinus1 = 5300
		u.BalanceT0 = 10600
		tMinus1.BalanceTotalStandard = 5300
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 100, distributionsForT0)
		require.EqualValues(t, 100, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 100, t0)
		require.EqualValues(t, 100, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 200, earners[1].Balance)
		require.EqualValues(t, 200, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, 100)
		require.EqualValues(t, 100, u.BalanceForT0Ethereum)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 100)
		require.EqualValues(t, 100, u.BalanceForTMinus1Ethereum)
		require.Equal(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Equal(t, now, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Equal(t, now, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("TMinus1 eligible for distribution, but had distribution recently", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 5300)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		tMinus1 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		u.ForTMinus1LastEthereumCoinDistributionProcessedAt = timeDelta(-1 * stdlibtime.Hour)
		u.BalanceForT0 = 5300
		u.BalanceForTMinus1 = 5300
		u.BalanceT0 = 10600
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 100, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 100, t0)
		require.EqualValues(t, 100, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 200, earners[1].Balance)
		require.EqualValues(t, 200, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, 100)
		require.EqualValues(t, 100, u.BalanceForT0Ethereum)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		require.EqualValues(t, 0, u.BalanceForTMinus1Ethereum)
		require.Equal(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Equal(t, now, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("TMinus1 eligible for distribution, but some balance already distributed", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 5300)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		tMinus1 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		u.BalanceSolo = 5300
		u.BalanceSoloEthereum = 2650
		u.ForTMinus1LastEthereumCoinDistributionProcessedAt = nil
		u.BalanceForT0 = 5300
		u.BalanceForT0Ethereum = 2650
		u.BalanceForTMinus1 = 10600
		u.BalanceForTMinus1Ethereum = 5300
		u.BalanceT0 = 10600
		u.BalanceTotalStandard = u.BalanceT0 + u.BalanceSolo
		u.BalanceT0Ethereum = 5300
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 50, distributionsForT0)
		require.EqualValues(t, 100, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 50, t0)
		require.EqualValues(t, 50+2650, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 100, earners[1].Balance)
		require.EqualValues(t, 100+5300, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, 50)
		require.EqualValues(t, 50+2650, u.BalanceForT0Ethereum)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 100)
		require.EqualValues(t, 100+5300, u.BalanceForTMinus1Ethereum)
		require.Equal(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Equal(t, now, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Equal(t, now, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("TMinus1 has prestaking, so it affects forTMinus1 distributed balance", func(t *testing.T) {
		testWithPrestaking(t, now, testCfg, 50, 0, 50)
	})
}

func testProcessEthereumCoinDistributionT0(t *testing.T) {
	testCfg := testCollectorConfig(1, 3)
	cfg.coinDistributionCollectorSettings.Store(testCfg)
	cfg.MiningSessionDuration.Max = 24 * stdlibtime.Hour
	cfg.MiningSessionDuration.Min = 12 * stdlibtime.Hour
	now := timeDelta(1 * stdlibtime.Hour)
	t.Run("user is not egilible but T0 is -> no distribution", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, notEligible, 0)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		u.IDT0 = t0.ID
		tMinus1 := newRef()
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 0, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 0, t0)
		require.EqualValues(t, 0, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 0, earners[1].Balance)
		require.EqualValues(t, 0, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, 0)
		require.EqualValues(t, 0, u.BalanceForT0Ethereum)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		require.EqualValues(t, 0, u.BalanceForTMinus1Ethereum)
		require.Nil(t, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("both T0 and user are eligible", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 10600)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		u.IDT0 = t0.ID
		tMinus1 := newRef()
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
		u.BalanceForT0 = 5300
		u.BalanceT0 = 10600
		forT0Expected := float64(100)
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, forT0Expected, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 200, t0)
		require.EqualValues(t, 200, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 200, earners[1].Balance)
		require.EqualValues(t, 200, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, forT0Expected)
		require.EqualValues(t, forT0Expected, u.BalanceForT0Ethereum)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		require.EqualValues(t, 0, u.BalanceForTMinus1Ethereum)
		require.Equal(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Equal(t, now, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("both T0 and user are eligible but T0 balance was distributed recently", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 10600)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		u.IDT0 = t0.ID
		tMinus1 := newRef()
		u.ForT0LastEthereumCoinDistributionProcessedAt = timeDelta(-1 * stdlibtime.Minute)
		u.BalanceForT0 = 5300
		u.BalanceT0 = 10600
		forT0Expected := float64(0)
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, forT0Expected, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 200, t0)
		require.EqualValues(t, 200, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 200, earners[1].Balance)
		require.EqualValues(t, 200, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, forT0Expected)
		require.EqualValues(t, 0, u.BalanceForT0Ethereum)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		require.EqualValues(t, 0, u.BalanceForTMinus1Ethereum)
		require.Equal(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("T0 & user are eligible but has prestaking", func(t *testing.T) {
		t.Run("user has prestaking, it affects balance earned by T0 and solo", func(t *testing.T) {
			testWithPrestaking(t, now, testCfg, 50, 0, 0)
		})
		t.Run("T0 has prestaking, it affects balance ForT0", func(t *testing.T) {
			testWithPrestaking(t, now, testCfg, 0, 50, 0)
		})
		t.Run("both have prestaking affects balance earned by T0, solo and ForT0", func(t *testing.T) {
			testWithPrestaking(t, now, testCfg, 50, 50, 0)
		})
	})

	t.Run("T0 balance is partially distributed", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 10600)
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
		u.IDT0 = t0.ID
		tMinus1 := newRef()
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
		u.BalanceForT0 = 5300
		u.BalanceT0 = 10600
		u.BalanceForT0Ethereum = 0.5 * u.BalanceForT0
		u.BalanceT0Ethereum = 0.5 * u.BalanceT0
		forT0Expected := float64(50)
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, forT0Expected, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 200, t0)
		require.EqualValues(t, 200, u.BalanceSoloEthereum)
		require.EqualValues(t, t0.UserID, earners[1].EarnerUserID)
		require.EqualValues(t, u.UserID, earners[1].UserID)
		require.EqualValues(t, now, earners[1].CreatedAt)
		require.EqualValues(t, 100, earners[1].Balance)
		require.EqualValues(t, u.BalanceSolo*0.5+100, u.BalanceT0Ethereum)
		compareEarnerForRef(t, u, t0, earners[2], now, forT0Expected)
		require.EqualValues(t, 2650+forT0Expected, u.BalanceForT0Ethereum)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		require.EqualValues(t, 0, u.BalanceForTMinus1Ethereum)
		require.Equal(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Equal(t, now, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
}

func testWithPrestaking(tb testing.TB, now *time.Time, testCfg *coindistribution.CollectorSettings, prestakingUser, prestakingT0, prestakingTMinus1 float64) {
	tb.Helper()
	u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 15900)
	t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
	u.IDT0 = t0.ID
	tMinus1 := refEligibleForDistribution(now, newRef(), testCfg, cfg, eligible)
	u.ForT0LastEthereumCoinDistributionProcessedAt = nil
	u.ForTMinus1LastEthereumCoinDistributionProcessedAt = nil
	u.BalanceForT0 = 5300
	u.BalanceT0 = 10600
	u.PreStakingBonus = 250
	if prestakingUser > 0 {
		u.PreStakingAllocation = prestakingUser
		u.PreStakingBonus = 250
	}
	if prestakingT0 > 0 {
		t0.PreStakingAllocation = prestakingT0
		t0.PreStakingBonus = 250
	}
	if prestakingTMinus1 > 0 {
		u.BalanceForTMinus1 = 5300
		tMinus1.PreStakingAllocation = prestakingTMinus1
		tMinus1.PreStakingBonus = 250
	}

	forT0Expected := float64(100) * ((100 - prestakingT0) / 100.0)
	forSoloExpected := float64(300) * ((100 - prestakingUser) / 100.0)
	forTMinus1Expected := float64(0)
	if prestakingTMinus1 > 0 {
		forTMinus1Expected = float64(100) * ((100 - prestakingTMinus1) / 100.0)
	}
	T0EarningsExpected := float64(200) * ((100 - prestakingUser) / 100.0)
	earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
	require.EqualValues(tb, forT0Expected, distributionsForT0)
	require.EqualValues(tb, forTMinus1Expected, distributionsForTMinus1)
	require.Len(tb, earners, 4)
	compareEarnerForUser(tb, u, earners[0], now, forSoloExpected, t0)
	require.EqualValues(tb, u.BalanceSoloEthereum, forSoloExpected)
	require.EqualValues(tb, t0.UserID, earners[1].EarnerUserID)
	require.EqualValues(tb, u.UserID, earners[1].UserID)
	require.EqualValues(tb, now, earners[1].CreatedAt)
	require.EqualValues(tb, T0EarningsExpected, earners[1].Balance)
	require.EqualValues(tb, u.BalanceT0Ethereum, T0EarningsExpected)
	compareEarnerForRef(tb, u, t0, earners[2], now, forT0Expected)
	require.EqualValues(tb, u.BalanceForT0Ethereum, forT0Expected)
	compareEarnerForRef(tb, u, tMinus1, earners[3], now, forTMinus1Expected)
	require.EqualValues(tb, u.BalanceForTMinus1Ethereum, forTMinus1Expected)
	require.Equal(tb, now, u.SoloLastEthereumCoinDistributionProcessedAt)
	require.Equal(tb, now, u.ForT0LastEthereumCoinDistributionProcessedAt)
	require.Equal(tb, now, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
}

func testProcessEthereumCoinDistributionSolo(t *testing.T) {
	testCfg := testCollectorConfig(1, 3)
	cfg.coinDistributionCollectorSettings.Store(testCfg)
	cfg.MiningSessionDuration.Max = 24 * stdlibtime.Hour
	cfg.MiningSessionDuration.Min = 12 * stdlibtime.Hour
	now := timeDelta(1 * stdlibtime.Hour)

	t.Run("not eligible -> balances are 0 and no refs", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, notEligible, 0)
		t0 := newRef()
		tMinus1 := newRef()
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 0, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 0, t0)
		require.EqualValues(t, 0, u.BalanceSoloEthereum)
		require.EqualValues(t, 0, earners[1].Balance)
		compareEarnerForRef(t, u, t0, earners[2], now, 0)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		require.Nil(t, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("had distribution recently -> balances are 0 and no refs", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 106)
		u.SoloLastEthereumCoinDistributionProcessedAt = timeDelta(-1 * stdlibtime.Minute)
		t0 := newRef()
		tMinus1 := newRef()
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 0, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, 0, t0)
		require.EqualValues(t, 0, u.BalanceSoloEthereum)
		require.EqualValues(t, 0, earners[1].Balance)
		compareEarnerForRef(t, u, t0, earners[2], now, 0)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		require.Nil(t, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("eligible solo -> nothing distributed before -> get distribution solo with balance", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 106)
		require.True(t, u.isEligibleForSelfForEthereumDistribution(now))
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, notEligible)
		tMinus1 := refEligibleForDistribution(now, newRef(), testCfg, cfg, notEligible)
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 0, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForRef(t, u, t0, earners[2], now, 0)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)
		expectedBalance := float64(2)
		compareEarnerForUser(t, u, earners[0], now, expectedBalance, t0)
		require.EqualValues(t, expectedBalance, u.BalanceSoloEthereum)
		require.EqualValues(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("with prestaked balance (not applied to distibution)", func(t *testing.T) {
		testWithPrestaking(t, now, testCfg, 50, 0, 0)
	})
	t.Run("eligible solo -> had distributions before -> get reduced amounts", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 10600)
		u.BalanceSoloEthereum = 0.5 * u.BalanceSolo
		expectedBalance := float64(100)
		require.True(t, u.isEligibleForSelfForEthereumDistribution(now))
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, notEligible)
		tMinus1 := refEligibleForDistribution(now, newRef(), testCfg, cfg, notEligible)
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, tMinus1)
		require.EqualValues(t, 0, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 4)
		compareEarnerForUser(t, u, earners[0], now, expectedBalance, t0)
		require.EqualValues(t, expectedBalance+5300, u.BalanceSoloEthereum)
		require.EqualValues(t, 0, earners[1].Balance)
		compareEarnerForRef(t, u, t0, earners[2], now, 0)
		compareEarnerForRef(t, u, tMinus1, earners[3], now, 0)

		require.EqualValues(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.EqualValues(t, 0, u.BalanceT1Ethereum)
		require.EqualValues(t, 0, u.BalanceForT0Ethereum)
		require.EqualValues(t, 0, u.BalanceT0Ethereum)
		require.Nil(t, u.BalanceT1EthereumPending)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
		require.EqualValues(t, 0, u.BalanceT2Ethereum)
		require.EqualValues(t, 0, u.BalanceForTMinus1Ethereum)
		require.Nil(t, u.BalanceT2EthereumPending)
	})

	t.Run("eligible solo -> no referrals set -> no referrals in earners but solo is distributed", func(t *testing.T) {
		u := userEligibleForDistribution(now, newUser(), testCfg, cfg, eligible, 106)
		require.True(t, u.isEligibleForSelfForEthereumDistribution(now))
		t0 := refEligibleForDistribution(now, newRef(), testCfg, cfg, notEligible)
		t0.UserID = u.UserID
		earners, distributionsForT0, distributionsForTMinus1 := u.processEthereumCoinDistribution(now, t0, nil)
		require.EqualValues(t, 0, distributionsForT0)
		require.EqualValues(t, 0, distributionsForTMinus1)
		require.Len(t, earners, 1)
		expectedBalance := coindistribution.CalculateEthereumDistributionICEBalance(
			u.BalanceSolo,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max,
			now, testCfg.EndDate,
		)
		compareEarnerForUser(t, u, earners[0], now, expectedBalance, t0)
		require.EqualValues(t, expectedBalance, u.BalanceSoloEthereum)
		require.EqualValues(t, now, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
}

func testProcessEthereumCoinDistributionDisabled(t *testing.T) {
	cfg.coinDistributionCollectorSettings.Store(&coindistribution.CollectorSettings{Enabled: false})
	now := testTime
	t.Run("with no pending balance", func(t *testing.T) {
		u := newUser()
		earners, t0, tMinus1 := u.processEthereumCoinDistribution(now, nil, nil)
		require.EqualValues(t, 0, t0)
		require.EqualValues(t, 0, tMinus1)
		require.Nil(t, earners)
		require.EqualValues(t, 0, u.BalanceT1Ethereum)
		require.EqualValues(t, 0, u.BalanceT2Ethereum)
		require.Nil(t, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})
	t.Run("with pending balance to be applied", func(t *testing.T) {
		const uint256max = 115792089237316195423570985008687907853269984665640564039457584007913129639935
		u := newUser()
		large := new(model.FlexibleFloat64)
		require.NoError(t, large.UnmarshalText([]byte("115792089237316195423570985008687907853269984665640564039457584007913129639935")))
		u.BalanceT1EthereumPending = large
		u.BalanceT2EthereumPending = large
		earners, t0, tMinus1 := u.processEthereumCoinDistribution(now, nil, nil)
		require.EqualValues(t, 0, t0)
		require.EqualValues(t, 0, tMinus1)
		require.Nil(t, earners)
		require.EqualValues(t, float64(uint256max), u.BalanceT1Ethereum)
		require.EqualValues(t, float64(uint256max), u.BalanceT2Ethereum)
		require.Nil(t, u.SoloLastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForT0LastEthereumCoinDistributionProcessedAt)
		require.Nil(t, u.ForTMinus1LastEthereumCoinDistributionProcessedAt)
	})

}

func compareEarnerForRef(tb testing.TB, user *user, ref *referral, earner *coindistribution.ByEarnerForReview, now *time.Time, expectedBalance float64) {
	tb.Helper()
	require.EqualValues(tb, user.UserID, earner.EarnerUserID)
	require.EqualValues(tb, ref.UserID, earner.UserID)
	require.EqualValues(tb, expectedBalance, earner.Balance)
	require.EqualValues(tb, now, earner.CreatedAt)
}
func compareEarnerForUser(tb testing.TB, u *user, earner *coindistribution.ByEarnerForReview, now *time.Time, expectedBalance float64, t0 *referral) {
	tb.Helper()
	require.EqualValues(tb, u.UserID, earner.UserID)
	require.EqualValues(tb, u.UserID, earner.EarnerUserID)
	require.EqualValues(tb, expectedBalance, earner.Balance)
	require.EqualValues(tb, now, earner.CreatedAt)
	require.EqualValues(tb, u.ID, earner.InternalID)
	require.EqualValues(tb, u.Username, earner.Username)
	if t0 != nil {
		require.EqualValues(tb, t0.Username, earner.ReferredByUsername)
	}
}
