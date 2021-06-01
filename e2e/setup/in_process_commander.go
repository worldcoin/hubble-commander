package setup

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
	cfg.Rollup.Prune = true
	cfg.Rollup.SyncBatches = false

	cmd := commander.NewCommander(cfg)

	endpoint := fmt.Sprintf("http://localhost:%s", cfg.API.Port)
	client := jsonrpc.NewClient(endpoint)

	return &InProcessCommander{
		client:    client,
		commander: cmd,
	}
}

func CreateInProcessCommanderWithConfig(cfg *config.Config) *InProcessCommander {
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

	cfg := config.GetConfig()
	cfg.Rollup.Prune = false
	cfg.Rollup.SyncBatches = false

	e.commander = commander.NewCommander(cfg)
	return e.commander.Start()
}

func (e *InProcessCommander) Client() jsonrpc.RPCClient {
	return e.client
}
