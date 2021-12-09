package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) GetTransactionWithBatchDetails(hash common.Hash) (
	transaction *models.TransactionWithBatchDetails,
	err error,
) {
	err = s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		transaction, err = txStorage.unsafeGetTransactionWithBatchDetails(hash)
		return err
	})
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *Storage) unsafeGetTransactionWithBatchDetails(hash common.Hash) (
	*models.TransactionWithBatchDetails,
	error,
) {
	storedTx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}

	typedTxInterface := storedTx.ToTypedTxInterface(txReceipt)

	result := &models.TransactionWithBatchDetails{Transaction: typedTxInterface}

	if txReceipt == nil || txReceipt.CommitmentID == nil {
		return result, nil
	}

	batch, err := s.GetBatch(txReceipt.CommitmentID.BatchID)
	if err != nil {
		return nil, err
	}

	result.BatchHash = batch.Hash
	result.BatchTime = batch.SubmissionTime

	return result, nil
}
