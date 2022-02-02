package commander

import (
	"encoding/json"
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
	"github.com/ethereum/go-ethereum/common"
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

type Commander struct {
	lifecycle
	workers
	rollupControls

	cfg        *config.Config
	blockchain chain.Connection
	metrics    *metrics.CommanderMetrics

	storage       *st.Storage
	client        *eth.Client
	apiServer     *http.Server
	metricsServer *http.Server

	stateMutex     sync.Mutex
	invalidBatchID *models.Uint256

	txsHashesChan chan common.Hash
}

func NewCommander(cfg *config.Config, blockchain chain.Connection, migrate bool) *Commander {
	return &Commander{
		lifecycle:      lifecycle{},
		workers:        makeWorkers(),
		rollupControls: makeRollupControls(migrate),
		cfg:            cfg,
		blockchain:     blockchain,
		metrics:        metrics.NewCommanderMetrics(),
	}
}

func (c *Commander) StartAndWait() error {
	if err := c.Start(); err != nil {
		return err
	}

	<-c.getStartAndWaitChan()
	return nil
}

func (c *Commander) Start() (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isActive() {
		return nil
	}

	c.storage, err = st.NewStorage(c.cfg)
	if err != nil {
		return err
	}

	c.txsHashesChan = make(chan common.Hash, 32)

	c.client, err = getClient(c.blockchain, c.storage, c.cfg, c.metrics, c.txsHashesChan)
	if err != nil {
		return err
	}

	err = c.addGenesisBatch()
	if err != nil {
		return err
	}

	c.metricsServer = c.metrics.NewServer(c.cfg.Metrics)

	c.apiServer, err = api.NewServer(c.cfg, c.storage, c.client, c.metrics, c.EnableBatchCreation)
	if err != nil {
		return err
	}

	c.startWorker("API Server", func() error {
		err := c.apiServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	c.startWorker("Metrics Server", func() error {
		err := c.metricsServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	c.startWorker("Tracking Txs", func() error { return c.txsTracking(c.txsHashesChan) })
	c.startWorker("New Block Loop", func() error { return c.newBlockLoop() })

	go c.handleWorkerError()

	log.Printf("Commander started and listening on port %s", c.cfg.API.Port)
	c.setActive(true)
	return nil
}

func (c *Commander) Stop() (err error) {
	if !c.isActive() {
		return nil
	}

	c.closeOnce.Do(func() {
		if err = c.stop(); err != nil {
			return
		}
		log.Warningln("Commander stopped.")
		c.releaseStartAndWait()
	})
	return err
}

func (c *Commander) EnableBatchCreation(enable bool) {
	c.batchCreationEnabled = enable
	if !enable {
		c.stopRollupLoop()
	}
}

func (c *Commander) handleWorkerError() {
	<-c.workersContext.Done()
	c.closeOnce.Do(func() {
		log.Warning("Stopping commander gracefully...")

		if err := c.stop(); err != nil {
			log.Panicf("Failed to stop commander gracefully: %+v", err)
		}
		log.Panicln("Commander stopped by worker error")
	})
}

func (c *Commander) stop() error {
	if err := c.apiServer.Close(); err != nil {
		return err
	}
	if err := c.metricsServer.Close(); err != nil {
		return err
	}
	c.stopWorkersAndWait()
	c.setActive(false)
	return c.storage.Close()
}

func getClient(
	blockchain chain.Connection,
	storage *st.Storage,
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
	txsHashesChan chan<- common.Hash,
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
		return createClientFromChainState(blockchain, dbChainState, cfg, commanderMetrics, txsHashesChan)
	}

	if cfg.Bootstrap.ChainSpecPath != nil {
		return bootstrapFromChainState(blockchain, storage, cfg, commanderMetrics, txsHashesChan)
	}
	if cfg.Bootstrap.BootstrapNodeURL != nil {
		log.Printf("Bootstrapping genesis state from node %s", *cfg.Bootstrap.BootstrapNodeURL)
		return bootstrapFromRemoteState(blockchain, storage, cfg, commanderMetrics, txsHashesChan)
	}

	return nil, errors.WithStack(errMissingBootstrapSource)
}

func bootstrapFromChainState(
	blockchain chain.Connection,
	storage *st.Storage,
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
	txsHashesChan chan<- common.Hash,
) (*eth.Client, error) {
	chainSpec, err := ReadChainSpecFile(*cfg.Bootstrap.ChainSpecPath)
	if err != nil {
		return nil, err
	}
	importedChainState := newChainStateFromChainSpec(chainSpec)
	return bootstrapChainStateAndCommander(blockchain, storage, importedChainState, cfg, commanderMetrics, txsHashesChan)
}

func bootstrapChainStateAndCommander(
	blockchain chain.Connection,
	storage *st.Storage,
	importedChainState *models.ChainState,
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
	txsHashesChan chan<- common.Hash,
) (*eth.Client, error) {
	chainID := blockchain.GetChainID()
	if chainID != importedChainState.ChainID {
		return nil, errors.WithStack(errInconsistentFileChainID)
	}

	log.Printf("Bootstrapping genesis state from chain spec file")
	return setGenesisStateAndCreateClient(blockchain, storage, importedChainState, cfg, commanderMetrics, txsHashesChan)
}

func bootstrapFromRemoteState(
	blockchain chain.Connection,
	storage *st.Storage,
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
	txsHashesChan chan<- common.Hash,
) (*eth.Client, error) {
	fetchedChainState, err := fetchChainStateFromRemoteNode(*cfg.Bootstrap.BootstrapNodeURL)
	if err != nil {
		return nil, err
	}

	if !fetchedChainState.ChainID.EqN(cfg.Ethereum.ChainID) {
		return nil, errors.WithStack(errInconsistentRemoteChainID)
	}
	return setGenesisStateAndCreateClient(blockchain, storage, fetchedChainState, cfg, commanderMetrics, txsHashesChan)
}

func setGenesisStateAndCreateClient(
	blockchain chain.Connection,
	storage *st.Storage,
	chainState *models.ChainState,
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
	txsHashesChan chan<- common.Hash,
) (*eth.Client, error) {
	err := PopulateGenesisAccounts(storage, chainState.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	err = storage.SetChainState(chainState)
	if err != nil {
		return nil, err
	}

	client, err := createClientFromChainState(blockchain, chainState, cfg, commanderMetrics, txsHashesChan)
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
	cfg *config.Config,
	commanderMetrics *metrics.CommanderMetrics,
	txsHashesChan chan<- common.Hash,
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
		TxsHashesChan:   txsHashesChan,
		ClientConfig: eth.ClientConfig{
			TransferBatchSubmissionGasLimit:  ref.Uint64(cfg.Rollup.TransferBatchSubmissionGasLimit),
			C2TBatchSubmissionGasLimit:       ref.Uint64(cfg.Rollup.C2TBatchSubmissionGasLimit),
			MMBatchSubmissionGasLimit:        ref.Uint64(cfg.Rollup.MMBatchSubmissionGasLimit),
			DepositBatchSubmissionGasLimit:   ref.Uint64(cfg.Rollup.DepositBatchSubmissionGasLimit),
			TransitionDisputeGasLimit:        ref.Uint64(cfg.Rollup.TransitionDisputeGasLimit),
			SignatureDisputeGasLimit:         ref.Uint64(cfg.Rollup.SignatureDisputeGasLimit),
			BatchAccountRegistrationGasLimit: ref.Uint64(cfg.Rollup.BatchAccountRegistrationGasLimit),
			TxMineTimeout:                    ref.Duration(cfg.Ethereum.MineTimeout),
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
