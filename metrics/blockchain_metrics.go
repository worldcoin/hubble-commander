package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (c *CommanderMetrics) initializeBlockchainMetrics() {
	blockchainCallDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: blockchainSubsystem,
			Name:      "method_call_duration_milliseconds",
			Help:      "Durations of blockchain calls",
			Buckets: []float64{
				0.0,
				1.0,
				2.0,
				3.0,
				4.0,
				5.0,
				10.0,
				15.0,
				20.0,
				30.0,
				40.0,
				50.0,
				100.0,
			},
		},
		[]string{"method"},
	)

	c.registry.MustRegister(
		blockchainCallDuration,
	)

	c.BlockchainCallDuration = blockchainCallDuration
}

func (c *CommanderMetrics) SaveBlockchainCallDurationMeasurement(
	duration time.Duration,
	contractEventMetricsLabel string,
) {
	c.BlockchainCallDuration.
		With(prometheus.Labels{
			"method": contractEventMetricsLabel,
		}).
		Observe(float64(duration.Milliseconds()))
}
