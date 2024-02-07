// SPDX-License-Identifier: ice License 1.0

package storage

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	stdlibtime "time"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/chpool"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/model"
	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

//nolint:gomnd,funlen // Default configs.
func MustConnect(ctx context.Context, applicationYAMLKey string) Client {
	var cfg config
	appCfg.MustLoadFromKey(applicationYAMLKey, &cfg)
	logger, err := zap.Config{
		Level: zap.NewAtomicLevel(),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	log.Panic(err)
	cl := new(db)
	cl.cfg = &cfg
	cl.settings = append(make([]ch.Setting, 0, 3),
		ch.SettingInt("async_insert", 1),
		ch.SettingInt("wait_for_async_insert", 1))
	cl.pools = make([]*chpool.Pool, 0, len(cfg.Storage.URLs))
	for _, url := range cfg.Storage.URLs {
		pool, dErr := chpool.Dial(ctx, chpool.Options{
			ClientOptions: ch.Options{
				Logger:           logger,
				Address:          url,
				Database:         cfg.Storage.DB,
				User:             cfg.Storage.Credentials.User,
				Password:         cfg.Storage.Credentials.Password,
				Compression:      ch.CompressionLZ4,
				ReadTimeout:      30 * stdlibtime.Second,
				DialTimeout:      30 * stdlibtime.Second,
				HandshakeTimeout: 30 * stdlibtime.Second,
				Settings:         cl.settings,
			},
			MaxConnLifetime:   24 * stdlibtime.Hour,
			MaxConnIdleTime:   30 * stdlibtime.Second,
			HealthCheckPeriod: 30 * stdlibtime.Second,
			MinConns:          1,
			MaxConns:          cfg.Storage.PoolSize,
		})
		log.Panic(dErr)
		if cfg.Storage.RunDDL {
			for _, query := range strings.Split(ddl, ";") {
				if strings.TrimSpace(query) != "" {
					log.Panic(pool.Do(ctx, ch.Query{Body: query}))
				}
			}
		}
		cl.pools = append(cl.pools, pool)
	}

	return cl
}

func (db *db) Close() error {
	for _, pool := range db.pools {
		pool.Close()
	}

	return nil
}

func (db *db) Ping(ctx context.Context) error {
	errChan := make(chan error, len(db.pools))
	wg := new(sync.WaitGroup)
	wg.Add(len(db.pools))
	for _, pool := range db.pools {
		go func(cl *chpool.Pool) {
			defer wg.Done()
			errChan <- cl.Ping(ctx)
		}(pool)
	}
	wg.Wait()
	close(errChan)
	errs := make([]error, 0, len(db.pools))
	for err := range errChan {
		errs = append(errs, err)
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}

func (db *db) Insert(ctx context.Context, columns *Columns, input InsertMetadata, usrs []*model.User) error {
	if len(usrs) == 0 {
		return nil
	}
	for _, column := range input {
		column.Data.(proto.Resettable).Reset()
	}

	now := time.Now()
	truncateDuration := stdlibtime.Minute
	if !db.cfg.Development {
		truncateDuration = stdlibtime.Hour
	}

	for _, usr := range usrs {
		if usr.MiningSessionSoloLastStartedAt.IsNil() {
			columns.miningSessionSoloLastStartedAt.Append(stdlibtime.Time{})
		} else {
			columns.miningSessionSoloLastStartedAt.Append(*usr.MiningSessionSoloLastStartedAt.Time)
		}
		if usr.MiningSessionSoloStartedAt.IsNil() {
			columns.miningSessionSoloStartedAt.Append(stdlibtime.Time{})
		} else {
			columns.miningSessionSoloStartedAt.Append(*usr.MiningSessionSoloStartedAt.Time)
		}
		if usr.MiningSessionSoloEndedAt.IsNil() {
			columns.miningSessionSoloEndedAt.Append(stdlibtime.Time{})
		} else {
			columns.miningSessionSoloEndedAt.Append(*usr.MiningSessionSoloEndedAt.Time)
		}
		if usr.MiningSessionSoloPreviouslyEndedAt.IsNil() {
			columns.miningSessionSoloPreviouslyEndedAt.Append(stdlibtime.Time{})
		} else {
			columns.miningSessionSoloPreviouslyEndedAt.Append(*usr.MiningSessionSoloPreviouslyEndedAt.Time)
		}
		if usr.ExtraBonusStartedAt.IsNil() {
			columns.extraBonusStartedAt.Append(stdlibtime.Time{})
		} else {
			columns.extraBonusStartedAt.Append(*usr.ExtraBonusStartedAt.Time)
		}
		if usr.ResurrectSoloUsedAt.IsNil() {
			columns.resurrectSoloUsedAt.Append(stdlibtime.Time{})
		} else {
			columns.resurrectSoloUsedAt.Append(*usr.ResurrectSoloUsedAt.Time)
		}
		if usr.ResurrectT0UsedAt.IsNil() {
			columns.resurrectT0UsedAt.Append(stdlibtime.Time{})
		} else {
			columns.resurrectT0UsedAt.Append(*usr.ResurrectT0UsedAt.Time)
		}
		if usr.ResurrectTMinus1UsedAt.IsNil() {
			columns.resurrectTminus1UsedAt.Append(stdlibtime.Time{})
		} else {
			columns.resurrectTminus1UsedAt.Append(*usr.ResurrectTMinus1UsedAt.Time)
		}
		if usr.MiningSessionSoloDayOffLastAwardedAt.IsNil() {
			columns.miningSessionSoloDayOffLastAwardedAt.Append(stdlibtime.Time{})
		} else {
			columns.miningSessionSoloDayOffLastAwardedAt.Append(*usr.MiningSessionSoloDayOffLastAwardedAt.Time)
		}
		if usr.ExtraBonusLastClaimAvailableAt.IsNil() {
			columns.extraBonusLastClaimAvailableAt.Append(stdlibtime.Time{})
		} else {
			columns.extraBonusLastClaimAvailableAt.Append(*usr.ExtraBonusLastClaimAvailableAt.Time)
		}
		if usr.SoloLastEthereumCoinDistributionProcessedAt.IsNil() {
			columns.soloLastEthereumCoinDistributionProcessedAt.Append(stdlibtime.Time{})
		} else {
			columns.soloLastEthereumCoinDistributionProcessedAt.Append(*usr.SoloLastEthereumCoinDistributionProcessedAt.Time)
		}
		if usr.ForT0LastEthereumCoinDistributionProcessedAt.IsNil() {
			columns.forT0LastEthereumCoinDistributionProcessedAt.Append(stdlibtime.Time{})
		} else {
			columns.forT0LastEthereumCoinDistributionProcessedAt.Append(*usr.ForT0LastEthereumCoinDistributionProcessedAt.Time)
		}
		if usr.ForTMinus1LastEthereumCoinDistributionProcessedAt.IsNil() {
			columns.forTMinus1LastEthereumCoinDistributionProcessedAt.Append(stdlibtime.Time{})
		} else {
			columns.forTMinus1LastEthereumCoinDistributionProcessedAt.Append(*usr.ForTMinus1LastEthereumCoinDistributionProcessedAt.Time)
		}
		columns.createdAt.Append(now.Truncate(truncateDuration))
		columns.country.Append(usr.Country)
		columns.profilePictureName.Append(usr.ProfilePictureName)
		columns.username.Append(usr.Username)
		columns.miningBlockchainAccountAddress.Append(usr.MiningBlockchainAccountAddress)
		columns.blockchainAccountAddress.Append(usr.BlockchainAccountAddress)
		columns.userID.Append(usr.UserID)
		columns.id.Append(usr.ID)
		columns.idT0.Append(usr.IDT0)
		columns.idTminus1.Append(usr.IDTMinus1)
		columns.balanceTotalStandard.Append(usr.BalanceTotalStandard)
		columns.balanceTotalPreStaking.Append(usr.BalanceTotalPreStaking)
		columns.balanceTotalMinted.Append(usr.BalanceTotalMinted)
		columns.balanceTotalSlashed.Append(usr.BalanceTotalSlashed)
		columns.balanceSoloPending.Append(usr.BalanceSoloPending)
		columns.balanceT1Pending.Append(usr.BalanceT1Pending)
		columns.balanceT2Pending.Append(usr.BalanceT2Pending)
		columns.balanceSoloPendingApplied.Append(usr.BalanceSoloPendingApplied)
		columns.balanceT1PendingApplied.Append(usr.BalanceT1PendingApplied)
		columns.balanceT2PendingApplied.Append(usr.BalanceT2PendingApplied)
		columns.balanceSolo.Append(usr.BalanceSolo)
		columns.balanceT0.Append(usr.BalanceT0)
		columns.balanceT1.Append(usr.BalanceT1)
		columns.balanceT2.Append(usr.BalanceT2)
		columns.balanceForT0.Append(usr.BalanceForT0)
		columns.balanceForTminus1.Append(usr.BalanceForTMinus1)
		columns.balanceSoloEthereum.Append(usr.BalanceSoloEthereum)
		columns.balanceT0Ethereum.Append(usr.BalanceT0Ethereum)
		columns.balanceT1Ethereum.Append(usr.BalanceT1Ethereum)
		columns.balanceT2Ethereum.Append(usr.BalanceT2Ethereum)
		columns.balanceForT0Ethereum.Append(usr.BalanceForT0Ethereum)
		columns.balanceForTMinus1Ethereum.Append(usr.BalanceForTMinus1Ethereum)
		columns.balanceSoloEthereumMainnetRewardPoolContribution.Append(usr.BalanceSoloEthereumMainnetRewardPoolContribution)
		columns.balanceT0EthereumMainnetRewardPoolContribution.Append(usr.BalanceT0EthereumMainnetRewardPoolContribution)
		columns.balanceT1EthereumMainnetRewardPoolContribution.Append(usr.BalanceT1EthereumMainnetRewardPoolContribution)
		columns.balanceT2EthereumMainnetRewardPoolContribution.Append(usr.BalanceT2EthereumMainnetRewardPoolContribution)
		columns.balanceForT0EthereumMainnetRewardPoolContribution.Append(usr.BalanceForT0EthereumMainnetRewardPoolContribution)
		columns.balanceForTMinus1EthereumMainnetRewardPoolContribution.Append(usr.BalanceForTMinus1EthereumMainnetRewardPoolContribution)
		columns.slashingRateSolo.Append(usr.SlashingRateSolo)
		columns.slashingRateT0.Append(usr.SlashingRateT0)
		columns.slashingRateT1.Append(usr.SlashingRateT1)
		columns.slashingRateT2.Append(usr.SlashingRateT2)
		columns.slashingRateForT0.Append(usr.SlashingRateForT0)
		columns.slashingRateForTminus1.Append(usr.SlashingRateForTMinus1)
		columns.activeT1Referrals.Append(usr.ActiveT1Referrals)
		columns.activeT2Referrals.Append(usr.ActiveT2Referrals)
		columns.preStakingBonus.Append(uint16(usr.PreStakingBonus))
		columns.preStakingAllocation.Append(uint16(usr.PreStakingAllocation))
		columns.extraBonus.Append(uint16(usr.ExtraBonus))
		columns.newsSeen.Append(usr.NewsSeen)
		columns.extraBonusDaysClaimNotAvailable.Append(usr.ExtraBonusDaysClaimNotAvailable)
		columns.utcOffset.Append(int16(usr.UTCOffset))
		columns.kycStepPassed.Append(uint8(usr.KYCStepPassed))
		columns.kycStepBlocked.Append(uint8(usr.KYCStepBlocked))
		columns.kycQuizCompleted.Append(usr.KYCQuizCompleted)
		columns.kycQuizDisabled.Append(usr.KYCQuizDisabled)
		columns.hideRanking.Append(usr.HideRanking)
		kycStepsCreatedAt := make([]stdlibtime.Time, 0, 6)
		if usr.KYCStepsCreatedAt != nil {
			for _, date := range *usr.KYCStepsCreatedAt {
				if date.IsNil() {
					kycStepsCreatedAt = append(kycStepsCreatedAt, stdlibtime.Time{})
				} else {
					kycStepsCreatedAt = append(kycStepsCreatedAt, *date.Time)
				}
			}
		}
		columns.kycStepsCreatedAt.Append(kycStepsCreatedAt)
		kycStepsLastUpdatedAt := make([]stdlibtime.Time, 0, 6)
		if usr.KYCStepsLastUpdatedAt != nil {
			for _, date := range *usr.KYCStepsLastUpdatedAt {
				if date.IsNil() {
					kycStepsLastUpdatedAt = append(kycStepsLastUpdatedAt, stdlibtime.Time{})
				} else {
					kycStepsLastUpdatedAt = append(kycStepsLastUpdatedAt, *date.Time)
				}
			}
		}
		columns.kycStepsLastUpdatedAt.Append(kycStepsLastUpdatedAt)
	}

	return db.pools[atomic.AddUint64(&db.currentIndex, 1)%uint64(len(db.pools))].Do(ctx, ch.Query{
		Body:     input.Into(tableName),
		Input:    input,
		Settings: db.settings,
	})
}

func InsertDDL(rows int) (*Columns, proto.Input) {
	var (
		miningSessionSoloLastStartedAt                         = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		miningSessionSoloStartedAt                             = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		miningSessionSoloEndedAt                               = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		miningSessionSoloPreviouslyEndedAt                     = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		extraBonusStartedAt                                    = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		resurrectSoloUsedAt                                    = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		resurrectT0UsedAt                                      = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		resurrectTminus1UsedAt                                 = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		miningSessionSoloDayOffLastAwardedAt                   = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		extraBonusLastClaimAvailableAt                         = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		soloLastEthereumCoinDistributionProcessedAt            = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		forT0LastEthereumCoinDistributionProcessedAt           = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		forTMinus1LastEthereumCoinDistributionProcessedAt      = &proto.ColDateTime64{Data: make([]proto.DateTime64, 0, rows), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}
		createdAt                                              = &proto.ColDateTime{Data: make([]proto.DateTime, 0, rows), Location: stdlibtime.UTC}
		country                                                = &proto.ColStr{Buf: make([]byte, 0, 3*rows), Pos: make([]proto.Position, 0, rows)}
		profilePictureName                                     = &proto.ColStr{Buf: make([]byte, 0, 50*rows), Pos: make([]proto.Position, 0, rows)}
		username                                               = &proto.ColStr{Buf: make([]byte, 0, 40*rows), Pos: make([]proto.Position, 0, rows)}
		miningBlockchainAccountAddress                         = &proto.ColStr{Buf: make([]byte, 0, 50*rows), Pos: make([]proto.Position, 0, rows)}
		blockchainAccountAddress                               = &proto.ColStr{Buf: make([]byte, 0, 50*rows), Pos: make([]proto.Position, 0, rows)}
		userID                                                 = &proto.ColStr{Buf: make([]byte, 0, 40*rows), Pos: make([]proto.Position, 0, rows)}
		id                                                     = make(proto.ColInt64, 0, rows)
		idT0                                                   = make(proto.ColInt64, 0, rows)
		idTminus1                                              = make(proto.ColInt64, 0, rows)
		balanceTotalStandard                                   = make(proto.ColFloat64, 0, rows)
		balanceTotalPreStaking                                 = make(proto.ColFloat64, 0, rows)
		balanceTotalMinted                                     = make(proto.ColFloat64, 0, rows)
		balanceTotalSlashed                                    = make(proto.ColFloat64, 0, rows)
		balanceSoloPending                                     = make(proto.ColFloat64, 0, rows)
		balanceT1Pending                                       = make(proto.ColFloat64, 0, rows)
		balanceT2Pending                                       = make(proto.ColFloat64, 0, rows)
		balanceSoloPendingApplied                              = make(proto.ColFloat64, 0, rows)
		balanceT1PendingApplied                                = make(proto.ColFloat64, 0, rows)
		balanceT2PendingApplied                                = make(proto.ColFloat64, 0, rows)
		balanceSolo                                            = make(proto.ColFloat64, 0, rows)
		balanceT0                                              = make(proto.ColFloat64, 0, rows)
		balanceT1                                              = make(proto.ColFloat64, 0, rows)
		balanceT2                                              = make(proto.ColFloat64, 0, rows)
		balanceForT0                                           = make(proto.ColFloat64, 0, rows)
		balanceForTminus1                                      = make(proto.ColFloat64, 0, rows)
		balanceSoloEthereum                                    = make(proto.ColFloat64, 0, rows)
		balanceT0Ethereum                                      = make(proto.ColFloat64, 0, rows)
		balanceT1Ethereum                                      = make(proto.ColFloat64, 0, rows)
		balanceT2Ethereum                                      = make(proto.ColFloat64, 0, rows)
		balanceForT0Ethereum                                   = make(proto.ColFloat64, 0, rows)
		balanceForTMinus1Ethereum                              = make(proto.ColFloat64, 0, rows)
		balanceSoloEthereumMainnetRewardPoolContribution       = make(proto.ColFloat64, 0, rows)
		balanceT0EthereumMainnetRewardPoolContribution         = make(proto.ColFloat64, 0, rows)
		balanceT1EthereumMainnetRewardPoolContribution         = make(proto.ColFloat64, 0, rows)
		balanceT2EthereumMainnetRewardPoolContribution         = make(proto.ColFloat64, 0, rows)
		balanceForT0EthereumMainnetRewardPoolContribution      = make(proto.ColFloat64, 0, rows)
		balanceForTMinus1EthereumMainnetRewardPoolContribution = make(proto.ColFloat64, 0, rows)
		slashingRateSolo                                       = make(proto.ColFloat64, 0, rows)
		slashingRateT0                                         = make(proto.ColFloat64, 0, rows)
		slashingRateT1                                         = make(proto.ColFloat64, 0, rows)
		slashingRateT2                                         = make(proto.ColFloat64, 0, rows)
		slashingRateForT0                                      = make(proto.ColFloat64, 0, rows)
		slashingRateForTminus1                                 = make(proto.ColFloat64, 0, rows)
		activeT1Referrals                                      = make(proto.ColInt32, 0, rows)
		activeT2Referrals                                      = make(proto.ColInt32, 0, rows)
		preStakingBonus                                        = make(proto.ColUInt16, 0, rows)
		preStakingAllocation                                   = make(proto.ColUInt16, 0, rows)
		extraBonus                                             = make(proto.ColUInt16, 0, rows)
		newsSeen                                               = make(proto.ColUInt16, 0, rows)
		extraBonusDaysClaimNotAvailable                        = make(proto.ColUInt16, 0, rows)
		utcOffset                                              = make(proto.ColInt16, 0, rows)
		kycStepPassed                                          = make(proto.ColUInt8, 0, rows)
		kycStepBlocked                                         = make(proto.ColUInt8, 0, rows)
		kycQuizCompleted                                       = make(proto.ColBool, 0, rows)
		kycQuizDisabled                                        = make(proto.ColBool, 0, rows)
		hideRanking                                            = make(proto.ColBool, 0, rows)
		kycStepsCreatedAt                                      = proto.NewArray[stdlibtime.Time](&proto.ColDateTime64{Data: make([]proto.DateTime64, 0, 6), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}) //nolint:lll // .
		kycStepsLastUpdatedAt                                  = proto.NewArray[stdlibtime.Time](&proto.ColDateTime64{Data: make([]proto.DateTime64, 0, 6), Location: stdlibtime.UTC, Precision: proto.PrecisionMax, PrecisionSet: true}) //nolint:lll // .
	)
	input := append(make(proto.Input, 0, 72),
		proto.InputColumn{Name: "mining_session_solo_last_started_at", Data: miningSessionSoloLastStartedAt},
		proto.InputColumn{Name: "mining_session_solo_started_at", Data: miningSessionSoloStartedAt},
		proto.InputColumn{Name: "mining_session_solo_ended_at", Data: miningSessionSoloEndedAt},
		proto.InputColumn{Name: "mining_session_solo_previously_ended_at", Data: miningSessionSoloPreviouslyEndedAt},
		proto.InputColumn{Name: "extra_bonus_started_at", Data: extraBonusStartedAt},
		proto.InputColumn{Name: "resurrect_solo_used_at", Data: resurrectSoloUsedAt},
		proto.InputColumn{Name: "resurrect_t0_used_at", Data: resurrectT0UsedAt},
		proto.InputColumn{Name: "resurrect_tminus1_used_at", Data: resurrectTminus1UsedAt},
		proto.InputColumn{Name: "mining_session_solo_day_off_last_awarded_at", Data: miningSessionSoloDayOffLastAwardedAt},
		proto.InputColumn{Name: "extra_bonus_last_claim_available_at", Data: extraBonusLastClaimAvailableAt},
		proto.InputColumn{Name: "solo_last_ethereum_coin_distribution_processed_at", Data: soloLastEthereumCoinDistributionProcessedAt},
		proto.InputColumn{Name: "for_t0_last_ethereum_coin_distribution_processed_at", Data: forT0LastEthereumCoinDistributionProcessedAt},
		proto.InputColumn{Name: "for_tminus1_last_ethereum_coin_distribution_processed_at", Data: forTMinus1LastEthereumCoinDistributionProcessedAt},
		proto.InputColumn{Name: "created_at", Data: createdAt},
		proto.InputColumn{Name: "country", Data: country},
		proto.InputColumn{Name: "profile_picture_name", Data: profilePictureName},
		proto.InputColumn{Name: "username", Data: username},
		proto.InputColumn{Name: "mining_blockchain_account_address", Data: miningBlockchainAccountAddress},
		proto.InputColumn{Name: "blockchain_account_address", Data: blockchainAccountAddress},
		proto.InputColumn{Name: "user_id", Data: userID},
		proto.InputColumn{Name: "balance_total_standard", Data: &balanceTotalStandard},
		proto.InputColumn{Name: "balance_total_pre_staking", Data: &balanceTotalPreStaking},
		proto.InputColumn{Name: "balance_total_minted", Data: &balanceTotalMinted},
		proto.InputColumn{Name: "balance_total_slashed", Data: &balanceTotalSlashed},
		proto.InputColumn{Name: "balance_solo_pending", Data: &balanceSoloPending},
		proto.InputColumn{Name: "balance_t1_pending", Data: &balanceT1Pending},
		proto.InputColumn{Name: "balance_t2_pending", Data: &balanceT2Pending},
		proto.InputColumn{Name: "balance_solo_pending_applied", Data: &balanceSoloPendingApplied},
		proto.InputColumn{Name: "balance_t1_pending_applied", Data: &balanceT1PendingApplied},
		proto.InputColumn{Name: "balance_t2_pending_applied", Data: &balanceT2PendingApplied},
		proto.InputColumn{Name: "balance_solo", Data: &balanceSolo},
		proto.InputColumn{Name: "balance_t0", Data: &balanceT0},
		proto.InputColumn{Name: "balance_t1", Data: &balanceT1},
		proto.InputColumn{Name: "balance_t2", Data: &balanceT2},
		proto.InputColumn{Name: "balance_for_t0", Data: &balanceForT0},
		proto.InputColumn{Name: "balance_for_tminus1", Data: &balanceForTminus1},
		proto.InputColumn{Name: "balance_solo_ethereum", Data: &balanceSoloEthereum},
		proto.InputColumn{Name: "balance_t0_ethereum", Data: &balanceT0Ethereum},
		proto.InputColumn{Name: "balance_t1_ethereum", Data: &balanceT1Ethereum},
		proto.InputColumn{Name: "balance_t2_ethereum", Data: &balanceT2Ethereum},
		proto.InputColumn{Name: "balance_for_t0_ethereum", Data: &balanceForT0Ethereum},
		proto.InputColumn{Name: "balance_for_tminus1_ethereum", Data: &balanceForTMinus1Ethereum},
		proto.InputColumn{Name: "balance_solo_ethereum_mainnet_reward_pool_contribution", Data: &balanceSoloEthereumMainnetRewardPoolContribution},
		proto.InputColumn{Name: "balance_t0_ethereum_mainnet_reward_pool_contribution", Data: &balanceT0EthereumMainnetRewardPoolContribution},
		proto.InputColumn{Name: "balance_t1_ethereum_mainnet_reward_pool_contribution", Data: &balanceT1EthereumMainnetRewardPoolContribution},
		proto.InputColumn{Name: "balance_t2_ethereum_mainnet_reward_pool_contribution", Data: &balanceT2EthereumMainnetRewardPoolContribution},
		proto.InputColumn{Name: "balance_for_t0_ethereum_mainnet_reward_pool_contribution", Data: &balanceForT0EthereumMainnetRewardPoolContribution},
		proto.InputColumn{Name: "balance_for_tminus1_ethereum_mainnet_reward_pool_contribution", Data: &balanceForTMinus1EthereumMainnetRewardPoolContribution},
		proto.InputColumn{Name: "slashing_rate_solo", Data: &slashingRateSolo},
		proto.InputColumn{Name: "slashing_rate_t0", Data: &slashingRateT0},
		proto.InputColumn{Name: "slashing_rate_t1", Data: &slashingRateT1},
		proto.InputColumn{Name: "slashing_rate_t2", Data: &slashingRateT2},
		proto.InputColumn{Name: "slashing_rate_for_t0", Data: &slashingRateForT0},
		proto.InputColumn{Name: "slashing_rate_for_tminus1", Data: &slashingRateForTminus1},
		proto.InputColumn{Name: "id", Data: &id},
		proto.InputColumn{Name: "id_t0", Data: &idT0},
		proto.InputColumn{Name: "id_tminus1", Data: &idTminus1},
		proto.InputColumn{Name: "active_t1_referrals", Data: &activeT1Referrals},
		proto.InputColumn{Name: "active_t2_referrals", Data: &activeT2Referrals},
		proto.InputColumn{Name: "pre_staking_bonus", Data: &preStakingBonus},
		proto.InputColumn{Name: "pre_staking_allocation", Data: &preStakingAllocation},
		proto.InputColumn{Name: "extra_bonus", Data: &extraBonus},
		proto.InputColumn{Name: "news_seen", Data: &newsSeen},
		proto.InputColumn{Name: "extra_bonus_days_claim_not_available", Data: &extraBonusDaysClaimNotAvailable},
		proto.InputColumn{Name: "utc_offset", Data: &utcOffset},
		proto.InputColumn{Name: "kyc_step_passed", Data: &kycStepPassed},
		proto.InputColumn{Name: "kyc_step_blocked", Data: &kycStepBlocked},
		proto.InputColumn{Name: "kyc_quiz_completed", Data: &kycQuizCompleted},
		proto.InputColumn{Name: "kyc_quiz_disabled", Data: &kycQuizDisabled},
		proto.InputColumn{Name: "hide_ranking", Data: &hideRanking},
		proto.InputColumn{Name: "kyc_steps_created_at", Data: kycStepsCreatedAt},
		proto.InputColumn{Name: "kyc_steps_last_updated_at", Data: kycStepsLastUpdatedAt})

	return &Columns{
		miningSessionSoloLastStartedAt:                    miningSessionSoloLastStartedAt,
		miningSessionSoloStartedAt:                        miningSessionSoloStartedAt,
		miningSessionSoloEndedAt:                          miningSessionSoloEndedAt,
		miningSessionSoloPreviouslyEndedAt:                miningSessionSoloPreviouslyEndedAt,
		extraBonusStartedAt:                               extraBonusStartedAt,
		resurrectSoloUsedAt:                               resurrectSoloUsedAt,
		resurrectT0UsedAt:                                 resurrectT0UsedAt,
		resurrectTminus1UsedAt:                            resurrectTminus1UsedAt,
		miningSessionSoloDayOffLastAwardedAt:              miningSessionSoloDayOffLastAwardedAt,
		extraBonusLastClaimAvailableAt:                    extraBonusLastClaimAvailableAt,
		soloLastEthereumCoinDistributionProcessedAt:       soloLastEthereumCoinDistributionProcessedAt,
		forT0LastEthereumCoinDistributionProcessedAt:      forT0LastEthereumCoinDistributionProcessedAt,
		forTMinus1LastEthereumCoinDistributionProcessedAt: forTMinus1LastEthereumCoinDistributionProcessedAt,
		createdAt:                      createdAt,
		country:                        country,
		profilePictureName:             profilePictureName,
		username:                       username,
		miningBlockchainAccountAddress: miningBlockchainAccountAddress,
		blockchainAccountAddress:       blockchainAccountAddress,
		userID:                         userID,
		id:                             &id,
		idT0:                           &idT0,
		idTminus1:                      &idTminus1,
		balanceTotalStandard:           &balanceTotalStandard,
		balanceTotalPreStaking:         &balanceTotalPreStaking,
		balanceTotalMinted:             &balanceTotalMinted,
		balanceTotalSlashed:            &balanceTotalSlashed,
		balanceSoloPending:             &balanceSoloPending,
		balanceT1Pending:               &balanceT1Pending,
		balanceT2Pending:               &balanceT2Pending,
		balanceSoloPendingApplied:      &balanceSoloPendingApplied,
		balanceT1PendingApplied:        &balanceT1PendingApplied,
		balanceT2PendingApplied:        &balanceT2PendingApplied,
		balanceSolo:                    &balanceSolo,
		balanceT0:                      &balanceT0,
		balanceT1:                      &balanceT1,
		balanceT2:                      &balanceT2,
		balanceForT0:                   &balanceForT0,
		balanceForTminus1:              &balanceForTminus1,
		balanceSoloEthereum:            &balanceSoloEthereum,
		balanceT0Ethereum:              &balanceT0Ethereum,
		balanceT1Ethereum:              &balanceT1Ethereum,
		balanceT2Ethereum:              &balanceT2Ethereum,
		balanceForT0Ethereum:           &balanceForT0Ethereum,
		balanceForTMinus1Ethereum:      &balanceForTMinus1Ethereum,
		balanceSoloEthereumMainnetRewardPoolContribution:       &balanceSoloEthereumMainnetRewardPoolContribution,
		balanceT0EthereumMainnetRewardPoolContribution:         &balanceT0EthereumMainnetRewardPoolContribution,
		balanceT1EthereumMainnetRewardPoolContribution:         &balanceT1EthereumMainnetRewardPoolContribution,
		balanceT2EthereumMainnetRewardPoolContribution:         &balanceT2EthereumMainnetRewardPoolContribution,
		balanceForT0EthereumMainnetRewardPoolContribution:      &balanceForT0EthereumMainnetRewardPoolContribution,
		balanceForTMinus1EthereumMainnetRewardPoolContribution: &balanceForTMinus1EthereumMainnetRewardPoolContribution,
		slashingRateSolo:                &slashingRateSolo,
		slashingRateT0:                  &slashingRateT0,
		slashingRateT1:                  &slashingRateT1,
		slashingRateT2:                  &slashingRateT2,
		slashingRateForT0:               &slashingRateForT0,
		slashingRateForTminus1:          &slashingRateForTminus1,
		activeT1Referrals:               &activeT1Referrals,
		activeT2Referrals:               &activeT2Referrals,
		preStakingBonus:                 &preStakingBonus,
		preStakingAllocation:            &preStakingAllocation,
		extraBonus:                      &extraBonus,
		newsSeen:                        &newsSeen,
		extraBonusDaysClaimNotAvailable: &extraBonusDaysClaimNotAvailable,
		utcOffset:                       &utcOffset,
		kycStepPassed:                   &kycStepPassed,
		kycStepBlocked:                  &kycStepBlocked,
		kycQuizCompleted:                &kycQuizCompleted,
		kycQuizDisabled:                 &kycQuizDisabled,
		hideRanking:                     &hideRanking,
		kycStepsCreatedAt:               kycStepsCreatedAt,
		kycStepsLastUpdatedAt:           kycStepsLastUpdatedAt,
	}, input
}

func (db *db) SelectBalanceHistory(ctx context.Context, id int64, createdAts []stdlibtime.Time) ([]*BalanceHistory, error) {
	var (
		createdAt           = proto.ColDateTime{Data: make([]proto.DateTime, 0, len(createdAts)), Location: stdlibtime.UTC}
		balanceTotalMinted  = make(proto.ColFloat64, 0, len(createdAts))
		balanceTotalSlashed = make(proto.ColFloat64, 0, len(createdAts))
		res                 = make([]*BalanceHistory, 0, len(createdAts))
	)
	createdAtArray := make([]string, 0, len(createdAts))
	for _, date := range createdAts {
		format := date.UTC().Format(stdlibtime.RFC3339)
		createdAtArray = append(createdAtArray, format[0:len(format)-1])
	}
	if err := db.pools[atomic.AddUint64(&db.currentIndex, 1)%uint64(len(db.pools))].Do(ctx, ch.Query{
		Body: fmt.Sprintf(`SELECT created_at,
								  balance_total_minted, 
								  balance_total_slashed 
						   FROM %[1]v
						   WHERE id = %[2]v
						     AND created_at IN ['%[3]v']`, tableName, id, strings.Join(createdAtArray, "','")),
		Result: append(make(proto.Results, 0, 3),
			proto.ResultColumn{Name: "created_at", Data: &createdAt},
			proto.ResultColumn{Name: "balance_total_minted", Data: &balanceTotalMinted},
			proto.ResultColumn{Name: "balance_total_slashed", Data: &balanceTotalSlashed}),
		OnResult: func(_ context.Context, block proto.Block) error {
			for ix := 0; ix < block.Rows; ix++ {
				res = append(res, &BalanceHistory{
					CreatedAt:           time.New((&createdAt).Row(ix)),
					BalanceTotalMinted:  (&balanceTotalMinted).Row(ix),
					BalanceTotalSlashed: (&balanceTotalSlashed).Row(ix),
				})
			}
			(&createdAt).Reset()
			(&balanceTotalMinted).Reset()
			(&balanceTotalSlashed).Reset()

			return nil
		},
		Secret:      "",
		InitialUser: "",
	}); err != nil {
		return nil, err
	}
	dedupedRes := make([]*BalanceHistory, 0, len(createdAts))
	for _, rowA := range res {
		found := false
		for _, rowB := range dedupedRes {
			if rowA.CreatedAt.Equal(*rowB.CreatedAt.Time) {
				found = true

				break
			}
		}
		if !found {
			dedupedRes = append(dedupedRes, rowA)
		}
	}
	res = dedupedRes

	return res, nil
}

func (db *db) SelectTotalCoins(ctx context.Context, createdAts []stdlibtime.Time) ([]*TotalCoins, error) {
	var (
		createdAt              = proto.ColDateTime{Data: make([]proto.DateTime, 0, len(createdAts)), Location: stdlibtime.UTC}
		balanceTotalStandard   = make(proto.ColFloat64, 0, len(createdAts))
		balanceTotalPreStaking = make(proto.ColFloat64, 0, len(createdAts))
		balanceTotalEthereum   = make(proto.ColFloat64, 0, len(createdAts))
		res                    = make([]*TotalCoins, 0, len(createdAts))
	)
	createdAtArray := make([]string, 0, len(createdAts))
	for _, date := range createdAts {
		format := date.UTC().Format(stdlibtime.RFC3339)
		createdAtArray = append(createdAtArray, format[0:len(format)-1])
	}
	if err := db.pools[atomic.AddUint64(&db.currentIndex, 1)%uint64(len(db.pools))].Do(ctx, ch.Query{
		Body: fmt.Sprintf(selectTotalCoinsSQL, tableName, strings.Join(createdAtArray, "','"), users.LivenessDetectionKYCStep),
		Result: append(make(proto.Results, 0, 4),
			proto.ResultColumn{Name: "created_at", Data: &createdAt},
			proto.ResultColumn{Name: "balance_total_standard", Data: &balanceTotalStandard},
			proto.ResultColumn{Name: "balance_total_pre_staking", Data: &balanceTotalPreStaking},
			proto.ResultColumn{Name: "balance_total_ethereum", Data: &balanceTotalEthereum}),
		OnResult: func(_ context.Context, block proto.Block) error {
			for ix := 0; ix < block.Rows; ix++ {
				resItem := &TotalCoins{
					CreatedAt:              time.New((&createdAt).Row(ix)),
					BalanceTotalStandard:   (&balanceTotalStandard).Row(ix),
					BalanceTotalPreStaking: (&balanceTotalPreStaking).Row(ix),
					BalanceTotalEthereum:   (&balanceTotalEthereum).Row(ix),
				}
				resItem.BalanceTotalStandard += resItem.BalanceTotalPreStaking - resItem.BalanceTotalEthereum
				resItem.BalanceTotal = resItem.BalanceTotalStandard + resItem.BalanceTotalEthereum
				res = append(res, resItem)
			}
			(&createdAt).Reset()
			(&balanceTotalStandard).Reset()
			(&balanceTotalPreStaking).Reset()
			(&balanceTotalEthereum).Reset()

			return nil
		},
		Secret:      "",
		InitialUser: "",
	}); err != nil {
		return nil, err
	}
	dedupedRes := make([]*TotalCoins, 0, len(createdAts))

	for _, cAt := range createdAts {
		var found *TotalCoins = nil
		for _, rowA := range res {
			if rowA.CreatedAt.Equal(cAt) {
				found = rowA
				break
			}
		}
		if found != nil {
			dedupedRes = append(dedupedRes, found)
		} else {
			dedupedRes = append(dedupedRes, &TotalCoins{
				CreatedAt:              time.New(cAt),
				BalanceTotalStandard:   0,
				BalanceTotalPreStaking: 0,
				BalanceTotalEthereum:   0,
				BalanceTotal:           0,
			})
		}
	}
	res = dedupedRes

	return res, nil
}

func (t *TotalCoins) Key() string {
	return fmt.Sprintf("totalCoinStats:%v", t.CreatedAt.Format(stdlibtime.RFC3339))
}
