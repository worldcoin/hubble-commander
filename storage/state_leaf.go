package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddLeaf(leaf *models.StateLeaf) error {
	_, err := s.QB.Insert("state_leaf").
		Values(
			leaf.DataHash,
			leaf.AccountIndex,
			leaf.TokenIndex,
			leaf.Balance,
			leaf.Nonce,
		).
		RunWith(s.DB).
		Exec()

	return err
}

func (s *Storage) GetLeaf(hash common.Hash) (*models.StateLeaf, error) {
	res := make([]models.StateLeaf, 0, 1)
	sql, args, err := s.QB.Select("*").
		From("state_leaf").
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
