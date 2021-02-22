package db

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/jmoiron/sqlx"
)

type TestDB struct {
	*sqlx.DB
	Teardown func() error
}

func GetTestDB() (*TestDB, error) {
	cfg := config.GetTestConfig()

	err := recreateDatabase(&cfg)
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

	dbInstance, err := GetDB(&cfg)
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
