package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) UpsertAccountNode(node *models.AccountNode) error {
	return s.Badger.Upsert(node.MerklePath, *node)
}

func (s *Storage) GetAccountNodeByPath(path *models.MerklePath) (*models.AccountNode, error) {
	node := models.AccountNode{MerklePath: *path}
	err := s.Badger.Get(*path, &node)
	if err == bh.ErrNotFound {
		return newZeroAccountNode(path), nil
	}
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func newZeroAccountNode(path *models.MerklePath) *models.AccountNode {
	return &models.AccountNode{
		MerklePath: *path,
		DataHash:   merkletree.GetZeroHash(StateTreeDepth - uint(path.Depth)),
	}
}
