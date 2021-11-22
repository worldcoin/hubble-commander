package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func SaveHistogramMeasurementInMilliseconds(
	duration *time.Duration,
	metric *prometheus.HistogramVec,
	labels prometheus.Labels,
) {
	metric.With(labels).Observe(float64(duration.Milliseconds()))
}
