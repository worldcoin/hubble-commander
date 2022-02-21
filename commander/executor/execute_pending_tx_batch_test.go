package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ExecutePendingTxBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	storage      *st.TestStorage
	cfg          *config.RollupConfig
	txsCtx       *TxsContext
	pendingBatch models.PendingBatch
}

func (s *ExecutePendingTxBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ExecutePendingTxBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		MaxTxsPerCommitment: 1,
	}

	setInitialUserStates(s.Assertions, s.storage.Storage)

	executionCtx := NewTestExecutionContext(s.storage.Storage, nil, s.cfg)
	s.txsCtx, err = NewTestTxsContext(executionCtx, batchtype.Transfer)
	s.NoError(err)

	s.pendingBatch = models.PendingBatch{
		ID:              models.MakeUint256(1),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
		PrevStateRoot:   utils.RandomHash(),
		Commitments: []models.PendingCommitment{
			{
				Commitment: &models.TxCommitment{
					CommitmentBase: models.CommitmentBase{
						ID: models.CommitmentID{
							BatchID:      models.MakeUint256(1),
							IndexInBatch: 0,
						},
						Type:          batchtype.Transfer,
						PostStateRoot: utils.RandomHash(),
					},
					FeeReceiver:       3,
					CombinedSignature: models.MakeRandomSignature(),
					BodyHash:          nil,
				},
			},
		},
	}
}

func (s *ExecutePendingTxBatchTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ExecutePendingTxBatchTestSuite) TestExecutePendingBatch_UpdatesUserBalances() {
	tx := testutils.MakeTransfer(1, 2, 0, 100)
	tx.CommitmentID = &s.pendingBatch.Commitments[0].GetCommitmentBase().ID
	s.pendingBatch.Commitments[0].Transactions = models.TransferArray{tx}

	prevSenderLeaf, err := s.storage.StateTree.Leaf(tx.FromStateID)
	s.NoError(err)
	prevReceiverLeaf, err := s.storage.StateTree.Leaf(tx.ToStateID)
	s.NoError(err)

	err = s.txsCtx.ExecutePendingBatch(&s.pendingBatch)
	s.NoError(err)

	senderLeaf, err := s.storage.StateTree.Leaf(tx.FromStateID)
	s.NoError(err)
	expectedSenderBalance := prevSenderLeaf.Balance.Sub(&tx.Amount).Sub(&tx.Fee)
	s.Equal(*expectedSenderBalance, senderLeaf.Balance)

	receiverLeaf, err := s.storage.StateTree.Leaf(tx.ToStateID)
	s.NoError(err)
	s.Equal(*prevReceiverLeaf.Balance.Add(&tx.Amount), receiverLeaf.Balance)
}

func (s *ExecutePendingTxBatchTestSuite) TestExecutePendingBatch_AddsPendingBatch() {
	tx := testutils.MakeTransfer(1, 2, 0, 100)
	tx.CommitmentID = &s.pendingBatch.Commitments[0].GetCommitmentBase().ID
	s.pendingBatch.Commitments[0].Transactions = models.MakeGenericArray(&tx)

	err := s.txsCtx.ExecutePendingBatch(&s.pendingBatch)
	s.NoError(err)

	expectedBatch := models.Batch{
		ID:              s.pendingBatch.ID,
		Type:            s.pendingBatch.Type,
		TransactionHash: s.pendingBatch.TransactionHash,
		PrevStateRoot:   &s.pendingBatch.PrevStateRoot,
	}

	batch, err := s.storage.GetBatch(s.pendingBatch.ID)
	s.NoError(err)
	s.Equal(expectedBatch, *batch)

	commitments, err := s.storage.GetCommitmentsByBatchID(s.pendingBatch.ID)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Equal(s.pendingBatch.Commitments[0].Commitment, commitments[0])

	txs, err := s.storage.GetTransactionsByCommitmentID(commitments[0].GetCommitmentBase().ID)
	s.NoError(err)
	s.Len(txs, 1)
	s.Equal(s.pendingBatch.Commitments[0].Transactions, txs)
}

func (s *ExecutePendingTxBatchTestSuite) TestExecutePendingBatch_ErrorsOnFailureToApplyAllPendingBatchTxs() {
	tx := testutils.MakeTransfer(1, 2, 1, 100)
	tx.CommitmentID = &s.pendingBatch.Commitments[0].GetCommitmentBase().ID
	s.pendingBatch.Commitments[0].Transactions = models.TransferArray{tx}

	err := s.txsCtx.ExecutePendingBatch(&s.pendingBatch)
	s.ErrorIs(err, applier.ErrNonceTooHigh)
}

func TestExecutePendingTxBatchTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutePendingTxBatchTestSuite))
}
