package storage

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type NewTestStorageTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
	batch   *models.Batch
}

func (s *NewTestStorageTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *NewTestStorageTestSuite) SetupTest() {
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
		PrevStateRoot:     utils.RandomHash(),
		MinedTime:         models.NewTimestamp(time.Unix(140, 0).UTC()),
	}
}

func (s *NewTestStorageTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *NewTestStorageTestSuite) TestClone() {
	err := s.storage.AddBatch(s.batch)
	s.NoError(err)

	stateLeaf, err := NewStateLeaf(1, &models.UserState{})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(stateLeaf.StateID, &stateLeaf.UserState)
	s.NoError(err)

	clonedStorage, err := s.storage.Clone()
	s.NoError(err)
	defer func() {
		err = clonedStorage.Teardown()
		s.NoError(err)
	}()

	clonedBatch, err := clonedStorage.GetBatch(s.batch.ID)
	s.NoError(err)
	s.Equal(s.batch, clonedBatch)

	clonedStateLeaf, err := clonedStorage.StateTree.Leaf(stateLeaf.StateID)
	s.NoError(err)
	s.Equal(stateLeaf, clonedStateLeaf)
}

func (s *NewTestStorageTestSuite) TestClone_ClonesFeeReceiverStateIDsByValue() {
	s.storage.feeReceiverStateIDs["abc"] = 123

	clonedStorage, err := s.storage.Clone()
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
	s.EqualValues(0, defID) // empty value
}

func TestNewTestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(NewTestStorageTestSuite))
}
