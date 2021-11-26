//go:build e2e
// +build e2e

package e2e

import (
	"errors"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommander(t *testing.T) {
	cfg := config.GetConfig().Rollup
	cfg.MinTxsPerCommitment = 32
	cfg.MaxTxsPerCommitment = 32
	cfg.MinCommitmentsPerBatch = 1

	commander, err := setup.NewConfiguredCommanderFromEnv(cfg)
	require.NoError(t, err)
	err = commander.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, commander.Stop())
	}()

	domain := GetDomain(t, commander.Client())

	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	feeReceiverWallet := wallets[0]
	senderWallet := wallets[1]

	testGetVersion(t, commander.Client())
	firstUserState := testGetUserStates(t, commander.Client(), senderWallet)
	testGetPublicKey(t, commander.Client(), &firstUserState, senderWallet)
	testSubmitTransferBatch(t, commander.Client(), senderWallet, 0)

	firstC2TWallet := wallets[len(wallets)-32]
	testSubmitC2TBatch(t, commander.Client(), senderWallet, wallets, firstC2TWallet.PublicKey(), 32)

	testSubmitDepositBatch(t, commander.Client())

	testSenderStateAfterTransfers(t, commander.Client(), senderWallet)
	testFeeReceiverStateAfterTransfers(t, commander.Client(), feeReceiverWallet)
	testGetBatches(t, commander.Client())

	testCommanderRestart(t, commander, senderWallet)
}

func testGetVersion(t *testing.T, client jsonrpc.RPCClient) {
	var version string
	err := client.CallFor(&version, "hubble_getVersion")
	require.NoError(t, err)
	require.Equal(t, config.GetConfig().API.Version, version)
}

func testGetUserStates(t *testing.T, client jsonrpc.RPCClient, wallet bls.Wallet) dto.UserStateWithID {
	var userStates []dto.UserStateWithID
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{wallet.PublicKey()})
	require.NoError(t, err)
	require.Len(t, userStates, 2)
	require.EqualValues(t, 1, userStates[0].StateID)
	require.Equal(t, models.MakeUint256(0), userStates[0].Nonce)
	require.EqualValues(t, 3, userStates[1].StateID)
	require.Equal(t, models.MakeUint256(0), userStates[1].Nonce)
	return userStates[0]
}

func testGetPublicKey(t *testing.T, client jsonrpc.RPCClient, state *dto.UserStateWithID, wallet bls.Wallet) {
	var publicKey models.PublicKey
	err := client.CallFor(&publicKey, "hubble_getPublicKeyByPubKeyID", []interface{}{state.PubKeyID})
	require.NoError(t, err)
	require.Equal(t, *wallet.PublicKey(), publicKey)

	err = client.CallFor(&publicKey, "hubble_getPublicKeyByStateID", []interface{}{state.StateID})
	require.NoError(t, err)
	require.Equal(t, *wallet.PublicKey(), publicKey)
}

func testSendTransfer(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, nonce uint64) common.Hash {
	transfer, err := api.SignTransfer(&senderWallet, dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(nonce),
	})
	require.NoError(t, err)

	var txHash common.Hash
	err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
	require.NoError(t, err)
	require.NotZero(t, txHash)
	return txHash
}

func testSendCreate2Transfer(
	t *testing.T,
	client jsonrpc.RPCClient,
	senderWallet bls.Wallet,
	targetPublicKey *models.PublicKey,
	nonce uint64,
) common.Hash {
	transfer, err := api.SignCreate2Transfer(&senderWallet, dto.Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToPublicKey: targetPublicKey,
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(nonce),
	})
	require.NoError(t, err)

	var txHash common.Hash
	err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
	require.NoError(t, err)
	require.NotZero(t, txHash)
	return txHash
}

func testGetTransaction(t *testing.T, client jsonrpc.RPCClient, txHash common.Hash) {
	var txReceipt dto.TransferReceipt
	err := client.CallFor(&txReceipt, "hubble_getTransaction", []interface{}{txHash})
	require.NoError(t, err)
	require.Equal(t, txHash, txReceipt.Hash)
	require.Equal(t, txstatus.Pending, txReceipt.Status)
}

