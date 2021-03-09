package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
)

func NewTestStorage(database *db.Database) *Storage {
	return &Storage{DB: database, QB: getQueryBuilder()}
}
