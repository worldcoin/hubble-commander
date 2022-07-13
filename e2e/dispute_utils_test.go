//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

func (s *DisputesE2ETestSuite) requireRollbackCompleted(triggerRollback func()) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := s.ETHClient.Rollup.WatchRollbackStatus(nil, sink)
	s.NoError(err)
	defer subscription.Unsubscribe()

	triggerRollback()

	s.Eventually(func() bool {
		for {
			select {
			case err := <-subscription.Err():
				s.NoError(err)
				return false
			case rollbackStatus := <-sink:
				if rollbackStatus.Completed {
					receipt, err := s.ETHClient.Blockchain.GetBackend().TransactionReceipt(context.Background(), rollbackStatus.Raw.TxHash)
					s.NoError(err)
					log.Infof("Rollback gas used: %d", receipt.GasUsed)
					return true
				}
			}
		}
	}, 30*time.Second, testutils.TryInterval)
}

func (s *DisputesE2ETestSuite) requireBatchesCount(expectedCount int) {
	batches := s.GetAllBatches()
	s.Len(batches, expectedCount)
}

func (s *DisputesE2ETestSuite) calculateWithdrawRoot(receiverBalance models.Uint256, txsAmount int) common.Hash {
	hash, err := encoder.HashUserState(&models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  receiverBalance,
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	hashes := make([]common.Hash, 0, txsAmount)
	for i := 0; i < txsAmount; i++ {
		hashes = append(hashes, *hash)
	}

	merkleTree, err := merkletree.NewMerkleTree(hashes)
	s.NoError(err)
	return merkleTree.Root()
}

func (s *DisputesE2ETestSuite) sendNTransfersBatchWithInvalidSignature(batchID uint64, txsAmount uint32, postStateRoot common.Hash) {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	s.NoError(err)

	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: bytes.Repeat(encodedTransfer, int(txsAmount)),
	}
	s.submitTransfersBatch([]models.CommitmentWithTxs{&commitment}, batchID)
}

func (s *DisputesE2ETestSuite) sendNTransfersBatchWithInvalidStateRoot(batchID uint64, txsAmount uint32) {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: 2,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	s.NoError(err)

	s.sendTransferCommitment(bytes.Repeat(encodedTransfer, int(txsAmount)), batchID)
}

func (s *DisputesE2ETestSuite) sendNC2TsBatchWithInvalidSignature(batchID uint64, txsAmount uint32, postStateRoot common.Hash) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(6),
	}

	encodedTransfers := s.registerAccountsAndEncodeC2Ts(&transfer, txsAmount)

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
	s.submitC2TBatch([]models.CommitmentWithTxs{&commitment}, batchID)
}

func (s *DisputesE2ETestSuite) sendNC2TsBatchWithInvalidStateRoot(batchID uint64, txsAmount uint32) {
	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		ToStateID: ref.Uint32(38),
	}

	encodedTransfers := s.registerAccountsAndEncodeC2Ts(&transfer, txsAmount)

	s.sendC2TCommitment(encodedTransfers, batchID)
}

func (s *DisputesE2ETestSuite) registerAccounts() []uint32 {
	publicKeyBatch := make([]models.PublicKey, 16)
	registeredPubKeyIDs := make([]uint32, 0, 32)
	walletIndex := len(s.Wallets) - 32
	for i := 0; i < 2; i++ {
		for j := range publicKeyBatch {
			publicKeyBatch[j] = *s.Wallets[walletIndex].PublicKey()
			walletIndex++
		}
		pubKeyIDs, err := s.ETHClient.RegisterBatchAccountAndWait(publicKeyBatch)
		s.NoError(err)
		registeredPubKeyIDs = append(registeredPubKeyIDs, pubKeyIDs...)
	}
	return registeredPubKeyIDs
}

func (s *DisputesE2ETestSuite) encodeCreate2Transfers(
	transfer *models.Create2Transfer,
	registeredPubKeyIDs []uint32,
	startStateID uint32,
) []byte {
	encodedTransfers := make([]byte, 0, encoder.Create2TransferLength*len(registeredPubKeyIDs))
	for i := range registeredPubKeyIDs {
		transfer.ToStateID = ref.Uint32(startStateID + uint32(i))
		encodedTx, err := encoder.EncodeCreate2TransferForCommitment(transfer, registeredPubKeyIDs[i])
		s.NoError(err)

		encodedTransfers = append(encodedTransfers, encodedTx...)
	}
	return encodedTransfers
}

