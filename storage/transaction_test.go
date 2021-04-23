package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
	tree    *StateTree
}

func (s *TransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
	s.tree = NewStateTree(s.storage)
}

func (s *TransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *TransactionTestSuite) TestGetTransaction_Transfer() {
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	res, err := s.storage.GetTransaction(transfer.Hash)
	s.NoError(err)

	s.Equal(&transfer, res)
}

func (s *TransactionTestSuite) TestGetTransaction_Create2Transfer() {
	err := s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	err = s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	res, err := s.storage.GetTransaction(create2Transfer.Hash)
	s.NoError(err)

	s.Equal(&create2Transfer, res)
}

func (s *TransactionTestSuite) TestGetTransaction_NonExistingTransaction() {
	res, err := s.storage.GetTransaction(transfer.Hash)
	s.True(IsNotFoundError(err))
	s.Nil(res)
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
