// +build e2e

package e2e

import (
	"errors"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommanderUsingDocker(t *testing.T) {
	commander, err := StartCommander(StartOptions{
		Image: "ghcr.io/worldcoin/hubble-commander:latest",
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, commander.Stop())
	}()
	runE2ETest(t, commander.Client)
}

func runE2ETest(t *testing.T, client jsonrpc.RPCClient) {
	wallets, err := createWallets()
	require.NoError(t, err)

	var version string
	err = client.CallFor(&version, "hubble_getVersion")
	require.NoError(t, err)
	require.Equal(t, "dev-0.1.0", version)

	var userStates []dto.UserState
	err = client.CallFor(&userStates, "hubble_getUserStates", []interface{}{wallets[0].PublicKey()})
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
	err = client.CallFor(&transferHash1, "hubble_sendTransaction", []interface{}{*transfer})
	require.NoError(t, err)
	require.NotNil(t, transferHash1)

	var sentTransfer dto.TransferReceipt
	err = client.CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash1})
	require.NoError(t, err)
	require.Equal(t, models.Pending, sentTransfer.Status)

	transfer2, err := api.SignTransfer(&wallets[1], dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(1),
	})
	require.NoError(t, err)

	var transferHash2 common.Hash
	err = client.CallFor(&transferHash2, "hubble_sendTransaction", []interface{}{*transfer2})
	require.NoError(t, err)
	require.NotNil(t, transferHash2)

	testutils.WaitToPass(func() bool {
		err = client.CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash1})
		require.NoError(t, err)
		return sentTransfer.Status == models.InBatch
	}, 10*time.Second)

	err = client.CallFor(&sentTransfer, "hubble_getTransfer", []interface{}{transferHash2})
	require.NoError(t, err)
	require.Equal(t, models.InBatch, sentTransfer.Status)

	err = client.CallFor(&userStates, "hubble_getUserStates", []interface{}{wallets[1].PublicKey()})
	require.NoError(t, err)
	require.Len(t, userStates, 2)

	userState, err := getUserState(userStates, 1)
	require.NoError(t, err)
	require.EqualValues(t, models.MakeUint256(2), userState.Nonce)
	require.EqualValues(t, models.MakeUint256(920), userState.Balance)

	var batches []models.Batch
	err = client.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})

	require.NoError(t, err)
	require.Len(t, batches, 1)
}

func createWallets() ([]bls.Wallet, error) {
	cfg := config.GetConfig().Rollup
	accounts := cfg.GenesisAccounts

	wallets := make([]bls.Wallet, 0, len(accounts))
	for i := range accounts {
		wallet, err := bls.NewWallet(accounts[i].PrivateKey[:], cfg.SignaturesDomain)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, *wallet)
	}
	return wallets, nil
}

func getUserState(userStates []dto.UserState, stateID uint32) (*dto.UserState, error) {
	for i := range userStates {
		if userStates[i].StateID == stateID {
			return &userStates[i], nil
		}
	}
	return nil, errors.New("user state with given stateID not found")
}
