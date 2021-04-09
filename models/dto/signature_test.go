package dto

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignature_MarshalJSON_NonEmpty(t *testing.T) {
	sig, err := hex.DecodeString("deadbeef")
	require.NoError(t, err)
	signature := Signature(sig)

	bytes, err := json.Marshal(signature)
	require.NoError(t, err)
	require.Equal(t, `"0xdeadbeef"`, string(bytes))
}

func TestSignature_MarshalJSON_Empty(t *testing.T) {
	signature := Signature{}
	bytes, err := json.Marshal(signature)
	require.NoError(t, err)
	require.Equal(t, `""`, string(bytes))
}

func TestSignature_UnmarshalJSON_NonEmpty(t *testing.T) {
	expected, err := hex.DecodeString("deadbeef")
	require.NoError(t, err)

	var dest Signature
	err = json.Unmarshal([]byte(`"0xdeadbeef"`), &dest)
	require.NoError(t, err)
	require.Equal(t, Signature(expected), dest)
}

func TestSignature_UnmarshalJSON_0x(t *testing.T) {
	var dest Signature
	err := json.Unmarshal([]byte(`"0x"`), &dest)
	require.NoError(t, err)
	require.Equal(t, Signature{}, dest)
}

func TestSignature_UnmarshalJSON_Empty(t *testing.T) {
	var dest Signature
	err := json.Unmarshal([]byte(`""`), &dest)
	require.NoError(t, err)
	require.Equal(t, Signature{}, dest)
}

func TestSignature_UnmarshalJSON_Error(t *testing.T) {
	var dest Signature
	err := json.Unmarshal([]byte(`"deadbeef"`), &dest)
	require.Errorf(t, err, "hex string must be 0x prepended")
}
