package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSliceDiff_Unsorted(t *testing.T) {
	a := []string{"all", "ducks", "are", "awesome"}
	b := []string{"all", "dogs", "are", "cuddly"}

	diff := StringSliceDiff(a, b)
	require.Contains(t, diff, "dogs", "ducks", "awesome", "cuddly")
}

func TestStringSliceDiff_Sorted(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"a", "c", "d"}

	diff := StringSliceDiff(a, b)
	require.Contains(t, diff, "b")
}