func (s *DisputesE2ETestSuite) registerAccountsAndEncodeC2Ts(transfer *models.Create2Transfer, txsAmount uint32) []byte {
	registeredPubKeyIDs := s.registerAccounts()
	startStateID := 6 + txsAmount // 6 = registered accounts with user states based on the genesis file
	encodedTransfers := s.encodeCreate2Transfers(transfer, registeredPubKeyIDs[:txsAmount], startStateID)

	return encodedTransfers
}

func (s *DisputesE2ETestSuite) sendNMMsBatchWithInvalidSignature(batchID uint64, txsAmount uint32, postStateRoot common.Hash) {
	tx := models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
		},
		SpokeID: 1,
	}

	encodedTx, err := encoder.EncodeMassMigrationForCommitment(&tx)
	s.NoError(err)

	commitment := models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: postStateRoot,
			},
			CombinedSignature: models.Signature{},
			Meta: &models.MassMigrationMeta{
				SpokeID:     tx.SpokeID,
				TokenID:     models.MakeUint256(0),
				Amount:      *tx.Amount.MulN(uint64(txsAmount)),
				FeeReceiver: 0,
			},
			WithdrawRoot: s.calculateWithdrawRoot(tx.Amount, int(txsAmount)),
		},
		Transactions: bytes.Repeat(encodedTx, int(txsAmount)),
	}

	s.submitMMBatch([]models.CommitmentWithTxs{&commitment}, batchID)
}

func (s *DisputesE2ETestSuite) sendNMMsBatchWithInvalidStateRoot(batchID uint64, txsAmount uint32) {
	tx := models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(90),
			Fee:         models.MakeUint256(10),
		},
		SpokeID: 1,
	}

	encodedTx, err := encoder.EncodeMassMigrationForCommitment(&tx)
	s.NoError(err)

	totalAmount := tx.Amount.MulN(uint64(txsAmount)).Uint64()
	withdrawRoot := s.calculateWithdrawRoot(tx.Amount, int(txsAmount))

	s.sendMMCommitment(bytes.Repeat(encodedTx, int(txsAmount)), withdrawRoot, totalAmount, batchID)
}

func (s *DisputesE2ETestSuite) sendTransferCommitment(encodedTransfer []byte, batchID uint64) {
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
	s.submitTransfersBatch([]models.CommitmentWithTxs{&commitment}, batchID)
}

func (s *DisputesE2ETestSuite) submitTransfersBatch(commitments []models.CommitmentWithTxs, batchID uint64) {
	transaction, err := s.ETHClient.SubmitTransfersBatch(context.Background(), models.NewUint256(batchID), commitments)
	s.NoError(err)

	s.waitForSubmittedBatch(transaction, batchID)
}

func (s *DisputesE2ETestSuite) sendC2TCommitment(encodedTransfer []byte, batchID uint64) {
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

	s.submitC2TBatch([]models.CommitmentWithTxs{&commitment}, batchID)
}

func (s *DisputesE2ETestSuite) submitC2TBatch(commitments []models.CommitmentWithTxs, batchID uint64) {
	transaction, err := s.ETHClient.SubmitCreate2TransfersBatch(context.Background(), models.NewUint256(batchID), commitments)
	s.NoError(err)

	s.waitForSubmittedBatch(transaction, batchID)
}

func (s *DisputesE2ETestSuite) sendMMCommitment(encodedTxs []byte, withdrawRoot common.Hash, totalAmount, batchID uint64) {
	commitment := models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				PostStateRoot: utils.RandomHash(),
			},
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

	s.submitMMBatch([]models.CommitmentWithTxs{&commitment}, batchID)
}

func (s *DisputesE2ETestSuite) submitMMBatch(
	commitments []models.CommitmentWithTxs,
	batchID uint64,
) {
	transaction, err := s.ETHClient.SubmitMassMigrationsBatch(models.NewUint256(batchID), commitments)
	s.NoError(err)

	s.waitForSubmittedBatch(transaction, batchID)
}

func (s *DisputesE2ETestSuite) waitForSubmittedBatch(transaction *types.Transaction, batchID uint64) {
	_, err := s.ETHClient.WaitToBeMined(transaction)
	s.NoError(err)

	_, err = s.ETHClient.GetContractBatch(models.NewUint256(batchID))
	s.NoError(err)
}
