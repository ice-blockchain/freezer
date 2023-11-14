// SPDX-License-Identifier: ice License 1.0

package model

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

type (
	User struct {
		BalanceLastUpdatedAtField
		MiningSessionSoloLastStartedAtField
		MiningSessionSoloStartedAtField
		MiningSessionSoloEndedAtField
		MiningSessionSoloPreviouslyEndedAtField
		ExtraBonusStartedAtField
		ResurrectSoloUsedAtField
		ResurrectT0UsedAtField
		ResurrectTMinus1UsedAtField
		MiningSessionSoloDayOffLastAwardedAtField
		ExtraBonusLastClaimAvailableAtField
		ProfilePictureNameField
		UsernameField
		MiningBlockchainAccountAddressField
		BlockchainAccountAddressField
		UserIDField
		DeserializedUsersKey
		BalanceTotalStandardField
		BalanceTotalPreStakingField
		BalanceTotalMintedField
		BalanceTotalSlashedField
		BalanceSoloPendingField
		BalanceT1PendingField
		BalanceT2PendingField
		BalanceSoloPendingAppliedField
		BalanceT1PendingAppliedField
		BalanceT2PendingAppliedField
		BalanceSoloField
		BalanceT0Field
		BalanceT1Field
		BalanceT2Field
		BalanceForT0Field
		BalanceForTMinus1Field
		SlashingRateSoloField
		SlashingRateT0Field
		SlashingRateT1Field
		SlashingRateT2Field
		SlashingRateForT0Field
		SlashingRateForTMinus1Field
		PreStakingBonusField
		PreStakingAllocationField
		ExtraBonusField
		IDT0Field
		IDTMinus1Field
		UTCOffsetField
		ActiveT1ReferralsField
		ActiveT2ReferralsField
		NewsSeenField
		ExtraBonusDaysClaimNotAvailableField
		HideRankingField
	}
	BalanceLastUpdatedAtField struct {
		BalanceLastUpdatedAt *time.Time `redis:"balance_last_updated_at,omitempty"`
	}
	MiningSessionSoloLastStartedAtField struct {
		MiningSessionSoloLastStartedAt *time.Time `redis:"mining_session_solo_last_started_at,omitempty"`
	}
	MiningSessionSoloStartedAtField struct {
		MiningSessionSoloStartedAt *time.Time `redis:"mining_session_solo_started_at,omitempty"`
	}
	MiningSessionSoloEndedAtField struct {
		MiningSessionSoloEndedAt *time.Time `redis:"mining_session_solo_ended_at,omitempty"`
	}
	MiningSessionSoloPreviouslyEndedAtField struct {
		MiningSessionSoloPreviouslyEndedAt *time.Time `redis:"mining_session_solo_previously_ended_at,omitempty"`
	}
	ExtraBonusStartedAtField struct {
		ExtraBonusStartedAt *time.Time `redis:"extra_bonus_started_at,omitempty"`
	}
	ResurrectSoloUsedAtField struct {
		ResurrectSoloUsedAt *time.Time `redis:"resurrect_solo_used_at,omitempty"`
	}
	ResurrectT0UsedAtField struct {
		ResurrectT0UsedAt *time.Time `redis:"resurrect_t0_used_at,omitempty"`
	}
	ResurrectTMinus1UsedAtField struct {
		ResurrectTMinus1UsedAt *time.Time `redis:"resurrect_tminus1_used_at,omitempty"`
	}
	MiningSessionSoloDayOffLastAwardedAtField struct {
		MiningSessionSoloDayOffLastAwardedAt *time.Time `redis:"mining_session_solo_day_off_last_awarded_at,omitempty"`
	}
	ExtraBonusLastClaimAvailableAtField struct {
		ExtraBonusLastClaimAvailableAt *time.Time `redis:"extra_bonus_last_claim_available_at,omitempty"`
	}
	BalancesBackupUsedAtField struct {
		BalancesBackupUsedAt *time.Time `redis:"balance_backup_used_at,omitempty"`
	}
	UserIDField struct {
		UserID string `redis:"user_id,omitempty"`
	}
	ProfilePictureNameField struct {
		ProfilePictureName string `redis:"profile_picture_name,omitempty"`
	}
	UsernameField struct {
		Username string `redis:"username,omitempty"`
	}
	MiningBlockchainAccountAddressField struct {
		MiningBlockchainAccountAddress string `redis:"mining_blockchain_account_address,omitempty"`
	}
	BlockchainAccountAddressField struct {
		BlockchainAccountAddress string `redis:"blockchain_account_address,omitempty"`
	}
	LatestDeviceField struct {
		LatestDevice string `redis:"latest_device,omitempty"`
	}
	BalanceTotalStandardField struct {
		BalanceTotalStandard float64 `redis:"balance_total_standard"`
	}
	BalanceTotalPreStakingField struct {
		BalanceTotalPreStaking float64 `redis:"balance_total_pre_staking"`
	}
	BalanceTotalMintedField struct {
		BalanceTotalMinted float64 `redis:"balance_total_minted"`
	}
	BalanceTotalSlashedField struct {
		BalanceTotalSlashed float64 `redis:"balance_total_slashed"`
	}
	BalanceSoloPendingField struct {
		BalanceSoloPending float64 `redis:"balance_solo_pending,omitempty"`
	}
	BalanceT1PendingField struct {
		BalanceT1Pending float64 `redis:"balance_t1_pending,omitempty"`
	}
	BalanceT2PendingField struct {
		BalanceT2Pending float64 `redis:"balance_t2_pending,omitempty"`
	}
	BalanceSoloPendingAppliedField struct {
		BalanceSoloPendingApplied float64 `redis:"balance_solo_pending_applied,omitempty"`
	}
	BalanceT1PendingAppliedField struct {
		BalanceT1PendingApplied float64 `redis:"balance_t1_pending_applied,omitempty"`
	}
	BalanceT2PendingAppliedField struct {
		BalanceT2PendingApplied float64 `redis:"balance_t2_pending_applied,omitempty"`
	}
	BalanceSoloField struct {
		BalanceSolo float64 `redis:"balance_solo"`
	}
	BalanceT0Field struct {
		BalanceT0 float64 `redis:"balance_t0"`
	}
	BalanceT1Field struct {
		BalanceT1 float64 `redis:"balance_t1"`
	}
	BalanceT2Field struct {
		BalanceT2 float64 `redis:"balance_t2"`
	}
	FirstRecalculatedBalanceT1Field struct {
		FirstRecalculatedBalanceT1 float64 `redis:"first_recalculated_balance_t1"`
	}
	FirstRecalculatedBalanceT2Field struct {
		FirstRecalculatedBalanceT2 float64 `redis:"first_recalculated_balance_t2"`
	}
	BalanceForT0Field struct {
		BalanceForT0 float64 `redis:"balance_for_t0"`
	}
	BalanceForTMinus1Field struct {
		BalanceForTMinus1 float64 `redis:"balance_for_tminus1"`
	}
	SlashingRateSoloField struct {
		SlashingRateSolo float64 `redis:"slashing_rate_solo"`
	}
	SlashingRateT0Field struct {
		SlashingRateT0 float64 `redis:"slashing_rate_t0"`
	}
	SlashingRateT1Field struct {
		SlashingRateT1 float64 `redis:"slashing_rate_t1"`
	}
	SlashingRateT2Field struct {
		SlashingRateT2 float64 `redis:"slashing_rate_t2"`
	}
	FirstRecalculatedSlashingRateT1Field struct {
		FirstRecalculatedSlashingRateT1 float64 `redis:"first_recalculated_slashing_rate_t1"`
	}
	FirstRecalculatedSlashingRateT2Field struct {
		FirstRecalculatedSlashingRateT2 float64 `redis:"first_recalculated_slashing_rate_t2"`
	}
	SlashingRateForT0Field struct {
		SlashingRateForT0 float64 `redis:"slashing_rate_for_t0"`
	}
	SlashingRateForTMinus1Field struct {
		SlashingRateForTMinus1 float64 `redis:"slashing_rate_for_tminus1"`
	}
	PreStakingBonusField struct {
		PreStakingBonus float64 `redis:"pre_staking_bonus,omitempty"`
	}
	PreStakingAllocationField struct {
		PreStakingAllocation float64 `redis:"pre_staking_allocation,omitempty"`
	}
	ExtraBonusField struct {
		ExtraBonus float64 `redis:"extra_bonus,omitempty"`
	}
	DeserializedUsersKey struct {
		ID int64 `redis:"-"`
	}
	IDT0Field struct {
		IDT0 int64 `redis:"id_t0,omitempty"`
	}
	IDTMinus1Field struct {
		IDTMinus1 int64 `redis:"id_tminus1,omitempty"`
	}
	IDT0ResettableField struct {
		IDT0 int64 `redis:"id_t0"`
	}
	IDTMinus1ResettableField struct {
		IDTMinus1 int64 `redis:"id_tminus1"`
	}
	ActiveT1ReferralsField struct {
		ActiveT1Referrals int32 `redis:"active_t1_referrals,omitempty"`
	}
	ActiveT2ReferralsField struct {
		ActiveT2Referrals int32 `redis:"active_t2_referrals,omitempty"`
	}
	FirstRecalculatedActiveT1ReferralsField struct {
		FirstRecalculatedActiveT1Referrals int32 `redis:"first_recalculated_active_t1_referrals,omitempty"`
	}
	FirstRecalculatedActiveT2ReferralsField struct {
		FirstRecalculatedActiveT2Referrals int32 `redis:"first_recalculated_active_t2_referrals,omitempty"`
	}
	NewsSeenField struct {
		NewsSeen uint16 `redis:"news_seen"`
	}
	ExtraBonusDaysClaimNotAvailableResettableField struct {
		ExtraBonusDaysClaimNotAvailable uint16 `redis:"extra_bonus_days_claim_not_available"`
	}
	ExtraBonusDaysClaimNotAvailableField struct {
		ExtraBonusDaysClaimNotAvailable uint16 `redis:"extra_bonus_days_claim_not_available,omitempty"`
	}
	UTCOffsetField struct {
		UTCOffset int64 `redis:"utc_offset"`
	}
	HideRankingField struct {
		HideRanking bool `redis:"hide_ranking"`
	}
	KYCStepsCreatedAtField struct {
		KYCStepsCreatedAt *TimeSlice `json:"kycStepsCreatedAt" redis:"kyc_steps_created_at"`
	}
	KYCStepsLastUpdatedAtField struct {
		KYCStepsLastUpdatedAt *TimeSlice `json:"kycStepsLastUpdatedAt" redis:"kyc_steps_last_updated_at"`
	}
	KYCStepPassedField struct {
		KYCStepPassed users.KYCStep `json:"kycStepPassed" redis:"kyc_step_passed"`
	}
	KYCStepBlockedField struct {
		KYCStepBlocked users.KYCStep `json:"kycStepBlocked" redis:"kyc_step_blocked"`
	}
	DeserializedBackupUsersKey struct {
		ID int64 `redis:"-"`
	}
)

