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
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommander(t *testing.T) {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = 32
	cfg.Rollup.MaxTxsPerCommitment = 32
	cfg.Rollup.MinCommitmentsPerBatch = 1
	cfg.Rollup.MaxTxnDelay = 2 * time.Second

	commander, err := setup.NewConfiguredCommanderFromEnv(cfg, nil)
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

	submitTxBatchAndWait(t, commander.Client(), func() common.Hash {
		return testSubmitTransferBatch(t, commander.Client(), senderWallet, 0)
	})

	submitTxBatchAndWait(t, commander.Client(), func() common.Hash {
		firstC2TWallet := wallets[len(wallets)-32]
		return testSubmitC2TBatch(t, commander.Client(), senderWallet, wallets, firstC2TWallet.PublicKey(), 32)
	})

	submitTxBatchAndWait(t, commander.Client(), func() common.Hash {
		return testSubmitMassMigrationBatch(t, commander.Client(), senderWallet, 64)
	})

	testMaxBatchDelay(t, commander.Client(), senderWallet, 96)

	testSubmitDepositBatchAndWait(t, commander.Client(), 5)

	testSenderStateAfterTransfers(t, commander.Client(), senderWallet,
		32*3+1,
		32*100*3+100,
	)
	testFeeReceiverStateAfterTransfers(
		t, commander.Client(), feeReceiverWallet,
		32*10*3+10,
	)

	testGetBatches(t, commander.Client(), 6)

	testCommanderRestart(t, commander, senderWallet, 97)
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

func testSendMassMigration(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, nonce uint64) common.Hash {
	massMigration, err := api.SignMassMigration(&senderWallet, dto.MassMigration{
		FromStateID: ref.Uint32(1),
		SpokeID:     ref.Uint32(1),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(nonce),
	})
	require.NoError(t, err)

	var txHash common.Hash
	err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*massMigration})
	require.NoError(t, err)
	require.NotZero(t, txHash)
	return txHash
}

func send31MoreMassMigrations(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, startNonce uint64) {
	for nonce := startNonce; nonce < startNonce+31; nonce++ {
		transfer, err := api.SignMassMigration(&senderWallet, dto.MassMigration{
			FromStateID: ref.Uint32(1),
			SpokeID:     ref.Uint32(1),
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

func testGetTransaction(t *testing.T, client jsonrpc.RPCClient, txHash common.Hash) {
	var txReceipt dto.TransactionReceipt
	err := client.CallFor(&txReceipt, "hubble_getTransaction", []interface{}{txHash})
	require.NoError(t, err)
	require.Equal(t, txstatus.Pending, txReceipt.Status)
	require.Equal(t, txHash, txReceipt.Hash)
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
		var txReceipt dto.TransactionReceipt
		err := client.CallFor(&txReceipt, "hubble_getTransaction", []interface{}{txHash})
		require.NoError(t, err)
		return txReceipt.Status == txstatus.Mined
	}, 30*time.Second, testutils.TryInterval)
}

func testSenderStateAfterTransfers(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, expectedNonce, expectedBalance uint64) {
	var userStates []dto.UserStateWithID
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{senderWallet.PublicKey()})
	require.NoError(t, err)

	senderState, err := getUserState(userStates, 1)
	require.NoError(t, err)

	initialBalance := models.MakeUint256(setup.InitialGenesisBalance)
	require.Equal(t, models.MakeUint256(expectedNonce), senderState.Nonce)
	require.Equal(t, *initialBalance.SubN(expectedBalance), senderState.Balance)
}

func testFeeReceiverStateAfterTransfers(t *testing.T, client jsonrpc.RPCClient, feeReceiverWallet bls.Wallet, expectedBalance uint64) {
	var userStates []dto.UserStateWithID
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{feeReceiverWallet.PublicKey()})
	require.NoError(t, err)

	feeReceiverState, err := getUserState(userStates, 0)
	require.NoError(t, err)

	initialBalance := models.MakeUint256(setup.InitialGenesisBalance)
	require.Equal(t, *initialBalance.AddN(expectedBalance), feeReceiverState.Balance)
	require.Equal(t, models.MakeUint256(0), feeReceiverState.Nonce)
}

func testGetBatches(t *testing.T, client jsonrpc.RPCClient, expectedBatchCount int) {
	var batches []dto.Batch
	err := client.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})
	require.NoError(t, err)
	require.Len(t, batches, expectedBatchCount)
	require.Equal(t, models.MakeUint256(1), batches[1].ID)

	batchTypes := map[batchtype.BatchType]bool{}
	for _, batch := range batches {
		batchTypes[batch.Type] = true
	}
	require.Contains(t, batchTypes, batchtype.Transfer)
	require.Contains(t, batchTypes, batchtype.Create2Transfer)
	require.Contains(t, batchTypes, batchtype.MassMigration)
	require.Contains(t, batchTypes, batchtype.Deposit)
}

func testCommanderRestart(t *testing.T, commander setup.Commander, senderWallet bls.Wallet, startNonce uint64) {
	err := commander.Restart()
	require.NoError(t, err)

	testSendTransfer(t, commander.Client(), senderWallet, startNonce)
}

func getUserState(userStates []dto.UserStateWithID, stateID uint32) (*dto.UserStateWithID, error) {
	for i := range userStates {
		if userStates[i].StateID == stateID {
			return &userStates[i], nil
		}
	}
	return nil, errors.New("user state with given stateID not found")
}

// confirms that batches smaller than the minimum will be submitted if any txn is left
// pending for too long
func testMaxBatchDelay(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, startNonce uint64) {
	txnHash := testSendTransfer(t, client, senderWallet, startNonce)
	require.NotZero(t, txnHash)

	time.Sleep(1 * time.Second)

	var txReceipt dto.TransactionReceipt
	err := client.CallFor(&txReceipt, "hubble_getTransaction", []interface{}{txnHash})
	require.NoError(t, err)
	require.NotEqual(t, txReceipt.Status, txstatus.Mined)

	log.Warn("txn is not yet in batch. waiting..")

	require.Eventually(t, func() bool {
		err = client.CallFor(&txReceipt, "hubble_getTransaction", []interface{}{txnHash})
		require.NoError(t, err)
		return txReceipt.Status == txstatus.Mined
	}, 10*time.Second, testutils.TryInterval)
}
