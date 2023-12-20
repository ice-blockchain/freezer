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

const (
	miningSessionSoloEndedAtNetworkDelayAdjustment = 20 * stdlibtime.Second
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
	now, miningSessionSoloStartedAt, miningSessionSoloEndedAt, ethereumDistributionEndDate *time.Time,
	kycState model.KYCState,
	miningSessionDuration, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax stdlibtime.Duration) bool {
	var countryAllowed bool
	if _, countryDenied := distributionDeniedCountries[strings.ToLower(country)]; country != "" && !countryDenied {
		countryAllowed = true
	}

	return countryAllowed &&
		!miningSessionSoloEndedAt.IsNil() && miningSessionSoloEndedAt.After(now.Add(miningSessionSoloEndedAtNetworkDelayAdjustment)) &&
		isEthereumAddressValid(ethAddress) &&
		CalculateEthereumDistributionICEBalance(standardBalance, ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax, now, ethereumDistributionEndDate) >= minEthereumDistributionICEBalanceRequired && //nolint:lll // .
		model.CalculateMiningStreak(now, miningSessionSoloStartedAt, miningSessionSoloEndedAt, miningSessionDuration) >= minMiningStreaksRequired &&
		kycState.KYCStepPassedCorrectly(users.QuizKYCStep)
}

//nolint:funlen // .
func IsEligibleForEthereumDistributionNow(id int64,
	now, lastEthereumCoinDistributionProcessedAt, coinDistributionStartDate *time.Time,
	ethereumDistributionFrequencyMin, ethereumDistributionFrequencyMax stdlibtime.Duration) bool {
	var (
		reservationIsTodayAndIsOutsideOfFirstCycle, secondReservationIsWithinFirstDistributionCycle bool
		truncatedLastEthereumCoinDistributionProcessedAt                                            *time.Time
		ethereumDistributionStartDate                                                               = coinDistributionStartDate.Truncate(ethereumDistributionFrequencyMin)          //nolint:lll // .
		ethereumDistributionDayAfterStartDate                                                       = ethereumDistributionStartDate.Add(ethereumDistributionFrequencyMin)           //nolint:lll // .
		ethereumDistributionFirstCycleEndDate                                                       = ethereumDistributionDayAfterStartDate.Add(ethereumDistributionFrequencyMax)   //nolint:lll // .
		userReservedDayForEthereumCoinDistributionIndex                                             = id % int64(ethereumDistributionFrequencyMax/ethereumDistributionFrequencyMin) //nolint:lll // .
		today                                                                                       = now.Truncate(ethereumDistributionFrequencyMin)
		neverDoneItBeforeAndTodayIsNotEthereumDistributionStartDateButReservationIsToday            = userReservedDayForEthereumCoinDistributionIndex == int64((today.Sub(ethereumDistributionDayAfterStartDate)%ethereumDistributionFrequencyMax)/ethereumDistributionFrequencyMin) //nolint:lll // .
	)
	if !lastEthereumCoinDistributionProcessedAt.IsNil() {
		truncatedLastEthereumCoinDistributionProcessedAt = time.New(lastEthereumCoinDistributionProcessedAt.Truncate(ethereumDistributionFrequencyMin))
		reservationIsTodayAndIsOutsideOfFirstCycle = today.Sub(*truncatedLastEthereumCoinDistributionProcessedAt.Time)%ethereumDistributionFrequencyMax == 0                                                                                                                                                              //nolint:lll // .
		secondReservationIsWithinFirstDistributionCycle = userReservedDayForEthereumCoinDistributionIndex == int64((truncatedLastEthereumCoinDistributionProcessedAt.Add(ethereumDistributionFrequencyMin).Sub(ethereumDistributionDayAfterStartDate)%ethereumDistributionFrequencyMax)/ethereumDistributionFrequencyMin) //nolint:lll // .
	}
	switch {
	case lastEthereumCoinDistributionProcessedAt.IsNil() && today.Equal(ethereumDistributionStartDate):
		return true
	case !lastEthereumCoinDistributionProcessedAt.IsNil() && today.Equal(ethereumDistributionStartDate):
		return false
	case lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate):
		return neverDoneItBeforeAndTodayIsNotEthereumDistributionStartDateButReservationIsToday
	case !lastEthereumCoinDistributionProcessedAt.IsNil() && !today.Equal(ethereumDistributionStartDate):
		switch {
		case today.Before(ethereumDistributionFirstCycleEndDate):
			return secondReservationIsWithinFirstDistributionCycle
		default:
			return reservationIsTodayAndIsOutsideOfFirstCycle
		}
	default:
		return false
	}
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
