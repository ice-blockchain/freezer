// SPDX-License-Identifier: ice License 1.0

package model

import (
	"context"
	"testing"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/time"
)

func TestKYCFields(t *testing.T) { //nolint:funlen // .
	t.Parallel()
	type (
		dummy struct {
			KYCStepsLastUpdatedAtField
			KYCStepPassedField
		}
	)
	value := &dummy{}
	resp := storage.SerializeValue(value)
	assert.EqualValues(t, []any{
		"kyc_steps_last_updated_at",
		"",
		"kyc_step_passed",
		"0",
	}, resp)

	value = &dummy{
		KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{
			time.New(stdlibtime.Date(1111, 11, 1, 1, 1, 1, 1, stdlibtime.UTC)),
			time.New(stdlibtime.Date(2222, 12, 2, 2, 2, 2, 2, stdlibtime.UTC)),
		}},
		KYCStepPassedField: KYCStepPassedField{KYCStepPassed: users.Social2KYCStep},
	}
	resp = storage.SerializeValue(value)
	assert.EqualValues(t, []any{
		"kyc_steps_last_updated_at",
		"1111-11-01T01:01:01.000000001Z,2222-12-02T02:02:02.000000002Z",
		"kyc_step_passed",
		"5",
	}, resp)
	args := append(make([]interface{}, 0, 2+len(resp)/2), "hmget", "boguskey")
	vals := make([]any, 0, len(resp)/2)
	for i, field := range resp {
		if i%2 == 0 {
			args = append(args, field)
		} else {
			vals = append(vals, field)
		}
	}
	cmd := redis.NewSliceCmd(context.Background(), args...)
	cmd.SetVal(vals)
	var res dummy
	assert.NoError(t, storage.DeserializeValue(&res, cmd.Scan))
	assert.EqualValues(t, value, &res)

	cmd = redis.NewSliceCmd(context.Background(), args...)
	cmd.SetVal([]any{
		"",
		"0",
	})
	var res2 dummy
	assert.NoError(t, storage.DeserializeValue(&res2, cmd.Scan))
	assert.EqualValues(t, &dummy{KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: new(TimeSlice)}}, &res2)

	cmd = redis.NewSliceCmd(context.Background(), "hmget", "boguskey", "bogus")
	cmd.SetVal([]any{
		"bogus value",
	})
	var res3 dummy
	assert.NoError(t, storage.DeserializeValue(&res3, cmd.Scan))
	assert.EqualValues(t, &dummy{}, &res3)
}

func TestTimeSliceEquals(t *testing.T) {
	t.Parallel()
	now := time.Now()
	assert.True(t, new(TimeSlice).Equals(new(TimeSlice)))
	assert.True(t, (*TimeSlice)(nil).Equals(nil))
	assert.True(t, (&TimeSlice{now}).Equals(&TimeSlice{now}))
	assert.True(t, (&TimeSlice{now, now}).Equals(&TimeSlice{now, now}))
	assert.True(t, (&TimeSlice{nil, now}).Equals(&TimeSlice{nil, now}))
	assert.True(t, (&TimeSlice{nil}).Equals(&TimeSlice{nil}))
	assert.True(t, (&TimeSlice{new(time.Time)}).Equals(&TimeSlice{new(time.Time)}))
	assert.True(t, (&TimeSlice{new(time.Time), now}).Equals(&TimeSlice{new(time.Time), now}))
	assert.True(t, (&TimeSlice{new(time.Time), now}).Equals(&TimeSlice{nil, now}))

	assert.False(t, new(TimeSlice).Equals(nil))
	assert.False(t, new(TimeSlice).Equals(&TimeSlice{now}))
	assert.False(t, new(TimeSlice).Equals(&TimeSlice{now, now}))
	assert.False(t, (*TimeSlice)(nil).Equals(new(TimeSlice)))
	assert.False(t, (*TimeSlice)(nil).Equals(&TimeSlice{now}))
	assert.False(t, (*TimeSlice)(nil).Equals(&TimeSlice{now, now}))
	assert.False(t, (&TimeSlice{now}).Equals(&TimeSlice{nil}))
	assert.False(t, (&TimeSlice{now}).Equals(&TimeSlice{now, now}))
	assert.False(t, (&TimeSlice{now, now}).Equals(&TimeSlice{now}))
	assert.False(t, (&TimeSlice{nil, now}).Equals(&TimeSlice{now, nil}))
	assert.False(t, (&TimeSlice{nil, nil}).Equals(&TimeSlice{now, now}))
}

