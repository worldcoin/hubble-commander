package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestPendingDepositSubTree_SetBytes_InvalidBytesLength(t *testing.T) {
	subTree := PendingDepositSubTree{}
	err := subTree.SetBytes([]byte{1, 2, 3})
	require.ErrorIs(t, err, ErrInvalidLength)
}

func TestPendingDepositSubTree_SetBytes_UpdatesDeposits(t *testing.T) {
	initialSubTree := PendingDepositSubTree{
		ID:   MakeUint256(1),
		Root: utils.RandomHash(),
		Deposits: []PendingDeposit{
			{
				ID: DepositID{
					SubtreeID:    MakeUint256(1),
					DepositIndex: MakeUint256(0),
				},
				ToPubKeyID: 3,
				TokenID:    MakeUint256(4),
				L2Amount:   MakeUint256(500),
			},
		},
	}

	newSubTree := PendingDepositSubTree{
		ID:   MakeUint256(1),
		Root: utils.RandomHash(),
	}

	err := initialSubTree.SetBytes(newSubTree.Bytes())
	require.NoError(t, err)
	require.Equal(t, newSubTree, initialSubTree)
}

func TestPendingDepositSubTree_Bytes(t *testing.T) {
	subTree := PendingDepositSubTree{
		Root: utils.RandomHash(),
		Deposits: []PendingDeposit{
			{
				ID: DepositID{
					SubtreeID:    MakeUint256(1),
					DepositIndex: MakeUint256(2),
				},
				ToPubKeyID: 3,
				TokenID:    MakeUint256(4),
				L2Amount:   MakeUint256(500),
			},
			{
				ID: DepositID{
					SubtreeID:    MakeUint256(6),
					DepositIndex: MakeUint256(7),
				},
				ToPubKeyID: 8,
				TokenID:    MakeUint256(9),
				L2Amount:   MakeUint256(1000),
			},
		},
	}

	bytes := subTree.Bytes()

	decodedSubTree := PendingDepositSubTree{
		ID:   MakeUint256(1),
		Root: utils.RandomHash(),
		Deposits: []PendingDeposit{
			{
				ID: DepositID{
					SubtreeID:    MakeUint256(999),
					DepositIndex: MakeUint256(999),
				},
				ToPubKeyID: 999,
				TokenID:    MakeUint256(999),
				L2Amount:   MakeUint256(999),
			},
		},
	}
	err := decodedSubTree.SetBytes(bytes)
	require.NoError(t, err)

	require.Equal(t, MakeUint256(1), decodedSubTree.ID)
	decodedSubTree.ID = subTree.ID
	require.Equal(t, subTree, decodedSubTree)
}
