package setup

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

type InProcessCommander struct {
	client     jsonrpc.RPCClient
	commander  *commander.Commander
	cfg        *config.Config
	blockchain chain.Connection
	chainSpec  *models.ChainSpec
}

func DeployAndCreateInProcessCommander(commanderConfig *config.Config, deployerConfig *config.DeployerConfig) (*InProcessCommander, error) {
	if commanderConfig == nil {
		commanderConfig = config.GetCommanderConfigAndSetupLogger()
	}

	commanderConfig.Badger.Path += "_e2e"
	commanderConfig.Bootstrap.Prune = true

	if deployerConfig == nil {
		deployerConfig = config.GetDeployerTestConfig()
	}

	return CreateInProcessCommander(commanderConfig, deployerConfig)
}

func CreateInProcessCommander(commanderConfig *config.Config, deployerConfig *config.DeployerConfig) (*InProcessCommander, error) {
	blockchain, err := commander.GetChainConnection(commanderConfig.Ethereum)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("http://localhost:%s", commanderConfig.API.Port)
	inProcessCmd := &InProcessCommander{
		client:     jsonrpc.NewClient(endpoint),
		commander:  commander.NewCommander(commanderConfig, blockchain),
		cfg:        commanderConfig,
		blockchain: blockchain,
	}

	if deployerConfig != nil {
		inProcessCmd.chainSpec, inProcessCmd.cfg.Bootstrap.ChainSpecPath, err = deployContracts(inProcessCmd.blockchain, deployerConfig)
		if err != nil {
			return nil, err
		}
	}

	return inProcessCmd, nil
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
	e.commander = commander.NewCommander(e.cfg, e.blockchain)
	return e.Start()
}

func (e *InProcessCommander) Client() jsonrpc.RPCClient {
	return e.client
}

func (e *InProcessCommander) ChainSpec() *models.ChainSpec {
	if e.chainSpec == nil {
		panic("call ChainSpec() on commander that deployed contracts")
	}
	return e.chainSpec
}
