package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeCommitmentIDPointer(t *testing.T) {
	id := &CommitmentID{
		BatchID:      MakeUint256(5),
		IndexInBatch: 2,
	}
	bytes := EncodeCommitmentIDPointer(id)
	require.EqualValues(t, 1, bytes[0])

	decodedID, err := DecodeCommitmentIDPointer(bytes)
	require.NoError(t, err)
	require.Equal(t, *id, *decodedID)
}

func TestEncodeCommitmentIDPointer_NilValue(t *testing.T) {
	var id *CommitmentID
	bytes := EncodeCommitmentIDPointer(id)
	require.EqualValues(t, 0, bytes[0])

	decodedID, err := DecodeCommitmentIDPointer(bytes)
	require.NoError(t, err)
	require.Nil(t, decodedID)
}
