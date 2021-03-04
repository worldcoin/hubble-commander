package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// Needed for migrator
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type DatabaseLike interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Close() error
	Ping() error
}

type Database struct {
	DatabaseLike
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	datasource := CreateDatasource(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, &cfg.DBName)
	db, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, err
	}
	return &Database{DatabaseLike: db}, nil
}

func GetMigrator(cfg *config.Config) (*migrate.Migrate, error) {
	datasource := CreateDatasource(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, &cfg.DBName)

	db, err := sql.Open("postgres", datasource)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	projectRoot := "../db"

	return migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migrations", projectRoot),
		"postgres",
		driver,
	)
}

func CreateDatasource(host, port, user, password, dbname *string) string {
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

	return strings.Join(datasource, " ")
}
