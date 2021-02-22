package db

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/jmoiron/sqlx"
)

func recreateDatabase(cfg *config.Config) error {
	datasource := CreateDatasource(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, nil)
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