func TestEskimoToFreezerKYCStateDeserialization(t *testing.T) {
	t.Parallel()
	stepA := users.LivenessDetectionKYCStep
	stepB := users.Social2KYCStep
	usr := &users.User{
		KYCStepsCreatedAt:     &[]*time.Time{time.Now()},
		KYCStepsLastUpdatedAt: &[]*time.Time{time.Now(), time.Now()},
		KYCStepPassed:         &stepB,
		KYCStepBlocked:        &stepA,
	}
	usr.Country = "XX"
	usr.MiningBlockchainAccountAddress = "YY"
	serializedUser, err := json.Marshal(usr)
	require.NoError(t, err)

	var deserializedUser struct {
		CountryField
		MiningBlockchainAccountAddressField
		KYCState
		DeserializedUsersKey
	}
	err = json.Unmarshal(serializedUser, &deserializedUser)
	require.NoError(t, err)
	assert.EqualValues(t, "XX", deserializedUser.Country)
	assert.EqualValues(t, "YY", deserializedUser.MiningBlockchainAccountAddress)
	assert.EqualValues(t, stepA, deserializedUser.KYCStepBlocked)
	assert.EqualValues(t, stepB, deserializedUser.KYCStepPassed)
	assert.EqualValues(t, 1, len(*usr.KYCStepsCreatedAt))
	assert.EqualValues(t, (*usr.KYCStepsCreatedAt)[0], (*deserializedUser.KYCStepsCreatedAt)[0])
	assert.EqualValues(t, 2, len(*usr.KYCStepsLastUpdatedAt))
	assert.EqualValues(t, (*usr.KYCStepsLastUpdatedAt)[0], (*deserializedUser.KYCStepsLastUpdatedAt)[0])
	assert.EqualValues(t, (*usr.KYCStepsLastUpdatedAt)[1], (*deserializedUser.KYCStepsLastUpdatedAt)[1])
}

