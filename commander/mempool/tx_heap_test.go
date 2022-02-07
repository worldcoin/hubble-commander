package mempool

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxHeapTestSuite struct {
	*require.Assertions
	suite.Suite
}

func (s *TxHeapTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxHeapTestSuite) TestPop_OrdersCorrectly() {
	txs := s.makeTestTxs()
	heap := NewTxHeap(txs...)

	initialSize := heap.Size()
	orderedFees := make([]uint64, initialSize)
	for i := 0; i < initialSize; i++ {
		orderedFees[i] = heap.Pop().GetBase().Fee.Uint64()
	}
	require.Equal(s.T(), []uint64{20, 10, 9, 6, 5, 5, 4, 3, 3, 2, 2, 1}, orderedFees)
}

func (s *TxHeapTestSuite) makeTestTxs() []models.GenericTransaction {
	fees := []uint64{3, 2, 20, 5, 3, 1, 2, 5, 6, 9, 10, 4}
	txs := make([]models.GenericTransaction, len(fees))
	for i, fee := range fees {
		txs[i] = s.newTx(fee)
	}
	return txs
}

func (s *TxHeapTestSuite) newTx(fee uint64) models.GenericTransaction {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	tx.Fee = models.MakeUint256(fee)
	return &tx
}

func TestTxHeapTestSuite(t *testing.T) {
	suite.Run(t, new(TxHeapTestSuite))
}
