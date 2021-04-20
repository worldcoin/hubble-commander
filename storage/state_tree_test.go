package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	updatedUserState = models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(800),
		Nonce:      models.MakeUint256(1),
	}
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
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = NewTestStorage(testDB.DB)
	s.tree = NewStateTree(s.storage)

	state := models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
	leaf, err := NewStateLeaf(&state)
	s.NoError(err)
	s.leaf = leaf

	err = s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)
}

func (s *StateTreeTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *StateTreeTestSuite) TestSet_StoresStateLeafRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	actualLeaf, err := s.storage.GetStateLeafByHash(s.leaf.DataHash)
	s.NoError(err)
	s.Equal(s.leaf, actualLeaf)
}

func (s *StateTreeTestSuite) TestSet_StoresLeafStateNodeRecord() {
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

func (s *StateTreeTestSuite) TestSet_UpdatesRootStateNodeRecord() {
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

func (s *StateTreeTestSuite) TestSet_CalculatesCorrectRootForLeafOfId1() {
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

func (s *StateTreeTestSuite) TestSet_CalculatesCorrectRootForTwoLeaves() {
	err := s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	err = s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	state := models.UserState{
		PubKeyID:   2,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
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

func (s *StateTreeTestSuite) TestSet_StoresStateUpdateRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	path := models.MerklePath{
		Path:  0,
		Depth: 32,
	}

	currentRoot := common.HexToHash("0xd8cb702fc833817dccdc3889282af96755b2909274ca2f1a3827a60d11d796eb")
	expectedUpdate := &models.StateUpdate{
		ID:          1,
		StateID:     path,
		CurrentHash: s.leaf.DataHash,
		CurrentRoot: currentRoot,
		PrevHash:    GetZeroHash(0),
		PrevRoot:    GetZeroHash(32),
	}

	update, err := s.storage.GetStateUpdateByRootHash(currentRoot)
	s.NoError(err)

	s.Equal(expectedUpdate, update)
}

func (s *StateTreeTestSuite) TestSet_UpdateExistingLeafCorrectRootStateNode() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	err = s.tree.Set(0, &updatedUserState)
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

func (s *StateTreeTestSuite) TestSet_UpdateExistingLeafCorrectLeafStateNode() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	leaf, err := NewStateLeaf(&updatedUserState)
	s.NoError(err)
	err = s.tree.Set(0, &updatedUserState)
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

func (s *StateTreeTestSuite) TestSet_UpdateExistingLeafNewStateLeafRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	expectedLeaf, err := NewStateLeaf(&updatedUserState)
	s.NoError(err)
	err = s.tree.Set(0, &updatedUserState)
	s.NoError(err)

	leaf, err := s.storage.GetStateLeafByHash(expectedLeaf.DataHash)
	s.NoError(err)
	s.Equal(expectedLeaf, leaf)
}

func (s *StateTreeTestSuite) TestSet_UpdateExistingLeafAddsStateUpdateRecord() {
	err := s.tree.Set(0, &s.leaf.UserState)
	s.NoError(err)

	updatedLeaf, err := NewStateLeaf(&updatedUserState)
	s.NoError(err)
	err = s.tree.Set(0, &updatedUserState)
	s.NoError(err)

	path := models.MerklePath{
		Path:  0,
		Depth: 32,
	}

	currentRoot := common.HexToHash("0x406515786640be8c51eacf1221f017e7f59e04ef59637a27dcb2b2f054b309bf")
	expectedUpdate := &models.StateUpdate{
		ID:          2,
		StateID:     path,
		CurrentHash: updatedLeaf.DataHash,
		CurrentRoot: currentRoot,
		PrevHash:    s.leaf.DataHash,
		PrevRoot:    common.HexToHash("0xd8cb702fc833817dccdc3889282af96755b2909274ca2f1a3827a60d11d796eb"),
	}

	update, err := s.storage.GetStateUpdateByRootHash(currentRoot)
	s.NoError(err)

	s.Equal(expectedUpdate, update)
}

func (s *StateTreeTestSuite) TestRevertTo() {
	err := s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)

	states := []models.UserState{
		{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   2,
			TokenIndex: models.MakeUint256(5),
			Balance:    models.MakeUint256(100),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(500),
			Nonce:      models.MakeUint256(0),
		},
	}

	err = s.tree.Set(0, &states[0])
	s.NoError(err)

	stateRoot, err := s.tree.Root()
	s.NoError(err)

	err = s.tree.Set(1, &states[1])
	s.NoError(err)
	err = s.tree.Set(0, &states[2])
	s.NoError(err)

	err = s.tree.RevertTo(*stateRoot)
	s.NoError(err)

	newStateRoot, err := s.tree.Root()
	s.NoError(err)
	s.Equal(stateRoot, newStateRoot)

	leaf, err := s.tree.Leaf(0)
	s.NoError(err)
	s.Equal(states[0], leaf.UserState)
}

func TestMerkleTreeTestSuite(t *testing.T) {
	suite.Run(t, new(StateTreeTestSuite))
}
