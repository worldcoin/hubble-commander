package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StateLeafTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
	tree    *StateTree
}

func (s *StateLeafTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateLeafTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
	s.tree = NewStateTree(s.storage)
}

func (s *StateLeafTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *StateLeafTestSuite) Test_AddStateLeaf_AddAndRetrieve() {
	leaf := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			AccountIndex: 1,
			TokenIndex:   models.MakeUint256(1),
			Balance:      models.MakeUint256(420),
			Nonce:        models.MakeUint256(0),
		},
	}
	err := s.storage.AddStateLeaf(leaf)
	s.NoError(err)

	res, err := s.storage.GetStateLeaf(leaf.DataHash)
	s.NoError(err)

	s.Equal(leaf, res)
}

func (s *StateLeafTestSuite) Test_GetStateLeaf_NonExistentLeaf() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetStateLeaf(hash)
	s.EqualError(err, "state leaf not found")
	s.Nil(res)
}

func (s *StateLeafTestSuite) Test_GetStateLeaves_NoLeaves() {
	res, err := s.storage.GetStateLeaves(1)
	s.EqualError(err, "no state leaves found")
	s.Nil(res)
}

func (s *StateLeafTestSuite) Test_GetStateLeaves() {
	var accountIndex uint32 = 1

	leaves := []models.StateLeaf{
		{
			DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
			UserState: models.UserState{
				AccountIndex: accountIndex,
				TokenIndex:   models.MakeUint256(1),
				Balance:      models.MakeUint256(420),
				Nonce:        models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
			UserState: models.UserState{
				AccountIndex: accountIndex,
				TokenIndex:   models.MakeUint256(2),
				Balance:      models.MakeUint256(500),
				Nonce:        models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{3, 4, 5, 6, 7}),
			UserState: models.UserState{
				AccountIndex: accountIndex,
				TokenIndex:   models.MakeUint256(2),
				Balance:      models.MakeUint256(500),
				Nonce:        models.MakeUint256(1),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{4, 5, 6, 7, 8}),
			UserState: models.UserState{
				AccountIndex: accountIndex,
				TokenIndex:   models.MakeUint256(1),
				Balance:      models.MakeUint256(500),
				Nonce:        models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{5, 6, 7, 8, 9}),
			UserState: models.UserState{
				AccountIndex: accountIndex,
				TokenIndex:   models.MakeUint256(1),
				Balance:      models.MakeUint256(505),
				Nonce:        models.MakeUint256(0),
			},
		},
	}

	for i := range leaves {
		err := s.storage.AddStateLeaf(&leaves[i])
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

	res, err := s.storage.GetStateLeaves(accountIndex)
	s.NoError(err)

	s.Len(res, 2)
	s.Equal(leaves[4].DataHash, res[0].DataHash)
	s.Equal(leaves[2].DataHash, res[1].DataHash)
}

func (s *StateLeafTestSuite) Test_GetUserStatesByPublicKey() {
	accounts := []models.Account{
		{
			AccountIndex: 1,
			PublicKey:    models.PublicKey{1, 2, 3},
		},
		{
			AccountIndex: 3,
			PublicKey:    models.PublicKey{1, 2, 3},
		},
	}

	for i := range accounts {
		err := s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	userStates := []models.UserState{
		{
			AccountIndex: accounts[0].AccountIndex,
			TokenIndex:   models.MakeUint256(1),
			Balance:      models.MakeUint256(420),
			Nonce:        models.MakeUint256(0),
		},
		{
			AccountIndex: 2,
			TokenIndex:   models.MakeUint256(2),
			Balance:      models.MakeUint256(500),
			Nonce:        models.MakeUint256(0),
		},
		{
			AccountIndex: accounts[0].AccountIndex,
			TokenIndex:   models.MakeUint256(25),
			Balance:      models.MakeUint256(1),
			Nonce:        models.MakeUint256(73),
		},
		{
			AccountIndex: accounts[1].AccountIndex,
			TokenIndex:   models.MakeUint256(25),
			Balance:      models.MakeUint256(1),
			Nonce:        models.MakeUint256(73),
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

func TestStateLeafTestSuite(t *testing.T) {
	suite.Run(t, new(StateLeafTestSuite))
}
