package storage

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/db/postgres"
)

type TestStorage struct {
	*Storage
	Teardown func() error
}

type TestStorageConfig struct {
	Postgres bool
	Badger   bool
}

type TeardownFunc func() error

func NewTestStorage() (*TestStorage, error) {
	return NewConfiguredTestStorage(TestStorageConfig{
		Postgres: true,
		Badger:   false,
	})
}

func NewTestStorageWithBadger() (*TestStorage, error) {
	return NewConfiguredTestStorage(TestStorageConfig{
		Postgres: true,
		Badger:   true,
	})
}

func NewTestStorageWithoutPostgres() (*TestStorage, error) {
	return NewConfiguredTestStorage(TestStorageConfig{
		Postgres: false,
		Badger:   true,
	})
}

func NewConfiguredTestStorage(cfg TestStorageConfig) (*TestStorage, error) {
	database := &Database{QB: getQueryBuilder()}
	teardown := make([]TeardownFunc, 0, 2)

	if cfg.Postgres {
		postgresTestDB, err := postgres.NewTestDB()
		if err != nil {
			return nil, err
		}
		database.Postgres = postgresTestDB.DB
		teardown = append(teardown, postgresTestDB.Teardown)
	}

	if cfg.Badger {
		badgerTestDB, err := badger.NewTestDB()
		if err != nil {
			return nil, err
		}
		database.Badger = badgerTestDB.DB
		teardown = append(teardown, badgerTestDB.Teardown)
	}

	storageBase := &StorageBase{
		Database:            database,
		feeReceiverStateIDs: make(map[string]uint32),
	}

	return &TestStorage{
		Storage: &Storage{
			StorageBase: storageBase,
			StateTree:   NewStateTree(database),
			AccountTree: NewAccountTree(database),
		},
		Teardown: toTeardownFunc(teardown),
	}, nil
}

func (s *TestStorage) Clone(currentConfig *config.PostgresConfig) (*TestStorage, error) {
	storageBase := *s.Storage.StorageBase
	database := *s.Storage.Database
	storageBase.Database = &database

	teardown := make([]TeardownFunc, 0, 2)
	initialTeardown := make([]TeardownFunc, 0, 2)

	if s.Database.Postgres != nil {
		testPostgres := postgres.TestDB{DB: s.Database.Postgres}
		clonedPostgres, err := testPostgres.Clone(currentConfig)
		if err != nil {
			return nil, err
		}
		storageBase.Database.Postgres = clonedPostgres.DB
		teardown = append(teardown, func() error {
			return clonedPostgres.Teardown()
		})
		initialTeardown = append(initialTeardown, testPostgres.Teardown)
	}

	if s.Database.Badger != nil {
		testBadger := badger.TestDB{DB: s.Database.Badger}
		clonedBadger, err := testBadger.Clone()
		if err != nil {
			return nil, err
		}
		storageBase.Database.Badger = clonedBadger.DB
		teardown = append(teardown, func() error {
			return clonedBadger.Teardown()
		})
		initialTeardown = append(initialTeardown, s.Database.Badger.Close)
	}
	s.Teardown = toTeardownFunc(initialTeardown)

	return &TestStorage{
		Storage: &Storage{
			StorageBase: &storageBase,
			StateTree:   NewStateTree(storageBase.Database),
			AccountTree: NewAccountTree(storageBase.Database),
		},
		Teardown: toTeardownFunc(teardown),
	}, nil
}

func toTeardownFunc(teardown []TeardownFunc) TeardownFunc {
	return func() error {
		for i := range teardown {
			if err := teardown[i](); err != nil {
				return err
			}
		}
		return nil
	}
}
