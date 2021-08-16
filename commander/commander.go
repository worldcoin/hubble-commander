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
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	ethRollup "github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc/v2"
)

var (
	ErrInvalidChainStates        = errors.New("database chain state and file chain state are not the same")
	ErrCannotBootstrap           = errors.New("cannot bootstrap - no chain spec file or bootstrap url specified")
	ErrChainSpecChainIDConflict  = errors.New("cannot bootstrap from chain spec - there's a conflict between config Chain ID and chain spec Chain ID")
	ErrDatabaseChainIDConflict   = errors.New("cannot bootstrap from chain spec - there's a conflict between config Chain ID and Chain ID saved in the database")
	ErrRemoteNodeChainIDConflict = errors.New("cannot bootstrap from remote node - there's a conflict between config Chain ID and Chain ID saved in the database")
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
	apiServer *http.Server
	domain    *bls.Domain
}

func NewCommander(cfg *config.Config) *Commander {
	return &Commander{
		cfg:                 cfg,
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

	chain, err := getChainConnection(c.cfg.Ethereum)
	if err != nil {
		return err
	}

	c.client, err = getClient(chain, c.storage, c.cfg)
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

	log.Printf("Commander started and listening on port %s.", c.cfg.API.Port)

	return nil
}

func (c *Commander) Deploy() (chainSpec *string, err error) {
	if c.IsRunning() {
		return nil, nil
	}

	c.storage, err = st.NewStorage(c.cfg)
	if err != nil {
		return nil, err
	}

	chain, err := getChainConnection(c.cfg.Ethereum)
	if err != nil {
		return nil, err
	}

	log.Printf(
		"Bootstrapping genesis state with %d accounts on chainId = %s",
		len(c.cfg.Bootstrap.GenesisAccounts),
		c.cfg.Ethereum.ChainID,
	)
	chainState, err := deployContractsAndSetupGenesisState(c.storage, chain, c.cfg.Bootstrap.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	chainSpec, err = GenerateChainSpec(chainState)
	if err != nil {
		return nil, err
	}

	return chainSpec, nil
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
	*c = *NewCommander(c.cfg)
}

func getChainConnection(cfg *config.EthereumConfig) (deployer.ChainConnection, error) {
	if cfg == nil {
		return simulator.NewAutominingSimulator()
	}
	return deployer.NewRPCChainConnection(cfg)
}

func getClient(chain deployer.ChainConnection, storage *st.Storage, cfg *config.Config) (*eth.Client, error) {
	if cfg.Bootstrap.ChainSpecPath != nil {
		return bootstrapFromChainState(chain, storage, cfg)

	}
	if cfg.Bootstrap.BootstrapNodeURL != nil {
		log.Printf("Bootstrapping genesis state from node %s", *cfg.Bootstrap.BootstrapNodeURL)
		return bootstrapFromRemoteState(chain, storage, cfg)
	}

	return nil, ErrCannotBootstrap
}

func bootstrapFromChainState(
	chain deployer.ChainConnection,
	storage *st.Storage,
	cfg *config.Config,
) (*eth.Client, error) {
	chainSpec, err := ReadChainSpecFile(*cfg.Bootstrap.ChainSpecPath)
	if err != nil {
		return nil, err
	}
	importedChainState := makeChainStateFromChainSpec(chainSpec)
	dbChainState, err := storage.GetChainState()
	if st.IsNotFoundError(err) {
		return bootstrapChainStateAndCommander(chain, storage, importedChainState)
	}
	if err != nil {
		return nil, err
	}

	if dbChainState.ChainID.String() != cfg.Ethereum.ChainID {
		return nil, ErrDatabaseChainIDConflict
	}

	err = compareChainStates(importedChainState, dbChainState)
	if err != nil {
		return nil, err
	}

	log.Printf("Continuing from saved state on ChainID = %s", importedChainState.ChainID.String())
	return createClientFromChainState(chain, importedChainState)
}

func bootstrapChainStateAndCommander(
	chain deployer.ChainConnection,
	storage *st.Storage,
	chainState *models.ChainState,
) (*eth.Client, error) {
	chainID := chain.GetChainID()
	if chainID != chainState.ChainID {
		return nil, ErrChainSpecChainIDConflict
	}
	err := storage.SetChainState(chainState)
	if err != nil {
		return nil, err
	}
	log.Printf("Bootstrapping genesis state from chain spec file")
	return createClientFromChainState(chain, chainState)
}

func compareChainStates(chainStateA, chainStateB *models.ChainState) error {
	if chainStateA.ChainID != chainStateB.ChainID ||
		chainStateA.DeploymentBlock != chainStateB.DeploymentBlock ||
		chainStateA.Rollup != chainStateB.Rollup ||
		chainStateA.AccountRegistry != chainStateB.AccountRegistry {
		return ErrInvalidChainStates
	}

	if len(chainStateA.GenesisAccounts) != len(chainStateB.GenesisAccounts) {
		return ErrInvalidChainStates
	}
	for i := range chainStateA.GenesisAccounts {
		if chainStateA.GenesisAccounts[i] != chainStateB.GenesisAccounts[i] {
			return ErrInvalidChainStates
		}
	}

	return nil
}

func bootstrapFromRemoteState(
	chain deployer.ChainConnection,
	storage *st.Storage,
	cfg *config.Config,
) (*eth.Client, error) {
	chainState, err := fetchChainStateFromRemoteNode(*cfg.Bootstrap.BootstrapNodeURL)
	if err != nil {
		return nil, err
	}

	if chainState.ChainID.String() != cfg.Ethereum.ChainID {
		return nil, ErrRemoteNodeChainIDConflict
	}

	err = PopulateGenesisAccounts(storage, chainState.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	err = storage.SetChainState(chainState)
	if err != nil {
		return nil, err
	}

	client, err := createClientFromChainState(chain, chainState)
	if err != nil {
		return nil, err
	}

	err = verifyCommitmentRoot(storage, client)
	if err != nil {
		return nil, err
	}
	return client, nil
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
		ChainID:         info.ChainID,
		AccountRegistry: info.AccountRegistry,
		DeploymentBlock: info.DeploymentBlock,
		Rollup:          info.Rollup,
		GenesisAccounts: genesisAccounts,
		SyncedBlock:     getInitialSyncedBlock(info.DeploymentBlock),
	}, nil
}

func createClientFromChainState(chain deployer.ChainConnection, chainState *models.ChainState) (*eth.Client, error) {
	err := logChainState(chainState)
	if err != nil {
		return nil, err
	}
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

func deployContractsAndSetupGenesisState(
	storage *st.Storage,
	chain deployer.ChainConnection,
	accounts []models.GenesisAccount,
) (*models.ChainState, error) {
	accountRegistryAddress, accountRegistryDeploymentBlock, accountRegistry, err := deployer.DeployAccountRegistry(chain)
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

	stateRoot, err := storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	contracts, err := ethRollup.DeployConfiguredRollup(chain, ethRollup.DeploymentConfig{
		Params:       ethRollup.Params{GenesisStateRoot: stateRoot},
		Dependencies: ethRollup.Dependencies{AccountRegistry: accountRegistryAddress},
	})
	if err != nil {
		return nil, err
	}

	chainState := &models.ChainState{
		ChainID:         chain.GetChainID(),
		AccountRegistry: *accountRegistryAddress,
		DeploymentBlock: *accountRegistryDeploymentBlock,
		Rollup:          contracts.RollupAddress,
		GenesisAccounts: populatedAccounts,
		SyncedBlock:     getInitialSyncedBlock(*accountRegistryDeploymentBlock),
	}

	return chainState, nil
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
