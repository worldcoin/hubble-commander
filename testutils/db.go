package testutils

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/jmoiron/sqlx"
)

type TestDB struct {
	*sqlx.DB
	Teardown func() error
}

func GetTestDB() (*TestDB, error) {
	cfg := config.GetTestConfig()

	err := RecreateDatabase(&cfg)
	if err != nil {
		return nil, err
	}

	migrator, err := db.GetMigrator(&cfg)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil {
		return nil, err
	}

	dbInstance, err := db.GetDB(&cfg)
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

func RecreateDatabase(cfg *config.Config) error {
	datasource := db.CreateDatasource(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, nil)
	dbInstance, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DBName)
	_, err = dbInstance.Exec(query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
	_, err = dbInstance.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
