package storage

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	Badger *db.Database
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	badgerDB, err := db.NewDatabase(cfg.Badger)
	if err != nil {
		return nil, err
	}

	database := &Database{
		Badger: badgerDB,
	}

	if cfg.Bootstrap.Prune {
		err = database.Badger.Prune()
		if err != nil {
			return nil, err
		}
		log.Debug("Badger database was pruned")
	}

	return database, nil
}

func (d *Database) BeginTransaction(opts TxOptions) (*db.TxController, *Database, error) {
	database := *d

	badgerTx, badgerDB := d.Badger.BeginTransaction(!opts.ReadOnly)
	database.Badger = badgerDB

	return badgerTx, &database, nil
}

func (d *Database) Close() error {
	return d.Badger.Close()
}
