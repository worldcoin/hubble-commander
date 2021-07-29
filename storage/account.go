package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *AccountTree) addAccountLeafIfNotExists(account *models.AccountLeaf) error {
	return s.storageBase.Badger.Insert(account.PubKeyID, *account)
}

func (s *AccountTree) getAccountLeaf(pubKeyID uint32) (*models.AccountLeaf, error) {
	var leaf models.AccountLeaf
	err := s.storageBase.Badger.Get(pubKeyID, &leaf)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("account leaf")
	}
	if err != nil {
		return nil, err
	}
	return &leaf, nil
}

func (s *AccountTree) getAccountLeaves(publicKey *models.PublicKey) ([]models.AccountLeaf, error) {
	accounts := make([]models.AccountLeaf, 0, 1)
	err := s.storageBase.Badger.Find(
		&accounts,
		bh.Where("PublicKey").Eq(publicKey).Index("PublicKey"),
	)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, NewNotFoundError("account leaves")
	}
	return accounts, nil
}

func (s *AccountTree) GetPublicKey(pubKeyID uint32) (*models.PublicKey, error) {
	var account models.AccountLeaf
	err := s.storageBase.Badger.Get(pubKeyID, &account)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("account")
	}
	if err != nil {
		return nil, err
	}
	return &account.PublicKey, nil
}

func (s *Storage) GetUnusedPubKeyID(publicKey *models.PublicKey, tokenID *models.Uint256) (*uint32, error) {
	accounts, err := s.AccountTree.Leaves(publicKey)
	if err != nil {
		return nil, err
	}

	for i := range accounts {
		stateLeaves := make([]models.FlatStateLeaf, 0, 1)
		err = s.Badger.Find(
			&stateLeaves,
			bh.Where("TokenID").Eq(tokenID).
				And("PubKeyID").Eq(accounts[i].PubKeyID).Index("PubKeyID"),
		)
		if err != nil {
			return nil, err
		}
		if len(stateLeaves) == 0 {
			return &accounts[i].PubKeyID, nil
		}
	}

	return nil, NewNotFoundError("pub key id")
}

func (s *Storage) GetPublicKeyByStateID(stateID uint32) (*models.PublicKey, error) {
	stateLeaf, err := s.StateTree.Leaf(stateID)
	if IsNotFoundError(err) {
		return nil, NewNotFoundError("account")
	}
	if err != nil {
		return nil, err
	}
	return s.AccountTree.GetPublicKey(stateLeaf.PubKeyID)
}
