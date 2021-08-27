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
		Deposits: []DepositID{{
			BlockNumber: 1,
			LogIndex:    2,
		}},
	}

	newSubTree := PendingDepositSubTree{
		ID:   MakeUint256(1),
		Root: utils.RandomHash(),
	}

	initialCopy := initialSubTree
	err := initialCopy.SetBytes(newSubTree.Bytes())
	require.NoError(t, err)
	require.Equal(t, newSubTree, initialCopy)
}

func TestPendingDepositSubTree_Bytes(t *testing.T) {
	subTree := PendingDepositSubTree{
		Root: utils.RandomHash(),
		Deposits: []DepositID{
			{
				BlockNumber: 1,
				LogIndex:    2,
			},
			{
				BlockNumber: 3,
				LogIndex:    4,
			},
		},
	}

	bytes := subTree.Bytes()

	decodedSubTree := PendingDepositSubTree{
		Root: utils.RandomHash(),
		Deposits: []DepositID{
			{
				BlockNumber: 999,
				LogIndex:    999,
			},
		},
	}
	err := decodedSubTree.SetBytes(bytes)
	require.NoError(t, err)

	require.Equal(t, Uint256{}, decodedSubTree.ID)
	decodedSubTree.ID = subTree.ID
	require.Equal(t, subTree, decodedSubTree)
}
