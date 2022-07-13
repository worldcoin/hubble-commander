package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RevertBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	storage      *st.TestStorage
	executionCtx *ExecutionContext
	txsCtx       *TxsContext
	transfer     models.Transfer
}

func (s *RevertBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RevertBatchesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.executionCtx = NewTestExecutionContext(s.storage.Storage, eth.DomainOnlyTestClient, &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
	})
	s.txsCtx, err = NewTestTxsContext(s.executionCtx, batchtype.Transfer)
	s.NoError(err)

	s.transfer = testutils.MakeTransfer(0, 1, 0, 400)
	err = populateAccounts(s.storage.Storage, []models.Uint256{models.MakeUint256(1000), models.MakeUint256(0)})
	s.NoError(err)
}

func (s *RevertBatchesTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *RevertBatchesTestSuite) TestRevertBatches_RevertsState() {
	initialStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	pendingBatch := s.addTxBatch(&s.transfer)
	err = s.executionCtx.RevertBatches(pendingBatch)
	s.NoError(err)

	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(*initialStateRoot, *stateRoot)

	state0, err := s.storage.StateTree.Leaf(s.transfer.FromStateID)
	s.NoError(err)
	s.Equal(uint64(1000), state0.Balance.Uint64())

	state1, err := s.storage.StateTree.Leaf(s.transfer.ToStateID)
	s.NoError(err)
	s.Equal(uint64(0), state1.Balance.Uint64())
}

func (s *RevertBatchesTestSuite) TestRevertBatches_ExcludesTransactionsFromCommitments() {
	pendingBatch := s.addTxBatch(&s.transfer)
	err := s.executionCtx.RevertBatches(pendingBatch)
	s.NoError(err)

	// TODO: eventually this should put the transaction back into the mempool
	transfer, err := s.storage.GetTransfer(s.transfer.Hash)
	s.Error(err)
	s.Nil(transfer)
}

func (s *RevertBatchesTestSuite) TestRevertBatches_DeletesCommitmentsAndBatches() {
	transfers := make([]models.Transfer, 2)
	transfers[0] = s.transfer
	transfers[1] = testutils.MakeTransfer(0, 1, 1, 200)

	pendingBatches := make([]models.Batch, 2)
	for i := range pendingBatches {
		pendingBatches[i] = *s.addTxBatch(&transfers[i])
	}

	latestCommitment, err := s.executionCtx.storage.GetLatestCommitment()
	s.NoError(err)
	s.Equal(models.MakeUint256(2), latestCommitment.ID.BatchID)

	err = s.executionCtx.RevertBatches(&pendingBatches[0])
	s.NoError(err)

	_, err = s.executionCtx.storage.GetLatestCommitment()
	s.ErrorIs(err, st.NewNotFoundError("commitment"))

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 0)
}

func (s *RevertBatchesTestSuite) TestRevertBatches_AddsPendingDepositSubtree() {
	subtree := &models.PendingDepositSubtree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: testutils.GetFourDeposits(),
	}
	pendingBatch := s.addDepositBatch(subtree)
	err := s.executionCtx.RevertBatches(pendingBatch)
	s.NoError(err)

	depositSubtree, err := s.storage.GetPendingDepositSubtree(subtree.ID)
	s.NoError(err)
	s.Equal(subtree.Root, depositSubtree.Root)
	s.Equal(subtree.Deposits, depositSubtree.Deposits)
}

func (s *RevertBatchesTestSuite) addTxBatch(tx *models.Transfer) *models.Batch {
	initTxs(s.Assertions, s.txsCtx, models.TransferArray{*tx})

	pendingBatch, err := s.txsCtx.NewPendingBatch(s.txsCtx.BatchType)
	s.NoError(err)

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)

	err = s.storage.AddBatch(pendingBatch)
	s.NoError(err)

	err = s.txsCtx.addCommitments(commitments)
	s.NoError(err)

	return pendingBatch
}

func (s *RevertBatchesTestSuite) addDepositBatch(subtree *models.PendingDepositSubtree) *models.Batch {
	pendingBatch, err := s.executionCtx.NewPendingBatch(batchtype.Deposit)
	s.NoError(err)

	deposits := testutils.GetFourDeposits()
	err = s.executionCtx.Applier.ApplyDeposits(2, deposits)
	s.NoError(err)

	root, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	err = s.storage.AddBatch(pendingBatch)
	s.NoError(err)

	err = s.executionCtx.storage.AddCommitment(&models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID: pendingBatch.ID,
			},
			Type:          batchtype.Deposit,
			PostStateRoot: *root,
		},
		SubtreeID:   subtree.ID,
		SubtreeRoot: subtree.Root,
		Deposits:    subtree.Deposits,
	})
	s.NoError(err)
	return pendingBatch
}

func TestRevertBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(RevertBatchesTestSuite))
}
