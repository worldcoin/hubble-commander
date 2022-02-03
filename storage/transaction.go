package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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

func (s *TransactionStorage) UpdateTransaction(tx models.GenericTransaction) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		txBase := tx.GetBase()
		receipt, err := txStorage.getStoredTxReceipt(txBase.Hash)
		if err != nil {
			return err
		}
		if receipt == nil {
			return NewNotFoundError("txReceipt")
		}
		if receipt.ErrorMessage == nil {
			return errors.WithStack(ErrAlreadyMinedTransaction)
		}

		err = txStorage.MarkTransactionsAsPending([]common.Hash{txBase.Hash})
		if err != nil {
			return err
		}
		return txStorage.updateStoredTx(stored.NewTx(tx))
	})
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
	result.MinedTime = batch.MinedTime

	return result, nil
}

func (s *TransactionStorage) AddTransaction(tx models.GenericTransaction) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		if tx.GetBase().CommitmentID != nil || tx.GetBase().ErrorMessage != nil {
			err := txStorage.addStoredTxReceipt(stored.NewTxReceipt(tx))
			if err != nil {
				return err
			}
		}
		return txStorage.addStoredTx(stored.NewTx(tx))
	})
}

func (s *TransactionStorage) BatchAddTransaction(txs models.GenericTransactionArray) error {
	if txs.Len() < 1 {
		return errors.WithStack(ErrNoRowsAffected)
	}

	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := 0; i < txs.Len(); i++ {
			err := txStorage.AddTransaction(txs.At(i))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
