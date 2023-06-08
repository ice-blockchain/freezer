// SPDX-License-Identifier: ice License 1.0

package extrabonusnotifier

import (
	"fmt"
	stdlibtime "time"

	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func IsExtraBonusAvailable(
	currentTime, extraBonusStartDate, extraBonusStartedAt *time.Time,
	extraBonusIndicesDistribution map[uint16]map[uint16]uint16,
	id int64, utcOffsetFactor int16,
	extraBonusIndex, extraBonusDaysClaimNotAvailable *uint16,
	extraBonusLastClaimAvailableAt **time.Time,
) (available, claimable bool) {
	const notifyHourStart, notifyHourEnd = 10, 20
	var (
		utcOffset                       = stdlibtime.Duration(utcOffsetFactor) * cfg.ExtraBonuses.UTCOffsetDuration
		location                        = stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
		now                             = currentTime.In(location)
		extraBonusStartDateWithLocation = extraBonusStartDate.In(location)
		currentExtraBonusIndex          = 1 + uint16(now.Sub(extraBonusStartDateWithLocation)/cfg.ExtraBonuses.Duration)
		chunkNumber                     = uint16(id % int64(len(extraBonusIndicesDistribution)))
		bonusValue, dayHasExtraBonus    = extraBonusIndicesDistribution[chunkNumber][currentExtraBonusIndex]
	)

	if !(*extraBonusLastClaimAvailableAt).IsNil() {
		*extraBonusIndex = 1 + uint16((*extraBonusLastClaimAvailableAt).In(location).Sub(extraBonusStartDateWithLocation)/cfg.ExtraBonuses.Duration)
	}

	if !dayHasExtraBonus ||
		bonusValue == 0 ||
		*extraBonusIndex >= currentExtraBonusIndex ||
		now.Hour() < notifyHourStart || now.Hour() >= notifyHourEnd ||
		(!extraBonusStartedAt.IsNil() && currentTime.Before(extraBonusStartedAt.Add(cfg.ExtraBonuses.Duration-cfg.ExtraBonuses.AvailabilityWindow))) ||
		(!(*extraBonusLastClaimAvailableAt).IsNil() && currentTime.Before((*extraBonusLastClaimAvailableAt).Add(cfg.ExtraBonuses.Duration-cfg.ExtraBonuses.AvailabilityWindow))) ||
		now.Before(stdlibtime.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), notifyHourStart, 0, 0, 0, location).
			Add(stdlibtime.Duration(chunkNumber)*(cfg.ExtraBonuses.AvailabilityWindow-cfg.ExtraBonuses.ClaimWindow)/stdlibtime.Duration(len(extraBonusIndicesDistribution)))) || //nolint:lll // .
		now.After(stdlibtime.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), notifyHourEnd, 0, 0, 0, location)) {
		log.Info(fmt.Sprintf("extraBonusAvailable:dayHasExtraBonus:%#v,bonusValue:%#v,extraBonusIndex:%#v,currentExtraBonusIndex:%#v,hour:%#v,extraBonusStartedAt:%#v,extraBonusLastClaimAvailableAt:%#v,now:%#v,currentTime:%#v,location:%#v,chunkNumber:%#v",
			dayHasExtraBonus, bonusValue, *extraBonusIndex, currentExtraBonusIndex, now.Hour(), extraBonusStartedAt, extraBonusLastClaimAvailableAt, now, currentTime, location, chunkNumber))
		return false, dayHasExtraBonus && bonusValue > 0 &&
			(*extraBonusIndex == currentExtraBonusIndex) &&
			(now.Hour() >= notifyHourStart && now.Hour() <= notifyHourEnd) &&
			((*extraBonusLastClaimAvailableAt).Before(*currentTime.Time) && (*extraBonusLastClaimAvailableAt).Add(cfg.ExtraBonuses.ClaimWindow).After(*currentTime.Time))
	}
	if !(*extraBonusLastClaimAvailableAt).IsNil() {
		*extraBonusDaysClaimNotAvailable = currentExtraBonusIndex - *extraBonusIndex - 1
		for ix := *extraBonusIndex + 1; ix < currentExtraBonusIndex && *extraBonusDaysClaimNotAvailable > 0; ix++ {
			if pastBonusValue, pastDayHasExtraBonus := extraBonusIndicesDistribution[chunkNumber][ix]; (!pastDayHasExtraBonus || pastBonusValue == 0) && *extraBonusDaysClaimNotAvailable > 0 { //nolint:lll // .
				*extraBonusDaysClaimNotAvailable--
			}
		}
	}
	*extraBonusLastClaimAvailableAt = currentTime
	*extraBonusIndex = currentExtraBonusIndex

	return true, true
}
