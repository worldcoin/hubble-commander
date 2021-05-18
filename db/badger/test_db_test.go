package badger

import (
	"testing"

	"github.com/stretchr/testify/require"
	bh "github.com/timshannon/badgerhold/v3"
)

type someStruct struct {
	Name string `badgerhold:"key"`
	Age  uint
}

func TestNewTestDB(t *testing.T) {
	bdg, err := NewTestDB()
	require.NoError(t, err)

	key := "Duck"
	value := someStruct{
		Name: key,
		Age:  4,
	}

	err = bdg.DB.Insert(key, value)
	require.NoError(t, err)

	var res someStruct
	err = bdg.DB.Get(key, &res)
	require.NoError(t, err)

	require.Equal(t, value, res)
}

func TestPrune(t *testing.T) {
	bdg, err := NewTestDB()
	require.NoError(t, err)

	key := "Duck"
	value := someStruct{
		Name: key,
		Age:  4,
	}

	err = bdg.DB.Insert(key, value)
	require.NoError(t, err)

	err = bdg.DB.Prune()
	require.NoError(t, err)

	var res someStruct
	err = bdg.DB.Get(key, &res)
	require.Equal(t, bh.ErrNotFound, err)
}
