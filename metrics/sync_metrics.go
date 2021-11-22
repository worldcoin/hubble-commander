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
				25.0,
				50.0,
				75.0,
				100.0,
				150.0,
				200.0,
				250.0,
				300.0,
				350.0,
				400.0,
				450.0,
				500.0,
				600.0,
				700.0,
				800.0,
				900.0,
				1000.0,
				2000.0,
				3000.0,
				4000.0,
				5000.0,
			},
		},
		[]string{"method"},
	)

	c.registry.MustRegister(
		syncMethodDuration,
	)

	c.SyncingMethodDuration = syncMethodDuration
}
