package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetLatestBlockNumber(t *testing.T) {
	storage := Storage{}
	currentBlockNumber := uint32(420)

	storage.SetLatestBlockNumber(currentBlockNumber)

	latestBlockNumber := storage.GetLatestBlockNumber()

	require.Equal(t, currentBlockNumber, latestBlockNumber)
	require.Equal(t, currentBlockNumber, storage.latestBlockNumber)
}
