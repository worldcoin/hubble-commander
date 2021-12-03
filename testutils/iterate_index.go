package testutils

import (
	"github.com/Worldcoin/hubble-commander/db"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"
	bh "github.com/timshannon/badgerhold/v4"
)

func IterateIndex(
	s *require.Assertions,
	badger *db.Database,
	typeName []byte,
	indexName string,
	handleIndex func(encodedKey []byte, keyList bh.KeyList),
) {
	indexPrefix := db.IndexKeyPrefix(typeName, indexName)
	err := badger.Iterator(indexPrefix, db.PrefetchIteratorOpts, func(item *bdg.Item) (finish bool, err error) {
		// Get key value
		encodedKeyValue := item.Key()[len(indexPrefix):]

		// Decode value
		var keyList bh.KeyList
		err = item.Value(func(val []byte) error {
			return db.Decode(val, &keyList)
		})
		s.NoError(err)

		handleIndex(encodedKeyValue, keyList)
		return false, nil
	})
	s.ErrorIs(err, db.ErrIteratorFinished)
}
