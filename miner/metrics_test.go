// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"fmt"
	"testing"
	stdlibtime "time"

	"github.com/stretchr/testify/assert"
)

type (
	expectedFunc func(worker, batch, iteration uint64) bool
)

//nolint:funlen
func TestShouldSynchronizeBalance(t *testing.T) {
	t.Parallel()
	t.Run("workers is greater than 1 batch in total", func(t *testing.T) {
		t.Parallel()
		iterateOverBatches(t, new(telemetry).mustInit(config{Workers: 800}), 1, 1000, trueOncePerWorkerIteration(t, 800, 1))
	})
	t.Run("workers is greater than 1 batch in total and processing is slowed", func(t *testing.T) {
		t.Parallel()
		tel := slowTelemetry(800)
		iterateOverBatches(t, tel, 1, 1000, trueOncePerWorkerIteration(t, 800, 1))
	})
	t.Run("workers is greater than 2 batch in total", func(t *testing.T) {
		t.Parallel()
		iterateOverBatches(t, new(telemetry).mustInit(config{Workers: 800}), 2, 1000, trueOncePerWorkerIteration(t, 800, 2))
	})
	t.Run("one iteration and multiple batches", func(t *testing.T) {
		t.Parallel()
		tel := slowTelemetry(800)
		iterateOverBatches(t, tel, 9, 1000, trueOncePerWorkerIterationPerBatch(t, 9, 800, 1, 9))
	})
	t.Run("multiple iterations and batches per worker", func(t *testing.T) {
		t.Parallel()
		tel := slowTelemetry(10)
		iterateOverBatches(t, tel, 25, 1000, trueOncePerWorkerIterationPerBatch(t, 25, 50, 5, 5))
	})
	t.Run("only one worker but a lot of batches", func(t *testing.T) {
		t.Parallel()
		iterateOverBatches(t, new(telemetry).mustInit(config{Workers: 1}), 100, 1000, trueOncePerWorkerIteration(t, 1, 100))
	})
	t.Run("only one worker and processing is slowed", func(t *testing.T) {
		t.Parallel()
		tel := slowTelemetry(1)
		iterateOverBatches(t, tel, 59, 1000, trueOncePerWorkerIteration(t, 1, 59))
	})
	t.Run("every batch is processed at least once", func(t *testing.T) {
		t.Parallel()
		maxWorkers := int64(70)
		tel := new(telemetry).mustInit(config{Workers: maxWorkers})
		tel.collectElapsed(0, stdlibtime.Now().Add(-2*stdlibtime.Second))
		count := 0
		for w := uint64(0); w < uint64(maxWorkers); w++ {
			for i := uint64(0); i < uint64(10000); i++ {
				shouldSync := tel.shouldSynchronizeBalanceFunc(w, 3, i)
				if shouldSync(0) {
					count += 1
				}
				if shouldSync(1) {
					count += 1
				}
				if shouldSync(2) {
					count += 1
				}
			}
		}
		assert.Equal(t, 10000, count)
	})
	t.Run("previous workers waits in queue until all next is processed", func(t *testing.T) {
		t.Parallel()
		tel := slowTelemetry(400)
		checkPerWorkerAndIteration(t, tel, 0, 2, map[uint64]bool{
			0: true, 1: false, 2: false, 3: false, 4: false,
			400: true, 401: false, 402: false,
			799: false,
			800: true, 801: false,
		})
		checkPerWorkerAndIteration(t, tel, 1, 2, map[uint64]bool{
			0: false, 1: true, 2: false, 3: false, 4: false,
			400: false, 401: true, 402: false,
			800: false, 801: true, 802: false, 803: false,
		})
		checkPerWorkerAndIteration(t, tel, 163, 2, map[uint64]bool{
			0: false, 1: false, 2: false, 3: false,
			162: false, 163: true, 164: false,
			400: false, 401: false, 402: false,
			562: false, 563: true,
			800: false, 801: false,
			963: true,
		})
		checkPerWorkerAndIteration(t, tel, 399, 2, map[uint64]bool{
			0: false, 1: false, 2: false, 3: false,
			398: false, 399: true, 400: false,
			797: false, 798: true, 799: false, 800: false,
		})
	})
}

