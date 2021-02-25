package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddStateNode(node *models.StateNode) error {
	_, err := s.QB.Insert("state_node").
		Values(
			node.MerklePath,
			node.DataHash,
		).
		RunWith(s.DB).
		Exec()

	return err
}

func (s *Storage) GetStateNodeByHash(hash common.Hash) (*models.StateNode, error) {
	res := make([]models.StateNode, 0, 1)
	sql, args, err := s.QB.Select("*").
		From("state_node").
		Where(squirrel.Eq{"data_hash": hash}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = s.DB.Select(&res, sql, args...)
	if err != nil {
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
	sql, args, err := s.QB.Select("*").
		From("state_node").
		Where(squirrel.Eq{"merkle_path": pathValue}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = s.DB.Select(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}
