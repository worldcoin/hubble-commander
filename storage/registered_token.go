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
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("registered token"))
	}
	if err != nil {
		return nil, err
	}
	registeredToken.ID = tokenID
	return &registeredToken, nil
}

func (s *RegisteredTokenStorage) DeleteRegisteredTokens(tokenIds ...models.Uint256) (err error) {
	tx, txDatabase := s.database.BeginTransaction(TxOptions{})
	defer tx.Rollback(&err)

	registeredToken := models.RegisteredToken{}
	for i := range tokenIds {
		err = txDatabase.Badger.Delete(tokenIds[i], registeredToken)
		if err == bh.ErrNotFound {
			return errors.WithStack(NewNotFoundError("registered token"))
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
