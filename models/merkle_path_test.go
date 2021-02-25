package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMerklePathLength(t *testing.T) {
	// 33 zeroes
	_, err := NewMerklePath("000000000000000000000000000000000")
	require.Error(t, err)
}

func TestValue(t *testing.T) {
	a, err := NewMerklePath("000111")
	require.NoError(t, err)

	b, err := a.Value()
	require.NoError(t, err)
	require.Equal(t, "000111", b)
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

func TestAddCreatesNewStruct(t *testing.T) {
	a, err := NewMerklePath("0000")
	require.NoError(t, err)

	b, err := a.Add(1)
	require.NoError(t, err)

	require.NotEqual(t, a, b)
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
	a, err := NewMerklePath("0000")
	require.NoError(t, err)

	_, err = a.Sub(1)
	require.Error(t, err)
}

func TestSubFromRoot(t *testing.T) {
	root, err := NewMerklePath("")
	require.NoError(t, err)

	_, err = root.Sub(1)
	require.Error(t, err)
}

func TestSubCreatesNewStruct(t *testing.T) {
	a, err := NewMerklePath("0010")
	require.NoError(t, err)

	b, err := a.Sub(1)
	require.NoError(t, err)

	require.NotEqual(t, a, b)
}

func TestSiblingOfLeft(t *testing.T) {
	a, err := NewMerklePath("0000")
	require.NoError(t, err)

	b, err := a.Sibling()
	require.NoError(t, err)

	expected, err := NewMerklePath("0001")
	require.NoError(t, err)

	require.Equal(t, expected, b)
}

func TestSiblingOfRight(t *testing.T) {
	a, err := NewMerklePath("0001")
	require.NoError(t, err)

	b, err := a.Sibling()
	require.NoError(t, err)

	expected, err := NewMerklePath("0000")
	require.NoError(t, err)

	require.Equal(t, expected, b)
}

func TestSiblingUint32(t *testing.T) {
	a, err := NewMerklePath("11111111111111111111111111111111")
	require.NoError(t, err)

	b, err := a.Sibling()
	require.NoError(t, err)

	expected, err := NewMerklePath("11111111111111111111111111111110")
	require.NoError(t, err)

	require.Equal(t, expected, b)
}

func TestSiblingOfRoot(t *testing.T) {
	root, err := NewMerklePath("")
	require.NoError(t, err)

	_, err = root.Sibling()
	require.Error(t, err)
}

func TestSiblingCreatesNewStruct(t *testing.T) {
	a, err := NewMerklePath("0010")
	require.NoError(t, err)

	b, err := a.Sibling()
	require.NoError(t, err)

	require.NotEqual(t, a, b)
}
