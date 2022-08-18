package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

// We need to "initialize" the indices on fields of pointer type to make them work with bh.Find operations.
// The problem originates in `indexExists` function in BadgerHold (https://github.com/timshannon/badgerhold/blob/v4.0.1/index.go#L148).
// Badger assumes that if there is a value for some data type, then there must exist at least one index entry.
// If you don't index nil values the way we did for stored.TxReceipt.ToStateID it can be the case that there is some
// stored.TxReceipt stored, but there is no index entry. To work around this we set an empty index entry.
// See:
//   - stored.TxReceipt Indexes() method
//   - InitializeIndexTestSuite.TestStoredTxReceipt_ToStateID_FindUsingIndexWorksWhenThereAreOnlyValuesWithThisFieldSetToNil
//
// Note, this does not create the index. In order to create the index we must implement
// Indexes() on the relevant type.
func initializeIndex(database *Database, typeName []byte, indexName string, zeroValue interface{}) error {
	encodedZeroValue, err := db.Encode(zeroValue)
	if err != nil {
		return err
	}
	zeroValueIndexKey := db.IndexKey(typeName, indexName, encodedZeroValue)

	initialized, err := indexAlreadyInitialised(database, zeroValueIndexKey)
	if initialized || err != nil {
		return err
	}

	emptyKeyList := make(bh.KeyList, 0)
	encodedEmptyKeyList, err := db.Encode(emptyKeyList)
	if err != nil {
		return err
	}

	return database.Badger.RawUpdate(func(txn *bdg.Txn) error {
		return txn.Set(zeroValueIndexKey, encodedEmptyKeyList)
	})
}

func indexAlreadyInitialised(database *Database, indexKey []byte) (bool, error) {
	err := database.Badger.View(func(txn *bdg.Txn) error {
		_, err := txn.Get(indexKey)
		return err
	})
	if errors.Is(err, bdg.ErrKeyNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
