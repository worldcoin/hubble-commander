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

	s.initialTransactions = []models.GenericTransaction{
		s.newTransfer(0, 13), // 0 - gap
		s.newTransfer(0, 11), // 1
		s.newTransfer(0, 10), // 2 - executable

		s.newTransfer(1, 12), // 3
		s.newTransfer(1, 11), // 4 - non-executable

		s.newTransfer(2, 16), // 5
		s.newTransfer(2, 15), // 6 - executable

		s.newC2T(3, 11), // 7
		s.newC2T(3, 10), // 8 - executable
	}
}

func (s *MempoolTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	err = s.storage.BatchAddTransaction(models.GenericArray(s.initialTransactions))
	s.NoError(err)

	s.addInitialStateLeaves(map[uint32]uint64{
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

	s.Len(mempool.userTxsMap, 4)
	s.Equal(mempool.userTxsMap[0].txs, []models.GenericTransaction{
		s.initialTransactions[2],
		s.initialTransactions[1],
		s.initialTransactions[0],
	})
	s.Len(mempool.userTxsMap[1].txs, 2)
	s.Len(mempool.userTxsMap[2].txs, 2)
	s.Len(mempool.userTxsMap[3].txs, 2)
}

func (s *MempoolTestSuite) TestGetExecutableTxs_ReturnsExecutableTxsOfCorrectType() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	executable := mempool.getExecutableTxs(txtype.Transfer)
	s.Len(executable, 2)
	s.Contains(executable, s.initialTransactions[2])
	s.Contains(executable, s.initialTransactions[6])

	executable = mempool.getExecutableTxs(txtype.Create2Transfer)
	s.Len(executable, 1)
	s.Contains(executable, s.initialTransactions[8])
}

//func (s *MempoolTestSuite) TestAddTransaction() {
//	mempool, err := NewMempool(s.storage.Storage)
//	s.NoError(err)
//
//	tx := createTx(3, 10)
//	mempool.addOrReplace(tx, 10)
//
//	executable := mempool.getExecutableTxs(txtype.Transfer)
//	s.Len(executable, 3)
//	s.Equal(s.initialTransactions[0], executable[0])
//	s.Equal(s.initialTransactions[5], executable[1])
//	s.Equal(tx, executable[2])
//}
//
//func (s *MempoolTestSuite) TestReplaceTransaction() {
//	mempool, err := NewMempool(s.storage.Storage)
//	s.NoError(err)
//
//	tx := createTx(0, 10)
//	mempool.addOrReplace(tx, 10)
//
//	executable := mempool.getExecutableTxs(txtype.Transfer)
//	s.Len(executable, 2)
//	s.Equal(tx, executable[0])
//	s.Equal(s.initialTransactions[5], executable[1])
//}

func (s *MempoolTestSuite) addInitialStateLeaves(nonces map[uint32]uint64) {
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
