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
	options.Encoder = Encode
	options.Decoder = Decode
	options.Options = badger.
		DefaultOptions(cfg.Path).
		WithLoggingLevel(badger.WARNING)
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

func (d *Database) View(fn func(txn *badger.Txn) error) error {
	if d.duringTransaction() {
		return fn(d.txn)
	}
	return d.store.Badger().View(fn)
}

func (d *Database) Find(result interface{}, query *bh.Query) error {
	if d.duringTransaction() {
		return d.store.TxFind(d.txn, result, query)
	}
	return d.store.Find(result, query)
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

func (d *Database) Delete(key, dataType interface{}) error {
	if d.duringUpdateTransaction() {
		return d.store.TxDelete(d.txn, key, dataType)
	}
	return d.store.Delete(key, dataType)
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

func (d *Database) Prune() error {
	return d.store.Badger().DropAll()
}

func PruneDatabase(cfg *config.BadgerConfig) error {
	database, err := NewDatabase(cfg)
	if err != nil {
		return err
	}
	defer func() {
		err = database.Close()
		if err != nil {
			panic(err)
		}
	}()
	return database.Prune()
}
