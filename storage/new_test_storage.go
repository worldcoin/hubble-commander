package storage

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/db/postgres"
)

// TODO-CLONE NewTestStorage() -> Storage, TeardownFunc
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
	storage := Storage{feeReceiverStateIDs: make(map[string]uint32)}
	teardown := make([]TeardownFunc, 0, 2)

	if cfg.Postgres {
		postgresTestDB, err := postgres.NewTestDB()
		if err != nil {
			return nil, err
		}
		storage.Postgres = postgresTestDB.DB
		storage.QB = getQueryBuilder()
		teardown = append(teardown, func() error {
			return postgresTestDB.Teardown()
		})
	}

	if cfg.Badger {
		badgerTestDB, err := badger.NewTestDB()
		if err != nil {
			return nil, err
		}
		storage.Badger = badgerTestDB.DB
		teardown = append(teardown, badgerTestDB.Teardown)
	}

	return &TestStorage{
		Storage:  &storage,
		Teardown: toTeardownFunc(teardown),
	}, nil
}

func (s *TestStorage) Clone(cfg *config.CloneConfig) (*TestStorage, error) {
	storage := *s.Storage
	teardown := make([]TeardownFunc, 0, 2)
	initialTeardown := make([]TeardownFunc, 0, 2)

	if s.Postgres != nil {
		testPostgres := postgres.TestDB{DB: s.Postgres}
		clonedPostgres, err := testPostgres.Clone(&cfg.PostgresConfig, cfg.PostgresSourceDB)
		if err != nil {
			return nil, err
		}
		storage.Postgres = clonedPostgres.DB
		teardown = append(teardown, func() error {
			return clonedPostgres.Teardown()
		})
		initialTeardown = append(initialTeardown, testPostgres.Teardown)
	}

	if s.Badger != nil {
		testBadger := badger.TestDB{DB: s.Badger}
		clonedBadger, err := testBadger.Clone()
		if err != nil {
			return nil, err
		}
		storage.Badger = clonedBadger.DB
		teardown = append(teardown, func() error {
			return clonedBadger.Teardown()
		})
		initialTeardown = append(initialTeardown, s.Badger.Close)
	}

	s.Teardown = toTeardownFunc(initialTeardown)

	return &TestStorage{
		Storage:  &storage,
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
