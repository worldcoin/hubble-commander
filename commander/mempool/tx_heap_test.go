package mempool

import (
	"fmt"
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

func (s *TxHeapTestSuite) TestHeap() {
	txs := s.makeTestTxs()
	heap := NewTxHeap(txs...)

	for i := 0; i < heap.Size(); i++ {
		fmt.Printf("%s,", heap.Pop().GetBase().Fee.String())
	}
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
