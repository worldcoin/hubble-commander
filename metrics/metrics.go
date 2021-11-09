package metrics

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics naming conventions https://prometheus.io/docs/practices/naming/
type CommanderMetrics struct {
	registry *prometheus.Registry

	// API
	ApiTotalRequestsCounter prometheus.Counter
	ApiRequestDuration      prometheus.Histogram
	ApiTotalTransactions    *prometheus.CounterVec
}

func (c *CommanderMetrics) NewMetricsServer(cfg *config.MetricsConfig) *http.Server {
	handler := promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{})

	mux := http.NewServeMux()
	mux.Handle(cfg.Endpoint, handler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	return &http.Server{Addr: addr, Handler: mux}
}

func NewCommanderMetrics() *CommanderMetrics {
	apiTotalRequestsCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "general",
		Name:      "requests_total",
		Help:      "Number of total requests made to the commander API",
	})

	apiRequestDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "api",
		Subsystem: "general",
		Name:      "request_duration_milliseconds",
		Help:      "Histogram of API requests duration",
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

	apiTotalTransactions := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "general",
		Name:      "transactions_total",
		Help:      "Total number of transactions sent to the commander API",
	},
		[]string{"type"},
	)

	// Makes total transactions metrics visible on the commander startup
	lowercaseTransfer := strings.ToLower(txtype.Transfer.String())
	apiTotalTransactions.With(prometheus.Labels{"type": lowercaseTransfer}).Add(0)
	lowercaseC2T := strings.ToLower(txtype.Create2Transfer.String())
	apiTotalTransactions.With(prometheus.Labels{"type": lowercaseC2T}).Add(0)

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		apiTotalRequestsCounter,
		apiRequestDuration,
		apiTotalTransactions,
	)

	return &CommanderMetrics{
		registry:                registry,
		ApiTotalRequestsCounter: apiTotalRequestsCounter,
		ApiRequestDuration:      apiRequestDuration,
		ApiTotalTransactions:    apiTotalTransactions,
	}
}

func (c *CommanderMetrics) CountTransaction(transactionType txtype.TransactionType) {
	lowercaseTxType := strings.ToLower(transactionType.String())
	c.ApiTotalTransactions.With(prometheus.Labels{"type": lowercaseTxType}).Inc()
}
