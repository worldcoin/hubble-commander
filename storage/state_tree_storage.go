package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
)

type StateTreeStorage struct {
	storage Storage
}

func (s *StateTreeStorage) UpsertTreeLeaf(leaf merkletree.TreeLeaf) error {
	stateLeaf := leaf.(*models.StateLeaf)
	return s.storage.UpsertStateLeaf(stateLeaf)
}

func (s *StateTreeStorage) UpsertTreeNode(node merkletree.TreeNode) error {
	stateNode := node.(*models.StateNode)
	return s.storage.UpsertStateNode(stateNode)
}

func (s *StateTreeStorage) GetTreeNode(path *models.MerklePath) (merkletree.TreeNode, error) {
	return s.storage.GetStateNodeByPath(path)
}
