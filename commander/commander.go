package commander

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/chain"
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
	errInconsistentChainState    = NewCannotBootstrapError("database chain state and file chain state are not the same")
	errMissingBootstrapSource    = NewCannotBootstrapError("no chain spec file or bootstrap url specified")
	errInconsistentDBChainID     = NewInconsistentChainIDError("database")
	errInconsistentFileChainID   = NewInconsistentChainIDError("chain spec file")
	errInconsistentRemoteChainID = NewInconsistentChainIDError("fetched chain state")
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

	storage    *st.Storage
	client     *eth.Client
	blockchain chain.Connection
	apiServer  *http.Server
}

func NewCommander(cfg *config.Config, blockchain chain.Connection) *Commander {
	return &Commander{
		cfg:                 cfg,
		blockchain:          blockchain,
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

	c.client, err = getClient(c.blockchain, c.storage, c.cfg)
	if err != nil {
		return err
	}

	err = c.addGenesisBatch()
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
	*c = *NewCommander(c.cfg, c.blockchain)
}

func GetChainConnection(cfg *config.EthereumConfig) (chain.Connection, error) {
	if cfg.RPCURL == "simulator" {
		return simulator.NewConfiguredSimulator(simulator.Config{
			FirstAccountPrivateKey: ref.String(cfg.PrivateKey),
			AutomineEnabled:        ref.Bool(true),
		})
	}
	return chain.NewRPCCConnection(cfg)
}

func getClient(blockchain chain.Connection, storage *st.Storage, cfg *config.Config) (*eth.Client, error) {
	if cfg.Ethereum == nil {
		log.Fatal("no Ethereum config")
	}

	dbChainState, err := storage.GetChainState()
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	}

	if dbChainState != nil {
		if dbChainState.ChainID.CmpN(cfg.Ethereum.ChainID) != 0 {
			return nil, errors.WithStack(errInconsistentDBChainID)
		}
	}

	if cfg.Bootstrap.ChainSpecPath != nil {
		return bootstrapFromChainState(blockchain, dbChainState, storage, cfg)
	}
	if cfg.Bootstrap.BootstrapNodeURL != nil {
		log.Printf("Bootstrapping genesis state from node %s", *cfg.Bootstrap.BootstrapNodeURL)
		return bootstrapFromRemoteState(blockchain, storage, cfg)
	}

	return nil, errors.WithStack(errMissingBootstrapSource)
}

func bootstrapFromChainState(
	blockchain chain.Connection,
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
		return bootstrapChainStateAndCommander(blockchain, storage, importedChainState, cfg.Rollup)
	}

	err = compareChainStates(importedChainState, dbChainState)
	if err != nil {
		return nil, err
	}

	log.Printf("Continuing from saved state on ChainID = %s", importedChainState.ChainID.String())
	return createClientFromChainState(blockchain, importedChainState, cfg.Rollup)
}

func compareChainStates(chainStateA, chainStateB *models.ChainState) error {
	if chainStateA.ChainID != chainStateB.ChainID ||
		chainStateA.AccountRegistryDeploymentBlock != chainStateB.AccountRegistryDeploymentBlock ||
		chainStateA.Rollup != chainStateB.Rollup ||
		chainStateA.AccountRegistry != chainStateB.AccountRegistry ||
		chainStateA.TokenRegistry != chainStateB.TokenRegistry ||
		chainStateA.DepositManager != chainStateB.DepositManager {
		return errors.WithStack(errInconsistentChainState)
	}

	if len(chainStateA.GenesisAccounts) != len(chainStateB.GenesisAccounts) {
		return errors.WithStack(errInconsistentChainState)
	}
	for i := range chainStateA.GenesisAccounts {
		if chainStateA.GenesisAccounts[i] != chainStateB.GenesisAccounts[i] {
			return errors.WithStack(errInconsistentChainState)
		}
	}

	return nil
}

func bootstrapChainStateAndCommander(
	blockchain chain.Connection,
	storage *st.Storage,
	importedChainState *models.ChainState,
	cfg *config.RollupConfig,
) (*eth.Client, error) {
	chainID := blockchain.GetChainID()
	if chainID != importedChainState.ChainID {
		return nil, errors.WithStack(errInconsistentFileChainID)
	}

	log.Printf("Bootstrapping genesis state from chain spec file")
	return setGenesisStateAndCreateClient(blockchain, storage, importedChainState, cfg)
}

func bootstrapFromRemoteState(
	blockchain chain.Connection,
	storage *st.Storage,
	cfg *config.Config,
) (*eth.Client, error) {
	fetchedChainState, err := fetchChainStateFromRemoteNode(*cfg.Bootstrap.BootstrapNodeURL)
	if err != nil {
		return nil, err
	}

	if fetchedChainState.ChainID.CmpN(cfg.Ethereum.ChainID) != 0 {
		return nil, errors.WithStack(errInconsistentRemoteChainID)
	}

	return setGenesisStateAndCreateClient(blockchain, storage, fetchedChainState, cfg.Rollup)
}

func setGenesisStateAndCreateClient(
	blockchain chain.Connection,
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

	client, err := createClientFromChainState(blockchain, chainState, cfg)
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
	blockchain chain.Connection,
	chainState *models.ChainState,
	cfg *config.RollupConfig,
) (*eth.Client, error) {
	err := logChainState(chainState)
	if err != nil {
		return nil, err
	}

	backend := blockchain.GetBackend()

	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, backend)
	if err != nil {
		return nil, err
	}

	tokenRegistry, err := tokenregistry.NewTokenRegistry(chainState.TokenRegistry, backend)
	if err != nil {
		return nil, err
	}

	depositManager, err := depositmanager.NewDepositManager(chainState.DepositManager, backend)
	if err != nil {
		return nil, err
	}

	rollupContract, err := rollup.NewRollup(chainState.Rollup, backend)
	if err != nil {
		return nil, err
	}

	client, err := eth.NewClient(blockchain, &eth.NewClientParams{
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
