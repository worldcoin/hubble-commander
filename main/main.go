package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	ethRollup "github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/golang-migrate/migrate/v4"
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
	prune := flag.Bool("prune", false, "drop database before running app")
	flag.Parse()

	cfg := config.GetConfig()

	migrator, err := db.GetMigrator(&cfg.DB)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if *prune {
		err = migrator.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("%+v", err)
	}

	storage, err := st.NewStorage(&cfg.DB)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	chain, err := getDeployer(cfg.Ethereum)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	client, err := getClient(storage, chain)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	go func() {
		err := commander.BlockNumberEndlessLoop(client, &cfg.Rollup)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}()
	go func() {
		err := commander.RollupEndlessLoop(storage, client, &cfg.Rollup)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}()
	go func() {
		err := commander.WatchAccounts(storage, client)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}()

	log.Fatal(api.StartAPIServer(&cfg, client))
}

func getClient(storage *st.Storage, chain deployer.ChainConnection) (*eth.Client, error) {
	chainState, err := storage.GetChainState(chain.GetChainID())

	if st.IsNotFoundError(err) {
		fmt.Println("Bootstrapping genesis state")
		stateTree := st.NewStateTree(storage)
		chainState, err = bootstrapState(stateTree, chain, genesisAccounts)
		if err != nil {
			return nil, err
		}

		err = storage.SetChainState(chainState)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		fmt.Println("Continuing from saved state")
	}

	return createClientFromChainState(chain, chainState)
}

func createClientFromChainState(chain deployer.ChainConnection, chainState *models.ChainState) (*eth.Client, error) {
	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, chain.GetBackend())
	if err != nil {
		return nil, err
	}

	rollupContract, err := rollup.NewRollup(chainState.Rollup, chain.GetBackend())
	if err != nil {
		return nil, err
	}

	client, err := eth.NewClient(chain, &eth.NewClientParams{
		ChainState:      *chainState,
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
	return deployer.NewRPCDeployer(cfg)
}

func bootstrapState(
	stateTree *st.StateTree,
	chain deployer.ChainConnection,
	accounts []commander.GenesisAccount,
) (*models.ChainState, error) {
	accountRegistryAddress, accountRegistry, err := deployer.DeployAccountRegistry(chain)
	if err != nil {
		return nil, err
	}

	registeredAccounts, err := commander.RegisterGenesisAccounts(chain.GetAccount(), accountRegistry, accounts)
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

	contracts, err := ethRollup.DeployConfiguredRollup(chain, ethRollup.DeploymentConfig{
		AccountRegistryAddress: accountRegistryAddress,
		GenesisStateRoot:       stateRoot,
	})
	if err != nil {
		return nil, err
	}

	chainState := &models.ChainState{
		ChainID:         chain.GetChainID(),
		AccountRegistry: *accountRegistryAddress,
		Rollup:          contracts.RollupAddress,
	}

	return chainState, nil
}
