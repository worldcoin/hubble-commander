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
	storage             *st.TestStorage
	initialTransactions []models.GenericTransaction
}

func (s *MempoolTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MempoolTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	// no need to shuffle initial transactions as they are retrieved from DB sorted by tx hashes which are random
	s.initialTransactions = []models.GenericTransaction{
		s.newTransfer(0, 10), // 0 - executable
		s.newTransfer(0, 11), // 1
		s.newTransfer(0, 13), // 2 - non-executable

		s.newTransfer(1, 11), // 3 - non-executable
		s.newTransfer(1, 12), // 4

		s.newTransfer(2, 15), // 5 - executable
		s.newTransfer(2, 16), // 6

		s.newC2T(3, 10), // 7 - executable
		s.newC2T(3, 11), // 8
	}

	err = s.storage.BatchAddTransaction(models.GenericArray(s.initialTransactions))
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

func (s *MempoolTestSuite) TestNewMempool_InitsBucketsCorrectly() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	s.Len(mempool.buckets, 4)
	s.Equal(mempool.buckets[0].txs, s.initialTransactions[0:3])
	s.Equal(mempool.buckets[1].txs, s.initialTransactions[3:5])
	s.Equal(mempool.buckets[2].txs, s.initialTransactions[5:7])
	s.Equal(mempool.buckets[3].txs, s.initialTransactions[7:9])

	s.EqualValues(mempool.buckets[0].nonce, 10)
	s.EqualValues(mempool.buckets[1].nonce, 10)
	s.EqualValues(mempool.buckets[2].nonce, 15)
	s.EqualValues(mempool.buckets[3].nonce, 10)

	s.Equal(mempool.buckets[0].executableIndex, 0)
	s.Equal(mempool.buckets[1].executableIndex, nonExecutableIndex)
	s.Equal(mempool.buckets[2].executableIndex, 0)
	s.Equal(mempool.buckets[3].executableIndex, 0)
}

func (s *MempoolTestSuite) TestGetExecutableTx_ReturnsNextExecutableTxOfGivenStateID() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	tx := mempool.GetExecutableTx(txtype.Transfer, 0)
	s.Equal(s.initialTransactions[0], tx)
}

func (s *MempoolTestSuite) TestGetExecutableTx_NoMoreTxs() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	_ = mempool.GetExecutableTx(txtype.Create2Transfer, 3)
	_ = mempool.GetExecutableTx(txtype.Create2Transfer, 3)
	tx := mempool.GetExecutableTx(txtype.Create2Transfer, 3)
	s.Nil(tx)
	s.Equal(nonExecutableIndex, mempool.buckets[3].executableIndex)
}

func (s *MempoolTestSuite) TestGetExecutableTx_NoMoreExecutableTxs() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	_ = mempool.GetExecutableTx(txtype.Transfer, 0)
	_ = mempool.GetExecutableTx(txtype.Transfer, 0)
	tx := mempool.GetExecutableTx(txtype.Transfer, 0)
	s.Nil(tx)
	s.Equal(nonExecutableIndex, mempool.buckets[0].executableIndex)
}

func (s *MempoolTestSuite) TestGetExecutableTx_NoExecutableTxsOfGivenType() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	tx := mempool.GetExecutableTx(txtype.Create2Transfer, 0)
	s.Nil(tx)
	s.Equal(0, mempool.buckets[0].executableIndex)
}

func (s *MempoolTestSuite) TestGetExecutableTxs_ReturnsExecutableTxsOfCorrectType() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	executable := mempool.GetExecutableTxs(txtype.Transfer)
	s.Len(executable, 2)
	s.Contains(executable, s.initialTransactions[0])
	s.Contains(executable, s.initialTransactions[5])

	executable = mempool.GetExecutableTxs(txtype.Create2Transfer)
	s.Len(executable, 1)
	s.Contains(executable, s.initialTransactions[7])
}

func (s *MempoolTestSuite) TestGetExecutableTxs_UpdatesExecutableIndices() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	mempool.GetExecutableTxs(txtype.Transfer)
	s.Equal(1, mempool.buckets[0].executableIndex)
	s.Equal(1, mempool.buckets[2].executableIndex)

	mempool.GetExecutableTxs(txtype.Create2Transfer)
	s.Equal(1, mempool.buckets[3].executableIndex)
}

func (s *MempoolTestSuite) TestAddOrReplace_AppendsNewTxToBucketList() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	tx := s.newTransfer(0, 14)
	err = mempool.AddOrReplace(tx, 10)
	s.NoError(err)

	bucket := mempool.buckets[0]
	lastTxInBucket := bucket.txs[len(bucket.txs)-1]
	s.Equal(tx, lastTxInBucket)
}

func (s *MempoolTestSuite) TestAddOrReplace_InsertsNewTxInTheMiddleOfBucketList() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	tx := s.newTransfer(0, 12)
	err = mempool.AddOrReplace(tx, 10)
	s.NoError(err)

	bucket := mempool.buckets[0]
	s.Equal(tx, bucket.txs[2])
}

func (s *MempoolTestSuite) TestAddOrReplace_ReplacesTx() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	tx := s.newTransfer(0, 11)
	tx.Fee = models.MakeUint256(20)
	err = mempool.AddOrReplace(tx, 10)
	s.NoError(err)

	bucket := mempool.buckets[0]
	s.Equal(tx, bucket.txs[1])
}

func (s *MempoolTestSuite) TestAddOrReplace_ReturnsErrorOnFeeTooLowToReplace() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	tx := s.newTransfer(0, 11)
	err = mempool.AddOrReplace(tx, 10)
	s.ErrorIs(err, ErrTxReplacementFailed)
}

func (s *MempoolTestSuite) TestIgnoreUserTxs() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	mempool.IgnoreUserTxs(0)
	s.Equal(nonExecutableIndex, mempool.buckets[0].executableIndex)
}

func (s *MempoolTestSuite) TestResetExecutableIndices() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	mempool.buckets[0].executableIndex = 2
	mempool.buckets[2].executableIndex = nonExecutableIndex

	mempool.ResetExecutableIndices()
	s.Equal(0, mempool.buckets[0].executableIndex)
	s.Equal(0, mempool.buckets[2].executableIndex)
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

func (s *MempoolTestSuite) newTransfer(from uint32, nonce uint64) *models.Transfer {
	return testutils.NewTransfer(from, 1, nonce, 100)
}

func (s *MempoolTestSuite) newC2T(from uint32, nonce uint64) *models.Create2Transfer {
	return testutils.NewCreate2Transfer(from, ref.Uint32(1), nonce, 100, nil)
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}
