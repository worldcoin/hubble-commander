package postgres

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func CreateDatabaseIfNotExist(cfg *config.PostgresConfig) (err error) {
	exists, err := databaseExists(cfg)
	if err != nil {
		return err
	}
	if *exists {
		return nil
	}

	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, nil)
	dbInstance, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return errors.Errorf("database %s does not exist and could not create it", cfg.Name)
	}
	defer closeDB(dbInstance, &err)

	log.Printf("Creating database %s", cfg.Name)
	return createDatabase(dbInstance, cfg.Name)
}

func RecreateDatabase(cfg *config.PostgresConfig) (err error) {
	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, nil)
	dbInstance, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return err
	}
	defer closeDB(dbInstance, &err)

	query := fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid) 
		FROM pg_stat_activity 
		WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();`,
		cfg.Name,
	)
	_, err = dbInstance.Exec(query)
	if err != nil {
		return err
	}

	_, err = dbInstance.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.Name))
	if err != nil {
		return err
	}

	return createDatabase(dbInstance, cfg.Name)
}

func databaseExists(cfg *config.PostgresConfig) (*bool, error) {
	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, &cfg.Name)
	dbInstance, err := sqlx.Connect("postgres", datasource)
	if err == nil {
		return ref.Bool(true), dbInstance.Close()
	}
	if isDBNotExistsError(cfg.Name, err) {
		return ref.Bool(false), nil
	}
	return nil, err
}

func isDBNotExistsError(dbName string, err error) bool {
	notExistsMsg := fmt.Sprintf("database \"%s\" does not exist", dbName)
	return strings.Contains(err.Error(), notExistsMsg)
}

// nolint:gocritic
func closeDB(dbInstance *sqlx.DB, err *error) {
	if closeErr := dbInstance.Close(); *err == nil {
		*err = closeErr
	}
}

func createDatabase(dbInstance *sqlx.DB, dbName string) error {
	_, err := dbInstance.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return err
	}
	return nil
}
