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
}

func (s *StateLeafTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateLeafTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
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
	s.EqualError(err, "state leaf not found", err.Error())
	s.Nil(res)
}

func (s *StateLeafTestSuite) Test_GetStateLeafs() {
	var accountIndex uint32 = 1

	leafs := []models.StateLeaf{
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
	}

	err := s.storage.AddStateLeaf(&leafs[0])
	s.NoError(err)
	err = s.storage.AddStateLeaf(&leafs[1])
	s.NoError(err)

	res, err := s.storage.GetStateLeafs(accountIndex)
	s.NoError(err)

	s.Equal(leafs, res)
}

func TestStateLeafTestSuite(t *testing.T) {
	suite.Run(t, new(StateLeafTestSuite))
}
