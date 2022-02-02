package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
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
	s.storage, err = NewTestStorage()
	s.NoError(err)
	s.batch = &models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
		PrevStateRoot:     utils.NewRandomHash(),
		MinedTime:         models.NewTimestamp(time.Unix(140, 0).UTC()),
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

	tx, txStorage := s.storage.BeginTransaction(TxOptions{})
	_, err = txStorage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)
	err = txStorage.AddBatch(s.batch)
	s.NoError(err)

	res, err := s.storage.StateTree.Leaf(leaf.StateID)
	s.ErrorIs(err, NewNotFoundError("state leaf"))
	s.Nil(res)

	batch, err := s.storage.GetBatch(s.batch.ID)
	s.ErrorIs(err, NewNotFoundError("batch"))
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

	tx, txStorage := s.storage.BeginTransaction(TxOptions{})
	_, err := txStorage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)
	err = txStorage.AddBatch(s.batch)
	s.NoError(err)

	tx.Rollback(&err)
	s.Nil(errors.Unwrap(err))

	res, err := s.storage.StateTree.Leaf(leaf.StateID)
	s.ErrorIs(err, NewNotFoundError("state leaf"))
	s.Nil(res)

	batch, err := s.storage.GetBatch(s.batch.ID)
	s.ErrorIs(err, NewNotFoundError("batch"))
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

	tx, txStorage := s.storage.BeginTransaction(TxOptions{})

	_, err = txStorage.StateTree.Set(leafOne.StateID, &leafOne.UserState)
	s.NoError(err)

	nestedTx, nestedStorage := txStorage.BeginTransaction(TxOptions{})

	_, err = nestedStorage.StateTree.Set(leafTwo.StateID, &leafTwo.UserState)
	s.NoError(err)

	nestedTx.Rollback(&err)
	s.NoError(err)

	res, err := s.storage.StateTree.Leaf(leafOne.StateID)
	s.ErrorIs(err, NewNotFoundError("state leaf"))
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

func (s *StorageTestSuite) TestExecuteInTransaction_RepliesTransactionOnConflict() {
	state := &models.ChainState{
		ChainID: models.MakeUint256(1),
	}
	err := s.storage.SetChainState(state)
	s.NoError(err)

	executions := 0
	err = s.storage.ExecuteInTransaction(TxOptions{}, func(txStorage *Storage) error {
		if executions == 0 {
			// Use s.storage to start a new DB transaction that commits before the tx started by ExecuteInTransaction
			// call above. This will cause a Transaction Conflict error.
			state.ChainID = models.MakeUint256(2)
			err = s.storage.SetChainState(state)
			s.NoError(err)
		}

		state.ChainID = models.MakeUint256(3)
		err = txStorage.SetChainState(state)
		s.NoError(err)

		executions++
		return nil
	})
	s.NoError(err)
	s.Equal(2, executions)

	retrievedState, err := s.storage.GetChainState()
	s.NoError(err)

	s.Equal(state, retrievedState)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
