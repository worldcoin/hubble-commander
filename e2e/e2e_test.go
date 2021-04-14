// +build e2e

package e2e

import (
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

	tx := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}

	var txHash1 common.Hash
	err = commander.Client.CallFor(&txHash1, "hubble_sendTransaction", []interface{}{tx})
	require.NoError(t, err)
	require.NotNil(t, txHash1)

	var sentTx models.TransactionReceipt
	err = commander.Client.CallFor(&sentTx, "hubble_getTransaction", []interface{}{txHash1})
	require.NoError(t, err)
	require.Equal(t, models.Pending, sentTx.Status)

	tx2 := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(1),
		Signature:   []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}

	var txHash2 common.Hash
	err = commander.Client.CallFor(&txHash2, "hubble_sendTransaction", []interface{}{tx2})
	require.NoError(t, err)
	require.NotNil(t, txHash2)

	testutils.WaitToPass(func() bool {
		err = commander.Client.CallFor(&sentTx, "hubble_getTransaction", []interface{}{txHash1})
		require.NoError(t, err)
		return sentTx.Status == models.InBatch
	}, 10*time.Second)

	err = commander.Client.CallFor(&sentTx, "hubble_getTransaction", []interface{}{txHash2})
	require.NoError(t, err)
	require.Equal(t, models.InBatch, sentTx.Status)

	err = commander.Client.CallFor(&userStates, "hubble_getUserStates", []interface{}{models.PublicKey{2, 3, 4}})
	require.NoError(t, err)
	require.Len(t, userStates, 2)
	require.EqualValues(t, models.MakeUint256(2), userStates[0].Nonce)
	require.EqualValues(t, models.MakeUint256(920), userStates[0].Balance)
}
