package postgres

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/jmoiron/sqlx"
)

func RecreateDatabase(cfg *config.PostgresConfig) (err error) {
	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, nil)
	dbInstance, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dbInstance.Close(); err == nil {
			err = closeErr
		}
	}()

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

	_, err = dbInstance.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Name))
	if err != nil {
		return err
	}
	return nil
}
