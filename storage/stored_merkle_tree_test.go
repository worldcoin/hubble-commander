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
	storage   *TestStorage
	treeDepth uint8
}

func (s *StoredMerkleTreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.treeDepth = 32
}

func (s *StoredMerkleTreeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *StoredMerkleTreeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StoredMerkleTreeTestSuite) TestRoot_InitialRoot() {
	tree := NewStoredMerkleTree("state", s.storage.database, s.treeDepth)

	root, err := tree.Root()
	s.NoError(err)
	s.Equal(merkletree.GetZeroHash(s.treeDepth), *root)
}

func (s *StoredMerkleTreeTestSuite) TestRoot_ChangesAfterSet() {
	tree := NewStoredMerkleTree("state", s.storage.database, s.treeDepth)

	newRoot, _, err := tree.SetNode(&models.MerklePath{
		Path:  0,
		Depth: s.treeDepth,
	}, utils.RandomHash())
	s.NoError(err)

	root, err := tree.Root()
	s.NoError(err)
	s.NotEqual(merkletree.GetZeroHash(s.treeDepth), *root)
	s.Equal(newRoot, root)
}

func (s *StoredMerkleTreeTestSuite) TestSetSingleNode_VerifiesDepth() {
	tree := NewStoredMerkleTree("state", s.storage.database, s.treeDepth)

	err := tree.SetSingleNode(&models.MerkleTreeNode{
		MerklePath: models.MerklePath{
			Path:  0,
			Depth: s.treeDepth,
		},
		DataHash: utils.RandomHash(),
	})
	s.NoError(err)

	err = tree.SetSingleNode(&models.MerkleTreeNode{
		MerklePath: models.MerklePath{
			Path:  0,
			Depth: 33,
		},
		DataHash: utils.RandomHash(),
	})
	s.ErrorIs(err, ErrExceededTreeDepth)
}

func (s *StoredMerkleTreeTestSuite) TestSetNode_VerifiesDepth() {
	tree := NewStoredMerkleTree("state", s.storage.database, s.treeDepth)

	_, _, err := tree.SetNode(&models.MerklePath{
		Path:  0,
		Depth: s.treeDepth,
	}, utils.RandomHash())
	s.NoError(err)

	_, _, err = tree.SetNode(&models.MerklePath{
		Path:  0,
		Depth: 33,
	}, utils.RandomHash())
	s.ErrorIs(err, ErrExceededTreeDepth)
}

func (s *StoredMerkleTreeTestSuite) TestTwoTreesWithDifferentNamespaces() {
	stateTree := NewStoredMerkleTree("state", s.storage.database, s.treeDepth)
	accountTree := NewStoredMerkleTree("account", s.storage.database, s.treeDepth)

	hash1 := utils.RandomHash()
	_, _, err := stateTree.SetNode(&models.MerklePath{
		Path:  0,
		Depth: s.treeDepth,
	}, hash1)
	s.NoError(err)

	hash2 := utils.RandomHash()
	_, _, err = accountTree.SetNode(&models.MerklePath{
		Path:  0,
		Depth: s.treeDepth,
	}, hash2)
	s.NoError(err)

	node1, err := stateTree.Get(models.MerklePath{
		Path:  0,
		Depth: s.treeDepth,
	})
	s.NoError(err)
	s.Equal(hash1, node1.DataHash)

	node2, err := accountTree.Get(models.MerklePath{
		Path:  0,
		Depth: s.treeDepth,
	})
	s.NoError(err)
	s.Equal(hash2, node2.DataHash)
}

func TestStoredMerkleTreeTestSuite(t *testing.T) {
	suite.Run(t, new(StoredMerkleTreeTestSuite))
}
