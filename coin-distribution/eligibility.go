// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"strings"
	stdlibtime "time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/time"
)

func IsCoinDistributionCollectorEnabled(now *time.Time, ethereumDistributionFrequencyMin stdlibtime.Duration, cs *CollectorSettings) bool {
	return cs.Enabled &&
		(cs.ForcedExecution ||
			(now.Hour() >= cs.StartHour &&
				now.After(*cs.StartDate.Time) &&
				now.Before(*cs.EndDate.Time) &&
				(cs.LatestDate.IsNil() ||
					!now.Truncate(ethereumDistributionFrequencyMin).Equal(cs.LatestDate.Truncate(ethereumDistributionFrequencyMin)))))
}

func CalculateEthereumDistributionICEBalance(
	standardBalance float64,
	ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax stdlibtime.Duration,
	now, ethereumDistributionEndDate *time.Time,
) float64 {
	delta := ethereumDistributionEndDate.Truncate(ethereumDistributionFrequencyMin).Sub(now.Truncate(ethereumDistributionFrequencyMin))
	if delta <= ethereumDistributionFrequencyMax {
		return standardBalance
	}

	return standardBalance / float64(int64(delta/ethereumDistributionFrequencyMax)+1)
}

func IsEligibleForEthereumDistribution(
	minMiningStreaksRequired uint64,
	standardBalance, minEthereumDistributionICEBalanceRequired float64,
	ethAddress, country string,
	distributionDeniedCountries map[string]struct{},
	now, collectingEndedAt, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate *time.Time,
	kycState model.KYCState,
	miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax stdlibtime.Duration) bool {
	var countryAllowed bool
	if _, countryDenied := distributionDeniedCountries[strings.ToLower(country)]; len(distributionDeniedCountries) == 0 || (country != "" && !countryDenied) {
		countryAllowed = true
	}
	distributedBalance := CalculateEthereumDistributionICEBalance(standardBalance, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax, now, ethereumDistributionEndDate) //nolint:lll // .

	return countryAllowed &&
		!miningSessionSoloEndedAt.IsNil() && miningSessionSoloEndedAt.After(*collectingEndedAt.Time) &&
		isEthereumAddressValid(ethAddress) &&
		((minEthereumDistributionICEBalanceRequired > 0 && distributedBalance >= minEthereumDistributionICEBalanceRequired) || (distributedBalance > 0)) &&
		model.CalculateMiningStreak(now, miningSessionSoloStartedAt, miningSessionSoloEndedAt, miningSessionDuration) >= minMiningStreaksRequired &&
		kycState.KYCStepPassedCorrectly(users.QuizKYCStep)
}

func IsEligibleForEthereumDistributionNow(id int64,
	now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate *time.Time,
	ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax stdlibtime.Duration) bool {
	return (lastEthereumCoinDistributionProcessedAt.IsNil() && now.Truncate(ethereumDistributionFrequencyMin).Equal(coinDistributionStartDate.Truncate(ethereumDistributionFrequencyMin))) || //nolint:lll // .
		(id%int64(ethereumDistributionFrequencyMax/ethereumDistributionFrequencyMin)) == int64((now.Truncate(ethereumDistributionFrequencyMin).Sub(coinDistributionStartDate.Truncate(ethereumDistributionFrequencyMin).Add(ethereumDistributionFrequencyMin))%ethereumDistributionFrequencyMax)/ethereumDistributionFrequencyMin) //nolint:lll // .
}

func isEthereumAddressValid(ethAddress string) bool {
	if ethAddress == "" {
		return false
	}
	if ethAddress == "skip" {
		return true
	}

	return common.IsHexAddress(ethAddress)
}
