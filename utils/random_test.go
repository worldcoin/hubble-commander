package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomBytes(t *testing.T) {
	bytes := RandomBytes(32)
	require.Len(t, bytes, 32)
}

func TestRandomHex(t *testing.T) {
	hex := RandomHex(32)
	require.Len(t, hex, 32)
}

func TestRandomHash(t *testing.T) {
	hash1 := RandomHash()
	hash2 := RandomHash()
	require.NotEqual(t, hash1, hash2)
}
