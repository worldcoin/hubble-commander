package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *TransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *TransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *TransactionTestSuite) Test_AddTransaction_AddAndRetrieve() {
	tx := &models.Transaction{
		Hash:                 common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		FromIndex:            models.MakeUint256(1),
		ToIndex:              models.MakeUint256(2),
		Amount:               models.MakeUint256(1000),
		Fee:                  models.MakeUint256(100),
		Nonce:                models.MakeUint256(0),
		Signature:            []byte{1, 2, 3, 4, 5},
		IncludedInCommitment: common.Hash{},
	}
	err := s.storage.AddTransaction(tx)
	s.NoError(err)

	res, err := s.storage.GetTransaction(tx.Hash)
	s.NoError(err)

	s.Equal(tx, res)
}

func (s *TransactionTestSuite) Test_GetTransaction_NonExistentTransaction() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetTransaction(hash)
	s.NoError(err)
	s.Nil(res)
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
