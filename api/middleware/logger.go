package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/utils"
	log "github.com/sirupsen/logrus"
)

var disabledAPIMethods = []string{
	"getVersion",
	"getNetworkInfo",
}

type payload struct {
	Method string `json:"method"`
}

func Logger(next http.Handler, commanderMetrics *metrics.CommanderMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commanderMetrics.APITotalRequests.Inc()

		start := time.Now()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("API: failed to read request body: %s", err)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		defer func() {
			duration := metrics.MeasureDuration(start, commanderMetrics.APIRequestDuration)
			logRequest(body, duration)
		}()
		next.ServeHTTP(w, r)
	})
}

func logRequest(body []byte, duration time.Duration) {
	var decoded payload

	err := json.Unmarshal(body, &decoded)
	if err != nil {
		logBatchRequest(body, duration)
		return
	}

	if shouldMethodBeLogged(decoded.Method) {
		log.Debugf("API: method: %v, duration: %v", decoded.Method, duration.String())
	}
}

func shouldMethodBeLogged(method string) bool {
	split := strings.Split(method, "_")

	if len(split) < 2 {
		return true
	}

	return !utils.StringInSlice(split[1], disabledAPIMethods)
}

func logBatchRequest(body []byte, duration time.Duration) {
	var decoded []payload
	err := json.Unmarshal(body, &decoded)
	if err != nil {
		log.Errorf("API: failed to unmarshal request body: %s", err)
		return
	}
	methodsArray := extractMethodNames(decoded)
	log.Debugf("API: batch call, methods: %s, duration: %v", methodsArray, duration.String())
}

func extractMethodNames(decoded []payload) string {
	methods := make([]string, 0, len(decoded))
	for i := range decoded {
		if decoded[i].Method == "" {
			decoded[i].Method = "invalid request"
		}
		methods = append(methods, decoded[i].Method)
	}
	return "[" + strings.Join(methods, ", ") + "]"
}
