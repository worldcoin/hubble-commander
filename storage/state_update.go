package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *Storage) AddStateUpdate(update *models.StateUpdate) error {
	return s.Badger.Insert(bh.NextSequence(), update)
}

func (s *Storage) GetStateUpdateByRootHash(stateRootHash common.Hash) (*models.StateUpdate, error) {
	updates := make([]models.StateUpdate, 0, 1)

	err := s.Badger.Find(&updates, bh.Where("CurrentRoot").Eq(stateRootHash).Index("CurrentRoot"))
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 {
		return nil, NewNotFoundError("state update")
	}

	return &updates[0], nil
}

func (s *Storage) DeleteStateUpdate(id uint64) error {
	return s.Badger.Delete(id, &models.StateUpdate{})
}
