package metrics

import (
	"strings"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *CommanderMetrics) initializeMetricsForAPI() {
	apiTotalRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "general",
		Name:      "requests_total",
		Help:      "Number of requests made to the commander API",
	})

	apiRequestDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "api",
		Subsystem: "general",
		Name:      "request_duration_milliseconds",
		Help:      "API requests duration",
		Buckets:   []float64{0.0, 25.0, 50.0, 75.0, 100.0, 150.0, 200.0, 250.0, 300.0, 350.0, 400.0, 450.0, 500.0, 600.0, 700.0, 800.0, 900.0, 1000.0},
	})

	apiTotalTransactions := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "general",
		Name:      "transactions_total",
		Help:      "Number of transactions sent to the commander API",
	},
		[]string{"type"},
	)

	// Makes total transactions metrics visible on the commander startup.
	lowercaseTransfer := strings.ToLower(txtype.Transfer.String())
	apiTotalTransactions.With(prometheus.Labels{"type": lowercaseTransfer}).Add(0)
	lowercaseC2T := strings.ToLower(txtype.Create2Transfer.String())
	apiTotalTransactions.With(prometheus.Labels{"type": lowercaseC2T}).Add(0)

	apiTotalFailedTransactions := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "general",
		Name:      "failed_transactions_total",
		Help:      "Number of transactions sent to the API which haven't passed the sanitization/validation",
		// There's a change that this label is used incorrectly. Verify when adding more metrics.
		ConstLabels: prometheus.Labels{
			"type": "transfer|create2transfer",
		},
	})

	c.registry.MustRegister(
		apiTotalRequests,
		apiRequestDuration,
		apiTotalTransactions,
		apiTotalFailedTransactions,
	)

	c.ApiTotalRequestsCounter = apiTotalRequests
	c.ApiRequestDuration = apiRequestDuration
	c.ApiTotalTransactions = apiTotalTransactions
	c.ApiTotalFailedTransactions = apiTotalFailedTransactions
}
