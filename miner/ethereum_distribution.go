// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"sync"
	stdlibtime "time"

	"github.com/pkg/errors"

	coindistribution "github.com/ice-blockchain/freezer/coin-distribution"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

const (
	ethereumDistributionDryRunModeEnabled = true
)

func (ref *referral) username() string {
	if ref != nil && ref.Username != "" {
		return ref.Username
	}

	return "icenetwork/bogus"
}

func (ref *referral) isEligibleForSelfForEthereumDistribution(now *time.Time) bool {
	coinDistributionCollectorSettings := cfg.coinDistributionCollectorSettings.Load()

	return ref != nil &&
		ref.ID != 0 &&
		coindistribution.IsEligibleForEthereumDistributionNow(
			ref.ID,
			now,
			ref.SoloLastEthereumCoinDistributionProcessedAt,
			cfg.coinDistributionCollectorSettings.Load().StartDate,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max) &&
		coindistribution.IsEligibleForEthereumDistribution(
			coinDistributionCollectorSettings.MinMiningStreaksRequired,
			ref.BalanceTotalStandard-ref.BalanceSoloEthereum-ref.BalanceT0Ethereum-ref.BalanceT1Ethereum-ref.BalanceT2Ethereum,
			coinDistributionCollectorSettings.MinBalanceRequired,
			ref.MiningBlockchainAccountAddress,
			ref.Country,
			coinDistributionCollectorSettings.DeniedCountries,
			now,
			ref.MiningSessionSoloStartedAt,
			ref.MiningSessionSoloEndedAt,
			coinDistributionCollectorSettings.EndDate,
			ref.KYCState,
			cfg.MiningSessionDuration.Max,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max)
}

func (ref *referral) isEligibleForReferralForEthereumDistribution(now *time.Time) bool {
	coinDistributionCollectorSettings := cfg.coinDistributionCollectorSettings.Load()

	return ref != nil &&
		ref.ID != 0 &&
		coindistribution.IsEligibleForEthereumDistribution(
			0,
			0,
			0,
			"skip",
			ref.Country,
			coinDistributionCollectorSettings.DeniedCountries,
			now,
			ref.MiningSessionSoloStartedAt,
			ref.MiningSessionSoloEndedAt,
			coinDistributionCollectorSettings.EndDate,
			ref.KYCState,
			cfg.MiningSessionDuration.Max,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max)
}

func (u *user) isEligibleForSelfForEthereumDistribution(now *time.Time) bool {
	coinDistributionCollectorSettings := cfg.coinDistributionCollectorSettings.Load()

	return u != nil &&
		u.ID != 0 &&
		coindistribution.IsEligibleForEthereumDistributionNow(
			u.ID,
			now,
			u.SoloLastEthereumCoinDistributionProcessedAt,
			cfg.coinDistributionCollectorSettings.Load().StartDate,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max) &&
		coindistribution.IsEligibleForEthereumDistribution(
			coinDistributionCollectorSettings.MinMiningStreaksRequired,
			u.BalanceTotalStandard-u.BalanceSoloEthereum-u.BalanceT0Ethereum-u.BalanceT1Ethereum-u.BalanceT2Ethereum,
			coinDistributionCollectorSettings.MinBalanceRequired,
			u.MiningBlockchainAccountAddress,
			u.Country,
			coinDistributionCollectorSettings.DeniedCountries,
			now,
			u.MiningSessionSoloStartedAt,
			u.MiningSessionSoloEndedAt,
			coinDistributionCollectorSettings.EndDate,
			u.KYCState,
			cfg.MiningSessionDuration.Max,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max)
}

func (u *user) isEligibleForT0ForEthereumDistribution(now *time.Time) bool {
	return u != nil &&
		u.ID != 0 &&
		coindistribution.IsEligibleForEthereumDistributionNow(
			u.ID,
			now,
			u.ForT0LastEthereumCoinDistributionProcessedAt,
			cfg.coinDistributionCollectorSettings.Load().StartDate,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max) &&
		u.isEligibleForReferralForEthereumDistribution(now)
}

