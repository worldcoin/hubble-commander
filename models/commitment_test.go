package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestCommitment_Bytes(t *testing.T) {
	commitment := Commitment{
		BatchID:           MakeUint256(123),
		IndexInBatch:      4,
		Type:              txtype.Transfer,
		FeeReceiver:       11,
		CombinedSignature: MakeRandomSignature(),
		PostStateRoot:     common.Hash{1, 2, 3, 4},
		IncludedInBatch:   nil,
		Transactions:      []byte{5, 6, 1, 2, 6},
	}

	bytes := commitment.Bytes()

	decodedCommitment := Commitment{}
	err := decodedCommitment.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, commitment, decodedCommitment)
}

func TestCommitmentKey_Bytes(t *testing.T) {
	commitmentKey := CommitmentKey{
		BatchID:      MakeUint256(24),
		IndexInBatch: 4,
	}

	bytes := commitmentKey.Bytes()

	var decodedCommitmentKey CommitmentKey
	err := decodedCommitmentKey.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, commitmentKey, decodedCommitmentKey)
}
