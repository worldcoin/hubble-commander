package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ExecuteTransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	storage     *st.TestStorage
	cfg         *config.RollupConfig
	txsCtx      *TxsContext
	feeReceiver *FeeReceiver
}

func (s *ExecuteTransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ExecuteTransfersTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		MaxTxsPerCommitment: 6,
	}

	setInitialUserStates(s.Assertions, s.storage.Storage)

	executionCtx := NewTestExecutionContext(s.storage.Storage, nil, s.cfg)
	s.txsCtx = NewTestTxsContext(executionCtx, batchtype.Transfer)

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func setInitialUserStates(s *require.Assertions, storage *st.Storage) {
	senderState := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	receiverState := models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	feeReceiverState := models.UserState{
		PubKeyID: 3,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	}

	_, err := storage.StateTree.Set(1, &senderState)
	s.NoError(err)
	_, err = storage.StateTree.Set(2, &receiverState)
	s.NoError(err)
	_, err = storage.StateTree.Set(3, &feeReceiverState)
	s.NoError(err)
}

func (s *ExecuteTransfersTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_AllValid() {
	generatedTransfers := testutils.GenerateValidTransfers(3)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 3)
	s.Len(executeTxsResult.InvalidTxs(), 0)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_SomeValid() {
	generatedTransfers := testutils.GenerateValidTransfers(2)
	generatedTransfers = append(generatedTransfers, testutils.GenerateInvalidTransfers(3)...)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 2)
	s.Len(executeTxsResult.InvalidTxs(), 3)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_ExecutesNoMoreThanLimit() {
	generatedTransfers := testutils.GenerateValidTransfers(13)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 6)
	s.Len(executeTxsResult.InvalidTxs(), 0)

	state, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_SavesTxErrors() {
	generatedTransfers := testutils.GenerateValidTransfers(3)
	generatedTransfers = append(generatedTransfers, testutils.GenerateInvalidTransfers(2)...)

	err := s.storage.BatchAddTransaction(generatedTransfers)
	s.NoError(err)

	result, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(result.AppliedTxs(), 3)
	s.Len(result.InvalidTxs(), 2)
	s.Len(s.txsCtx.txErrorsToStore, 2)

	for i := 0; i < result.InvalidTxs().Len(); i++ {
		s.Equal(generatedTransfers[i+3].Hash, s.txsCtx.txErrorsToStore[i].TxHash)
		s.Equal(applier.ErrNonceTooLow.Error(), s.txsCtx.txErrorsToStore[i].ErrorMessage)
	}
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_AppliesFee() {
	generatedTransfers := testutils.GenerateValidTransfers(3)

	_, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.txsCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_SkipsNonceTooHighTx() {
	generatedTransfers := testutils.GenerateValidTransfers(2)
	generatedTransfers[1].Nonce = models.MakeUint256(21)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 1)
	s.Len(executeTxsResult.SkippedTxs(), 1)
	s.Equal(*executeTxsResult.SkippedTxs().At(0).ToTransfer(), generatedTransfers[1])
}

func TestExecuteTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteTransfersTestSuite))
}
