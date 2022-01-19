package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/stretchr/testify/require"
)

type Result struct {
	// nolint:revive,stylecheck
	JsonRpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  string `json:"result"`
}

func TestStartApiServer(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc": "2.0", "method": "hubble_getVersion", "id": "1"}`)
	req, err := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	cfg := config.APIConfig{Version: "v0123"}
	commanderMetrics := metrics.NewCommanderMetrics()
	server, err := getAPIServer(&cfg, nil, eth.DomainOnlyTestClient, commanderMetrics, false, func(enable bool) {})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	actual := &Result{}
	err = json.Unmarshal(w.Body.Bytes(), actual)
	require.NoError(t, err)

	require.Equal(t, &Result{"2.0", "1", "v0123"}, actual)
}
