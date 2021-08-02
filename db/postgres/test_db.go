package postgres

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/pkg/errors"
)

type TestDB struct {
	DB       *Database
	Teardown func() error
}

func NewTestDB() (*TestDB, error) {
	cfg := config.GetTestConfig().Postgres
	err := RecreateDatabase(cfg)
	if err != nil {
		return nil, err
	}

	err = migrateUp(cfg)
	if err != nil {
		return nil, err
	}

	dbInstance, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	testDB := &TestDB{
		DB:       dbInstance,
		Teardown: newTeardown(dbInstance, cfg),
	}

	return testDB, err
}

func newTeardown(database *Database, cfg *config.PostgresConfig) func() error {
	return func() error {
		err := migrateDown(cfg)
		if err != nil {
			return err
		}
		return database.Close()
	}
}

func migrateUp(cfg *config.PostgresConfig) error {
	migrator, err := GetMigrator(cfg)
	if err != nil {
		return err
	}
	err = migrator.Up()
	if err != nil {
		return err
	}
	srcErr, dbErr := migrator.Close()
	if srcErr != nil {
		return srcErr
	}
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func migrateDown(cfg *config.PostgresConfig) error {
	migrator, err := GetMigrator(cfg)
	if err != nil {
		return err
	}
	err = migrator.Down()
	if err != nil {
		return err
	}
	srcErr, dbErr := migrator.Close()
	if srcErr != nil {
		return srcErr
	}
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func (d *TestDB) Clone(currentConfig *config.PostgresConfig) (testDB *TestDB, err error) {
	clonedDB, err := d.DB.Clone(currentConfig)
	if err != nil {
		return nil, err
	}

	clonedConfig := *currentConfig
	clonedConfig.Name = currentConfig.Name + clonedDBSuffix

	return &TestDB{
		DB:       clonedDB,
		Teardown: newTeardown(clonedDB, &clonedConfig),
	}, nil
}

func disconnectUsers(database DatabaseLike, dbName string) error {
	_, err := database.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid) 
		FROM pg_stat_activity 
		WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid()`,
		dbName,
	))
	return err
}

func cloneDatabase(database DatabaseLike, cfg *config.PostgresConfig, clonedDBName string) (*Database, error) {
	_, err := database.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", clonedDBName))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = database.Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s", clonedDBName, cfg.Name, *cfg.User))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clonedCfg := *cfg
	clonedCfg.Name = clonedDBName
	return NewDatabase(&clonedCfg)
}
