package dto

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHexString_MarshalJSON_NonEmpty(t *testing.T) {
	hexBytes, err := hex.DecodeString("deadbeef")
	require.NoError(t, err)
	hexString := HexString(hexBytes)

	bytes, err := json.Marshal(hexString)
	require.NoError(t, err)
	require.Equal(t, `"0xdeadbeef"`, string(bytes))
}

func TestHexString_MarshalJSON_Empty(t *testing.T) {
	bytes, err := json.Marshal(HexString{})
	require.NoError(t, err)
	require.Equal(t, `""`, string(bytes))
}

func TestHexString_UnmarshalJSON_NonEmpty(t *testing.T) {
	expected, err := hex.DecodeString("deadbeef")
	require.NoError(t, err)

	var dest HexString
	err = json.Unmarshal([]byte(`"0xdeadbeef"`), &dest)
	require.NoError(t, err)
	require.Equal(t, HexString(expected), dest)
}

func TestHexString_UnmarshalJSON_0x(t *testing.T) {
	var dest HexString
	err := json.Unmarshal([]byte(`"0x"`), &dest)
	require.NoError(t, err)
	require.Equal(t, HexString{}, dest)
}

func TestHexString_UnmarshalJSON_Empty(t *testing.T) {
	var dest HexString
	err := json.Unmarshal([]byte(`""`), &dest)
	require.NoError(t, err)
	require.Equal(t, HexString{}, dest)
}

func TestHexString_UnmarshalJSON_Error(t *testing.T) {
	var dest HexString
	err := json.Unmarshal([]byte(`"deadbeef"`), &dest)
	require.Equal(t, ErrHexStringNotPrepended, err)
}