func (k *DeserializedUsersKey) Key() string {
	if k == nil || k.ID == 0 {
		return ""
	}

	return SerializedUsersKey(k.ID)
}

func (k *DeserializedUsersKey) SetKey(val string) {
	if val == "" || val == "users:" {
		return
	}
	if val[0] == 'u' {
		val = val[6:]
	}
	var err error
	k.ID, err = strconv.ParseInt(val, 10, 64)
	log.Panic(err)
}

func SerializedUsersKey(val any) string {
	switch typedVal := val.(type) {
	case string:
		if typedVal == "" {
			return ""
		}

		return "users:" + typedVal
	case int64:
		if typedVal == 0 {
			return ""
		}

		return "users:" + strconv.FormatInt(typedVal, 10)
	default:
		panic(fmt.Sprintf("%#v cannot be used as users key", val))
	}
}

func (k *DeserializedBackupUsersKey) Key() string {
	if k == nil || k.ID == 0 {
		return ""
	}

	return SerializedBackupUsersKey(k.ID)
}

func (k *DeserializedBackupUsersKey) SetKey(val string) {
	if val == "" || val == "backup:" {
		return
	}
	if val[0] == 'b' {
		val = val[7:]
	}
	var err error
	k.ID, err = strconv.ParseInt(val, 10, 64)
	log.Panic(err)
}

