// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"

	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/time"
)

func TestRepositoryCalculateMiningSession(t *testing.T) {
	t.Parallel()
	var cfg config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	repo := &repository{cfg: &cfg}

	now := time.Now()
	start := time.New(now.Add(-1 * stdlibtime.Second))
	end := time.New(now.Add(repo.cfg.MiningSessionDuration.Max).Add(-1 * stdlibtime.Second))
	actual := repo.calculateMiningSession(now, start, end)
	assert.EqualValues(t, start, actual.StartedAt)
	assert.False(t, *actual.Free)

	start = time.New(now.Add(-1 - repo.cfg.MiningSessionDuration.Min))
	end = time.New(now.Add(repo.cfg.MiningSessionDuration.Max).Add(-1 - repo.cfg.MiningSessionDuration.Min))
	actual = repo.calculateMiningSession(now, start, end)
	assert.EqualValues(t, start, actual.StartedAt)
	assert.False(t, *actual.Free)

	start = time.New(now.Add(-1 - repo.cfg.MiningSessionDuration.Max))
	end = time.New(now.Add(repo.cfg.MiningSessionDuration.Max).Add(repo.cfg.MiningSessionDuration.Min).Add(-1 - repo.cfg.MiningSessionDuration.Max))
	actual = repo.calculateMiningSession(now, start, end)
	assert.EqualValues(t, time.New(start.Add(repo.cfg.MiningSessionDuration.Max)), actual.StartedAt)
	assert.True(t, *actual.Free)
}
