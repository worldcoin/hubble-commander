package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/db"
)

func NewTestStorage(database *db.Database) *Storage {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &Storage{DB: database, QB: queryBuilder}
}
