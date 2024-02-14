// SPDX-License-Identifier: ice License 1.0

package detailed_coin_metrics //nolint:revive,nosnakecase,stylecheck //.

import (
	"context"
	"net/http"
	"net/url"

	"github.com/imroc/req/v3"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/log"
)

func newAPIClient(key string) *apiClientImpl {
	return &apiClientImpl{
		Key: key,
	}
}

func fetchFromAPI[T any](ctx context.Context, key, target string) (T, error) {
	const retryCount = 3
	var response struct {
		Data   T                 `json:"data"`
		Status apiResponseStatus `json:"status"`
	}
	resp, err := req.DefaultClient().R().SetContext(ctx).SetRetryCount(retryCount).
		SetHeader("X-CMC_PRO_API_KEY", key).
		SetHeader("Accept", "application/json").
		SetRetryHook(func(resp *req.Response, err error) {
			if err != nil {
				log.Error(errors.Wrap(err, "API: fetch failed"))
			} else {
				log.Warn("API: fetch failed: unexpected status code: " + resp.Status)
			}
		}).
		SetRetryCondition(func(resp *req.Response, err error) bool {
			return !(err == nil && resp.GetStatusCode() == http.StatusOK)
		}).
		SetErrorResult(&response).SetSuccessResult(&response).Get(target)
	switch {
	case err != nil:
		return response.Data, errors.Wrap(err, "cannot fetch data")
	case response.Status.ErrorCode != 0:
		return response.Data, errors.Wrapf(ErrAPIFailed, "message: %s (code %d)", response.Status.ErrorMessage,
			response.Status.ErrorCode)
	case resp.StatusCode != http.StatusOK:
		return response.Data, errors.Wrapf(ErrAPIFailed, "unexpected status code %d", resp.StatusCode)
	}

	return response.Data, nil
}

func (a *apiClientImpl) GetLatestQuote(ctx context.Context, slug, currency string) (map[int]apiResponseQuoteData, error) {
	const targetURL = `https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest`

	parsed, err := url.Parse(targetURL)
	log.Panic(errors.Wrap(err, "cannot parse target URL")) //nolint:revive // False positive.

	query := parsed.Query()
	query.Set("slug", slug)
	if currency != "" {
		query.Set("convert", currency)
	}
	parsed.RawQuery = query.Encode()

	return fetchFromAPI[map[int]apiResponseQuoteData](ctx, a.Key, parsed.String())
}