func (u *user) isEligibleForTMinus1ForEthereumDistribution(now *time.Time) bool {
	return u != nil &&
		u.ID != 0 &&
		coindistribution.IsEligibleForEthereumDistributionNow(
			u.ID,
			now,
			u.ForTMinus1LastEthereumCoinDistributionProcessedAt,
			cfg.coinDistributionCollectorSettings.Load().StartDate,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max) &&
		u.isEligibleForReferralForEthereumDistribution(now)
}

func (u *user) isEligibleForReferralForEthereumDistribution(now *time.Time) bool {
	coinDistributionCollectorSettings := cfg.coinDistributionCollectorSettings.Load()
	return coindistribution.IsEligibleForEthereumDistribution(
		0,
		0,
		0,
		"skip",
		u.Country,
		coinDistributionCollectorSettings.DeniedCountries,
		now,
		u.MiningSessionSoloStartedAt,
		u.MiningSessionSoloEndedAt,
		coinDistributionCollectorSettings.EndDate,
		u.KYCState,
		cfg.MiningSessionDuration.Max,
		cfg.EthereumDistributionFrequency.Min,
		cfg.EthereumDistributionFrequency.Max)
}

func (u *user) couldHaveBeenEligibleForEthereumDistributionRecently(now *time.Time) bool {
	return u != nil && !u.MiningSessionSoloEndedAt.IsNil() && u.MiningSessionSoloEndedAt.After(now.Add(-(cfg.MiningSessionDuration.Max / 8)))
}

func (ref *referral) couldHaveBeenEligibleForEthereumDistributionRecently(now *time.Time) bool {
	return ref != nil && !ref.MiningSessionSoloEndedAt.IsNil() && ref.MiningSessionSoloEndedAt.After(now.Add(-(cfg.MiningSessionDuration.Max / 8)))
}

