package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StateUpdateTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *StateUpdateTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateUpdateTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)
}

func (s *StateUpdateTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StateUpdateTestSuite) TestAddStateUpdate_AddAndRetrieve() {
	update := &models.StateUpdate{
		ID:          0,
		StateID:     12,
		CurrentRoot: common.BytesToHash([]byte{1, 2, 3}),
		PrevRoot:    common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		PrevStateLeaf: models.StateLeaf{
			StateID:   12,
			DataHash:  [32]byte{1,2,3,4},
			UserState: models.UserState{
				PubKeyID:   1,
				TokenIndex: models.MakeUint256(1),
				Balance:    models.MakeUint256(100),
				Nonce:      models.MakeUint256(0),
			},
		},
	}
	err := s.storage.AddStateUpdate(update)
	s.NoError(err)

	res, err := s.storage.GetStateUpdateByRootHash(common.BytesToHash([]byte{1, 2, 3}))
	s.NoError(err)
	s.Equal(update, res)
}

func (s *StateUpdateTestSuite) TestGetStateUpdateByRootHash_NonExistentUpdate() {
	res, err := s.storage.GetStateUpdateByRootHash(common.BytesToHash([]byte{9, 4, 1, 2}))
	s.Equal(NewNotFoundError("state update"), err)
	s.Nil(res)
}

func (s *StateUpdateTestSuite) TestDeleteStateUpdate() {
	updates := []models.StateUpdate{
		{
			ID:          0,
			StateID:     1,
			CurrentRoot: common.BytesToHash([]byte{1}),
			PrevRoot:    common.BytesToHash([]byte{2}),
			PrevStateLeaf: models.StateLeaf{
				StateID:   12,
				DataHash:  [32]byte{1,2,3,4},
				UserState: models.UserState{
					PubKeyID:   1,
					TokenIndex: models.MakeUint256(1),
					Balance:    models.MakeUint256(100),
					Nonce:      models.MakeUint256(0),
				},
			},
		},
		{
			ID:          1,
			StateID:     1,
			CurrentRoot: common.BytesToHash([]byte{2}),
			PrevRoot:    common.BytesToHash([]byte{2}),
			PrevStateLeaf: models.StateLeaf{
				StateID:   12,
				DataHash:  [32]byte{1,2,3,4},
				UserState: models.UserState{
					PubKeyID:   1,
					TokenIndex: models.MakeUint256(1),
					Balance:    models.MakeUint256(100),
					Nonce:      models.MakeUint256(0),
				},
			},
		},
	}
	err := s.storage.AddStateUpdate(&updates[0])
	s.NoError(err)
	err = s.storage.AddStateUpdate(&updates[1])
	s.NoError(err)

	err = s.storage.DeleteStateUpdate(1)
	s.NoError(err)

	_, err = s.storage.GetStateUpdateByRootHash(updates[1].CurrentRoot)
	s.Equal(NewNotFoundError("state update"), err)

	res, err := s.storage.GetStateUpdateByRootHash(updates[0].CurrentRoot)
	s.NoError(err)
	s.Equal(&updates[0], res)
}

func TestStateUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(StateUpdateTestSuite))
}
