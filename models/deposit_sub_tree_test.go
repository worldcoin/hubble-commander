package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestPendingDepositSubtree_SetBytes_InvalidBytesLength(t *testing.T) {
	subtree := PendingDepositSubtree{}
	err := subtree.SetBytes([]byte{1, 2, 3})
	require.ErrorIs(t, err, ErrInvalidLength)
}

func TestPendingDepositSubtree_SetBytes_UpdatesDeposits(t *testing.T) {
	initialSubtree := PendingDepositSubtree{
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

	newSubtree := PendingDepositSubtree{
		ID:   MakeUint256(1),
		Root: utils.RandomHash(),
	}

	err := initialSubtree.SetBytes(newSubtree.Bytes())
	require.NoError(t, err)
	require.Equal(t, newSubtree, initialSubtree)
}

func TestPendingDepositSubtree_Bytes(t *testing.T) {
	subtree := PendingDepositSubtree{
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

	bytes := subtree.Bytes()

	decodedSubtree := PendingDepositSubtree{
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
	err := decodedSubtree.SetBytes(bytes)
	require.NoError(t, err)

	require.Equal(t, MakeUint256(1), decodedSubtree.ID)
	decodedSubtree.ID = subtree.ID
	require.Equal(t, subtree, decodedSubtree)
}
