package storage

import (
	"sort"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (s *TransactionStorage) AddMassMigration(m *models.MassMigration) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		if m.CommitmentID != nil || m.ErrorMessage != nil {
			err := txStorage.addStoredTxReceipt(models.NewStoredTxReceiptFromMassMigration(m))
			if err != nil {
				return err
			}
		}
		return txStorage.addStoredTx(models.NewStoredTxFromMassMigration(m))
	})
}

func (s *TransactionStorage) BatchAddMassMigration(txs []models.MassMigration) error {
	if len(txs) < 1 {
		return errors.WithStack(ErrNoRowsAffected)
	}

	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			err := txStorage.AddMassMigration(&txs[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) GetMassMigration(hash common.Hash) (*models.MassMigration, error) {
	tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.MassMigration {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return tx.ToMassMigration(txReceipt), nil
}

func (s *TransactionStorage) GetPendingMassMigrations() (txs models.MassMigrationArray, err error) {
	err = s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		txs, err = txStorage.unsafeGetPendingMassMigrations()
		return err
	})
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionStorage) unsafeGetPendingMassMigrations() ([]models.MassMigration, error) {
	txs := make([]models.MassMigration, 0, 32)
	var storedTx models.StoredTx
	err := s.database.Badger.Iterator(models.StoredTxPrefix, db.KeyIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			skip, err := s.getStoredTxFromItem(item, &storedTx)
			if err != nil || skip {
				return false, err
			}
			if storedTx.TxType == txtype.MassMigration {
				txs = append(txs, *storedTx.ToMassMigration(nil))
			}
			return false, nil
		})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}

	sort.SliceStable(txs, func(i, j int) bool {
		return txs[i].Nonce.Cmp(&txs[j].Nonce) < 0
	})

	return txs, nil
}
