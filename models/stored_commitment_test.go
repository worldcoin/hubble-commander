package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestStoredCommitment_Bytes_Tx(t *testing.T) {
	commitment := &Commitment{
		ID: CommitmentID{
			BatchID:      MakeUint256(1),
			IndexInBatch: 4,
		},
		Type:              batchtype.Transfer,
		FeeReceiver:       3,
		CombinedSignature: Signature{1, 2, 3, 4, 5},
		PostStateRoot:     utils.RandomHash(),
		Transactions:      []byte{3, 2, 1},
	}

	storedCommitment := MakeStoredCommitmentFromTxCommitment(commitment)
	bytes := storedCommitment.Bytes()

	var decodedStoredCommitment StoredCommitment
	err := decodedStoredCommitment.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, storedCommitment, decodedStoredCommitment)

	decodedCommitment := decodedStoredCommitment.ToTxCommitment()
	require.Equal(t, *commitment, *decodedCommitment)
}

func TestStoredCommitment_Bytes_Deposit(t *testing.T) {
	commitment := &DepositCommitment{
		CommitmentBase: CommitmentBase{
			ID: CommitmentID{
				BatchID:      MakeUint256(1),
				IndexInBatch: 4,
			},
			Type:          batchtype.Deposit,
			PostStateRoot: utils.RandomHash(),
		},
		SubTreeID:   MakeUint256(5),
		SubTreeRoot: utils.RandomHash(),
		Deposits: []PendingDeposit{
			{
				ID: DepositID{
					BlockNumber: 10,
					LogIndex:    3,
				},
				ToPubKeyID: 5,
				TokenID:    MakeUint256(2),
				L2Amount:   MakeUint256(100),
			},
			{
				ID: DepositID{
					BlockNumber: 11,
					LogIndex:    2,
				},
				ToPubKeyID: 2,
				TokenID:    MakeUint256(3),
				L2Amount:   MakeUint256(50),
			},
		},
	}

	storedCommitment := MakeStoredCommitmentFromDepositCommitment(commitment)
	bytes := storedCommitment.Bytes()

	var decodedStoredCommitment StoredCommitment
	err := decodedStoredCommitment.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, storedCommitment, decodedStoredCommitment)

	decodedCommitment := decodedStoredCommitment.ToDepositCommitment()
	require.Equal(t, *commitment, *decodedCommitment)
}

func TestStoredCommitment_ToTxCommitment_InvalidType(t *testing.T) {
	commitment := MakeStoredCommitmentFromDepositCommitment(&DepositCommitment{})

	require.Panics(t, func() {
		commitment.ToTxCommitment()
	})
}

func TestStoredCommitment_ToDepositCommitment_InvalidType(t *testing.T) {
	commitment := MakeStoredCommitmentFromTxCommitment(&Commitment{})

	require.Panics(t, func() {
		commitment.ToDepositCommitment()
	})
}
