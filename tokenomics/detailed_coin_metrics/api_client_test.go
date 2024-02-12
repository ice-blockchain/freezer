// SPDX-License-Identifier: ice License 1.0

package detailed_coin_metrics //nolint:revive,nosnakecase,stylecheck //.

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiClientGetQuote(t *testing.T) {
	t.Parallel()

	key := os.Getenv("CMC_API_KEY")
	if key == "" {
		t.Skip("CMC_API_KEY is not set")
	}

	client := newAPIClient(key)
	require.NotNil(t, client)

	data, err := client.GetLatestQuote(context.Background(), iceSlug, targetCurrency)
	require.NoError(t, err)

	t.Logf("%+v", data)
	require.NotEmpty(t, data)
	require.Len(t, data, 1)
	require.Contains(t, data, iceID)
	require.Equal(t, iceSlug, data[iceID].Slug)
	require.Contains(t, data[iceID].Quote, targetCurrency)
	require.NotZero(t, data[iceID].Quote[targetCurrency].Price)
	require.NotZero(t, data[iceID].Quote[targetCurrency].Volume24h)
}

func TestApiClientBadCall(t *testing.T) {
	t.Parallel()

	client := newAPIClient("1234567890")
	require.NotNil(t, client)

	_, err := client.GetLatestQuote(context.Background(), iceSlug, targetCurrency)
	require.ErrorIs(t, err, ErrAPIFailed)
}