//nolint:funlen // .
func TestEskimoToFreezerKYCStateDeserialization_WithNullsInsideSlice(t *testing.T) {
	t.Parallel()
	stepA := users.LivenessDetectionKYCStep
	stepB := users.Social2KYCStep
	t1 := time.New(stdlibtime.Date(2221, 11, 11, 11, 11, 11, 11, stdlibtime.UTC))
	t2 := time.New(stdlibtime.Date(2222, 11, 11, 11, 11, 11, 11, stdlibtime.UTC))
	t3 := time.New(stdlibtime.Date(2223, 11, 11, 11, 11, 11, 11, stdlibtime.UTC))
	usr := &users.User{
		KYCStepsCreatedAt:     &[]*time.Time{new(time.Time), t1, nil, new(time.Time), time.New(stdlibtime.Unix(0, 0)), new(time.Time)},
		KYCStepsLastUpdatedAt: &[]*time.Time{nil, t2, nil, t3, {Time: new(stdlibtime.Time)}},
		KYCStepPassed:         &stepB,
		KYCStepBlocked:        &stepA,
	}
	serializedUser, err := json.MarshalContext(context.Background(), usr)
	require.NoError(t, err)
	assert.Equal(t, `{"kycStepsLastUpdatedAt":[null,"2222-11-11T11:11:11.000000011Z",null,"2223-11-11T11:11:11.000000011Z","0001-01-01T00:00:00Z"],"kycStepsCreatedAt":[null,"2221-11-11T11:11:11.000000011Z",null,null,null,null],"kycStepPassed":5,"kycStepBlocked":2}`, string(serializedUser)) //nolint:lll // .

	var deserializedEskimoUser users.User
	err = json.UnmarshalContext(context.Background(), serializedUser, &deserializedEskimoUser)
	require.NoError(t, err)
	assert.EqualValues(t, users.User{
		KYCStepsCreatedAt:     &[]*time.Time{nil, t1, nil, nil, nil, nil},
		KYCStepsLastUpdatedAt: &[]*time.Time{nil, t2, nil, t3, {Time: new(stdlibtime.Time)}},
		KYCStepPassed:         &stepB,
		KYCStepBlocked:        &stepA,
	}, deserializedEskimoUser)

	type (
		KYCState struct {
			KYCStepsCreatedAtField
			KYCStepsLastUpdatedAtField
			KYCStepPassedField
			KYCStepBlockedField
			KYCQuizCompletedField
			KYCQuizDisabledField
		}
	)
	var deserializedUser KYCState
	err = json.Unmarshal(serializedUser, &deserializedUser)
	require.NoError(t, err)
	assert.EqualValues(t, KYCState{
		KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{nil, t1, nil, nil, nil, nil}}, //nolint:lll // .
		KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{nil, t2, nil, t3, {Time: new(stdlibtime.Time)}}},
		KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: stepB},
		KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: stepA},
		KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: false},
		KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: false},
	}, deserializedUser)
	serializedUser2, err := json.Marshal(deserializedUser)
	require.NoError(t, err)
	assert.Equal(t, `{"kycStepsCreatedAt":[null,"2221-11-11T11:11:11.000000011Z",null,null,null,null],"kycStepsLastUpdatedAt":[null,"2222-11-11T11:11:11.000000011Z",null,"2223-11-11T11:11:11.000000011Z","0001-01-01T00:00:00Z"],"kycStepPassed":5,"kycStepBlocked":2,"kycQuizCompleted":false,"kycQuizDisabled":false}`, string(serializedUser2)) //nolint:lll // .

	serializedKYCStepsCreatedAt, err := deserializedUser.KYCStepsCreatedAt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, ",2221-11-11T11:11:11.000000011Z,,,,", string(serializedKYCStepsCreatedAt))
	serializedKYCStepsLastUpdatedAtField, err := deserializedUser.KYCStepsLastUpdatedAt.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, ",2222-11-11T11:11:11.000000011Z,,2223-11-11T11:11:11.000000011Z,0001-01-01T00:00:00Z", string(serializedKYCStepsLastUpdatedAtField))

	deserializedTimeSlice := new(TimeSlice)
	require.NoError(t, deserializedTimeSlice.UnmarshalBinary(serializedKYCStepsLastUpdatedAtField))
	assert.EqualValues(t, &TimeSlice{nil, t2, nil, t3, {Time: new(stdlibtime.Time)}}, deserializedTimeSlice)

	err = json.Unmarshal(serializedUser, &deserializedUser)
	require.NoError(t, err)

	assert.EqualValues(t, KYCState{
		KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{nil, t1, nil, nil, nil, nil}}, //nolint:lll // .
		KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{nil, t2, nil, t3, {Time: new(stdlibtime.Time)}}},
		KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: stepB},
		KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: stepA},
		KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: false},
		KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: false},
	}, deserializedUser)
	deserializedUser.KYCQuizCompleted = true
	deserializedUser.KYCQuizDisabled = true
	serializedUser3, err := json.Marshal(deserializedUser)
	require.NoError(t, err)
	assert.Equal(t, `{"kycStepsCreatedAt":[null,"2221-11-11T11:11:11.000000011Z",null,null,null,null],"kycStepsLastUpdatedAt":[null,"2222-11-11T11:11:11.000000011Z",null,"2223-11-11T11:11:11.000000011Z","0001-01-01T00:00:00Z"],"kycStepPassed":5,"kycStepBlocked":2,"kycQuizCompleted":true,"kycQuizDisabled":true}`, string(serializedUser3)) //nolint:lll // .

}

