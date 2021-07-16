package postgres

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/golang-migrate/migrate/v4"
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

	migrator, err := GetMigrator(cfg)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil {
		return nil, err
	}

	dbInstance, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	testDB := &TestDB{
		DB:       dbInstance,
		Teardown: newTeardown(dbInstance, migrator),
	}

	return testDB, err
}

func newTeardown(database *Database, migrator *migrate.Migrate) func() error {
	return func() error {
		err := migrator.Down()
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
		return database.Close()
	}
}

// TODO-CLONE do we really need this or can we only use Clone on the actual Database object
func (d *TestDB) Clone(currentConfig *config.PostgresConfig) (testDB *TestDB, err error) {
	clonedDB, err := d.DB.Clone(currentConfig)
	if err != nil {
		return nil, err
	}

	clonedConfig := *currentConfig
	clonedConfig.Name = currentConfig.Name + clonedDBSuffix
	migrator, err := GetMigrator(&clonedConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	d.Teardown = newTeardown(d.DB, migrator)

	return &TestDB{
		DB: clonedDB,
		Teardown: func() error {
			return clonedDB.Close()
		},
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

	return NewDatabase(cfg)
}
