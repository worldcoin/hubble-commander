package storage

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddStateLeaf(leaf *models.StateLeaf) error {
	_, err := s.DB.ExecBuilder(
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
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("state leaf not found")
	}
	return &res[0], nil
}

func (s *Storage) GetStateLeafs(accountIndex uint32) ([]models.StateLeaf, error) {
	res := make([]models.StateLeaf, 0, 1)
	err := s.DB.Query(
		squirrel.Select("*").
			From("state_leaf").
			Where(squirrel.Eq{"account_index": accountIndex}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return res, nil
}