//nolint:funlen // .
func TestEskimoToFreezerKYCStateDeserialization_WithEmptySlices(t *testing.T) {
	t.Parallel()
	var nilSlice []*time.Time
	usr := &users.User{
		KYCStepsCreatedAt:     &([]*time.Time{}),
		KYCStepsLastUpdatedAt: &nilSlice,
	}
	serializedUser, err := json.MarshalContext(context.Background(), usr)
	require.NoError(t, err)
	assert.Equal(t, `{"kycStepsLastUpdatedAt":null,"kycStepsCreatedAt":[]}`, string(serializedUser))

	var deserializedEskimoUser users.User
	err = json.UnmarshalContext(context.Background(), serializedUser, &deserializedEskimoUser)
	require.NoError(t, err)
	assert.EqualValues(t, users.User{
		KYCStepsCreatedAt: &([]*time.Time{}),
	}, deserializedEskimoUser)

	serializedUserTmp, err := json.MarshalContext(context.Background(), new(users.User))
	require.NoError(t, err)
	assert.Equal(t, `{}`, string(serializedUserTmp))

	var deserializedEskimoUserTmp users.User
	err = json.UnmarshalContext(context.Background(), serializedUserTmp, &deserializedEskimoUserTmp)
	require.NoError(t, err)
	assert.EqualValues(t, users.User{}, deserializedEskimoUserTmp)

	type (
		KYCState struct {
			KYCStepsCreatedAtField
			KYCStepsLastUpdatedAtField
			KYCStepPassedField
			KYCStepBlockedField
			KYCQuizCompletedField
			KYCQuizDisabledField
		}
	)
	var deserializedUser KYCState
	err = json.Unmarshal(serializedUser, &deserializedUser)
	require.NoError(t, err)
	assert.EqualValues(t, KYCState{
		KYCStepsCreatedAtField: KYCStepsCreatedAtField{KYCStepsCreatedAt: new(TimeSlice)},
	}, deserializedUser)
	serializedUser2, err := json.Marshal(deserializedUser)
	require.NoError(t, err)
	assert.Equal(t, `{"kycStepsCreatedAt":null,"kycStepsLastUpdatedAt":null,"kycStepPassed":0,"kycStepBlocked":0,"kycQuizCompleted":false,"kycQuizDisabled":false}`, string(serializedUser2))

	var deserializedUserTmp KYCState
	err = json.Unmarshal(serializedUserTmp, &deserializedUserTmp)
	require.NoError(t, err)
	assert.EqualValues(t, KYCState{}, deserializedUserTmp)
	serializedUserTmp2, err := json.Marshal(deserializedUserTmp)
	require.NoError(t, err)
	assert.Equal(t, `{"kycStepsCreatedAt":null,"kycStepsLastUpdatedAt":null,"kycStepPassed":0,"kycStepBlocked":0,"kycQuizCompleted":false,"kycQuizDisabled":false}`, string(serializedUserTmp2))

	var nilTs *TimeSlice
	t1, err := nilTs.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(t1))
	t2, err := new(TimeSlice).MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(t2))
	t3, err := (&TimeSlice{}).MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(t3))
	t4, err := (&TimeSlice{nil}).MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(t4))
	t5, err := (&TimeSlice{new(time.Time)}).MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(t5))

	deserializedTimeSliceA := new(TimeSlice)
	require.NoError(t, deserializedTimeSliceA.UnmarshalJSON([]byte("")))
	assert.EqualValues(t, new(TimeSlice), deserializedTimeSliceA)
	deserializedTimeSliceB := new(TimeSlice)
	require.NoError(t, deserializedTimeSliceB.UnmarshalJSON([]byte("[]")))
	assert.EqualValues(t, new(TimeSlice), deserializedTimeSliceB)
	deserializedTimeSliceC := new(TimeSlice)
	require.NoError(t, deserializedTimeSliceC.UnmarshalJSON(nil))
	assert.EqualValues(t, new(TimeSlice), deserializedTimeSliceC)

	var deserializedTimeSliceD *TimeSlice
	require.NoError(t, deserializedTimeSliceD.UnmarshalJSON([]byte("")))
	assert.EqualValues(t, (*TimeSlice)(nil), deserializedTimeSliceD)
	var deserializedTimeSliceE *TimeSlice
	require.NoError(t, deserializedTimeSliceE.UnmarshalJSON([]byte("[]")))
	assert.EqualValues(t, (*TimeSlice)(nil), deserializedTimeSliceE)
	var deserializedTimeSliceF *TimeSlice
	require.NoError(t, deserializedTimeSliceF.UnmarshalJSON(nil))
	assert.EqualValues(t, (*TimeSlice)(nil), deserializedTimeSliceF)

	deserializedTimeSlice := new(TimeSlice)
	require.NoError(t, deserializedTimeSlice.UnmarshalBinary([]byte("")))
	assert.EqualValues(t, new(TimeSlice), deserializedTimeSlice)
	deserializedTimeSlice2 := new(TimeSlice)
	require.NoError(t, deserializedTimeSlice2.UnmarshalBinary(nil))
	assert.EqualValues(t, new(TimeSlice), deserializedTimeSlice2)
	var deserializedTimeSlice3 *TimeSlice
	require.NoError(t, deserializedTimeSlice3.UnmarshalBinary([]byte("")))
	assert.EqualValues(t, (*TimeSlice)(nil), deserializedTimeSlice3)
	var deserializedTimeSlice4 *TimeSlice
	require.Panics(t, func() {
		_ = deserializedTimeSlice4.UnmarshalBinary([]byte("2221-11-11T11:11:11.000000011Z"))
	})
	assert.EqualValues(t, (*TimeSlice)(nil), deserializedTimeSlice4)
}

