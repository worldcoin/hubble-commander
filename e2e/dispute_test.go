//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommanderDispute(t *testing.T) {
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
	receiverWallet := wallets[len(wallets)-1]

	ethClient := newEthClient(t, cmd.Client())

	testDisputeSignatureTransfer(t, cmd.Client(), ethClient)
	testDisputeSignatureC2T(t, cmd.Client(), ethClient, receiverWallet)

	testDisputeTransitionTransfer(t, cmd.Client(), ethClient, senderWallet)
	testDisputeTransitionC2T(t, cmd.Client(), ethClient, senderWallet, wallets)

	testDisputeTransitionTransferInvalidStateRoot(t, cmd.Client(), ethClient)
	testDisputeTransitionC2TInvalidStateRoot(t, cmd.Client(), ethClient, receiverWallet)
}

func testDisputeSignatureTransfer(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	sendTransferBatchWithInvalidSignature(t, ethClient)
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 1)
}

func testDisputeSignatureC2T(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, receiverWallet bls.Wallet) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	sendC2TBatchWithInvalidSignature(t, ethClient, receiverWallet.PublicKey())
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 1)
}

func testDisputeTransitionTransfer(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, senderWallet bls.Wallet) {
	testSubmitTransferBatch(t, client, senderWallet, 0)

	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	sendTransferBatchWithInvalidAmount(t, ethClient)
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 2)

	testSubmitTransferBatch(t, client, senderWallet, 32)
}

func testDisputeTransitionTransferInvalidStateRoot(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	sendTransferBatchWithInvalidStateRoot(t, ethClient)
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 4)
}

func testDisputeTransitionC2T(
	t *testing.T,
	client jsonrpc.RPCClient,
	ethClient *eth.Client,
	senderWallet bls.Wallet,
	wallets []bls.Wallet,
) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	firstC2TWallet := wallets[len(wallets)-32]
	sendC2TBatchWithInvalidAmount(t, ethClient, firstC2TWallet.PublicKey())
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 3)

	testSubmitC2TBatch(t, client, senderWallet, wallets, firstC2TWallet.PublicKey(), 64)
}

func testDisputeTransitionC2TInvalidStateRoot(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, receiverWallet bls.Wallet) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	sendC2TBatchWithInvalidStateRoot(t, ethClient, receiverWallet.PublicKey())
	testRollbackCompletion(t, ethClient, sink, subscription)

	testBatchesAfterDispute(t, client, 4)
}

func testSubmitTransferBatch(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, startNonce uint64) {
	firstTransferHash := testSendTransfer(t, client, senderWallet, startNonce)
	testGetTransaction(t, client, firstTransferHash)
	send31MoreTransfers(t, client, senderWallet, startNonce+1)

	waitForTxToBeIncludedInBatch(t, client, firstTransferHash)
}

func testSubmitC2TBatch(
	t *testing.T,
	client jsonrpc.RPCClient,
	senderWallet bls.Wallet,
	wallets []bls.Wallet,
	targetPublicKey *models.PublicKey,
	startNonce uint64,
) {
	firstTransferHash := testSendCreate2Transfer(t, client, senderWallet, targetPublicKey, startNonce)
	testGetTransaction(t, client, firstTransferHash)
	send31MoreCreate2Transfers(t, client, senderWallet, wallets, startNonce+1)

	waitForTxToBeIncludedInBatch(t, client, firstTransferHash)
}

func testBatchesAfterDispute(t *testing.T, client jsonrpc.RPCClient, expectedLength int) {
	var batches []dto.Batch
	err := client.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})

	require.NoError(t, err)
	require.Len(t, batches, expectedLength)
}

func testRollbackCompletion(
	t *testing.T,
	ethClient *eth.Client,
	sink chan *rollup.RollupRollbackStatus,
	subscription event.Subscription,
) {
	require.Eventually(t, func() bool {
		for {
			select {
			case err := <-subscription.Err():
				require.NoError(t, err)
				return false
			case rollbackStatus := <-sink:
				if rollbackStatus.Completed {
					receipt, err := ethClient.Blockchain.GetBackend().TransactionReceipt(context.Background(), rollbackStatus.Raw.TxHash)
					require.NoError(t, err)
					logrus.Infof("Rollback gas used: %d", receipt.GasUsed)
					return true
				}
			}
		}
	}, 30*time.Second, testutils.TryInterval)
}

func sendTransferBatchWithInvalidSignature(t *testing.T, ethClient *eth.Client) {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	require.NoError(t, err)
	postStateRoot := common.Hash{223, 216, 112, 56, 113, 248, 202, 217, 207, 95, 115, 189, 153, 14, 156, 202, 27, 160, 133, 14, 118, 218,
		161, 109, 17, 61, 77, 118, 7, 252, 42, 18}

	commitment := models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfer,
	}
	submitTransfersBatch(t, ethClient, []models.CommitmentWithTxs{commitment}, 1)
}

func sendTransferBatchWithInvalidAmount(t *testing.T, ethClient *eth.Client) {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(2_000_000_000_000_000_000),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	require.NoError(t, err)

	sendTransferCommitment(t, ethClient, encodedTransfer, 2)
}

func sendTransferBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client) {
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

	sendTransferCommitment(t, ethClient, encodedTransfer, 4)
}

func sendC2TBatchWithInvalidAmount(t *testing.T, ethClient *eth.Client, toPublicKey *models.PublicKey) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(2_000_000_000_000_000_000),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(6),
	}

	pubKeyID, err := ethClient.RegisterAccountAndWait(toPublicKey)
	require.NoError(t, err)

	encodedTransfer, err := encoder.EncodeCreate2TransferForCommitment(&transfer, *pubKeyID)
	require.NoError(t, err)

	sendC2TCommitment(t, ethClient, encodedTransfer, 3)
}

func sendC2TBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client, toPublicKey *models.PublicKey) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(38),
	}

	pubKeyID, err := ethClient.RegisterAccountAndWait(toPublicKey)
	require.NoError(t, err)

	encodedTransfer, err := encoder.EncodeCreate2TransferForCommitment(&transfer, *pubKeyID)
	require.NoError(t, err)

	sendC2TCommitment(t, ethClient, encodedTransfer, 4)
}

func sendC2TBatchWithInvalidSignature(t *testing.T, ethClient *eth.Client, toPublicKey *models.PublicKey) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(6),
	}

	pubKeyID, err := ethClient.RegisterAccountAndWait(toPublicKey)
	require.NoError(t, err)

	encodedTransfer, err := encoder.EncodeCreate2TransferForCommitment(&transfer, *pubKeyID)
	require.NoError(t, err)
	postStateRoot := common.Hash{5, 64, 118, 3, 181, 231, 59, 98, 230, 215, 146, 132, 59, 141, 73, 132, 133, 23, 149, 118, 59, 118, 88, 153,
		150, 65, 112, 215, 128, 132, 47, 58}

	commitment := models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfer,
	}
	submitC2TBatch(t, ethClient, []models.CommitmentWithTxs{commitment}, 1)
}

func sendTransferCommitment(t *testing.T, ethClient *eth.Client, encodedTransfer []byte, batchID uint64) {
	commitment := models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfer,
	}
	submitTransfersBatch(t, ethClient, []models.CommitmentWithTxs{commitment}, batchID)
}

func submitTransfersBatch(t *testing.T, ethClient *eth.Client, commitments []models.CommitmentWithTxs, batchID uint64) {
	transaction, err := ethClient.SubmitTransfersBatch(models.NewUint256(batchID), commitments)
	require.NoError(t, err)

	waitForSubmittedBatch(t, ethClient, transaction, batchID)
}

func sendC2TCommitment(t *testing.T, ethClient *eth.Client, encodedTransfer []byte, batchID uint64) {
	commitment := models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfer,
	}

	submitC2TBatch(t, ethClient, []models.CommitmentWithTxs{commitment}, batchID)
}

func submitC2TBatch(t *testing.T, ethClient *eth.Client, commitments []models.CommitmentWithTxs, batchID uint64) {
	transaction, err := ethClient.SubmitCreate2TransfersBatch(models.NewUint256(batchID), commitments)
	require.NoError(t, err)

	waitForSubmittedBatch(t, ethClient, transaction, batchID)
}

func waitForSubmittedBatch(t *testing.T, ethClient *eth.Client, transaction *types.Transaction, batchID uint64) {
	_, err := chain.WaitToBeMined(ethClient.Blockchain.GetBackend(), transaction)
	require.NoError(t, err)

	_, err = ethClient.GetBatch(models.NewUint256(batchID))
	require.NoError(t, err)
}

func newEthClient(t *testing.T, client jsonrpc.RPCClient) *eth.Client {
	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	require.NoError(t, err)

	chainState := models.ChainState{
		ChainID:                        info.ChainID,
		AccountRegistry:                info.AccountRegistry,
		AccountRegistryDeploymentBlock: info.AccountRegistryDeploymentBlock,
		TokenRegistry:                  info.TokenRegistry,
		DepositManager:                 info.DepositManager,
		Rollup:                         info.Rollup,
	}

	cfg := config.GetConfig()
	blockchain, err := chain.NewRPCCConnection(cfg.Ethereum)
	require.NoError(t, err)

	backend := blockchain.GetBackend()

	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, backend)
	require.NoError(t, err)

	tokenRegistry, err := tokenregistry.NewTokenRegistry(chainState.TokenRegistry, backend)
	require.NoError(t, err)

	depositManager, err := depositmanager.NewDepositManager(chainState.DepositManager, backend)
	require.NoError(t, err)

	rollupContract, err := rollup.NewRollup(chainState.Rollup, backend)
	require.NoError(t, err)

	ethClient, err := eth.NewClient(blockchain, metrics.NewCommanderMetrics(), &eth.NewClientParams{
		ChainState:      chainState,
		AccountRegistry: accountRegistry,
		TokenRegistry:   tokenRegistry,
		DepositManager:  depositManager,
		Rollup:          rollupContract,
	})
	require.NoError(t, err)
	return ethClient
}
