package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// Needed for migrator
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DatabaseLike interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Database struct {
	DatabaseLike
}

func NewDatabase(cfg *config.PostgresConfig) (*Database, error) {
	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, &cfg.Name)
	database, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, err
	}
	return &Database{DatabaseLike: database}, nil
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

func (d *Database) BeginTransaction() (*db.TxController, *Database, error) {
	switch v := d.DatabaseLike.(type) {
	case *sqlx.DB:
		tx, err := v.Beginx()
		if err != nil {
			return nil, nil, err
		}

		database := &Database{DatabaseLike: tx}
		controller := db.NewTxController(tx, false)
		return controller, database, nil

	case *sqlx.Tx:
		// Already in a transaction
		database := &Database{DatabaseLike: d.DatabaseLike}
		controller := db.NewTxController(v, true)
		return controller, database, nil
	}
	return nil, nil, fmt.Errorf("database object created with unsupported DatabaseLike implementation")
}

func GetMigrator(cfg *config.PostgresConfig) (*migrate.Migrate, error) {
	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, &cfg.Name)

	database, err := sql.Open("postgres", datasource)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})
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

func (d *Database) Clone(cfg *config.PostgresConfig, templateName string) (*Database, error) {
	err := d.Close()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	datasource := CreateDatasource(cfg.Host, cfg.Port, cfg.User, cfg.Password, nil)
	database, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		err = database.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = disconnectUsers(database, templateName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clonedDB, err := cloneDatabase(database, cfg, templateName)
	if err != nil {
		return nil, err
	}

	oldDatabase, err := NewDatabase(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	*d = *oldDatabase

	return clonedDB, nil
}

func disconnectUsers(database DatabaseLike, dbName string) error {
	_, err := database.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid) 
		FROM pg_stat_activity 
		WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid()`,
		dbName,
	))
	return err
}

func cloneDatabase(database DatabaseLike, cfg *config.PostgresConfig, templateName string) (*Database, error) {
	_, err := database.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.Name))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = database.Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s", cfg.Name, templateName, *cfg.User))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewDatabase(cfg)
}
