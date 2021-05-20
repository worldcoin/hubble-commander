package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type payload struct {
	JSONRpc interface{}   `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("API: failed to read request body: %s", err)
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		defer logRequest(body, start)
		next.ServeHTTP(w, r)
	})
}

func logRequest(body []byte, start time.Time) {
	var decoded payload
	err := json.Unmarshal(body, &decoded)
	if err != nil {
		log.Printf("API: failed to unmarshal request body: %s", err)
		return
	}
	log.Printf("API: method: %v, duration: %v", decoded.Method, time.Since(start))
}
