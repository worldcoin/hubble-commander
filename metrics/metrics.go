package metrics

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics naming conventions https://prometheus.io/docs/practices/naming/
type CommanderMetrics struct {
	registry *prometheus.Registry

	// API
	ApiTotalRequestsCounter prometheus.Counter
	ApiRequestDuration      prometheus.Histogram
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

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		apiTotalRequestsCounter,
		apiRequestDuration,
	)

	return &CommanderMetrics{
		registry:                registry,
		ApiTotalRequestsCounter: apiTotalRequestsCounter,
		ApiRequestDuration:      apiRequestDuration,
	}
}
