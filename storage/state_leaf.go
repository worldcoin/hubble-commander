package storage

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddStateLeaf(leaf *models.StateLeaf) error {
	_, err := s.DB.Query(
		s.QB.Insert("state_leaf").
			Values(
				leaf.DataHash,
				leaf.AccountIndex,
				leaf.TokenIndex,
				leaf.Balance,
				leaf.Nonce,
			).
			Suffix("ON CONFLICT DO NOTHING"),
	).Exec()

	return err
}

func (s *Storage) GetStateLeaf(hash common.Hash) (*models.StateLeaf, error) {
	res := make([]models.StateLeaf, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
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

func (s *Storage) GetStateLeaves(accountIndex uint32) ([]models.StateLeaf, error) {
	query := `
	SELECT state_leaf.*
	FROM state_leaf
	NATURAL JOIN state_node
	WHERE account_index = $1`

	res := make([]models.StateLeaf, 0, 1)
	err := s.DB.Select(&res, query, accountIndex)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no state leaves found")
	}
	return res, nil
}

type userStateWithPath struct {
	MerklePath models.MerklePath `db:"merkle_path"`
	models.UserState
}

func (s *Storage) GetUserStatesByPublicKey(publicKey *models.PublicKey) ([]models.UserStateWithID, error) {
	res := make([]userStateWithPath, 0, 1)
	err := s.DB.Query(
		s.QB.
			Select(
				"state_leaf.account_index",
				"state_leaf.token_index",
				"state_leaf.balance",
				"state_leaf.nonce",
				"state_node.merkle_path",
			).
			From("account").
			InnerJoin("state_leaf on state_leaf.account_index = account.account_index").
			InnerJoin("state_node on state_node.data_hash = state_leaf.data_hash").
			Where(squirrel.Eq{"account.public_key": publicKey}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no state leaves found")
	}
	return toUserStateWithID(res), nil
}

func toUserStateWithID(userStates []userStateWithPath) []models.UserStateWithID {
	res := make([]models.UserStateWithID, 0, len(userStates))
	for i := range userStates {
		res = append(res, models.UserStateWithID{
			StateID:   userStates[i].MerklePath.Path,
			UserState: userStates[i].UserState,
		})
	}
	return res
}
