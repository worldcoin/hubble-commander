package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytes_ReturnsACopy(t *testing.T) {
	key := PublicKey{1, 2, 3}
	bytes := key.Bytes()
	bytes[0] = 9
	require.Equal(t, PublicKey{1, 2, 3}, key)
}

func TestValue_ReturnsACopy(t *testing.T) {
	key := PublicKey{1, 2, 3}
	value, err := key.Value()
	require.NoError(t, err)
	bytes, ok := value.([]byte)
	require.True(t, ok)
	bytes[0] = 9
	require.Equal(t, PublicKey{1, 2, 3}, key)
}

func TestPublicKey_JSONMarshaling(t *testing.T) {
	key := PublicKey{1, 2, 3}
	data, err := json.Marshal(key)
	require.NoError(t, err)

	var unmarshalled PublicKey
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, key, unmarshalled)
}
