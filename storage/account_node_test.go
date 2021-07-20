package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountNodeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *AccountNodeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *AccountNodeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithoutPostgres()
	s.NoError(err)
}

func (s *AccountNodeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *AccountNodeTestSuite) TestUpsertAccountNode_AddAndRetrieve() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.AccountNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
	}
	err = s.storage.UpsertAccountNode(node)
	s.NoError(err)

	res, err := s.storage.GetAccountNodeByPath(path)
	s.NoError(err)

	s.Equal(node, res)
}

func (s *AccountNodeTestSuite) TestUpsertStateNode_UpdateAndRetrieve() {
	path, err := models.NewMerklePath("0000111")
	s.NoError(err)
	node := &models.AccountNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err = s.storage.UpsertAccountNode(node)
	s.NoError(err)

	s.NoError(err)
	expectedNode := &models.AccountNode{
		MerklePath: *path,
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
	}
	err = s.storage.UpsertAccountNode(expectedNode)
	s.NoError(err)

	res, err := s.storage.GetAccountNodeByPath(path)
	s.NoError(err)

	s.Equal(expectedNode, res)
}

func (s *AccountNodeTestSuite) TestGetAccountNodeByPath_NonExistentLeaf() {
	path := models.MerklePath{
		Path:  0,
		Depth: StateTreeDepth,
	}

	expected := &models.AccountNode{
		MerklePath: path,
		DataHash:   merkletree.GetZeroHash(0),
	}

	res, err := s.storage.GetAccountNodeByPath(&path)
	s.NoError(err)
	s.Equal(expected, res)
}

func (s *AccountNodeTestSuite) TestGetAccountNodeByPath_NonExistentRoot() {
	path := models.MerklePath{
		Path:  0,
		Depth: 0,
	}

	expected := &models.AccountNode{
		MerklePath: path,
		DataHash:   merkletree.GetZeroHash(StateTreeDepth),
	}

	res, err := s.storage.GetAccountNodeByPath(&path)
	s.NoError(err)
	s.Equal(expected, res)
}

func TestAccountNodeTestSuite(t *testing.T) {
	suite.Run(t, new(AccountNodeTestSuite))
}
