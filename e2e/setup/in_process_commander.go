package setup

import (
	"fmt"
	"os"
	"time"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
	"gopkg.in/yaml.v2"
)

const EthClientPrivateKey = "c216d5eef9c83c9d6f4629fff79e8e90d73b4beb9921de18f974f0d2c6d4e9b0"

type InProcessCommander struct {
	client     jsonrpc.RPCClient
	commander  *commander.Commander
	cfg        *config.Config
	blockchain chain.Connection
	chainSpec  *models.ChainSpec
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

	endpoint := fmt.Sprintf("http://localhost:%s", commanderConfig.API.Port)
	inProcessCommander := &InProcessCommander{
		client:     jsonrpc.NewClient(endpoint),
		commander:  commander.NewCommander(commanderConfig, blockchain),
		cfg:        commanderConfig,
		blockchain: blockchain,
	}

	if deployerConfig != nil {
		err = inProcessCommander.deployContracts(deployerConfig)
		if err != nil {
			return nil, err
		}
	}

	return inProcessCommander, nil
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
	return e.chainSpec
}

func deployChooser(blockchain chain.Connection, deployerConfig *config.DeployerConfig) error {
	e2eAccountAddress, err := privateKeyToAddress(EthClientPrivateKey)
	if err != nil {
		return err
	}
	poaAddress, _, err := deployer.DeployProofOfAuthority(
		blockchain,
		deployerConfig.Ethereum.MineTimeout,
		[]common.Address{blockchain.GetAccount().From, *e2eAccountAddress},
	)
	if err != nil {
		return err
	}
	deployerConfig.Bootstrap.Chooser = poaAddress
	return nil
}

func (e *InProcessCommander) deployContracts(deployerConfig *config.DeployerConfig) error {
	err := deployChooser(e.blockchain, deployerConfig)
	if err != nil {
		return err
	}

	e.chainSpec, e.cfg.Bootstrap.ChainSpecPath, err = deployRemainingContracts(e.blockchain, deployerConfig)
	return err
}

func deployRemainingContracts(blockchain chain.Connection, deployerConfig *config.DeployerConfig) (*models.ChainSpec, *string, error) {
	file, err := os.CreateTemp("", "in_process_commander")
	if err != nil {
		return nil, nil, err
	}

	chainSpecPath := file.Name()
	chainSpecStr, err := commander.Deploy(deployerConfig, blockchain)
	if err != nil {
		return nil, nil, err
	}

	err = os.WriteFile(chainSpecPath, []byte(*chainSpecStr), 0600)
	if err != nil {
		return nil, nil, err
	}

	var chainSpec models.ChainSpec
	err = yaml.Unmarshal([]byte(*chainSpecStr), &chainSpec)
	if err != nil {
		return nil, nil, err
	}

	return &chainSpec, &chainSpecPath, nil
}

func privateKeyToAddress(privateKey string) (*common.Address, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	return &address, nil
}
