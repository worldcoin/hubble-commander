package storage

import (
	"reflect"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	bh "github.com/timshannon/badgerhold/v3"
)

var flatStateLeafPrefix = []byte("bh_" + reflect.TypeOf(models.FlatStateLeaf{}).Name())

func (s *Storage) UpsertStateNode(node *models.StateNode) error {
	return s.Badger.Upsert(node.MerklePath, *node)
}

func (s *Storage) AddStateNode(node *models.StateNode) error {
	return s.Badger.Insert(node.MerklePath, *node)
}

func (s *Storage) GetStateNodeByPath(path *models.MerklePath) (*models.StateNode, error) {
	node := models.StateNode{MerklePath: *path}
	err := s.Badger.Get(*path, &node)
	if err == bh.ErrNotFound {
		return newZeroStateNode(path), nil
	}
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func newZeroStateNode(path *models.MerklePath) *models.StateNode {
	return &models.StateNode{
		MerklePath: *path,
		DataHash:   merkletree.GetZeroHash(StateTreeDepth - uint(path.Depth)),
	}
}

func (s *Storage) GetStateNodes(paths []models.MerklePath) (nodes []models.StateNode, err error) {
	tx, storage, err := s.BeginTransaction(TxOptions{Badger: true, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	nodes = make([]models.StateNode, 0)
	for i := range paths {
		var node *models.StateNode
		node, err = storage.GetStateNodeByPath(&paths[i])
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, *node)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
