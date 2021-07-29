package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

// TODO-acc merge with Leaves
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

// TODO-acc remove and use AccountTree.Leaf instead
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

// TODO-acc move this method
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

// TODO-acc move this method
func (s *Storage) GetPublicKeyByStateID(stateID uint32) (*models.PublicKey, error) {
	stateLeaf, err := s.GetStateLeaf(stateID)
	if IsNotFoundError(err) {
		return nil, NewNotFoundError("account")
	}
	if err != nil {
		return nil, err
	}
	return s.AccountTree.GetPublicKey(stateLeaf.PubKeyID)
}
