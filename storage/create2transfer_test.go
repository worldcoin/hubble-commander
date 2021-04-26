package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	create2Transfer = models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 common.BigToHash(big.NewInt(1234)),
			TxType:               txtype.Create2Transfer,
			FromStateID:          1,
			Amount:               models.MakeUint256(1000),
			Fee:                  models.MakeUint256(100),
			Nonce:                models.MakeUint256(0),
			Signature:            []byte{1, 2, 3, 4, 5},
			IncludedInCommitment: nil,
		},
		ToStateID:  2,
		ToPubKeyID: 2,
	}
)

type Create2TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
	tree    *StateTree
}

func (s *Create2TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *Create2TransferTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
	s.tree = NewStateTree(s.storage)

	err = s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)
}

func (s *Create2TransferTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *Create2TransferTestSuite) TestAddCreate2Transfer_AddAndRetrieve() {
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	res, err := s.storage.GetCreate2Transfer(create2Transfer.Hash)
	s.NoError(err)

	s.Equal(create2Transfer, *res)
}

func (s *Create2TransferTestSuite) TestGetCreate2Transfer_NonExistentTransaction() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetCreate2Transfer(hash)
	s.Equal(NewNotFoundError("transaction"), err)
	s.Nil(res)
}

func TestCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferTestSuite))
}
