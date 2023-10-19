// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"testing"
	stdlibtime "time"

	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/time"
	"github.com/stretchr/testify/assert"
)

func TestGetAdoptionsRange_1AdoptionPerRange(t *testing.T) {
	var adoptions []*tokenomics.Adoption[float64]
	adoptions = append(adoptions, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-58 * stdlibtime.Minute)),
		BaseMiningRate:   16.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-55 * stdlibtime.Minute)),
		BaseMiningRate:   8.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-30 * stdlibtime.Minute)),
		BaseMiningRate:   4.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-10 * stdlibtime.Minute)),
		BaseMiningRate:   2.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       nil,
		BaseMiningRate:   1.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       nil,
		BaseMiningRate:   0.5,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       nil,
		BaseMiningRate:   0.25,
		Milestone:        1,
		TotalActiveUsers: 1,
	})
	startedAt := time.New(time.Now().Add(-1 * stdlibtime.Minute))
	endedAt := time.Now()
	ranges := splitByAdoptionTimeRanges(adoptions, startedAt, endedAt)
	assert.Equal(t, 2., ranges[0].BaseMiningRate)
	assert.Equal(t, 2., ranges[1].BaseMiningRate)
}

func TestGetAdoptionsRange_2AdoptionsPerRange(t *testing.T) {
	var adoptions []*tokenomics.Adoption[float64]
	adoptions = append(adoptions, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-58 * stdlibtime.Minute)),
		BaseMiningRate:   16.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-55 * stdlibtime.Minute)),
		BaseMiningRate:   8.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-30 * stdlibtime.Minute)),
		BaseMiningRate:   4.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-10 * stdlibtime.Minute)),
		BaseMiningRate:   2.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-30 * stdlibtime.Second)),
		BaseMiningRate:   1.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       nil,
		BaseMiningRate:   0.5,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       nil,
		BaseMiningRate:   0.25,
		Milestone:        1,
		TotalActiveUsers: 1,
	})
	startedAt := time.New(time.Now().Add(-1 * stdlibtime.Minute))
	endedAt := time.Now()
	ranges := splitByAdoptionTimeRanges(adoptions, startedAt, endedAt)
	assert.Equal(t, 2., ranges[0].BaseMiningRate)
	assert.Equal(t, 1., ranges[1].BaseMiningRate)
	assert.Equal(t, 1., ranges[2].BaseMiningRate)
}

func TestGetAdoptionsRange_AdoptionDemotionPerRange(t *testing.T) {
	var adoptions []*tokenomics.Adoption[float64]
	adoptions = append(adoptions, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-30 * stdlibtime.Second)),
		BaseMiningRate:   16.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-55 * stdlibtime.Minute)),
		BaseMiningRate:   8.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-30 * stdlibtime.Minute)),
		BaseMiningRate:   4.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       time.New(time.Now().Add(-1 * stdlibtime.Minute)),
		BaseMiningRate:   2.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       nil,
		BaseMiningRate:   1.0,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       nil,
		BaseMiningRate:   0.5,
		Milestone:        1,
		TotalActiveUsers: 1,
	}, &tokenomics.Adoption[float64]{
		AchievedAt:       nil,
		BaseMiningRate:   0.25,
		Milestone:        1,
		TotalActiveUsers: 1,
	})
	startedAt := time.New(time.Now().Add(-2 * stdlibtime.Minute))
	endedAt := time.Now()
	ranges := splitByAdoptionTimeRanges(adoptions, startedAt, endedAt)
	assert.Equal(t, 4., ranges[0].BaseMiningRate)
	assert.Equal(t, 2., ranges[1].BaseMiningRate)
	assert.Equal(t, 16., ranges[2].BaseMiningRate)
	assert.Equal(t, 16., ranges[3].BaseMiningRate)
}
