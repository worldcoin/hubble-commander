package commander

import (
	"log"
	"net/http"
	"sync"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	ethRollup "github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ybbus/jsonrpc/v2"
)

type Commander struct {
	cfg               *config.Config
	workers           sync.WaitGroup
	rollupLoopRunning bool
	stateMutex        sync.Mutex

	stopChannel      chan bool
	storage          *st.Storage
	client           *eth.Client
	apiServer        *http.Server
	signaturesDomain *bls.Domain
}

func NewCommander(cfg *config.Config) *Commander {
	return &Commander{
		cfg:        cfg,
		workers:    sync.WaitGroup{},
		stateMutex: sync.Mutex{},
	}
}

func (c *Commander) IsRunning() bool {
	return c.stopChannel != nil
}

func (c *Commander) Start() (err error) {
	if c.IsRunning() {
		return nil
	}

	c.storage, err = st.NewConfiguredStorage(c.cfg)
	if err != nil {
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

	c.signaturesDomain, err = c.storage.GetDomain(c.client.ChainState.ChainID)
	if err != nil {
		return err
	}

	c.apiServer, err = api.NewAPIServer(c.cfg.API, c.storage, c.client)
	if err != nil {
		return err
	}

	stopChannel := make(chan bool)
	c.startWorker(func() error {
		err = c.apiServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	c.startWorker(func() error { return c.newBlockLoop() })
	c.stopChannel = stopChannel

	log.Printf("Commander started and listening on port %s.\n", c.cfg.API.Port)

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
	err := c.storage.Close()
	if err != nil {
		return err
	}

	log.Println("Commander stopped.")

	return nil
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
	chainID := chain.GetChainID()
	chainState, err := storage.GetChainState(chainID)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	}

	if st.IsNotFoundError(err) {
		if cfg.BootstrapNodeURL != nil {
			log.Printf("Bootstrapping genesis state from node %s", *cfg.BootstrapNodeURL)
			chainState, err = fetchChainStateFromRemoteNode(*cfg.BootstrapNodeURL)
			if err != nil {
				return nil, err
			}

			err = PopulateGenesisAccounts(storage, chainState.GenesisAccounts)
			if err != nil {
				return nil, err
			}
		} else {
			log.Printf("Bootstrapping genesis state with %d accounts on chainId=%s.\n", len(cfg.GenesisAccounts), chainID.String())
			chainState, err = bootstrapState(storage, chain, cfg.GenesisAccounts)
			if err != nil {
				return nil, err
			}
		}

		err = storage.SetChainState(chainState)
		if err != nil {
			return nil, err
		}
	} else {
		log.Printf("Continuing from saved state on chainId=%s.\n", chainID.String())
	}

	return createClientFromChainState(chain, chainState)

	// TODO: Verify commitment root of batch #0 (for multi-operator sync).
}

func fetchChainStateFromRemoteNode(url string) (*models.ChainState, error) {
	client := jsonrpc.NewClient(url)
	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	if err != nil {
		return nil, err
	}

	err = client.CallFor(&info.ChainState.GenesisAccounts, "hubble_getGenesisAccounts")
	if err != nil {
		return nil, err
	}

	return &info.ChainState, nil
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

	populatedAccounts := AssignStateIDs(registeredAccounts)
	err = PopulateGenesisAccounts(storage, populatedAccounts)
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
		GenesisAccounts: populatedAccounts,
		SyncedBlock:     0,
	}

	return chainState, nil
}
