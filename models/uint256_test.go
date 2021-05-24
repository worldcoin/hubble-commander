package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUint256_JSONMarshaling(t *testing.T) {
	value := MakeUint256(5)
	data, err := json.Marshal(value)
	require.NoError(t, err)

	var unmarshalled Uint256
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, value, unmarshalled)
}

func TestUint256_JSONPtrMarshaling(t *testing.T) {
	value := NewUint256(5)
	data, err := json.Marshal(value)
	require.NoError(t, err)

	var unmarshalled *Uint256
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, value, unmarshalled)
}

func TestUint256_UnmarshalString(t *testing.T) {
	var unmarshalled Uint256
	err := json.Unmarshal([]byte("\"5\""), &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, MakeUint256(5), unmarshalled)

	var unmarshalledPtr *Uint256
	err = json.Unmarshal([]byte("\"5\""), &unmarshalledPtr)
	require.NoError(t, err)

	require.Equal(t, MakeUint256(5), *unmarshalledPtr)
}

func TestUint256_UnmarshalNumber(t *testing.T) {
	var unmarshalled Uint256
	err := json.Unmarshal([]byte("5123123"), &unmarshalled)
	require.Error(t, err)
}

func TestUint256_Add(t *testing.T) {
	one := NewUint256(1)
	two := NewUint256(2)
	three := NewUint256(3)
	require.Equal(t, three, one.Add(two))
}

func TestUint256_Sub(t *testing.T) {
	one := NewUint256(1)
	two := NewUint256(2)
	three := NewUint256(3)
	require.Equal(t, one, three.Sub(two))
}

func TestUint256_Mul(t *testing.T) {
	two := NewUint256(2)
	four := NewUint256(4)
	require.Equal(t, four, two.Mul(two))
}

func TestUint256_Div(t *testing.T) {
	two := NewUint256(2)
	four := NewUint256(4)
	require.Equal(t, two, four.Div(two))
}
