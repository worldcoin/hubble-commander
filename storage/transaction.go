package storage

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *Storage) GetTransactionWithBatchDetails(hash common.Hash) (tx *models.TransactionWithBatchDetails, err error) {
	err = s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		tx, err = txStorage.unsafeGetTransactionWithBatchDetails(hash)
		return err
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *Storage) unsafeGetTransactionWithBatchDetails(hash common.Hash) (
	*models.TransactionWithBatchDetails,
	error,
) {
	generic, err := s.getTransactionByHash(hash)
	if err != nil {
		return nil, err
	}

	result := &models.TransactionWithBatchDetails{Transaction: generic}

	base := generic.GetBase()
	if base.CommitmentSlot == nil {
		return result, nil
	}

	batch, err := s.GetBatch(base.CommitmentSlot.BatchID)
	if err != nil {
		return nil, err
	}

	result.BatchHash = batch.Hash
	result.MinedTime = batch.MinedTime

	return result, nil
}

func (s *TransactionStorage) GetTransactionsByCommitmentID(id models.CommitmentID) (models.GenericTransactionArray, error) {
	batchedTxs := make([]stored.BatchedTx, 0, 32)

	// BatchedTx are stored with CommitmentSlot as their primary key:
	// - CommitmentSlot is (BatchID, IndexInBatch, IndexInCommitment)
	// - CommitmentID   is (BatchId, IndexInBatch)
	// This means that if a CommitmentID and a CommitmentSlot represent the same
	// commitment then the serialization of the CommitmentID will be a prefix of the
	// serialization of the CommitmentSlot. This makes it easy to look for all
	// BatchedTxs for `id`

	// nolint: gocritic
	seekPrefix := append(stored.BatchedTxPrefix, id.Bytes()...)

	err := s.database.Badger.Iterator(seekPrefix, db.PrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
		var batchedTx stored.BatchedTx
		err := item.Value(batchedTx.SetBytes)
		if err != nil {
			return db.Continue, err
		}

		batchedTxs = append(batchedTxs, batchedTx)
		return db.Continue, nil
	})
	if err != nil && !errors.Is(err, db.ErrIteratorFinished) {
		return nil, err
	}

	txs := make(models.GenericArray, 0, len(batchedTxs))
	for i := range batchedTxs {
		txs = append(txs, batchedTxs[i].ToGenericTransaction())
	}

	return txs, nil
}

func (s *TransactionStorage) AddTransaction(tx models.GenericTransaction) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		return txStorage.unsafeAddTransaction(tx)
	})
}

func (s *TransactionStorage) unsafeAddTransaction(tx models.GenericTransaction) error {
	badger := s.database.Badger

	base := tx.GetBase()
	hash := base.Hash

	if base.ErrorMessage != nil {
		// This is a FailedTx

		err := s.checkNoTx(&hash, &stored.PendingTx{})
		if err != nil {
			return err
		}

		err = s.checkNoBatchedTx(&hash)
		if err != nil {
			return err
		}

		failedTx := stored.NewFailedTx(tx)
		return badger.Insert(hash, *failedTx)
	} else if base.CommitmentSlot != nil {
		// This is a BatchedTx

		err := s.checkNoTx(&hash, &stored.PendingTx{})
		if err != nil {
			return err
		}

		err = s.checkNoTx(&hash, &stored.FailedTx{})
		if err != nil {
			return err
		}

		batchedTx := stored.NewBatchedTx(tx)
		err = s.insertBatchedTx(batchedTx)
		if err != nil {
			return err
		}
		s.incrementTransactionCount()
		return nil
	} else {
		// This is a PendingTx

		return errors.WithStack(
			fmt.Errorf("Use AddMempoolTx to insert pending txns"),
		)
	}
}

func (s *TransactionStorage) checkNoBatchedTx(hash *common.Hash) error {
	_, err := s.getBatchedTxByHash(*hash)
	if errors.Is(err, bh.ErrNotFound) {
		// there is no tx, so we're free to insert a tx!
		return nil
	}
	if err == nil {
		// we successfully fetched a tx, so our caller should fail
		return bh.ErrKeyExists
	}
	return err
}

func (s *TransactionStorage) checkNoTx(hash *common.Hash, result interface{}) error {
	err := s.database.Badger.Get(*hash, result)
	if errors.Is(err, bh.ErrNotFound) {
		// there is no tx, so we're free to insert a tx!
		return nil
	}
	if err == nil {
		// we successfully fetched a tx, so our caller should fail
		return bh.ErrKeyExists
	}
	return err
}

func (s *TransactionStorage) BatchAddTransaction(txs models.GenericTransactionArray) error {
	if txs.Len() < 1 {
		return errors.WithStack(ErrNoRowsAffected)
	}

	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := 0; i < txs.Len(); i++ {
			// TODO: this can call unsafeAddTransaction directly, we're
			//       already inside a txn
			err := txStorage.AddTransaction(txs.At(i))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// TODO: This needs to be fixed.
//       It used to take pending/failed transactions and add them to batches (as a result
//       of the sync process.
func (s *TransactionStorage) BatchUpsertTransaction(txs models.GenericTransactionArray) error {
	if txs.Len() < 1 {
		return errors.WithStack(ErrNoRowsAffected)
	}

	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := 0; i < txs.Len(); i++ {
			err := txStorage.AddTransaction(txs.At(i))
			if errors.Is(err, bh.ErrKeyExists) {
				err = txStorage.MarkTransactionAsIncluded(
					txs.At(i), txs.At(i).GetBase().CommitmentSlot,
				)
				if err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) AddFailedTransactions(txs models.GenericTransactionArray) error {
	return s.addTxsInMultipleDBTransactions(txs, "failed")
}

func (s *TransactionStorage) addTxsInMultipleDBTransactions(txs models.GenericTransactionArray, status string) error {
	if txs.Len() == 0 {
		return nil
	}

	operations := make([]dbOperation, 0, txs.Len())
	for i := 0; i < txs.Len(); i++ {
		tx := txs.At(i)
		operations = append(operations, func(txStorage *TransactionStorage) error {
			return txStorage.AddTransaction(tx)
		})
	}

	dbTxsCount, err := s.updateInMultipleTransactions(operations)
	if err != nil {
		return errors.Wrapf(err, "storing %d %s tx(s) failed during database transaction #%d", txs.Len(), status, dbTxsCount)
	}
	return nil
}
