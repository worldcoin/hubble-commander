package commitmentstatus

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
	"github.com/stretchr/testify/require"
)

func TestCommitmentStatus_UnmarshalJSON_SupportedStatus(t *testing.T) {
	input := `"PENDING"`
	var res CommitmentStatus
	err := json.Unmarshal([]byte(input), &res)
	require.NoError(t, err)
	require.Equal(t, Pending, res)
}

func TestCommitmentStatus_UnmarshalJSON_UnsupportedStatus(t *testing.T) {
	input := `"NOT_SUPPORTED"`
	var res CommitmentStatus
	err := json.Unmarshal([]byte(input), &res)
	require.Error(t, err)
	require.Equal(t, enumerr.NewUnsupportedError("commitment status"), err)
	require.True(t, enumerr.IsUnsupportedError(err))
}

func TestCommitmentStatus_MarshalJSON_SupportedStatus(t *testing.T) {
	input := Mined
	expected := fmt.Sprintf(`%q`, CommitmentStatuses[input])
	bytes, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))
}

func TestCommitmentStatus_MarshalJSON_UnsupportedStatus(t *testing.T) {
	input := CommitmentStatus(0)
	bytes, err := json.Marshal(input)
	require.Error(t, err)
	require.Nil(t, bytes)
	require.Equal(t, enumerr.NewUnsupportedError("commitment status"), errors.Unwrap(err))
	require.True(t, enumerr.IsUnsupportedError(err))
}
