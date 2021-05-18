package storage

import (
	"strings"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	userState1 = &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
	userState2 = &models.UserState{
		PubKeyID:   2,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
)

type StateLeafTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
	tree    *StateTree
}

func (s *StateLeafTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateLeafTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)
	s.tree = NewStateTree(s.storage.Storage)
}

func (s *StateLeafTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StateLeafTestSuite) TestUpsertStateLeaf_AddAndRetrieve() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	leaf := &models.StateLeaf{
		StateID:  0,
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	err = s.storage.UpsertStateLeaf(leaf)
	s.NoError(err)

	res, err := s.storage.GetStateLeaf(leaf.StateID)
	s.NoError(err)

	s.Equal(leaf, res)
}

func (s *StateLeafTestSuite) TestUpsertStateLeaf_UpdateAndRetrieve() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	leaf := &models.StateLeaf{
		StateID:  0,
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	err = s.storage.UpsertStateLeaf(leaf)
	s.NoError(err)

	leaf.UserState.Balance = models.MakeUint256(320)
	err = s.storage.UpsertStateLeaf(leaf)
	s.NoError(err)

	res, err := s.storage.GetStateLeaf(leaf.StateID)
	s.NoError(err)

	s.Equal(leaf, res)
}

func (s *StateLeafTestSuite) TestGetStateLeaf_ReturnsCorrectStruct() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	leaf := &models.StateLeaf{
		StateID:  0,
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	path, err := models.NewMerklePath(strings.Repeat("0", 32))
	s.NoError(err)

	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   leaf.DataHash,
	}

	err = s.storage.UpsertStateLeaf(leaf)
	s.NoError(err)

	err = s.storage.AddStateNode(node)
	s.NoError(err)

	actual, err := s.storage.GetStateLeaf(leaf.StateID)
	s.NoError(err)
	s.Equal(leaf, actual)
}

func (s *StateLeafTestSuite) TestGetStateLeaf_NonExistentLeaf() {
	_, err := s.storage.GetStateLeaf(0)
	s.Equal(NewNotFoundError("state leaf"), err)
}

func (s *StateLeafTestSuite) TestGetNextAvailableStateID_NoLeavesInStateTree() {
	path, err := s.storage.GetNextAvailableStateID()
	s.NoError(err)
	s.Equal(uint32(0), *path)
}

func (s *StateLeafTestSuite) TestGetNextAvailableStateID() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)
	err = s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	err = s.tree.Set(0, userState1)
	s.NoError(err)
	err = s.tree.Set(1, userState2)
	s.NoError(err)

	path, err := s.storage.GetNextAvailableStateID()
	s.NoError(err)
	s.Equal(uint32(2), *path)
}

func (s *StateLeafTestSuite) TestGetUserStatesByPublicKey() {
	accounts := []models.Account{
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
		err := s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	userStates := []models.UserState{
		{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   2,
			TokenIndex: models.MakeUint256(2),
			Balance:    models.MakeUint256(500),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(25),
			Balance:    models.MakeUint256(1),
			Nonce:      models.MakeUint256(73),
		},
		{
			PubKeyID:   3,
			TokenIndex: models.MakeUint256(25),
			Balance:    models.MakeUint256(1),
			Nonce:      models.MakeUint256(73),
		},
	}

	for i := range userStates {
		err := s.tree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}

	returnUserStates, err := s.storage.GetUserStatesByPublicKey(&accounts[0].PublicKey)
	s.NoError(err)

	s.Len(returnUserStates, 3)
	s.Contains(returnUserStates, models.UserStateWithID{
		StateID:   0,
		UserState: userStates[0],
	})
	s.Contains(returnUserStates, models.UserStateWithID{
		StateID:   2,
		UserState: userStates[2],
	})
	s.Contains(returnUserStates, models.UserStateWithID{
		StateID:   3,
		UserState: userStates[3],
	})
}

func (s *StateLeafTestSuite) TestGetFeeReceiverStateLeaf() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)
	err = s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	err = s.tree.Set(0, userState1)
	s.NoError(err)

	err = s.tree.Set(1, userState2)
	s.NoError(err)

	stateLeaf, err := s.storage.GetFeeReceiverStateLeaf(userState1.PubKeyID, userState1.TokenIndex)
	s.NoError(err)
	s.Equal(*userState1, stateLeaf.UserState)
	s.Equal(uint32(0), stateLeaf.StateID)
	s.Equal(uint32(0), s.storage.feeReceiver[userState1.TokenIndex.String()])
}

func (s *StateLeafTestSuite) TestGetUserStateByID() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)
	err = s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	err = s.tree.Set(0, userState1)
	s.NoError(err)
	err = s.tree.Set(1, userState2)
	s.NoError(err)

	userStateWithID, err := s.storage.GetUserStateByID(0)
	s.NoError(err)
	s.Equal(*userState1, userStateWithID.UserState)
	s.Equal(uint32(0), userStateWithID.StateID)
}

func (s *StateLeafTestSuite) TestGetUserStateByID_NonExistentUserState() {
	userStateWithID, err := s.storage.GetUserStateByID(10)
	s.Equal(NewNotFoundError("user state"), err)
	s.Nil(userStateWithID)
}

func TestStateLeafTestSuite(t *testing.T) {
	suite.Run(t, new(StateLeafTestSuite))
}
