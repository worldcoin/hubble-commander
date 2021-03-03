package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
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

func (s *StateTreeTestSuite) Test_Set_StoresLeafStateNodeRecord() {
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

func (s *StateTreeTestSuite) Test_Set_UpdatesRootStateNodeRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	rootPath := models.MerklePath{
		Path:  0,
		Depth: 0,
	}

	expectedRoot := &models.StateNode{
		MerklePath: rootPath,
		DataHash:   common.HexToHash("0xd8cb702fc833817dccdc3889282af96755b2909274ca2f1a3827a60d11d796eb"),
	}

	root, err := s.storage.GetStateNodeByPath(&rootPath)
	s.NoError(err)
	s.Equal(expectedRoot, root)
}

func (s *StateTreeTestSuite) Test_Set_CalculatesCorrectRootForLeafOfIndex1() {
	err := s.tree.Set(1, &s.leaf.UserState)
	s.NoError(err)

	rootPath := models.MerklePath{
		Path:  0,
		Depth: 0,
	}

	expectedRoot := &models.StateNode{
		MerklePath: rootPath,
		DataHash:   common.HexToHash("0xbec68099063e1499a5144a2d5b41f6a3e005ceac77caef6a171d77573570a000"),
	}

	root, err := s.storage.GetStateNodeByPath(&rootPath)
	s.NoError(err)
	s.Equal(expectedRoot, root)
}

func (s *StateTreeTestSuite) Test_Set_CalculatesCorrectRootForTwoLeaves() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	state := models.UserState{
		AccountIndex: models.MakeUint256(2),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	err = s.tree.Set(1, &state)
	s.NoError(err)

	rootPath := models.MerklePath{
		Path:  0,
		Depth: 0,
	}

	expectedRoot := &models.StateNode{
		MerklePath: rootPath,
		DataHash:   common.HexToHash("0x7b1b0382bdffda7f4a6b24d974189c60797b87ce76836de6f18039e1dc73c050"),
	}

	root, err := s.storage.GetStateNodeByPath(&rootPath)
	s.NoError(err)
	s.Equal(expectedRoot, root)
}

func (s *StateTreeTestSuite) Test_Set_StoresStateUpdateRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	path := models.MerklePath{
		Path:  0,
		Depth: 32,
	}

	expectedUpdate := &models.StateUpdate{
		ID:          1,
		MerklePath:  path,
		CurrentHash: s.leaf.DataHash,
		CurrentRoot: common.HexToHash("0xd8cb702fc833817dccdc3889282af96755b2909274ca2f1a3827a60d11d796eb"),
		PrevHash:    GetZeroHash(0),
		PrevRoot:    GetZeroHash(32),
	}

	update, err := s.storage.GetStateUpdate(1)
	s.NoError(err)

	s.Equal(expectedUpdate, update)
}

func (s *StateTreeTestSuite) Test_Set_UpdateExistingLeaf_CorrectRootStateNode() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	state := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(800),
		Nonce:        models.MakeUint256(1),
	}
	err = s.tree.Set(0, &state)
	s.NoError(err)

	rootPath := models.MerklePath{
		Path:  0,
		Depth: 0,
	}

	expectedRoot := &models.StateNode{
		MerklePath: rootPath,
		DataHash:   common.HexToHash("0x406515786640be8c51eacf1221f017e7f59e04ef59637a27dcb2b2f054b309bf"),
	}

	root, err := s.storage.GetStateNodeByPath(&rootPath)
	s.NoError(err)
	s.Equal(expectedRoot, root)
}

func (s *StateTreeTestSuite) Test_Set_UpdateExistingLeaf_CorrectLeafStateNode() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	state := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(800),
		Nonce:        models.MakeUint256(1),
	}
	leaf, err := NewStateLeaf(&state)
	s.NoError(err)
	err = s.tree.Set(0, &state)
	s.NoError(err)

	leafPath := models.MerklePath{
		Path:  0,
		Depth: 32,
	}

	expectedLeaf := &models.StateNode{
		MerklePath: leafPath,
		DataHash:   leaf.DataHash,
	}

	leafNode, err := s.storage.GetStateNodeByPath(&leafPath)
	s.NoError(err)
	s.Equal(expectedLeaf, leafNode)
}

func (s *StateTreeTestSuite) Test_Set_UpdateExistingLeaf_NewStateLeafRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	state := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(800),
		Nonce:        models.MakeUint256(1),
	}
	expectedLeaf, err := NewStateLeaf(&state)
	s.NoError(err)
	err = s.tree.Set(0, &state)
	s.NoError(err)

	leaf, err := s.storage.GetStateLeaf(expectedLeaf.DataHash)
	s.NoError(err)
	s.Equal(expectedLeaf, leaf)
}

func (s *StateTreeTestSuite) Test_Set_UpdateExistingLeaf_AddsStateUpdateRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	state := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(800),
		Nonce:        models.MakeUint256(1),
	}
	updatedLeaf, err := NewStateLeaf(&state)
	s.NoError(err)
	err = s.tree.Set(0, &state)
	s.NoError(err)

	path := models.MerklePath{
		Path:  0,
		Depth: 32,
	}

	expectedUpdate := &models.StateUpdate{
		ID:          2,
		MerklePath:  path,
		CurrentHash: updatedLeaf.DataHash,
		CurrentRoot: common.HexToHash("0x406515786640be8c51eacf1221f017e7f59e04ef59637a27dcb2b2f054b309bf"),
		PrevHash:    s.leaf.DataHash,
		PrevRoot:    common.HexToHash("0xd8cb702fc833817dccdc3889282af96755b2909274ca2f1a3827a60d11d796eb"),
	}

	update, err := s.storage.GetStateUpdate(2)
	s.NoError(err)

	s.Equal(expectedUpdate, update)
}

func TestMerkleTreeTestSuite(t *testing.T) {
	suite.Run(t, new(StateTreeTestSuite))
}
