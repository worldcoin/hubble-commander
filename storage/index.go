package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	bdg "github.com/dgraph-io/badger/v3"
	bh "github.com/timshannon/badgerhold/v4"
)

func initializeIndex(database *Database, typeName []byte, indexName string, zeroValue interface{}) error {
	encodedZeroValue, err := db.Encode(zeroValue)
	if err != nil {
		return err
	}
	zeroValueIndexKey := db.IndexKey(typeName, indexName, encodedZeroValue)

	emptyKeyList := make(bh.KeyList, 0)
	encodedEmptyKeyList, err := db.Encode(emptyKeyList)
	if err != nil {
		return err
	}

	return database.Badger.RawUpdate(func(txn *bdg.Txn) error {
		return txn.Set(zeroValueIndexKey, encodedEmptyKeyList)
	})
}
