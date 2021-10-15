package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestCommitment_Bytes(t *testing.T) {
	commitment := Commitment{
		CommitmentBase: CommitmentBase{
			ID: CommitmentID{
				BatchID:      MakeUint256(123),
				IndexInBatch: 4,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: common.Hash{1, 2, 3, 4},
		},
		FeeReceiver:       11,
		CombinedSignature: MakeRandomSignature(),
		Transactions:      []byte{5, 6, 1, 2, 6},
	}

	bytes := commitment.Bytes()

	decodedCommitment := Commitment{
		CommitmentBase: CommitmentBase{
			ID: CommitmentID{
				BatchID:      MakeUint256(123),
				IndexInBatch: 4,
			},
		},
	}
	err := decodedCommitment.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, commitment, decodedCommitment)
}

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
