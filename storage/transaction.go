package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) GetTransactionWithBatchDetails(hash common.Hash) (
	transaction *models.TransactionWithBatchDetails,
	txType *txtype.TransactionType,
	err error,
) {
	err = s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		transaction, txType, err = txStorage.unsafeGetTransactionWithBatchDetails(hash)
		return err
	})
	if err != nil {
		return nil, nil, err
	}
	return transaction, txType, nil
}

func (s *Storage) unsafeGetTransactionWithBatchDetails(hash common.Hash) (
	*models.TransactionWithBatchDetails,
	*txtype.TransactionType,
	error,
) {
	storedTx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, nil, err
	}

	typedTxInterface := storedTx.ToTypedTxInterface(txReceipt)

	result := &models.TransactionWithBatchDetails{Transaction: typedTxInterface}

	if txReceipt == nil || txReceipt.CommitmentID == nil {
		return result, &storedTx.TxType, nil
	}

	batch, err := s.GetBatch(txReceipt.CommitmentID.BatchID)
	if err != nil {
		return nil, nil, err
	}

	result.BatchHash = batch.Hash
	result.BatchTime = batch.SubmissionTime

	return result, &storedTx.TxType, nil
}
