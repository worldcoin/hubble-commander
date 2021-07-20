package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
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
	err = s.storage.UpsertStateNode(node)
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
		Depth: StateTreeDepth,
	}

	expected := &models.StateNode{
		MerklePath: path,
		DataHash:   merkletree.GetZeroHash(0),
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
		DataHash:   merkletree.GetZeroHash(StateTreeDepth),
	}

	res, err := s.storage.GetStateNodeByPath(&path)
	s.NoError(err)
	s.Equal(expected, res)
}

func TestStateNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StateNodeTestSuite))
}