func TestKYCStepPassedCorrectly(t *testing.T) {
	t.Parallel()

	t.Run("blocked = None, step = quiz(4), quizCompleted = true && !quizDisabled, stepPassed >= kycStep, arg = QuizKYCStep", func(t *testing.T) {
		t.Parallel()
		kycState := KYCState{
			KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: false},
			KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
			KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
		}
		assert.Equal(t, true, kycState.KYCStepPassedCorrectly(users.QuizKYCStep))
	})

	t.Run("blocked = FacialRecognition, step = quiz(4), quizCompleted = true && !quizDisabled, stepPassed >= kycStep, arg = QuizKYCStep", func(t *testing.T) {
		t.Parallel()
		kycState := KYCState{
			KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: users.FacialRecognitionKYCStep},
			KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: false},
			KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
			KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
		}
		assert.Equal(t, false, kycState.KYCStepPassedCorrectly(users.QuizKYCStep))
	})

	t.Run("blocked = None, step = quiz(4), quizCompleted = false && !quizDisabled, stepPassed >= kycStep, arg = QuizKYCStep", func(t *testing.T) {
		t.Parallel()
		kycState := KYCState{
			KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: false},
			KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: false},
			KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
			KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
		}
		assert.Equal(t, false, kycState.KYCStepPassedCorrectly(users.QuizKYCStep))
	})

	t.Run("blocked = None, step = quiz(4), quizCompleted = true && quizDisabled, stepPassed >= kycStep, arg = QuizKYCStep", func(t *testing.T) {
		t.Parallel()
		kycState := KYCState{
			KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: true},
			KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
			KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
		}
		assert.Equal(t, false, kycState.KYCStepPassedCorrectly(users.QuizKYCStep))
	})

	t.Run("blocked = None, step = Social1KYCStep(3), quizCompleted = true && !quizDisabled, stepPassed >= kycStep, arg = Social1KYCStep", func(t *testing.T) {
		t.Parallel()
		kycState := KYCState{
			KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: users.Social1KYCStep},
			KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: true},
			KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
			KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
		}
		assert.Equal(t, true, kycState.KYCStepPassedCorrectly(users.Social1KYCStep))
	})

	t.Run("blocked = None, step = Social1KYCStep(3), quizCompleted = false && quizDisabled, stepPassed >= kycStep, arg = Social1KYCStep", func(t *testing.T) {
		t.Parallel()
		kycState := KYCState{
			KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: false},
			KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: users.QuizKYCStep},
			KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: true},
			KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
			KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
		}
		assert.Equal(t, true, kycState.KYCStepPassedCorrectly(users.Social1KYCStep))
	})

	t.Run("blocked = None, step = Social2KYCStep(5), quizCompleted = true && !quizDisabled, stepPassed >= kycStep, arg = QuizKYCStep", func(t *testing.T) {
		t.Parallel()
		kycState := KYCState{
			KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: true},
			KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: users.Social2KYCStep},
			KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: false},
			KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
			KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), time.Now()}},
		}
		assert.Equal(t, true, kycState.KYCStepPassedCorrectly(users.QuizKYCStep))
	})

	t.Run("blocked = None, step = Social2KYCStep(5), quizCompleted = false && !quizDisabled, stepPassed >= kycStep, arg = Social2KYCStep", func(t *testing.T) {
		t.Parallel()
		kycState := KYCState{
			KYCQuizCompletedField:      KYCQuizCompletedField{KYCQuizCompleted: false},
			KYCStepPassedField:         KYCStepPassedField{KYCStepPassed: users.Social2KYCStep},
			KYCStepBlockedField:        KYCStepBlockedField{KYCStepBlocked: users.NoneKYCStep},
			KYCQuizDisabledField:       KYCQuizDisabledField{KYCQuizDisabled: false},
			KYCStepsCreatedAtField:     KYCStepsCreatedAtField{KYCStepsCreatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), nil, time.Now()}},
			KYCStepsLastUpdatedAtField: KYCStepsLastUpdatedAtField{KYCStepsLastUpdatedAt: &TimeSlice{time.Now(), time.Now(), time.Now(), nil, time.Now()}},
		}
		assert.Equal(t, true, kycState.KYCStepPassedCorrectly(users.Social2KYCStep))
	})
}
