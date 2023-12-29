// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"errors"
	"io"
	"math/big"
	"sync"
	stdlibtime "time"

	"github.com/alitto/pond"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/time"
)

// Public API.

type (
	Client interface {
		io.Closer
		CheckHealth(ctx context.Context) error
	}

	Repository interface {
		io.Closer
		GetCoinDistributionsForReview(ctx context.Context, arg *GetCoinDistributionsForReviewArg) (*CoinDistributionsForReview, error)
		CheckHealth(ctx context.Context) error
		ReviewCoinDistributions(ctx context.Context, reviewerUserID string, decision string) error
		NotifyCoinDistributionCollectionCycleEnded(ctx context.Context) error
		GetCollectorSettings(ctx context.Context) (*CollectorSettings, error)
		CollectCoinDistributionsForReview(ctx context.Context, records []*ByEarnerForReview) error
	}
	CollectorSettings struct {
		DeniedCountries          map[string]struct{}
		LatestDate               *time.Time
		StartDate                *time.Time
		EndDate                  *time.Time
		MinBalanceRequired       float64
		StartHour                int
		MinMiningStreaksRequired uint64
		Enabled                  bool
		ForcedExecution          bool
	}

	CoinDistributionsForReview struct {
		Distributions []*PendingReview `json:"distributions"`
		Cursor        uint64           `json:"cursor" example:"5065"`
		TotalRows     uint64           `json:"totalRows" example:"5065"`
		TotalIce      float64          `json:"totalIce" example:"5065.3"`
	}

	GetCoinDistributionsForReviewArg struct {
		CreatedAtOrderBy          string `form:"createdAtOrderBy" example:"asc"`
		IceOrderBy                string `form:"iceOrderBy" example:"asc"`
		UsernameOrderBy           string `form:"usernameOrderBy" example:"asc"`
		ReferredByUsernameOrderBy string `form:"referredByUsernameOrderBy" example:"asc"`
		UsernameKeyword           string `form:"usernameKeyword" example:"jdoe"`
		ReferredByUsernameKeyword string `form:"referredByUsernameKeyword" example:"jdoe"`
		Cursor                    uint64 `form:"cursor" example:"5065"`
		Limit                     uint64 `form:"limit" example:"5000"`
	}

	PendingReview struct {
		CreatedAt          *time.Time `json:"time" swaggertype:"string" example:"2022-01-03T16:20:52.156534Z"`
		Iceflakes          string     `json:"iceflakes" swaggertype:"string" example:"100000000000000"`
		Username           string     `json:"username" swaggertype:"string" example:"myusername"`
		ReferredByUsername string     `json:"referredByUsername" swaggertype:"string" example:"myrefusername"`
		UserID             string     `json:"userId" swaggertype:"string" example:"12746386-03de-44d7-91c7-856fa66b6ed6"`
		EthAddress         string     `json:"ethAddress" swaggertype:"string" example:"0x43...."`
		Ice                float64    `json:"ice" db:"-" example:"1000"`
		IceInternal        int64      `json:"-" db:"ice" swaggerignore:"true"`
	}

	ByEarnerForReview struct {
		CreatedAt          *time.Time
		Username           string
		ReferredByUsername string
		UserID             string
		EarnerUserID       string
		EthAddress         string
		InternalID         int64
		Balance            float64
	}
)

// Private API.

const (
	applicationYamlKey = "coin-distribution"
	requestDeadline    = 25 * stdlibtime.Second

	batchSize = 700

	gasPriceCacheTTL = stdlibtime.Minute

	workerActionRun      workerAction = 0
	workerActionBlocked  workerAction = 1
	workerActionDisabled workerAction = 2
	workerActionOnDemand workerAction = 3

	ethApiStatusNew      ethApiStatus = "NEW"
	ethApiStatusPending  ethApiStatus = "PENDING"
	ethApiStatusAccepted ethApiStatus = "ACCEPTED"
	ethApiStatusRejected ethApiStatus = "REJECTED"

	ethTxStatusSuccessful ethTxStatus = "SUCCESSFUL"
	ethTxStatusFailed     ethTxStatus = "FAILED"
	ethTxStatusPending    ethTxStatus = "PENDING"

	configKeyCoinDistributerEnabled  = "coin_distributer_enabled"
	configKeyCoinDistributerOnDemand = "coin_distributer_forced_execution"
	configKeyCoinDistributerGasLimit = "coin_distributer_gas_limit_units"
	configKeyCoinDistributerGasPrice = "coin_distributer_gas_price_override"
)

