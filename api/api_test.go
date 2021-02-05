package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Result struct {
	JsonRpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  string `json:"result"`
}

func TestApiServer(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc": "2.0", "method": "hubble_sayHello", "id": "1"}`)
	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	server, _ := getApiServer()
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	actual := new(Result)
	_ = json.Unmarshal(w.Body.Bytes(), actual)

	expected := &Result{"2.0", "1", "Hello World!"}

	assert.Equal(t, expected, actual)
}
