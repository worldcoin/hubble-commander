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
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
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
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommanderDispute(t *testing.T) {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = 32
	cfg.Rollup.MaxTxsPerCommitment = 32
	cfg.Rollup.MinCommitmentsPerBatch = 1

	cmd, err := setup.NewConfiguredCommanderFromEnv(cfg, nil)
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
	testDisputeSignatureMM(t, cmd.Client(), ethClient)

	testDisputeTransitionTransfer(t, cmd.Client(), ethClient, senderWallet)
	testDisputeTransitionC2T(t, cmd.Client(), ethClient, senderWallet, receiverWallet, wallets)
	testDisputeTransitionMM(t, cmd.Client(), ethClient, senderWallet)
}

func testDisputeSignatureTransfer(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	requireRollbackCompleted(t, ethClient, func() {
		sendTransferBatchWithInvalidSignature(t, ethClient, 1)
	})

	requireBatchesCount(t, client, 1)
}

func testDisputeSignatureC2T(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, receiverWallet bls.Wallet) {
	requireRollbackCompleted(t, ethClient, func() {
		sendC2TBatchWithInvalidSignature(t, ethClient, receiverWallet.PublicKey())
	})

	requireBatchesCount(t, client, 1)
}

func testDisputeSignatureMM(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	requireRollbackCompleted(t, ethClient, func() {
		sendMMBatchWithInvalidSignature(t, ethClient, 1)
	})

	requireBatchesCount(t, client, 1)
}

func testDisputeTransitionTransfer(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, senderWallet bls.Wallet) {
	submitTxBatchAndWait(t, client, func() common.Hash {
		return testSubmitTransferBatch(t, client, senderWallet, 0)
	})

	requireRollbackCompleted(t, ethClient, func() {
		sendTransferBatchWithInvalidStateRoot(t, ethClient, 2)
	})

	requireBatchesCount(t, client, 2)

	submitTxBatchAndWait(t, client, func() common.Hash {
		return testSubmitTransferBatch(t, client, senderWallet, 32)
	})
}

func testDisputeTransitionC2T(
	t *testing.T,
	client jsonrpc.RPCClient,
	ethClient *eth.Client,
	senderWallet bls.Wallet,
	receiverWallet bls.Wallet,
	wallets []bls.Wallet,
) {
	requireRollbackCompleted(t, ethClient, func() {
		sendC2TBatchWithInvalidStateRoot(t, ethClient, receiverWallet.PublicKey(), 3)
	})

	requireBatchesCount(t, client, 3)

	firstC2TWallet := wallets[len(wallets)-32]
	submitTxBatchAndWait(t, client, func() common.Hash {
		return testSubmitC2TBatch(t, client, senderWallet, wallets, firstC2TWallet.PublicKey(), 64)
	})
}

func testDisputeTransitionMM(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, senderWallet bls.Wallet) {
	requireRollbackCompleted(t, ethClient, func() {
		sendMMBatchWithInvalidStateRoot(t, ethClient, 4)
	})

	requireBatchesCount(t, client, 4)

	submitTxBatchAndWait(t, client, func() common.Hash {
		return testSubmitMassMigrationBatch(t, client, senderWallet, 96)
	})
}

func requireBatchesCount(t *testing.T, client jsonrpc.RPCClient, expectedCount int) {
	var batches []dto.Batch
	err := client.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})

	require.NoError(t, err)
	require.Len(t, batches, expectedCount)
}

func requireRollbackCompleted(t *testing.T, ethClient *eth.Client, triggerRollback func()) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := ethClient.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	require.NoError(t, err)
	defer subscription.Unsubscribe()

	triggerRollback()

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

func sendTransferBatchWithInvalidSignature(t *testing.T, ethClient *eth.Client, batchID uint64) {
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

	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfer,
	}
	submitTransfersBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, batchID)
}

func sendTransferBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client, batchID uint64) {
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

	sendTransferCommitment(t, ethClient, encodedTransfer, batchID)
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

	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfer,
	}
	submitC2TBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, 1)
}

func sendC2TBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client, toPublicKey *models.PublicKey, batchID uint64) {
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

	sendC2TCommitment(t, ethClient, encodedTransfer, batchID)
}

