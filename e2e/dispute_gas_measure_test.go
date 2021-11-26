//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestMeasureDisputeGasUsage(t *testing.T) {
	cfg := config.GetConfig().Rollup
	cfg.MinTxsPerCommitment = 32
	cfg.MaxTxsPerCommitment = 32
	cfg.MinCommitmentsPerBatch = 1

	cmd, err := setup.NewConfiguredCommanderFromEnv(cfg)
	require.NoError(t, err)
	err = cmd.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cmd.Stop())
	}()

	domain := GetDomain(t, cmd.Client())
	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	senderWallet := wallets[1]

	ethClient := newEthClient(t, cmd.Client())

	measureDisputeSignatureTransfer(t, cmd.Client(), ethClient)
	measureDisputeSignatureC2T(t, cmd.Client(), ethClient, wallets)

	testSubmitTransferBatch(t, cmd.Client(), senderWallet, 0)

	measureDisputeTransitionTransferInvalidStateRoot(t, cmd.Client(), ethClient)
	measureDisputeTransitionC2TInvalidStateRoot(t, cmd.Client(), ethClient, wallets)
}

func measureDisputeSignatureTransfer(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	send32TransfersBatchWithInvalidSignature(t, ethClient)
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 1)
}

func measureDisputeSignatureC2T(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, wallets []bls.Wallet) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	send32C2TBatchWithInvalidSignature(t, ethClient, wallets)
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 1)
}

func measureDisputeTransitionTransferInvalidStateRoot(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	send32TransfersBatchWithInvalidStateRoot(t, ethClient)
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 2)
}

func measureDisputeTransitionC2TInvalidStateRoot(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, wallets []bls.Wallet) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	send32C2TBatchWithInvalidStateRoot(t, ethClient, wallets)
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 2)
}

func send32TransfersBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client) {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	require.NoError(t, err)

	sendTransferCommitment(t, ethClient, bytes.Repeat(encodedTransfer, 32), 2)
}

func send32C2TBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client, wallets []bls.Wallet) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(38),
	}

	registeredPubKeyIDs := register32Accounts(t, ethClient, wallets)
	encodedTransfers := encodeCreate2Transfers(t, &transfer, registeredPubKeyIDs, 38)

	sendC2TCommitment(t, ethClient, encodedTransfers, 2)
}

func send32TransfersBatchWithInvalidSignature(t *testing.T, ethClient *eth.Client) {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	require.NoError(t, err)
	postStateRoot := common.Hash{45, 76, 35, 230, 155, 178, 7, 67, 241, 86, 195, 114, 225, 244, 169, 166, 182, 213, 46, 60, 106, 107, 252,
		125, 107, 78, 157, 106, 126, 38, 160, 137}

	commitment := models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: bytes.Repeat(encodedTransfer, 32),
	}
	submitTransfersBatch(t, ethClient, []models.CommitmentWithTxs{commitment}, 1)
}

func send32C2TBatchWithInvalidSignature(t *testing.T, ethClient *eth.Client, wallets []bls.Wallet) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(6),
	}

	registeredPubKeyIDs := register32Accounts(t, ethClient, wallets)
	encodedTransfers := encodeCreate2Transfers(t, &transfer, registeredPubKeyIDs, 6)

	postStateRoot := common.Hash{9, 165, 135, 45, 162, 158, 64, 129, 26, 232, 17, 209, 169, 198, 175, 189, 42, 40, 119, 15, 11, 78, 238,
		158, 35, 163, 205, 164, 23, 120, 249, 253}

	commitment := models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfers,
	}
	submitC2TBatch(t, ethClient, []models.CommitmentWithTxs{commitment}, 1)
}

func encodeCreate2Transfers(t *testing.T, transfer *models.Create2Transfer, registeredPubKeyIDs []uint32, startStateID uint32) []byte {
	encodedTransfers := make([]byte, 0, encoder.Create2TransferLength*len(registeredPubKeyIDs))
	for i := range registeredPubKeyIDs {
		transfer.ToStateID = ref.Uint32(startStateID + uint32(i))
		encodedTx, err := encoder.EncodeCreate2TransferForCommitment(transfer, registeredPubKeyIDs[i])
		require.NoError(t, err)

		encodedTransfers = append(encodedTransfers, encodedTx...)
	}
	return encodedTransfers
}

func register32Accounts(t *testing.T, ethClient *eth.Client, wallets []bls.Wallet) []uint32 {
	publicKeyBatch := make([]models.PublicKey, 16)
	registeredPubKeyIDs := make([]uint32, 0, 32)
	walletIndex := len(wallets) - 32
	for i := 0; i < 2; i++ {
		for j := range publicKeyBatch {
			publicKeyBatch[j] = *wallets[walletIndex].PublicKey()
			walletIndex++
		}
		pubKeyIDs, err := ethClient.RegisterBatchAccountAndWait(publicKeyBatch)
		require.NoError(t, err)
		registeredPubKeyIDs = append(registeredPubKeyIDs, pubKeyIDs...)
	}
	return registeredPubKeyIDs
}
