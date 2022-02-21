package mempool

import (
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
	mempool *Mempool
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
		s.newTransfer(0, 10), // 0 - executable
		s.newTransfer(0, 11), // 1
		s.newTransfer(0, 13), // 2 - non-executable

		s.newTransfer(1, 11), // 3 - non-executable
		s.newTransfer(1, 12), // 4

		s.newTransfer(2, 15), // 5 - executable
		s.newTransfer(2, 16), // 6

		s.newC2T(3, 10),      // 7 - executable
		s.newC2T(3, 11),      // 8
		s.newTransfer(3, 12), // 9
	}

	err = s.storage.BatchAddTransaction(models.GenericArray(s.txs))
	s.NoError(err)

	setUserStates(s.Assertions, s.storage.StateTree, map[uint32]uint64{
		0: 10,
		1: 10,
		2: 15,
		3: 10,
	})

	s.mempool, err = NewMempool(s.storage.Storage)
	s.NoError(err)
}

func (s *MempoolTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MempoolTestSuite) TestNewMempool_InitsBucketsCorrectly() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	s.Len(mempool.buckets, 4)
	s.Equal(mempool.buckets[0].txs, s.txs[0:3])
	s.Equal(mempool.buckets[1].txs, s.txs[3:5])
	s.Equal(mempool.buckets[2].txs, s.txs[5:7])
	s.Equal(mempool.buckets[3].txs, s.txs[7:10])

	s.EqualValues(mempool.buckets[0].nonce, 10)
	s.EqualValues(mempool.buckets[1].nonce, 10)
	s.EqualValues(mempool.buckets[2].nonce, 15)
	s.EqualValues(mempool.buckets[3].nonce, 10)
}

