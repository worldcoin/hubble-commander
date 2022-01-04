package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *Storage) GetUnusedPubKeyID(publicKey *models.PublicKey, tokenID *models.Uint256) (*uint32, error) {
	accounts, err := s.AccountTree.Leaves(publicKey)
	if err != nil {
		return nil, err
	}

	for i := range accounts {
		stateLeaves := make([]models.FlatStateLeaf, 0, 1)
		err = s.database.Badger.Find(
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
	if err != nil {
		return nil, err
	}
	accountLeaf, err := s.AccountTree.Leaf(stateLeaf.PubKeyID)
	if err != nil {
		return nil, err
	}
	return &accountLeaf.PublicKey, nil
}
