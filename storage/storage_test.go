package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *StorageTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StorageTestSuite) SetupTest() {
	testDB, err := db.GetTestDB()
	s.NoError(err)
	s.storage = &Storage{DB: testDB.DB}
	s.db = testDB
}

func (s *StorageTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *StorageTestSuite) TestAddTransaction() {
	tx := &models.Transaction{
		Hash:      common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		FromIndex: models.MakeUint256(1),
		ToIndex:   models.MakeUint256(2),
		Amount:    models.MakeUint256(1000),
		Fee:       models.MakeUint256(100),
		Nonce:     models.MakeUint256(0),
		Signature: []byte{1, 2, 3, 4, 5},
	}
	err := s.storage.AddTransaction(tx)
	s.NoError(err)

	res, err := s.storage.GetTransaction(tx.Hash)
	s.NoError(err)

	s.Equal(tx, res)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
