package commander

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	ethRollup "github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/golang-migrate/migrate/v4"
)

type Commander struct {
	cfg     *config.Config
	workers sync.WaitGroup

	stopChannel chan bool
	apiServer   *http.Server
}

func NewCommander(cfg *config.Config) *Commander {
	return &Commander{
		cfg:     cfg,
		workers: sync.WaitGroup{},
	}
}

func (c *Commander) IsRunning() bool {
	return c.stopChannel != nil
}

func (c *Commander) Start() error {
	if c.IsRunning() {
		return nil
	}
	migrator, err := postgres.GetMigrator(&c.cfg.DB)
	if err != nil {
		return err
	}

	if c.cfg.Rollup.Prune {
		err = migrator.Down()
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	storage, err := st.NewStorage(&c.cfg.DB)
	if err != nil {
		return err
	}

	chain, err := getDeployer(c.cfg.Ethereum)
	if err != nil {
		return err
	}

	client, err := getClient(storage, chain, &c.cfg.Rollup)
	if err != nil {
		return err
	}

	c.apiServer, err = api.NewAPIServer(&c.cfg.API, storage, client)
	if err != nil {
		return err
	}

	stopChannel := make(chan bool)
	c.startWorker(func() error {
		err := c.apiServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	c.startWorker(func() error { return BlockNumberLoop(storage, client, &c.cfg.Rollup, stopChannel) })
	c.startWorker(func() error { return RollupLoop(storage, client, &c.cfg.Rollup, stopChannel) })
	c.startWorker(func() error { return WatchAccounts(storage, client, stopChannel) })
	c.stopChannel = stopChannel
	return nil
}

func (c *Commander) startWorker(fn func() error) {
	startGoroutine(func() {
		c.workers.Add(1)
		defer c.workers.Done()
		if err := fn(); err != nil {
			log.Fatalf("%+v", err)
		}
	})
}

func startGoroutine(fn func()) {
	var start sync.WaitGroup
	start.Add(1)
	go func() {
		start.Done()
		fn()
	}()
	start.Wait()
}

func (c *Commander) StartAndWait() error {
	if err := c.Start(); err != nil {
		return err
	}
	c.workers.Wait()
	return nil
}

func (c *Commander) Stop() error {
	if !c.IsRunning() {
		return nil
	}
	close(c.stopChannel)
	err := c.apiServer.Close()
	c.workers.Wait()
	c.clearState()
	return err
}

func (c *Commander) clearState() {
	c.stopChannel = nil
	c.apiServer = nil
}

func getDeployer(cfg *config.EthereumConfig) (deployer.ChainConnection, error) {
	if cfg == nil {
		return simulator.NewAutominingSimulator()
	}
	return deployer.NewRPCChainConnection(cfg)
}

func getClient(storage *st.Storage, chain deployer.ChainConnection, cfg *config.RollupConfig) (*eth.Client, error) {
	chainState, err := storage.GetChainState(chain.GetChainID())

	if st.IsNotFoundError(err) {
		fmt.Println("Bootstrapping genesis state")
		chainState, err = bootstrapState(storage, chain, cfg.GenesisAccounts)
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

func bootstrapState(
	storage *st.Storage,
	chain deployer.ChainConnection,
	accounts []models.GenesisAccount,
) (*models.ChainState, error) {
	accountRegistryAddress, accountRegistry, err := deployer.DeployAccountRegistry(chain)
	if err != nil {
		return nil, err
	}

	registeredAccounts, err := RegisterGenesisAccounts(chain.GetAccount(), accountRegistry, accounts)
	if err != nil {
		return nil, err
	}

	err = PopulateGenesisAccounts(storage, registeredAccounts)
	if err != nil {
		return nil, err
	}

	stateRoot, err := st.NewStateTree(storage).Root()
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
