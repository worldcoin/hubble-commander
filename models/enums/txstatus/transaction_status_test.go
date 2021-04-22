package txstatus

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSON_SupportedStatus(t *testing.T) {
	input := `"PENDING"`
	var res TransactionStatus
	err := json.Unmarshal([]byte(input), &res)
	require.NoError(t, err)
	require.Equal(t, Pending, res)
}

func TestUnmarshalJSON_UnsupportedStatus(t *testing.T) {
	input := `"NOT_SUPPORTED"`
	var res TransactionStatus
	err := json.Unmarshal([]byte(input), &res)
	require.Error(t, err)
	require.Equal(t, enumerr.NewUnsupportedError("transaction status"), err)
	require.True(t, enumerr.IsUnsupportedError(err))
}

func TestMarshalJSON_SupportedStatus(t *testing.T) {
	input := InBatch
	expected := fmt.Sprintf(`"%s"`, TransactionStatuses[input])
	bytes, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))
}

func TestMarshalJSON_UnsupportedStatus(t *testing.T) {
	input := TransactionStatus(0)
	bytes, err := json.Marshal(input)
	require.Error(t, err)
	require.Nil(t, bytes)
	require.Equal(t, enumerr.NewUnsupportedError("transaction status"), errors.Unwrap(err))
	require.True(t, enumerr.IsUnsupportedError(err))
}
