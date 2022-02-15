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
}

func (s *TxsTrackingTestSuite) setupTestWithClientConfig(conf *eth.TestClientConfig) {
	s.cmd = NewCommander(s.cfg, nil)
	// pass txs channels to testClient to use commander tracking worker
	conf.TxsChannels = s.cmd.txsTrackingChannels

	s.client = newClientWithGenesisStateWithClientConfig(s.T(), s.storage, conf)

	setStateLeaves(s.T(), s.storage.Storage)
	s.cmd.client = s.client.Client
	s.cmd.blockchain = s.client.Blockchain
	s.cmd.storage = s.storage.Storage

	err := s.cmd.addGenesisBatch()
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)
	s.setAccountsAndChainState()

	s.startWorkers()
	s.waitForLatestBlockSync()
}

func (s *TxsTrackingTestSuite) setupTestWithFailedTxs() {
	lowGasLimit := uint64(40_000)
	s.setupTestWithClientConfig(&eth.TestClientConfig{
		ClientConfig: eth.ClientConfig{
			TransferBatchSubmissionGasLimit:  &lowGasLimit,
			C2TBatchSubmissionGasLimit:       &lowGasLimit,
			MMBatchSubmissionGasLimit:        &lowGasLimit,
			BatchAccountRegistrationGasLimit: &lowGasLimit,
			StakeWithdrawalGasLimit:          &lowGasLimit,
			DepositBatchSubmissionGasLimit:   &lowGasLimit,
		},
	})
}

func (s *TxsTrackingTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxsTrackingTestSuite) TestStartFailedTxsTracking_TransferTransaction() {
	s.setupTestWithFailedTxs()

	transfer := testutils.MakeTransfer(0, 1, 0, 400)
	s.submitBatchInTransaction(&transfer, batchtype.Transfer)

	s.waitForWorkersCancellation()
}

func (s *TxsTrackingTestSuite) TestStartFailedTxsTracking_Create2TransfersTransaction() {
	s.setupTestWithFailedTxs()

	transfer := testutils.MakeCreate2Transfer(
		0,
		ref.Uint32(1),
		0,
		50,
		&models.PublicKey{2, 3, 4},
	)
	s.submitBatchInTransaction(&transfer, batchtype.Create2Transfer)

	s.waitForWorkersCancellation()
}

func (s *TxsTrackingTestSuite) TestStartFailedTxsTracking_MassMigrationTransaction() {
	s.setupTestWithFailedTxs()

	massMigration := testutils.MakeMassMigration(0, 2, 0, 50)
	s.submitBatchInTransaction(&massMigration, batchtype.MassMigration)

	s.waitForWorkersCancellation()
}

func (s *TxsTrackingTestSuite) TestStartFailedTxsTracking_BatchAccountRegistrationTransaction() {
	s.setupTestWithFailedTxs()

	publicKeys := make([]models.PublicKey, st.AccountBatchSize)
	_, err := s.client.Client.RegisterBatchAccount(publicKeys)
	s.NoError(err)

	s.waitForWorkersCancellation()
}

func (s *TxsTrackingTestSuite) TestStartFailedTxsTracking_WithdrawStake() {
	s.setupTestWithFailedTxs()

	transfer := testutils.MakeTransfer(0, 1, 0, 400)
	batch := s.submitBatchInTransaction(&transfer, batchtype.Transfer)

	_, err := s.client.Client.WithdrawStake(&batch.ID)
	s.NoError(err)

	s.waitForWorkersCancellation()
}

func (s *TxsTrackingTestSuite) TestStartFailedTxsTracking_SubmitDepositBatch() {
	s.setupTestWithFailedTxs()

	err := s.storage.AddPendingDepositSubtree(&models.PendingDepositSubtree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: testutils.GetFourDeposits(),
	})
	s.NoError(err)
	executionCtx := executor.NewTestExecutionContext(s.storage.Storage, s.client.Client, nil)
	depositsCtx := executor.NewTestDepositsContext(executionCtx)

	_, _, err = depositsCtx.CreateAndSubmitBatch()
	s.NoError(err)

	s.waitForWorkersCancellation()
}

func (s *TxsTrackingTestSuite) waitForWorkersCancellation() {
	s.Eventually(func() bool {
		err := s.cmd.workersContext.Err()
		return err == context.Canceled
	}, time.Second, time.Millisecond*300)
}

func (s *TxsTrackingTestSuite) submitBatchInTransaction(
	tx models.GenericTransaction,
	batchType batchtype.BatchType,
) (batch *models.Batch) {
	s.runInTransaction(batchType, func(txStorage *st.Storage, txsCtx *executor.TxsContext) {
		err := txStorage.AddTransaction(tx)
		s.NoError(err)

		batchData, err := txsCtx.CreateCommitments()
		s.NoError(err)
		s.Len(batchData, 1)

		batch, err = txsCtx.NewPendingBatch(batchType)
		s.NoError(err)
		err = txsCtx.SubmitBatch(batch, batchData)
		s.NoError(err)
		s.client.GetBackend().Commit()
	})
	return batch
}

func (s *TxsTrackingTestSuite) runInTransaction(
	batchType batchtype.BatchType,
	handler func(*st.Storage, *executor.TxsContext),
) {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	defer txController.Rollback(nil)

	executionCtx := executor.NewTestExecutionContext(txStorage, s.client.Client, s.cfg.Rollup)
	txsCtx := executor.NewTestTxsContext(executionCtx, batchType)
	handler(txStorage, txsCtx)
}

func (s *TxsTrackingTestSuite) startWorkers() {
	s.cmd.startWorker("Test Txs Requests Sending", func() error {
		err := s.cmd.startTxsRequestsSending()
		s.NoError(err)
		return nil
	})
	s.cmd.startWorker("Test Failed Txs Tracking", func() error {
		err := s.cmd.startFailedTxsTracking()
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

func newClientWithGenesisStateWithClientConfig(t *testing.T, storage *st.TestStorage, conf *eth.TestClientConfig) *eth.TestClient {
	setStateLeaves(t, storage.Storage)
	genesisRoot, err := storage.StateTree.Root()
	require.NoError(t, err)

	client, err := eth.NewConfiguredTestClient(&rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: genesisRoot,
		},
	}, conf)

	require.NoError(t, err)

	return client
}

func TestTxsTrackingTestSuite(t *testing.T) {
	suite.Run(t, new(TxsTrackingTestSuite))
}
