package storage

import (
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
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

func (s *Storage) GetUnusedPubKeyID(publicKey *models.PublicKey) (*uint32, error) {
	accounts, err := s.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, NewNotFoundError("pub key id")
	}

	allPubKeyIDs := make([]string, 0, len(accounts))

	for i := range accounts {
		allPubKeyIDs = append(allPubKeyIDs, strconv.Itoa(int(accounts[i].PubKeyID)))
	}

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.Badger.Find(&leaves, nil)
	if err != nil {
		return nil, err
	}
	if len(leaves) == 0 {
		return nil, NewNotFoundError("pub key id")
	}

	usedPubKeyIDs := make([]string, 0, len(accounts))

	for i := range leaves {
		usedPubKeyIDs = append(usedPubKeyIDs, strconv.Itoa(int(leaves[i].PubKeyID)))
	}

	availablePubKeyIDs := utils.StringSliceDiff(allPubKeyIDs, usedPubKeyIDs)
	if len(availablePubKeyIDs) == 0 {
		return nil, NewNotFoundError("pub key id")
	}

	firstAvailablePubKeyIDUint64, err := strconv.ParseUint(availablePubKeyIDs[0], 10, 32)
	firstAvailablePubKeyID := uint32(firstAvailablePubKeyIDUint64)

	return &firstAvailablePubKeyID, nil
}
