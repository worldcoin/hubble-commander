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
	storage     *st.Storage
	client      *eth.Client
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
	migrator, err := postgres.GetMigrator(c.cfg.Postgres)
	if err != nil {
		return err
	}

	c.storage, err = st.NewStorage(c.cfg.Postgres, c.cfg.Badger)
	if err != nil {
		return err
	}

	if c.cfg.Rollup.Prune {
		err = c.storage.Prune(migrator)
		if err != nil {
			return err
		}
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	chain, err := getChainConnection(c.cfg.Ethereum)
	if err != nil {
		return err
	}

	c.client, err = getClientOrBootstrapChainState(chain, c.storage, c.cfg.Rollup)
	if err != nil {
		return err
	}

	c.apiServer, err = api.NewAPIServer(c.cfg.API, c.storage, c.client)
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
	c.startWorker(func() error { return c.newBlockLoop() })
	c.startWorker(func() error { return c.rollupLoop() })
	c.startWorker(func() error { return WatchAccounts(c.storage, c.client, stopChannel) })
	c.stopChannel = stopChannel
	return nil
}

func (c *Commander) startWorker(fn func() error) {
	c.workers.Add(1)
	go func() {
		if err := fn(); err != nil {
			log.Fatalf("%+v", err)
		}
		c.workers.Done()
	}()
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
	defer c.clearState()
	close(c.stopChannel)
	if err := c.apiServer.Close(); err != nil {
		return err
	}
	c.workers.Wait()
	return c.storage.Close()
}

func (c *Commander) clearState() {
	c.stopChannel = nil
	c.storage = nil
	c.apiServer = nil
}

func getChainConnection(cfg *config.EthereumConfig) (deployer.ChainConnection, error) {
	if cfg == nil {
		return simulator.NewAutominingSimulator()
	}
	return deployer.NewRPCChainConnection(cfg)
}

func getClientOrBootstrapChainState(chain deployer.ChainConnection, storage *st.Storage, cfg *config.RollupConfig) (*eth.Client, error) {
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
