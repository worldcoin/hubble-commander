package middleware

import (
	"net/http"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
)

func DefaultHandler(next http.Handler, commanderMetrics *metrics.CommanderMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commanderMetrics.APITotalRequests.Inc()

		start := time.Now()
		defer measureRequestDuration(start, commanderMetrics)

		next.ServeHTTP(w, r)
	})
}
