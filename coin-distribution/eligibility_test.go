// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/time"
)

func TestIsEligibleForEthereumDistributionNow(t *testing.T) {
	t.Parallel()
	coinDistributionStartDate := time.New(stdlibtime.Date(2024, 1, 16, 0, 0, 0, 0, stdlibtime.UTC))
	lastCoinDistributionProcessedAt := time.New(stdlibtime.Date(2024, 1, 19, 0, 0, 0, 0, stdlibtime.UTC))
	now := time.New(stdlibtime.Date(2024, 1, 22, 14, 17, 33, 0, stdlibtime.UTC))

	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 22, 0, 0, 0, 0, stdlibtime.UTC))
	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 21, 0, 0, 0, 0, stdlibtime.UTC))
	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 16, 0, 0, 0, 0, stdlibtime.UTC))
	assert.True(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCoinDistributionProcessedAt))
}

func TestFinalDistribution(t *testing.T) {
	now := time.New(stdlibtime.Date(2024, 2, 27, 18, 0, 1, 0, stdlibtime.UTC))
	coinDistributionStartDate := time.New(stdlibtime.Date(2024, 2, 27, 0, 0, 0, 0, stdlibtime.UTC))
	coinDistributionEndDate := time.New(stdlibtime.Date(2024, 2, 27, 0, 0, 0, 0, stdlibtime.UTC))
	var lastCollectingProcessedAt *time.Time = new(time.Time)

	assert.True(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCollectingProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCollectingProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCollectingProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCollectingProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCollectingProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionStartDate, lastCollectingProcessedAt))

	firstTimeUserLastEthereumProcessedAt := new(time.Time)
	assert.True(t, IsEligibleForEthereumDistributionNow(0, now, firstTimeUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(1, now, firstTimeUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(2, now, firstTimeUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(3, now, firstTimeUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(4, now, firstTimeUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(5, now, firstTimeUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))

	previouslyDistributedUserLastEthereumProcessedAt := time.New(stdlibtime.Date(2024, 2, 15, 16, 0, 0, 0, stdlibtime.UTC))
	assert.True(t, IsEligibleForEthereumDistributionNow(0, now, previouslyDistributedUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(1, now, previouslyDistributedUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(2, now, previouslyDistributedUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(3, now, previouslyDistributedUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(4, now, previouslyDistributedUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.True(t, IsEligibleForEthereumDistributionNow(5, now, previouslyDistributedUserLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))

	userAlreadyProcessedLastEthereumProcessedAt := time.New(stdlibtime.Date(2024, 2, 27, 0, 0, 0, 0, stdlibtime.UTC))
	assert.False(t, IsEligibleForEthereumDistributionNow(0, now, userAlreadyProcessedLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistributionNow(1, now, userAlreadyProcessedLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistributionNow(2, now, userAlreadyProcessedLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistributionNow(3, now, userAlreadyProcessedLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistributionNow(4, now, userAlreadyProcessedLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistributionNow(5, now, userAlreadyProcessedLastEthereumProcessedAt, coinDistributionStartDate, lastCollectingProcessedAt, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))

	balanceRequired := float64(0)
	var kycPassed model.KYCState = buildKYC(true, now)
	var kycNotPassed model.KYCState = buildKYC(false, now)
	var collectingEndedAt *time.Time
	if lastCollectingProcessedAt.IsNil() {
		collectingEndedAt = time.New(time.Now().Add(-1 * stdlibtime.Millisecond).Add(20 * stdlibtime.Minute))
	}
	miningSessionDuration := 24 * stdlibtime.Hour
	activeMiningSessionStarted := time.New(stdlibtime.Date(2024, 2, 25, 0, 0, 0, 0, stdlibtime.UTC))
	nonActiveMiningSessionEnded := time.New(stdlibtime.Date(2024, 2, 26, 0, 0, 0, 0, stdlibtime.UTC))
	assert.True(t, IsEligibleForEthereumDistribution(uint64(0), 0.1, balanceRequired, "skip", "US", make(map[string]struct{}), now, collectingEndedAt, activeMiningSessionStarted, nonActiveMiningSessionEnded, coinDistributionEndDate, kycPassed, miningSessionDuration, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistribution(uint64(0), 1, balanceRequired, "skip", "US", make(map[string]struct{}), now, collectingEndedAt, activeMiningSessionStarted, nonActiveMiningSessionEnded, coinDistributionEndDate, kycNotPassed, miningSessionDuration, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistribution(uint64(0), 1, balanceRequired, "bogusInvalidAddress", "US", make(map[string]struct{}), now, collectingEndedAt, activeMiningSessionStarted, nonActiveMiningSessionEnded, coinDistributionEndDate, kycPassed, miningSessionDuration, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistribution(uint64(0), 1, balanceRequired, "", "US", make(map[string]struct{}), now, collectingEndedAt, activeMiningSessionStarted, nonActiveMiningSessionEnded, coinDistributionEndDate, kycPassed, miningSessionDuration, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))

	activeMiningSessionEnded := time.New(stdlibtime.Date(2024, 2, 29, 0, 0, 0, 0, stdlibtime.UTC))
	assert.True(t, IsEligibleForEthereumDistribution(uint64(0), 0.1, balanceRequired, "skip", "US", make(map[string]struct{}), now, collectingEndedAt, activeMiningSessionStarted, activeMiningSessionEnded, coinDistributionEndDate, kycPassed, miningSessionDuration, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistribution(uint64(0), 1, balanceRequired, "skip", "US", make(map[string]struct{}), now, collectingEndedAt, activeMiningSessionStarted, activeMiningSessionEnded, coinDistributionEndDate, kycNotPassed, miningSessionDuration, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistribution(uint64(0), 1, balanceRequired, "bogusInvalidAddress", "US", make(map[string]struct{}), now, collectingEndedAt, activeMiningSessionEnded, nonActiveMiningSessionEnded, coinDistributionEndDate, kycPassed, miningSessionDuration, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))
	assert.False(t, IsEligibleForEthereumDistribution(uint64(0), 1, balanceRequired, "", "US", make(map[string]struct{}), now, collectingEndedAt, activeMiningSessionStarted, activeMiningSessionEnded, coinDistributionEndDate, kycPassed, miningSessionDuration, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour))

	assert.Equal(t, float64(0), CalculateEthereumDistributionICEBalance(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionEndDate))
	assert.Equal(t, float64(100), CalculateEthereumDistributionICEBalance(100, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, now, coinDistributionEndDate))

	finalDistributionSettings := &CollectorSettings{
		DeniedCountries:          nil,
		LatestDate:               nil,
		StartDate:                coinDistributionStartDate,
		EndDate:                  coinDistributionEndDate,
		MinBalanceRequired:       0,
		StartHour:                18,
		MinMiningStreaksRequired: 0,
		Enabled:                  true,
		ForcedExecution:          false,
	}
	assert.False(t, IsCoinDistributionCollectorEnabled(time.New(now.Add(-1*24*stdlibtime.Hour)), 24*stdlibtime.Hour, finalDistributionSettings))
	assert.False(t, IsCoinDistributionCollectorEnabled(time.New(now.Add(-2*stdlibtime.Minute)), 24*stdlibtime.Hour, finalDistributionSettings))
	assert.False(t, IsCoinDistributionCollectorEnabled(now, 24*stdlibtime.Hour, finalDistributionSettings))
	assert.True(t, IsCoinDistributionCollectorEnabled(time.New(now.Add(20*stdlibtime.Minute)), 24*stdlibtime.Hour, finalDistributionSettings))
	assert.False(t, IsCoinDistributionCollectorEnabled(time.New(now.Add(1*stdlibtime.Hour)), 24*stdlibtime.Hour, finalDistributionSettings))
	assert.True(t, IsCoinDistributionCollectorEnabled(time.New(now.Add(1*stdlibtime.Hour).Add(20*stdlibtime.Minute)), 24*stdlibtime.Hour, finalDistributionSettings))
	assert.False(t, IsCoinDistributionCollectorEnabled(time.New(now.Add(1*24*stdlibtime.Hour)), 24*stdlibtime.Hour, finalDistributionSettings))

}

func buildKYC(passed bool, now *time.Time) model.KYCState {
	var kyc model.KYCState
	if passed {
		kyc.KYCStepPassed = users.QuizKYCStep
		kyc.KYCQuizCompleted = true
	} else {
		kyc.KYCStepPassed = users.LivenessDetectionKYCStep
	}
	kycPassedTimes := model.TimeSlice([]*time.Time{
		now, now, now, now,
	})
	kyc.KYCStepsCreatedAt = &kycPassedTimes
	kyc.KYCStepsLastUpdatedAt = &kycPassedTimes

	return kyc
}
