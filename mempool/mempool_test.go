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

type MempoolTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	txs     []models.GenericTransaction
}

func (s *MempoolTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MempoolTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	// no need to shuffle initial transactions as they are retrieved from DB sorted by tx hashes which are random
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

	s.setUserStates(map[uint32]uint64{
		0: 10,
		1: 10,
		2: 15,
		3: 10,
	})
}

func (s *MempoolTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MempoolTestSuite) TestMempool() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	txs := mempool.GetExecutableTxs(txtype.Transfer)
	heap := NewTxHeap(txs...)

	txController, txMempool := mempool.BeginTransaction()
	s.createBatch(heap, txMempool)
	txController.Commit()

	s.Equal(0, heap.Size())
	s.Equal([]models.GenericTransaction{s.txs[1], s.txs[2]}, mempool.buckets[0].txs)
	s.Equal([]models.GenericTransaction{s.txs[3], s.txs[4]}, mempool.buckets[1].txs)
	s.Equal([]models.GenericTransaction{}, mempool.buckets[2].txs)
	s.Equal([]models.GenericTransaction{s.txs[7], s.txs[8]}, mempool.buckets[3].txs)

	s.EqualValues(10, mempool.buckets[0].nonce)
	s.EqualValues(10, mempool.buckets[1].nonce)
	s.EqualValues(17, mempool.buckets[2].nonce)
	s.EqualValues(10, mempool.buckets[3].nonce)
}

func (s *MempoolTestSuite) createBatch(heap *TxHeap, mempool *TxMempool) {
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

func (s *MempoolTestSuite) createCommitment(heap *TxHeap, mempool *TxMempool) {
	tx := heap.Peek() // 5
	s.Equal(s.txs[5], tx)
	// execute 5 -- success
	tx = mempool.GetNextExecutableTx(txtype.Transfer, tx.GetFromStateID()) // 6
	s.Equal(s.txs[6], tx)
	heap.Replace(tx)

	tx = heap.Peek() // 0
	s.Equal(s.txs[0], tx)
	// execute 0 -- failure
	mempool.RemoveFailedTx(tx.GetFromStateID())
	heap.Pop()

	tx = heap.Peek() // 6
	s.Equal(s.txs[6], tx)
	// execute 6 -- success
	tx = mempool.GetNextExecutableTx(txtype.Transfer, tx.GetFromStateID()) // nil
	s.Nil(tx)
	heap.Pop()

	// MaxTxsPerCommitment = 2
}

func (s *MempoolTestSuite) tryCreatingSecondCommitment(heap *TxHeap, mempool *TxMempool) error {
	tx := heap.Peek() // 3
	s.Equal(s.txs[3], tx)
	// execute 3 -- success
	tx = mempool.GetNextExecutableTx(txtype.Transfer, tx.GetFromStateID()) // nil
	s.Nil(tx)
	heap.Pop()

	tx = heap.Peek() // nil
	s.Nil(tx)

	// no more txs and MinTxsPerCommitment = 2
	return fmt.Errorf("not enough transacitons")
}

func (s *MempoolTestSuite) setUserStates(nonces map[uint32]uint64) {
	for stateID, nonce := range nonces {
		_, err := s.storage.StateTree.Set(stateID, &models.UserState{
			PubKeyID: 0,
			TokenID:  models.MakeUint256(uint64(stateID)),
			Balance:  models.MakeUint256(1000),
			Nonce:    models.MakeUint256(nonce),
		})
		s.NoError(err)
	}
}

func (s *MempoolTestSuite) newTransfer(from uint32, nonce, fee uint64) *models.Transfer {
	transfer := testutils.NewTransfer(from, 1, nonce, 100)
	transfer.Fee = models.MakeUint256(fee)
	return transfer
}

func (s *MempoolTestSuite) newC2T(from uint32, nonce, fee uint64) *models.Create2Transfer {
	c2t := testutils.NewCreate2Transfer(from, ref.Uint32(1), nonce, 100, nil)
	c2t.Fee = models.MakeUint256(fee)
	return c2t
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}
