package setup

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

type InProcessCommander struct {
	client    jsonrpc.RPCClient
	commander *commander.Commander
	cfg       *config.Config
	chain     chain.Connection
}

func CreateInProcessCommander() (*InProcessCommander, error) {
	cfg := config.GetConfig()
	cfg.Bootstrap.Prune = true
	return CreateInProcessCommanderWithConfig(cfg, true)
}

func CreateInProcessCommanderWithConfig(cfg *config.Config, deployContracts bool) (*InProcessCommander, error) {
	cfg.Rollup.MinTxsPerCommitment = cfg.Rollup.MaxTxsPerCommitment
	blockchain, err := commander.GetChainConnection(cfg.Ethereum)
	if err != nil {
		return nil, err
	}

	cmd := commander.NewCommander(cfg, blockchain)
	endpoint := fmt.Sprintf("http://localhost:%s", cfg.API.Port)
	client := jsonrpc.NewClient(endpoint)

	if deployContracts {
		file, err := ioutil.TempFile("", "in_process_commander")
		if err != nil {
			return nil, err
		}

		chainSpecPath := file.Name()
		cfg.Bootstrap.ChainSpecPath = &chainSpecPath

		chainSpec, err := commander.Deploy(cfg, blockchain)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(*cfg.Bootstrap.ChainSpecPath, []byte(*chainSpec), 0600)
		if err != nil {
			return nil, err
		}
	}

	return &InProcessCommander{
		client:    client,
		commander: cmd,
		cfg:       cfg,
		chain:     blockchain,
	}, nil
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
	e.commander = commander.NewCommander(e.cfg, e.chain)
	return e.Start()
}

func (e *InProcessCommander) Client() jsonrpc.RPCClient {
	return e.client
}
