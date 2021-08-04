package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
	batch   *models.Batch
}

func (s *StorageTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StorageTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)
	s.batch = &models.Batch{
		ID:                models.MakeUint256(1),
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
		PrevStateRoot:     utils.NewRandomHash(),
		SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
	}
}

func (s *StorageTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StorageTestSuite) TestBeginTransaction_Commit() {
	leaf, err := NewStateLeaf(0, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	tx, txStorage, err := s.storage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	s.NoError(err)
	_, err = txStorage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)
	err = txStorage.AddBatch(s.batch)
	s.NoError(err)

	res, err := s.storage.StateTree.Leaf(leaf.StateID)
	s.Equal(NewNotFoundError("state leaf"), err)
	s.Nil(res)

	batch, err := s.storage.GetBatch(s.batch.ID)
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(batch)

	err = tx.Commit()
	s.NoError(err)

	res, err = s.storage.StateTree.Leaf(leaf.StateID)
	s.NoError(err)
	s.Equal(leaf, res)

	batch, err = s.storage.GetBatch(s.batch.ID)
	s.NoError(err)
	s.Equal(s.batch, batch)
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

	tx, txStorage, err := s.storage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	s.NoError(err)
	_, err = txStorage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)
	err = txStorage.AddBatch(s.batch)
	s.NoError(err)

	tx.Rollback(&err)
	s.Nil(errors.Unwrap(err))

	res, err := s.storage.StateTree.Leaf(leaf.StateID)
	s.Equal(NewNotFoundError("state leaf"), err)
	s.Nil(res)

	batch, err := s.storage.GetBatch(s.batch.ID)
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(batch)
}

func (s *StorageTestSuite) TestBeginTransaction_Lock() {
	leafOne, err := NewStateLeaf(0, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	leafTwo, err := NewStateLeaf(1, &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	tx, txStorage, err := s.storage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	s.NoError(err)

	_, err = txStorage.StateTree.Set(leafOne.StateID, &leafOne.UserState)
	s.NoError(err)

	nestedTx, nestedStorage, err := txStorage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	s.NoError(err)

	_, err = nestedStorage.StateTree.Set(leafTwo.StateID, &leafTwo.UserState)
	s.NoError(err)

	nestedTx.Rollback(&err)
	s.NoError(err)

	res, err := s.storage.StateTree.Leaf(leafOne.StateID)
	s.Equal(NewNotFoundError("state leaf"), err)
	s.Nil(res)

	err = tx.Commit()
	s.NoError(err)

	res, err = s.storage.StateTree.Leaf(leafOne.StateID)
	s.NoError(err)
	s.Equal(leafOne, res)

	res, err = s.storage.StateTree.Leaf(leafTwo.StateID)
	s.NoError(err)
	s.Equal(leafTwo, res)
}

func (s *StorageTestSuite) TestClone() {
	testConfig := config.GetTestConfig().Postgres

	batch := models.Batch{
		ID:              models.MakeUint256(1),
		Type:            txtype.Transfer,
		TransactionHash: utils.RandomHash(),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	stateLeaf, err := NewStateLeaf(1, &models.UserState{})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(stateLeaf.StateID, &stateLeaf.UserState)
	s.NoError(err)

	clonedStorage, err := s.storage.Clone(testConfig)
	s.NoError(err)
	defer func() {
		err = clonedStorage.Teardown()
		s.NoError(err)
	}()

	clonedBatch, err := clonedStorage.GetBatch(batch.ID)
	s.NoError(err)
	s.Equal(batch, *clonedBatch)

	clonedStateLeaf, err := clonedStorage.StateTree.Leaf(stateLeaf.StateID)
	s.NoError(err)
	s.Equal(stateLeaf, clonedStateLeaf)
}

func (s *StorageTestSuite) TestClone_ClonesFeeReceiverStateIDsByValue() {
	testConfig := config.GetTestConfig().Postgres

	s.storage.feeReceiverStateIDs["abc"] = 123

	clonedStorage, err := s.storage.Clone(testConfig)
	s.NoError(err)
	defer func() {
		err = clonedStorage.Teardown()
		s.NoError(err)
	}()

	abcID, ok := clonedStorage.feeReceiverStateIDs["abc"]
	s.True(ok)
	s.EqualValues(123, abcID)

	clonedStorage.feeReceiverStateIDs["def"] = 456

	defID, ok := s.storage.feeReceiverStateIDs["def"]
	s.False(ok)
	s.Nil(defID)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
