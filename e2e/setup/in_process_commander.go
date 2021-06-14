package setup

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

type InProcessCommander struct {
	client    jsonrpc.RPCClient
	commander *commander.Commander
	cfg       *config.Config
}

func CreateInProcessCommander() *InProcessCommander {
	cfg := config.GetConfig()
	cfg.Bootstrap.Prune = true
	return CreateInProcessCommanderWithConfig(cfg)
}

func CreateInProcessCommanderWithConfig(cfg *config.Config) *InProcessCommander {
	cmd := commander.NewCommander(cfg)

	endpoint := fmt.Sprintf("http://localhost:%s", cfg.API.Port)
	client := jsonrpc.NewClient(endpoint)

	return &InProcessCommander{
		client:    client,
		commander: cmd,
		cfg:       cfg,
	}
}

func (e *InProcessCommander) Start() error {
	err := e.commander.Start()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			var version string
			err = e.client.CallFor(&version, "hubble_getVersion")
			if err == nil {
				return nil
			}
		case <-timeout:
			return errors.Errorf("In-process commander start timed out: %s", err.Error())
		}
	}
}

func (e *InProcessCommander) Stop() error {
	return e.commander.Stop()
}

func (e *InProcessCommander) Restart() error {
	err := e.Stop()
	if err != nil {
		return err
	}
	e.cfg.Bootstrap.Prune = false
	e.commander = commander.NewCommander(e.cfg)
	return e.Start()
}

func (e *InProcessCommander) Client() jsonrpc.RPCClient {
	return e.client
}
