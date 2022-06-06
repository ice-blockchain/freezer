package adoption

import "github.com/framey-io/go-tarantool"

// Public API.

type (
	Repository interface{}
)

// Private API.
type (
	adoptionSource struct {
		r *repository
	}
	repository struct {
		db tarantool.Connector
	}

	withCount struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack struct{} `msgpack:",asArray"`
		Count    uint64
	}

	adoptionHistory struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack         struct{} `msgpack:",asArray"`
		MinuteTimestamp  uint64
		HoursTimestamp   uint64
		TotalActiveUsers uint64
	}

	global struct {
		//nolint:unused // Because it is used by the msgpack library for marshalling/unmarshalling.
		_msgpack struct{} `msgpack:",asArray"`
		Key      string
		Value    uint64
	}
)

const (
	spaceAdoptionHistory = "ADOPTION_HISTORY"
	spaceGlobal          = "GLOBAL"
	fieldGlobalValue     = 1
	keyTotalActiveUsers  = "TOTAL_ACTIVE_USERS"
	hoursInDay           = 24
	secsInMinute         = 60
	minsInHour           = 60
	secsInHour           = secsInMinute * minsInHour
)
