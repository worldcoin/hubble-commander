package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func MeasureDuration(start time.Time, histogram prometheus.Histogram) time.Duration {
	duration := time.Since(start).Round(time.Millisecond)
	histogram.Observe(float64(duration.Milliseconds()))
	return duration
}
