package commander

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	create2transfer = models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 common.BigToHash(big.NewInt(1234)),
			FromStateID:          0,
			Amount:               models.MakeUint256(1000),
			Fee:                  models.MakeUint256(100),
			Nonce:                models.MakeUint256(0),
			Signature:            []byte{1, 2, 3, 4, 5},
			IncludedInCommitment: nil,
		},
		ToPubkeyID: 2,
	}
	senderTokenIndex = models.MakeUint256(5)
)

type ApplyCreate2TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.Storage
	db      *db.TestDB
	tree    *st.StateTree
}

func (s *ApplyCreate2TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyCreate2TransferTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = st.NewTestStorage(testDB.DB)
	s.db = testDB
	s.tree = st.NewStateTree(s.storage)

	s.tree.Set(0, &models.UserState{
		PubkeyID:   0,
		TokenIndex: senderTokenIndex,
		Balance:    models.MakeUint256(10000),
		Nonce:      models.MakeUint256(0),
	})
	s.tree.Set(1, &models.UserState{
		PubkeyID:   1,
		TokenIndex: senderTokenIndex,
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	})
}

func (s *ApplyCreate2TransferTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransferTestSuite) Test_ApplyCreate2Transfer_InsertsNewEmptyStateLeaf() {
	transferError, appError := ApplyCreate2Transfer(s.storage, &create2transfer, senderTokenIndex)
	s.NoError(appError)
	s.NoError(transferError)

	leaf, err := s.tree.Leaf(2)
	s.NoError(err)
	s.NotNil(leaf)
	s.Equal(models.MakeUint256(0), leaf.Nonce)
	s.Equal(senderTokenIndex, leaf.TokenIndex)
}

func (s *ApplyCreate2TransferTestSuite) Test_ApplyCreate2Transfer_ApplyTransfer() {
	transferError, appError := ApplyCreate2Transfer(s.storage, &create2transfer, senderTokenIndex)
	s.NoError(appError)
	s.NoError(transferError)

	receiverLeaf, err := s.tree.Leaf(2)
	s.NoError(err)
	senderLeaf, err := s.tree.Leaf(create2transfer.FromStateID)
	s.NoError(err)

	s.Equal(int64(8900), senderLeaf.Balance.Int64())
	s.Equal(int64(1000), receiverLeaf.Balance.Int64())
}

func TestApplyCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransferTestSuite))
}
