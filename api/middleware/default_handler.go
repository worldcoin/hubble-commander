package middleware

import (
	"net/http"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
)

func DefaultHandler(next http.Handler, commanderMetrics *metrics.CommanderMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		countRequest(commanderMetrics)

		start := time.Now()
		defer measureRequestDuration(start, commanderMetrics)

		next.ServeHTTP(w, r)
	})
}

func measureRequestDuration(start time.Time, commanderMetrics *metrics.CommanderMetrics) time.Duration {
	duration := time.Since(start).Round(time.Millisecond)
	commanderMetrics.ApiRequestDuration.Observe(float64(duration.Milliseconds()))
	return duration
}
