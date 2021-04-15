// +build e2e

package e2e

import (
	"errors"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_Commander(t *testing.T) {
	commander, err := StartCommander(StartOptions{
		Image: "ghcr.io/worldcoin/hubble-commander:latest",
	})
	require.NoError(t, err)
	defer func() {
		err = commander.Stop()
		require.NoError(t, err)
	}()

	var version string
	err = commander.Client.CallFor(&version, "hubble_getVersion")
	require.NoError(t, err)
	require.Equal(t, "dev-0.1.0", version)

	var userStates []dto.UserState
	err = commander.Client.CallFor(&userStates, "hubble_getUserStates", []interface{}{models.PublicKey{1, 2, 3}})
	require.NoError(t, err)
	require.Len(t, userStates, 1)
	require.EqualValues(t, models.MakeUint256(0), userStates[0].Nonce)

	transfer := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}

	var transferHash1 common.Hash
	err = commander.Client.CallFor(&transferHash1, "hubble_sendTransaction", []interface{}{transfer})
	require.NoError(t, err)
	require.NotNil(t, transferHash1)

	var sentTransfer models.TransferReceipt
	err = commander.Client.CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash1})
	require.NoError(t, err)
	require.Equal(t, models.Pending, sentTransfer.Status)

	transfer2 := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(1),
		Signature:   []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}

	var transferHash2 common.Hash
	err = commander.Client.CallFor(&transferHash2, "hubble_sendTransaction", []interface{}{transfer2})
	require.NoError(t, err)
	require.NotNil(t, transferHash2)

	testutils.WaitToPass(func() bool {
		err = commander.Client.CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash1})
		require.NoError(t, err)
		return sentTransfer.Status == models.InBatch
	}, 10*time.Second)

	err = commander.Client.CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash2})
	require.NoError(t, err)
	require.Equal(t, models.InBatch, sentTransfer.Status)

	err = commander.Client.CallFor(&userStates, "hubble_getUserStates", []interface{}{models.PublicKey{2, 3, 4}})
	require.NoError(t, err)
	require.Len(t, userStates, 2)

	userState, err := getUserState(userStates, 1)
	require.NoError(t, err)
	require.EqualValues(t, models.MakeUint256(2), userState.Nonce)
	require.EqualValues(t, models.MakeUint256(920), userState.Balance)
}

func getUserState(userStates []dto.UserState, stateID uint) (*dto.UserState, error) {
	for i := range userStates {
		if userStates[i].StateID == 1 {
			return &userStates[i], nil
		}
	}
	return nil, errors.New("user state with given stateID not found")
}
