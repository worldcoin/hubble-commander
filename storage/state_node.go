package storage

import (
	"reflect"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
)

var flatStateLeafPrefix = []byte("bh_" + reflect.TypeOf(models.FlatStateLeaf{}).Name())

// TODO-ST move to stored_merkle_tree.go
func newZeroStateNode(path *models.MerklePath) *models.MerkleTreeNode {
	return &models.MerkleTreeNode{
		MerklePath: *path,
		DataHash:   merkletree.GetZeroHash(StateTreeDepth - uint(path.Depth)),
	}
}