//nolint:funlen // .
func (u *user) processEthereumCoinDistribution(
	now *time.Time, t0, tMinus1 *referral,
) (records []*coindistribution.ByEarnerForReview, balanceDistributedForT0, balanceDistributedForTMinus1 float64) {
	if !isCoinDistributionCollectorEnabled(now) {
		if u.BalanceT1EthereumPending != nil {
			u.BalanceT1Ethereum += float64(*u.BalanceT1EthereumPending)
		}
		if u.BalanceT2EthereumPending != nil {
			u.BalanceT2Ethereum += float64(*u.BalanceT2EthereumPending)
		}
		u.BalanceT1EthereumPending = new(model.FlexibleFloat64)
		u.BalanceT2EthereumPending = new(model.FlexibleFloat64)
		u.SoloLastEthereumCoinDistributionProcessedAt = nil
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
		u.ForTMinus1LastEthereumCoinDistributionProcessedAt = nil

		return nil, 0, 0
	}
	u.BalanceT1EthereumPending = nil
	u.BalanceT2EthereumPending = nil
	records = make([]*coindistribution.ByEarnerForReview, 0, 1+1+1+1)
	var (
		t0CD         *coindistribution.ByEarnerForReview
		forT0CD      *coindistribution.ByEarnerForReview
		forTMinus1CD *coindistribution.ByEarnerForReview
		soloCD       *coindistribution.ByEarnerForReview
	)
	if u.couldHaveBeenEligibleForEthereumDistributionRecently(now) {
		soloCD = &coindistribution.ByEarnerForReview{
			CreatedAt:          now,
			Username:           u.Username,
			ReferredByUsername: t0.username(),
			UserID:             u.UserID,
			EarnerUserID:       u.UserID,
			EthAddress:         u.MiningBlockchainAccountAddress,
			InternalID:         u.ID,
			Balance:            0,
		}
		records = append(records, soloCD)
	}
	if u.couldHaveBeenEligibleForEthereumDistributionRecently(now) && t0.couldHaveBeenEligibleForEthereumDistributionRecently(now) && t0 != nil && t0.UserID != u.UserID && (tMinus1 == nil || (tMinus1.UserID != u.UserID && tMinus1.UserID != t0.UserID)) { //nolint:lll // .
		t0CD = &coindistribution.ByEarnerForReview{
			CreatedAt:    now,
			UserID:       u.UserID,
			EarnerUserID: t0.UserID,
			Balance:      0,
		}
		forT0CD = &coindistribution.ByEarnerForReview{
			CreatedAt:    now,
			UserID:       t0.UserID,
			EarnerUserID: u.UserID,
			Balance:      0,
		}
		records = append(records, t0CD, forT0CD)
	}
	if u.couldHaveBeenEligibleForEthereumDistributionRecently(now) && tMinus1.couldHaveBeenEligibleForEthereumDistributionRecently(now) && tMinus1 != nil && tMinus1.UserID != u.UserID && t0 != nil && tMinus1.UserID != t0.UserID { //nolint:lll // .
		forTMinus1CD = &coindistribution.ByEarnerForReview{
			CreatedAt:    now,
			UserID:       tMinus1.UserID,
			EarnerUserID: u.UserID,
			Balance:      0,
		}
		records = append(records, forTMinus1CD)
	}

	if u.isEligibleForSelfForEthereumDistribution(now) {
		// Amount I've earned for myself.
		soloCD.Balance = u.processEthereumCoinDistributionForSolo(now)

		if t0 != nil && t0.UserID != u.UserID && (tMinus1 == nil || (tMinus1.UserID != u.UserID && tMinus1.UserID != t0.UserID)) && t0.isEligibleForReferralForEthereumDistribution(now) {
			// Amount my T0 earned for me.
			t0CD.Balance = u.processEthereumCoinDistributionForT0(now)
		}

		if !ethereumDistributionDryRunModeEnabled {
			u.SoloLastEthereumCoinDistributionProcessedAt = now
		}
	} else {
		u.SoloLastEthereumCoinDistributionProcessedAt = nil
	}

	if t0 != nil && t0.UserID != u.UserID && (tMinus1 == nil || (tMinus1.UserID != u.UserID && tMinus1.UserID != t0.UserID)) && u.isEligibleForT0ForEthereumDistribution(now) && t0.isEligibleForSelfForEthereumDistribution(now) { //nolint:lll // .
		// Amount I've earned for my T0.
		balanceDistributedForT0 = u.processEthereumCoinDistributionForForT0(t0, now)
		forT0CD.Balance = balanceDistributedForT0

		if !ethereumDistributionDryRunModeEnabled {
			u.ForT0LastEthereumCoinDistributionProcessedAt = now
		} else {
			balanceDistributedForT0 = 0
		}
	} else {
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
	}

	if tMinus1 != nil && tMinus1.UserID != u.UserID && t0 != nil && tMinus1.UserID != t0.UserID && u.isEligibleForTMinus1ForEthereumDistribution(now) && tMinus1.isEligibleForSelfForEthereumDistribution(now) { //nolint:lll // .
		// Amount I've earned for my T-1.
		balanceDistributedForTMinus1 = u.processEthereumCoinDistributionForForTMinus1(tMinus1, now)
		forTMinus1CD.Balance = balanceDistributedForTMinus1

		if !ethereumDistributionDryRunModeEnabled {
			u.ForTMinus1LastEthereumCoinDistributionProcessedAt = now
		} else {
			balanceDistributedForTMinus1 = 0
		}
	} else {
		u.ForTMinus1LastEthereumCoinDistributionProcessedAt = nil
	}

	return records, balanceDistributedForT0, balanceDistributedForTMinus1
}

