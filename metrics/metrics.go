package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics and labels naming conventions https://prometheus.io/docs/practices/naming/.
type CommanderMetrics struct {
	registry *prometheus.Registry

	// API
	APITotalRequests     prometheus.Counter
	APIRequestDuration   prometheus.Histogram
	APITotalTransactions *prometheus.CounterVec

	// Rollup
	CommitmentBuildDuration         *prometheus.HistogramVec
	BatchBuildAndSubmissionDuration *prometheus.HistogramVec

	// Syncing
	SyncingMethodDuration *prometheus.HistogramVec

	// Blockchain
	BlockchainCallDuration *prometheus.HistogramVec
}

func NewCommanderMetrics() *CommanderMetrics {
	commanderMetrics := &CommanderMetrics{
		registry: prometheus.NewRegistry(),
	}

	commanderMetrics.initializeAPIMetrics()
	commanderMetrics.initializeRollupLoopMetrics()
	commanderMetrics.initializeSyncingMetrics()
	commanderMetrics.initializeBlockchainMetrics()

	return commanderMetrics
}