func sendMMBatchWithInvalidSignature(t *testing.T, ethClient *eth.Client, batchID uint64) {
	tx := models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
		},
		SpokeID: 1,
	}

	encodedTx, err := encoder.EncodeMassMigrationForCommitment(&tx)
	require.NoError(t, err)

	postStateRoot := common.Hash{25, 2, 167, 141, 141, 223, 41, 53, 199, 36, 50, 52, 166, 110, 139, 144, 117, 71, 15, 68, 65, 127, 115, 174,
		77, 40, 231, 185, 228, 186, 225, 136}

	commitment := models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
			Meta: &models.MassMigrationMeta{
				SpokeID:     tx.SpokeID,
				TokenID:     models.MakeUint256(0),
				Amount:      tx.Amount,
				FeeReceiver: 0,
			},
			WithdrawRoot: calculateWithdrawRoot(t, tx.Amount, 1),
		},
		Transactions: encodedTx,
	}

	submitMMBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, batchID)
}

func sendMMBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client, batchID uint64) {
	tx := models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		SpokeID: 1,
	}

	encodedTx, err := encoder.EncodeMassMigrationForCommitment(&tx)
	require.NoError(t, err)

	hash, err := encoder.HashUserState(&models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  tx.Amount,
		Nonce:    models.MakeUint256(0),
	})
	require.NoError(t, err)

	merkleTree, err := merkletree.NewMerkleTree([]common.Hash{*hash})
	require.NoError(t, err)

	sendMMCommitment(t, ethClient, encodedTx, merkleTree.Root(), tx.Amount.Uint64(), batchID)
}

func sendTransferCommitment(t *testing.T, ethClient *eth.Client, encodedTransfer []byte, batchID uint64) {
	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfer,
	}
	submitTransfersBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, batchID)
}

func submitTransfersBatch(t *testing.T, ethClient *eth.Client, commitments []models.CommitmentWithTxs, batchID uint64) {
	transaction, err := ethClient.SubmitTransfersBatch(models.NewUint256(batchID), commitments)
	require.NoError(t, err)

	waitForSubmittedBatch(t, ethClient, transaction, batchID)
}

func sendC2TCommitment(t *testing.T, ethClient *eth.Client, encodedTransfer []byte, batchID uint64) {
	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfer,
	}

	submitC2TBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, batchID)
}

func submitC2TBatch(t *testing.T, ethClient *eth.Client, commitments []models.CommitmentWithTxs, batchID uint64) {
	transaction, err := ethClient.SubmitCreate2TransfersBatch(models.NewUint256(batchID), commitments)
	require.NoError(t, err)

	waitForSubmittedBatch(t, ethClient, transaction, batchID)
}

func sendMMCommitment(t *testing.T, ethClient *eth.Client, encodedTxs []byte, withdrawRoot common.Hash, totalAmount, batchID uint64) {
	commitment := models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
			Meta: &models.MassMigrationMeta{
				SpokeID:     1,
				TokenID:     models.MakeUint256(0),
				Amount:      models.MakeUint256(totalAmount),
				FeeReceiver: 0,
			},
			WithdrawRoot: withdrawRoot,
		},
		Transactions: encodedTxs,
	}

	submitMMBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, batchID)
}

func submitMMBatch(
	t *testing.T,
	ethClient *eth.Client,
	commitments []models.CommitmentWithTxs,
	batchID uint64,
) {
	transaction, err := ethClient.SubmitMassMigrationsBatch(models.NewUint256(batchID), commitments)
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
		SpokeRegistry:                  info.SpokeRegistry,
		DepositManager:                 info.DepositManager,
		WithdrawManager:                info.WithdrawManager,
		Rollup:                         info.Rollup,
	}

	cfg := config.GetConfig()
	blockchain, err := chain.NewRPCConnection(cfg.Ethereum)
	require.NoError(t, err)

	backend := blockchain.GetBackend()

	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, backend)
	require.NoError(t, err)

	spokeRegistry, err := spokeregistry.NewSpokeRegistry(chainState.SpokeRegistry, backend)
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
		SpokeRegistry:   spokeRegistry,
		TokenRegistry:   tokenRegistry,
		DepositManager:  depositManager,
		Rollup:          rollupContract,
	})
	require.NoError(t, err)
	return ethClient
}
