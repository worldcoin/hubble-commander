package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	if base.CommitmentID == nil {
		return result, nil
	}

	batch, err := s.GetBatch(base.CommitmentID.BatchID)
	if err != nil {
		return nil, err
	}

	result.BatchHash = batch.Hash
	result.MinedTime = batch.MinedTime

	return result, nil
}

func (s *TransactionStorage) GetTransactionsByCommitmentID(id models.CommitmentID) (models.GenericTransactionArray, error) {
	batchedTxs := make([]stored.BatchedTx, 0, 32)

	query := bh.Where("CommitmentID").Eq(id).Index("CommitmentID")

	err := s.database.Badger.Find(&batchedTxs, query)
	if err != nil {
		return nil, err
	}

	txs := make(models.GenericArray, 0, len(batchedTxs))
	for i := range batchedTxs {
		txs = append(txs, batchedTxs[i].ToGenericTransaction())
	}

	return txs, nil
}

// returns error if the tranasaction is not a FailedTx
func (s *TransactionStorage) ReplaceFailedTransaction(tx models.GenericTransaction) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		return txStorage.unsafeReplaceFailedTransaction(tx)
	})
}

func (s *TransactionStorage) unsafeReplaceFailedTransaction(tx models.GenericTransaction) error {
	txHash := tx.GetBase().Hash

	_, err := s.getBatchedTxByHash(txHash)
	if err == nil {
		return errors.WithStack(ErrAlreadyMinedTransaction)
	}
	if !errors.Is(err, bh.ErrNotFound) {
		return err
	}

	var failedTx stored.FailedTx
	err = s.getAndDelete(txHash, &failedTx)
	if errors.Is(err, bh.ErrNotFound) {
		return NewNotFoundError("FailedTx")
	}
	if err != nil {
		return err
	}

	// It seems worthwhile to record previous errors somewhere and if we did
	// not log then they would be lost forever
	log.Warnf(
		"Replacing failed transaction. Hash=%x ErrorMessage=%+q",
		txHash,
		failedTx.ErrorMessage,
	)

	err = s.database.Badger.Insert(txHash, *stored.NewPendingTx(tx))
	if err != nil {
		return err
	}

	return nil
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

		err = s.checkNoTx(&hash, &stored.BatchedTx{})
		if err != nil {
			return err
		}

		failedTx := stored.NewFailedTx(tx)
		return badger.Insert(hash, *failedTx)
	} else if base.CommitmentID != nil {
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
		err = badger.Insert(hash, *batchedTx)
		if err != nil {
			return err
		}
		s.incrementTransactionCount()
		return nil
	} else {
		// This is a PendingTx

		err := s.checkNoTx(&hash, &stored.FailedTx{})
		if err != nil {
			return err
		}
		err = s.checkNoTx(&hash, &stored.BatchedTx{})
		if err != nil {
			return err
		}

		pendingTx := stored.NewPendingTx(tx)
		return badger.Insert(pendingTx.Hash, *pendingTx)
	}
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

func (s *TransactionStorage) ReplacePendingTransaction(hash *common.Hash, newTx models.GenericTransaction) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		err := txStorage.database.Badger.Delete(*hash, &stored.PendingTx{})
		if errors.Is(err, bh.ErrNotFound) {
			return errors.WithStack(NewNotFoundError("transaction"))
		}
		if err != nil {
			return err
		}
		return txStorage.database.Badger.Insert(newTx.GetBase().Hash, *stored.NewPendingTx(newTx))
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

func (s *TransactionStorage) BatchUpsertTransaction(txs models.GenericTransactionArray) error {
	if txs.Len() < 1 {
		return errors.WithStack(ErrNoRowsAffected)
	}

	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := 0; i < txs.Len(); i++ {
			err := txStorage.AddTransaction(txs.At(i))
			if errors.Is(err, bh.ErrKeyExists) {
				err = txStorage.MarkTransactionsAsIncluded(models.GenericArray{txs.At(i)}, txs.At(i).GetBase().CommitmentID)
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

func (s *TransactionStorage) AddPendingTransactions(txs models.GenericTransactionArray) error {
	return s.addTxsInMultipleDBTransactions(txs, "pending")
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
