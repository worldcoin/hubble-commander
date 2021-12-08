package storage

import (
	"sort"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (s *TransactionStorage) AddMassMigration(m *models.MassMigration) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		if m.CommitmentID != nil || m.ErrorMessage != nil {
			err := txStorage.addStoredTxReceipt(stored.NewTxReceiptFromMassMigration(m))
			if err != nil {
				return err
			}
		}
		return txStorage.addStoredTx(stored.NewTxFromMassMigration(m))
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
	var storedTx stored.Tx
	err := s.database.Badger.Iterator(stored.TxPrefix, db.KeyIteratorOpts,
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

func (s *TransactionStorage) GetMassMigrationsByCommitmentID(id models.CommitmentID) ([]models.MassMigration, error) {
	massMigrations := make([]models.MassMigration, 0, 1)

	err := s.iterateTxsByCommitmentID(id, func(storedTx *stored.Tx, storedTxReceipt *stored.TxReceipt) {
		if storedTx.TxType == txtype.MassMigration {
			massMigrations = append(massMigrations, *storedTx.ToMassMigration(storedTxReceipt))
		}
	})
	if err != nil {
		return nil, err
	}

	return massMigrations, nil
}

func (s *TransactionStorage) MarkMassMigrationsAsIncluded(txs []models.MassMigration, commitmentID *models.CommitmentID) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			storedTxReceipt := stored.NewTxReceiptFromMassMigration(&txs[i])
			storedTxReceipt.CommitmentID = commitmentID
			err := txStorage.addStoredTxReceipt(storedTxReceipt)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) GetMassMigrationWithBatchDetails(hash common.Hash) (*models.MassMigrationWithBatchDetails, error) {
	var massMigrations []models.MassMigrationWithBatchDetails
	err := s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		tx, txReceipt, err := txStorage.getStoredTxWithReceipt(hash)
		if err != nil {
			return err
		}
		if tx.TxType != txtype.MassMigration {
			return errors.WithStack(NewNotFoundError("transaction"))
		}

		massMigrations, err = txStorage.massMigrationsToMassMigrationsWithBatchDetails(
			[]stored.Tx{*tx},
			[]*stored.TxReceipt{txReceipt},
		)
		return err
	})

	if err != nil {
		return nil, err
	}
	return &massMigrations[0], nil
}

func (s *Storage) massMigrationsToMassMigrationsWithBatchDetails(txs []stored.Tx, txReceipts []*stored.TxReceipt) (
	result []models.MassMigrationWithBatchDetails,
	err error,
) {
	result = make([]models.MassMigrationWithBatchDetails, 0, len(txs))
	batchIDs := make(map[models.Uint256]*models.Batch)
	for i := range txs {
		massMigration := txs[i].ToMassMigration(txReceipts[i])
		if massMigration.CommitmentID == nil {
			result = append(result, models.MassMigrationWithBatchDetails{MassMigration: *massMigration})
			continue
		}
		batch, ok := batchIDs[massMigration.CommitmentID.BatchID]
		if !ok {
			batch, err = s.GetBatch(massMigration.CommitmentID.BatchID)
			if err != nil {
				return nil, err
			}
			batchIDs[massMigration.CommitmentID.BatchID] = batch
		}

		result = append(result, models.MassMigrationWithBatchDetails{
			MassMigration: *massMigration,
			BatchHash:     batch.Hash,
			BatchTime:     batch.SubmissionTime,
		})
	}
	return result, nil
}
