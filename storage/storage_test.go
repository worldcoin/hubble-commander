package storage

import (
	"errors"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *StorageTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StorageTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)
}

func (s *StorageTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StorageTestSuite) TestBeginTransaction_Commit() {
	leaf := &models.StateLeaf{
		StateID:  0,
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(420),
			Nonce:    models.MakeUint256(0),
		},
	}

	tx, storage, err := s.storage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	s.NoError(err)
	err = storage.UpsertStateLeaf(leaf)
	s.NoError(err)
	err = storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	res, err := s.storage.GetStateLeaf(leaf.StateID)
	s.Equal(NewNotFoundError("state leaf"), err)
	s.Nil(res)

	accounts, err := s.storage.GetAccounts(&account2.PublicKey)
	s.Equal(NewNotFoundError("accounts"), err)
	s.Nil(accounts)

	err = tx.Commit()
	s.NoError(err)

	res, err = s.storage.GetStateLeaf(leaf.StateID)
	s.NoError(err)
	s.Equal(leaf, res)

	accounts, err = s.storage.GetAccounts(&account2.PublicKey)
	s.NoError(err)
	s.Len(accounts, 1)
}

func (s *StorageTestSuite) TestBeginTransaction_Rollback() {
	leaf := &models.StateLeaf{
		StateID:  0,
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(420),
			Nonce:    models.MakeUint256(0),
		},
	}

	tx, storage, err := s.storage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	s.NoError(err)
	err = storage.UpsertStateLeaf(leaf)
	s.NoError(err)
	err = storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	tx.Rollback(&err)
	s.Nil(errors.Unwrap(err))

	res, err := s.storage.GetStateLeaf(leaf.StateID)
	s.Equal(NewNotFoundError("state leaf"), err)
	s.Nil(res)

	accounts, err := s.storage.GetAccounts(&account2.PublicKey)
	s.Equal(NewNotFoundError("accounts"), err)
	s.Nil(accounts)
}

func (s *StorageTestSuite) TestBeginTransaction_Lock() {
	err := s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	leafOne := &models.StateLeaf{
		StateID:  0,
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(420),
			Nonce:    models.MakeUint256(0),
		},
	}
	leafTwo := &models.StateLeaf{
		StateID:  1,
		DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
		UserState: models.UserState{
			PubKeyID: 2,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(1000),
			Nonce:    models.MakeUint256(0),
		},
	}

	tx, storage, err := s.storage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	s.NoError(err)

	err = storage.UpsertStateLeaf(leafOne)
	s.NoError(err)

	nestedTx, nestedStorage, err := storage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	s.NoError(err)

	err = nestedStorage.UpsertStateLeaf(leafTwo)
	s.NoError(err)

	nestedTx.Rollback(&err)
	s.NoError(err)

	res, err := s.storage.GetStateLeaf(leafOne.StateID)
	s.Equal(NewNotFoundError("state leaf"), err)
	s.Nil(res)

	err = tx.Commit()
	s.NoError(err)

	res, err = s.storage.GetStateLeaf(leafOne.StateID)
	s.NoError(err)
	s.Equal(leafOne, res)

	res, err = s.storage.GetStateLeaf(leafTwo.StateID)
	s.NoError(err)
	s.Equal(leafTwo, res)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
