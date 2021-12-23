package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type RegisteredSpokeStorage struct {
	database *Database
}

func NewRegisteredSpokeStorage(database *Database) *RegisteredSpokeStorage {
	return &RegisteredSpokeStorage{
		database: database,
	}
}

func (s *RegisteredSpokeStorage) copyWithNewDatabase(database *Database) *RegisteredSpokeStorage {
	newRegisteredSpokeStorage := *s
	newRegisteredSpokeStorage.database = database

	return &newRegisteredSpokeStorage
}

func (s *RegisteredSpokeStorage) AddRegisteredSpoke(registeredSpoke *models.RegisteredSpoke) error {
	return s.database.Badger.Insert(registeredSpoke.ID, *registeredSpoke)
}

func (s *RegisteredSpokeStorage) GetRegisteredSpoke(spokeID models.Uint256) (*models.RegisteredSpoke, error) {
	var registeredSpoke models.RegisteredSpoke
	err := s.database.Badger.Get(spokeID, &registeredSpoke)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("registered spoke"))
	}
	if err != nil {
		return nil, err
	}
	registeredSpoke.ID = spokeID
	return &registeredSpoke, nil
}
