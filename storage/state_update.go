package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) AddStateUpdate(update *models.StateUpdate) error {
	return s.Badger.Insert(bh.NextSequence(), update)
}

func (s *Storage) GetStateUpdate(ID uint64) (*models.StateUpdate, error) {
	var stateUpdate models.StateUpdate
	err := s.Badger.Get(ID, &stateUpdate)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("state update")
	}
	if err != nil {
		return nil, err
	}
	return &stateUpdate, nil
}

func (s *Storage) DeleteStateUpdate(id uint64) error {
	return s.Badger.Delete(id, &models.StateUpdate{})
}
