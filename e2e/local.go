package e2e

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/ybbus/jsonrpc/v2"
)

type E2ELocalCommander struct {
	client jsonrpc.RPCClient
}

func (e *E2ELocalCommander) Start() error {
	return nil
}

func (e *E2ELocalCommander) Stop() error {
	return nil
}

func (e *E2ELocalCommander) Client() jsonrpc.RPCClient {
	return e.client
}

func ConnectToLocalCommander() *E2ELocalCommander {
	cfg := config.GetConfig()
	endpoint := fmt.Sprintf("http://localhost:%s", cfg.API.Port)
	client := jsonrpc.NewClient(endpoint)
	return &E2ELocalCommander{client}
}
