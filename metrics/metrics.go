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
	LatestBlockNumber     prometheus.Gauge
	SyncedBlockNumber     prometheus.Gauge

	// Blockchain
	BlockchainCallDuration *prometheus.HistogramVec
	BlockchainGasSpend     prometheus.Counter

	// Mempool
	MempoolSize                prometheus.Gauge
	MempoolSizeTransfer        prometheus.Gauge
	MempoolSizeCreate2Transfer prometheus.Gauge
	MempoolSizeMassMigration   prometheus.Gauge
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
		Name:      "mempool_size",
		Help:      "number of pending transactions",
	})
	commanderMetrics.registry.MustRegister(commanderMetrics.MempoolSize)

	commanderMetrics.MempoolSizeTransfer = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "mempool_size_transfer",
		Help:      "number of pending transfers",
	})
	commanderMetrics.registry.MustRegister(commanderMetrics.MempoolSizeTransfer)

	commanderMetrics.MempoolSizeCreate2Transfer = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "mempool_size_create2transfer",
		Help:      "number of pending create2transfers",
	})
	commanderMetrics.registry.MustRegister(commanderMetrics.MempoolSizeCreate2Transfer)

	commanderMetrics.MempoolSizeMassMigration = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "mempool_size_massmigration",
		Help:      "number of pending massmigrations",
	})
	commanderMetrics.registry.MustRegister(commanderMetrics.MempoolSizeMassMigration)

	commanderMetrics.LatestBlockNumber = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: syncingSubsystem,
		Name:      "latest_block_number",
		Help:      "The most recent canonical block we know of",
	})
	commanderMetrics.registry.MustRegister(commanderMetrics.LatestBlockNumber)

	commanderMetrics.SyncedBlockNumber = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: syncingSubsystem,
		Name:      "synced_block_number",
		Help:      "The most recent block we have synced",
	})
	commanderMetrics.registry.MustRegister(commanderMetrics.SyncedBlockNumber)

	commanderMetrics.BlockchainGasSpend = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: blockchainSubsystem,
		Name:      "gas_spent",
		Help:      "How much gas we have spent",
	})
	commanderMetrics.registry.MustRegister(commanderMetrics.BlockchainGasSpend)

	return commanderMetrics
}
