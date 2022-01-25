package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (c *CommanderMetrics) initializeAPIMetrics() {
	totalRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: apiSubsystem,
		Name:      "requests_total",
		Help:      "Number of requests made to the commander API",
	})

	requestsDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: apiSubsystem,
		Name:      "request_duration_milliseconds",
		Help:      "API requests duration",
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
	})

	totalTransactions := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: apiSubsystem,
			Name:      "transactions_total",
			Help:      "Number of transactions sent to the commander API",
		},
		[]string{"type", "status"},
	)

	c.registry.MustRegister(
		totalRequests,
		requestsDuration,
		totalTransactions,
	)

	c.APITotalRequests = totalRequests
	c.APIRequestDuration = requestsDuration
	c.APITotalTransactions = totalTransactions
}
