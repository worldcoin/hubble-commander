package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/golang-migrate/migrate/v4"
)

type Storage struct {
	Postgres            *postgres.Database
	Badger              *badger.Database
	QB                  squirrel.StatementBuilderType
	domain              *bls.Domain
	feeReceiverStateIDs map[string]uint32 // token index => state id
	latestBlockNumber   uint32
	syncedBlock         *uint64
}

type TxOptions struct {
	Postgres bool
	Badger   bool
	ReadOnly bool
}

func NewStorage(postgresConfig *config.PostgresConfig, badgerConfig *config.BadgerConfig) (*Storage, error) {
	postgresDB, err := postgres.NewDatabase(postgresConfig)
	if err != nil {
		return nil, err
	}

	badgerDB, err := badger.NewDatabase(badgerConfig)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Postgres:            postgresDB,
		Badger:              badgerDB,
		QB:                  getQueryBuilder(),
		feeReceiverStateIDs: make(map[string]uint32),
	}, nil
}

func NewConfiguredStorage(cfg *config.Config) (storage *Storage, err error) {
	err = postgres.CreateDatabaseIfNotExist(cfg.Postgres)
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

	storage, err = NewStorage(cfg.Postgres, cfg.Badger)
	if err != nil {
		return nil, err
	}

	if cfg.Rollup.Prune {
		err = storage.Prune(migrator)
		if err != nil {
			return nil, err
		}
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) BeginTransaction(opts TxOptions) (*db.TxController, *Storage, error) {
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

func (s *Storage) Close() error {
	err := s.Postgres.Close()
	if err != nil {
		return err
	}
	return s.Badger.Close()
}

func (s *Storage) Prune(migrator *migrate.Migrate) error {
	err := migrator.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return s.Badger.Prune()
}

func getQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
