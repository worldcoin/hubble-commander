package storage

import (
	"fmt"

	"github.com/Masterminds/squirrel"
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

func (s *Storage) UpdateStateNode(node *models.StateNode) error {
	res, err := s.Postgres.Query(
		s.QB.Update("state_node").
			Set("data_hash", squirrel.Expr("?", node.DataHash)).
			Where("merkle_path = ?", node.MerklePath),
	).Exec()
	if err != nil {
		return err
	}

	numUpdatedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numUpdatedRows == 0 {
		return fmt.Errorf("no rows were affected by the update")
	}

	return s.Badger.Update(node.MerklePath, node)
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

func (s *Storage) getStateNodes(witnessPaths []models.MerklePath) ([]models.StateNode, error) {
	res := make([]models.StateNode, 0, len(witnessPaths))
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("state_node").
			Where(squirrel.Eq{"merkle_path": witnessPaths}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
