package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddStateLeaf(leaf *models.StateLeaf) error {
	_, err := s.DB.Insert(
		s.QB.Insert("state_leaf").
			Values(
				leaf.DataHash,
				leaf.AccountIndex,
				leaf.TokenIndex,
				leaf.Balance,
				leaf.Nonce,
			),
	)

	return err
}

func (s *Storage) GetStateLeaf(hash common.Hash) (*models.StateLeaf, error) {
	res := make([]models.StateLeaf, 0, 1)
	err := s.DB.Query(
		squirrel.Select("*").
			From("state_leaf").
			Where(squirrel.Eq{"data_hash": hash}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}
