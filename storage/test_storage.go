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

	storage := newStorageFromDatabase(database)

	return &TestStorage{
		Storage:  storage,
		Teardown: toTeardownFunc(teardown),
	}, nil
}

func (s *TestStorage) Clone(currentConfig *config.PostgresConfig) (*TestStorage, error) {
	database := *s.database
	teardown := make([]TeardownFunc, 0, 2)

	if s.database.Postgres != nil {
		testPostgres := postgres.TestDB{DB: s.database.Postgres}
		clonedPostgres, err := testPostgres.Clone(currentConfig)
		if err != nil {
			return nil, err
		}
		database.Postgres = clonedPostgres.DB
		teardown = append(teardown, clonedPostgres.Teardown)
	}

	if s.database.Badger != nil {
		testBadger := badger.TestDB{DB: s.database.Badger}
		clonedBadger, err := testBadger.Clone()
		if err != nil {
			return nil, err
		}
		database.Badger = clonedBadger.DB
		teardown = append(teardown, clonedBadger.Teardown)
	}

	storage := s.copyWithNewDatabase(&database)

	return &TestStorage{
		Storage:  storage,
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
