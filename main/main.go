package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
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
	{
		PublicKey: models.PublicKey{2, 3, 4},
		Balance:   models.MakeUint256(1000),
	},
}

func main() {
	cfg := config.GetConfig()

	storage, err := st.NewStorage(&cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	dep, err := getDeployer(cfg.Ethereum)
	if err != nil {
		log.Fatal(err)
	}

	client, err := getClient(storage, dep)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := commander.BlockNumberEndlessLoop(client, &cfg.Rollup)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		err := commander.CommitmentsEndlessLoop(storage, &cfg.Rollup)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		err := commander.BatchesEndlessLoop(storage, client, &cfg.Rollup)
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

func getClient(storage *st.Storage, dep deployer.ChainConnection) (*eth.Client, error) {
	chainState, err := storage.GetChainState(dep.GetChainID())
	if err != nil {
		return nil, err
	}

	if chainState == nil {
		fmt.Println("Bootstrapping genesis state")
		stateTree := st.NewStateTree(storage)
		chainState, err = bootstrapState(stateTree, dep, genesisAccounts)
		if err != nil {
			return nil, err
		}

		err = storage.SetChainState(chainState)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Continuing from saved state")
	}

	return createClientFromChainState(dep, chainState)
}

func createClientFromChainState(dep deployer.ChainConnection, chainState *models.ChainState) (*eth.Client, error) {
	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, dep.GetBackend())
	if err != nil {
		return nil, err
	}

	rollupContract, err := rollup.NewRollup(chainState.Rollup, dep.GetBackend())
	if err != nil {
		return nil, err
	}

	client, err := eth.NewClient(dep, eth.NewClientParams{
		Rollup:          rollupContract,
		AccountRegistry: accountRegistry,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getDeployer(cfg *config.EthereumConfig) (deployer.ChainConnection, error) {
	if cfg == nil {
		return simulator.NewAutominingSimulator()
	}

	chainID, ok := big.NewInt(0).SetString(cfg.ChainID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid chain id")
	}

	key, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	account, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	return deployer.NewRPCDeployer(cfg.RPCURL, chainID, account)
}

func bootstrapState(
	stateTree *st.StateTree,
	d deployer.ChainConnection,
	accounts []commander.GenesisAccount,
) (*models.ChainState, error) {
	accountRegistryAddress, accountRegistry, err := deployer.DeployAccountRegistry(d)
	if err != nil {
		return nil, err
	}

	registeredAccounts, err := commander.RegisterGenesisAccounts(d.GetAccount(), accountRegistry, accounts)
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

	chainState := &models.ChainState{
		ChainID:         d.GetChainID(),
		AccountRegistry: *accountRegistryAddress,
		Rollup:          contracts.RollupAddress,
	}

	return chainState, nil
}
