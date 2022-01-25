package stored

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestCommitment_Bytes_TxCommitment(t *testing.T) {
	commitment := &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 4,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       3,
		CombinedSignature: models.Signature{1, 2, 3, 4, 5},
		BodyHash:          utils.NewRandomHash(),
	}

	storedCommitment := MakeCommitmentFromTxCommitment(commitment)
	bytes := storedCommitment.Bytes()

	var decodedStoredCommitment Commitment
	err := decodedStoredCommitment.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, storedCommitment, decodedStoredCommitment)

	decodedCommitment := decodedStoredCommitment.ToTxCommitment()
	require.Equal(t, *commitment, *decodedCommitment)
}

func TestCommitment_Bytes_MMCommitment(t *testing.T) {
	commitment := &models.MMCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 4,
			},
			Type:          batchtype.MassMigration,
			PostStateRoot: utils.RandomHash(),
		},
		CombinedSignature: models.Signature{1, 2, 3, 4, 5},
		BodyHash:          utils.NewRandomHash(),
		Meta: &models.MassMigrationMeta{
			SpokeID:     5,
			TokenID:     models.MakeUint256(6),
			Amount:      models.MakeUint256(7),
			FeeReceiver: 8,
		},
		WithdrawRoot: utils.RandomHash(),
	}

	storedCommitment := MakeCommitmentFromMMCommitment(commitment)
	bytes := storedCommitment.Bytes()

	var decodedStoredCommitment Commitment
	err := decodedStoredCommitment.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, storedCommitment, decodedStoredCommitment)

	decodedCommitment := decodedStoredCommitment.ToMMCommitment()
	require.Equal(t, *commitment, *decodedCommitment)
}

func TestCommitment_Bytes_DepositCommitment(t *testing.T) {
	commitment := &models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 4,
			},
			Type:          batchtype.Deposit,
			PostStateRoot: utils.RandomHash(),
		},
		SubtreeID:   models.MakeUint256(5),
		SubtreeRoot: utils.RandomHash(),
		Deposits: []models.PendingDeposit{
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(32),
					DepositIndex: models.MakeUint256(1),
				},
				ToPubKeyID: 5,
				TokenID:    models.MakeUint256(2),
				L2Amount:   models.MakeUint256(100),
			},
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(11),
					DepositIndex: models.MakeUint256(2),
				},
				ToPubKeyID: 2,
				TokenID:    models.MakeUint256(3),
				L2Amount:   models.MakeUint256(50),
			},
		},
	}

	storedCommitment := MakeCommitmentFromDepositCommitment(commitment)
	bytes := storedCommitment.Bytes()

	var decodedStoredCommitment Commitment
	err := decodedStoredCommitment.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, storedCommitment, decodedStoredCommitment)

	decodedCommitment := decodedStoredCommitment.ToDepositCommitment()
	require.Equal(t, *commitment, *decodedCommitment)
}

func TestCommitment_Bytes_DepositCommitmentWithoutPendingDeposits(t *testing.T) {
	commitment := &models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 4,
			},
			Type:          batchtype.Deposit,
			PostStateRoot: utils.RandomHash(),
		},
		SubtreeID:   models.MakeUint256(5),
		SubtreeRoot: utils.RandomHash(),
		Deposits:    []models.PendingDeposit{},
	}

	storedCommitment := MakeCommitmentFromDepositCommitment(commitment)
	bytes := storedCommitment.Bytes()

	var decodedStoredCommitment Commitment
	err := decodedStoredCommitment.SetBytes(bytes)
	require.ErrorIs(t, err, models.ErrInvalidLength)
}

func TestCommitment_ToTxCommitment_InvalidType(t *testing.T) {
	commitment := MakeCommitmentFromDepositCommitment(&models.DepositCommitment{})

	require.Panics(t, func() {
		commitment.ToTxCommitment()
	})
}

func TestCommitment_ToDepositCommitment_InvalidType(t *testing.T) {
	commitment := MakeCommitmentFromTxCommitment(&models.TxCommitment{})

	require.Panics(t, func() {
		commitment.ToDepositCommitment()
	})
}