func (u *user) processEthereumCoinDistributionForSolo(now *time.Time) float64 {
	standard, _ := tokenomics.ApplyPreStaking(u.BalanceSolo, u.PreStakingAllocation, u.PreStakingBonus)
	ethIce := coindistribution.CalculateEthereumDistributionICEBalance(standard-u.BalanceSoloEthereum, cfg.EthereumDistributionFrequency.Min, cfg.EthereumDistributionFrequency.Max, now, cfg.coinDistributionCollectorSettings.Load().EndDate) //nolint:lll // .
	if !ethereumDistributionDryRunModeEnabled {
		u.BalanceSoloEthereum += ethIce
	}

	return ethIce
}

func (u *user) processEthereumCoinDistributionForT0(now *time.Time) float64 {
	standard, _ := tokenomics.ApplyPreStaking(u.BalanceT0, u.PreStakingAllocation, u.PreStakingBonus)
	ethIce := coindistribution.CalculateEthereumDistributionICEBalance(standard-u.BalanceT0Ethereum, cfg.EthereumDistributionFrequency.Min, cfg.EthereumDistributionFrequency.Max, now, cfg.coinDistributionCollectorSettings.Load().EndDate) //nolint:lll // .
	if !ethereumDistributionDryRunModeEnabled {
		u.BalanceT0Ethereum += ethIce
	}

	return ethIce
}

// The double `For` is intended, cuz it's ForXX, where XX can be Solo/T0/ForT1/ForTMinus1.
func (u *user) processEthereumCoinDistributionForForT0(t0 *referral, now *time.Time) float64 {
	standard, _ := tokenomics.ApplyPreStaking(u.BalanceForT0, t0.PreStakingAllocation, t0.PreStakingBonus)
	ethIce := coindistribution.CalculateEthereumDistributionICEBalance(standard-u.BalanceForT0Ethereum, cfg.EthereumDistributionFrequency.Min, cfg.EthereumDistributionFrequency.Max, now, cfg.coinDistributionCollectorSettings.Load().EndDate) //nolint:lll // .
	if !ethereumDistributionDryRunModeEnabled {
		u.BalanceForT0Ethereum += ethIce
	}

	return ethIce
}

// The double `For` is intended, cuz it's ForXX, where XX can be Solo/T0/ForT1/ForTMinus1.
func (u *user) processEthereumCoinDistributionForForTMinus1(tMinus1 *referral, now *time.Time) float64 {
	standard, _ := tokenomics.ApplyPreStaking(u.BalanceForTMinus1, tMinus1.PreStakingAllocation, tMinus1.PreStakingBonus)
	ethIce := coindistribution.CalculateEthereumDistributionICEBalance(standard-u.BalanceForTMinus1Ethereum, cfg.EthereumDistributionFrequency.Min, cfg.EthereumDistributionFrequency.Max, now, cfg.coinDistributionCollectorSettings.Load().EndDate) //nolint:lll // .
	if !ethereumDistributionDryRunModeEnabled {
		u.BalanceForTMinus1Ethereum += ethIce
	}

	return ethIce
}

func isCoinDistributionCollectorEnabled(now *time.Time) bool {
	return coindistribution.IsCoinDistributionCollectorEnabled(now, cfg.EthereumDistributionFrequency.Min, cfg.coinDistributionCollectorSettings.Load())
}

