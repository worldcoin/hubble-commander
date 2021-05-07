package storage

import (
	"github.com/Worldcoin/hubble-commander/db/postgres"
)

func NewTestStorage(database *postgres.Database) *Storage {
	return &Storage{DB: database, QB: getQueryBuilder()}
}
