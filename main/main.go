package main

import (
	"log"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
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

	sim, err := simulator.NewAutominingSimulator()
	if err != nil {
		log.Fatal(err)
	}

	client, err := DeployContracts(stateTree, sim, genesisAccounts)
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
