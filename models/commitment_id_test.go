package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommitmentID_Bytes(t *testing.T) {
	commitmentID := CommitmentID{
		BatchID:      MakeUint256(24),
		IndexInBatch: 4,
	}

	bytes := commitmentID.Bytes()

	var decodedCommitmentID CommitmentID
	err := decodedCommitmentID.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, commitmentID, decodedCommitmentID)
}
