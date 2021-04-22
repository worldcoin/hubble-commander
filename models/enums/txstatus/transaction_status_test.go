package txstatus

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSON_SupportedType(t *testing.T) {
	input := `"PENDING"`
	var res TransactionStatus
	err := json.Unmarshal([]byte(input), &res)
	require.NoError(t, err)
	require.Equal(t, Pending, res)
}

func TestUnmarshalJSON_UnsupportedType(t *testing.T) {
	input := `"NOT_SUPPORTED"`
	var res TransactionStatus
	err := json.Unmarshal([]byte(input), &res)
	require.Error(t, err)
	require.Equal(t, ErrUnsupportedTransactionStatus, err)
}

func TestMarshalJSON_SupportedType(t *testing.T) {
	input := Finalised
	expected := fmt.Sprintf(`"%s"`, TransactionStatuses[input])
	bytes, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))
}

func TestMarshalJSON_UnsupportedType(t *testing.T) {
	input := TransactionStatus(0)
	bytes, err := json.Marshal(input)
	require.Error(t, err)
	require.Nil(t, bytes)
	require.Equal(t, ErrUnsupportedTransactionStatus, errors.Unwrap(err))
}
