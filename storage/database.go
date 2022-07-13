package storage

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
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

func (d *Database) BeginTransaction(opts TxOptions) (*db.TxController, *Database) {
	database := *d

	badgerTx, badgerDB := d.Badger.BeginTransaction(!opts.ReadOnly)
	database.Badger = badgerDB

	return badgerTx, &database
}

// all errors are already wrapped w stack traces, except errors fn returns
func (d *Database) ExecuteInTransaction(opts TxOptions, fn func(txDatabase *Database) error) error {
	err := d.unsafeExecuteInTransaction(opts, fn)
	if errors.Is(err, bdg.ErrConflict) {
		// nb. if we were already inside a transaction when this function is
		//     called then we run inside the outer transaction, so the `Commit`
		//     call is a no-op and this retry logic will never get a chance to
		//     fire
		log.Debug("ExecuteInTransaction transaction conflicted, trying again")
		return d.ExecuteInTransaction(opts, fn)
	}
	return err
}

func (d *Database) unsafeExecuteInTransaction(opts TxOptions, fn func(txDatabase *Database) error) (err error) {
	txController, txDatabase := d.BeginTransaction(opts)
	defer txController.Rollback(&err)

	err = fn(txDatabase)
	if err != nil {
		return err
	}

	return txController.Commit()
}

func (d *Database) Close() error {
	return d.Badger.Close()
}
