package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) AddAccountIfNotExists(account *models.Account) error {
	return s.Badger.Upsert(account.PubKeyID, account)
}

func (s *Storage) GetAccounts(publicKey *models.PublicKey) ([]models.Account, error) {
	accounts := make([]models.Account, 0, 1)
	err := s.Badger.Find(
		&accounts,
		bh.Where("PublicKey").Eq(publicKey).Index("PublicKey"),
	)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, NewNotFoundError("accounts")
	}
	return accounts, nil
}

func (s *Storage) GetPublicKey(pubKeyID uint32) (*models.PublicKey, error) {
	var account models.Account
	err := s.Badger.Get(pubKeyID, &account)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("account")
	}
	if err != nil {
		return nil, err
	}
	return &account.PublicKey, nil
}

func (s *Storage) GetUnusedPubKeyID(publicKey *models.PublicKey, tokenID *models.Uint256) (*uint32, error) {
	accounts, err := s.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}

	for i := range accounts {
		leaves := make([]models.FlatStateLeaf, 0, 1)
		err = s.Badger.Find(
			&leaves,
			bh.Where("TokenID").Eq(tokenID).
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

func (s *Storage) GetPublicKeyByStateID(stateID uint32) (*models.PublicKey, error) {
	stateLeaf, err := s.GetStateLeaf(stateID)
	if IsNotFoundError(err) {
		return nil, NewNotFoundError("account")
	}
	if err != nil {
		return nil, err
	}
	return s.GetPublicKey(stateLeaf.PubKeyID)
}