func checkPerWorkerAndIteration(tb testing.TB, tel *telemetry, worker, totalBatches uint64, iterations map[uint64]bool) {
	tb.Helper()
	for i := uint64(0); i < uint64(len(iterations)); i++ {
		shouldSync := tel.shouldSynchronizeBalanceFunc(worker, totalBatches, i)
		for b := uint64(0); b < totalBatches; b++ {
			assert.Equal(tb, iterations[i], shouldSync(b),
				"iteration %v on worker %v should be %v (batch %v)", i, worker, iterations[i], b)
		}
	}

}

func iterateOverBatches(t testing.TB, tel *telemetry, totalBatches, iterations uint64, expected expectedFunc) {
	t.Helper()
	for w := uint64(0); w < uint64(tel.cfg.Workers); w++ {
		for i := uint64(0); i < iterations; i++ {
			shouldSyncBalance := tel.shouldSynchronizeBalanceFunc(w, totalBatches, i)
			for b := uint64(0); b < totalBatches; b++ {
				assert.Equal(t, expected(w, b, i), shouldSyncBalance(b), fmt.Sprintf("worker %v, batch %v, iteration %v", w, b, i))
			}
		}
	}
}

func trueOncePerWorkerIterationPerBatch(tb testing.TB, totalBatches, totalIterations, iterationsPerWorker, batchesPerIteration uint64) expectedFunc {
	tb.Helper()
	iterations := make(map[string]bool)
	return func(worker, batch, iteration uint64) bool {
		iterationMatch := iteration%(totalIterations) >= (worker*iterationsPerWorker) &&
			iteration%(totalIterations) < ((worker+1)*iterationsPerWorker)
		maxBatches := ((iteration + 1) % iterationsPerWorker * batchesPerIteration)
		if maxBatches == 0 {
			maxBatches = totalBatches
		}
		batchMatch := batch >= ((iteration%iterationsPerWorker)*batchesPerIteration)%totalBatches && batch < maxBatches

		res := iterationMatch && batchMatch
		if res {
			key := fmt.Sprintf("%v~%v", iteration, batch)
			_, dupl := iterations[key]
			assert.False(tb, dupl, fmt.Sprintf("duplicated true for call on iteration %v and batch %v (worker %v)", iteration, batch, worker))
			iterations[key] = true
		}

		return res
	}
}

func trueOncePerWorkerIteration(tb testing.TB, totalWorkers, totalBatches uint64) expectedFunc {
	tb.Helper()
	iterations := make(map[uint64]bool)
	return func(worker, batch, iteration uint64) bool {
		res := (worker*totalBatches+batch)%(totalWorkers*totalBatches) == iteration%(totalWorkers*totalBatches)
		if res {
			_, dupl := iterations[iteration]
			assert.False(tb, dupl, fmt.Sprintf("%v", iteration))
			iterations[iteration] = true
		}
		return res
	}
}

func slowTelemetry(workers int64) *telemetry {
	tel := new(telemetry).mustInit(config{Workers: workers})
	tel.collectElapsed(0, stdlibtime.Now().Add(-60*stdlibtime.Second))
	tel.collectElapsed(1, stdlibtime.Now().Add(-50*stdlibtime.Second))
	tel.collectElapsed(2, stdlibtime.Now().Add(-40*stdlibtime.Second))
	tel.collectElapsed(3, stdlibtime.Now().Add(-30*stdlibtime.Second))
	tel.collectElapsed(4, stdlibtime.Now().Add(-20*stdlibtime.Second))
	tel.collectElapsed(5, stdlibtime.Now().Add(-10*stdlibtime.Second))
	tel.collectElapsed(6, stdlibtime.Now().Add(-1*stdlibtime.Second))
	tel.collectElapsed(7, stdlibtime.Now().Add(-1*stdlibtime.Second))
	tel.collectElapsed(8, stdlibtime.Now().Add(-1*stdlibtime.Second))

	return tel
}
