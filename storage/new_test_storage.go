package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func NewTestStorage(db *sqlx.DB) *Storage {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &Storage{DB: db, QB: queryBuilder}
}
