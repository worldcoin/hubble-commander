package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestStoredCommitment_Bytes_TxCommitment(t *testing.T) {
	commitment := &TxCommitment{
		CommitmentBase: CommitmentBase{
			ID: CommitmentID{
				BatchID:      MakeUint256(1),
				IndexInBatch: 4,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       3,
		CombinedSignature: Signature{1, 2, 3, 4, 5},
		BodyHash:          utils.NewRandomHash(),
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

func TestStoredCommitment_Bytes_DepositCommitment(t *testing.T) {
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
					SubtreeID:    MakeUint256(32),
					DepositIndex: MakeUint256(1),
				},
				ToPubKeyID: 5,
				TokenID:    MakeUint256(2),
				L2Amount:   MakeUint256(100),
			},
			{
				ID: DepositID{
					SubtreeID:    MakeUint256(11),
					DepositIndex: MakeUint256(2),
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

func TestStoredCommitment_Bytes_DepositCommitmentWithoutPendingDeposits(t *testing.T) {
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
		Deposits:    []PendingDeposit{},
	}

	storedCommitment := MakeStoredCommitmentFromDepositCommitment(commitment)
	bytes := storedCommitment.Bytes()

	var decodedStoredCommitment StoredCommitment
	err := decodedStoredCommitment.SetBytes(bytes)
	require.ErrorIs(t, err, ErrInvalidLength)
}

func TestStoredCommitment_ToTxCommitment_InvalidType(t *testing.T) {
	commitment := MakeStoredCommitmentFromDepositCommitment(&DepositCommitment{})

	require.Panics(t, func() {
		commitment.ToTxCommitment()
	})
}

func TestStoredCommitment_ToDepositCommitment_InvalidType(t *testing.T) {
	commitment := MakeStoredCommitmentFromTxCommitment(&TxCommitment{})

	require.Panics(t, func() {
		commitment.ToDepositCommitment()
	})
}
