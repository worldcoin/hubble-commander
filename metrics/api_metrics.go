package metrics

import (
	"strings"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultAPIMetricsSubsystem = "general"
)

func (c *CommanderMetrics) initializeMetricsForAPI() {
	totalRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: DefaultAPIMetricsSubsystem,
		Name:      "requests_total",
		Help:      "Number of requests made to the commander API",
	})

	requestsDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: DefaultAPIMetricsSubsystem,
		Name:      "request_duration_milliseconds",
		Help:      "API requests duration",
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

	totalTransactions := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: DefaultAPIMetricsSubsystem,
		Name:      "transactions_total",
		Help:      "Number of transactions sent to the commander API",
	},
		[]string{"type"},
	)

	// Makes total transactions metrics visible on the commander startup.
	lowercaseTransferType := strings.ToLower(txtype.Transfer.String())
	totalTransactions.With(prometheus.Labels{"type": lowercaseTransferType}).Add(0)
	lowercaseC2TType := strings.ToLower(txtype.Create2Transfer.String())
	totalTransactions.With(prometheus.Labels{"type": lowercaseC2TType}).Add(0)

	totalFailedTransactions := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: DefaultAPIMetricsSubsystem,
		Name:      "failed_transactions_total",
		Help:      "Number of transactions sent to the API which haven't passed the sanitization/validation",
		// There's a chance that this label is used incorrectly. Verify when adding more metrics.
		ConstLabels: prometheus.Labels{
			"type": "transfer|create2transfer",
		},
	})

	c.registry.MustRegister(
		totalRequests,
		requestsDuration,
		totalTransactions,
		totalFailedTransactions,
	)

	c.APITotalRequests = totalRequests
	c.APIRequestDuration = requestsDuration
	c.APITotalTransactions = totalTransactions
	c.APITotalFailedTransactions = totalFailedTransactions
}
