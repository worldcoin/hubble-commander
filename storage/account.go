package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *Storage) GetFirstPubKeyID(publicKey *models.PublicKey) (*uint32, error) {
	var accounts []models.AccountLeaf
	// We're not using FindOne here because of inefficient underlying implementation
	err := s.database.Badger.Find(
		&accounts,
		bh.Where("PublicKey").Eq(*publicKey).Index("PublicKey").Limit(1),
	)
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, errors.WithStack(NewNotFoundError("pub key id"))
	}

	return &accounts[0].PubKeyID, nil
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
