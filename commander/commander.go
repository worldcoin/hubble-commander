package commander

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc/v2"
)

var (
	inconsistentChainStateErr    = NewCannotBootstrapError("database chain state and file chain state are not the same")
	missingBootstrapSourceErr    = NewCannotBootstrapError("no chain spec file or bootstrap url specified")
	inconsistentDBChainIDErr     = NewInconsistentChainIDError("database")
	inconsistentFileChainIDErr   = NewInconsistentChainIDError("chain spec file")
	inconsistentRemoteChainIDErr = NewInconsistentChainIDError("fetched chain state")
)

type Commander struct {
	cfg                 *config.Config
	workersContext      context.Context
	stopWorkers         context.CancelFunc
	workers             sync.WaitGroup
	releaseStartAndWait context.CancelFunc

	invalidBatchID    *models.Uint256
	rollupLoopRunning bool
	stateMutex        sync.Mutex

	storage   *st.Storage
	client    *eth.Client
	chain     deployer.ChainConnection
	apiServer *http.Server
	domain    *bls.Domain
}

func NewCommander(cfg *config.Config, chain deployer.ChainConnection) *Commander {
	return &Commander{
		cfg:                 cfg,
		chain:               chain,
		releaseStartAndWait: func() {}, // noop
	}
}

func (c *Commander) IsRunning() bool {
	return c.workersContext != nil
}

