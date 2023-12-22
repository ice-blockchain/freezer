// SPDX-License-Identifier: ice License 1.0

package storage

import (
	"context"
	_ "embed"
	"io"
	stdlibtime "time"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/chpool"
	"github.com/ClickHouse/ch-go/proto"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

type (
	Client interface {
		io.Closer
		Ping(ctx context.Context) error
		Insert(ctx context.Context, columns *Columns, input InsertMetadata, usrs []*model.User) error
		SelectBalanceHistory(ctx context.Context, id int64, createdAts []stdlibtime.Time) ([]*BalanceHistory, error)
		SelectTotalCoins(ctx context.Context, createdAts []stdlibtime.Time) ([]*TotalCoins, error)
	}
	BalanceHistory struct {
		CreatedAt                               *time.Time
		BalanceTotalMinted, BalanceTotalSlashed float64
	}
	TotalCoins struct {
		CreatedAt              *time.Time `redis:"created_at"`
		BalanceTotalStandard   float64    `redis:"standard"`
		BalanceTotalPreStaking float64    `redis:"pre_staking"`
		BalanceTotalEthereum   float64    `redis:"blockchain"`
	}
	InsertMetadata = proto.Input
	Columns        struct {
		miningSessionSoloLastStartedAt                    *proto.ColDateTime64
		miningSessionSoloStartedAt                        *proto.ColDateTime64
		miningSessionSoloEndedAt                          *proto.ColDateTime64
		miningSessionSoloPreviouslyEndedAt                *proto.ColDateTime64
		extraBonusStartedAt                               *proto.ColDateTime64
		resurrectSoloUsedAt                               *proto.ColDateTime64
		resurrectT0UsedAt                                 *proto.ColDateTime64
		resurrectTminus1UsedAt                            *proto.ColDateTime64
		miningSessionSoloDayOffLastAwardedAt              *proto.ColDateTime64
		extraBonusLastClaimAvailableAt                    *proto.ColDateTime64
		soloLastEthereumCoinDistributionProcessedAt       *proto.ColDateTime64
		forT0LastEthereumCoinDistributionProcessedAt      *proto.ColDateTime64
		forTMinus1LastEthereumCoinDistributionProcessedAt *proto.ColDateTime64
		createdAt                                         *proto.ColDateTime
		country                                           *proto.ColStr
		profilePictureName                                *proto.ColStr
		username                                          *proto.ColStr
		miningBlockchainAccountAddress                    *proto.ColStr
		blockchainAccountAddress                          *proto.ColStr
		userID                                            *proto.ColStr
		id                                                *proto.ColInt64
		idT0                                              *proto.ColInt64
		idTminus1                                         *proto.ColInt64
		balanceTotalStandard                              *proto.ColFloat64
		balanceTotalPreStaking                            *proto.ColFloat64
		balanceTotalMinted                                *proto.ColFloat64
		balanceTotalSlashed                               *proto.ColFloat64
		balanceSoloPending                                *proto.ColFloat64
		balanceT1Pending                                  *proto.ColFloat64
		balanceT2Pending                                  *proto.ColFloat64
		balanceSoloPendingApplied                         *proto.ColFloat64
		balanceT1PendingApplied                           *proto.ColFloat64
		balanceT2PendingApplied                           *proto.ColFloat64
		balanceSolo                                       *proto.ColFloat64
		balanceT0                                         *proto.ColFloat64
		balanceT1                                         *proto.ColFloat64
		balanceT2                                         *proto.ColFloat64
		balanceForT0                                      *proto.ColFloat64
		balanceForTminus1                                 *proto.ColFloat64
		balanceSoloEthereum                               *proto.ColFloat64
		balanceT0Ethereum                                 *proto.ColFloat64
		balanceT1Ethereum                                 *proto.ColFloat64
		balanceT2Ethereum                                 *proto.ColFloat64
		balanceForT0Ethereum                              *proto.ColFloat64
		balanceForTMinus1Ethereum                         *proto.ColFloat64
		slashingRateSolo                                  *proto.ColFloat64
		slashingRateT0                                    *proto.ColFloat64
		slashingRateT1                                    *proto.ColFloat64
		slashingRateT2                                    *proto.ColFloat64
		slashingRateForT0                                 *proto.ColFloat64
		slashingRateForTminus1                            *proto.ColFloat64
		activeT1Referrals                                 *proto.ColInt32
		activeT2Referrals                                 *proto.ColInt32
		preStakingBonus                                   *proto.ColUInt16
		preStakingAllocation                              *proto.ColUInt16
		extraBonus                                        *proto.ColUInt16
		newsSeen                                          *proto.ColUInt16
		extraBonusDaysClaimNotAvailable                   *proto.ColUInt16
		utcOffset                                         *proto.ColInt16
		kycStepPassed                                     *proto.ColUInt8
		kycStepBlocked                                    *proto.ColUInt8
		hideRanking                                       *proto.ColBool
		kycStepsCreatedAt                                 *proto.ColArr[stdlibtime.Time]
		kycStepsLastUpdatedAt                             *proto.ColArr[stdlibtime.Time]
	}
)

// Private API.

const (
	tableName                = "freezer_user_history"
	kycStepToCalculateTotals = 2
)

// .
var (
	//go:embed ddl.sql
	ddl string
)

type (
	db struct {
		cfg          *config
		pools        []*chpool.Pool
		settings     []ch.Setting
		currentIndex uint64
	}
	config struct {
		Storage struct {
			Credentials struct {
				User     string `yaml:"user"`
				Password string `yaml:"password"`
			} `yaml:"credentials" mapstructure:"credentials"`
			DB       string   `yaml:"db" mapstructure:"db"`
			URLs     []string `yaml:"urls" mapstructure:"urls"`
			PoolSize int32    `yaml:"poolSize" mapstructure:"poolSize"`
			RunDDL   bool     `yaml:"runDDL" mapstructure:"runDDL"`
		} `yaml:"bookkeeper/storage" mapstructure:"bookkeeper/storage"`
		Development bool `yaml:"development" mapstructure:"development"`
	}
)
