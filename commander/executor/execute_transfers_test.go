package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ExecuteTransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	storage     *storage.TestStorage
	cfg         *config.RollupConfig
	rollupCtx   *RollupContext
	feeReceiver *FeeReceiver
}

func (s *ExecuteTransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ExecuteTransfersTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorage()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		MaxTxsPerCommitment: 6,
	}

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

	_, err = s.storage.StateTree.Set(1, &senderState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(2, &receiverState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(3, &feeReceiverState)
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, nil, s.cfg)
	s.rollupCtx = NewTestRollupContext(executionCtx, batchtype.Transfer)

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *ExecuteTransfersTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_AllValid() {
	generatedTransfers := testutils.GenerateValidTransfers(3)

	executeTxsResult, err := s.rollupCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 3)
	s.Len(executeTxsResult.InvalidTxs(), 0)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_SomeValid() {
	generatedTransfers := testutils.GenerateValidTransfers(2)
	generatedTransfers = append(generatedTransfers, testutils.GenerateInvalidTransfers(3)...)

	executeTxsResult, err := s.rollupCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 2)
	s.Len(executeTxsResult.InvalidTxs(), 3)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_ExecutesNoMoreThanLimit() {
	generatedTransfers := testutils.GenerateValidTransfers(13)

	executeTxsResult, err := s.rollupCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 6)
	s.Len(executeTxsResult.InvalidTxs(), 0)

	state, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_SavesTransferErrors() {
	generatedTransfers := testutils.GenerateValidTransfers(3)
	generatedTransfers = append(generatedTransfers, testutils.GenerateInvalidTransfers(2)...)

	for i := range generatedTransfers {
		err := s.storage.AddTransfer(&generatedTransfers[i])
		s.NoError(err)
	}

	executeTxsResult, err := s.rollupCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 3)
	s.Len(executeTxsResult.InvalidTxs(), 2)

	for i := range generatedTransfers {
		transfer, err := s.storage.GetTransfer(generatedTransfers[i].Hash)
		s.NoError(err)
		if i < 3 {
			s.Nil(transfer.ErrorMessage)
		} else {
			s.Equal(*transfer.ErrorMessage, applier.ErrNonceTooLow.Error())
		}
	}
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_AppliesFee() {
	generatedTransfers := testutils.GenerateValidTransfers(3)

	_, err := s.rollupCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.rollupCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func TestExecuteTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteTransfersTestSuite))
}
