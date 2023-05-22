// SPDX-License-Identifier: ice License 1.0

//go:build xxx

package tokenomics

import (
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"

	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/time"
)

func TestRepositoryNewStartOrExtendMiningSession_CloseToMin(t *testing.T) { //nolint:funlen // .
	t.Parallel()
	var cfg Config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	repo := &repository{cfg: &cfg}
	actual := make([]*startOrExtendMiningSession, 0, 27)
	startDates := make([]*time.Time, 0, 27)
	old := new(startOrExtendMiningSession)
	now := time.Now()
	for ii := 0; ii < 27; ii++ {
		newMS := repo.newStartOrExtendMiningSession(old, now)
		startDates = append(startDates, now)
		actual = append(actual, newMS)
		old = newMS
		delta := stdlibtime.Duration(0)
		if ii == 24 {
			delta = repo.cfg.MiningSessionDuration.Min + repo.cfg.MiningSessionDuration.Max + 1
		}
		if ii == 25 {
			delta = repo.cfg.MiningSessionDuration.Min + 1
		}
		now = time.New(now.Add(repo.cfg.MiningSessionDuration.Min).Add(delta))
	}
	assert.Len(t, actual, 27)
	ix := 0
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               0,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               0,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               1,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               1,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               2,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               2,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               3,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               3,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               4,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               4,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               5,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               5,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   6,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-1],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   6,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-2],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   7,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-3],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   7,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-4],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   8,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-5],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   8,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-6],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   9,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-7],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   9,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-8],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   10,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-9],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   10,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-10],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   11,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-11],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   11,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   12,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-1],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   14,
	}, actual[ix])
	delta := 1 * stdlibtime.Millisecond
	assert.EqualValues(t, uint64(16), repo.calculateMiningStreak(time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(delta)), startDates[0], time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(delta)))) //nolint:lll // .
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[ix],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               0,
	}, actual[ix])
}

func TestRepositoryNewstartOrExtendMiningSession_CloseToMax(t *testing.T) { //nolint:funlen,maintidx // .
	t.Parallel()
	var cfg Config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	repo := &repository{cfg: &cfg}
	actual := make([]*startOrExtendMiningSession, 0, 27)
	startDates := make([]*time.Time, 0, 27)
	old := new(startOrExtendMiningSession)
	now := time.Now()
	for ii := 0; ii < 27; ii++ {
		newMS := repo.newStartOrExtendMiningSession(old, now)
		startDates = append(startDates, now)
		actual = append(actual, newMS)
		old = newMS
		delta := stdlibtime.Duration(0)
		if ii == 24 {
			delta = 3*repo.cfg.MiningSessionDuration.Max - 1
		}
		if ii == 25 {
			delta = 2*repo.cfg.MiningSessionDuration.Max - 1
		}
		now = time.New(now.Add(repo.cfg.MiningSessionDuration.Max - 1).Add(delta))
	}
	assert.Len(t, actual, 27)
	ix := 0
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               0,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               0,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               1,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               2,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               3,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               4,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[0],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               5,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   6,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-1],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   7,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-2],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   8,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-3],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   9,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-4],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   10,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-5],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   11,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-6],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   12,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   13,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-1],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   14,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-2],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   15,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-3],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   16,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-4],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   17,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-5],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   18,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-6],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   19,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   20,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-1],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   21,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-2],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   22,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix-3],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)), //nolint:lll // .
		MiningStreak:                   23,
	}, actual[ix])
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt:     startDates[ix],
		LastMiningStartedAt:            startDates[0],
		LastFreeMiningSessionAwardedAt: startDates[ix],
		LastMiningEndedAt:              time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:                   27,
	}, actual[ix])
	delta := 1 * stdlibtime.Millisecond
	assert.EqualValues(t, uint64(30), repo.calculateMiningStreak(time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(delta)), startDates[0], time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(delta)))) //nolint:lll // .
	ix++
	assert.EqualValues(t, &startOrExtendMiningSession{
		LastNaturalMiningStartedAt: startDates[ix],
		LastMiningStartedAt:        startDates[ix],
		LastMiningEndedAt:          time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max)),
		MiningStreak:               0,
	}, actual[ix])
}
