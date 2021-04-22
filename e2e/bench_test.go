// +build e2e

package e2e

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBenchCommander(t *testing.T) {
	commander, err := CreateCommanderFromEnv()
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

	stateIds := make([]uint32, 0)

	for _, wallet := range wallets {
		var userStates []dto.UserState
		err = commander.Client().CallFor(&userStates, "hubble_getUserStates", []interface{}{wallet.PublicKey()})
		for _, state := range userStates {
			stateIds = append(stateIds, state.StateID)
		}
	}

	transfersSent := 0
	startTime := time.Now()

	for _, wallet := range wallets {
		wallet := wallet
		go func() {
			for {
				to := stateIds[rand.Intn(len(stateIds))]
				sendTransfer(t, commander, wallet, to)

				transfersSent += 1
				fmt.Printf("Transfers sent: %d, throughput: %f tx/s\n", transfersSent, float64(transfersSent)/(time.Since(startTime).Seconds()))
			}
		}()
	}

	for {
		time.Sleep(time.Second)
	}
}

func sendTransfer(t *testing.T, commander E2ECommander, wallet bls.Wallet, to uint32) {
	var userStates []dto.UserState
	err := commander.Client().CallFor(&userStates, "hubble_getUserStates", []interface{}{wallet.PublicKey()})
	require.NoError(t, err)

	transfer, err := api.SignTransfer(&wallet, dto.Transfer{
		FromStateID: &userStates[0].StateID,
		ToStateID:   &to,
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       &userStates[0].Nonce,
	})
	require.NoError(t, err)

	var transferHash common.Hash
	err = commander.Client().CallFor(&transferHash, "hubble_sendTransaction", []interface{}{*transfer})
	require.NoError(t, err)
	require.NotNil(t, transferHash)

	var sentTransfer dto.TransferReceipt
	err = commander.Client().CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash})
	require.NoError(t, err)
	require.Equal(t, txstatus.Pending, sentTransfer.Status)

	testutils.WaitToPass(func() bool {
		err = commander.Client().CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash})
		require.NoError(t, err)
		return sentTransfer.Status == txstatus.InBatch
	}, 10*time.Second)
}
