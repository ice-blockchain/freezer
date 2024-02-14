// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"

	detailedCoinMetrics "github.com/ice-blockchain/freezer/tokenomics/detailed_coin_metrics"
	appCfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
)

type mockedDetailedCoinMetrics struct{}

func (*mockedDetailedCoinMetrics) ReadDetails(context.Context) (*detailedCoinMetrics.Details, error) {
	return &detailedCoinMetrics.Details{
		CurrentPrice: 1.42,
		Volume24h:    42.42,
	}, nil
}

func helperCreateRepoWithRedisOnly(t *testing.T) *repository {
	t.Helper()

	var cfg Config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	db := storage.MustConnect(context.TODO(), applicationYamlKey)
	repo := &repository{
		cfg: &cfg,
		shutdown: func() error {
			return multierror.Append(db.Close()).ErrorOrNil()
		},
		db:                  db,
		detailedMetricsRepo: new(mockedDetailedCoinMetrics),
	}

	return repo
}

func TestGetCoinStatsBlockchainDetails(t *testing.T) {
	t.Parallel()

	repo := helperCreateRepoWithRedisOnly(t)

	t.Run("InvalidConfig", func(t *testing.T) {
		repo.cfg.CoinStats.RefreshInterval = 0
		require.Panics(t, func() {
			repo.keepBlockchainDetailsCacheUpdated(context.Background())
		})
	})

	t.Run("ReadFromEmptyCache", func(t *testing.T) {
		_, err := repo.db.Del(context.TODO(), totalCoinStatsDetailsKey).Result()
		require.NoError(t, err)

		data, err := repo.loadCachedBlockchainDetails(context.TODO())
		require.NoError(t, err)
		require.Nil(t, data)
	})

	t.Run("FillFromKeeper", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		repo.cfg.CoinStats.RefreshInterval = time.Minute
		repo.keepBlockchainDetailsCacheUpdated(ctx)
	})

	t.Run("CheckTimestampNoUpdate", func(t *testing.T) {
		err := repo.updateCachedBlockchainDetails(context.TODO())
		require.NoError(t, err)
	})

	t.Run("ReadCache", func(t *testing.T) {
		data, err := repo.loadCachedBlockchainDetails(context.TODO())
		require.NoError(t, err)
		require.NotNil(t, data)
		require.Equal(t, 1.42, data.CurrentPrice)
		require.Equal(t, 42.42, data.Volume24h)
		require.NotNil(t, data.Timestamp)
	})

	require.NoError(t, repo.Close())
}
