package merkletree

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type TreeLeaf interface {
	Index() uint32
}

type TreeNode interface {
	Path() *models.MerklePath
	Hash() *common.Hash
}

type MerkleTreeStorage interface {
	UpsertTreeLeaf(leaf TreeLeaf) error
	UpsertTreeNode(node TreeNode) error
	GetTreeNode(path *models.MerklePath) (TreeNode, error)
}

type StoredMerkleTree struct {
	storage MerkleTreeStorage
}

func (t *StoredMerkleTree) Set(index uint32, leaf TreeLeaf) (*common.Hash, models.Witness, error) {
	// TODO possibly start a DB tx (just in case caller hasn't done it already)
	panic("not implemented")
}
