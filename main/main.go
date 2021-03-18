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

	client, err := NewSimulatedClient(stateTree, genesisAccounts)
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

func NewSimulatedClient(stateTree *st.StateTree, accounts []commander.GenesisAccount) (*eth.Client, error) {
	sim, err := simulator.NewAutominingSimulator()
	if err != nil {
		return nil, err
	}
	accountRegistryAddress, accountRegistry, err := deployer.DeployAccountRegistry(sim)
	if err != nil {
		return nil, err
	}

	registeredAccounts, err := commander.RegisterGenesisAccounts(sim.Account, accountRegistry, accounts)
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

	contracts, err := deployer.DeployConfiguredRollup(sim, deployer.DeploymentConfig{
		AccountRegistryAddress: accountRegistryAddress,
		GenesisStateRoot:       stateRoot,
	})
	if err != nil {
		return nil, err
	}

	client := eth.NewTestClient(sim.Account, contracts.Rollup, contracts.AccountRegistry)
	return client, nil
}
