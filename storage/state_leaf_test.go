package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	userState1 = &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	userState2 = &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	userState3 = &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(2),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
)

type StateLeafTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *StateLeafTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateLeafTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *StateLeafTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StateLeafTestSuite) TestUpsertStateLeaf_AddAndRetrieve() {
	leaf, err := NewStateLeaf(0, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	err = s.storage.StateTree.upsertStateLeaf(leaf)
	s.NoError(err)

	res, err := s.storage.StateTree.Leaf(leaf.StateID)
	s.NoError(err)

	s.Equal(leaf, res)
}

func (s *StateLeafTestSuite) TestUpsertStateLeaf_UpdateAndRetrieve() {
	leaf, err := NewStateLeaf(0, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	err = s.storage.StateTree.upsertStateLeaf(leaf)
	s.NoError(err)

	updatedLeaf, err := NewStateLeaf(0, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(320),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(leaf.StateID, &updatedLeaf.UserState)
	s.NoError(err)

	res, err := s.storage.StateTree.Leaf(leaf.StateID)
	s.NoError(err)

	s.Equal(updatedLeaf, res)
}

func (s *StateLeafTestSuite) TestGetStateLeavesByPublicKey() {
	accounts := []models.AccountLeaf{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{3, 4, 5},
		},
		{
			PubKeyID:  3,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}

	for i := range accounts {
		err := s.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}

	userStates := []models.UserState{
		{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(420),
			Nonce:    models.MakeUint256(0),
		},
		{
			PubKeyID: 2,
			TokenID:  models.MakeUint256(2),
			Balance:  models.MakeUint256(500),
			Nonce:    models.MakeUint256(0),
		},
		{
			PubKeyID: 3,
			TokenID:  models.MakeUint256(25),
			Balance:  models.MakeUint256(1),
			Nonce:    models.MakeUint256(73),
		},
		{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(25),
			Balance:  models.MakeUint256(1),
			Nonce:    models.MakeUint256(73),
		},
	}

	stateLeaves := make([]models.StateLeaf, 0, len(userStates))

	for i := range userStates {
		stateLeaf, err := NewStateLeaf(uint32(i), &userStates[i])
		s.NoError(err)
		stateLeaves = append(stateLeaves, *stateLeaf)

		_, err = s.storage.StateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}

	returnUserStates, err := s.storage.GetStateLeavesByPublicKey(&accounts[0].PublicKey)
	s.NoError(err)

	s.Len(returnUserStates, 3)
	s.Equal(returnUserStates[0], stateLeaves[0])
	s.Equal(returnUserStates[1], stateLeaves[2])
	s.Equal(returnUserStates[2], stateLeaves[3])
}

func (s *StateLeafTestSuite) TestGetStateLeavesByPublicKey_SortsStateLeaves() {
	accounts := []models.AccountLeaf{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}

	for i := range accounts {
		err := s.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}

	stateLeaves := []models.StateLeaf{
		{
			StateID: 2,
			UserState: models.UserState{
				PubKeyID: 1,
				TokenID:  models.MakeUint256(1),
				Balance:  models.MakeUint256(420),
				Nonce:    models.MakeUint256(0),
			},
		},
		{
			StateID: 0,
			UserState: models.UserState{
				PubKeyID: 2,
				TokenID:  models.MakeUint256(2),
				Balance:  models.MakeUint256(500),
				Nonce:    models.MakeUint256(0),
			},
		},
		{
			StateID: 1,
			UserState: models.UserState{
				PubKeyID: 1,
				TokenID:  models.MakeUint256(25),
				Balance:  models.MakeUint256(1),
				Nonce:    models.MakeUint256(73),
			},
		},
	}

	for i := range stateLeaves {
		_, err := s.storage.StateTree.Set(stateLeaves[i].StateID, &stateLeaves[i].UserState)
		s.NoError(err)
	}

	returnUserStates, err := s.storage.GetStateLeavesByPublicKey(&accounts[0].PublicKey)
	s.NoError(err)

	s.Len(returnUserStates, 3)
	s.EqualValues(0, returnUserStates[0].StateID)
	s.EqualValues(1, returnUserStates[1].StateID)
	s.EqualValues(2, returnUserStates[2].StateID)
}

func (s *StateLeafTestSuite) TestGetFeeReceiverStateLeaf() {
	_, err := s.storage.StateTree.Set(0, userState1)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(1, userState2)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(3, userState3)
	s.NoError(err)

	stateLeaf, err := s.storage.GetFeeReceiverStateLeaf(1, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(*userState1, stateLeaf.UserState)
	s.Equal(uint32(0), stateLeaf.StateID)
	s.Equal(uint32(0), s.storage.feeReceiverStateIDs[userState1.TokenID.String()])
}

func (s *StateLeafTestSuite) TestGetFeeReceiverStateLeaf_WorkWithCachedValue() {
	_, err := s.storage.StateTree.Set(0, userState1)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(1, userState2)
	s.NoError(err)

	_, err = s.storage.GetFeeReceiverStateLeaf(userState2.PubKeyID, userState2.TokenID)
	s.NoError(err)
	s.Equal(uint32(1), s.storage.feeReceiverStateIDs[userState2.TokenID.String()])

	stateLeaf, err := s.storage.GetFeeReceiverStateLeaf(userState2.PubKeyID, userState2.TokenID)
	s.NoError(err)
	s.Equal(*userState2, stateLeaf.UserState)
	s.Equal(uint32(1), stateLeaf.StateID)
}

func TestStateLeafTestSuite(t *testing.T) {
	suite.Run(t, new(StateLeafTestSuite))
}
