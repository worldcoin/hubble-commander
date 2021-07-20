package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StoredMerkleTreeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *StoredMerkleTreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StoredMerkleTreeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithoutPostgres()
	s.NoError(err)
}

func (s *StoredMerkleTreeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StoredMerkleTreeTestSuite) TestInitialRoot() {
	tree := NewStoredMerkleTree("state", s.storage.Storage)

	root, err := tree.Root()
	s.NoError(err)
	s.Equal(merkletree.GetZeroHash(StateTreeDepth), *root)
}

func (s *StoredMerkleTreeTestSuite) TestRootAfterSet() {
	tree := NewStoredMerkleTree("state", s.storage.Storage)

	newRoot, _, err := tree.SetNode(&models.MerklePath{
		Path:  0,
		Depth: StateTreeDepth,
	}, utils.RandomHash())
	s.NoError(err)

	root, err := tree.Root()
	s.NoError(err)
	s.NotEqual(merkletree.GetZeroHash(StateTreeDepth), *root)
	s.Equal(newRoot, root)
}

func TestStoredMerkleTreeTestSuite(t *testing.T) {
	suite.Run(t, new(StoredMerkleTreeTestSuite))
}
