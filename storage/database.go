package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/db/postgres"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	Postgres *postgres.Database
	Badger   *badger.Database
	QB       squirrel.StatementBuilderType
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	badgerDB, err := badger.NewDatabase(cfg.Badger)
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
