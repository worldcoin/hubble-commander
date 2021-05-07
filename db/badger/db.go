package badger

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/dgraph-io/badger/v3"
)

type Database struct {
	*badger.DB
}

func NewDatabase(cfg *config.BadgerConfig) (*Database, error) {
	db, err := badger.Open(badger.DefaultOptions(cfg.Path))
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
