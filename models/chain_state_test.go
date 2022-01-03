package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestChainState_Bytes_ReturnsACopy(t *testing.T) {
	chainState := ChainState{
		ChainID: MakeUint256(1337),
	}
	bytes := chainState.Bytes()
	bytes[0] = 9
	require.Equal(t, MakeUint256(1337), chainState.ChainID)
}

func TestChainState_SetBytes(t *testing.T) {
	chainState := ChainState{
		ChainID:                        MakeUint256(1337),
		AccountRegistry:                utils.RandomAddress(),
		AccountRegistryDeploymentBlock: 7392,
		TokenRegistry:                  utils.RandomAddress(),
		SpokeRegistry:                  utils.RandomAddress(),
		DepositManager:                 utils.RandomAddress(),
		Rollup:                         utils.RandomAddress(),
		SyncedBlock:                    8001,
		GenesisAccounts: GenesisAccounts{
			{
				PublicKey: PublicKey{1, 2, 0, 5, 4},
				PubKeyID:  7,
				StateID:   44,
				Balance:   MakeUint256(4314),
			},
			{
				PublicKey: PublicKey{3, 2, 1, 1},
				PubKeyID:  83,
				StateID:   99,
				Balance:   MakeUint256(173212),
			},
		},
	}
	bytes := chainState.Bytes()
	newChainState := ChainState{}
	err := newChainState.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, chainState, newChainState)
}

func TestChainState_SetBytes_InvalidLength(t *testing.T) {
	chainState := ChainState{}

	data40 := make([]byte, 40)
	err := chainState.SetBytes(data40)
	require.ErrorIs(t, err, ErrInvalidLength)

	data257 := make([]byte, 257)
	err = chainState.SetBytes(data257)
	require.ErrorIs(t, err, ErrInvalidLength)
}
