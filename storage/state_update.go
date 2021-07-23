package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *StorageBase) AddStateUpdate(update *models.StateUpdate) error {
	return s.Badger.Insert(bh.NextSequence(), *update)
}

func (s *StorageBase) GetStateUpdate(id uint64) (*models.StateUpdate, error) {
	var stateUpdate models.StateUpdate
	err := s.Badger.Get(id, &stateUpdate)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("state update")
	}
	if err != nil {
		return nil, err
	}
	return &stateUpdate, nil
}

func (s *StorageBase) DeleteStateUpdate(id uint64) error {
	return s.Badger.Delete(id, models.StateUpdate{})
}
