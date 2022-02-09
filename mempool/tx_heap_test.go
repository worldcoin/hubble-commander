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

func (s *TxHeapTestSuite) TestPeek() {
	heap := NewTxHeap(s.makeTestTxs()...)
	element, err := heap.Peek()
	s.NoError(err)
	s.EqualValues(20, element.GetBase().Fee.Uint64())
}

func (s *TxHeapTestSuite) TestPeek_EmptyHeap() {
	heap := NewTxHeap()
	element, err := heap.Peek()
	s.ErrorIs(err, ErrEmptyHeap)
	s.Nil(element)
}

func (s *TxHeapTestSuite) TestPush() {
	heap := NewTxHeap(s.makeTestTxs()...)
	heap.Push(s.newTx(7))
	s.Equal([]uint64{20, 10, 9, 7, 6, 5, 5, 4, 3, 3, 2, 2, 1}, s.popAll(heap))
}

func (s *TxHeapTestSuite) TestPop() {
	heap := NewTxHeap(s.makeTestTxs()...)
	s.Equal([]uint64{20, 10, 9, 6, 5, 5, 4, 3, 3, 2, 2, 1}, s.popAll(heap))
}

func (s *TxHeapTestSuite) TestPop_EmptyHeap() {
	heap := NewTxHeap()
	element, err := heap.Pop()
	s.ErrorIs(err, ErrEmptyHeap)
	s.Nil(element)
}

func (s *TxHeapTestSuite) TestReplace() {
	heap := NewTxHeap(s.makeTestTxs()...)
	newTx := s.newTx(7)
	s.EqualValues(20, heap.Replace(newTx).GetBase().Fee.Uint64())
	s.Equal([]uint64{10, 9, 7, 6, 5, 5, 4, 3, 3, 2, 2, 1}, s.popAll(heap))
}

func (s *TxHeapTestSuite) TestReplace_EmptyHeap() {
	heap := NewTxHeap()
	newTx := s.newTx(7)
	replacedElement := heap.Replace(newTx)
	s.Nil(replacedElement)
	s.Equal([]uint64{7}, s.popAll(heap))
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

func (s *TxHeapTestSuite) popAll(heap *TxHeap) []uint64 {
	initialSize := heap.Size()
	orderedFees := make([]uint64, initialSize)
	for i := 0; i < initialSize; i++ {
		element, err := heap.Pop()
		s.NoError(err)

		orderedFees[i] = element.GetBase().Fee.Uint64()
	}
	return orderedFees
}

func TestTxHeapTestSuite(t *testing.T) {
	suite.Run(t, new(TxHeapTestSuite))
}
