package metrics

import "github.com/prometheus/client_golang/prometheus"

func (c *CommanderMetrics) initializeSyncingMetrics() {
	c.initializeAccountSyncingMetrics()
}

func (c *CommanderMetrics) initializeAccountSyncingMetrics() {
	syncingDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: syncingSubsystem,
		Name:      "duration_milliseconds",
		Help:      "Account syncing duration (single accounts and batch accounts)",
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
		},
	})

	c.registry.MustRegister(
		syncingDuration,
	)

	c.SyncingAccountsDuration = syncingDuration
}
