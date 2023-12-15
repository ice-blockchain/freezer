// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	stdlibtime "time"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/time"
)

const (
	miningSessionSoloEndedAtNetworkDelayAdjustment = 20 * stdlibtime.Second
)

func CalculateEthereumDistributionICEBalance(
	standardBalance float64,
	ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax stdlibtime.Duration,
	now, ethereumDistributionEndDate *time.Time,
) float64 {
	delta := ethereumDistributionEndDate.Truncate(ethereumDistributionFrequencyMin).Sub(now.Truncate(ethereumDistributionFrequencyMin))
	if delta <= ethereumDistributionFrequencyMax {
		return standardBalance
	}

	//TODO: should this be fractional or natural?
	return standardBalance / (float64(delta.Nanoseconds()) / float64(ethereumDistributionFrequencyMax.Nanoseconds()))
}

func IsEligibleForEthereumDistribution(
	minMiningStreaksRequired uint64,
	standardBalance, minEthereumDistributionICEBalanceRequired float64,
	ethAddress, country string,
	distributionDeniedCountries map[string]struct{},
	now, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate *time.Time,
	kycState model.KYCState,
	miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax stdlibtime.Duration) bool {
	var countryAllowed bool
	if _, countryDenied := distributionDeniedCountries[country]; country != "" && !countryDenied {
		countryAllowed = true
	}

	return countryAllowed &&
		!miningSessionSoloEndedAt.IsNil() && miningSessionSoloEndedAt.After(now.Add(miningSessionSoloEndedAtNetworkDelayAdjustment)) &&
		isEthereumAddressValid(ethAddress) &&
		CalculateEthereumDistributionICEBalance(standardBalance, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax, now, ethereumDistributionEndDate) >= minEthereumDistributionICEBalanceRequired && //nolint:lll // .
		model.CalculateMiningStreak(now, miningSessionSoloStartedAt, miningSessionSoloEndedAt, miningSessionDuration) >= minMiningStreaksRequired &&
		kycState.KYCStepPassedCorrectly(users.QuizKYCStep)
}

func isEthereumAddressValid(ethAddress string) bool {
	if ethAddress == "" {
		return false
	}
	if ethAddress == "skip" {
		return true
	}

	return true
}
