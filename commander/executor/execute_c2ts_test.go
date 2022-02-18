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
	s.txsCtx.mempool = newMempool(s.Assertions, s.storage, generatedTransfers)
	_, txMempool := s.txsCtx.mempool.BeginTransaction()

	transfers, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 3)
	s.Len(transfers.InvalidTxs(), 0)
	s.Len(transfers.AddedPubKeyIDs(), 3)
	s.Len(transfers.PendingAccounts(), 1)
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_SomeValid() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(2)
	generatedTransfers = append(generatedTransfers, generateInvalidCreate2Transfers(3)...)
	s.txsCtx.mempool = newMempool(s.Assertions, s.storage, generatedTransfers)
	_, txMempool := s.txsCtx.mempool.BeginTransaction()

	transfers, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 2)
	s.Len(transfers.InvalidTxs(), 1)
	s.Len(transfers.AddedPubKeyIDs(), 2)
	s.Len(transfers.PendingAccounts(), 1)
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_ExecutesNoMoreThanLimit() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(7)
	s.txsCtx.mempool = newMempool(s.Assertions, s.storage, generatedTransfers)
	_, txMempool := s.txsCtx.mempool.BeginTransaction()

	transfers, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.AppliedTxs(), 6)
	s.Len(transfers.InvalidTxs(), 0)
	s.Len(transfers.AddedPubKeyIDs(), 6)
	s.Len(transfers.PendingAccounts(), 1)
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_SavesTxErrors() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)
	generatedTransfers = append(generatedTransfers, generateInvalidCreate2Transfers(1)...)
	s.txsCtx.mempool = newMempool(s.Assertions, s.storage, generatedTransfers)
	_, txMempool := s.txsCtx.mempool.BeginTransaction()

	result, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(result.AppliedTxs(), 3)
	s.Len(result.InvalidTxs(), 1)
	s.Len(result.AddedPubKeyIDs(), 3)
	s.Len(result.PendingAccounts(), 1)
	s.Len(s.txsCtx.txErrorsToStore, 1)
	s.Equal(generatedTransfers[3].Hash, s.txsCtx.txErrorsToStore[0].TxHash)
	s.Equal(applier.ErrBalanceTooLow.Error(), s.txsCtx.txErrorsToStore[0].ErrorMessage)
}

func (s *ExecuteCreate2TransfersTestSuite) TestExecuteTxs_AppliesFee() {
	generatedTransfers := testutils.GenerateValidCreate2Transfers(3)
	s.txsCtx.mempool = newMempool(s.Assertions, s.storage, generatedTransfers)
	_, txMempool := s.txsCtx.mempool.BeginTransaction()

	_, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
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
	s.txsCtx.mempool = newMempool(s.Assertions, s.storage, generatedTransfers)
	_, txMempool := s.txsCtx.mempool.BeginTransaction()

	transfers, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
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
	s.txsCtx.mempool = newMempool(s.Assertions, s.storage, generatedTransfers)
	_, txMempool := s.txsCtx.mempool.BeginTransaction()

	executeTxsResult, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 1)
}

// TODO: change GenerateInvalidCreate2Transfers FromStateID
func generateInvalidCreate2Transfers(transfersAmount uint64) []models.Create2Transfer {
	txs := testutils.GenerateInvalidCreate2Transfers(transfersAmount)
	for i := range txs {
		txs[i].FromStateID = 3
		txs[i].Amount = models.MakeUint256(1_000_000)
	}
	return txs
}

func TestExecuteCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteCreate2TransfersTestSuite))
}
