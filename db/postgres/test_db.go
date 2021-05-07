package postgres

import (
	"github.com/Worldcoin/hubble-commander/config"
)

type TestDB struct {
	DB       *Database
	Teardown func() error
}

func NewTestDB() (*TestDB, error) {
	cfg := config.GetTestConfig().DB

	err := RecreateDatabase(&cfg)
	if err != nil {
		return nil, err
	}

	migrator, err := GetMigrator(&cfg)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil {
		return nil, err
	}

	dbInstance, err := NewDatabase(&cfg)
	if err != nil {
		return nil, err
	}

	teardown := func() error {
		err := migrator.Down() // nolint:govet
		if err != nil {
			return err
		}
		return dbInstance.Close()
	}

	testDB := &TestDB{
		DB:       dbInstance,
		Teardown: teardown,
	}

	return testDB, err
}
