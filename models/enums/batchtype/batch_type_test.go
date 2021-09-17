package batchtype

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
	"github.com/stretchr/testify/require"
)

func TestBatchType_UnmarshalJSON_SupportedType(t *testing.T) {
	input := `"TRANSFER"`
	var res BatchType
	err := json.Unmarshal([]byte(input), &res)
	require.NoError(t, err)
	require.Equal(t, Transfer, res)
}

func TestBatchType_UnmarshalJSON_UnsupportedType(t *testing.T) {
	input := `"NOT_SUPPORTED"`
	var res BatchType
	err := json.Unmarshal([]byte(input), &res)
	require.Error(t, err)
	require.Equal(t, enumerr.NewUnsupportedError("batch type"), err)
	require.True(t, enumerr.IsUnsupportedError(err))
}

func TestBatchType_MarshalJSON_SupportedType(t *testing.T) {
	input := Create2Transfer
	expected := fmt.Sprintf(`"%s"`, BatchTypes[input])
	bytes, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))
}

func TestBatchType_MarshalJSON_UnsupportedType(t *testing.T) {
	input := BatchType(10)
	bytes, err := json.Marshal(input)
	require.Error(t, err)
	require.Nil(t, bytes)
	require.Equal(t, enumerr.NewUnsupportedError("batch type"), errors.Unwrap(err))
	require.True(t, enumerr.IsUnsupportedError(err))
}
