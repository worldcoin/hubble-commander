package badger

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/dgraph-io/badger/v3"
	bh "github.com/timshannon/badgerhold/v3"
)

type Database struct {
	store             *bh.Store
	txn               *badger.Txn
	updateTransaction bool
}

func NewDatabase(cfg *config.BadgerConfig) (*Database, error) {
	options := bh.DefaultOptions
	options.Options = badger.DefaultOptions(cfg.Path)

	store, err := bh.Open(options)
	if err != nil {
		return nil, err
	}
	return &Database{store: store}, nil
}

func (d *Database) Close() error {
	return d.store.Close()
}

func (d *Database) duringTransaction() bool {
	return d.txn != nil
}

func (d *Database) duringUpdateTransaction() bool {
	return d.duringTransaction() && d.updateTransaction
}

func (d *Database) Get(key, result interface{}) error {
	if d.duringTransaction() {
		return d.store.TxGet(d.txn, key, result)
	}
	return d.store.Get(key, result)
}

func (d *Database) Insert(key, data interface{}) error {
	if d.duringUpdateTransaction() {
		return d.store.TxInsert(d.txn, key, data)
	}
	return d.store.Insert(key, data)
}

func (d *Database) Upsert(key, data interface{}) error {
	if d.duringUpdateTransaction() {
		return d.store.TxUpsert(d.txn, key, data)
	}
	return d.store.Upsert(key, data)
}

func (d *Database) Update(key, data interface{}) error {
	if d.duringUpdateTransaction() {
		return d.store.TxUpdate(d.txn, key, data)
	}
	return d.store.Update(key, data)
}

func (d *Database) BeginTransaction(update bool) (*db.TxController, *Database) {
	if d.duringTransaction() {
		return db.NewTxController(&ControllerAdapter{d.txn}, true), d
	}
	txn := d.store.Badger().NewTransaction(update)
	dbDuringTx := &Database{
		store:             d.store,
		txn:               txn,
		updateTransaction: update,
	}
	return db.NewTxController(&ControllerAdapter{txn}, false), dbDuringTx
}
