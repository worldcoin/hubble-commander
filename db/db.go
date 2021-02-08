package db

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/golang-migrate/migrate/v3"
	_ "github.com/golang-migrate/migrate/v3/database/postgres"
	_ "github.com/golang-migrate/migrate/v3/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDB(cfg *config.Config) (*sqlx.DB, error) {
	datasource := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, err
	}

	return db, err
}

func GetTestDB(cfg *config.Config) (*sqlx.DB, error) {
	datasource := fmt.Sprintf(
		"user=%s password=%s sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
	)

	db, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DBName)
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return GetDB(cfg)
}

func GetMigrator(cfg *config.Config) (*migrate.Migrate, error) {
	datasource := fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	return migrate.New(
		"file://./migrations",
		datasource,
	)
}
