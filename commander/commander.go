package commander

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc/v2"
)

var (
	errMissingBootstrapSource    = NewCannotBootstrapError("no chain spec file or bootstrap url specified")
	errInconsistentDBChainID     = NewInconsistentChainIDError("database")
	errInconsistentFileChainID   = NewInconsistentChainIDError("chain spec file")
	errInconsistentRemoteChainID = NewInconsistentChainIDError("fetched chain state")
)

// nolint:structcheck
type lifecycle struct {
	isRunning           bool
	releaseStartAndWait context.CancelFunc
	manualStop          bool

	workersContext     context.Context
	stopWorkersContext context.CancelFunc
	workersWaitGroup   sync.WaitGroup
}

type Commander struct {
	lifecycle

	cfg        *config.Config
	blockchain chain.Connection

	metrics       *metrics.CommanderMetrics
	storage       *st.Storage
	client        *eth.Client
	apiServer     *http.Server
	metricsServer *http.Server

	stateMutex        sync.Mutex
	rollupLoopRunning bool
	invalidBatchID    *models.Uint256
}

func NewCommander(cfg *config.Config, blockchain chain.Connection) *Commander {
	return &Commander{
		cfg:        cfg,
		blockchain: blockchain,
		lifecycle: lifecycle{
			releaseStartAndWait: func() {}, // noop
		},
	}
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

func (c *Commander) Start() (err error) {
	if c.isRunning {
		return nil
	}

	c.storage, err = st.NewStorage(c.cfg)
	if err != nil {
		return err
	}

	c.metrics = metrics.NewCommanderMetrics()

	c.client, err = getClient(c.blockchain, c.storage, c.cfg, c.metrics)
	if err != nil {
		return err
	}

	err = c.addGenesisBatch()
	if err != nil {
		return err
	}

	c.metricsServer = c.metrics.NewServer(c.cfg.Metrics)

	c.apiServer, err = api.NewServer(c.cfg, c.storage, c.client, c.metrics)
	if err != nil {
		return err
	}

	c.workersContext, c.stopWorkersContext = context.WithCancel(context.Background())

	c.startWorker("API Server", func() error {
		err = c.apiServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	c.startWorker("Metrics Server", func() error {
		err = c.metricsServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	c.startWorker("New Block Loop", func() error { return c.newBlockLoop() })

	go c.handleWorkerError()

	log.Printf("Commander started and listening on port %s", c.cfg.API.Port)
	c.isRunning = true
	return nil
}

func (c *Commander) Stop() error {
	if !c.isRunning {
		return nil
	}

	c.manualStop = true

	if err := c.stop(); err != nil {
		return err
	}

	log.Warningln("Commander stopped.")

	c.releaseStartAndWait()
	return nil
}

func (c *Commander) startWorker(name string, fn func() error) {
	c.workersWaitGroup.Add(1)
	go func() {
		var err error
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				var ok bool
				err, ok = recoverErr.(error)
				if !ok {
					err = fmt.Errorf("%+v", recoverErr)
				}
			}
			if err != nil {
				log.Errorf("%s worker failed with: %+v", name, err)
				c.stopWorkersContext()
			}
			c.workersWaitGroup.Done()
		}()

		err = fn()
	}()
}

func (c *Commander) handleWorkerError() {
	<-c.workersContext.Done()
	if c.manualStop {
		return
	}
	log.Warning("Stopping commander gracefully...")

	if err := c.stop(); err != nil {
		log.Panicf("Failed to stop commander gracefully: %+v", err)
	}
	log.Panicln("Commander stopped by worker error")
}

func (c *Commander) stop() error {
	if err := c.apiServer.Close(); err != nil {
		return err
	}
	if err := c.metricsServer.Close(); err != nil {
		return err
	}
	c.stopWorkersContext()
	c.workersWaitGroup.Wait()
	return c.storage.Close()
}

func getClient(
	blockchain chain.Connection,
	storage *st.Storage,
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
) (*eth.Client, error) {
	if cfg.Ethereum == nil {
		return nil, errors.Errorf("Ethereum config missing")
	}

	dbChainState, err := storage.GetChainState()
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	}

	if dbChainState != nil {
		if !dbChainState.ChainID.EqN(cfg.Ethereum.ChainID) {
			return nil, errors.WithStack(errInconsistentDBChainID)
		}

		log.Printf("Continuing from saved state on ChainID = %s", dbChainState.ChainID.String())
		return createClientFromChainState(blockchain, dbChainState, cfg.Rollup, commanderMetrics)
	}

	if cfg.Bootstrap.ChainSpecPath != nil {
		return bootstrapFromChainState(blockchain, storage, cfg, commanderMetrics)
	}
	if cfg.Bootstrap.BootstrapNodeURL != nil {
		log.Printf("Bootstrapping genesis state from node %s", *cfg.Bootstrap.BootstrapNodeURL)
		return bootstrapFromRemoteState(blockchain, storage, cfg, commanderMetrics)
	}

	return nil, errors.WithStack(errMissingBootstrapSource)
}

func bootstrapFromChainState(
	blockchain chain.Connection,
	storage *st.Storage,
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
) (*eth.Client, error) {
	chainSpec, err := ReadChainSpecFile(*cfg.Bootstrap.ChainSpecPath)
	if err != nil {
		return nil, err
	}
	importedChainState := newChainStateFromChainSpec(chainSpec)
	return bootstrapChainStateAndCommander(blockchain, storage, importedChainState, cfg.Rollup, commanderMetrics)
}

func bootstrapChainStateAndCommander(
	blockchain chain.Connection,
	storage *st.Storage,
	importedChainState *models.ChainState,
	cfg *config.RollupConfig,
	commanderMetrics *metrics.CommanderMetrics,
) (*eth.Client, error) {
	chainID := blockchain.GetChainID()
	if chainID != importedChainState.ChainID {
		return nil, errors.WithStack(errInconsistentFileChainID)
	}

	log.Printf("Bootstrapping genesis state from chain spec file")
	return setGenesisStateAndCreateClient(blockchain, storage, importedChainState, cfg, commanderMetrics)
}

func bootstrapFromRemoteState(
	blockchain chain.Connection,
	storage *st.Storage,
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
) (*eth.Client, error) {
	fetchedChainState, err := fetchChainStateFromRemoteNode(*cfg.Bootstrap.BootstrapNodeURL)
	if err != nil {
		return nil, err
	}

	if !fetchedChainState.ChainID.EqN(cfg.Ethereum.ChainID) {
		return nil, errors.WithStack(errInconsistentRemoteChainID)
	}
	return setGenesisStateAndCreateClient(blockchain, storage, fetchedChainState, cfg.Rollup, commanderMetrics)
}

func setGenesisStateAndCreateClient(
	blockchain chain.Connection,
	storage *st.Storage,
	chainState *models.ChainState,
	cfg *config.RollupConfig,
	commanderMetrics *metrics.CommanderMetrics,
) (*eth.Client, error) {
	err := PopulateGenesisAccounts(storage, chainState.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	err = storage.SetChainState(chainState)
	if err != nil {
		return nil, err
	}

	client, err := createClientFromChainState(blockchain, chainState, cfg, commanderMetrics)
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
		SpokeRegistry:                  info.SpokeRegistry,
		DepositManager:                 info.DepositManager,
		WithdrawManager:                info.WithdrawManager,
		Rollup:                         info.Rollup,
		GenesisAccounts:                genesisAccounts,
		SyncedBlock:                    getInitialSyncedBlock(info.AccountRegistryDeploymentBlock),
	}, nil
}

func createClientFromChainState(
	blockchain chain.Connection,
	chainState *models.ChainState,
	cfg *config.RollupConfig,
	commanderMetrics *metrics.CommanderMetrics,
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

	spokeRegistry, err := spokeregistry.NewSpokeRegistry(chainState.SpokeRegistry, backend)
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

	client, err := eth.NewClient(blockchain, commanderMetrics, &eth.NewClientParams{
		ChainState:      *chainState,
		Rollup:          rollupContract,
		AccountRegistry: accountRegistry,
		TokenRegistry:   tokenRegistry,
		SpokeRegistry:   spokeRegistry,
		DepositManager:  depositManager,
		ClientConfig: eth.ClientConfig{
			TransferBatchSubmissionGasLimit:  ref.Uint64(cfg.TransferBatchSubmissionGasLimit),
			C2TBatchSubmissionGasLimit:       ref.Uint64(cfg.C2TBatchSubmissionGasLimit),
			MMBatchSubmissionGasLimit:        ref.Uint64(cfg.MMBatchSubmissionGasLimit),
			DepositBatchSubmissionGasLimit:   ref.Uint64(cfg.DepositBatchSubmissionGasLimit),
			TransitionDisputeGasLimit:        ref.Uint64(cfg.TransitionDisputeGasLimit),
			SignatureDisputeGasLimit:         ref.Uint64(cfg.SignatureDisputeGasLimit),
			BatchAccountRegistrationGasLimit: ref.Uint64(cfg.BatchAccountRegistrationGasLimit),
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
