package middleware

import (
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
)

func measureRequestDuration(start time.Time, commanderMetrics *metrics.CommanderMetrics) time.Duration {
	duration := time.Since(start).Round(time.Millisecond)
	commanderMetrics.APIRequestDuration.Observe(float64(duration.Milliseconds()))
	return duration
}
