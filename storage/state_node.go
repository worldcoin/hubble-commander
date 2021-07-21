package storage

import (
	"reflect"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	bh "github.com/timshannon/badgerhold/v3"
)

var flatStateLeafPrefix = []byte("bh_" + reflect.TypeOf(models.FlatStateLeaf{}).Name())

func (s *Storage) UpsertStateNode(node *models.MerkleTreeNode) error {
	return s.Badger.Upsert(node.MerklePath, *node)
}

func (s *Storage) GetStateNodeByPath(path *models.MerklePath) (*models.MerkleTreeNode, error) {
	node := models.MerkleTreeNode{MerklePath: *path}
	err := s.Badger.Get(*path, &node)
	if err == bh.ErrNotFound {
		return newZeroStateNode(path), nil
	}
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func newZeroStateNode(path *models.MerklePath) *models.MerkleTreeNode {
	return &models.MerkleTreeNode{
		MerklePath: *path,
		DataHash:   merkletree.GetZeroHash(StateTreeDepth - uint(path.Depth)),
	}
}
