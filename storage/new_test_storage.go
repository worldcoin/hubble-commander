package storage

import (
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

func NewConfiguredTestStorage(cfg TestStorageConfig) (*TestStorage, error) {
	storage := Storage{feeReceiver: make(map[string]uint32)}
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
