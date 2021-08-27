package models

import (
	"encoding/json"
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestPublicKey_Bytes_ReturnsACopy(t *testing.T) {
	key := PublicKey{1, 2, 3}
	bytes := key.Bytes()
	bytes[0] = 9
	require.Equal(t, PublicKey{1, 2, 3}, key)
}

func TestPublicKey_SetBytes(t *testing.T) {
	key := PublicKey{1, 2, 3}
	bytes := key.Bytes()
	newKey := PublicKey{}
	err := newKey.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, key, newKey)
}

func TestPublicKey_SetBytes_InvalidLength(t *testing.T) {
	bytes := utils.PadLeft([]byte{1, 2, 3}, 130)
	key := PublicKey{}
	err := key.SetBytes(bytes)
	require.NotNil(t, err)
	require.ErrorIs(t, err, ErrInvalidPublicKeyLength)
}

func TestPublicKey_Value_ReturnsACopy(t *testing.T) {
	key := PublicKey{1, 2, 3}
	value, err := key.Value()
	require.NoError(t, err)
	bytes, ok := value.([]byte)
	require.True(t, ok)
	bytes[0] = 9
	require.Equal(t, PublicKey{1, 2, 3}, key)
}

func TestPublicKeyPublicKey_JSONMarshaling(t *testing.T) {
	key := PublicKey{1, 2, 3}
	data, err := json.Marshal(key)
	require.NoError(t, err)

	var unmarshalled PublicKey
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, key, unmarshalled)
}

func TestPublicKey_YAMLMarshaling(t *testing.T) {
	key := PublicKey{1, 2, 3}
	data, err := yaml.Marshal(key)
	require.NoError(t, err)

	var unmarshalled PublicKey
	err = yaml.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, key, unmarshalled)
}
