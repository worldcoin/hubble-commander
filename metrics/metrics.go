package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics and labels naming conventions https://prometheus.io/docs/practices/naming/.
type CommanderMetrics struct {
	registry *prometheus.Registry

	// API
	APITotalRequests           prometheus.Counter
	APIRequestDuration         prometheus.Histogram
	APITotalTransactions       *prometheus.CounterVec
	APITotalFailedTransactions prometheus.Counter
}

func NewCommanderMetrics() *CommanderMetrics {
	registry := prometheus.NewRegistry()

	commanderMetrics := &CommanderMetrics{
		registry: registry,
	}

	commanderMetrics.initializeMetricsForAPI()

	return commanderMetrics
}
