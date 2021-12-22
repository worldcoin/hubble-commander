package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
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

func (s *TransactionStorage) UpdateTransaction(tx models.GenericTransaction) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		txBase := tx.GetBase()
		receipt, err := txStorage.getStoredTxReceipt(txBase.Hash)
		if err != nil {
			return err
		}
		if receipt == nil || receipt.ErrorMessage == nil {
			return nil
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
	result.BatchTime = batch.SubmissionTime

	return result, nil
}

func (s *TransactionStorage) AddTransaction(tx models.GenericTransaction) error {
	switch tx.Type() {
	case txtype.Transfer:
		return s.AddTransfer(tx.ToTransfer())
	case txtype.Create2Transfer:
		return s.AddCreate2Transfer(tx.ToCreate2Transfer())
	case txtype.MassMigration:
		return s.AddMassMigration(tx.ToMassMigration())
	}
	return nil
}

func (s *TransactionStorage) BatchAddTransaction(txs models.GenericTransactionArray) error {
	switch txs.Type() {
	case txtype.Transfer:
		return s.BatchAddTransfer(txs.ToTransferArray())
	case txtype.Create2Transfer:
		return s.BatchAddCreate2Transfer(txs.ToCreate2TransferArray())
	case txtype.MassMigration:
		return s.BatchAddMassMigration(txs.ToMassMigrationArray())
	}
	return nil
}
