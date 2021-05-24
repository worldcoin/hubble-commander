package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StateNodeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *StateNodeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateNodeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithoutPostgres()
	s.NoError(err)
}

func (s *StateNodeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StateNodeTestSuite) TestAddStateNode_AddAndRetrieve() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err = s.storage.AddStateNode(node)
	s.NoError(err)

	res, err := s.storage.GetStateNodeByPath(path)
	s.NoError(err)
	s.Equal(node, res)
}

func (s *StateNodeTestSuite) TestAddStateNode_AddAndRetrieveRoot() {
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

func (s *StateNodeTestSuite) TestUpsertStateNode_AddAndRetrieve() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
	}
	err = s.storage.UpsertStateNode(node)
	s.NoError(err)

	res, err := s.storage.GetStateNodeByPath(path)
	s.NoError(err)

	s.Equal(node, res)
}

func (s *StateNodeTestSuite) TestUpsertStateNode_UpdateAndRetrieve() {
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
	err = s.storage.UpsertStateNode(expectedNode)
	s.NoError(err)

	res, err := s.storage.GetStateNodeByPath(path)
	s.NoError(err)

	s.Equal(expectedNode, res)
}

func (s *StateNodeTestSuite) TestGetStateNodeByPath_NonExistentLeaf() {
	path := models.MerklePath{
		Path:  0,
		Depth: leafDepth,
	}

	expected := &models.StateNode{
		MerklePath: path,
		DataHash:   GetZeroHash(0),
	}

	res, err := s.storage.GetStateNodeByPath(&path)
	s.NoError(err)
	s.Equal(expected, res)
}

func (s *StateNodeTestSuite) TestGetStateNodeByPath_NonExistentRoot() {
	path := models.MerklePath{
		Path:  0,
		Depth: 0,
	}

	expected := &models.StateNode{
		MerklePath: path,
		DataHash:   GetZeroHash(leafDepth),
	}

	res, err := s.storage.GetStateNodeByPath(&path)
	s.NoError(err)
	s.Equal(expected, res)
}

func (s *StateNodeTestSuite) TestGetStateNodes() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)

	node := &models.StateNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err = s.storage.AddStateNode(node)
	s.NoError(err)

	nodes, err := s.storage.GetStateNodes([]models.MerklePath{*path})
	s.NoError(err)
	s.Len(nodes, 1)
}

func (s *StateNodeTestSuite) TestBatchUpsertStateNode_AddAndRetrieve() {
	paths, nodes := getPathsAndNodes()
	err := s.storage.BatchUpsertStateNodes(nodes)
	s.NoError(err)

	res, err := s.storage.GetStateNodes(paths)
	s.NoError(err)
	s.Len(res, 2)
	s.Contains(res, nodes[0])
	s.Contains(res, nodes[1])
}

func (s *StateNodeTestSuite) TestBatchUpsertStateNode_UpdateAndRetrieve() {
	paths, nodes := getPathsAndNodes()
	node := models.StateNode{
		MerklePath: paths[0],
		DataHash:   common.BytesToHash([]byte{8, 7, 6, 5, 4}),
	}
	err := s.storage.AddStateNode(&node)
	s.NoError(err)

	err = s.storage.BatchUpsertStateNodes(nodes)
	s.NoError(err)

	res, err := s.storage.GetStateNodes(paths)
	s.NoError(err)
	s.Len(res, 2)
	s.Contains(res, nodes[0])
	s.Contains(res, nodes[1])
}

func TestStateNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StateNodeTestSuite))
}

func getPathsAndNodes() ([]models.MerklePath, []models.StateNode) {
	paths := []models.MerklePath{
		{
			Path:  7,
			Depth: 7,
		},
		{
			Path:  6,
			Depth: 7,
		},
	}
	nodes := []models.StateNode{
		{
			MerklePath: paths[0],
			DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
		},
		{
			MerklePath: paths[1],
			DataHash:   common.BytesToHash([]byte{1, 2, 3, 5, 6}),
		},
	}
	return paths, nodes
}
