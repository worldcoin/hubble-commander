// +build e2e

package e2e

import (
	"errors"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestCommander(t *testing.T) {
	commander, err := NewCommanderFromEnv()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, commander.Stop())
	}()

	wallets, err := createWallets()
	require.NoError(t, err)

	var version string
	err = commander.Client().CallFor(&version, "hubble_getVersion")
	require.NoError(t, err)
	require.Equal(t, "dev-0.1.0", version)

	var userStates []dto.UserState
	err = commander.Client().CallFor(&userStates, "hubble_getUserStates", []interface{}{wallets[0].PublicKey()})
	require.NoError(t, err)
	require.Len(t, userStates, 1)
	require.EqualValues(t, models.MakeUint256(0), userStates[0].Nonce)

	transfer, err := api.SignTransfer(&wallets[1], dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	})
	require.NoError(t, err)

	var transferHash1 common.Hash
	err = commander.Client().CallFor(&transferHash1, "hubble_sendTransaction", []interface{}{*transfer})
	require.NoError(t, err)
	require.NotNil(t, transferHash1)

	var sentTransfer dto.TransferReceipt
	err = commander.Client().CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash1})
	require.NoError(t, err)
	require.Equal(t, txstatus.Pending, sentTransfer.Status)

	transfer2, err := api.SignTransfer(&wallets[1], dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(1),
	})
	require.NoError(t, err)

	var transferHash2 common.Hash
	err = commander.Client().CallFor(&transferHash2, "hubble_sendTransaction", []interface{}{*transfer2})
	require.NoError(t, err)
	require.NotNil(t, transferHash2)

	testutils.WaitToPass(func() bool {
		err = commander.Client().CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash1})
		require.NoError(t, err)
		return sentTransfer.Status == txstatus.InBatch
	}, 10*time.Second)

	err = commander.Client().CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash2})
	require.NoError(t, err)
	require.Equal(t, txstatus.InBatch, sentTransfer.Status)

	err = commander.Client().CallFor(&userStates, "hubble_getUserStates", []interface{}{wallets[1].PublicKey()})
	require.NoError(t, err)
	require.Len(t, userStates, 2)

	userState, err := getUserState(userStates, 1)
	require.NoError(t, err)
	require.EqualValues(t, models.MakeUint256(2), userState.Nonce)
	require.EqualValues(t, models.MakeUint256(920), userState.Balance)

	var batches []models.Batch
	err = commander.Client().CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})

	require.NoError(t, err)
	require.Len(t, batches, 1)
	require.Equal(t, txtype.Transfer, batches[0].Type)
	require.Equal(t, models.MakeUint256(1), batches[0].ID)
}

func getUserState(userStates []dto.UserState, stateID uint32) (*dto.UserState, error) {
	for i := range userStates {
		if userStates[i].StateID == stateID {
			return &userStates[i], nil
		}
	}
	return nil, errors.New("user state with given stateID not found")
}
