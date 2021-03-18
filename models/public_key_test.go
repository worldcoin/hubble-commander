package models

import (
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
