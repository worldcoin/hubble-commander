package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
)

type TemporaryStorage struct {
	*Storage
}

func NewTemporaryStorage() (*TemporaryStorage, error) {
	badgerDB, err := db.NewInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	database := &Database{
		Badger: badgerDB,
	}

	storage := newStorageFromDatabase(database)

	tempStorage := TemporaryStorage{storage}

	return &tempStorage, nil
}

func (s *TemporaryStorage) Close() error {
	return s.database.Badger.Close()
}
