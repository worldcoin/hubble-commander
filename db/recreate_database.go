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

	_, err = dbInstance.Exec(fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();", cfg.DBName))
	if err != nil {
		return err
	}

	_, err = dbInstance.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DBName))
	if err != nil {
		return err
	}

	_, err = dbInstance.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	return nil
}
