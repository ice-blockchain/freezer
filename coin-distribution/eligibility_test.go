// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"github.com/ice-blockchain/wintr/time"
	"github.com/stretchr/testify/assert"
	"testing"
	stdlibtime "time"
)

func TestIsEligibleForEthereumDistributionNow(t *testing.T) {
	t.Parallel()
	coinDistributionStartDate := time.New(stdlibtime.Date(2024, 1, 16, 0, 0, 0, 0, stdlibtime.UTC))
	lastCoinDistributionProcessedAt := time.New(stdlibtime.Date(2024, 1, 19, 0, 0, 0, 0, stdlibtime.UTC))

	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 22, 0, 0, 0, 0, stdlibtime.UTC))
	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 21, 0, 0, 0, 0, stdlibtime.UTC))
	assert.False(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.False(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))

	lastCoinDistributionProcessedAt = time.New(stdlibtime.Date(2024, 1, 16, 0, 0, 0, 0, stdlibtime.UTC))
	assert.True(t, isEligibleForEthereumDistributionNow(0, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(1, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(2, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(3, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(4, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
	assert.True(t, isEligibleForEthereumDistributionNow(5, 24*stdlibtime.Hour, 24*28*stdlibtime.Hour, time.Now(), coinDistributionStartDate, lastCoinDistributionProcessedAt))
}
