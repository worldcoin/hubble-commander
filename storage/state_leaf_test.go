package storage

import (
	"strings"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	s.storage, err = NewTestStorage()
	s.NoError(err)
	s.tree = NewStateTree(s.storage.Storage)
}

func (s *StateLeafTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StateLeafTestSuite) TestAddStateLeaf_AddAndRetrieve() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	leaf := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	err = s.storage.AddStateLeaf(leaf)
	s.NoError(err)

	res, err := s.storage.GetStateLeafByHash(leaf.DataHash)
	s.NoError(err)

	s.Equal(leaf, res)
}

func (s *StateLeafTestSuite) TestGetStateLeafByHash_NonExistentLeaf() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetStateLeafByHash(hash)
	s.Equal(NewNotFoundError("state leaf"), err)
	s.Nil(res)
}

func (s *StateLeafTestSuite) TestGetStateLeafByPath_ReturnsCorrectStruct() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	leaf := &models.StateLeaf{
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

	err = s.storage.AddStateLeaf(leaf)
	s.NoError(err)

	err = s.storage.AddStateNode(node)
	s.NoError(err)

	actual, err := s.storage.GetStateLeafByPath(path)
	s.NoError(err)
	s.Equal(leaf, actual)
}

func (s *StateLeafTestSuite) TestGetStateLeafByPath_NonExistentLeaf() {
	path, err := models.NewMerklePath(strings.Repeat("0", 32))
	s.NoError(err)
	_, err = s.storage.GetStateLeafByPath(path)
	s.Equal(NewNotFoundError("state leaf"), err)
}

func (s *StateLeafTestSuite) TestGetStateLeaves_NoLeaves() {
	res, err := s.storage.GetStateLeaves(1)
	s.Equal(NewNotFoundError("state leaves"), err)
	s.Nil(res)
}

func (s *StateLeafTestSuite) TestGetStateLeaves() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	leaves := []models.StateLeaf{
		{
			DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
			UserState: models.UserState{
				PubKeyID:   account1.PubKeyID,
				TokenIndex: models.MakeUint256(1),
				Balance:    models.MakeUint256(420),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
			UserState: models.UserState{
				PubKeyID:   account1.PubKeyID,
				TokenIndex: models.MakeUint256(2),
				Balance:    models.MakeUint256(500),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{3, 4, 5, 6, 7}),
			UserState: models.UserState{
				PubKeyID:   account1.PubKeyID,
				TokenIndex: models.MakeUint256(2),
				Balance:    models.MakeUint256(500),
				Nonce:      models.MakeUint256(1),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{4, 5, 6, 7, 8}),
			UserState: models.UserState{
				PubKeyID:   account1.PubKeyID,
				TokenIndex: models.MakeUint256(1),
				Balance:    models.MakeUint256(500),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{5, 6, 7, 8, 9}),
			UserState: models.UserState{
				PubKeyID:   account1.PubKeyID,
				TokenIndex: models.MakeUint256(1),
				Balance:    models.MakeUint256(505),
				Nonce:      models.MakeUint256(0),
			},
		},
	}

	for i := range leaves {
		err = s.storage.AddStateLeaf(&leaves[i])
		s.NoError(err)
	}

	path, err := models.NewMerklePath("01")
	s.NoError(err)
	err = s.storage.UpsertStateNode(&models.StateNode{
		DataHash:   common.BytesToHash([]byte{5, 6, 7, 8, 9}),
		MerklePath: *path,
	})
	s.NoError(err)

	path, err = models.NewMerklePath("10")
	s.NoError(err)
	err = s.storage.UpsertStateNode(&models.StateNode{
		DataHash:   common.BytesToHash([]byte{3, 4, 5, 6, 7}),
		MerklePath: *path,
	})
	s.NoError(err)

	res, err := s.storage.GetStateLeaves(account1.PubKeyID)
	s.NoError(err)

	s.Len(res, 2)
	s.Equal(common.BytesToHash([]byte{5, 6, 7, 8, 9}), res[0].DataHash)
	s.Equal(common.BytesToHash([]byte{3, 4, 5, 6, 7}), res[1].DataHash)
}

func (s *StateLeafTestSuite) TestGetNextAvailableStateID_NoLeavesInStateTree() {
	path, err := s.storage.GetNextAvailableStateID()
	s.NoError(err)
	s.Equal(uint32(0), *path)
}

func (s *StateLeafTestSuite) TestGetNextAvailableStateID() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}

	err = s.tree.Set(0, userState)
	s.NoError(err)
	err = s.tree.Set(1, userState)
	s.NoError(err)
	err = s.tree.Set(2, userState)
	s.NoError(err)

	path, err := s.storage.GetNextAvailableStateID()
	s.NoError(err)
	s.Equal(uint32(3), *path)
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

func (s *StateLeafTestSuite) TestGetStateLeafByPubKeyIDAndTokenIndex() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
	err = s.tree.Set(0, userState)
	s.NoError(err)

	userStateWithID, err := s.storage.GetUserStateByPubKeyIDAndTokenIndex(userState.PubKeyID, userState.TokenIndex)
	s.NoError(err)
	s.Equal(*userState, userStateWithID.UserState)
	s.Equal(uint32(0), userStateWithID.StateID)
}

func (s *StateLeafTestSuite) TestGetUserStateByID() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
	err = s.tree.Set(0, userState)
	s.NoError(err)

	userStateWithID, err := s.storage.GetUserStateByID(0)
	s.NoError(err)
	s.Equal(*userState, userStateWithID.UserState)
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
