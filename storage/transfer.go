package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *TransactionStorage) AddTransfer(tx *models.Transfer) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		if tx.CommitmentID != nil || tx.ErrorMessage != nil {
			err := txStorage.addStoredTxReceipt(stored.NewTxReceiptFromTransfer(tx))
			if err != nil {
				return err
			}
		}
		return txStorage.addStoredTx(stored.NewTxFromTransfer(tx))
	})
}

func (s *TransactionStorage) BatchAddTransfer(txs []models.Transfer) error {
	if len(txs) < 1 {
		return errors.WithStack(ErrNoRowsAffected)
	}

	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			err := txStorage.AddTransfer(&txs[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Transfer {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return tx.ToTransfer(txReceipt), nil
}

func (s *TransactionStorage) GetPendingTransfers() (txs models.TransferArray, err error) {
	err = s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		txs, err = txStorage.unsafeGetPendingTransfers()
		return err
	})
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionStorage) unsafeGetPendingTransfers() ([]models.Transfer, error) {
	txs := make([]models.Transfer, 0, 32)
	var storedTx stored.Tx
	err := s.database.Badger.Iterator(stored.TxPrefix, db.KeyIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			skip, err := s.getStoredTxFromItem(item, &storedTx)
			if err != nil || skip {
				return false, err
			}
			if storedTx.TxType == txtype.Transfer {
				txs = append(txs, *storedTx.ToTransfer(nil))
			}
			return false, nil
		})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}

	return txs, nil
}

func (s *TransactionStorage) GetTransfersByCommitmentID(id models.CommitmentID) ([]models.Transfer, error) {
	transfers := make([]models.Transfer, 0, 1)

	err := s.iterateTxsByCommitmentID(id, func(storedTx *stored.Tx, storedTxReceipt *stored.TxReceipt) {
		if storedTx.TxType == txtype.Transfer {
			transfers = append(transfers, *storedTx.ToTransfer(storedTxReceipt))
		}
	})
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func (s *TransactionStorage) iterateTxsByCommitmentID(
	id models.CommitmentID,
	handleTx func(storedTx *stored.Tx, storedTxReceipt *stored.TxReceipt),
) error {
	return s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		return txStorage.unsafeIterateTxsByCommitmentID(id, handleTx)
	})
}

func (s *TransactionStorage) unsafeIterateTxsByCommitmentID(
	id models.CommitmentID,
	handleTx func(storedTx *stored.Tx, storedTxReceipt *stored.TxReceipt),
) error {
	receipts := make([]stored.TxReceipt, 0, 1)
	err := s.database.Badger.Find(
		&receipts,
		bh.Where("CommitmentID").Eq(id).Index("CommitmentID"),
	)
	if err != nil {
		return err
	}

	for i := range receipts {
		storedTx, storedTxReceipt, err := s.getStoredTxWithReceipt(receipts[i].Hash)
		if err != nil {
			return err
		}
		handleTx(storedTx, storedTxReceipt)
	}
	return nil
}

func (s *TransactionStorage) MarkTransfersAsIncluded(txs []models.Transfer, commitmentID *models.CommitmentID) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			txReceipt := stored.NewTxReceiptFromTransfer(&txs[i])
			txReceipt.CommitmentID = commitmentID
			err := txStorage.addStoredTxReceipt(txReceipt)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
