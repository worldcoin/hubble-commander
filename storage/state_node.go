package storage

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) UpsertStateNode(node *models.StateNode) error {
	_, err := s.DB.Query(
		s.QB.Insert("state_node").
			Values(
				node.MerklePath,
				node.DataHash,
			).Suffix("ON CONFLICT (merkle_path) DO UPDATE SET data_hash = ?", node.DataHash),
	).Exec()

	return err
}

func (s *Storage) AddStateNode(node *models.StateNode) error {
	_, err := s.DB.Query(
		s.QB.Insert("state_node").
			Values(
				node.MerklePath,
				node.DataHash,
			),
	).Exec()

	return err
}

func (s *Storage) UpdateStateNode(node *models.StateNode) error {
	res, err := s.DB.Query(
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
	return nil
}

func (s *Storage) GetStateNodeByHash(hash common.Hash) (*models.StateNode, error) {
	res := make([]models.StateNode, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("state_node").
			Where(squirrel.Eq{"data_hash": hash}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("state node")
	}
	return &res[0], nil
}

func (s *Storage) GetStateNodeByPath(path *models.MerklePath) (*models.StateNode, error) {
	res := make([]models.StateNode, 0, 1)
	pathValue, err := path.Value()
	if err != nil {
		return nil, err
	}
	err = s.DB.Query(
		s.QB.Select("*").
			From("state_node").
			Where(squirrel.Eq{"merkle_path": pathValue}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return newZeroStateNode(path), nil
	}

	return &res[0], nil
}

func newZeroStateNode(path *models.MerklePath) *models.StateNode {
	return &models.StateNode{
		MerklePath: *path,
		DataHash:   GetZeroHash(32 - uint(path.Depth)),
	}
}
