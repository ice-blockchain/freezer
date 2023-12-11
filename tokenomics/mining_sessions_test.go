// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"

	"github.com/ice-blockchain/freezer/model"
	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/time"
)

func TestRepositoryNewStartOrExtendMiningSession_CloseToMin(t *testing.T) { //nolint:funlen // .
	t.Parallel()
	var cfg Config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	repo := &repository{cfg: &cfg}
	actual := make([]*StartOrExtendMiningSession, 0, 27)
	startDates := make([]*time.Time, 0, 27)
	extensions := make([]stdlibtime.Duration, 0, 27)
	old := new(StartOrExtendMiningSession)
	now := time.Now()
	for ii := 0; ii < 27; ii++ {
		newMS, ext := repo.newStartOrExtendMiningSession(old, now)
		if newMS.MiningSessionSoloDayOffLastAwardedAt == nil {
			newMS.MiningSessionSoloDayOffLastAwardedAt = old.MiningSessionSoloDayOffLastAwardedAt
		} else if newMS.MiningSessionSoloDayOffLastAwardedAt.IsNil() {
			newMS.MiningSessionSoloDayOffLastAwardedAt = nil
		}
		if newMS.MiningSessionSoloStartedAt == nil {
			newMS.MiningSessionSoloStartedAt = old.MiningSessionSoloStartedAt
		}
		if newMS.MiningSessionSoloPreviouslyEndedAt == nil {
			newMS.MiningSessionSoloPreviouslyEndedAt = old.MiningSessionSoloPreviouslyEndedAt
		}
		startDates = append(startDates, now)
		actual = append(actual, newMS)
		extensions = append(extensions, ext)
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
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:     model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:         model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:           model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
		ReferralsCountChangeGuardUpdatedAtField: model.ReferralsCountChangeGuardUpdatedAtField{ReferralsCountChangeGuardUpdatedAt: startDates[0]},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min+repo.cfg.MiningSessionDuration.Max, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-1]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-2]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-3]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-4]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-5]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-6]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-7]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-8]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-9]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-10]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-11]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max+repo.cfg.MiningSessionDuration.Min, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-1]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, 1, extensions[ix])
	delta := 1 * stdlibtime.Millisecond
	assert.EqualValues(t, uint64(16), repo.calculateMiningStreak(time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(delta)), startDates[0], time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(delta)))) //nolint:lll // .
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:     model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:         model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[ix]},
		MiningSessionSoloEndedAtField:           model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
		MiningSessionSoloPreviouslyEndedAtField: model.MiningSessionSoloPreviouslyEndedAtField{MiningSessionSoloPreviouslyEndedAt: time.New(startDates[ix-1].Add(repo.cfg.MiningSessionDuration.Max))},
		ReferralsCountChangeGuardUpdatedAtField: model.ReferralsCountChangeGuardUpdatedAtField{ReferralsCountChangeGuardUpdatedAt: startDates[ix]},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max, extensions[ix])
}

func TestRepositoryNewStartOrExtendMiningSession_CloseToMax(t *testing.T) { //nolint:funlen,maintidx // .
	t.Parallel()
	var cfg Config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	repo := &repository{cfg: &cfg}
	actual := make([]*StartOrExtendMiningSession, 0, 27)
	startDates := make([]*time.Time, 0, 27)
	extensions := make([]stdlibtime.Duration, 0, 27)
	old := new(StartOrExtendMiningSession)
	now := time.Now()
	for ii := 0; ii < 27; ii++ {
		newMS, ext := repo.newStartOrExtendMiningSession(old, now)
		if newMS.MiningSessionSoloDayOffLastAwardedAt == nil {
			newMS.MiningSessionSoloDayOffLastAwardedAt = old.MiningSessionSoloDayOffLastAwardedAt
		} else if newMS.MiningSessionSoloDayOffLastAwardedAt.IsNil() {
			newMS.MiningSessionSoloDayOffLastAwardedAt = nil
		}
		if newMS.MiningSessionSoloStartedAt == nil {
			newMS.MiningSessionSoloStartedAt = old.MiningSessionSoloStartedAt
		}
		if newMS.MiningSessionSoloPreviouslyEndedAt == nil {
			newMS.MiningSessionSoloPreviouslyEndedAt = old.MiningSessionSoloPreviouslyEndedAt
		}
		startDates = append(startDates, now)
		actual = append(actual, newMS)
		extensions = append(extensions, ext)
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
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:     model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:         model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:           model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
		ReferralsCountChangeGuardUpdatedAtField: model.ReferralsCountChangeGuardUpdatedAtField{ReferralsCountChangeGuardUpdatedAt: startDates[0]},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField: model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:     model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloEndedAtField:       model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, 2*repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-1]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-2]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-3]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-4]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-5]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-6]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, 2*repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-1]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-2]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-3]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-4]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-5]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-6]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, 2*repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-1]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-2]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix-3]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))}, //nolint:lll // .
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max-1, extensions[ix])
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:       model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:           model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[0]},
		MiningSessionSoloDayOffLastAwardedAtField: model.MiningSessionSoloDayOffLastAwardedAtField{MiningSessionSoloDayOffLastAwardedAt: startDates[ix]},
		MiningSessionSoloEndedAtField:             model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
	}, actual[ix])
	assert.EqualValues(t, 2*repo.cfg.MiningSessionDuration.Max-1-1, extensions[ix])
	delta := 1 * stdlibtime.Millisecond
	assert.EqualValues(t, uint64(30), repo.calculateMiningStreak(time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(delta)), startDates[0], time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max).Add(delta)))) //nolint:lll // .
	ix++
	assert.EqualValues(t, &StartOrExtendMiningSession{
		MiningSessionSoloLastStartedAtField:     model.MiningSessionSoloLastStartedAtField{MiningSessionSoloLastStartedAt: startDates[ix]},
		MiningSessionSoloStartedAtField:         model.MiningSessionSoloStartedAtField{MiningSessionSoloStartedAt: startDates[ix]},
		MiningSessionSoloEndedAtField:           model.MiningSessionSoloEndedAtField{MiningSessionSoloEndedAt: time.New(startDates[ix].Add(repo.cfg.MiningSessionDuration.Max))},
		MiningSessionSoloPreviouslyEndedAtField: model.MiningSessionSoloPreviouslyEndedAtField{MiningSessionSoloPreviouslyEndedAt: time.New(startDates[ix-1].Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Max))},
		ReferralsCountChangeGuardUpdatedAtField: model.ReferralsCountChangeGuardUpdatedAtField{ReferralsCountChangeGuardUpdatedAt: startDates[ix]},
	}, actual[ix])
	assert.EqualValues(t, repo.cfg.MiningSessionDuration.Max, extensions[ix])
}
