package db

import (
	"reflect"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type Database struct {
	store             *bh.Store
	txn               *badger.Txn
	updateTransaction bool
}

func NewDatabase(cfg *config.BadgerConfig) (*Database, error) {
	options := badger.DefaultOptions(cfg.Path).
		WithLoggingLevel(badger.WARNING).
		WithMemTableSize(64 << 22)
	return newConfiguredDatabase(&options)
}

func NewInMemoryDatabase() (*Database, error) {
	options := badger.DefaultOptions("").
		WithInMemory(true).
		WithLoggingLevel(badger.WARNING)
	return newConfiguredDatabase(&options)
}

func newConfiguredDatabase(opts *badger.Options) (*Database, error) {
	bhOptions := bh.DefaultOptions
	bhOptions.Encoder = Encode
	bhOptions.Decoder = Decode
	bhOptions.Options = *opts

	store, err := bh.Open(bhOptions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Database{store: store}, nil
}

// returns a second database which references the same underlying badgerhold store
// this allows us to open a second transaction
func (d *Database) Clone() *Database {
	return &Database{
		store: d.store,
	}
}

func (d *Database) Close() error {
	err := d.store.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) TriggerGC() error {
	// We're garbage collecting the value log files, which are a combination WAL and
	// mechamism for keeping large values out of the LSM levels which are frequently
	// copied between files as levels are compacted. Each value log file is a
	// collection of values, some number of which refer to inaccessible versions which
	// should be garbage collected. The garbage collector works at the granularity of
	// individual files, in order to garbage collect a file it reads all the records
	// and creates a new file without all the discardable ones. The float we're
	// passing in is the proportion of the records in a value log file which must be
	// discardable for that value log file to be rewritten. e.g. if we pass 1.0 then
	// value log files will be rewritten if they contain even a single discardable
	// record. 0.5 is the recommended value and it seems fine, worth revisiting this
	// number if we end up using too much disk space or if we're doing too much disk
	// I/O.
	err := d.store.Badger().RunValueLogGC(0.5)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) duringTransaction() bool {
	return d.txn != nil
}

func (d *Database) duringReadOnlyTransaction() bool {
	return d.duringTransaction() && !d.updateTransaction
}

func (d *Database) duringUpdateTransaction() bool {
	return d.duringTransaction() && d.updateTransaction
}

func (d *Database) View(fn func(txn *badger.Txn) error) error {
	if d.duringTransaction() {
		return fn(d.txn)
	}
	err := d.store.Badger().View(fn)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) RawUpdate(fn func(txn *badger.Txn) error) error {
	if d.duringReadOnlyTransaction() {
		panic("RawUpdate called during ReadOnly transaction")
	}
	if d.duringUpdateTransaction() {
		return fn(d.txn)
	}
	err := d.store.Badger().Update(fn)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) Count(result interface{}, query *bh.Query) (uint64, error) {
	if d.duringTransaction() {
		count, err := d.store.TxCount(d.txn, result, query)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		return count, nil
	}
	count, err := d.store.Count(result, query)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return count, nil
}

func (d *Database) Find(result interface{}, query *bh.Query) error {
	if d.duringTransaction() {
		err := d.store.TxFind(d.txn, result, query)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	err := d.store.Find(result, query)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) FindOne(result interface{}, query *bh.Query) error {
	if d.duringTransaction() {
		err := d.store.TxFindOne(d.txn, result, query)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	err := d.store.FindOne(result, query)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// The BadgerHold implementation of FindOne does not use indexes, this is a bit of a hack
// which adapts our call into a call to Find, which does use indexes.
func (d *Database) FindOneUsingIndex(result, key interface{}, index string) error {
	typeResult := reflect.TypeOf(result)
	if typeResult.Kind() != reflect.Ptr {
		panic("`result` must be a pointer")
	}
	typeResult = typeResult.Elem()
	valResults := reflect.New(reflect.SliceOf(typeResult))

	err := d.Find(
		valResults.Interface(),
		bh.Where(index).Eq(key).Index(index).Limit(1),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	if valResults.Elem().Len() == 0 {
		return errors.WithStack(bh.ErrNotFound)
	}

	valResult := reflect.ValueOf(result)
	valResult.Elem().Set(valResults.Elem().Index(0))
	return nil
}

func (d *Database) Get(key, result interface{}) error {
	if d.duringTransaction() {
		err := d.store.TxGet(d.txn, key, result)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	err := d.store.Get(key, result)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) Insert(key, data interface{}) error {
	if d.duringReadOnlyTransaction() {
		panic("Insert called during ReadOnly transaction")
	}
	if d.duringUpdateTransaction() {
		err := d.store.TxInsert(d.txn, key, data)
		if errors.Is(err, bh.ErrKeyExists) {
			return errors.Wrapf(err, "duplicate key: %x", key)
		}
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	err := d.store.Insert(key, data)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) Upsert(key, data interface{}) error {
	if d.duringReadOnlyTransaction() {
		panic("Upsert called during ReadOnly transaction")
	}
	if d.duringUpdateTransaction() {
		err := d.store.TxUpsert(d.txn, key, data)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	err := d.store.Upsert(key, data)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) Update(key, data interface{}) error {
	if d.duringReadOnlyTransaction() {
		panic("Update called during ReadOnly transaction")
	}
	if d.duringUpdateTransaction() {
		err := d.store.TxUpdate(d.txn, key, data)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	err := d.store.Update(key, data)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) Delete(key, dataType interface{}) error {
	if d.duringReadOnlyTransaction() {
		panic("Delete called during ReadOnly transaction")
	}
	if d.duringUpdateTransaction() {
		err := d.store.TxDelete(d.txn, key, dataType)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	err := d.store.Delete(key, dataType)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *Database) BeginTransaction(update bool) (*TxController, *Database) {
	if d.duringTransaction() {
		return NewTxController(&ControllerAdapter{d.txn}, true), d
	}
	txn := d.store.Badger().NewTransaction(update)
	dbDuringTx := &Database{
		store:             d.store,
		txn:               txn,
		updateTransaction: update,
	}
	return NewTxController(&ControllerAdapter{txn}, false), dbDuringTx
}

func (d *Database) Prune() error {
	err := d.store.Badger().DropAll()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
