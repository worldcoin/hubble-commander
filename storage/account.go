package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *Storage) GetFirstPubKeyID(publicKey *models.PublicKey) (*uint32, error) {
	var account models.AccountLeaf
	err := s.database.Badger.FindOneUsingIndex(&account, *publicKey, "PublicKey")
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("pub key id"))
	}
	if err != nil {
		return nil, err
	}

	return &account.PubKeyID, nil
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
