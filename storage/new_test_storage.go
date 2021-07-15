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
	var teardown = func() error {
		return nil
	}

	if cfg.Postgres {
		postgresTestDB, err := postgres.NewTestDB()
		if err != nil {
			return nil, err
		}
		storage.Postgres = postgresTestDB.DB
		storage.QB = getQueryBuilder()
		teardown = func() error {
			return postgresTestDB.Teardown()
		}
	}

	if cfg.Badger {
		badgerTestDB, err := badger.NewTestDB()
		if err != nil {
			return nil, err
		}
		storage.Badger = badgerTestDB.DB
		parentTeardown := teardown
		teardown = func() error {
			if err := parentTeardown(); err != nil {
				return err
			}
			return badgerTestDB.Teardown()
		}
	}

	return &TestStorage{
		Storage:  &storage,
		Teardown: teardown,
	}, nil
}

func (s *TestStorage) Clone(cfg *config.CloneConfig) (*TestStorage, error) {
	// TODO-CLONE: rethink it
	// - especially these ugly teardowns
	storage := *s.Storage
	var teardown = func() error {
		return nil
	}
	var oldTeardown = func() error {
		return nil
	}

	if s.Postgres != nil {
		testPostgres := postgres.TestDB{DB: s.Postgres}
		clonedPostgres, err := testPostgres.Clone(&cfg.PostgresConfig, cfg.PostgresSourceDB)
		if err != nil {
			return nil, err
		}
		storage.Postgres = clonedPostgres.DB
		teardown = func() error {
			return clonedPostgres.Teardown()
		}
		oldTeardown = testPostgres.Teardown
	}

	if s.Badger != nil {
		testBadger := badger.TestDB{DB: s.Badger}
		clonedBadger, err := testBadger.Clone()
		if err != nil {
			return nil, err
		}
		storage.Badger = clonedBadger.DB
		parentTeardown := teardown
		teardown = func() error {
			if err := parentTeardown(); err != nil {
				return err
			}
			return clonedBadger.Teardown()
		}
		oldParentTeardown := oldTeardown
		oldTeardown = func() error {
			if err := oldParentTeardown(); err != nil {
				return err
			}
			return s.Badger.Close()
		}
	}

	s.Teardown = oldTeardown

	return &TestStorage{
		Storage:  &storage,
		Teardown: teardown,
	}, nil
}
