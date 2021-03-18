package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
)

func (s *Storage) AddAccount(account *models.Account) error {
	_, err := s.DB.ExecBuilder(
		s.QB.Insert("account").
			Values(
				account.AccountIndex,
				account.PublicKey,
			).
			Suffix("ON CONFLICT DO NOTHING"),
	)

	return err
}

func (s *Storage) GetAccounts(publicKey *models.PublicKey) ([]models.Account, error) {
	res := make([]models.Account, 0, 1)
	err := s.DB.Query(
		squirrel.Select("*").
			From("account").
			Where(squirrel.Eq{"public_key": publicKey}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
