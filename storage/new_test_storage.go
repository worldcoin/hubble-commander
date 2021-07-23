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
	internalStorage := &InternalStorage{feeReceiverStateIDs: make(map[string]uint32)}
	teardown := make([]TeardownFunc, 0, 2)

	if cfg.Postgres {
		postgresTestDB, err := postgres.NewTestDB()
		if err != nil {
			return nil, err
		}
		internalStorage.Postgres = postgresTestDB.DB
		internalStorage.QB = getQueryBuilder()
		teardown = append(teardown, func() error {
			return postgresTestDB.Teardown()
		})
	}

	if cfg.Badger {
		badgerTestDB, err := badger.NewTestDB()
		if err != nil {
			return nil, err
		}
		internalStorage.Badger = badgerTestDB.DB
		teardown = append(teardown, badgerTestDB.Teardown)
	}

	return &TestStorage{
		Storage: &Storage{
			InternalStorage: internalStorage,
			StateTree:       NewStateTree(internalStorage),
			AccountTree:     NewAccountTree(internalStorage),
		},
		Teardown: toTeardownFunc(teardown),
	}, nil
}

func (s *TestStorage) Clone(currentConfig *config.PostgresConfig) (*TestStorage, error) {
	storage := *s.Storage
	teardown := make([]TeardownFunc, 0, 2)
	initialTeardown := make([]TeardownFunc, 0, 2)

	if s.Postgres != nil {
		testPostgres := postgres.TestDB{DB: s.Postgres}
		clonedPostgres, err := testPostgres.Clone(currentConfig)
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

	storage.StateTree = NewStateTree(storage.InternalStorage)
	storage.AccountTree = NewAccountTree(storage.InternalStorage)

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
