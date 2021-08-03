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

func (s *StoredMerkleTreeTestSuite) TestRoot_InitialRoot() {
	tree := NewStoredMerkleTree("state", s.storage.database)

	root, err := tree.Root()
	s.NoError(err)
	s.Equal(merkletree.GetZeroHash(StateTreeDepth), *root)
}

func (s *StoredMerkleTreeTestSuite) TestRoot_ChangesAfterSet() {
	tree := NewStoredMerkleTree("state", s.storage.database)

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

func (s *StoredMerkleTreeTestSuite) TestTwoTreesWithDifferentNamespaces() {
	stateTree := NewStoredMerkleTree("state", s.storage.database)
	accountTree := NewStoredMerkleTree("account", s.storage.database)

	hash1 := utils.RandomHash()
	_, _, err := stateTree.SetNode(&models.MerklePath{
		Path:  0,
		Depth: StateTreeDepth,
	}, hash1)
	s.NoError(err)

	hash2 := utils.RandomHash()
	_, _, err = accountTree.SetNode(&models.MerklePath{
		Path:  0,
		Depth: StateTreeDepth,
	}, hash2)
	s.NoError(err)

	node1, err := stateTree.Get(models.MerklePath{
		Path:  0,
		Depth: StateTreeDepth,
	})
	s.NoError(err)
	s.Equal(hash1, node1.DataHash)

	node2, err := accountTree.Get(models.MerklePath{
		Path:  0,
		Depth: StateTreeDepth,
	})
	s.NoError(err)
	s.Equal(hash2, node2.DataHash)
}

func TestStoredMerkleTreeTestSuite(t *testing.T) {
	suite.Run(t, new(StoredMerkleTreeTestSuite))
}
