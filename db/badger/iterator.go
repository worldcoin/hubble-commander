package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

var ErrIteratorFinished = errors.New("iterator finished")

func (d *Database) ReverseIterator(prefix []byte, filter func(item *badger.Item) (bool, error)) error {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Reverse = true
	return d.Iterator(prefix, opts, filter)
}

func (d *Database) Iterator(prefix []byte, opts badger.IteratorOptions, filter func(item *badger.Item) (bool, error)) error {
	return d.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(opts)
		defer it.Close()

		seekPrefix := newSeekPrefix(prefix, opts)
		for it.Seek(seekPrefix); it.ValidForPrefix(prefix); it.Next() {
			finish, err := filter(it.Item())
			if err != nil {
				return err
			}
			if finish {
				return nil
			}
		}
		return ErrIteratorFinished
	})
}

func newSeekPrefix(prefix []byte, opts badger.IteratorOptions) []byte {
	if opts.Reverse {
		return append(prefix, 0xFF) // Required to loop backwards
	}
	return prefix
}