// .
var (
	//nolint:gochecknoglobals // Singleton & global config mounted only during bootstrap.
	cfg config
	//go:embed DDL.sql
	ddl                  string
	errNotEnoughData     = errors.New("not enough data")
	errClientUncoverable = errors.New("uncoverable error")
)

type (
	ethTxStatus  string
	ethApiStatus string
	workerAction uint
	gasGetter    interface {
		GetGasOptions(ctx context.Context) (price *big.Int, limit uint64, err error)
	}
	ethClient interface {
		SuggestGasPrice(ctx context.Context) (*big.Int, error)
		TransactionsStatus(ctx context.Context, hashes []*string) (statuses map[ethTxStatus][]string, err error)
		TransactionStatus(ctx context.Context, hash string) (status ethTxStatus, err error)
		Airdrop(ctx context.Context, chanID *big.Int, gas gasGetter, recipients []common.Address, amounts []*big.Int) (string, error)
		io.Closer
	}
	airDropper interface {
		AirdropToWallets(opts *bind.TransactOpts, recipients []common.Address, amounts []*big.Int) (*types.Transaction, error)
	}
	batchRecord struct {
		CreatedAt  *time.Time   `db:"created_at"`
		Day        *time.Time   `db:"day"`
		EthTX      *string      `db:"eth_tx"`
		UserID     string       `db:"user_id"`
		EthAddress string       `db:"eth_address"`
		EthStatus  ethApiStatus `db:"eth_status"`
		Iceflakes  string       `db:"iceflakes"`
		InternalID int64        `db:"internal_id"`
	}
	batch struct {
		ID      string
		TX      string
		Records []*batchRecord
	}
	databaseConfig struct {
		DB *storage.DB
	}
	coinTracker struct {
		*databaseConfig
		Client       ethClient
		Conf         *config
		Workers      *pond.WorkerPool
		CancelSignal chan struct{}
	}
	coinProcessor struct {
		*databaseConfig
		Client        ethClient
		Conf          *config
		WG            *sync.WaitGroup
		CancelSignal  chan struct{}
		gasPriceCache struct {
			price *big.Int
			time  *time.Time
			mu    *sync.RWMutex
		}
	}
	ethClientImpl struct {
		RPC        *ethclient.Client
		Mutex      *sync.Mutex
		Key        *ecdsa.PrivateKey
		AirDropper airDropper
	}
	coinDistributer struct {
		Client    ethClient
		DB        *storage.DB
		Processor *coinProcessor
		Tracker   *coinTracker
	}
	repository struct {
		cfg *config
		db  *storage.DB
	}
	config struct {
		AlertSlackWebhook string `yaml:"alert-slack-webhook" mapstructure:"alert-slack-webhook"`
		Environment       string `yaml:"environment"         mapstructure:"environment"`
		ReviewURL         string `yaml:"review-url"          mapstructure:"review-url"`
		Ethereum          struct {
			RPC             string `yaml:"rpc"             mapstructure:"rpc"`
			PrivateKey      string `yaml:"privateKey"      mapstructure:"private-key"`
			ContractAddress string `yaml:"contractAddress" mapstructure:"contract-address"`
			ChainID         int64  `yaml:"chainId"         mapstructure:"chain-id"`
		} `yaml:"ethereum" mapstructure:"ethereum"`
		StartHours  int   `yaml:"startHours"  mapstructure:"start-hours"`
		EndHours    int   `yaml:"endHours"    mapstructure:"end-hours"`
		Workers     int64 `yaml:"workers"     mapstructure:"workers"`
		Development bool  `yaml:"development" mapstructure:"development"`
	}
)
