// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"fmt"
	stdlog "log"
	"strings"
	stdlibtime "time"

	"github.com/rcrowley/go-metrics"

	"github.com/ice-blockchain/wintr/log"
)

func init() {
	stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime | stdlog.Lmsgprefix | stdlog.LUTC | stdlog.Lmicroseconds)
}

type telemetry struct {
	registry        metrics.Registry
	steps           [10]string
	currentStepName string
	cfg             config
}

func (t *telemetry) mustInit(cfg config) *telemetry {
	const (
		decayAlpha    = 0.015
		reservoirSize = 10_000
	)
	t.cfg = cfg
	t.registry = metrics.NewRegistry()
	t.steps = [10]string{"mine[full iteration]", "mine", "get_users", "get_referrals", "send_messages", "get_history", "sync_quiz_status", "insert_history", "collect_coin_distributions", "update_users"} //nolint:lll // .
	for ix := range &t.steps {
		if ix > 1 {
			t.steps[ix] = fmt.Sprintf("[%v]mine.%v", ix-1, t.steps[ix])
		}
		log.Panic(t.registry.Register(t.steps[ix], metrics.NewCustomTimer(metrics.NewHistogram(metrics.NewExpDecaySample(reservoirSize, decayAlpha)), metrics.NewMeter()))) //nolint:lll // .
	}

	go metrics.LogScaled(t.registry, 15*stdlibtime.Minute, stdlibtime.Millisecond, t) //nolint:gomnd // .

	return t
}

func (t *telemetry) collectElapsed(step int, since stdlibtime.Time) {
	t.registry.Get(t.steps[step]).(metrics.Timer).UpdateSince(since)
}

func (t *telemetry) shouldSynchronizeBalanceFunc(workerNumber, totalBatches, iteration uint64) func(batchNumber uint64) bool {
	var deadline float64
	if t.cfg.Development {
		deadline = float64(stdlibtime.Minute)
	} else {
		deadline = float64(stdlibtime.Hour)
	}
	timingPrevStep := t.registry.Get(t.steps[0]).(metrics.Timer).Percentile(0.99) // nolint:forcetypeassert
	targetIterations := uint64(deadline / timingPrevStep)
	targetIterations = (targetIterations / uint64(t.cfg.Workers)) * uint64(t.cfg.Workers)
	if targetIterations <= 0 {
		targetIterations = 1
	}
	iterationsOwnedBy1Worker := targetIterations / uint64(t.cfg.Workers)
	if iterationsOwnedBy1Worker <= 0 {
		iterationsOwnedBy1Worker = 1
	}
	batchesPerIterationsOwnedBy1Worker := totalBatches / iterationsOwnedBy1Worker
	if batchesPerIterationsOwnedBy1Worker <= 0 {
		batchesPerIterationsOwnedBy1Worker = 1
	}
	if totalBatches <= iterationsOwnedBy1Worker {
		iterationsOwnedBy1Worker = totalBatches
		if targetIterations >= 1 {
			targetIterations = iterationsOwnedBy1Worker * uint64(t.cfg.Workers)
		}
	}
	var (
		currentIteration = iteration % targetIterations
		left             = workerNumber * iterationsOwnedBy1Worker
		right            = (workerNumber + 1) * iterationsOwnedBy1Worker
	)
	if targetIterations == 1 {
		targetIterations = iterationsOwnedBy1Worker * uint64(t.cfg.Workers)
		if currentIteration == 0 {
			currentIteration = iteration % targetIterations
		}
	}
	if currentIteration < left || currentIteration >= right {
		return func(batchNumber uint64) bool {
			return false
		}
	}
	if t.cfg.Development {
		return func(batchNumber uint64) bool {
			return currentIteration == left
		}
	}

	return func(batchNumber uint64) bool {
		for i := uint64(0); i < batchesPerIterationsOwnedBy1Worker; i++ {
			if batchNumber == (((currentIteration - left) * batchesPerIterationsOwnedBy1Worker) + i) {
				return true
			}
		}
		if currentIteration == right-1 {
			for expectedBatchNumber := (batchesPerIterationsOwnedBy1Worker) * iterationsOwnedBy1Worker; expectedBatchNumber < totalBatches; expectedBatchNumber++ {
				if batchNumber == expectedBatchNumber {
					return true
				}
			}
		}

		return false
	}
}

func (t *telemetry) Printf(format string, args ...interface{}) {
	const prefixMarker = "timer "
	if strings.HasPrefix(format, prefixMarker) {
		t.currentStepName = fmt.Sprintf(format[len(prefixMarker):strings.IndexRune(format, '\n')], args...)
	}
	stdlog.Printf("["+t.currentStepName+"]"+strings.ReplaceAll(format, prefixMarker, ""), args...)
}
