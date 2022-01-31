package db

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"
)

func TestNewSeekPrefix_CopiesToNewMemoryLocationToAvoidRaceConditions(t *testing.T) {
	prefix := make([]byte, 3, 4)
	prefix[0] = 1
	prefix[1] = 2
	prefix[2] = 3

	underlyingArrayAddress := &prefix[0]

	newPrefix := newSeekPrefix(prefix, badger.IteratorOptions{
		Reverse: true,
	})

	require.Same(t, underlyingArrayAddress, &prefix[0])
	require.NotSame(t, underlyingArrayAddress, &newPrefix[0])
}
