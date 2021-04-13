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
}

type Database struct {
	DatabaseLike
}

func NewDatabase(cfg *config.DBConfig) (*Database, error) {
	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, &cfg.Name)
	db, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, err
	}
	return &Database{DatabaseLike: db}, nil
}

func (d *Database) Close() error {
	switch v := d.DatabaseLike.(type) {
	case *sqlx.DB:
		return v.Close()
	default:
		return fmt.Errorf("cannot close Database in transaction mode")
	}
}

func (d *Database) Ping() error {
	switch v := d.DatabaseLike.(type) {
	case *sqlx.DB:
		return v.Ping()
	default:
		return fmt.Errorf("cannot ping Database in transaction mode")
	}
}

func (d *Database) BeginTransaction() (*TransactionController, *Database, error) {
	switch v := d.DatabaseLike.(type) {
	case *sqlx.DB:
		tx, err := v.Beginx()
		if err != nil {
			return nil, nil, err
		}

		db := &Database{DatabaseLike: tx}
		controller := &TransactionController{tx: tx, isLocked: false}
		return controller, db, nil

	case *sqlx.Tx:
		// Already in a transaction
		db := &Database{DatabaseLike: d.DatabaseLike}
		controller := &TransactionController{tx: v, isLocked: true}
		return controller, db, nil
	}
	return nil, nil, fmt.Errorf("database object created with unsupported DatabaseLike implementation")
}

func GetMigrator(cfg *config.DBConfig) (*migrate.Migrate, error) {
	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, &cfg.Name)

	db, err := sql.Open("postgres", datasource)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		"postgres",
		driver,
	)
}

func CreateDatasource(host, port, user, password, dbname *string) string {
	datasource := make([]string, 0, 6)
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
