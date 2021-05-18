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
	Postgres    *postgres.Database
	Badger      *badger.Database
	QB          squirrel.StatementBuilderType
	domain      *bls.Domain
	feeReceiver map[string]uint32
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
		Postgres:    postgresDB,
		Badger:      badgerDB,
		QB:          getQueryBuilder(),
		feeReceiver: make(map[string]uint32),
	}, nil
}

func (s *Storage) BeginTransaction(opts TxOptions) (*db.TxController, *Storage, error) {
	var txController *db.TxController
	var storage Storage
	storage.Postgres = s.Postgres
	storage.QB = getQueryBuilder()

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
