package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *StateTree) addStateUpdate(update *models.StateUpdate) error {
	return s.storageBase.Badger.Insert(bh.NextSequence(), *update)
}

func (s *StateTree) getStateUpdate(id uint64) (*models.StateUpdate, error) {
	var stateUpdate models.StateUpdate
	err := s.storageBase.Badger.Get(id, &stateUpdate)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("state update")
	}
	if err != nil {
		return nil, err
	}
	return &stateUpdate, nil
}

func (s *StateTree) deleteStateUpdate(id uint64) error {
	return s.storageBase.Badger.Delete(id, models.StateUpdate{})
}
