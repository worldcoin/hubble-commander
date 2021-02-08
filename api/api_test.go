package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/assert"
)

type Result struct {
	JsonRpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  string `json:"result"`
}

func TestStartApiServer(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc": "2.0", "method": "hubble_getVersion", "id": "1"}`)
	req, err := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	cfg := config.Config{Version: "v0123"}
	server, err := getApiServer(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	actual := &Result{}
	err = json.Unmarshal(w.Body.Bytes(), actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := &Result{"2.0", "1", "v0123"}

	assert.Equal(t, expected, actual)
}
