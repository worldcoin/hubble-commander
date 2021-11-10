package middleware

import "github.com/Worldcoin/hubble-commander/metrics"

func countRequest(commanderMetrics *metrics.CommanderMetrics) {
	commanderMetrics.APITotalRequestsCounter.Inc()
}
