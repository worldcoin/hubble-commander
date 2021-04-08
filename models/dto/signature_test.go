package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignature_MarshalJSON(t *testing.T) {
	signature := Signature{}
	bytes, err := json.Marshal(signature)
	require.NoError(t, err)
	require.Equal(t, `""`, string(bytes))
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
