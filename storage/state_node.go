package storage

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddOrUpdateStateNode(node *models.StateNode) error {
	err := s.UpdateStateNode(node)
	if err != nil {
		isConstraintError := strings.Contains(err.Error(), "no rows were affected by the update")
		if isConstraintError {
			err = s.AddStateNode(node)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (s *Storage) AddStateNode(node *models.StateNode) error {
	_, err := s.DB.Insert(
		s.QB.Insert("state_node").
			Values(
				node.MerklePath,
				node.DataHash,
			),
	)

	return err
}

func (s *Storage) UpdateStateNode(node *models.StateNode) error {
	sql, args, err := s.QB.Update("state_node").
		Set("data_hash", squirrel.Expr("?", node.DataHash)).
		Where("merkle_path = ?", node.MerklePath).ToSql()
	if err != nil {
		return err
	}

	result, err := s.DB.Exec(sql, args...)
	if err != nil {
		return err
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if updatedRows == 0 {
		return fmt.Errorf("no rows were affected by the update")
	}

	return err
}

func (s *Storage) GetStateNodeByHash(hash common.Hash) (*models.StateNode, error) {
	res := make([]models.StateNode, 0, 1)
	err := s.DB.Query(
		squirrel.Select("*").
			From("state_node").
			Where(squirrel.Eq{"data_hash": hash}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
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
		squirrel.Select("*").
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
