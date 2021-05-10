package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/db/postgres"
)

type Storage struct {
	Postgres *postgres.Database
	QB       squirrel.StatementBuilderType
}

func NewStorage(cfg *config.DBConfig) (*Storage, error) {
	dbInstance, err := postgres.NewDatabase(cfg)
	if err != nil {
		return nil, err
	}
	return &Storage{Postgres: dbInstance, QB: getQueryBuilder()}, nil
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

func getQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
