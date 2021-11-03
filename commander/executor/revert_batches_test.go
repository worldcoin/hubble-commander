package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RevertBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	storage         *st.TestStorage
	executionCtx    *ExecutionContext
	transactionsCtx *TransactionsContext
	transfer        models.Transfer
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
	s.transactionsCtx = NewTestTransactionsContext(s.executionCtx, batchtype.Transfer)

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

	pendingBatch := s.addBatch(&s.transfer)
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
	pendingBatch := s.addBatch(&s.transfer)
	err := s.executionCtx.RevertBatches(pendingBatch)
	s.NoError(err)

	transfer, err := s.storage.GetTransfer(s.transfer.Hash)
	s.NoError(err)
	s.Nil(transfer.CommitmentID)
}

func (s *RevertBatchesTestSuite) TestRevertBatches_DeletesCommitmentsAndBatches() {
	transfers := make([]models.Transfer, 2)
	transfers[0] = s.transfer
	transfers[1] = testutils.MakeTransfer(0, 1, 1, 200)

	pendingBatches := make([]models.Batch, 2)
	for i := range pendingBatches {
		pendingBatches[i] = *s.addBatch(&transfers[i])
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

func (s *RevertBatchesTestSuite) addBatch(tx *models.Transfer) *models.Batch {
	err := s.transactionsCtx.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := s.transactionsCtx.NewPendingBatch(s.transactionsCtx.BatchType)
	s.NoError(err)

	commitmentID, err := s.transactionsCtx.NextCommitmentID()
	s.NoError(err)
	result, err := s.transactionsCtx.createCommitment(models.TransferArray{*tx}, commitmentID)
	s.NoError(err)

	err = s.storage.AddBatch(pendingBatch)
	s.NoError(err)

	err = s.transactionsCtx.addCommitments([]models.CommitmentWithTxs{*result.Commitment()})
	s.NoError(err)

	return pendingBatch
}

func TestRevertBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(RevertBatchesTestSuite))
}