func (s *MempoolTestSuite) TestGetExecutableTxs_ReturnsAllExecutableTxsOfGivenType() {
	executable := s.mempool.GetExecutableTxs(txtype.Transfer)
	s.Len(executable, 2)
	s.Contains(executable, s.txs[0])
	s.Contains(executable, s.txs[5])

	executable = s.mempool.GetExecutableTxs(txtype.Create2Transfer)
	s.Len(executable, 1)
	s.Contains(executable, s.txs[7])
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_ReturnsNextTx() {
	_, txMempool := s.mempool.BeginTransaction()

	tx, err := txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.NoError(err)
	s.Equal(s.txs[1], tx)
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_IncrementsNonce() {
	_, txMempool := s.mempool.BeginTransaction()

	_, err := txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.NoError(err)
	s.EqualValues(11, txMempool.buckets[0].nonce)
	_, err = txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.NoError(err)
	s.EqualValues(12, txMempool.buckets[0].nonce)
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_NoMoreTransactionsInSlice() {
	_, txMempool := s.mempool.BeginTransaction()

	_, err := txMempool.GetNextExecutableTx(txtype.Transfer, 2)
	s.NoError(err)
	tx, err := txMempool.GetNextExecutableTx(txtype.Transfer, 2)
	s.NoError(err)
	s.Nil(tx)
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_NoMoreExecutableTxs() {
	_, txMempool := s.mempool.BeginTransaction()

	_, err := txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.NoError(err)
	tx, err := txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.NoError(err)
	s.Nil(tx)
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_NoExecutableTxsOfGivenType() {
	_, txMempool := s.mempool.BeginTransaction()

	_, err := txMempool.GetNextExecutableTx(txtype.Create2Transfer, 0)
	s.NoError(err)
	tx, err := txMempool.GetNextExecutableTx(txtype.Create2Transfer, 0)
	s.NoError(err)
	s.Nil(tx)
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_RemovesEmptyBuckets() {
	txController, txMempool := s.mempool.BeginTransaction()

	_, err := txMempool.GetNextExecutableTx(txtype.Transfer, 2)
	s.NoError(err)
	_, err = txMempool.GetNextExecutableTx(txtype.Transfer, 2)
	s.NoError(err)
	s.Contains(txMempool.buckets, uint32(2))
	s.Nil(txMempool.buckets[2])

	txController.Commit()
	s.NotContains(s.mempool.buckets, uint32(2))
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_BucketDoesNotExist() {
	_, txMempool := s.mempool.BeginTransaction()

	tx, err := txMempool.GetNextExecutableTx(txtype.Transfer, 10)
	s.ErrorIs(err, ErrNonexistentBucket)
	s.Nil(tx)
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_DecrementsTxCount() {
	txController, txMempool := s.mempool.BeginTransaction()

	_, err := txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.NoError(err)

	s.Equal(9, txMempool.TxCount())

	txController.Commit()
	s.Equal(9, s.mempool.TxCount())
}

func (s *MempoolTestSuite) TestRemoveFailedTx_RemovesTxFromMempool() {
	_, txMempool := s.mempool.BeginTransaction()

	err := txMempool.RemoveFailedTx(0)
	s.NoError(err)
	s.Equal(txMempool.buckets[0].txs, s.txs[1:3])
}

func (s *MempoolTestSuite) TestRemoveFailedTx_MakesTheNextTxNonExecutable() {
	_, txMempool := s.mempool.BeginTransaction()

	err := txMempool.RemoveFailedTx(0)
	s.NoError(err)
	tx, err := txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.NoError(err)
	s.Nil(tx)
}

func (s *MempoolTestSuite) TestRemoveFailedTx_RemovesEmptyBuckets() {
	txController, txMempool := s.mempool.BeginTransaction()

	_, err := txMempool.GetNextExecutableTx(txtype.Transfer, 2)
	s.NoError(err)
	err = txMempool.RemoveFailedTx(2)
	s.NoError(err)
	s.Contains(txMempool.buckets, uint32(2))
	s.Nil(txMempool.buckets[2])

	txController.Commit()
	s.NotContains(s.mempool.buckets, uint32(2))
}

func (s *MempoolTestSuite) TestRemoveFailedTx_BucketDoesNotExist() {
	_, txMempool := s.mempool.BeginTransaction()

	err := txMempool.RemoveFailedTx(10)
	s.ErrorIs(err, ErrNonexistentBucket)
}

func (s *MempoolTestSuite) TestRemoveFailedTx_DecrementsTxCount() {
	txController, txMempool := s.mempool.BeginTransaction()

	err := txMempool.RemoveFailedTx(0)
	s.NoError(err)

	s.Equal(9, txMempool.TxCount())

	txController.Commit()
	s.Equal(9, s.mempool.TxCount())
}

func (s *MempoolTestSuite) TestAddOrReplace_AppendsNewTxToBucketList() {
	tx := s.newTransfer(0, 14)
	prevTxHash, err := s.mempool.AddOrReplace(s.storage.Storage, tx)
	s.NoError(err)
	s.Nil(prevTxHash)

	bucket := s.mempool.buckets[0]
	lastTxInBucket := bucket.txs[len(bucket.txs)-1]
	s.Equal(tx, lastTxInBucket)
}

func (s *MempoolTestSuite) TestAddOrReplace_InsertsNewTxInTheMiddleOfBucketList() {
	tx := s.newTransfer(0, 12)
	prevTxHash, err := s.mempool.AddOrReplace(s.storage.Storage, tx)
	s.NoError(err)
	s.Nil(prevTxHash)

	bucket := s.mempool.buckets[0]
	s.Equal(tx, bucket.txs[2])
}

func (s *MempoolTestSuite) TestAddOrReplace_ReplacesTx() {
	tx := s.newTransfer(0, 11)
	tx.Fee = models.MakeUint256(20)
	prevTxHash, err := s.mempool.AddOrReplace(s.storage.Storage, tx)
	s.NoError(err)
	s.Equal(s.txs[1].GetBase().Hash, *prevTxHash)

	bucket := s.mempool.buckets[0]
	s.Equal(tx, bucket.txs[1])
}

func (s *MempoolTestSuite) TestAddOrReplace_ReturnsErrorOnFeeTooLowToReplace() {
	tx := s.newTransfer(0, 11)
	prevTxHash, err := s.mempool.AddOrReplace(s.storage.Storage, tx)
	s.ErrorIs(err, ErrTxReplacementFailed)
	s.Nil(prevTxHash)
}

func (s *MempoolTestSuite) TestAddOrReplace_IncrementsTxCountOnInsertion() {
	tx := s.newTransfer(0, 14)
	_, err := s.mempool.AddOrReplace(s.storage.Storage, tx)
	s.NoError(err)

	s.Equal(11, s.mempool.TxCount())
}

func (s *MempoolTestSuite) TestAddOrReplace_DoesNotIncrementTxCountOnReplacement() {
	tx := s.newTransfer(0, 11)
	tx.Fee = models.MakeUint256(20)
	_, err := s.mempool.AddOrReplace(s.storage.Storage, tx)
	s.NoError(err)

	s.Equal(10, s.mempool.TxCount())
}

func (s *MempoolTestSuite) TestTxCount() {
	s.Equal(10, s.mempool.TxCount())
}

func (s *MempoolTestSuite) newTransfer(from uint32, nonce uint64) *models.Transfer {
	return testutils.NewTransfer(from, 1, nonce, 100)
}

func (s *MempoolTestSuite) newC2T(from uint32, nonce uint64) *models.Create2Transfer {
	return testutils.NewCreate2Transfer(from, ref.Uint32(1), nonce, 100, nil)
}

func setUserStates(s *require.Assertions, stateTree *st.StateTree, nonces map[uint32]uint64) {
	for stateID, nonce := range nonces {
		_, err := stateTree.Set(stateID, &models.UserState{
			PubKeyID: 0,
			TokenID:  models.MakeUint256(uint64(stateID)),
			Balance:  models.MakeUint256(1000),
			Nonce:    models.MakeUint256(nonce),
		})
		s.NoError(err)
	}
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}
