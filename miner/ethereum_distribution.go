// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"fmt"
	"sync"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	coindistribution "github.com/ice-blockchain/freezer/coin-distribution"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (ref *referral) username() string {
	if ref != nil && ref.Username != "" {
		return ref.Username
	}

	return "icenetwork/bogus"
}

func (ref *referral) isEligibleForSelfForEthereumDistribution(now, lastEthereumCoinDistributionProcessedAt *time.Time) bool {
	coinDistributionCollectorSettings := cfg.coinDistributionCollectorSettings.Load()

	return ref != nil &&
		ref.ID != 0 &&
		coindistribution.IsEligibleForEthereumDistributionNow(
			ref.ID,
			now,
			lastEthereumCoinDistributionProcessedAt,
			cfg.coinDistributionCollectorSettings.Load().StartDate,
			cfg.coinDistributionCollectorSettings.Load().LatestDate,
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
			0.1,
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
			cfg.coinDistributionCollectorSettings.Load().LatestDate,
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

func (u *user) isEligibleForT0ForEthereumDistribution(now *time.Time, idT0 int64) bool {
	return u != nil &&
		u.ID != 0 &&
		coindistribution.IsEligibleForEthereumDistributionNow(
			idT0,
			now,
			u.ForT0LastEthereumCoinDistributionProcessedAt,
			cfg.coinDistributionCollectorSettings.Load().StartDate,
			cfg.coinDistributionCollectorSettings.Load().LatestDate,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max) &&
		u.isEligibleForReferralForEthereumDistribution(now)
}

func (u *user) isEligibleForTMinus1ForEthereumDistribution(now *time.Time, idTMinus1 int64) bool {
	return u != nil &&
		u.ID != 0 &&
		coindistribution.IsEligibleForEthereumDistributionNow(
			idTMinus1,
			now,
			u.ForTMinus1LastEthereumCoinDistributionProcessedAt,
			cfg.coinDistributionCollectorSettings.Load().StartDate,
			cfg.coinDistributionCollectorSettings.Load().LatestDate,
			cfg.EthereumDistributionFrequency.Min,
			cfg.EthereumDistributionFrequency.Max) &&
		u.isEligibleForReferralForEthereumDistribution(now)
}

func (u *user) isEligibleForReferralForEthereumDistribution(now *time.Time) bool {
	coinDistributionCollectorSettings := cfg.coinDistributionCollectorSettings.Load()
	return coindistribution.IsEligibleForEthereumDistribution(
		0,
		0.1,
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
	return u != nil && !u.MiningSessionSoloEndedAt.IsNil() && u.KYCStepPassedCorrectly(users.QuizKYCStep) && (u.MiningSessionSoloEndedAt.After(now.Add(-(cfg.MiningSessionDuration.Max / 8))) || coindistribution.AllowInactiveUsers) //nolint:lll // .
}

func (ref *referral) couldHaveBeenEligibleForEthereumDistributionRecently(now *time.Time) bool {
	return ref != nil && !ref.MiningSessionSoloEndedAt.IsNil() && ref.KYCStepPassedCorrectly(users.QuizKYCStep) && (ref.MiningSessionSoloEndedAt.After(now.Add(-(cfg.MiningSessionDuration.Max / 8))) || coindistribution.AllowInactiveUsers) //nolint:lll // .
}

//nolint:funlen // .
func (u *user) processEthereumCoinDistribution(
	enabled bool, now *time.Time, t0, tMinus1 *referral,
) (records []*coindistribution.ByEarnerForReview, balanceDistributedForT0, balanceDistributedForTMinus1 float64) {
	if !enabled {
		if u.BalanceSoloEthereumPending != nil {
			u.BalanceSoloEthereum += float64(*u.BalanceSoloEthereumPending)
			u.BalanceSoloEthereumMainnetRewardPoolContribution += float64(*u.BalanceSoloEthereumPending) * cfg.MainnetRewardPoolContributionPercentage
			u.BalanceSoloEthereumPending = new(model.FlexibleFloat64)
		}
		if u.BalanceT0EthereumPending != nil {
			u.BalanceT0Ethereum += float64(*u.BalanceT0EthereumPending)
			u.BalanceT0EthereumMainnetRewardPoolContribution += float64(*u.BalanceT0EthereumPending) * cfg.MainnetRewardPoolContributionPercentage
			u.BalanceT0EthereumPending = new(model.FlexibleFloat64)
		}
		if u.BalanceT1EthereumPending != nil {
			u.BalanceT1Ethereum += float64(*u.BalanceT1EthereumPending)
			u.BalanceT1EthereumMainnetRewardPoolContribution += float64(*u.BalanceT1EthereumPending) * cfg.MainnetRewardPoolContributionPercentage
			u.BalanceT1EthereumPending = new(model.FlexibleFloat64)
		}
		if u.BalanceT2EthereumPending != nil {
			u.BalanceT2Ethereum += float64(*u.BalanceT2EthereumPending)
			u.BalanceT2EthereumMainnetRewardPoolContribution += float64(*u.BalanceT2EthereumPending) * cfg.MainnetRewardPoolContributionPercentage
			u.BalanceT2EthereumPending = new(model.FlexibleFloat64)
		}
		u.SoloLastEthereumCoinDistributionProcessedAt = nil
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
		u.ForTMinus1LastEthereumCoinDistributionProcessedAt = nil

		return nil, 0, 0
	}
	u.BalanceSoloEthereumPending = nil
	u.BalanceT0EthereumPending = nil
	u.BalanceT1EthereumPending = nil
	u.BalanceT2EthereumPending = nil
	records = make([]*coindistribution.ByEarnerForReview, 0, 1+1+1+1+1+1+1+1)
	var (
		t0CD, t0MainnetRewardPoolContributionCD,
		forT0CD, forT0MainnetRewardPoolContributionCD,
		forTMinus1CD, forTMinus1MainnetRewardPoolContributionCD,
		soloCD, soloMainnetRewardPoolContributionCD *coindistribution.ByEarnerForReview
	)
	const mainnetRewardPoolContributionIdentifier = "mainnet/reward/pool/contribution"
	if u != nil && !u.MiningSessionSoloEndedAt.IsNil() && wasNotProcessedToday(now, u.SoloLastEthereumCoinDistributionProcessedAt) {
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
		soloMainnetRewardPoolContributionCD = &coindistribution.ByEarnerForReview{
			CreatedAt:          now,
			Username:           mainnetRewardPoolContributionIdentifier,
			ReferredByUsername: mainnetRewardPoolContributionIdentifier,
			UserID:             mainnetRewardPoolContributionIdentifier,
			EarnerUserID:       fmt.Sprintf("(%v,%v)", u.UserID, u.UserID),
			EthAddress:         cfg.MainnetRewardPoolContributionEthAddress,
			InternalID:         999999999,
			Balance:            0,
		}
		records = append(records, soloCD, soloMainnetRewardPoolContributionCD)
	}
	if u != nil && !u.MiningSessionSoloEndedAt.IsNil() && t0 != nil && !t0.MiningSessionSoloEndedAt.IsNil() && t0.UserID != u.UserID && (tMinus1 == nil || (tMinus1.UserID != u.UserID && tMinus1.UserID != t0.UserID)) { //nolint:lll // .
		if wasNotProcessedToday(now, u.SoloLastEthereumCoinDistributionProcessedAt) {
			t0CD = &coindistribution.ByEarnerForReview{
				CreatedAt:    now,
				UserID:       u.UserID,
				EarnerUserID: t0.UserID,
				Balance:      0,
			}
			t0MainnetRewardPoolContributionCD = &coindistribution.ByEarnerForReview{
				CreatedAt:    now,
				UserID:       mainnetRewardPoolContributionIdentifier,
				EarnerUserID: fmt.Sprintf("(%v,%v)", u.UserID, t0.UserID),
				Balance:      0,
			}
			records = append(records, t0CD, t0MainnetRewardPoolContributionCD)
		}
		if wasNotProcessedToday(now, u.ForT0LastEthereumCoinDistributionProcessedAt) {
			forT0CD = &coindistribution.ByEarnerForReview{
				CreatedAt:    now,
				UserID:       t0.UserID,
				EarnerUserID: u.UserID,
				Balance:      0,
			}
			forT0MainnetRewardPoolContributionCD = &coindistribution.ByEarnerForReview{
				CreatedAt:    now,
				UserID:       mainnetRewardPoolContributionIdentifier,
				EarnerUserID: fmt.Sprintf("(%v,%v)", t0.UserID, u.UserID),
				Balance:      0,
			}
			records = append(records, forT0CD, forT0MainnetRewardPoolContributionCD)
		}
	}
	if u != nil && !u.MiningSessionSoloEndedAt.IsNil() && tMinus1 != nil && !tMinus1.MiningSessionSoloEndedAt.IsNil() && tMinus1.UserID != u.UserID && t0 != nil && tMinus1.UserID != t0.UserID && wasNotProcessedToday(now, u.ForTMinus1LastEthereumCoinDistributionProcessedAt) { //nolint:lll // .
		forTMinus1CD = &coindistribution.ByEarnerForReview{
			CreatedAt:    now,
			UserID:       tMinus1.UserID,
			EarnerUserID: u.UserID,
			Balance:      0,
		}
		forTMinus1MainnetRewardPoolContributionCD = &coindistribution.ByEarnerForReview{
			CreatedAt:    now,
			UserID:       mainnetRewardPoolContributionIdentifier,
			EarnerUserID: fmt.Sprintf("(%v,%v)", tMinus1.UserID, u.UserID),
			Balance:      0,
		}
		records = append(records, forTMinus1CD, forTMinus1MainnetRewardPoolContributionCD)
	}

	if u.isEligibleForSelfForEthereumDistribution(now) {
		// Amount I've earned for myself.
		soloCD.Balance = u.processEthereumCoinDistributionForSolo(now)
		soloMainnetRewardPoolContributionCD.Balance = soloCD.Balance * cfg.MainnetRewardPoolContributionPercentage
		soloCD.Balance -= soloMainnetRewardPoolContributionCD.Balance
		totalForSelf := soloCD.Balance

		if t0 != nil && t0.UserID != u.UserID && (tMinus1 == nil || (tMinus1.UserID != u.UserID && tMinus1.UserID != t0.UserID)) && t0.isEligibleForReferralForEthereumDistribution(now) { //nolint:lll // .
			// Amount my T0 earned for me.
			t0CD.Balance = u.processEthereumCoinDistributionForT0(now)
			t0MainnetRewardPoolContributionCD.Balance = t0CD.Balance * cfg.MainnetRewardPoolContributionPercentage
			t0CD.Balance -= t0MainnetRewardPoolContributionCD.Balance
			totalForSelf += t0CD.Balance
		}

		if totalForSelf > 0 {
			u.SoloLastEthereumCoinDistributionProcessedAt = now
		} else {
			u.SoloLastEthereumCoinDistributionProcessedAt = nil
		}
	} else {
		u.SoloLastEthereumCoinDistributionProcessedAt = nil
	}

	if t0 != nil && t0.UserID != u.UserID && (tMinus1 == nil || (tMinus1.UserID != u.UserID && tMinus1.UserID != t0.UserID)) && u.isEligibleForT0ForEthereumDistribution(now, t0.ID) && t0.isEligibleForSelfForEthereumDistribution(now, u.ForT0LastEthereumCoinDistributionProcessedAt) { //nolint:lll // .
		// Amount I've earned for my T0.
		balanceDistributedForT0 = u.processEthereumCoinDistributionForForT0(t0, now)
		forT0CD.Balance = balanceDistributedForT0
		forT0MainnetRewardPoolContributionCD.Balance = forT0CD.Balance * cfg.MainnetRewardPoolContributionPercentage
		forT0CD.Balance -= forT0MainnetRewardPoolContributionCD.Balance

		if forT0CD.Balance > 0 {
			u.ForT0LastEthereumCoinDistributionProcessedAt = now
		} else {
			u.ForT0LastEthereumCoinDistributionProcessedAt = nil
			balanceDistributedForT0 = 0
		}
	} else {
		u.ForT0LastEthereumCoinDistributionProcessedAt = nil
	}

	if tMinus1 != nil && tMinus1.UserID != u.UserID && t0 != nil && tMinus1.UserID != t0.UserID && u.isEligibleForTMinus1ForEthereumDistribution(now, tMinus1.ID) && tMinus1.isEligibleForSelfForEthereumDistribution(now, u.ForTMinus1LastEthereumCoinDistributionProcessedAt) { //nolint:lll // .
		// Amount I've earned for my T-1.
		balanceDistributedForTMinus1 = u.processEthereumCoinDistributionForForTMinus1(tMinus1, now)
		forTMinus1CD.Balance = balanceDistributedForTMinus1
		forTMinus1MainnetRewardPoolContributionCD.Balance = forTMinus1CD.Balance * cfg.MainnetRewardPoolContributionPercentage
		forTMinus1CD.Balance -= forTMinus1MainnetRewardPoolContributionCD.Balance

		if forTMinus1CD.Balance > 0 {
			u.ForTMinus1LastEthereumCoinDistributionProcessedAt = now
		} else {
			u.ForTMinus1LastEthereumCoinDistributionProcessedAt = nil
			balanceDistributedForTMinus1 = 0
		}
	} else {
		u.ForTMinus1LastEthereumCoinDistributionProcessedAt = nil
	}

	return records, balanceDistributedForT0, balanceDistributedForTMinus1
}

func wasNotProcessedToday(now, lastEthereumCoinDistributionProcessedAt *time.Time) bool {
	return lastEthereumCoinDistributionProcessedAt.IsNil() ||
		!lastEthereumCoinDistributionProcessedAt.Truncate(cfg.EthereumDistributionFrequency.Min).Equal(now.Truncate(cfg.EthereumDistributionFrequency.Min))
}

func (u *user) processEthereumCoinDistributionForSolo(now *time.Time) float64 {
	standard, _ := tokenomics.ApplyPreStaking(u.BalanceSolo, u.PreStakingAllocation, u.PreStakingBonus)
	ethIce := coindistribution.CalculateEthereumDistributionICEBalance(standard-u.BalanceSoloEthereum, cfg.EthereumDistributionFrequency.Min, cfg.EthereumDistributionFrequency.Max, now, cfg.coinDistributionCollectorSettings.Load().EndDate) //nolint:lll // .
	if ethIce <= 0 {
		return 0
	}

	val := model.FlexibleFloat64(ethIce)
	u.BalanceSoloEthereumPending = &val

	return ethIce
}

func (u *user) processEthereumCoinDistributionForT0(now *time.Time) float64 {
	standard, _ := tokenomics.ApplyPreStaking(u.BalanceT0, u.PreStakingAllocation, u.PreStakingBonus)
	ethIce := coindistribution.CalculateEthereumDistributionICEBalance(standard-u.BalanceT0Ethereum, cfg.EthereumDistributionFrequency.Min, cfg.EthereumDistributionFrequency.Max, now, cfg.coinDistributionCollectorSettings.Load().EndDate) //nolint:lll // .
	if ethIce <= 0 {
		return 0
	}

	val := model.FlexibleFloat64(ethIce)
	u.BalanceT0EthereumPending = &val

	return ethIce
}

// The double `For` is intended, cuz it's ForXX, where XX can be Solo/T0/ForT1/ForTMinus1.
func (u *user) processEthereumCoinDistributionForForT0(t0 *referral, now *time.Time) float64 {
	standard, _ := tokenomics.ApplyPreStaking(u.BalanceForT0, t0.PreStakingAllocation, t0.PreStakingBonus)
	ethIce := coindistribution.CalculateEthereumDistributionICEBalance(standard-u.BalanceForT0Ethereum, cfg.EthereumDistributionFrequency.Min, cfg.EthereumDistributionFrequency.Max, now, cfg.coinDistributionCollectorSettings.Load().EndDate) //nolint:lll // .
	if ethIce <= 0 {
		return 0
	}

	u.BalanceForT0Ethereum += ethIce

	return ethIce
}

// The double `For` is intended, cuz it's ForXX, where XX can be Solo/T0/ForT1/ForTMinus1.
func (u *user) processEthereumCoinDistributionForForTMinus1(tMinus1 *referral, now *time.Time) float64 {
	standard, _ := tokenomics.ApplyPreStaking(u.BalanceForTMinus1, tMinus1.PreStakingAllocation, tMinus1.PreStakingBonus)
	ethIce := coindistribution.CalculateEthereumDistributionICEBalance(standard-u.BalanceForTMinus1Ethereum, cfg.EthereumDistributionFrequency.Min, cfg.EthereumDistributionFrequency.Max, now, cfg.coinDistributionCollectorSettings.Load().EndDate) //nolint:lll // .
	if ethIce <= 0 {
		return 0
	}

	u.BalanceForTMinus1Ethereum += ethIce

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
			before := time.Now()
			cfg.coinDistributionCollectorStartedAt.Store(before)
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
			after := time.Now()
			reqCtx, cancel = context.WithTimeout(context.Background(), requestDeadline)
			m.notifyCoinDistributionCollectionCycleEnded(reqCtx)
			cancel()
			log.Info(fmt.Sprintf("finished collecting coin distributions in %v", after.Sub(*before.Time)))
			cfg.coinDistributionCollectorStartedAt.Store(new(time.Time))
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

func currentCoinDistributionCollectingEndedAt() *time.Time {
	coinDistributionCollectorStartedAt := cfg.coinDistributionCollectorStartedAt.Load()
	if coinDistributionCollectorStartedAt.IsNil() {
		coinDistributionCollectorStartedAt = time.New(time.Now().Add(-1 * stdlibtime.Millisecond))
	}
	var startingWindow stdlibtime.Duration
	if cfg.Development {
		startingWindow = 10 * stdlibtime.Second
	} else {
		startingWindow = 20 * stdlibtime.Minute
	}

	return time.New(coinDistributionCollectorStartedAt.Add(startingWindow))
}
