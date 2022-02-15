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
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	executable := mempool.GetExecutableTxs(txtype.Transfer)
	s.Len(executable, 2)
	s.Contains(executable, s.txs[0])
	s.Contains(executable, s.txs[5])

	executable = mempool.GetExecutableTxs(txtype.Create2Transfer)
	s.Len(executable, 1)
	s.Contains(executable, s.txs[7])
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_ReturnsNextTx() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)
	_, txMempool := mempool.BeginTransaction()

	tx := txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.Equal(s.txs[1], tx)
}

func (s *MempoolTestSuite) TestGetNextExecutableTx_IncrementsNonce() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)
	_, txMempool := mempool.BeginTransaction()

	_ = txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.EqualValues(11, txMempool.buckets[0].nonce)
	_ = txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.EqualValues(12, txMempool.buckets[0].nonce)
}

func (s *MempoolTestSuite) TestGetExecutableTx_NoMoreTransactionsInSlice() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)
	_, txMempool := mempool.BeginTransaction()

	_ = txMempool.GetNextExecutableTx(txtype.Transfer, 2)
	tx := txMempool.GetNextExecutableTx(txtype.Transfer, 2)
	s.Nil(tx)
}

func (s *MempoolTestSuite) TestGetExecutableTx_NoMoreExecutableTxs() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)
	_, txMempool := mempool.BeginTransaction()

	_ = txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	tx := txMempool.GetNextExecutableTx(txtype.Transfer, 0)
	s.Nil(tx)
}

func (s *MempoolTestSuite) TestGetExecutableTx_NoExecutableTxsOfGivenType() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)
	_, txMempool := mempool.BeginTransaction()

	_ = txMempool.GetNextExecutableTx(txtype.Create2Transfer, 0)
	tx := txMempool.GetNextExecutableTx(txtype.Create2Transfer, 0)
	s.Nil(tx)
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
