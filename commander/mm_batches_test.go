package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/result"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MMBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd     *Commander
	client  *eth.TestClient
	storage *st.TestStorage
	cfg     *config.Config
}

func (s *MMBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinTxsPerCommitment = 1
}

func (s *MMBatchesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client = newClientWithGenesisState(s.T(), s.storage)

	s.cmd = NewCommander(s.cfg, nil, false)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.storage.Storage

	err = s.cmd.addGenesisBatch()
	s.NoError(err)
}

func (s *MMBatchesTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MMBatchesTestSuite) TestSyncRemoteBatch_SyncsBatch() {
	tx := testutils.MakeMassMigration(0, 1, 0, 100)
	err := s.storage.AddTransaction(&tx)
	s.NoError(err)

	s.submitBatch(s.storage.Storage)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.NoError(err)

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	depositBatch := remoteBatches[0].ToDecodedTxBatch()
	s.Equal(depositBatch.Hash, *batches[1].Hash)
	s.Equal(batchtype.MassMigration, batches[1].Type)
}

func (s *MMBatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithMismatchedTotalAmount() {
	tx := testutils.MakeMassMigration(0, 1, 0, 100)

	s.submitInvalidBatch(&tx, func(commitments []models.CommitmentWithTxs) {
		commitments[0].ToMMCommitmentWithTxs().Meta.Amount = models.MakeUint256(110)
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.MismatchedAmount, getDisputeResult(s.Assertions, s.client))
}

func (s *MMBatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithInvalidWithdrawRoot() {
	tx := testutils.MakeMassMigration(0, 1, 0, 100)

	s.submitInvalidBatch(&tx, func(commitments []models.CommitmentWithTxs) {
		commitments[0].ToMMCommitmentWithTxs().WithdrawRoot = common.Hash{1, 2, 3}
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.BadWithdrawRoot, getDisputeResult(s.Assertions, s.client))
}

func (s *MMBatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithInvalidTokenID() {
	tx := testutils.MakeMassMigration(0, 1, 0, 100)

	s.submitInvalidBatch(&tx, func(commitments []models.CommitmentWithTxs) {
		commitments[0].ToMMCommitmentWithTxs().Meta.TokenID = models.MakeUint256(110)
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.BadFromTokenID, getDisputeResult(s.Assertions, s.client))
}

func (s *MMBatchesTestSuite) submitInvalidBatch(tx *models.MassMigration, modifier func(commitments []models.CommitmentWithTxs)) {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	defer txController.Rollback(nil)

	err := txStorage.AddTransaction(tx)
	s.NoError(err)

	executionCtx := executor.NewTestExecutionContext(txStorage, s.client.Client, s.cfg.Rollup)
	txsCtx := executor.NewTestTxsContext(executionCtx, batchtype.MassMigration)

	pendingBatch, err := txsCtx.NewPendingBatch(txsCtx.BatchType)
	s.NoError(err)

	commitments, err := txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)

	modifier(commitments)

	err = txsCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.GetBackend().Commit()
}

func (s *MMBatchesTestSuite) submitBatch(storage *st.Storage) *models.Batch {
	txsCtx := executor.NewTxsContext(
		storage,
		s.client.Client,
		s.cfg.Rollup,
		metrics.NewCommanderMetrics(),
		context.Background(),
		batchtype.MassMigration,
	)
	defer txsCtx.Rollback(nil)

	batch, _, err := txsCtx.CreateAndSubmitBatch()
	s.NoError(err)

	s.client.GetBackend().Commit()
	return batch
}

func TestMMBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(MMBatchesTestSuite))
}
