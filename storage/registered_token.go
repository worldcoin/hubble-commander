package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type RegisteredTokenStorage struct {
	database *Database
}

func NewRegisteredTokenStorage(database *Database) *RegisteredTokenStorage {
	return &RegisteredTokenStorage{
		database: database,
	}
}

func (s *RegisteredTokenStorage) copyWithNewDatabase(database *Database) *RegisteredTokenStorage {
	newRegisteredTokenStorage := *s
	newRegisteredTokenStorage.database = database

	return &newRegisteredTokenStorage
}

func (s *RegisteredTokenStorage) AddRegisteredToken(registeredToken *models.RegisteredToken) error {
	return s.database.Badger.Insert(registeredToken.ID, *registeredToken)
}

func (s *RegisteredTokenStorage) GetRegisteredToken(tokenID models.Uint256) (*models.RegisteredToken, error) {
	var registeredToken models.RegisteredToken
	err := s.database.Badger.Get(tokenID, &registeredToken)
	if errors.Is(err, bh.ErrNotFound) {
		return nil, errors.WithStack(NewNotFoundError("registered token"))
	}
	if err != nil {
		return nil, err
	}
	registeredToken.ID = tokenID
	return &registeredToken, nil
}
