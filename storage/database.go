package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/golang-migrate/migrate/v4"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	Postgres *postgres.Database
	Badger   *badger.Database
	QB       squirrel.StatementBuilderType
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	err := postgres.CreateDatabaseIfNotExist(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	migrator, err := postgres.GetMigrator(cfg.Postgres)
	if err != nil {
		return nil, err
	}
	defer func() {
		srcErr, dbErr := migrator.Close()
		if err == nil {
			if srcErr != nil {
				err = srcErr
			} else {
				err = dbErr
			}
		}
	}()

	postgresDB, err := postgres.NewDatabase(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	badgerDB, err := badger.NewDatabase(cfg.Badger)
	if err != nil {
		return nil, err
	}

	database := &Database{
		Postgres: postgresDB,
		Badger:   badgerDB,
		QB:       getQueryBuilder(),
	}

	if cfg.Bootstrap.Prune {
		err = database.Prune(migrator)
		if err != nil {
			return nil, err
		}
		log.Debug("Badger and Postgres databases were pruned")
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return database, nil
}

func (s *Database) BeginTransaction(opts TxOptions) (*db.TxController, *Database, error) {
	var txController *db.TxController
	storage := *s

	if opts.Postgres && !opts.ReadOnly {
		postgresTx, postgresDB, err := s.Postgres.BeginTransaction()
		if err != nil {
			return nil, nil, err
		}
		txController = postgresTx
		storage.Postgres = postgresDB
	}

	if opts.Badger {
		badgerTx, badgerDB := s.Badger.BeginTransaction(!opts.ReadOnly)
		if txController != nil {
			combinedController := NewCombinedController(txController, badgerTx)
			txController = db.NewTxController(combinedController, txController.IsLocked())
		} else {
			txController = badgerTx
		}
		storage.Badger = badgerDB
	}

	return txController, &storage, nil
}

func (s *Database) Close() error {
	err := s.Postgres.Close()
	if err != nil {
		return err
	}
	return s.Badger.Close()
}

func (s *Database) Prune(migrator *migrate.Migrate) error {
	err := migrator.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return s.Badger.Prune()
}

func getQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
