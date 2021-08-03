package postgres

import (
	"fmt"
	"strings"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func CreateDatabaseIfNotExist(cfg *config.PostgresConfig) (err error) {
	exists, err := databaseExists(cfg)
	if err != nil {
		return err
	}
	if *exists {
		return nil
	}
	log.Debug("Postgres database does not exist, attempting to create it")

	dbInstance, err := connect(cfg.Host, cfg.Port, cfg.User, cfg.Password, nil)
	if err != nil {
		return errors.Errorf("database %s does not exist and could not create it", cfg.Name)
	}
	defer closeDB(dbInstance, &err)

	log.Printf("Creating database %s", cfg.Name)
	return createDatabase(dbInstance, cfg.Name)
}

func RecreateDatabase(cfg *config.PostgresConfig) (err error) {
	dbInstance, err := connect(cfg.Host, cfg.Port, cfg.User, cfg.Password, nil)
	if err != nil {
		return err
	}
	defer closeDB(dbInstance, &err)

	err = disconnectUsers(dbInstance, cfg.Name)
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
	dbInstance, err := connect(cfg.Host, cfg.Port, cfg.User, cfg.Password, &cfg.Name)
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
	return err
}
