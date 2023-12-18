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
		SoloLastEthereumCoinDistributionProcessedAtField
		ForT0LastEthereumCoinDistributionProcessedAtField
		ForTMinus1LastEthereumCoinDistributionProcessedAtField
		KYCState
		ProfilePictureNameField
		UsernameField
		MiningBlockchainAccountAddressField
		BlockchainAccountAddressField
		UserIDField
		CountryField
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
		BalanceSoloEthereumField
		BalanceT0EthereumField
		BalanceT1EthereumField
		BalanceT2EthereumField
		BalanceForT0EthereumField
		BalanceForTMinus1EthereumField
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
	KYCState struct {
		KYCStepsCreatedAtField
		KYCStepsLastUpdatedAtField
		KYCStepPassedField
		KYCStepBlockedField
	}
	SoloLastEthereumCoinDistributionProcessedAtField struct {
		SoloLastEthereumCoinDistributionProcessedAt *time.Time `redis:"solo_last_ethereum_coin_distribution_processed_at,omitempty"`
	}
	ForT0LastEthereumCoinDistributionProcessedAtField struct {
		ForT0LastEthereumCoinDistributionProcessedAt *time.Time `redis:"for_t0_last_ethereum_coin_distribution_processed_at,omitempty"`
	}
	ForTMinus1LastEthereumCoinDistributionProcessedAtField struct {
		ForTMinus1LastEthereumCoinDistributionProcessedAt *time.Time `redis:"for_tminus1_last_ethereum_coin_distribution_processed_at,omitempty"`
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
	ReferralsCountChangeGuardUpdatedAtField struct {
		ReferralsCountChangeGuardUpdatedAt *time.Time `redis:"referrals_count_change_guard_updated_at,omitempty"`
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
		MiningBlockchainAccountAddress string `redis:"mining_blockchain_account_address" json:"miningBlockchainAccountAddress"`
	}
	BlockchainAccountAddressField struct {
		BlockchainAccountAddress string `redis:"blockchain_account_address"`
	}
	LatestDeviceField struct {
		LatestDevice string `redis:"latest_device,omitempty"`
	}
	CountryField struct {
		Country string `redis:"country" json:"country"`
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
	BalanceSoloEthereumField struct {
		BalanceSoloEthereum float64 `redis:"balance_solo_ethereum"`
	}
	BalanceT0Field struct {
		BalanceT0 float64 `redis:"balance_t0"`
	}
	BalanceT0EthereumField struct {
		BalanceT0Ethereum float64 `redis:"balance_t0_ethereum"`
	}
	BalanceT1Field struct {
		BalanceT1 float64 `redis:"balance_t1"`
	}
	BalanceT1EthereumField struct {
		BalanceT1Ethereum float64 `redis:"balance_t1_ethereum"`
	}
	BalanceT1EthereumPendingField struct {
		BalanceT1EthereumPending float64 `redis:"balance_t1_ethereum_pending"`
	}
	BalanceT2Field struct {
		BalanceT2 float64 `redis:"balance_t2"`
	}
	BalanceT2EthereumField struct {
		BalanceT2Ethereum float64 `redis:"balance_t2_ethereum"`
	}
	BalanceT2EthereumPendingField struct {
		BalanceT2EthereumPending float64 `redis:"balance_t2_ethereum_pending"`
	}
	BalanceForT0Field struct {
		BalanceForT0 float64 `redis:"balance_for_t0"`
	}
	BalanceForT0EthereumField struct {
		BalanceForT0Ethereum float64 `redis:"balance_for_t0_ethereum"`
	}
	BalanceForTMinus1Field struct {
		BalanceForTMinus1 float64 `redis:"balance_for_tminus1"`
	}
	BalanceForTMinus1EthereumField struct {
		BalanceForTMinus1Ethereum float64 `redis:"balance_for_tminus1_ethereum"`
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
		ID int64 `redis:"-" json:"-"`
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
	RecalculatedBalanceForTMinus1AtField struct {
		RecalculatedBalanceForTMinus1At *time.Time `redis:"recalculated_balance_for_tminus1_at,omitempty"`
	}
	RecalculatedBalanceT2AtField struct {
		RecalculatedBalanceT2At *time.Time `redis:"recalculated_balance_t2_at,omitempty"`
	}
	DeserializedRecalculatedUsersKey struct {
		ID int64 `redis:"-"`
	}

	DeserializedDryRunUsersKey struct {
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

func CalculateMiningStreak(now, start, end *time.Time, miningSessionDuration stdlibtime.Duration) uint64 {
	if start.IsNil() || end.IsNil() || now.After(*end.Time) || now.Before(*start.Time) {
		return 0
	}

	return uint64(now.Sub(*start.Time) / miningSessionDuration)
}

func (kyc *KYCState) KYCStepPassedCorrectly(kycStep users.KYCStep) bool {
	return (kyc.KYCStepBlocked == users.NoneKYCStep || kyc.KYCStepBlocked > kycStep) &&
		kycStep == kyc.KYCStepPassed &&
		kyc.KYCStepAttempted(kycStep)
}

func (kyc *KYCState) KYCStepNotAttempted(kycStep users.KYCStep) bool {
	return !kyc.KYCStepAttempted(kycStep)
}

func (kyc *KYCState) KYCStepAttempted(kycStep users.KYCStep) bool {
	return kyc.KYCStepsLastUpdatedAt != nil && len(*kyc.KYCStepsLastUpdatedAt) >= int(kycStep) && !(*kyc.KYCStepsLastUpdatedAt)[kycStep-1].IsNil()
}

func (kyc *KYCState) DelayPassedSinceLastKYCStepAttempt(kycStep users.KYCStep, duration stdlibtime.Duration) bool {
	return kyc.KYCStepAttempted(kycStep) && time.Now().Sub(*(*kyc.KYCStepsLastUpdatedAt)[kycStep-1].Time) >= duration
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
		if unmarshalledTime.IsNil() {
			unmarshalledTime = nil
		}
		timeSlice = append(timeSlice, unmarshalledTime)
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
		if unmarshalledTime.IsNil() {
			unmarshalledTime = nil
		}
		timeSlice = append(timeSlice, unmarshalledTime)
	}
	*t = timeSlice

	return nil
}
func (t *TimeSlice) MarshalJSON() ([]byte, error) {
	if t == nil || *t == nil {
		return []byte("null"), nil
	} else if len(*t) == 0 {
		return []byte("[]"), nil
	}
	timeSlice := *t
	data := make([]byte, 0, (len(timeSlice)*(len(stdlibtime.RFC3339Nano)+len(`""`)+1))+2)
	data = append(data, '[')
	for ix, elem := range timeSlice {
		marshalled, err := elem.MarshalJSON(context.Background())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshall %#v", elem)
		}
		data = append(data, marshalled...)
		if ix != len(timeSlice)-1 {
			data = append(data, ',')
		}
	}
	data = append(data, ']')

	return data, nil
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

func (k *DeserializedRecalculatedUsersKey) Key() string {
	if k == nil || k.ID == 0 {
		return ""
	}

	return SerializedRecalculatedUsersKey(k.ID)
}

func (k *DeserializedRecalculatedUsersKey) SetKey(val string) {
	if val == "" || val == "recalculated:" {
		return
	}
	if val[0] == 'r' {
		val = val[13:]
	}
	var err error
	k.ID, err = strconv.ParseInt(val, 10, 64)
	log.Panic(err)
}

func SerializedRecalculatedUsersKey(val any) string {
	switch typedVal := val.(type) {
	case string:
		if typedVal == "" {
			return ""
		}

		return "recalculated:" + typedVal
	case int64:
		if typedVal == 0 {
			return ""
		}

		return "recalculated:" + strconv.FormatInt(typedVal, 10)
	default:
		panic(fmt.Sprintf("%#v cannot be used as recalculated key", val))
	}
}

func (k *DeserializedDryRunUsersKey) Key() string {
	if k == nil || k.ID == 0 {
		return ""
	}

	return SerializedDryRunUsersKey(k.ID)
}

func (k *DeserializedDryRunUsersKey) SetKey(val string) {
	if val == "" || val == "dryrun:" {
		return
	}
	if val[0] == 'd' {
		val = val[7:]
	}
	var err error
	k.ID, err = strconv.ParseInt(val, 10, 64)
	log.Panic(err)
}

func SerializedDryRunUsersKey(val any) string {
	switch typedVal := val.(type) {
	case string:
		if typedVal == "" {
			return ""
		}

		return "dryrun:" + typedVal
	case int64:
		if typedVal == 0 {
			return ""
		}

		return "dryrun:" + strconv.FormatInt(typedVal, 10)
	default:
		panic(fmt.Sprintf("%#v cannot be used as dryrun key", val))
	}
}
