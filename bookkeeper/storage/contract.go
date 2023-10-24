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
		GetAdjustUserInformation(ctx context.Context, userIDs map[int64]struct{}, limit, offset int64) ([]*AdjustUserInfo, error)
	}
	AdjustUserInfo struct {
		MiningSessionSoloStartedAt         *time.Time
		MiningSessionSoloEndedAt           *time.Time
		MiningSessionSoloLastStartedAt     *time.Time
		MiningSessionSoloPreviouslyEndedAt *time.Time
		CreatedAt                          *time.Time
		ResurrectSoloUsedAt                *time.Time
		ID                                 int64
		SlashingRateSolo                   float64
		SlashingRateT1                     float64
		SlashingRateT2                     float64
		BalanceSolo                        float64
		BalanceT0                          float64
		BalanceT1Pending                   float64
		BalanceT1PendingApplied            float64
		BalanceT2Pending                   float64
		BalanceT2PendingApplied            float64
		PrestakingAllocation               uint16
		PrestakingBonus                    uint16
	}
	BalanceHistory struct {
		CreatedAt                               *time.Time
		BalanceTotalMinted, BalanceTotalSlashed float64
	}
	InsertMetadata = proto.Input
	Columns        struct {
		miningSessionSoloLastStartedAt       *proto.ColDateTime64
		miningSessionSoloStartedAt           *proto.ColDateTime64
		miningSessionSoloEndedAt             *proto.ColDateTime64
		miningSessionSoloPreviouslyEndedAt   *proto.ColDateTime64
		extraBonusStartedAt                  *proto.ColDateTime64
		resurrectSoloUsedAt                  *proto.ColDateTime64
		resurrectT0UsedAt                    *proto.ColDateTime64
		resurrectTminus1UsedAt               *proto.ColDateTime64
		miningSessionSoloDayOffLastAwardedAt *proto.ColDateTime64
		extraBonusLastClaimAvailableAt       *proto.ColDateTime64
		createdAt                            *proto.ColDateTime
		profilePictureName                   *proto.ColStr
		username                             *proto.ColStr
		miningBlockchainAccountAddress       *proto.ColStr
		blockchainAccountAddress             *proto.ColStr
		userID                               *proto.ColStr
		id                                   *proto.ColInt64
		idT0                                 *proto.ColInt64
		idTminus1                            *proto.ColInt64
		balanceTotalStandard                 *proto.ColFloat64
		balanceTotalPreStaking               *proto.ColFloat64
		balanceTotalMinted                   *proto.ColFloat64
		balanceTotalSlashed                  *proto.ColFloat64
		balanceSoloPending                   *proto.ColFloat64
		balanceT1Pending                     *proto.ColFloat64
		balanceT2Pending                     *proto.ColFloat64
		balanceSoloPendingApplied            *proto.ColFloat64
		balanceT1PendingApplied              *proto.ColFloat64
		balanceT2PendingApplied              *proto.ColFloat64
		balanceSolo                          *proto.ColFloat64
		balanceT0                            *proto.ColFloat64
		balanceT1                            *proto.ColFloat64
		balanceT2                            *proto.ColFloat64
		balanceForT0                         *proto.ColFloat64
		balanceForTminus1                    *proto.ColFloat64
		slashingRateSolo                     *proto.ColFloat64
		slashingRateT0                       *proto.ColFloat64
		slashingRateT1                       *proto.ColFloat64
		slashingRateT2                       *proto.ColFloat64
		slashingRateForT0                    *proto.ColFloat64
		slashingRateForTminus1               *proto.ColFloat64
		activeT1Referrals                    *proto.ColInt32
		activeT2Referrals                    *proto.ColInt32
		preStakingBonus                      *proto.ColUInt16
		preStakingAllocation                 *proto.ColUInt16
		extraBonus                           *proto.ColUInt16
		newsSeen                             *proto.ColUInt16
		extraBonusDaysClaimNotAvailable      *proto.ColUInt16
		utcOffset                            *proto.ColInt16
		hideRanking                          *proto.ColBool
	}
)

// Private API.

const (
	tableName = "freezer_user_history"
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
