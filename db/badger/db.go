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
	options.Options = badger.
		DefaultOptions(cfg.Path).
		WithLoggingLevel(badger.WARNING).
		WithMemTableSize(64 << 21) // TODO: Bench to see if there are performance degradations and remove if so.

	store, err := bh.Open(options)
	if err != nil {
		return nil, err
	}
	return &Database{store: store}, nil
}

func (d *Database) Close() error {
	return d.store.Close()
}

func (d *Database) BadgerInstance() *badger.DB {
	return d.store.Badger()
}

func (d *Database) Tx() *badger.Txn {
	return d.txn
}
 
func (d *Database) DuringTransaction() bool {
	return d.txn != nil
}

func (d *Database) duringUpdateTransaction() bool {
	return d.DuringTransaction() && d.updateTransaction
}

func (d *Database) Find(result interface{}, query *bh.Query) error {
	if d.DuringTransaction() {
		return d.store.TxFind(d.txn, result, query)
	}
	return d.store.Find(result, query)
}

func (d *Database) Get(key, result interface{}) error {
	if d.DuringTransaction() {
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

func (d *Database) Delete(key, dataType interface{}) error {
	if d.duringUpdateTransaction() {
		return d.store.TxDelete(d.txn, key, dataType)
	}
	return d.store.Delete(key, dataType)
}

func (d *Database) BeginTransaction(update bool) (*db.TxController, *Database) {
	if d.DuringTransaction() {
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

func (d *Database) Prune() error {
	return d.store.Badger().DropAll()
}
