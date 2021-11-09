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
	ApiTotalRequestsCounter *prometheus.Counter
}

func NewMetricsServer(cfg *config.MetricsConfig) (*http.Server, *prometheus.Registry) {
	registry := prometheus.NewRegistry()

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	mux := http.NewServeMux()
	mux.Handle(cfg.Endpoint, handler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{Addr: addr, Handler: mux}

	return server, registry
}

func NewCommanderMetrics(metrics *CommanderMetrics, registry *prometheus.Registry) {
	apiTotalRequestsCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "general",
		Name:      "requests_total",
		Help:      "Number of total requests made to the commander API",
	})

	registry.MustRegister(
		apiTotalRequestsCounter,
	)

	metrics.registry = registry
	metrics.ApiTotalRequestsCounter = &apiTotalRequestsCounter
}

func IncrementCounter(counter *prometheus.Counter) {
	dereferencedCounter := *counter
	dereferencedCounter.Inc()
}
