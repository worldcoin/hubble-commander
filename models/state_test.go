package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestStateUpdate_ByteEncoding(t *testing.T) {
	stateUpdate := StateUpdate{
		ID:          1,
		CurrentRoot: utils.RandomHash(),
		PrevRoot:    utils.RandomHash(),
		PrevStateLeaf: StateLeaf{
			StateID:  2,
			DataHash: utils.RandomHash(),
			UserState: UserState{
				PubKeyID: 3,
				TokenID:  MakeUint256(4),
				Balance:  MakeUint256(5),
				Nonce:    MakeUint256(6),
			},
		},
	}

	var decodedUpdate StateUpdate
	_ = decodedUpdate.SetBytes(stateUpdate.Bytes())
	require.Equal(t, stateUpdate, decodedUpdate)
}
