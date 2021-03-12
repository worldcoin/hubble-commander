package main

import (
	"log"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/deployer"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
)

var genesisAccounts = []commander.GenesisAccount{
	{
		AccountIndex: 0,
		Balance:      models.MakeUint256(1000),
	},
	{
		AccountIndex: 1,
		Balance:      models.MakeUint256(1000),
	},
	{
		AccountIndex: 2,
		Balance:      models.MakeUint256(1000),
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

	go commander.RollupLoop(storage, client, &cfg)

	log.Fatal(api.StartAPIServer(&cfg))
}

func NewSimulatedClient(stateTree *st.StateTree, accounts []commander.GenesisAccount) (*eth.Client, error) {
	err := commander.PopulateGenesisAccounts(stateTree, accounts)
	if err != nil {
		return nil, err
	}

	simulator, err := simulator.NewAutominingSimulator()
	if err != nil {
		return nil, err
	}

	stateRoot, err := stateTree.Root()
	if err != nil {
		return nil, err
	}

	contracts, err := deployer.DeployConfiguredRollup(simulator, deployer.DeploymentConfig{
		GenesisStateRoot: stateRoot,
	})
	if err != nil {
		return nil, err
	}

	client := eth.NewTestClient(simulator.Account, contracts.Rollup)
	return client, nil
}
