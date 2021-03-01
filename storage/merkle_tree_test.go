package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StateTreeTestSuite struct {
	*require.Assertions
	suite.Suite
	db      *db.TestDB
	storage *Storage
	tree    *StateTree
	leaf    *models.StateLeaf
}

func (s *StateTreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateTreeTestSuite) SetupTest() {
	testDB, err := db.GetTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = NewTestStorage(testDB.DB)
	s.tree = NewStateTree(s.storage)

	state := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	leaf, err := NewStateLeaf(&state)
	s.NoError(err)
	s.leaf = leaf
}

func (s *StateTreeTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *StateTreeTestSuite) Test_Set_StoresStateLeafRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	actualLeaf, err := s.storage.GetStateLeaf(s.leaf.DataHash)
	s.NoError(err)
	s.Equal(s.leaf, actualLeaf)
}

func (s *StateTreeTestSuite) Test_Set_StoresStateNodeRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	expectedNode := &models.StateNode{
		MerklePath: models.MerklePath{
			Path:  0,
			Depth: 32,
		},
		DataHash: s.leaf.DataHash,
	}

	node, err := s.storage.GetStateNodeByHash(s.leaf.DataHash)
	s.NoError(err)
	s.Equal(expectedNode, node)
}

func TestMerkleTreeTestSuite(t *testing.T) {
	suite.Run(t, new(StateTreeTestSuite))
}
