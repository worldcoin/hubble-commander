package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) UpsertStateNode(node *models.StateNode) error {
	_, err := s.Postgres.Query(
		s.QB.Insert("state_node").
			Values(
				node.MerklePath,
				node.DataHash,
			).Suffix("ON CONFLICT (merkle_path) DO UPDATE SET data_hash = ?", node.DataHash),
	).Exec()
	if err != nil {
		return err
	}

	return s.Badger.Upsert(node.MerklePath, node)
}

func (s *Storage) BatchUpsertStateNodes(nodes []models.StateNode) (err error) {
	tx, storage, err := s.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)
	for i := range nodes {
		err = storage.UpsertStateNode(&nodes[i])
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) AddStateNode(node *models.StateNode) error {
	_, err := s.Postgres.Query(
		s.QB.Insert("state_node").
			Values(
				node.MerklePath,
				node.DataHash,
			),
	).Exec()
	if err != nil {
		return err
	}
	return s.Badger.Insert(node.MerklePath, node)
}

func (s *Storage) GetStateNodeByPath(path *models.MerklePath) (*models.StateNode, error) {
	var node models.StateNode
	err := s.Badger.Get(path, &node)
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
		DataHash:   GetZeroHash(leafDepth - uint(path.Depth)),
	}
}

// TODO consider rewriting to badgerhold.Find()
func (s *Storage) getStateNodes(witnessPaths []models.MerklePath) (nodes []models.StateNode, err error) {
	tx, storage, err := s.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	nodes = make([]models.StateNode, 0)
	for i := range witnessPaths {
		var node *models.StateNode
		node, err = storage.GetStateNodeByPath(&witnessPaths[i])
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
