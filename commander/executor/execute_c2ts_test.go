package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type ExecuteCreate2TransfersTestSuite struct {
	testSuiteWithTxsContext
	feeReceiver *FeeReceiver
}

func (s *ExecuteCreate2TransfersTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTestWithConfig(batchtype.Create2Transfer, &config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		MaxTxsPerCommitment: 6,
	})

	setInitialUserStates(s.Assertions, s.storage.Storage)

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_AllValid() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)

	transfers, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 3)
	s.Len(transfers.InvalidTxs(), 0)
	s.Len(transfers.AddedPubKeyIDs(), 3)
	s.Len(transfers.PendingAccounts(), 1)
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_SomeValid() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(2)
	generatedTransfers = append(generatedTransfers, testutils.GenerateInvalidCreate2Transfers(3)...)

	transfers, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 2)
	s.Len(transfers.InvalidTxs(), 3)
	s.Len(transfers.AddedPubKeyIDs(), 2)
	s.Len(transfers.PendingAccounts(), 1)
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_ExecutesNoMoreThanLimit() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(7)

	transfers, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 6)
	s.Len(transfers.InvalidTxs(), 0)
	s.Len(transfers.AddedPubKeyIDs(), 6)
	s.Len(transfers.PendingAccounts(), 1)
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_SavesTxErrors() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)
	generatedTransfers = append(generatedTransfers, testutils.GenerateInvalidCreate2Transfers(2)...)

	for i := range generatedTransfers {
		err := s.storage.AddTransaction(&generatedTransfers[i])
		s.NoError(err)
	}

	result, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(result.AppliedTxs(), 3)
	s.Len(result.InvalidTxs(), 2)
	s.Len(result.AddedPubKeyIDs(), 3)
	s.Len(result.PendingAccounts(), 1)
	s.Len(s.txsCtx.txErrorsToStore, 2)

	for i := 0; i < result.InvalidTxs().Len(); i++ {
		s.Equal(generatedTransfers[i+3].Hash, s.txsCtx.txErrorsToStore[i].TxHash)
		s.Equal(applier.ErrNonceTooLow.Error(), s.txsCtx.txErrorsToStore[i].ErrorMessage)
	}
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_AppliesFee() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)

	_, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.executionCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_AddsAccountsToAccountTree() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)
	generatedTransfers[0].ToPublicKey = models.PublicKey{1, 1, 1}
	generatedTransfers[1].ToPublicKey = models.PublicKey{2, 2, 2}
	generatedTransfers[2].ToPublicKey = models.PublicKey{3, 3, 3}

	transfers, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 3)
	s.Len(transfers.InvalidTxs(), 0)
	s.Len(transfers.AddedPubKeyIDs(), 3)
	s.Len(transfers.PendingAccounts(), 3)

	for i := range generatedTransfers {
		s.Equal(transfers.PendingAccounts()[i], models.AccountLeaf{
			PubKeyID:  transfers.AddedPubKeyIDs()[i],
			PublicKey: generatedTransfers[i].ToPublicKey,
		})
	}
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_SkipsNonceTooHighTx() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(2)
	generatedTransfers[1].Nonce = models.MakeUint256(21)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(generatedTransfers, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 1)
	s.Len(executeTxsResult.SkippedTxs(), 1)
	s.Equal(*executeTxsResult.SkippedTxs().At(0).ToCreate2Transfer(), generatedTransfers[1])
}

func TestExecuteCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteCreate2TransfersTestSuite))
}