func SerializedBackupUsersKey(val any) string {
	switch typedVal := val.(type) {
	case string:
		if typedVal == "" {
			return ""
		}

		return "backup:" + typedVal
	case int64:
		if typedVal == 0 {
			return ""
		}

		return "backup:" + strconv.FormatInt(typedVal, 10)
	default:
		panic(fmt.Sprintf("%#v cannot be used as backup key", val))
	}
}

type (
	TimeSlice []*time.Time
)

func (t *TimeSlice) Equals(other *TimeSlice) bool {
	if t != nil && other != nil && len(*t) == (len(*other)) {
		var equals int
		for ix, thisVal := range *t {
			if otherVal := (*other)[ix]; (thisVal.IsNil() && otherVal.IsNil()) || (!thisVal.IsNil() && !otherVal.IsNil() && thisVal.Equal(*otherVal.Time)) {
				equals++
			}
		}

		return equals == len(*t)
	}

	return t == nil && other == nil
}

func (t *TimeSlice) UnmarshalBinary(text []byte) error {
	return t.UnmarshalText(text)
}

func (t *TimeSlice) UnmarshalJSON(text []byte) error {
	if len(text) == 0 || (len(text) == 2 && string(text) == "[]") {
		return nil
	}
	sep := make([]byte, 1)
	sep[0] = ','
	elems := bytes.Split(text[1:len(text)-1], sep)
	timeSlice := make(TimeSlice, 0, len(elems))
	for _, val := range elems {
		unmarshalledTime := new(time.Time)
		if err := unmarshalledTime.UnmarshalJSON(context.Background(), val); err != nil {
			return errors.Wrapf(err, "failed to UnmarshalJSON %#v:%v", unmarshalledTime, string(val))
		}
		if !unmarshalledTime.IsNil() {
			timeSlice = append(timeSlice, unmarshalledTime)
		}
	}
	*t = timeSlice

	return nil
}

