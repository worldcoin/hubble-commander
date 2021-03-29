package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StateNodeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *StateNodeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateNodeTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *StateNodeTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *StateNodeTestSuite) Test_AddStateNode_AddAndRetrieve() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err = s.storage.AddStateNode(node)
	s.NoError(err)

	res, err := s.storage.GetStateNodeByHash(node.DataHash)
	s.NoError(err)

	s.Equal(node, res)

	res, err = s.storage.GetStateNodeByPath(path)
	s.NoError(err)

	s.Equal(node, res)
}

func (s *StateNodeTestSuite) Test_AddStateNode_AddAndRetrieveRoot() {
	pathRoot, err := models.NewMerklePath("")
	s.NoError(err)
	pathNode, err := models.NewMerklePath("0")
	s.NoError(err)
	root := &models.StateNode{
		MerklePath: *pathRoot,
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	node := &models.StateNode{
		MerklePath: *pathNode,
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
	}
	err = s.storage.AddStateNode(root)
	s.NoError(err)
	err = s.storage.AddStateNode(node)
	s.NoError(err)

	res, err := s.storage.GetStateNodeByPath(pathRoot)
	s.NoError(err)

	s.Equal(root, res)

	res, err = s.storage.GetStateNodeByPath(pathNode)
	s.NoError(err)

	s.Equal(node, res)
}

func (s *StateNodeTestSuite) Test_UpdateStateNode_UpdateAndRetrieve() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err = s.storage.AddStateNode(node)
	s.NoError(err)

	expectedNode := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
	}

	err = s.storage.UpdateStateNode(expectedNode)
	s.NoError(err)

	res, err := s.storage.GetStateNodeByPath(path)
	s.NoError(err)

	s.Equal(expectedNode, res)
}

func (s *StateNodeTestSuite) Test_UpdateStateNode_NotExistentNode() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
	}

	err = s.storage.UpdateStateNode(node)
	s.Error(err)
}

func (s *StateNodeTestSuite) Test_AddOrUpdateStateNode_AddAndRetrieve() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
	}
	err = s.storage.AddOrUpdateStateNode(node)
	s.NoError(err)

	res, err := s.storage.GetStateNodeByPath(path)
	s.NoError(err)

	s.Equal(node, res)
}

func (s *StateNodeTestSuite) Test_AddOrUpdateStateNode_UpdateAndRetrieve() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err = s.storage.AddStateNode(node)
	s.NoError(err)

	s.NoError(err)
	expectedNode := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
	}
	err = s.storage.AddOrUpdateStateNode(expectedNode)
	s.NoError(err)

	res, err := s.storage.GetStateNodeByPath(path)
	s.NoError(err)

	s.Equal(expectedNode, res)
}

func (s *StateNodeTestSuite) Test_GetStateNodeByHash_NonExistentNode() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetStateNodeByHash(hash)
	s.EqualError(err, "state node not found")
	s.Nil(res)
}

func (s *StateNodeTestSuite) Test_GetStateNodeByPath_NonExistentLeaf() {
	path := models.MerklePath{
		Path:  0,
		Depth: 32,
	}

	expected := &models.StateNode{
		MerklePath: path,
		DataHash:   GetZeroHash(0),
	}

	res, err := s.storage.GetStateNodeByPath(&path)
	s.NoError(err)
	s.Equal(expected, res)
}

func (s *StateNodeTestSuite) Test_GetStateNodeByPath_NonExistentRoot() {
	path := models.MerklePath{
		Path:  0,
		Depth: 0,
	}

	expected := &models.StateNode{
		MerklePath: path,
		DataHash:   GetZeroHash(32),
	}

	res, err := s.storage.GetStateNodeByPath(&path)
	s.NoError(err)
	s.Equal(expected, res)
}

func TestStateNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StateNodeTestSuite))
}