func send31MoreTransfers(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, startNonce uint64) {
	for nonce := startNonce; nonce < startNonce+31; nonce++ {
		transfer, err := api.SignTransfer(&senderWallet, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(nonce),
		})
		require.NoError(t, err)

		var txHash common.Hash
		err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
		require.NoError(t, err)
		require.NotZero(t, txHash)
	}
}

func send31MoreCreate2Transfers(
	t *testing.T,
	client jsonrpc.RPCClient,
	senderWallet bls.Wallet,
	wallets []bls.Wallet,
	startNonce uint64,
) {
	walletIndex := len(wallets) - 31
	for nonce := startNonce; nonce < startNonce+31; nonce++ {
		receiverWallet := wallets[walletIndex]
		transfer, err := api.SignCreate2Transfer(&senderWallet, dto.Create2Transfer{
			FromStateID: ref.Uint32(1),
			ToPublicKey: receiverWallet.PublicKey(),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(nonce),
		})
		require.NoError(t, err)

		var txHash common.Hash
		err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
		require.NoError(t, err)
		require.NotZero(t, txHash)

		walletIndex++
	}
}

func waitForTxToBeIncludedInBatch(t *testing.T, client jsonrpc.RPCClient, txHash common.Hash) {
	require.Eventually(t, func() bool {
		var txReceipt dto.TransferReceipt
		err := client.CallFor(&txReceipt, "hubble_getTransaction", []interface{}{txHash})
		require.NoError(t, err)
		return txReceipt.Status == txstatus.InBatch
	}, 30*time.Second, testutils.TryInterval)
}

func testSenderStateAfterTransfers(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet) {
	var userStates []dto.UserStateWithID
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{senderWallet.PublicKey()})
	require.NoError(t, err)

	senderState, err := getUserState(userStates, 1)
	require.NoError(t, err)

	initialBalance := config.GetDeployerConfig().Bootstrap.GenesisAccounts[1].Balance
	require.Equal(t, models.MakeUint256(32+32), senderState.Nonce)
	require.Equal(t, *initialBalance.SubN(32*100 + 32*100), senderState.Balance)
}

func testFeeReceiverStateAfterTransfers(t *testing.T, client jsonrpc.RPCClient, feeReceiverWallet bls.Wallet) {
	var userStates []dto.UserStateWithID
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{feeReceiverWallet.PublicKey()})
	require.NoError(t, err)

	feeReceiverState, err := getUserState(userStates, 0)
	require.NoError(t, err)

	initialBalance := config.GetDeployerConfig().Bootstrap.GenesisAccounts[1].Balance
	require.Equal(t, *initialBalance.AddN(32*10 + 32*10), feeReceiverState.Balance)
	require.Equal(t, models.MakeUint256(0), feeReceiverState.Nonce)
}

func testGetBatches(t *testing.T, client jsonrpc.RPCClient) {
	var batches []dto.Batch
	err := client.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})

	require.NoError(t, err)
	require.Len(t, batches, 4)
	require.Equal(t, models.MakeUint256(1), batches[1].ID)
	batchTypes := []batchtype.BatchType{batches[1].Type, batches[2].Type, batches[3].Type}
	require.Contains(t, batchTypes, batchtype.Transfer)
	require.Contains(t, batchTypes, batchtype.Create2Transfer)
	require.Contains(t, batchTypes, batchtype.Deposit)
}

func testCommanderRestart(t *testing.T, commander setup.Commander, senderWallet bls.Wallet) {
	err := commander.Restart()
	require.NoError(t, err)

	testSendTransfer(t, commander.Client(), senderWallet, 64)
}

func getUserState(userStates []dto.UserStateWithID, stateID uint32) (*dto.UserStateWithID, error) {
	for i := range userStates {
		if userStates[i].StateID == stateID {
			return &userStates[i], nil
		}
	}
	return nil, errors.New("user state with given stateID not found")
}
