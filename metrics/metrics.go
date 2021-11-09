package metrics

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CommanderMetrics struct {
	registry *prometheus.Registry
}

func NewMetricsServer(cfg *config.MetricsConfig) (*http.Server, *CommanderMetrics) {
	registry := prometheus.NewRegistry()

	metrics := &CommanderMetrics{registry: registry}

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	mux := http.NewServeMux()
	mux.Handle(cfg.Endpoint, handler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{Addr: addr, Handler: mux}

	return server, metrics
}