func (m *miner) startCoinDistributionCollectionWorkerManager(ctx context.Context) {
	defer func() { m.stopCoinDistributionCollectionWorkerManager <- struct{}{} }()

	for ctx.Err() == nil {
		select {
		case <-m.coinDistributionStartedSignaler:
			m.coinDistributionWorkerMX.Lock()
			log.Info("started collecting coin distributions")
			reqCtx, cancel := context.WithTimeout(context.Background(), requestDeadline)
			log.Error(errors.Wrap(coindistribution.SendNewCoinDistributionCollectionCycleStartedSlackMessage(reqCtx),
				"failed to SendNewCoinDistributionCollectionCycleStartedSlackMessage"))
			cancel()
			workersStarted := int64(1)
		outerStarted:
			for ctx.Err() == nil {
				select {
				case <-m.coinDistributionStartedSignaler:
					workersStarted++
					if workersStarted == cfg.Workers {
						break outerStarted
					}
				case <-ctx.Done():
					reqCtx, cancel = context.WithTimeout(context.Background(), requestDeadline)
					log.Error(errors.Wrap(coindistribution.SendNewCoinDistributionCollectionCycleEndedPrematurelySlackMessage(reqCtx),
						"failed to SendNewCoinDistributionCollectionCycleEndedPrematurelySlackMessage"))
					cancel()
					m.coinDistributionWorkerMX.Unlock()

					return
				}
			}
			workersEnded := int64(0)
		outerEnded:
			for ctx.Err() == nil {
				select {
				case <-m.coinDistributionEndedSignaler:
					workersEnded++
					if workersEnded == cfg.Workers {
						break outerEnded
					}
				case <-ctx.Done():
					reqCtx, cancel = context.WithTimeout(context.Background(), requestDeadline)
					log.Error(errors.Wrap(coindistribution.SendNewCoinDistributionCollectionCycleEndedPrematurelySlackMessage(reqCtx),
						"failed to SendNewCoinDistributionCollectionCycleEndedPrematurelySlackMessage"))
					cancel()
					m.coinDistributionWorkerMX.Unlock()

					return
				}
			}
			if ctx.Err() != nil {
				reqCtx, cancel = context.WithTimeout(context.Background(), requestDeadline)
				log.Error(errors.Wrap(coindistribution.SendNewCoinDistributionCollectionCycleEndedPrematurelySlackMessage(reqCtx),
					"failed to SendNewCoinDistributionCollectionCycleEndedPrematurelySlackMessage"))
				cancel()
				m.coinDistributionWorkerMX.Unlock()

				return
			}
			reqCtx, cancel = context.WithTimeout(context.Background(), requestDeadline)
			m.notifyCoinDistributionCollectionCycleEnded(reqCtx)
			cancel()
			log.Info("stopped collecting coin distributions")
			m.coinDistributionWorkerMX.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (m *miner) notifyCoinDistributionCollectionCycleEnded(ctx context.Context) {
	for ctx.Err() == nil {
		reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
		if err := m.coinDistributionRepository.NotifyCoinDistributionCollectionCycleEnded(reqCtx); err != nil {
			cancel()
			log.Error(errors.Wrap(err, "failed to NotifyCoinDistributionCollectionCycleEnded"))
		} else {
			cancel()

			break
		}
	}
	for ctx.Err() == nil {
		reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
		if settings, err := m.coinDistributionRepository.GetCollectorSettings(reqCtx); err != nil {
			cancel()
			log.Error(errors.Wrap(err, "failed to GetCollectorSettings"))
		} else {
			cancel()
			cfg.coinDistributionCollectorSettings.Store(settings)

			break
		}
	}
}

func (m *miner) mustInitCoinDistributionCollector(ctx context.Context) {
	settings, err := m.coinDistributionRepository.GetCollectorSettings(ctx)
	log.Panic(err)
	cfg.coinDistributionCollectorSettings.Store(settings)
	m.coinDistributionStartedSignaler = make(chan struct{}, cfg.Workers)
	m.coinDistributionEndedSignaler = make(chan struct{}, cfg.Workers)
	m.stopCoinDistributionCollectionWorkerManager = make(chan struct{})
	m.coinDistributionWorkerMX = new(sync.Mutex)

	go m.startCoinDistributionCollectionWorkerManager(ctx)
	go m.startSynchronizingCoinDistributionCollectorSettings(ctx)
}

func (m *miner) startSynchronizingCoinDistributionCollectorSettings(ctx context.Context) {
	ticker := stdlibtime.NewTicker(30 * stdlibtime.Second) //nolint:gomnd // .
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
			if settings, err := m.coinDistributionRepository.GetCollectorSettings(reqCtx); err != nil {
				log.Error(errors.Wrap(err, "failed to GetCollectorSettings"))
			} else {
				cfg.coinDistributionCollectorSettings.Store(settings)
			}
			cancel()
		case <-ctx.Done():
			return
		}
	}
}
