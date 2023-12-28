// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"testing"
	stdlibtime "time"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/time"
	"github.com/stretchr/testify/assert"
)

var (
	testTime = time.New(stdlibtime.Date(2023, 1, 2, 10, 4, 5, 6, stdlibtime.UTC))
)

func TestIsEligibleForEthereumDistributionNow(t *testing.T) {
	t.Parallel()
	coinDistributionStartDate := time.New(stdlibtime.Date(2024, 1, 16, 0, 0, 0, 0, stdlibtime.UTC))
	lastCoinDistributionProcessedAt := time.New(stdlibtime.Date(2024, 1, 19, 0, 0, 0, 0, stdlibtime.UTC))

	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 22, 0, 0, 0, 0, stdlibtime.UTC))
	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 21, 0, 0, 0, 0, stdlibtime.UTC))
	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 16, 0, 0, 0, 0, stdlibtime.UTC))
	assert.True(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
}

func TestIsCoinDistributionCollectorEnabled(t *testing.T) {
	t.Parallel()

	t.Run("enabled && forced execution", func(t *testing.T) {
		t.Parallel()

		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      nil,
			StartDate:       nil,
			EndDate:         nil,
			StartHour:       0,
			Enabled:         true,
			ForcedExecution: true,
		}
		assert.Equal(t, true, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("!enabled && forced execution", func(t *testing.T) {
		t.Parallel()

		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      nil,
			StartDate:       nil,
			EndDate:         nil,
			StartHour:       5,
			Enabled:         false,
			ForcedExecution: true,
		}
		assert.Equal(t, false, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("enabled && !forced execution && now.Hour < cs.StartHour", func(t *testing.T) {
		t.Parallel()

		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      testTime,
			StartDate:       time.New(testTime.Add(-2 * stdlibtime.Hour)),
			EndDate:         time.New(testTime.Add(2 * stdlibtime.Hour)),
			Enabled:         true,
			ForcedExecution: false,
		}
		assert.Equal(t, false, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("enabled && !forced execution && now.Hour >= cs.StartHour && now.Before(start)", func(t *testing.T) {
		t.Parallel()

		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      nil,
			StartDate:       time.New(testTime.Add(1 * stdlibtime.Hour)),
			EndDate:         time.New(testTime.Add(2 * stdlibtime.Hour)),
			StartHour:       9,
			Enabled:         true,
			ForcedExecution: false,
		}
		assert.Equal(t, false, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("enabled && !forced execution && now.Hour >= cs.StartHour && now.After(end)", func(t *testing.T) {
		t.Parallel()
		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      nil,
			StartDate:       time.New(testTime.Add(-3 * stdlibtime.Hour)),
			EndDate:         time.New(testTime.Add(-1 * stdlibtime.Hour)),
			StartHour:       9,
			Enabled:         true,
			ForcedExecution: false,
		}
		assert.Equal(t, false, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("enabled && hour && start < hour < end && latest date is nil", func(t *testing.T) {
		t.Parallel()
		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      nil,
			StartDate:       time.New(testTime.Add(-1 * stdlibtime.Hour)),
			EndDate:         time.New(testTime.Add(1 * stdlibtime.Hour)),
			StartHour:       9,
			Enabled:         true,
			ForcedExecution: false,
		}
		assert.Equal(t, true, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("enabled && hour && start < hour < end && !truncate(now ~ latest date)", func(t *testing.T) {
		t.Parallel()
		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      time.New(testTime.Add(2 * stdlibtime.Hour)),
			StartDate:       time.New(testTime.Add(-3 * stdlibtime.Hour)),
			EndDate:         time.New(testTime.Add(3 * stdlibtime.Hour)),
			StartHour:       9,
			Enabled:         true,
			ForcedExecution: false,
		}

		assert.Equal(t, true, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("enabled && hour && start < hour < end && truncate(now ~ latest date)", func(t *testing.T) {
		t.Parallel()
		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      time.New(testTime.Add(30 * stdlibtime.Minute)),
			StartDate:       time.New(testTime.Add(-3 * stdlibtime.Hour)),
			EndDate:         time.New(testTime.Add(3 * stdlibtime.Hour)),
			StartHour:       9,
			Enabled:         true,
			ForcedExecution: false,
		}
		assert.Equal(t, false, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})
}

func TestCalculateEthereumDistributionICEBalance(t *testing.T) {
	t.Parallel()

	t.Run("delta < ethereumDistributionFrequencyMax", func(t *testing.T) {
		t.Parallel()
		standardBalance := 1000.0
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		now := testTime
		ethereumDistributionEndDate := time.New(now.Add(5 * stdlibtime.Hour))
		result := CalculateEthereumDistributionICEBalance(standardBalance, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax, now, ethereumDistributionEndDate)
		assert.Equal(t, 1000.0, result)
	})

	t.Run("delta == ethereumDistributionFrequencyMax", func(t *testing.T) {
		t.Parallel()
		standardBalance := 1000.0
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		now := testTime
		ethereumDistributionEndDate := time.New(now.Add(24 * stdlibtime.Hour))
		result := CalculateEthereumDistributionICEBalance(standardBalance, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax, now, ethereumDistributionEndDate)
		assert.Equal(t, 1000.0, result)
	})

	t.Run("delta > ethereumDistributionFrequencyMax", func(t *testing.T) {
		t.Parallel()
		standardBalance := 1000.0
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		now := testTime
		ethereumDistributionEndDate := time.New(now.Add(48 * stdlibtime.Hour))
		result := CalculateEthereumDistributionICEBalance(standardBalance, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax, now, ethereumDistributionEndDate)
		assert.Equal(t, 333.3333333333333, result)
	})
}

func TestIsEligibleForEthereumDistribution(t *testing.T) {
	t.Parallel()

	t.Run("Denied country", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(0)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{
			"France":  {},
			"Germany": {},
			"spain":   {},
		}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.Now(), time.New(now.Add(48*stdlibtime.Hour)), time.New(now.Add(48*stdlibtime.Hour)), time.New(now.Add(-24*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)

		country = "france"
		res = IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)

		country = "Germany"
		res = IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)

		country = "germany"
		res = IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)

		country = "Spain"
		res = IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("Positive case", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(5*stdlibtime.Minute))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, true, res)
	})

	t.Run("Wrong eth address", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "1234"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(12*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("Not enough mining streaks", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(3)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(12*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("Mining session solo ended at is nil", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(-24*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		var miningSessionSoloEndedAt *time.Time
		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("Mining session solo ended before collecting enabled", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, ethereumDistributionEndDate, miningSessionSoloEndedAt, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(96*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("Mining session solo started at is nil", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(-24*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		var miningSessionSoloStartedAt *time.Time
		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("Not enough kyc step", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(12*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.Social1KYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)

		kycState = model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.Social2KYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		res = IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)

		kycState = model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.Social3KYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		res = IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)

		kycState = model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.FacialRecognitionKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		res = IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)

		kycState = model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.LivenessDetectionKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		res = IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("KYC last updated at is nil", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(12*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: nil},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("KYC previous last updated at is nil", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(12*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), nil}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("KYC blocked < QiuzKYCStep", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour)), time.New(now.Add(12*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.FacialRecognitionKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("Balance less than minimum distribution balance required, delta == ethereumDistributionFrequencyMax", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 1000.0
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		now := testTime
		minEthereumDistributionICEBalanceRequired := 1100.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(24*stdlibtime.Hour)), time.New(now.Add(12*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})

	t.Run("Balance less than minimum distribution balance required, delta > ethereumDistributionFrequencyMax", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 1000.0
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		now := testTime
		minEthereumDistributionICEBalanceRequired := 1000.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, collectingEndedAt := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(48*stdlibtime.Hour)), time.New(now.Add(12*stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		res := IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now, collectingEndedAt,
			miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax)
		assert.Equal(t, false, res)
	})
}

func TestIsEligibleForEthereumDistributionNow(t *testing.T) {
	t.Parallel()

	t.Run("lastEthereumCoinDistributionProcessedAt.IsNil() && today.Equal(ethereumDistributionStartDate)", func(t *testing.T) {
		var lastEthereumCoinDistributionProcessedAt *time.Time
		id := int64(1)
		now := testTime
		coinDistributionStartDate := testTime
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		assert.Equal(t, true, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("!lastEthereumCoinDistributionProcessedAt.IsNil() && today.Equal(ethereumDistributionStartDate)", func(t *testing.T) {
		id := int64(1)
		now := testTime
		lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate := testTime, testTime
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour

		assert.Equal(t, false, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate)", func(t *testing.T) {
		var lastEthereumCoinDistributionProcessedAt *time.Time
		id := int64(1)
		now := testTime
		coinDistributionStartDate := time.New(testTime.Add(24 * stdlibtime.Hour))
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		assert.Equal(t, false, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("!lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate) && secondReservationIsWithinFirstDistributionCycle == false", func(t *testing.T) {
		id := int64(1)
		now := testTime
		lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate := testTime, time.New(testTime.Add(24*stdlibtime.Hour))
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		assert.Equal(t, false, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("!lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate) && secondReservationIsWithinFirstDistributionCycle == true", func(t *testing.T) {
		id := int64(1)
		now := testTime
		lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate := time.New(testTime.Add(97*stdlibtime.Hour)), time.New(testTime.Add(48*stdlibtime.Hour))
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		assert.Equal(t, true, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		id = int64(2)
		lastEthereumCoinDistributionProcessedAt = time.New(testTime.Add(98 * stdlibtime.Hour))
		assert.Equal(t, true, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		id = int64(3)
		lastEthereumCoinDistributionProcessedAt = time.New(testTime.Add(99 * stdlibtime.Hour))
		assert.Equal(t, true, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("!lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate) && reservationIsTodayAndIsOutsideOfFirstCycle == false", func(t *testing.T) {
		id := int64(1)
		now := time.New(testTime.Add(48 * stdlibtime.Hour))
		lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate := time.New(testTime.Add(97*stdlibtime.Hour)), testTime
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		assert.Equal(t, false, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		id = int64(2)
		assert.Equal(t, false, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		id = int64(3)
		assert.Equal(t, false, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("!lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate) && reservationIsTodayAndIsOutsideOfFirstCycle == true", func(t *testing.T) {
		id := int64(1)
		now := time.New(testTime.Add(48 * stdlibtime.Hour))
		lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate := time.New(testTime.Add(96*stdlibtime.Hour)), testTime
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		assert.Equal(t, true, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		id = int64(2)
		lastEthereumCoinDistributionProcessedAt = time.New(testTime.Add(120 * stdlibtime.Hour))
		assert.Equal(t, true, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		id = int64(3)
		lastEthereumCoinDistributionProcessedAt = time.New(testTime.Add(144 * stdlibtime.Hour))
		assert.Equal(t, true, IsEligibleForEthereumDistributionNow(id, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})
}

func TestIsEthereumAddressValid(t *testing.T) {
	t.Parallel()

	t.Run("Positive case", func(t *testing.T) {
		t.Parallel()
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		assert.Equal(t, true, isEthereumAddressValid(ethAddress))
	})

	t.Run("Empty eth address", func(t *testing.T) {
		t.Parallel()
		ethAddress := ""
		assert.Equal(t, false, isEthereumAddressValid(ethAddress))
	})

	t.Run("Skip", func(t *testing.T) {
		t.Parallel()
		ethAddress := "skip"
		assert.Equal(t, true, isEthereumAddressValid(ethAddress))
	})

	t.Run("Not hex", func(t *testing.T) {
		t.Parallel()
		ethAddress := "abcdefghi"
		assert.Equal(t, false, isEthereumAddressValid(ethAddress))
	})
}
