package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/erc20"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd            *Commander
	client         *eth.TestClient
	storage        *st.TestStorage
	depositSubtree models.PendingDepositSubtree
	cfg            *config.Config
}

func (s *DepositBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinTxsPerCommitment = 1
}

func (s *DepositBatchesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client = newClientWithGenesisState(s.T(), s.storage)

	s.cmd = NewCommander(s.cfg, s.client.Blockchain)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.storage.Storage
	s.cmd.txsTrackingChannels = s.client.TxsChannels

	err = s.cmd.addGenesisBatch()
	s.NoError(err)

	s.depositSubtree = models.PendingDepositSubtree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: testutils.GetFourDeposits(),
	}
}

func newClientWithGenesisState(t *testing.T, storage *st.TestStorage) *eth.TestClient {
	setStateLeaves(t, storage.Storage)
	genesisRoot, err := storage.StateTree.Root()
	require.NoError(t, err)

	client, err := eth.NewConfiguredTestClient(&rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: genesisRoot,
		},
	}, &eth.TestClientConfig{})
	require.NoError(t, err)

	return client
}

func (s *DepositBatchesTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositBatchesTestSuite) TestSyncRemoteBatch_SyncsDepositBatch() {
	s.prepareDeposits()
	s.submitDepositBatch(s.storage.Storage)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.NoError(err)

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	depositBatch := remoteBatches[0].ToDecodedDepositBatch()
	s.Equal(depositBatch.Hash, *batches[1].Hash)
	s.Equal(depositBatch.Type, batches[1].Type)
}

func (s *DepositBatchesTestSuite) TestUnsafeSyncBatches_OmitsRolledBackBatch() {
	s.prepareDeposits()
	s.submitInvalidBatches()

	latestBlock, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	// trigger dispute on fraudulent batch
	err = s.cmd.unsafeSyncBatches(0, *latestBlock)
	s.ErrorIs(err, ErrRollbackInProgress)

	depositBatch := s.submitDepositBatch(s.storage.Storage)
	latestBlock, err = s.client.GetLatestBlockNumber()
	s.NoError(err)

	// try syncing already rolled back batch
	err = s.cmd.unsafeSyncBatches(0, *latestBlock)
	s.NoError(err)

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)
	s.Equal(depositBatch.TransactionHash, batches[1].TransactionHash)
}

func (s *DepositBatchesTestSuite) submitInvalidBatches() {
	previousRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	defer txController.Rollback(nil)

	executionCtx := executor.NewTestExecutionContext(txStorage, s.client.Client, s.cfg.Rollup)
	txsCtx := executor.NewTestTxsContext(executionCtx, batchtype.Transfer)
	invalidTransfer := testutils.MakeTransfer(0, 1, 0, 100)
	submitInvalidTxsBatch(s.Assertions, txStorage, txsCtx, &invalidTransfer, func(_ *st.Storage, commitment *models.TxCommitmentWithTxs) {
		commitment.Transactions = append(commitment.Transactions, commitment.Transactions...)
	})
	s.client.Blockchain.GetBackend().Commit()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	txsSyncCtx := syncer.NewTestTxsContext(txStorage, s.client.Client, s.cfg.Rollup, txtype.Transfer)
	err = txsSyncCtx.UpdateExistingBatch(remoteBatches[0], *previousRoot)
	s.NoError(err)

	s.submitDepositBatch(txStorage)
}

func (s *DepositBatchesTestSuite) submitDepositBatch(storage *st.Storage) *models.Batch {
	s.queueFourDeposits()

	depositsCtx := executor.NewDepositsContext(
		storage,
		s.client.Client,
		s.cfg.Rollup,
		metrics.NewCommanderMetrics(),
		context.Background(),
	)
	defer depositsCtx.Rollback(nil)

	batch, _, err := depositsCtx.CreateAndSubmitBatch()
	s.NoError(err)

	s.client.GetBackend().Commit()
	return batch
}

func (s *DepositBatchesTestSuite) prepareDeposits() {
	err := s.storage.AddPendingDepositSubtree(&s.depositSubtree)
	s.NoError(err)

	s.approveTokens()
}

func (s *DepositBatchesTestSuite) approveTokens() {
	token, err := erc20.NewERC20(s.client.ExampleTokenAddress, s.client.GetBackend())
	s.NoError(err)

	_, err = token.Approve(s.client.GetAccount(), s.client.ChainState.DepositManager, utils.ParseEther("100"))
	s.NoError(err)

	s.client.GetBackend().Commit()
}

func (s *DepositBatchesTestSuite) queueFourDeposits() {
	for i := 0; i < 4; i++ {
		s.queueDeposit()
	}
}

func (s *DepositBatchesTestSuite) queueDeposit() {
	toPubKeyID := models.NewUint256(1)
	tokenID := models.NewUint256(0)
	l1Amount := models.NewUint256FromBig(*utils.ParseEther("10"))
	_, _, err := s.client.QueueDepositAndWait(toPubKeyID, l1Amount, tokenID)
	s.NoError(err)
}

func TestDepositBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(DepositBatchesTestSuite))
}
