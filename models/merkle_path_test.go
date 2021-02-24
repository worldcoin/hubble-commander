package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMerklePathLength(t *testing.T) {
	// 33 zeroes
	_, err := NewMerklePath("000000000000000000000000000000000")
	require.Error(t, err)
}

func TestAddOne(t *testing.T) {
	a, err := NewMerklePath("0000")
	require.NoError(t, err)

	b, err := a.Add(1)
	require.NoError(t, err)

	expected, err := NewMerklePath("0001")
	require.NoError(t, err)

	require.Equal(t, expected, b)
}

func TestAddIndexOverflow(t *testing.T) {
	a, err := NewMerklePath("1111")
	require.NoError(t, err)

	_, err = a.Add(1)
	require.Error(t, err)
}

func TestAddUint32Overflow(t *testing.T) {
	a, err := NewMerklePath("11111111111111111111111111111111")
	require.NoError(t, err)

	_, err = a.Add(1)
	require.Error(t, err)
}

func TestAddToRoot(t *testing.T) {
	root, err := NewMerklePath("")
	require.NoError(t, err)

	_, err = root.Add(1)
	require.Error(t, err)
}

func TestSubOne(t *testing.T) {
	a, err := NewMerklePath("0010")
	require.NoError(t, err)

	b, err := a.Sub(1)
	require.NoError(t, err)

	expected, err := NewMerklePath("0001")
	require.NoError(t, err)

	require.Equal(t, expected, b)
}

func TestSubUnderflow(t *testing.T) {
	a, err := NewMerklePath("0001")
	require.NoError(t, err)

	_, err = a.Sub(1)
	require.Error(t, err)
}
