// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"fmt"
	stdlog "log"
	"strings"
	stdlibtime "time"

	"github.com/rcrowley/go-metrics"
)

func init() {
	stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime | stdlog.Lmsgprefix | stdlog.LUTC | stdlog.Lmicroseconds)
}

type telemetry struct {
	registry metrics.Registry
	steps    [8]string
}

func (t *telemetry) mustInit() *telemetry {
	const (
		decayAlpha    = 0.015
		reservoirSize = 1000_000
	)
	t.registry = metrics.NewRegistry()
	t.steps = [8]string{"mine[full iteration]", "mine", "get_users", "get_referrals", "send_messages", "get_history", "insert_history", "update_users"}
	for ix := range &t.steps {
		if ix > 1 {
			t.steps[ix] = fmt.Sprintf("[%v]mine.%v", ix-1, t.steps[ix])
		}
		if err := t.registry.Register(t.steps[ix], metrics.NewCustomTimer(metrics.NewHistogram(metrics.NewExpDecaySample(reservoirSize, decayAlpha)), metrics.NewMeter())); err != nil { //nolint:lll // .
			panic(err)
		}
	}

	go metrics.LogScaled(t.registry, 60*stdlibtime.Minute, stdlibtime.Millisecond, t)

	return t
}

func (t *telemetry) collectElapsed(step int, since stdlibtime.Time) {
	t.registry.Get(t.steps[step]).(metrics.Timer).UpdateSince(since)
}

func (*telemetry) Printf(format string, args ...interface{}) {
	stdlog.Printf(strings.ReplaceAll(format, "timer ", ""), args...)
}
