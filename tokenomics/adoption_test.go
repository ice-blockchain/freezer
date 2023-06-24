// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"

	"github.com/ice-blockchain/wintr/time"
)

func TestDetectIncrTotalActiveUsersKeys(t *testing.T) {
	t.Parallel()
	var cfg Config
	cfg.AdoptionMilestoneSwitch.Duration = stdlibtime.Hour
	cfg.MiningSessionDuration.Min = 12 * stdlibtime.Hour
	cfg.MiningSessionDuration.Max = 24 * stdlibtime.Hour
	//cfg.AdoptionMilestoneSwitch.Duration = stdlibtime.Minute
	//cfg.MiningSessionDuration.Min =  30 * stdlibtime.Second
	//cfg.MiningSessionDuration.Max =  stdlibtime.Minute
	repo := &repository{cfg: &cfg}
	now := time.Now()
	ms := &MiningSession{
		LastNaturalMiningStartedAt: now,
		StartedAt:                  now,
		EndedAt:                    time.New(now.Add(cfg.MiningSessionDuration.Max)),
		PreviouslyEndedAt:          nil,
		Extension:                  cfg.MiningSessionDuration.Max,
	}
	actual := ms.detectIncrTotalActiveUsersKeys(repo)
	assert.EqualValues(t, []string{
		repo.totalActiveUsersKey(now.Add(0 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(1 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(2 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(3 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(4 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(5 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(6 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(7 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(8 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(9 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(10 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(11 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(12 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(13 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(14 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(15 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(16 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(17 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(18 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(19 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(20 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(21 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(22 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(23 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(24 * cfg.AdoptionMilestoneSwitch.Duration)),
	}, actual)
	// --
	ms = &MiningSession{
		LastNaturalMiningStartedAt: time.New(now.Add(cfg.MiningSessionDuration.Min)),
		StartedAt:                  now,
		EndedAt:                    time.New(now.Add(cfg.MiningSessionDuration.Min).Add(cfg.MiningSessionDuration.Max)),
		PreviouslyEndedAt:          nil,
		Extension:                  cfg.MiningSessionDuration.Min,
	}
	actual = ms.detectIncrTotalActiveUsersKeys(repo)
	assert.EqualValues(t, []string{
		repo.totalActiveUsersKey(now.Add(25 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(26 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(27 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(28 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(29 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(30 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(31 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(32 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(33 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(34 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(35 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(36 * cfg.AdoptionMilestoneSwitch.Duration)),
	}, actual)
	// --
	ms = &MiningSession{
		LastNaturalMiningStartedAt: time.New(now.Add(cfg.MiningSessionDuration.Min).Add(cfg.MiningSessionDuration.Max + 1)),
		StartedAt:                  time.New(now.Add(cfg.MiningSessionDuration.Min).Add(cfg.MiningSessionDuration.Max + 1)),
		EndedAt:                    time.New(now.Add(cfg.MiningSessionDuration.Min).Add(cfg.MiningSessionDuration.Max + 1).Add(cfg.MiningSessionDuration.Max)),
		PreviouslyEndedAt:          time.New(now.Add(cfg.MiningSessionDuration.Min).Add(cfg.MiningSessionDuration.Max)),
		Extension:                  cfg.MiningSessionDuration.Max,
	}
	actual = ms.detectIncrTotalActiveUsersKeys(repo)
	assert.EqualValues(t, []string{
		repo.totalActiveUsersKey(now.Add(37 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(38 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(39 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(40 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(41 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(42 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(43 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(44 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(45 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(46 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(47 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(48 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(49 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(50 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(51 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(52 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(53 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(54 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(55 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(56 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(57 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(58 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(59 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(60 * cfg.AdoptionMilestoneSwitch.Duration)),
	}, actual)
	offset := 80 * stdlibtime.Minute
	ms = &MiningSession{
		LastNaturalMiningStartedAt: time.New(now.Add(cfg.MiningSessionDuration.Min).Add(2*cfg.MiningSessionDuration.Max + 1).Add(offset)),
		StartedAt:                  time.New(now.Add(cfg.MiningSessionDuration.Min).Add(cfg.MiningSessionDuration.Max + 1)),
		EndedAt:                    time.New(now.Add(2 * cfg.MiningSessionDuration.Min).Add(2*cfg.MiningSessionDuration.Max + 1).Add(offset)),
		PreviouslyEndedAt:          time.New(now.Add(cfg.MiningSessionDuration.Min).Add(cfg.MiningSessionDuration.Max)),
		Extension:                  cfg.MiningSessionDuration.Min + offset,
	}
	actual = ms.detectIncrTotalActiveUsersKeys(repo)
	assert.EqualValues(t, []string{
		repo.totalActiveUsersKey(now.Add(61 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(62 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(63 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(64 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(65 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(66 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(67 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(68 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(69 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(70 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(71 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(72 * cfg.AdoptionMilestoneSwitch.Duration)),
		repo.totalActiveUsersKey(now.Add(73 * cfg.AdoptionMilestoneSwitch.Duration)),
	}, actual)
}