func (c *Commander) Start() (err error) {
	if c.IsRunning() {
		return nil
	}

	c.storage, err = st.NewStorage(c.cfg)
	if err != nil {
		return err
	}

	c.client, err = getClient(c.chain, c.storage, c.cfg)
	if err != nil {
		return err
	}

	err = c.addGenesisBatch()
	if err != nil {
		return err
	}

	c.domain, err = c.client.GetDomain()
	if err != nil {
		return err
	}

	c.apiServer, err = api.NewAPIServer(c.cfg, c.storage, c.client)
	if err != nil {
		return err
	}

	c.workersContext, c.stopWorkers = context.WithCancel(context.Background())

	c.startWorker(func() error {
		err = c.apiServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	c.startWorker(func() error { return c.newBlockLoop() })

	log.Printf("Commander started and listening on port %s", c.cfg.API.Port)

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
	var stopContext context.Context
	stopContext, c.releaseStartAndWait = context.WithCancel(context.Background())

	<-stopContext.Done()
	return nil
}

func (c *Commander) Stop() error {
	if !c.IsRunning() {
		return nil
	}

	if err := c.apiServer.Close(); err != nil {
		return err
	}
	c.stopWorkers()
	c.workers.Wait()
	if err := c.storage.Close(); err != nil {
		return err
	}

	log.Warningln("Commander stopped.")

	c.releaseStartAndWait()
	c.resetCommander()
	return nil
}

func (c *Commander) resetCommander() {
	*c = *NewCommander(c.cfg, c.chain)
}

func GetChainConnection(cfg *config.EthereumConfig) (deployer.ChainConnection, error) {
	if cfg.RPCURL == "simulator" {
		return simulator.NewConfiguredSimulator(simulator.Config{
			FirstAccountPrivateKey: ref.String(cfg.PrivateKey),
			AutomineEnabled:        ref.Bool(true),
		})
	}
	return deployer.NewRPCChainConnection(cfg)
}

func getClient(chain deployer.ChainConnection, storage *st.Storage, cfg *config.Config) (*eth.Client, error) {
	if cfg.Ethereum == nil {
		log.Fatal("no Ethereum config")
	}

	dbChainState, err := storage.GetChainState()
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	}

	if dbChainState != nil {
		if dbChainState.ChainID.String() != cfg.Ethereum.ChainID {
			return nil, inconsistentDBChainIDErr
		}
	}

	if cfg.Bootstrap.ChainSpecPath != nil {
		return bootstrapFromChainState(chain, dbChainState, storage, cfg)
	}
	if cfg.Bootstrap.BootstrapNodeURL != nil {
		log.Printf("Bootstrapping genesis state from node %s", *cfg.Bootstrap.BootstrapNodeURL)
		return bootstrapFromRemoteState(chain, storage, cfg)
	}

	return nil, missingBootstrapSourceErr
}

func bootstrapFromChainState(
	chain deployer.ChainConnection,
	dbChainState *models.ChainState,
	storage *st.Storage,
	cfg *config.Config,
) (*eth.Client, error) {
	chainSpec, err := ReadChainSpecFile(*cfg.Bootstrap.ChainSpecPath)
	if err != nil {
		return nil, err
	}
	importedChainState := newChainStateFromChainSpec(chainSpec)

	if dbChainState == nil {
		return bootstrapChainStateAndCommander(chain, storage, importedChainState, cfg.Rollup)
	}

	err = compareChainStates(importedChainState, dbChainState)
	if err != nil {
		return nil, err
	}

	log.Printf("Continuing from saved state on ChainID = %s", importedChainState.ChainID.String())
	return createClientFromChainState(chain, importedChainState, cfg.Rollup)
}

func compareChainStates(chainStateA, chainStateB *models.ChainState) error {
	compareError := inconsistentChainStateErr

	if chainStateA.ChainID != chainStateB.ChainID ||
		chainStateA.AccountRegistryDeploymentBlock != chainStateB.AccountRegistryDeploymentBlock ||
		chainStateA.Rollup != chainStateB.Rollup ||
		chainStateA.AccountRegistry != chainStateB.AccountRegistry ||
		chainStateA.TokenRegistry != chainStateB.TokenRegistry ||
		chainStateA.DepositManager != chainStateB.DepositManager {
		return compareError
	}

	if len(chainStateA.GenesisAccounts) != len(chainStateB.GenesisAccounts) {
		return compareError
	}
	for i := range chainStateA.GenesisAccounts {
		if chainStateA.GenesisAccounts[i] != chainStateB.GenesisAccounts[i] {
			return compareError
		}
	}

	return nil
}

func bootstrapChainStateAndCommander(
	chain deployer.ChainConnection,
	storage *st.Storage,
	importedChainState *models.ChainState,
	cfg *config.RollupConfig,
) (*eth.Client, error) {
	chainID := chain.GetChainID()
	if chainID != importedChainState.ChainID {
		return nil, inconsistentFileChainIDErr
	}

	log.Printf("Bootstrapping genesis state from chain spec file")
	return setGenesisStateAndCreateClient(chain, storage, importedChainState, cfg)
}

func bootstrapFromRemoteState(
	chain deployer.ChainConnection,
	storage *st.Storage,
	cfg *config.Config,
) (*eth.Client, error) {
	fetchedChainState, err := fetchChainStateFromRemoteNode(*cfg.Bootstrap.BootstrapNodeURL)
	if err != nil {
		return nil, err
	}

	if fetchedChainState.ChainID.String() != cfg.Ethereum.ChainID {
		return nil, inconsistentRemoteChainIDErr
	}

	return setGenesisStateAndCreateClient(chain, storage, fetchedChainState, cfg.Rollup)
}

func setGenesisStateAndCreateClient(
	chain deployer.ChainConnection,
	storage *st.Storage,
	chainState *models.ChainState,
	cfg *config.RollupConfig,
) (*eth.Client, error) {
	err := PopulateGenesisAccounts(storage, chainState.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	err = storage.SetChainState(chainState)
	if err != nil {
		return nil, err
	}

	client, err := createClientFromChainState(chain, chainState, cfg)
	if err != nil {
		return nil, err
	}

	err = verifyCommitmentRoot(storage, client)
	if err != nil {
		return nil, err
	}

	return client, err
}

func fetchChainStateFromRemoteNode(url string) (*models.ChainState, error) {
	client := jsonrpc.NewClient(url)

	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	if err != nil {
		return nil, err
	}

	var genesisAccounts models.GenesisAccounts
	err = client.CallFor(&genesisAccounts, "hubble_getGenesisAccounts")
	if err != nil {
		return nil, err
	}

	return &models.ChainState{
		ChainID:                        info.ChainID,
		AccountRegistry:                info.AccountRegistry,
		AccountRegistryDeploymentBlock: info.AccountRegistryDeploymentBlock,
		TokenRegistry:                  info.TokenRegistry,
		DepositManager:                 info.DepositManager,
		Rollup:                         info.Rollup,
		GenesisAccounts:                genesisAccounts,
		SyncedBlock:                    getInitialSyncedBlock(info.AccountRegistryDeploymentBlock),
	}, nil
}

func createClientFromChainState(
	chain deployer.ChainConnection,
	chainState *models.ChainState,
	cfg *config.RollupConfig,
) (*eth.Client, error) {
	err := logChainState(chainState)
	if err != nil {
		return nil, err
	}

	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, chain.GetBackend())
	if err != nil {
		return nil, err
	}

	tokenRegistry, err := tokenregistry.NewTokenRegistry(chainState.TokenRegistry, chain.GetBackend())
	if err != nil {
		return nil, err
	}

	depositManager, err := depositmanager.NewDepositManager(chainState.DepositManager, chain.GetBackend())
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
		TokenRegistry:   tokenRegistry,
		DepositManager:  depositManager,
		ClientConfig: eth.ClientConfig{
			TransitionDisputeGasLimit: ref.Uint64(cfg.TransitionDisputeGasLimit),
			SignatureDisputeGasLimit:  ref.Uint64(cfg.SignatureDisputeGasLimit),
		},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getInitialSyncedBlock(deploymentBlock uint64) uint64 {
	return deploymentBlock - 1
}

func logChainState(chainState *models.ChainState) error {
	jsonState, err := json.Marshal(*chainState)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Debugf("Creating ethereum client from chain state: %s", string(jsonState))
	return nil
}
