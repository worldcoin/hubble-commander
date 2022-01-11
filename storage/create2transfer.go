package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (s *TransactionStorage) BatchAddCreate2Transfer(txs []models.Create2Transfer) error {
	return s.BatchAddTransaction(models.MakeCreate2TransferArray(txs...))
}

func (s *TransactionStorage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Create2Transfer {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return tx.ToCreate2Transfer(txReceipt), nil
}

func (s *TransactionStorage) GetPendingCreate2Transfers() (txs models.Create2TransferArray, err error) {
	err = s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		txs, err = txStorage.unsafeGetPendingCreate2Transfers()
		return err
	})
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionStorage) unsafeGetPendingCreate2Transfers() ([]models.Create2Transfer, error) {
	txs := make([]models.Create2Transfer, 0, 32)
	var storedTx stored.Tx
	err := s.database.Badger.Iterator(stored.TxPrefix, db.KeyIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			skip, err := s.getStoredTxFromItem(item, &storedTx)
			if err != nil || skip {
				return false, err
			}
			if storedTx.TxType == txtype.Create2Transfer {
				txs = append(txs, *storedTx.ToCreate2Transfer(nil))
			}
			return false, nil
		})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}

	return txs, nil
}

func (s *TransactionStorage) GetCreate2TransfersByCommitmentID(id models.CommitmentID) ([]models.Create2Transfer, error) {
	transfers := make([]models.Create2Transfer, 0, 1)

	err := s.iterateTxsByCommitmentID(id, func(storedTx *stored.Tx, storedTxReceipt *stored.TxReceipt) {
		if storedTx.TxType == txtype.Create2Transfer {
			transfers = append(transfers, *storedTx.ToCreate2Transfer(storedTxReceipt))
		}
	})
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func (s *TransactionStorage) MarkCreate2TransfersAsIncluded(txs []models.Create2Transfer, commitmentID *models.CommitmentID) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			storedTxReceipt := stored.NewTxReceiptFromCreate2Transfer(&txs[i])
			storedTxReceipt.CommitmentID = commitmentID
			err := txStorage.addStoredTxReceipt(storedTxReceipt)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
