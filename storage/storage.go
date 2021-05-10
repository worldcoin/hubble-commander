package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/db/postgres"
)

type Storage struct {
	Postgres *postgres.Database
	Badger   *badger.Database
	QB       squirrel.StatementBuilderType
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

	return &Storage{Postgres: postgresDB, Badger: badgerDB, QB: getQueryBuilder()}, nil
}

func (s *Storage) BeginTransaction() (*db.TxController, *Storage, error) {
	tx, txDB, err := s.Postgres.BeginTransaction()
	if err != nil {
		return nil, nil, err
	}

	storage := &Storage{
		Postgres: txDB,
		QB:       getQueryBuilder(),
	}

	return tx, storage, nil
}

func (s *Storage) Close() error {
	err := s.Postgres.Close()
	if err != nil {
		return err
	}
	return s.Badger.Close()
}

func getQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
