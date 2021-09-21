package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeposit_InvalidBytesLength(t *testing.T) {
	deposit := PendingDeposit{}
	err := deposit.SetBytes([]byte{1, 2, 3})
	require.ErrorIs(t, err, ErrInvalidLength)
}

func TestPendingDeposit_Bytes(t *testing.T) {
	deposit := PendingDeposit{
		ID: DepositID{
			BlockNumber: 4321,
			LogIndex:    63452,
		},
		ToPubKeyID: 16,
		TokenID:    MakeUint256(44),
		L2Amount:   MakeUint256(235),
	}

	bytes := deposit.Bytes()

	decodedDeposit := PendingDeposit{
		ToPubKeyID: 333,
		TokenID:    MakeUint256(222),
	}
	err := decodedDeposit.SetBytes(bytes)
	require.NoError(t, err)

	require.Equal(t, deposit, decodedDeposit)
}

func TestPendingDepositID_InvalidBytesLength(t *testing.T) {
	depositID := DepositID{}
	err := depositID.SetBytes([]byte{1, 2, 3})
	require.ErrorIs(t, err, ErrInvalidLength)
}

func TestPendingDepositID_Bytes(t *testing.T) {
	depositID := DepositID{
		BlockNumber: 68593,
		LogIndex:    49102,
	}

	bytes := depositID.Bytes()

	var decodedDepositID DepositID
	err := decodedDepositID.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, depositID, decodedDepositID)
}
