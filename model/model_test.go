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
	serializedUser, err := json.Marshal(usr)
	require.NoError(t, err)

	type (
		KYCState struct {
			KYCStepsCreatedAtField
			KYCStepsLastUpdatedAtField
			KYCStepPassedField
			KYCStepBlockedField
		}
	)
	var deserializedUser KYCState
	err = json.Unmarshal(serializedUser, &deserializedUser)
	require.NoError(t, err)
	assert.EqualValues(t, stepA, deserializedUser.KYCStepBlocked)
	assert.EqualValues(t, stepB, deserializedUser.KYCStepPassed)
	assert.EqualValues(t, 1, len(*usr.KYCStepsCreatedAt))
	assert.EqualValues(t, (*usr.KYCStepsCreatedAt)[0], (*deserializedUser.KYCStepsCreatedAt)[0])
	assert.EqualValues(t, 2, len(*usr.KYCStepsLastUpdatedAt))
	assert.EqualValues(t, (*usr.KYCStepsLastUpdatedAt)[0], (*deserializedUser.KYCStepsLastUpdatedAt)[0])
	assert.EqualValues(t, (*usr.KYCStepsLastUpdatedAt)[1], (*deserializedUser.KYCStepsLastUpdatedAt)[1])
}
