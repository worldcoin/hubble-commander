package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPadLeft(t *testing.T) {
	require.Equal(t, []byte{0, 0, 1}, PadLeft([]byte{1}, 3))
	require.Equal(t, []byte{0, 0, 1}, PadLeft([]byte{0, 0, 1}, 3))
	require.Equal(t, []byte{0, 0, 1}, PadLeft([]byte{0, 0, 1}, 2))
}

func TestByteSliceTo32ByteSlice(t *testing.T) {
	require.Equal(t, [32]byte{0, 0, 1}, ByteSliceTo32ByteSlice([]byte{0, 0, 1}))
	require.Equal(t, [32]byte{1, 2, 3}, ByteSliceTo32ByteSlice([]byte{1, 2, 3}))
}
