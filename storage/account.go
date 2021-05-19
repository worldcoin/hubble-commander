package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) AddAccountIfNotExists(account *models.Account) error {
	_, err := s.Postgres.Query(
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
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("account").
			Where(squirrel.Eq{"public_key": publicKey}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("accounts")
	}
	return res, nil
}

func (s *Storage) GetPublicKey(pubKeyID uint32) (*models.PublicKey, error) {
	res := make([]models.PublicKey, 0, 1)
	err := s.Postgres.Query(
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

func (s *Storage) GetUnusedPubKeyID(publicKey *models.PublicKey, tokenIndex *models.Uint256) (*uint32, error) {
	accounts, err := s.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}

	for i := range accounts {
		leaves := make([]models.FlatStateLeaf, 0, 1)
		err = s.Badger.Find(
			&leaves,
			bh.Where("TokenIndex").Eq(tokenIndex).Index("TokenIndex").
				And("PubKeyID").Eq(accounts[i].PubKeyID).Index("PubKeyID"),
		)
		if err != nil {
			return nil, err
		}
		if len(leaves) == 0 {
			return &accounts[i].PubKeyID, nil
		}
	}

	return nil, NewNotFoundError("pub key id")
}
