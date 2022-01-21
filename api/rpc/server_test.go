package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type Result struct {
	// nolint:revive,stylecheck
	JsonRpc string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result"`
}

type API struct{}

func (a *API) GetAuthKey(ctx context.Context) interface{} {
	return ctx.Value(AuthKey)
}

func TestServer_ServeHTTP(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc": "2.0", "method": "test_getAuthKey", "id": "1"}`)
	req, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(jsonStr))
	require.NoError(t, err)

	keyValue := "secret key"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authKeyHeader, keyValue)

	server := NewServer()
	err = server.RegisterName("test", &API{})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	actual := &Result{}
	err = json.Unmarshal(w.Body.Bytes(), actual)
	require.NoError(t, err)
	require.Equal(t, keyValue, actual.Result)
}
