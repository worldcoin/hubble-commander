package testutils

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
)

func GetTestStorage() (*db.Storage, error) {
	cfg := config.GetTestConfig()
	dbInstance, err := GetTestDB(&cfg)
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

	return &db.Storage{DB: dbInstance.DB}, nil
}
