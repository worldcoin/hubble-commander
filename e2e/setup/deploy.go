package setup

import (
	"os"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const TestEthClientPrivateKey = "ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82"

func deployContracts(blockchain chain.Connection, deployerConfig *config.DeployerConfig) (*models.ChainSpec, *string, error) {
	err := deployChooser(blockchain, deployerConfig)
	if err != nil {
		return nil, nil, err
	}

	return deployRemainingContracts(blockchain, deployerConfig)
}

func deployChooser(blockchain chain.Connection, deployerConfig *config.DeployerConfig) error {
	e2eAccountAddress, err := privateKeyToAddress(TestEthClientPrivateKey)
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

func deployRemainingContracts(blockchain chain.Connection, deployerConfig *config.DeployerConfig) (*models.ChainSpec, *string, error) {
	file, err := os.CreateTemp("", "e2e_chain_spec")
	if err != nil {
		return nil, nil, err
	}

	chainSpecStr, err := commander.Deploy(deployerConfig, blockchain)
	if err != nil {
		return nil, nil, err
	}

	chainSpecPath := file.Name()
	err = utils.StoreChainSpec(chainSpecPath, *chainSpecStr)
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
