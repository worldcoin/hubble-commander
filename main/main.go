package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

var genesisAccounts = []commander.GenesisAccount{
	{
		PublicKey: models.PublicKey{1, 2, 3},
		Balance:   models.MakeUint256(1000),
	},
	{
		PublicKey: models.PublicKey{2, 3, 4},
		Balance:   models.MakeUint256(1000),
	},
	{
		PublicKey: models.PublicKey{3, 4, 5},
		Balance:   models.MakeUint256(1000),
	},
}

func main() {
	cfg := config.GetConfig()

	storage, err := st.NewStorage(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	stateTree := st.NewStateTree(storage)

	dep, err := GetDeployer(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	client, err := DeployContracts(stateTree, dep, genesisAccounts)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := commander.RollupLoop(storage, client, &cfg)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		err := commander.WatchAccounts(storage, client)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Fatal(api.StartAPIServer(&cfg))
}

func GetDeployer(cfg *config.Config) (deployer.Deployer, error) {
	if cfg.EthereumRPCURL == nil {
		sim, err := simulator.NewAutominingSimulator()
		if err != nil {
			return nil, err
		}

		return sim, nil
	}

	if cfg.EthereumChainID == nil {
		return nil, fmt.Errorf("chain id should be specified in the config when connecting to remote ethereum RPC")
	}

	chainID, ok := big.NewInt(0).SetString(*cfg.EthereumChainID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid chain id")
	}

	if cfg.EthereumPrivateKey == nil {
		return nil, fmt.Errorf("private key should be specified in the config when connecting to remote ethereum RPC")
	}

	key, err := crypto.HexToECDSA(*cfg.EthereumPrivateKey)
	if err != nil {
		return nil, err
	}

	account, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	return deployer.NewRPCDeployer(*cfg.EthereumRPCURL, account)
}

func DeployContracts(stateTree *st.StateTree, d deployer.Deployer, accounts []commander.GenesisAccount) (*eth.Client, error) {
	accountRegistryAddress, accountRegistry, err := deployer.DeployAccountRegistry(d)
	if err != nil {
		return nil, err
	}

	registeredAccounts, err := commander.RegisterGenesisAccounts(d.TransactionOpts(), accountRegistry, accounts)
	if err != nil {
		return nil, err
	}

	err = commander.PopulateGenesisAccounts(stateTree, registeredAccounts)
	if err != nil {
		return nil, err
	}

	stateRoot, err := stateTree.Root()
	if err != nil {
		return nil, err
	}

	contracts, err := deployer.DeployConfiguredRollup(d, deployer.DeploymentConfig{
		AccountRegistryAddress: accountRegistryAddress,
		GenesisStateRoot:       stateRoot,
	})
	if err != nil {
		return nil, err
	}

	client, err := eth.NewClient(d.TransactionOpts(), eth.NewClientParams{
		Rollup:          contracts.Rollup,
		AccountRegistry: contracts.AccountRegistry,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
