package setup

import (
	"fmt"
	"os"
	"time"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

const EthClientPrivateKey = "c216d5eef9c83c9d6f4629fff79e8e90d73b4beb9921de18f974f0d2c6d4e9b0"

type InProcessCommander struct {
	client     jsonrpc.RPCClient
	commander  *commander.Commander
	cfg        *config.Config
	blockchain chain.Connection
}

func DeployAndCreateInProcessCommander(commanderConfig *config.Config, deployerConfig *config.DeployerConfig) (*InProcessCommander, error) {
	if commanderConfig == nil {
		commanderConfig = config.GetConfig()
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

	cmd := commander.NewCommander(commanderConfig, blockchain)
	endpoint := fmt.Sprintf("http://localhost:%s", commanderConfig.API.Port)
	client := jsonrpc.NewClient(endpoint)

	if deployerConfig != nil {
		err = deployChooser(blockchain, deployerConfig)
		if err != nil {
			return nil, err
		}

		commanderConfig.Bootstrap.ChainSpecPath, err = deployRemainingContracts(blockchain, deployerConfig)
		if err != nil {
			return nil, err
		}
	}

	return &InProcessCommander{
		client:     client,
		commander:  cmd,
		cfg:        commanderConfig,
		blockchain: blockchain,
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
	e.commander = commander.NewCommander(e.cfg, e.blockchain)
	return e.Start()
}

func (e *InProcessCommander) Client() jsonrpc.RPCClient {
	return e.client
}

func deployChooser(blockchain chain.Connection, deployerConfig *config.DeployerConfig) error {
	e2eAccountAddress, err := privateKeyToAddress(EthClientPrivateKey)
	if err != nil {
		return err
	}
	poaAddress, _, err := deployer.DeployProofOfAuthority(blockchain, deployerConfig.Ethereum.MineTimeout, []common.Address{blockchain.GetAccount().From, *e2eAccountAddress})
	if err != nil {
		return err
	}
	deployerConfig.Bootstrap.Chooser = poaAddress
	return nil
}

func deployRemainingContracts(blockchain chain.Connection, deployerConfig *config.DeployerConfig) (*string, error) {
	file, err := os.CreateTemp("", "in_process_commander")
	if err != nil {
		return nil, err
	}

	chainSpecPath := file.Name()
	chainSpec, err := commander.Deploy(deployerConfig, blockchain)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(chainSpecPath, []byte(*chainSpec), 0600)
	if err != nil {
		return nil, err
	}

	return &chainSpecPath, nil
}

func privateKeyToAddress(privateKey string) (*common.Address, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	return &address, nil
}