func (t *TimeSlice) UnmarshalText(text []byte) error {
	if len(text) == 0 || (len(text) == 1 && string(text) == "") {
		return nil
	}
	sep := make([]byte, 1)
	sep[0] = ','
	elems := bytes.Split(text, sep)
	timeSlice := make(TimeSlice, 0, len(elems))
	for _, val := range elems {
		unmarshalledTime := new(time.Time)
		if err := unmarshalledTime.UnmarshalText(val); err != nil {
			return errors.Wrapf(err, "failed to unmarshall %#v:%v", unmarshalledTime, string(val))
		}
		if !unmarshalledTime.IsNil() {
			timeSlice = append(timeSlice, unmarshalledTime)
		}
	}
	*t = timeSlice

	return nil
}

func (t *TimeSlice) MarshalText() ([]byte, error) {
	return t.MarshalBinary()
}

func (t *TimeSlice) MarshalBinary() ([]byte, error) {
	if t == nil || len(*t) == 0 {
		return nil, nil
	}
	timeSlice := *t
	text := make([]byte, 0, len(timeSlice)*(len(stdlibtime.RFC3339Nano)+1))
	for ix, val := range timeSlice {
		marshalledTime, err := val.MarshalText()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshall: %#v", val)
		}
		if len(marshalledTime) == 0 {
			continue
		}
		text = append(text, marshalledTime...)
		if ix != len(timeSlice)-1 {
			text = append(text, ',')
		}
	}
	if len(text) == 0 {
		return nil, nil
	}

	return text, nil
}
