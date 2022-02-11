package mempool

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MempoolTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.TestStorage
	initialTransactions []models.GenericTransaction
	initialNonces       map[uint32]uint64
}

func (s *MempoolTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.initialTransactions = []models.GenericTransaction{
		testutils.NewTransfer(0, 1, 13, 1), // gap
		testutils.NewTransfer(0, 1, 11, 1),
		testutils.NewTransfer(0, 1, 10, 1), // executable

		testutils.NewTransfer(1, 1, 13, 1),
		testutils.NewTransfer(1, 1, 12, 1), // non-executable

		testutils.NewTransfer(2, 1, 16, 1),
		testutils.NewTransfer(2, 1, 15, 1), // executable
	}
	s.initialNonces = map[uint32]uint64{}
	s.initialNonces[0] = 10
	s.initialNonces[1] = 11
	s.initialNonces[2] = 15
}

func (s *MempoolTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	err = s.storage.BatchAddTransaction(models.GenericArray(s.initialTransactions))
	s.NoError(err)

	s.addInitialStateLeaves(map[uint32]uint64{
		0: 10,
		1: 11,
		2: 15,
	})
}

func (s *MempoolTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MempoolTestSuite) TestNewMempool() {
	mempool, err := NewMempool(s.storage.Storage)
	s.NoError(err)

	executable := mempool.getExecutableTxs(txtype.Transfer)
	s.Len(executable, 2)
	s.Contains(executable, s.initialTransactions[2])
	s.Contains(executable, s.initialTransactions[6])
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

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}
