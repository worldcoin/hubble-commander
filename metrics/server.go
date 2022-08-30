package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (c *CommanderMetrics) NewServer(cfg *config.MetricsConfig) *http.Server {
	handler := promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{})

	mux := http.NewServeMux()
	mux.Handle(cfg.Endpoint, handler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	return &http.Server{
		ReadHeaderTimeout: time.Second * 5,
		Addr:              addr,
		Handler:           mux,
	}
}
