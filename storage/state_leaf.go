package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddStateLeaf(leaf *models.StateLeaf) error {
	_, err := s.DB.Query(
		s.QB.Insert("state_leaf").
			Values(
				leaf.DataHash,
				leaf.PubkeyID,
				leaf.TokenIndex,
				leaf.Balance,
				leaf.Nonce,
			).
			Suffix("ON CONFLICT DO NOTHING"),
	).Exec()

	return err
}

func (s *Storage) GetStateLeafByHash(hash common.Hash) (*models.StateLeaf, error) {
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
		return nil, NewNotFoundError("state leaf")
	}
	return &res[0], nil
}

func (s *Storage) GetStateLeafByPath(path *models.MerklePath) (*models.StateLeaf, error) {
	res := make([]models.StateLeaf, 0, 1)
	err := s.DB.Query(
		s.QB.Select("state_leaf.*").
			From("state_node").
			Join("state_leaf ON state_leaf.data_hash = state_node.data_hash").
			Where(squirrel.Eq{"merkle_path": *path}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrStateLeafNotFound
	}
	return &res[0], nil
}

func (s *Storage) GetStateLeaves(pubkeyID uint32) ([]models.StateLeaf, error) {
	res := make([]models.StateLeaf, 0, 1)
	err := s.DB.Query(
		s.QB.Select("state_leaf.*").
			From("state_leaf").
			JoinClause("NATURAL JOIN state_node").
			Where(squirrel.Eq{"pubkey_id": pubkeyID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("state leaves")
	}
	return res, nil
}

func (s *Storage) GetNextAvailableLeafPath() (*uint32, error) {
	res := make([]uint32, 0, 1)
	err := s.DB.Query(
		s.QB.Select("lpad(merkle_path::text, 33, '0')::bit(33)::bigint + 1 AS next_available_leaf_slot").
			From("state_leaf").
			JoinClause("NATURAL JOIN state_node").
			OrderBy("next_available_leaf_slot DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return ref.Uint32(0), nil
	}

	return &res[0], nil
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
				"state_leaf.pubkey_id",
				"state_leaf.token_index",
				"state_leaf.balance",
				"state_leaf.nonce",
				"state_node.merkle_path",
			).
			From("account").
			JoinClause("NATURAL JOIN state_leaf").
			JoinClause("NATURAL JOIN state_node").
			Where(squirrel.Eq{"account.public_key": publicKey}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("user states")
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
