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
	query := `
		select 
			* 
		from 
			state_leaf 
		natural join (
			select 
				data_hash 
			from (
				select 
					data_hash, token_index, nonce 
				from 
					state_leaf 
				where 
					account_index = $1
			) l1 
			join (
				select 
					token_index, max(nonce) as nonce 
				from 
					state_leaf 
				where 
					account_index = $2 
				group by 
					token_index
			) l2 
			on 
				l1.token_index = l2.token_index 
			and 
				l1.nonce = l2.nonce
		) grouped_leafs`

	res := []models.StateLeaf{}
	err := s.DB.Select(&res, query, accountIndex, accountIndex)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no state leafs found")
	}
	return res, nil
}
