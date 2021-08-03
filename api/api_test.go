package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
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

	storage, err := st.NewTestStorage()
	require.NoError(t, err)
	client, err := eth.NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	cfg := config.APIConfig{Version: "v0123"}
	server, err := getAPIServer(&cfg, storage.Storage, client.Client, false)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	actual := &Result{}
	err = json.Unmarshal(w.Body.Bytes(), actual)
	require.NoError(t, err)

	require.Equal(t, &Result{"2.0", "1", "v0123"}, actual)
}
