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

func TestIncomingTransaction_JSONMarshaling(t *testing.T) {
	tx := IncomingTransaction{
		FromIndex: NewUint256(1),
		ToIndex:   NewUint256(2),
		Amount:    NewUint256(50),
		Fee:       NewUint256(10),
		Nonce:     NewUint256(0),
		Signature: []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}
	data, err := json.Marshal(tx)
	require.NoError(t, err)

	var unmarshalled IncomingTransaction
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, tx, unmarshalled)
}
