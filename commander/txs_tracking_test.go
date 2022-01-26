package commander

import (
	"context"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxsTrackingTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd     *Commander
	storage *st.TestStorage
	client  *eth.TestClient
	cfg     *config.Config
	wallets []bls.Wallet
}

func (s *TxsTrackingTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DisableSignatures = true
}

func (s *TxsTrackingTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	lowGasLimit := uint64(25_000)

	s.client = newClientWithGenesisStateWithClientConfig(s.T(), s.storage, &eth.ClientConfig{
		TransferBatchSubmissionGasLimit: &lowGasLimit,
		C2TBatchSubmissionGasLimit:      &lowGasLimit,
		MMBatchSubmissionGasLimit:       &lowGasLimit,
	})

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)

	setStateLeaves(s.T(), s.storage.Storage)
	s.cmd = NewCommander(s.cfg, s.client.Blockchain)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.storage.Storage
	s.cmd.metrics = metrics.NewCommanderMetrics()
	s.cmd.workersContext, s.cmd.stopWorkersContext = context.WithCancel(context.Background())

	err = s.cmd.addGenesisBatch()
	s.NoError(err)

	s.setAccountsAndChainState()
}

func (s *TxsTrackingTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxsTrackingTestSuite) TestTxsTracking_FailedTransferTransaction() {
	s.startWorkers()
	s.waitForLatestBlockSync()

	transfer := testutils.MakeTransfer(0, 1, 0, 400)
	s.submitBatchInTransaction(&transfer, batchtype.Transfer)

	s.Eventually(func() bool {
		err := s.cmd.workersContext.Err()
		return err == context.Canceled
	}, time.Second, time.Millisecond*300)
}

func (s *TxsTrackingTestSuite) TestTxsTracking_FailedCreate2TransfersTransaction() {
	s.startWorkers()
	s.waitForLatestBlockSync()

	transfer := testutils.MakeCreate2Transfer(0, ref.Uint32(1), 0, 50, &models.PublicKey{2, 3, 4})
	s.submitBatchInTransaction(&transfer, batchtype.Create2Transfer)

	s.Eventually(func() bool {
		err := s.cmd.workersContext.Err()
		return err == context.Canceled
	}, time.Second, time.Millisecond*300)
}

func (s *TxsTrackingTestSuite) TestTxsTracking_FailedMassMigrationTransaction() {
	s.startWorkers()
	s.waitForLatestBlockSync()

	commitment := models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				ID: models.CommitmentID{
					BatchID:      models.MakeUint256(1),
					IndexInBatch: 0,
				},
				Type: batchtype.MassMigration,
			},
			Meta: &models.MassMigrationMeta{
				SpokeID:     1,
				TokenID:     models.MakeUint256(0),
				Amount:      models.MakeUint256(400),
				FeeReceiver: 0,
			},
			WithdrawRoot: utils.RandomHash(),
		},
		Transactions: []uint8{0, 0, 0, 0, 32, 4, 0, 0},
	}

	_, err := s.client.Client.SubmitMassMigrationsBatch(models.NewUint256(1),
		[]models.CommitmentWithTxs{&commitment})
	s.NoError(err)

	s.Eventually(func() bool {
		err := s.cmd.workersContext.Err()
		return err == context.Canceled
	}, time.Second, time.Millisecond*300)
}

func (s *TxsTrackingTestSuite) submitBatchInTransaction(tx models.GenericTransaction, batchType batchtype.BatchType) {
	s.runInTransaction(batchType, func(txStorage *st.Storage, txsCtx *executor.TxsContext) {
		err := txStorage.AddTransaction(tx)
		s.NoError(err)

		batchData, err := txsCtx.CreateCommitments()
		s.NoError(err)
		s.Len(batchData, 1)

		batch, err := txsCtx.NewPendingBatch(batchType)
		s.NoError(err)
		err = txsCtx.SubmitBatch(batch, batchData)
		s.NoError(err)
		s.client.GetBackend().Commit()
	})
}

func (s *TxsTrackingTestSuite) runInTransaction(batchType batchtype.BatchType, handler func(*st.Storage, *executor.TxsContext)) {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	defer txController.Rollback(nil)

	executionCtx := executor.NewTestExecutionContext(txStorage, s.client.Client, s.cfg.Rollup)
	txsCtx := executor.NewTestTxsContext(executionCtx, batchType)
	handler(txStorage, txsCtx)
}

func (s *TxsTrackingTestSuite) startWorkers() {
	s.cmd.startWorker("Test Txs Tracking", func() error {
		err := s.cmd.txsTracking()
		s.NoError(err)
		return nil
	})
	s.cmd.startWorker("Test New Block Loop", func() error {
		err := s.cmd.newBlockLoop()
		s.NoError(err)
		return nil
	})
}

func (s *TxsTrackingTestSuite) waitForLatestBlockSync() {
	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	s.Eventually(func() bool {
		syncedBlock, err := s.cmd.storage.GetSyncedBlock()
		s.NoError(err)
		return *syncedBlock >= *latestBlockNumber
	}, time.Hour, 100*time.Millisecond, "timeout when waiting for latest block sync")
}

func (s *TxsTrackingTestSuite) setAccountsAndChainState() {
	setChainState(s.T(), s.storage)
	setAccountLeaves(s.T(), s.storage.Storage, s.wallets)
}

func newClientWithGenesisStateWithClientConfig(t *testing.T, storage *st.TestStorage, config *eth.ClientConfig) *eth.TestClient {
	if config == nil {
		config = &eth.ClientConfig{}
	}
	setStateLeaves(t, storage.Storage)
	genesisRoot, err := storage.StateTree.Root()
	require.NoError(t, err)

	client, err := eth.NewConfiguredTestClient(rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: genesisRoot,
		},
	}, *config)
	require.NoError(t, err)

	return client
}

func TestTxsTrackingTestSuite(t *testing.T) {
	suite.Run(t, new(TxsTrackingTestSuite))
}
