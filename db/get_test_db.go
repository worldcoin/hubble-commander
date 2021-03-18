package db

import (
	"github.com/Worldcoin/hubble-commander/config"
)

type TestDB struct {
	DB       *Database
	Teardown func() error
}

func NewTestDB() (*TestDB, error) {
	cfg := config.GetTestConfig()

	err := recreateDatabase(&cfg.DB)
	if err != nil {
		return nil, err
	}

	migrator, err := GetMigrator(&cfg.DB)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil {
		return nil, err
	}

	dbInstance, err := NewDatabase(&cfg.DB)
	if err != nil {
		return nil, err
	}

	teardown := func() error {
		err = migrator.Down()
		if err != nil {
			return err
		}
		err = dbInstance.Close()
		if err != nil {
			return err
		}
		return nil
	}

	testDB := &TestDB{
		DB:       dbInstance,
		Teardown: teardown,
	}

	return testDB, err
}
