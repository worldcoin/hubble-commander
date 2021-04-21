package txtype

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSON_SupportedType(t *testing.T) {
	input := `"TRANSFER"`
	var res TransactionType
	err := json.Unmarshal([]byte(input), &res)
	require.NoError(t, err)
	require.Equal(t, Transfer, res)
}

func TestUnmarshalJSON_UnsupportedType(t *testing.T) {
	input := `"NOT_SUPPORTED"`
	var res TransactionType
	err := json.Unmarshal([]byte(input), &res)
	require.Error(t, err)
	require.Equal(t, ErrUnsupportedTransactionType, err)
}

func TestMarshalJSON_SupportedType(t *testing.T) {
	input := Create2Transfer
	expected := fmt.Sprintf(`"%s"`, TransactionTypes[input])
	bytes, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))
}

func TestMarshalJSON_UnsupportedType(t *testing.T) {
	input := TransactionType(0)
	bytes, err := json.Marshal(input)
	require.Error(t, err)
	require.Nil(t, bytes)
	require.Equal(t, ErrUnsupportedTransactionType, errors.Unwrap(err))
}
