package storage

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

var ErrAccountAlreadyExists = errors.New("account already exists")

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

func (s *Storage) GetUnusedPubKeyID(publicKey *models.PublicKey, tokenIndex models.Uint256) (*uint32, error) {
	accounts, err := s.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, NewNotFoundError("pub key id")
	}

	userPubKeyIDs := make([]interface{}, 0, len(accounts))
	usedPubKeyIDs := make(map[uint32]bool, len(accounts))
	for i := range accounts {
		usedPubKeyIDs[accounts[i].PubKeyID] = false
		userPubKeyIDs = append(userPubKeyIDs, accounts[i].PubKeyID)
	}

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.Badger.Find(&leaves, bh.Where("PubKeyID").In(userPubKeyIDs...))
	if err != nil {
		return nil, err
	}

	for i := range leaves {
		if leaves[i].TokenIndex.Cmp(&tokenIndex) == 0 {
			usedPubKeyIDs[leaves[i].PubKeyID] = true
			continue
		}
	}

	for pubKeyID, used := range usedPubKeyIDs {
		if !used {
			return &pubKeyID, nil
		}
	}
	return nil, NewNotFoundError("pub key id")
}
