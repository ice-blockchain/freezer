// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	stdlibtime "time"

	"github.com/alitto/pond"
	"github.com/bsm/redislock"
	"github.com/goccy/go-json"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetTotalCoinsSummary(ctx context.Context, days uint64, _ stdlibtime.Duration) (*TotalCoinsSummary, error) {
	var (
		dates []stdlibtime.Time
		res   = new(TotalCoinsSummary)
		now   = time.Now()
	)

	dates, res.TimeSeries = r.totalCoinsDates(now, days)
	totalCoins, err := r.getCachedTotalCoins(ctx, dates)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getCachedTotalCoins for createdAts:%#v", dates)
	}
	for _, child := range res.TimeSeries {
		for _, stats := range totalCoins {
			if stats.CreatedAt.Equal(child.Date) {
				child.Standard = stats.BalanceTotalStandard
				child.PreStaking = stats.BalanceTotalPreStaking
				child.Blockchain = stats.BalanceTotalEthereum
				child.Total = stats.BalanceTotal
				break
			}
		}
		child.Date = child.Date.Add(-1 * stdlibtime.Nanosecond)

	}
	res.TotalCoins = res.TimeSeries[0].TotalCoins
	res.BlockchainDetails = r.detailedMetricsData.Load()

	return r.enhanceWithBlockchainCoinStats(res), nil
}

func (r *repository) totalCoinsDates(now *time.Time, days uint64) ([]stdlibtime.Time, []*TotalCoinsTimeSeriesDataPoint) {
	var (
		truncationInterval = r.cfg.GlobalAggregationInterval.Child
		dates              = make([]stdlibtime.Time, 0, days)
		timeSeries         = make([]*TotalCoinsTimeSeriesDataPoint, 0, days)
		dayInterval        = r.cfg.GlobalAggregationInterval.Parent
		start              = now.Add(-1 * truncationInterval).Truncate(truncationInterval)
	)
	dates = append(dates, start)
	timeSeries = append(timeSeries, &TotalCoinsTimeSeriesDataPoint{Date: start})
	for day := uint64(0); day < days-1; day++ {
		date := now.Add(dayInterval * -1 * stdlibtime.Duration(day)).Truncate(dayInterval)
		dates = append(dates, date)
		timeSeries = append(timeSeries, &TotalCoinsTimeSeriesDataPoint{Date: date})
	}

	return dates, timeSeries
}

func (r *repository) cacheTotalCoins(ctx context.Context, coins []*dwh.TotalCoins) error {
	val := make([]interface{ Key() string }, 0, len(coins))
	for _, v := range coins {
		val = append(val, v)
	}

	return errors.Wrapf(storage.Set(ctx, r.db, val...), "failed to set cache value for total coins: %#v", coins)
}

func (r *repository) getCachedTotalCoins(ctx context.Context, dates []stdlibtime.Time) ([]*dwh.TotalCoins, error) {
	keys := make([]string, 0, len(dates))
	for _, d := range dates {
		keys = append(keys, r.totalCoinsCacheKey(d))
	}
	cached, err := storage.Get[dwh.TotalCoins](ctx, r.db, keys...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cached coinStats for dates %#v", dates)
	}

	return cached, nil
}

func (r *repository) totalCoinsCacheKey(date stdlibtime.Time) string {
	return (&dwh.TotalCoins{CreatedAt: time.New(date.Truncate(r.cfg.GlobalAggregationInterval.Child))}).Key()
}

func (r *repository) updateDetailedCoinMetrics(ctx context.Context) error {
	detailedCoinMetrics, err := r.detailedMetricsRepo.ReadDetails(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to read detailedCoinMetrics")
	}

	r.detailedMetricsData.Store(detailedCoinMetrics)

	return nil
}

