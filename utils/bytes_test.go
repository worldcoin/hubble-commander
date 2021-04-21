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
