package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	bh "github.com/timshannon/badgerhold/v4"
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
		*failedTx.ErrorMessage,
	)

	err = s.database.Badger.Insert(txHash, *stored.NewPendingTx(tx))
	if err != nil {
		return err
	}

	return nil
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

func (s *TransactionStorage) AddTransaction(tx models.GenericTransaction) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		badger := txStorage.database.Badger

		base := tx.GetBase()
		hash := base.Hash

		if base.ErrorMessage != nil {
			// This is a FailedTx

			err := txStorage.checkNoTx(&hash, &stored.PendingTx{})
			if err != nil {
				return err
			}

			err = txStorage.checkNoTx(&hash, &stored.BatchedTx{})
			if err != nil {
				return err
			}

			failedTx := stored.NewFailedTx(tx)
			return badger.Insert(hash, *failedTx)
		} else if base.CommitmentID != nil {
			// This is a BatchedTx

			err := txStorage.checkNoTx(&hash, &stored.PendingTx{})
			if err != nil {
				return err
			}

			err = txStorage.checkNoTx(&hash, &stored.FailedTx{})
			if err != nil {
				return err
			}

			batchedTx := stored.NewBatchedTx(tx)
			return badger.Insert(hash, *batchedTx)
		} else {
			// This is a PendingTx

			err := txStorage.checkNoTx(&hash, &stored.FailedTx{})
			if err != nil {
				return err
			}
			err = txStorage.checkNoTx(&hash, &stored.BatchedTx{})
			if err != nil {
				return err
			}

			pendingTx := stored.NewPendingTx(tx)
			return badger.Insert(pendingTx.Hash, *pendingTx)
		}
	})
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
			err := txStorage.AddTransaction(txs.At(i))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
