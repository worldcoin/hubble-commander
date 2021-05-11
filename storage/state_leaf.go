package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

var userStateWithIDCols = []string{
	"state_leaf.pub_key_id",
	"state_leaf.token_index",
	"state_leaf.balance",
	"state_leaf.nonce",
	"lpad(merkle_path::text, 33, '0')::bit(33)::bigint AS stateID",
}

func (s *Storage) AddStateLeaf(leaf *models.StateLeaf) error {
	_, err := s.Postgres.Query(
		s.QB.Insert("state_leaf").
			Values(
				leaf.DataHash,
				leaf.PubKeyID,
				leaf.TokenIndex,
				leaf.Balance,
				leaf.Nonce,
			).
			Suffix("ON CONFLICT DO NOTHING"),
	).Exec()
	if err != nil {
		return err
	}

	return s.Badger.Insert(leaf.DataHash, leaf)
}

func (s *Storage) GetStateLeafByHash(hash common.Hash) (*models.StateLeaf, error) {
	var leaf models.StateLeaf
	err := s.Badger.Get(hash, &leaf)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("state leaf")
	}
	if err != nil {
		return nil, err
	}
	return &leaf, nil
}

func (s *Storage) GetStateLeafByPath(path *models.MerklePath) (*models.StateLeaf, error) {
	var node []models.StateNode
	err := s.Badger.Find(&node, bh.Where("MerklePath").Eq(path))
	if err == bh.ErrNotFound || len(node) == 0 {
		return nil, NewNotFoundError("state leaf")
	}
	if err != nil {
		return nil, err
	}

	var leaf []models.StateLeaf
	err = s.Badger.Find(&leaf, bh.Where("DataHash").Eq(node[0].DataHash))
	if err == bh.ErrNotFound || len(leaf) == 0 {
		return nil, NewNotFoundError("state leaf")
	}
	if err != nil {
		return nil, err
	}
	return &leaf[0], nil
}

func (s *Storage) GetStateLeaves(pubKeyID uint32) ([]models.StateLeaf, error) {
	res := make([]models.StateLeaf, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("state_leaf.*").
			From("state_leaf").
			JoinClause("NATURAL JOIN state_node").
			Where(squirrel.Eq{"pub_key_id": pubKeyID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("state leaves")
	}
	return res, nil
}

func (s *Storage) GetNextAvailableStateID() (*uint32, error) {
	res := make([]uint32, 0, 1)
	err := s.Postgres.Query(
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

func (s *Storage) GetUserStatesByPublicKey(publicKey *models.PublicKey) ([]models.UserStateWithID, error) {
	res := make([]models.UserStateWithID, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select(userStateWithIDCols...).
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
	return res, nil
}

func (s *Storage) GetUserStateByPubKeyIDAndTokenIndex(pubKeyID uint32, tokenIndex models.Uint256) (*models.UserStateWithID, error) {
	res := make([]models.UserStateWithID, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select(userStateWithIDCols...).
			From("state_leaf").
			JoinClause("NATURAL JOIN state_node").
			Where(squirrel.Eq{"pub_key_id": pubKeyID}).
			Where(squirrel.Eq{"token_index": tokenIndex}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("state leaf")
	}
	return &res[0], nil
}

func (s *Storage) GetUserStateByID(stateID uint32) (*models.UserStateWithID, error) {
	path := models.MerklePath{
		Path:  stateID,
		Depth: leafDepth,
	}
	res := make([]models.UserStateWithID, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select(userStateWithIDCols...).
			From("state_leaf").
			JoinClause("NATURAL JOIN state_node").
			Where(squirrel.Eq{"merkle_path": path}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("user state")
	}
	return &res[0], nil
}
