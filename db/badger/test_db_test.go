package badger

import (
	"testing"

	"github.com/stretchr/testify/require"
	bh "github.com/timshannon/badgerhold/v4"
)

type someStruct struct {
	Name string `badgerhold:"key"`
	Age  uint
}

var testStruct = someStruct{
	Name: "Duck",
	Age:  4,
}

func TestNewTestDB(t *testing.T) {
	bdg, err := NewTestDB()
	require.NoError(t, err)

	err = bdg.DB.Insert(testStruct.Name, testStruct)
	require.NoError(t, err)

	var res someStruct
	err = bdg.DB.Get(testStruct.Name, &res)
	require.NoError(t, err)

	require.Equal(t, testStruct, res)
}

func TestPrune(t *testing.T) {
	bdg, err := NewTestDB()
	require.NoError(t, err)

	err = bdg.DB.Insert(testStruct.Name, testStruct)
	require.NoError(t, err)

	err = bdg.DB.Prune()
	require.NoError(t, err)

	var res someStruct
	err = bdg.DB.Get(testStruct.Name, &res)
	require.Equal(t, bh.ErrNotFound, err)
}

func TestTestDB_Clone(t *testing.T) {
	primary, err := NewTestDB()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, primary.Teardown())
	}()

	err = primary.DB.Insert(testStruct.Name, testStruct)
	require.NoError(t, err)

	cloned, err := primary.Clone()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cloned.Teardown())
	}()

	var value someStruct
	err = cloned.DB.Get(testStruct.Name, &value)
	require.NoError(t, err)
	require.Equal(t, testStruct, value)
}
