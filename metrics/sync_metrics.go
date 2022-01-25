package metrics

import "github.com/prometheus/client_golang/prometheus"

func (c *CommanderMetrics) initializeSyncingMetrics() {
	syncMethodDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: syncingSubsystem,
			Name:      "method_duration_milliseconds",
			Help:      "Durations of syncing methods",
			Buckets: []float64{
				0.0,
				5.0,
				10.0,
				25.0,
				50.0,
				75.0,
				100.0,
				150.0,
				200.0,
				250.0,
			},
		},
		[]string{"method"},
	)

	c.registry.MustRegister(
		syncMethodDuration,
	)

	c.SyncingMethodDuration = syncMethodDuration
}
