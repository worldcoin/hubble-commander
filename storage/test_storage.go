package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
)

type TestStorage struct {
	*Storage
	Database *Database
	Teardown func() error
}

type TeardownFunc func() error

func NewTestStorage() (*TestStorage, error) {
	badgerTestDB, err := db.NewTestDB()
	if err != nil {
		return nil, err
	}

	database := &Database{
		Badger: badgerTestDB.DB,
	}

	return &TestStorage{
		Storage:  newStorageFromDatabase(database),
		Database: database,
		Teardown: badgerTestDB.Teardown,
	}, nil
}

func (s *TestStorage) Clone() (*TestStorage, error) {
	database := *s.database

	testBadger := db.TestDB{DB: s.database.Badger}
	clonedBadger, err := testBadger.Clone()
	if err != nil {
		return nil, err
	}
	database.Badger = clonedBadger.DB

	return &TestStorage{
		Storage:  s.copyWithNewDatabase(&database),
		Teardown: clonedBadger.Teardown,
	}, nil
}
