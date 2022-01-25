package batchstatus

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
	"github.com/stretchr/testify/require"
)

func TestBatchStatus_UnmarshalJSON_SupportedStatus(t *testing.T) {
	input := `"PENDING"`
	var res BatchStatus
	err := json.Unmarshal([]byte(input), &res)
	require.NoError(t, err)
	require.Equal(t, Pending, res)
}

func TestBatchStatus_UnmarshalJSON_UnsupportedStatus(t *testing.T) {
	input := `"NOT_SUPPORTED"`
	var res BatchStatus
	err := json.Unmarshal([]byte(input), &res)
	require.Error(t, err)
	require.Equal(t, enumerr.NewUnsupportedError("batch status"), err)
	require.True(t, enumerr.IsUnsupportedError(err))
}

func TestBatchStatus_MarshalJSON_SupportedStatus(t *testing.T) {
	input := Mined
	expected := fmt.Sprintf(`%q`, BatchStatuses[input])
	bytes, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))
}

func TestBatchStatus_MarshalJSON_UnsupportedStatus(t *testing.T) {
	input := BatchStatus(0)
	bytes, err := json.Marshal(input)
	require.Error(t, err)
	require.Nil(t, bytes)
	require.Equal(t, enumerr.NewUnsupportedError("batch status"), errors.Unwrap(err))
	require.True(t, enumerr.IsUnsupportedError(err))
}
