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
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommander(t *testing.T) {
	commander, err := NewCommanderFromEnv(true)
	require.NoError(t, err)
	err = commander.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, commander.Stop())
	}()

	domain := getDomain(t, commander.Client())

	wallets, err := createWallets(domain)
	require.NoError(t, err)

	feeReceiverWallet := wallets[0]
	senderWallet := wallets[1]

	testGetVersion(t, commander.Client())
	testGetUserStates(t, commander.Client(), senderWallet)
	firstTransferHash := testSendTransfer(t, commander.Client(), senderWallet, models.NewUint256(0))
	testGetTransaction(t, commander.Client(), firstTransferHash)
	send31MoreTransfers(t, commander.Client(), senderWallet)

	firstC2TWallet := wallets[len(wallets)-32]
	firstCreate2TransferHash := testSendCreate2Transfer(t, commander.Client(), senderWallet, *firstC2TWallet.PublicKey())
	testGetTransaction(t, commander.Client(), firstCreate2TransferHash)
	send31MoreCreate2Transfers(t, commander.Client(), senderWallet, wallets)

	waitForTxToBeIncludedInBatch(t, commander.Client(), firstTransferHash)
	waitForTxToBeIncludedInBatch(t, commander.Client(), firstCreate2TransferHash)

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

func testGetUserStates(t *testing.T, client jsonrpc.RPCClient, wallet bls.Wallet) {
	var userStates []dto.UserState
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{wallet.PublicKey()})
	require.NoError(t, err)
	require.Len(t, userStates, 2)
	require.EqualValues(t, 1, userStates[0].StateID)
	require.Equal(t, models.MakeUint256(0), userStates[0].Nonce)
	require.EqualValues(t, 3, userStates[1].StateID)
	require.Equal(t, models.MakeUint256(0), userStates[1].Nonce)
}

func testSendTransfer(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, nonce *models.Uint256) common.Hash {
	transfer, err := api.SignTransfer(&senderWallet, dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       nonce,
	})
	require.NoError(t, err)

	var txHash common.Hash
	err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
	require.NoError(t, err)
	require.NotZero(t, txHash)
	return txHash
}

func testSendCreate2Transfer(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, targetPublicKey models.PublicKey) common.Hash {
	transfer, err := api.SignCreate2Transfer(&senderWallet, dto.Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToPublicKey: &targetPublicKey,
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(32),
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

func send31MoreTransfers(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet) {
	for nonce := uint64(1); nonce < 32; nonce++ {
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

func send31MoreCreate2Transfers(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, wallets []bls.Wallet) {
	for nonce := 1; nonce < 32; nonce++ {
		receiverWallet := wallets[len(wallets)-32+nonce]
		transfer, err := api.SignCreate2Transfer(&senderWallet, dto.Create2Transfer{
			FromStateID: ref.Uint32(1),
			ToPublicKey: receiverWallet.PublicKey(),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(uint64(32) + uint64(nonce)),
		})
		require.NoError(t, err)

		var txHash common.Hash
		err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*transfer})
		require.NoError(t, err)
		require.NotZero(t, txHash)
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
	var userStates []dto.UserState
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{senderWallet.PublicKey()})
	require.NoError(t, err)

	senderState, err := getUserState(userStates, 1)
	require.NoError(t, err)

	initialBalance := config.GetConfig().Rollup.GenesisAccounts[1].Balance
	require.Equal(t, models.MakeUint256(32+32), senderState.Nonce)
	require.Equal(t, *initialBalance.SubN(32*100 + 32*100), senderState.Balance)
}

func testFeeReceiverStateAfterTransfers(t *testing.T, client jsonrpc.RPCClient, feeReceiverWallet bls.Wallet) {
	var userStates []dto.UserState
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{feeReceiverWallet.PublicKey()})
	require.NoError(t, err)

	feeReceiverState, err := getUserState(userStates, 0)
	require.NoError(t, err)

	initialBalance := config.GetConfig().Rollup.GenesisAccounts[1].Balance
	require.Equal(t, *initialBalance.AddN(32*10 + 32*10), feeReceiverState.Balance)
	require.Equal(t, models.MakeUint256(0), feeReceiverState.Nonce)
}

func testGetBatches(t *testing.T, client jsonrpc.RPCClient) {
	var batches []models.Batch
	err := client.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})

	require.NoError(t, err)
	require.Len(t, batches, 2)
	require.Equal(t, models.MakeUint256(1), batches[0].ID)
	batchTypes := []txtype.TransactionType{batches[0].Type, batches[1].Type}
	require.Contains(t, batchTypes, txtype.Transfer)
	require.Contains(t, batchTypes, txtype.Create2Transfer)
}

func testCommanderRestart(t *testing.T, commander Commander, senderWallet bls.Wallet) {
	err := commander.Restart()
	require.NoError(t, err)

	testSendTransfer(t, commander.Client(), senderWallet, models.NewUint256(64))
}

func getDomain(t *testing.T, client jsonrpc.RPCClient) bls.Domain {
	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	require.NoError(t, err)

	domain, err := bls.DomainFromBytes(crypto.Keccak256(info.Rollup.Bytes()))
	require.NoError(t, err)
	return *domain
}

func getUserState(userStates []dto.UserState, stateID uint32) (*dto.UserState, error) {
	for i := range userStates {
		if userStates[i].StateID == stateID {
			return &userStates[i], nil
		}
	}
	return nil, errors.New("user state with given stateID not found")
}
