package e2e

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/ybbus/jsonrpc/v2"
)

type LocalCommander struct {
	client jsonrpc.RPCClient
}

func ConnectToLocalCommander() *LocalCommander {
	cfg := config.GetViperConfig()
	endpoint := fmt.Sprintf("http://localhost:%s", cfg.API.Port)
	client := jsonrpc.NewClient(endpoint)
	return &LocalCommander{client}
}

func (e *LocalCommander) Start() error {
	return nil
}

func (e *LocalCommander) Stop() error {
	return nil
}

func (e *LocalCommander) Restart() error {
	return nil
}

func (e *LocalCommander) Client() jsonrpc.RPCClient {
	return e.client
}
