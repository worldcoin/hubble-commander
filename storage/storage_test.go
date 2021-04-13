package storage

import (
	"errors"
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
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *StorageTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *StorageTestSuite) Test_BeginTransaction_Commit() {
	leaf := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			AccountIndex: 1,
			TokenIndex:   models.MakeUint256(1),
			Balance:      models.MakeUint256(420),
			Nonce:        models.MakeUint256(0),
		},
	}

	tx, storage, err := s.storage.BeginTransaction()
	s.NoError(err)
	err = storage.AddStateLeaf(leaf)
	s.NoError(err)

	res, err := s.storage.GetStateLeaf(leaf.DataHash)
	s.Error(err)
	s.Nil(res)

	err = tx.Commit()
	s.NoError(err)

	res, err = s.storage.GetStateLeaf(leaf.DataHash)
	s.NoError(err)
	s.Equal(leaf, res)
}

func (s *StorageTestSuite) Test_BeginTransaction_Rollback() {
	leaf := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			AccountIndex: 1,
			TokenIndex:   models.MakeUint256(1),
			Balance:      models.MakeUint256(420),
			Nonce:        models.MakeUint256(0),
		},
	}

	tx, storage, err := s.storage.BeginTransaction()
	s.NoError(err)
	err = storage.AddStateLeaf(leaf)
	s.NoError(err)

	tx.Rollback(&err)
	s.Nil(errors.Unwrap(err))

	res, err := s.storage.GetStateLeaf(leaf.DataHash)
	s.Error(err)
	s.Nil(res)
}

func (s *StorageTestSuite) Test_BeginTransaction_Lock() {
	leafOne := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			AccountIndex: 1,
			TokenIndex:   models.MakeUint256(1),
			Balance:      models.MakeUint256(420),
			Nonce:        models.MakeUint256(0),
		},
	}
	leafTwo := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
		UserState: models.UserState{
			AccountIndex: 2,
			TokenIndex:   models.MakeUint256(1),
			Balance:      models.MakeUint256(1000),
			Nonce:        models.MakeUint256(0),
		},
	}

	tx, storage, err := s.storage.BeginTransaction()
	s.NoError(err)

	err = storage.AddStateLeaf(leafOne)
	s.NoError(err)

	nestedTx, nestedStorage, err := storage.BeginTransaction()
	s.NoError(err)

	err = nestedStorage.AddStateLeaf(leafTwo)
	s.NoError(err)

	nestedTx.Rollback(&err)
	s.NoError(err)

	res, err := s.storage.GetStateLeaf(leafOne.DataHash)
	s.Error(err)
	s.Nil(res)

	err = tx.Commit()
	s.NoError(err)

	res, err = s.storage.GetStateLeaf(leafOne.DataHash)
	s.NoError(err)
	s.Equal(leafOne, res)

	res, err = s.storage.GetStateLeaf(leafTwo.DataHash)
	s.NoError(err)
	s.Equal(leafTwo, res)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
