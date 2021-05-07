package badger

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"
)

func TestNewTestDB(t *testing.T) {
	bdg, err := NewTestDB()
	require.NoError(t, err)

	key := []byte{1, 2, 3}
	value := []byte{2, 3, 4}

	err = bdg.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	require.NoError(t, err)

	err = bdg.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key) // nolint:govet
		require.NoError(t, err)
		actual, err := item.ValueCopy(nil)
		require.NoError(t, err)
		require.Equal(t, value, actual)
		return nil
	})
	require.NoError(t, err)
}
