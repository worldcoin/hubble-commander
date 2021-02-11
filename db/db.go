package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/golang-migrate/migrate/v3"
	"github.com/golang-migrate/migrate/v3/database/postgres"
	_ "github.com/golang-migrate/migrate/v3/database/postgres"
	_ "github.com/golang-migrate/migrate/v3/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func createDatasource(host, port, user, password, dbname *string) string {
	datasource := make([]string, 5)
	datasource = append(datasource, "sslmode=disable")

	if host != nil {
		datasource = append(datasource, fmt.Sprintf("host=%s", *host))
	}
	if port != nil {
		datasource = append(datasource, fmt.Sprintf("port=%s", *port))
	}
	if user != nil {
		datasource = append(datasource, fmt.Sprintf("user=%s", *user))
	}
	if password != nil {
		datasource = append(datasource, fmt.Sprintf("password=%s", *password))
	}
	if dbname != nil {
		datasource = append(datasource, fmt.Sprintf("dbname=%s", *dbname))
	}

	return strings.Join(datasource[:], " ")
}

func GetDB(cfg *config.Config) (*sqlx.DB, error) {
	datasource := createDatasource(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	return sqlx.Connect("postgres", datasource)
}

func GetTestDB(cfg *config.Config) (*sqlx.DB, error) {
	datasource := createDatasource(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, nil)
	db, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", *cfg.DBName)
	_, _ = db.Exec(query) // ignore errors

	query = fmt.Sprintf("CREATE DATABASE %s", *cfg.DBName)
	_, _ = db.Exec(query) // ignore errors

	return GetDB(cfg)
}

func GetMigrator(cfg *config.Config) (*migrate.Migrate, error) {
	datasource := createDatasource(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", datasource)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)
}
