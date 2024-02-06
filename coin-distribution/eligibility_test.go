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
	testTime = time.New(stdlibtime.Date(2023, 1, 2, 10, 45, 5, 6, stdlibtime.UTC))
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
		assert.True(t, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
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
		assert.False(t, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
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
		assert.False(t, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
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
		assert.False(t, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
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
		assert.False(t, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
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
		assert.True(t, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("enabled && hour && start < hour < end && truncate(now ~ latest date)", func(t *testing.T) {
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

		assert.True(t, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
	})

	t.Run("enabled && hour && start < hour < end && !truncate(now ~ latest date)", func(t *testing.T) {
		t.Parallel()
		now := testTime
		ethereumDistributionFrequencyMin := 1 * stdlibtime.Hour
		cs := CollectorSettings{
			LatestDate:      time.New(testTime.Add(-30 * stdlibtime.Minute)),
			StartDate:       time.New(testTime.Add(-3 * stdlibtime.Hour)),
			EndDate:         time.New(testTime.Add(3 * stdlibtime.Hour)),
			StartHour:       9,
			Enabled:         true,
			ForcedExecution: false,
		}
		assert.False(t, IsCoinDistributionCollectorEnabled(now, ethereumDistributionFrequencyMin, &cs))
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
		}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.Now(), time.New(now.Add(48*stdlibtime.Hour)), time.New(now.Add(48*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(24 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.FacialRecognitionKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		country = "france"
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		country = "Germany"
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		country = "germany"
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
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
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.True(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
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
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
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
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
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
		miningSessionSoloStartedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		var miningSessionSoloEndedAt *time.Time
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
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
		miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour

		var miningSessionSoloStartedAt *time.Time
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
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
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.Social1KYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		kycState = model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.Social2KYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		kycState = model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), nil}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), nil}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.Social3KYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("Last updated at is nil", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: nil},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("Previous kyc last updated at is nil", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), nil, time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), nil, time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("KYC blocked > kycStep", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.FacialRecognitionKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
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
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(24*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
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
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(48*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: false},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("Quiz completed is true/false && Quiz disabled is true", func(t *testing.T) {
		t.Parallel()

		minMiningStreaksRequired := uint64(1)
		standardBalance := 100.0
		minEthereumDistributionICEBalanceRequired := 10.0
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		country := "France"
		distributionDeniedCountries := map[string]struct{}{}
		now := testTime
		miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate := time.New(now.Add(-25*stdlibtime.Hour)), time.New(now.Add(13*stdlibtime.Hour)), time.New(now.Add(64*stdlibtime.Hour))
		collectingEndedAt := time.New(now.Add(5 * stdlibtime.Hour))
		kycState := model.KYCState{
			KYCStepsCreatedAtField:     model.KYCStepsCreatedAtField{KYCStepsCreatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepsLastUpdatedAtField: model.KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &model.TimeSlice{time.New(now.Add(1 * stdlibtime.Hour)), time.New(now.Add(2 * stdlibtime.Hour)), time.New(now.Add(3 * stdlibtime.Hour)), time.New(now.Add(4 * stdlibtime.Hour))}},
			KYCStepPassedField:         model.KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        model.KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizCompletedField:      model.KYCQuizCompletedField{KYCQuizCompleted: false},
			KYCQuizDisabledField:       model.KYCQuizDisabledField{KYCQuizDisabled: true},
		}
		miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		kycState.KYCQuizCompleted = true
		assert.False(t, IsEligibleForEthereumDistribution(minMiningStreaksRequired, standardBalance, minEthereumDistributionICEBalanceRequired, ethAddress, country, distributionDeniedCountries, now,
			collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate, kycState, miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})
}

func TestIsEligibleForEthereumDistributionNow1(t *testing.T) {
	t.Parallel()

	t.Run("lastEthereumCoinDistributionProcessedAt.IsNil() && today.Equal(ethereumDistributionStartDate)", func(t *testing.T) {
		var lastEthereumCoinDistributionProcessedAt *time.Time
		now := testTime
		coinDistributionStartDate := testTime
		latestCoinDistributionCollectingDate := testTime
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		assert.True(t, IsEligibleForEthereumDistributionNow(1, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(2, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(3, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(4, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(5, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("!lastEthereumCoinDistributionProcessedAt.IsNil() && today.Equal(ethereumDistributionStartDate)", func(t *testing.T) {
		now := testTime
		lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate := testTime, testTime
		latestCoinDistributionCollectingDate := testTime
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := stdlibtime.Hour, 24*stdlibtime.Hour
		assert.False(t, IsEligibleForEthereumDistributionNow(1, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(2, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(3, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(4, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(5, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate) && !latestCoinDistributionCollectingDate.isNil() && now.After(latestCoinDistributionCollectingDay)", func(t *testing.T) {
		lastEthereumCoinDistributionProcessedAt, now := testTime, time.New(testTime.Add(48*stdlibtime.Hour))
		coinDistributionStartDate := time.New(testTime.Add(-24 * stdlibtime.Hour))
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, 24*28*stdlibtime.Hour
		latestCoinDistributionCollectingDate := time.New(testTime.Add(24 * stdlibtime.Hour))
		assert.False(t, IsEligibleForEthereumDistributionNow(1, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(2, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(3, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(4, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(5, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		lastEthereumCoinDistributionProcessedAt = now
		now = time.New(testTime.Add(72 * stdlibtime.Hour))
		latestCoinDistributionCollectingDate = time.New(testTime.Add(48 * stdlibtime.Hour))
		assert.False(t, IsEligibleForEthereumDistributionNow(1, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(2, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(3, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(4, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(5, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))

		lastEthereumCoinDistributionProcessedAt = now
		now = time.New(testTime.Add(96 * stdlibtime.Hour))
		latestCoinDistributionCollectingDate = time.New(testTime.Add(72 * stdlibtime.Hour))
		assert.False(t, IsEligibleForEthereumDistributionNow(1, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(2, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(3, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(4, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(5, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate) && !latestCoinDistributionCollectingDate.isNil() && now.Before(latestCoinDistributionCollectingDay)", func(t *testing.T) {
		lastEthereumCoinDistributionProcessedAt, now := testTime, testTime
		coinDistributionStartDate := time.New(testTime.Add(-24 * stdlibtime.Hour))
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, 24*28*stdlibtime.Hour
		latestCoinDistributionCollectingDate := time.New(testTime.Add(24 * stdlibtime.Hour))
		assert.False(t, IsEligibleForEthereumDistributionNow(1, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(2, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(3, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(4, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(5, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate) && latestCoinDistributionCollectingDate.isNil() && today ~ coinDistributionStartDate", func(t *testing.T) {
		lastEthereumCoinDistributionProcessedAt, now := testTime, time.New(testTime.Add(52*stdlibtime.Hour))
		coinDistributionStartDate := time.New(testTime.Add(48 * stdlibtime.Hour))
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, 24*28*stdlibtime.Hour
		var latestCoinDistributionCollectingDate *time.Time
		assert.True(t, IsEligibleForEthereumDistributionNow(1, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(2, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(3, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(4, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.True(t, IsEligibleForEthereumDistributionNow(5, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})

	t.Run("lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate) && latestCoinDistributionCollectingDate.isNil() && today !~ coinDistributionStartDate", func(t *testing.T) {
		lastEthereumCoinDistributionProcessedAt, now := testTime, time.New(testTime.Add(72*stdlibtime.Hour))
		coinDistributionStartDate := time.New(testTime.Add(48 * stdlibtime.Hour))
		ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax := 24*stdlibtime.Hour, 24*28*stdlibtime.Hour
		var latestCoinDistributionCollectingDate *time.Time
		assert.False(t, IsEligibleForEthereumDistributionNow(1, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(2, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(3, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(4, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
		assert.False(t, IsEligibleForEthereumDistributionNow(5, now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate, latestCoinDistributionCollectingDate, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax))
	})
}

func TestIsEthereumAddressValid(t *testing.T) {
	t.Parallel()

	t.Run("Positive case", func(t *testing.T) {
		t.Parallel()
		ethAddress := "0x111Fc57e1e4e7687c9195F7856C45227f269323B"
		assert.True(t, isEthereumAddressValid(ethAddress))
	})

	t.Run("Empty eth address", func(t *testing.T) {
		t.Parallel()
		ethAddress := ""
		assert.False(t, isEthereumAddressValid(ethAddress))
	})

	t.Run("Skip", func(t *testing.T) {
		t.Parallel()
		ethAddress := "skip"
		assert.True(t, isEthereumAddressValid(ethAddress))
	})

	t.Run("Not hex", func(t *testing.T) {
		t.Parallel()
		ethAddress := "abcdefghi"
		assert.False(t, isEthereumAddressValid(ethAddress))
	})
}
