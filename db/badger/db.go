package badger

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/dgraph-io/badger/v3"
)

type Database struct {
	badger            *badger.DB
	txn               *badger.Txn
	updateTransaction bool
}

func NewDatabase(cfg *config.BadgerConfig) (*Database, error) {
	database, err := badger.Open(badger.DefaultOptions(cfg.Path))
	if err != nil {
		return nil, err
	}
	return &Database{badger: database}, nil
}

func (d *Database) Close() error {
	return d.badger.Close()
}

func (d *Database) duringTransaction() bool {
	return d.txn != nil
}

func (d *Database) duringUpdateTransaction() bool {
	return d.duringTransaction() && d.updateTransaction
}

func (d *Database) View(fn func(txn *badger.Txn) error) error {
	if d.duringTransaction() {
		return fn(d.txn)
	}
	return d.badger.View(fn)
}

func (d *Database) Update(fn func(txn *badger.Txn) error) error {
	if d.duringUpdateTransaction() {
		return fn(d.txn)
	}
	return d.badger.Update(fn)
}

func (d *Database) BeginTransaction(update bool) (*db.TxController, *Database) {
	if d.duringTransaction() {
		return db.NewTxController(&ControllerAdapter{d.txn}, true), d
	}
	txn := d.badger.NewTransaction(update)
	dbDuringTx := &Database{
		badger:            d.badger,
		txn:               txn,
		updateTransaction: update,
	}
	return db.NewTxController(&ControllerAdapter{txn}, false), dbDuringTx
}
