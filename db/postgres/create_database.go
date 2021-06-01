package postgres

import (
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/jmoiron/sqlx"
)

func CreateDatabaseIfNotExist(cfg *config.PostgresConfig) (err error) {
	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, nil)
	dbInstance, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return err
	}
	defer closeDB(dbInstance, &err)

	exists, err := databaseExists(dbInstance, cfg.Name)
	if err != nil {
		return err
	}

	if *exists {
		return nil
	}

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

// nolint:gocritic
func closeDB(dbInstance *sqlx.DB, err *error) {
	if closeErr := dbInstance.Close(); *err == nil {
		*err = closeErr
	}
}

func databaseExists(dbInstance *sqlx.DB, dbName string) (*bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", dbName)
	row := dbInstance.QueryRow(query)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return nil, err
	}
	return &exists, nil
}

func createDatabase(dbInstance *sqlx.DB, dbName string) error {
	_, err := dbInstance.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return err
	}
	return nil
}
