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

	// Mempool
	MempoolSize prometheus.Gauge
}

func NewCommanderMetrics() *CommanderMetrics {
	commanderMetrics := &CommanderMetrics{
		registry: prometheus.NewRegistry(),
	}

	commanderMetrics.initializeAPIMetrics()
	commanderMetrics.initializeRollupLoopMetrics()
	commanderMetrics.initializeSyncingMetrics()
	commanderMetrics.initializeBlockchainMetrics()

	commanderMetrics.MempoolSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name: "mempool_size",
		Help: "number of pending transactions",
	})
	commanderMetrics.registry.MustRegister(commanderMetrics.MempoolSize)

	return commanderMetrics
}
