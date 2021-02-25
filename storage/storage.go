package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB *sqlx.DB
	QB squirrel.StatementBuilderType
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	dbInstance, err := db.GetDB(cfg)
	if err != nil {
		return nil, err
	}
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Storage{DB: dbInstance, QB: queryBuilder}, nil
}
