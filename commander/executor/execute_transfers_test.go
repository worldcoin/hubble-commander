package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/mempool"
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
	txMempool := newMempool(s.Assertions, s.txsCtx, generatedTransfers)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 3)
	s.Len(executeTxsResult.InvalidTxs(), 0)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_SomeValid() {
	generatedTransfers := testutils.GenerateValidTransfers(2)
	generatedTransfers = append(generatedTransfers, generateInvalidTransfers(3)...)
	txMempool := newMempool(s.Assertions, s.txsCtx, generatedTransfers)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 2)
	s.Len(executeTxsResult.InvalidTxs(), 1)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_ExecutesNoMoreThanLimit() {
	generatedTransfers := testutils.GenerateValidTransfers(13)
	txMempool := newMempool(s.Assertions, s.txsCtx, generatedTransfers)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 6)
	s.Len(executeTxsResult.InvalidTxs(), 0)

	state, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_SavesTxErrors() {
	generatedTransfers := testutils.GenerateValidTransfers(3)
	generatedTransfers = append(generatedTransfers, generateInvalidTransfers(1)...)
	txMempool := newMempool(s.Assertions, s.txsCtx, generatedTransfers)

	result, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(result.AppliedTxs(), 3)
	s.Len(result.InvalidTxs(), 1)
	s.Len(s.txsCtx.txErrorsToStore, 1)
	s.Equal(generatedTransfers[3].Hash, s.txsCtx.txErrorsToStore[0].TxHash)
	s.Equal(applier.ErrBalanceTooLow.Error(), s.txsCtx.txErrorsToStore[0].ErrorMessage)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_AppliesFee() {
	generatedTransfers := testutils.GenerateValidTransfers(3)
	txMempool := newMempool(s.Assertions, s.txsCtx, generatedTransfers)

	_, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.txsCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ExecuteTransfersTestSuite) TestExecuteTxs_SkipsNonceTooHighTx() {
	generatedTransfers := testutils.GenerateValidTransfers(2)
	generatedTransfers[1].Nonce = models.MakeUint256(21)
	txMempool := newMempool(s.Assertions, s.txsCtx, generatedTransfers)

	executeTxsResult, err := s.txsCtx.ExecuteTxs(txMempool, s.feeReceiver)
	s.NoError(err)

	s.Len(executeTxsResult.AppliedTxs(), 1)
}

func newMempool(s *require.Assertions, txsCtx *TxsContext, txs models.GenericTransactionArray) *mempool.TxMempool {
	initMempool(s, txsCtx, txs)
	txsCtx.heap = txsCtx.newHeap()

	_, txMempool := txsCtx.mempool.BeginTransaction()
	return txMempool
}

func initMempool(s *require.Assertions, txsCtx *TxsContext, txs models.GenericTransactionArray) {
	if txs.Len() > 0 {
		err := txsCtx.storage.BatchAddTransaction(txs)
		s.NoError(err)
	}

	pool, err := mempool.NewMempool(txsCtx.storage)
	s.NoError(err)
	txsCtx.mempool = pool
}

// TODO: change GenerateInvalidTransfers FromStateID
func generateInvalidTransfers(transfersAmount uint64) []models.Transfer {
	txs := testutils.GenerateInvalidTransfers(transfersAmount)
	for i := range txs {
		txs[i].FromStateID = 3
	}
	return txs
}

func TestExecuteTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteTransfersTestSuite))
}
