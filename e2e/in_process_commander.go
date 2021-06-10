package e2e

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/ybbus/jsonrpc/v2"
)

type InProcessCommander struct {
	client    jsonrpc.RPCClient
	commander *commander.Commander
}

func CreateInProcessCommander() *InProcessCommander {
	cfg := config.GetConfig()
	cfg.Bootstrap.Prune = true

	cmd := commander.NewCommander(cfg)

	endpoint := fmt.Sprintf("http://localhost:%s", cfg.API.Port)
	client := jsonrpc.NewClient(endpoint)

	return &InProcessCommander{
		client:    client,
		commander: cmd,
	}
}

func (e *InProcessCommander) Start() error {
	return e.commander.Start()
}

func (e *InProcessCommander) Stop() error {
	return e.commander.Stop()
}

func (e *InProcessCommander) Restart() error {
	err := e.commander.Stop()
	if err != nil {
		return err
	}
	return e.commander.Start()
}

func (e *InProcessCommander) Client() jsonrpc.RPCClient {
	return e.client
}
