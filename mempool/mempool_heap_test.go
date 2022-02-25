package mempool

import (
	"fmt"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MempoolHeapTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	txs     []models.GenericTransaction
}

func (s *MempoolHeapTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MempoolHeapTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.txs = []models.GenericTransaction{
		s.newTransfer(0, 10, 100), // 0 - executable
		s.newTransfer(0, 11, 100), // 1
		s.newTransfer(0, 13, 100), // 2 - non-executable

		s.newTransfer(1, 10, 80), // 3 - executable
		s.newTransfer(1, 12, 80), // 4 - non-executable

		s.newTransfer(2, 15, 120), // 5 - executable
		s.newTransfer(2, 16, 90),  // 6

		s.newC2T(3, 10, 130), // 7 - executable
		s.newC2T(3, 11, 130), // 8
	}

	err = s.storage.BatchAddTransaction(models.GenericArray(s.txs))
	s.NoError(err)

	setUserStates(s.Assertions, s.storage.StateTree, map[uint32]uint64{
		0: 10,
		1: 10,
		2: 15,
		3: 10,
	})
}

func (s *MempoolHeapTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MempoolHeapTestSuite) Test_MempoolAndHeapRealUsage() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)
	s.Equal(9, mempool.TxCount())

	txs := mempool.GetExecutableTxs(txtype.Transfer)
	heap := NewTxHeap(txs...)

	txController, txMempool := mempool.BeginTransaction()
	s.createBatch(heap, txMempool)
	txController.Commit()

	s.Equal(0, heap.Size())
	s.Equal(6, mempool.TxCount())
	s.NotContains(mempool.buckets, 2)

	s.Equal([]models.GenericTransaction{s.txs[1], s.txs[2]}, mempool.buckets[0].txs)
	s.Equal([]models.GenericTransaction{s.txs[3], s.txs[4]}, mempool.buckets[1].txs)
	s.Equal([]models.GenericTransaction{s.txs[7], s.txs[8]}, mempool.buckets[3].txs)

	s.EqualValues(10, mempool.buckets[0].nonce)
	s.EqualValues(10, mempool.buckets[1].nonce)
	s.EqualValues(10, mempool.buckets[3].nonce)
}

func (s *MempoolHeapTestSuite) createBatch(heap *TxHeap, mempool *TxMempool) {
	txController, txMempool := mempool.BeginTransaction()
	s.createCommitment(heap, txMempool)
	txController.Commit()

	txController, txMempool = mempool.BeginTransaction()
	err := s.tryCreatingSecondCommitment(heap, txMempool)
	if err != nil {
		txController.Rollback()
	}
	txController.Commit()
}

func (s *MempoolHeapTestSuite) createCommitment(heap *TxHeap, mempool *TxMempool) {
	tx := heap.Peek() // 5
	s.Equal(s.txs[5], tx)
	// execute 5 -- success
	tx, err := mempool.GetNextExecutableTx(txtype.Transfer, tx.GetFromStateID()) // 6
	s.NoError(err)
	s.Equal(s.txs[6], tx)
	heap.Replace(tx)

	tx = heap.Peek() // 0
	s.Equal(s.txs[0], tx)
	// execute 0 -- failure
	err = mempool.RemoveFailedTx(tx.GetFromStateID())
	s.NoError(err)
	heap.Pop()

	tx = heap.Peek() // 6
	s.Equal(s.txs[6], tx)
	// execute 6 -- success
	tx, err = mempool.GetNextExecutableTx(txtype.Transfer, tx.GetFromStateID()) // nil
	s.NoError(err)
	s.Nil(tx)
	heap.Pop()

	// finishing because MaxTxsPerCommitment = 2
}

func (s *MempoolHeapTestSuite) tryCreatingSecondCommitment(heap *TxHeap, mempool *TxMempool) error {
	tx := heap.Peek() // 3
	s.Equal(s.txs[3], tx)
	// execute 3 -- success
	tx, err := mempool.GetNextExecutableTx(txtype.Transfer, tx.GetFromStateID()) // nil
	s.NoError(err)
	s.Nil(tx)
	heap.Pop()

	tx = heap.Peek() // nil
	s.Nil(tx)

	// no more txs and MinTxsPerCommitment = 2
	return fmt.Errorf("not enough transacitons")
}

func (s *MempoolHeapTestSuite) newTransfer(from uint32, nonce, fee uint64) *models.Transfer {
	transfer := testutils.NewTransfer(from, 1, nonce, 100)
	transfer.Fee = models.MakeUint256(fee)
	return transfer
}

func (s *MempoolHeapTestSuite) newC2T(from uint32, nonce, fee uint64) *models.Create2Transfer {
	c2t := testutils.NewCreate2Transfer(from, ref.Uint32(1), nonce, 100, nil)
	c2t.Fee = models.MakeUint256(fee)
	return c2t
}

func TestMempoolHeapTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolHeapTestSuite))
}
