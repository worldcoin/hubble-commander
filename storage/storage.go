package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db/postgres"
)

type Storage struct {
	DB *postgres.Database
	QB squirrel.StatementBuilderType
}

func NewStorage(cfg *config.DBConfig) (*Storage, error) {
	dbInstance, err := postgres.NewDatabase(cfg)
	if err != nil {
		return nil, err
	}
	return &Storage{DB: dbInstance, QB: getQueryBuilder()}, nil
}

func (s *Storage) BeginTransaction() (*postgres.TransactionController, *Storage, error) {
	tx, txDB, err := s.DB.BeginTransaction()
	if err != nil {
		return nil, nil, err
	}

	storage := &Storage{
		DB: txDB,
		QB: getQueryBuilder(),
	}

	return tx, storage, nil
}

func getQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
