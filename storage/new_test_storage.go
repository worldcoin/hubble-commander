package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/db"
)

func NewTestStorage(db *db.Database) *Storage {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &Storage{DB: db, QB: queryBuilder}
}
