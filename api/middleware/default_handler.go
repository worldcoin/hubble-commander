package middleware

import (
	"net/http"

	"github.com/Worldcoin/hubble-commander/metrics"
)

func DefaultHandler(next http.Handler, commanderMetrics *metrics.CommanderMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		countRequest(commanderMetrics)

		next.ServeHTTP(w, r)
	})
}
