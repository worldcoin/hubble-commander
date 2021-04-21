package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
)

func (s *Storage) AddAccountIfNotExists(account *models.Account) error {
	_, err := s.DB.Query(
		s.QB.Insert("account").
			Values(
				account.PubKeyID,
				account.PublicKey,
			).
			Suffix("ON CONFLICT DO NOTHING"),
	).Exec()

	return err
}

func (s *Storage) GetAccounts(publicKey *models.PublicKey) ([]models.Account, error) {
	res := make([]models.Account, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("account").
			Where(squirrel.Eq{"public_key": publicKey}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetPublicKey(pubKeyID uint32) (*models.PublicKey, error) {
	res := make([]models.PublicKey, 0, 1)
	err := s.DB.Query(
		s.QB.Select("public_key").
			From("account").
			Where(squirrel.Eq{"pub_key_id": pubKeyID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("account")
	}
	return &res[0], nil
}

func (s *Storage) GetUnusedPubKeyID(publicKey *models.PublicKey) (*uint32, error) {
	res := make([]uint32, 0, 1)
	err := s.DB.Query(
		s.QB.Select("account.pub_key_id").
			From("account").
			JoinClause("FULL OUTER JOIN state_leaf ON account.pub_key_id = state_leaf.pub_key_id").
			Where(squirrel.Eq{"public_key": publicKey}).
			Where(squirrel.Eq{"state_leaf.pub_key_id": nil}).
			OrderBy("account.pub_key_id ASC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("account")
	}
	return &res[0], nil
}