func (r *repository) keepBlockchainDetailsCacheUpdated(ctx context.Context) {
	signals := make(chan struct{}, 1)
	ticker := stdlibtime.NewTicker(r.cfg.GlobalAggregationInterval.Child / 6) //nolint:gomnd // Not an issue.
	defer ticker.Stop()

	// Send initial signal now without waiting for the first tick.
	signals <- struct{}{}

	go func() {
		for range ticker.C {
			select {
			case signals <- struct{}{}:
			default:
			}
		}
	}()

	for {
		select {
		case <-signals:
			reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
			if err := r.updateDetailedCoinMetrics(reqCtx); err != nil {
				log.Error(errors.Wrap(err, "failed to update detailedCoinMetrics"))
			}
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

func (r *repository) keepTotalCoinsCacheUpdated(ctx context.Context, initialNow *time.Time) {
	ticker := stdlibtime.NewTicker(stdlibtime.Duration(1+rand.Intn(10)) * (r.cfg.GlobalAggregationInterval.Child / 60)) //nolint:gosec,gomnd // Not an  issue.
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var (
				now                    = time.Now()
				newDate                = now.Truncate(r.cfg.GlobalAggregationInterval.Child)
				historyGenerationDelta = stdlibtime.Duration(float64(r.cfg.GlobalAggregationInterval.Child) * 0.75) //nolint:gomnd // .
			)
			lastDateCached, err := r.getLastDateCached(ctx)
			if err != nil {
				log.Error(errors.Wrapf(err, "failed to get last date cached"))
			}
			if lastDateCached.IsNil() || (!lastDateCached.Equal(newDate) && now.Sub(newDate) >= historyGenerationDelta) {
				dwhCtx, cancel := context.WithTimeout(ctx, 1*stdlibtime.Minute)
				lock, err := redislock.Obtain(dwhCtx, r.db, totalCoinStatsCacheLockKey, totalCoinStatsCacheLockDuration, &redislock.Options{RetryStrategy: redislock.NoRetry()})
				if err != nil && errors.Is(err, redislock.ErrNotObtained) {
					cancel()
					continue
				} else if err != nil {
					log.Error(errors.Wrapf(err, "failed to init total coin stats cache (aquire lock totalCoinStatsCache)")) //nolint:revive // Nope.
					cancel()
					continue
				}
				if err = r.buildTotalCoinCache(dwhCtx, newDate); err != nil {
					log.Error(errors.Wrapf(err, "failed to update total coin stats cache for date %v", *now.Time))
				} else {
					if err := r.setLastDateCached(ctx, time.New(newDate)); err != nil {
						log.Error(errors.Wrapf(err, "can't set last date cached: %v", time.New(newDate)))
					}
				}
				if lock != nil {
					log.Error(errors.Wrapf(lock.Release(dwhCtx), "error releasing lock, key: totalCoinStatsCache"))
				}
				cancel()
			}
		case <-ctx.Done():
			return
		}
	}
}

func (r *repository) buildTotalCoinCache(ctx context.Context, dates ...stdlibtime.Time) error {
	totalCoins, err := r.dwh.SelectTotalCoins(ctx, dates)
	if err != nil {
		return errors.Wrapf(err, "failed to read total coin stats cacheable values for dates %#v", dates)
	}

	return errors.Wrapf(
		r.cacheTotalCoins(ctx, totalCoins),
		"failed to save total coin stats cache for dates %#v", dates)
}

func (r *repository) mustInitTotalCoinsCache(ctx context.Context, now *time.Time) {
	dates, _ := r.totalCoinsDates(now, daysCountToInitCoinsCacheOnStartup)
	alreadyCached, err := r.getCachedTotalCoins(ctx, dates)
	log.Panic(errors.Wrapf(err, "failed to init total coin stats cache")) //nolint:revive // Nope.
	for _, cached := range alreadyCached {
		for dateIdx, date := range dates {
			if cached.CreatedAt.Equal(date) {
				dates = append(dates[:dateIdx], dates[dateIdx+1:]...)

				break
			}
		}
	}
	workerPool := pond.New(routinesCountToInitCoinsCacheOnStartup, 0, pond.MinWorkers(routinesCountToInitCoinsCacheOnStartup))
	if len(dates) > 0 {
		lockCtx, cancel := context.WithTimeout(ctx, 1*stdlibtime.Minute)
		defer cancel()
		lock, err := redislock.Obtain(lockCtx, r.db, totalCoinStatsCacheLockKey, totalCoinStatsCacheLockDuration, &redislock.Options{RetryStrategy: redislock.NoRetry()})
		if err != nil && errors.Is(err, redislock.ErrNotObtained) {
			return
		} else if err != nil {
			log.Panic(errors.Wrapf(err, "failed to init total coin stats cache (aquire lock totalCoinStatsCache)")) //nolint:revive // Nope.
		}
		defer func() {
			log.Error(errors.Wrapf(lock.Release(lockCtx), "error releasing lock, key: totalCoinStatsCache"))
		}()
		for _, date := range dates {
			fetchDate := date
			workerPool.Submit(func() {
				for err = errors.New("first try"); err != nil; {
					log.Info(fmt.Sprintf("Building total coins cache for `%v`", fetchDate))
					err = errors.Wrapf(r.buildTotalCoinCache(ctx, fetchDate), "failed to build/init total coins cache for %v", fetchDate)
					log.Error(err)
				}
			})
		}
		workerPool.StopAndWait()
	}
}

func (r *repository) setLastDateCached(ctx context.Context, lastDateCached *time.Time) error {
	_, err := r.db.Set(ctx, "total_coins_last_date_cached", lastDateCached, 0).Result()

	return errors.Wrapf(err, "failed to set total_coins_last_date_cached for: %v", lastDateCached)
}

func (r *repository) getLastDateCached(ctx context.Context) (lastDateCached *time.Time, err error) {
	lastDateCachedString, err := r.db.Get(ctx, "total_coins_last_date_cached").Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if lastDateCachedString != "" {
		lastDateCached = new(time.Time)
		if err := lastDateCached.UnmarshalText([]byte(lastDateCachedString)); err != nil {
			return nil, errors.Wrapf(err, "failed to parse total_coins_last_date_cached `%v`", lastDateCachedString)
		}
		lastDateCached = time.New(lastDateCached.UTC())

		return
	}

	return nil, errors.Wrap(err, "failed to get last date cached value")
}

func (r *repository) enhanceWithBlockchainCoinStats(res *TotalCoinsSummary) *TotalCoinsSummary {
	c := r.cfg.blockchainCoinStatsJSON.Load()
	if len(c.CoinsAddedHistory) == 0 {
		return res
	}
	cpy := *res
	for _, coinsAddedHistoryEntry := range c.CoinsAddedHistory {
		addedBeforeFirst := coinsAddedHistoryEntry.Date.Before(res.TimeSeries[0].Date)
		if addedBeforeFirst {
			cpy.Blockchain += coinsAddedHistoryEntry.CoinsAdded
		}
		for i := range res.TimeSeries {
			if coinsAddedHistoryEntry.Date.Before(res.TimeSeries[i].Date) {
				cpy.TimeSeries[i].Blockchain += coinsAddedHistoryEntry.CoinsAdded
			}
		}
	}

	return &cpy
}

func (r *repository) startBlockchainCoinStatsJSONSyncer(ctx context.Context) {
	ticker := stdlibtime.NewTicker(10 * stdlibtime.Minute) //nolint:gomnd // .
	defer ticker.Stop()
	r.cfg.blockchainCoinStatsJSON = new(atomic.Pointer[blockchainCoinStatsJSON])
	log.Panic(errors.Wrap(r.syncBlockchainCoinStatsJSON(ctx), "failed to syncBlockchainCoinStatsJSON"))

	for {
		select {
		case <-ticker.C:
			reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
			log.Error(errors.Wrap(r.syncBlockchainCoinStatsJSON(reqCtx), "failed to syncBlockchainCoinStatsJSON"))
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

//nolint:funlen,gomnd // .
func (r *repository) syncBlockchainCoinStatsJSON(ctx context.Context) error {
	if resp, err := req.
		SetContext(ctx).
		SetRetryCount(25).
		SetRetryBackoffInterval(10*stdlibtime.Millisecond, 1*stdlibtime.Second).
		SetRetryHook(func(resp *req.Response, err error) {
			if err != nil {
				log.Error(errors.Wrap(err, "failed to fetch BlockchainCoinStatsJSON, retrying..."))
			} else {
				log.Error(errors.Errorf("failed to fetch BlockchainCoinStatsJSON with status code:%v, retrying...", resp.GetStatusCode()))
			}
		}).
		SetRetryCondition(func(resp *req.Response, err error) bool {
			return err != nil || resp.GetStatusCode() != http.StatusOK
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Cache-Control", "no-cache, no-store, must-revalidate").
		SetHeader("Pragma", "no-cache").
		SetHeader("Expires", "0").
		Get(r.cfg.BlockchainCoinStatsJSONURL); err != nil {
		return errors.Wrapf(err, "failed to get fetch `%v`", r.cfg.BlockchainCoinStatsJSONURL)
	} else if data, err2 := resp.ToBytes(); err2 != nil {
		return errors.Wrapf(err2, "failed to read body of `%v`", r.cfg.BlockchainCoinStatsJSONURL)
	} else {
		var val blockchainCoinStatsJSON
		if err = json.UnmarshalContext(ctx, data, &val); err != nil {
			return errors.Wrapf(err, "failed to unmarshal into %#v, data: %v", val, string(data))
		}
		if body := string(data); !strings.Contains(body, "coinsAddedHistory") {
			return errors.Errorf("there's something wrong with the blockchainCoinStatsJSON body: %v", body)
		}
		r.cfg.blockchainCoinStatsJSON.Swap(&val)

		return nil
	}
}
