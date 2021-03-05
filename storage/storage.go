package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
)

type Storage struct {
	DB *db.Database
	QB squirrel.StatementBuilderType
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	dbInstance, err := db.NewDatabase(cfg)
	if err != nil {
		return nil, err
	}
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Storage{DB: dbInstance, QB: queryBuilder}, nil
}

func (s *Storage) BeginTransaction() (*db.TransactionController, *Storage, error) {
	tx, txDB, err := s.DB.BeginTransaction()
	if err != nil {
		return nil, nil, err
	}

	storage := &Storage{
		DB: txDB,
		QB: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	return tx, storage, nil
}
