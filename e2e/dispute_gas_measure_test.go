//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

const maxTxsPerCommitment = 32

func TestMeasureDisputeGasUsage(t *testing.T) {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = maxTxsPerCommitment
	cfg.Rollup.MaxTxsPerCommitment = maxTxsPerCommitment
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

	ethClient := newEthClient(t, cmd.Client())

	measureDisputeSignatureTransfer(t, cmd.Client(), ethClient)
	measureDisputeSignatureC2T(t, cmd.Client(), ethClient, wallets)
	measureDisputeSignatureMM(t, cmd.Client(), ethClient)

	submitTxBatchAndWait(t, cmd.Client(), func() common.Hash {
		return testSubmitTransferBatch(t, cmd.Client(), senderWallet, 0)
	})

	measureDisputeTransitionTransfer(t, cmd.Client(), ethClient)
	measureDisputeTransitionC2T(t, cmd.Client(), ethClient, wallets)
	measureDisputeTransitionMM(t, cmd.Client(), ethClient)
}

func measureDisputeSignatureTransfer(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	requireRollbackCompleted(t, ethClient, func() {
		send32TransfersBatchWithInvalidSignature(t, ethClient)
	})

	requireBatchesCount(t, client, 1)
}

func measureDisputeSignatureC2T(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, wallets []bls.Wallet) {
	requireRollbackCompleted(t, ethClient, func() {
		send32C2TBatchWithInvalidSignature(t, ethClient, wallets)
	})

	requireBatchesCount(t, client, 1)
}

func measureDisputeSignatureMM(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	requireRollbackCompleted(t, ethClient, func() {
		send32MMBatchWithInvalidSignature(t, ethClient)
	})

	requireBatchesCount(t, client, 1)
}

func measureDisputeTransitionTransfer(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	requireRollbackCompleted(t, ethClient, func() {
		send32TransfersBatchWithInvalidStateRoot(t, ethClient)
	})

	requireBatchesCount(t, client, 2)
}

func measureDisputeTransitionC2T(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client, wallets []bls.Wallet) {
	requireRollbackCompleted(t, ethClient, func() {
		send32C2TBatchWithInvalidStateRoot(t, ethClient, wallets)
	})

	requireBatchesCount(t, client, 2)
}

func measureDisputeTransitionMM(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) {
	requireRollbackCompleted(t, ethClient, func() {
		send32MMBatchWithInvalidStateRoot(t, ethClient)
	})

	requireBatchesCount(t, client, 2)
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

	sendTransferCommitment(t, ethClient, bytes.Repeat(encodedTransfer, maxTxsPerCommitment), 2)
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

func send32MMBatchWithInvalidStateRoot(t *testing.T, ethClient *eth.Client) {
	tx := models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		SpokeID: 1,
	}

	encodedTransfer, err := encoder.EncodeMassMigrationForCommitment(&tx)
	require.NoError(t, err)

	totalAmount := tx.Amount.MulN(uint64(maxTxsPerCommitment)).Uint64()
	withdrawRoot := calculateWithdrawRoot(t, tx.Amount, maxTxsPerCommitment)

	sendMMCommitment(t, ethClient, bytes.Repeat(encodedTransfer, maxTxsPerCommitment), withdrawRoot, totalAmount, 2)
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

	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: bytes.Repeat(encodedTransfer, maxTxsPerCommitment),
	}
	submitTransfersBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, 1)
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

	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: encodedTransfers,
	}
	submitC2TBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, 1)
}

func send32MMBatchWithInvalidSignature(t *testing.T, ethClient *eth.Client) {
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

	postStateRoot := common.Hash{68, 198, 251, 28, 54, 95, 42, 7, 136, 120, 20, 253, 146, 124, 84, 119, 183, 52, 27, 44, 225, 192, 165, 206,
		154, 69, 207, 53, 239, 253, 79, 216}

	commitment := models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			CombinedSignature: models.Signature{},
			Meta: &models.MassMigrationMeta{
				SpokeID:     tx.SpokeID,
				TokenID:     models.MakeUint256(0),
				Amount:      *tx.Amount.MulN(uint64(maxTxsPerCommitment)),
				FeeReceiver: 0,
			},
			WithdrawRoot: calculateWithdrawRoot(t, tx.Amount, maxTxsPerCommitment),
		},
		Transactions: bytes.Repeat(encodedTx, maxTxsPerCommitment),
	}

	submitMMBatch(t, ethClient, []models.CommitmentWithTxs{&commitment}, 1)
}

func calculateWithdrawRoot(t *testing.T, receiverBalance models.Uint256, txCount int) common.Hash {
	hash, err := encoder.HashUserState(&models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  receiverBalance,
		Nonce:    models.MakeUint256(0),
	})
	require.NoError(t, err)

	hashes := make([]common.Hash, 0, txCount)
	for i := 0; i < txCount; i++ {
		hashes = append(hashes, *hash)
	}

	merkleTree, err := merkletree.NewMerkleTree(hashes)
	require.NoError(t, err)
	return merkleTree.Root()
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
	registeredPubKeyIDs := make([]uint32, 0, maxTxsPerCommitment)
	walletIndex := len(wallets) - maxTxsPerCommitment
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
