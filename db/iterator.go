package db

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

var (
	ErrIteratorFinished = errors.New("iterator finished")

	ReverseKeyIteratorOpts = badger.IteratorOptions{
		PrefetchValues: false,
		Reverse:        true,
	}
	ReversePrefetchIteratorOpts = badger.IteratorOptions{
		Reverse:        true,
		PrefetchSize:   100,
		PrefetchValues: true,
	}
	PrefetchIteratorOpts = badger.IteratorOptions{
		PrefetchSize:   100,
		PrefetchValues: true,
	}
	KeyIteratorOpts = badger.IteratorOptions{
		PrefetchValues: false,
		Reverse:        false,
	}

	Continue = false
	Break    = true
)

type IteratorFilter func(item *badger.Item) (finish bool, err error)

// Iterator calls filter function for every element matching the prefix.
// First return value of the filter function is the finish flag.
func (d *Database) Iterator(prefix []byte, opts badger.IteratorOptions, filter IteratorFilter) error {
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
		newPrefix := make([]byte, 0, len(prefix)+1)
		newPrefix = append(newPrefix, prefix...)
		return append(newPrefix, 0xFF) // Required to loop backwards
	}
	return prefix
}
