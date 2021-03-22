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
		SELECT
			*
		FROM
			state_leaf
			NATURAL JOIN (
				SELECT
					data_hash
				FROM (
					SELECT
						data_hash,
						token_index,
						nonce
					FROM
						state_leaf
					WHERE
						account_index = $1) l1
					JOIN (
						SELECT
							token_index, max(nonce) AS nonce
						FROM
							state_leaf
						WHERE
							account_index = $2
						GROUP BY
							token_index) l2 ON l1.token_index = l2.token_index
						AND l1.nonce = l2.nonce) grouped_leafs`

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
