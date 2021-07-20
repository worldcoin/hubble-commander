package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMerklePathLength(t *testing.T) {
	_, err := NewMerklePath(strings.Repeat("0", 33))
	require.Error(t, err)
}

func TestValue(t *testing.T) {
	a, err := NewMerklePath("000111")
	require.NoError(t, err)

	b, err := a.Value()
	require.NoError(t, err)
	require.Equal(t, "0000111", b)
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
	a, err := NewMerklePath(strings.Repeat("1", 32))
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

func TestGetWitnessesOfRoot(t *testing.T) {
	path, err := NewMerklePath("")
	require.NoError(t, err)

	_, err = path.GetWitnessPaths()
	require.Error(t, err)
}

func TestGetWitnessesOfDepth1(t *testing.T) {
	path, err := NewMerklePath("1")
	require.NoError(t, err)

	witnesses, err := path.GetWitnessPaths()
	require.NoError(t, err)

	p, err := NewMerklePath("0")
	require.NoError(t, err)
	expected := []MerklePath{*p}

	require.Equal(t, expected, witnesses)
}

func TestGetWitnessesOfDepth3(t *testing.T) {
	path, err := NewMerklePath("101")
	require.NoError(t, err)

	witnesses, err := path.GetWitnessPaths()
	require.NoError(t, err)

	expected := make([]MerklePath, 0, 3)
	expectedPaths := []string{"100", "11", "0"}
	for _, paths := range expectedPaths {
		p, err := NewMerklePath(paths)
		require.NoError(t, err)
		expected = append(expected, *p)
	}

	require.Equal(t, expected, witnesses)
}

func TestNamespaceBytesSetBytes(t *testing.T) {
	path, err := NewMerklePath("101")
	require.NoError(t, err)
	ns := NamespacedMerklePath{Namespace: "foo", Path: *path}

	bytes := ns.Bytes()
	require.Len(t, bytes, 8)

	parsed := NamespacedMerklePath{}
	err = parsed.SetBytes(bytes)
	require.NoError(t, err)

	require.Equal(t, ns, parsed)
}
