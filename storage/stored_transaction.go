package storage

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	bh "github.com/timshannon/badgerhold/v4"
)

type TransactionStorage struct {
	database *Database
}

type dbOperation func(txStorage *TransactionStorage) error

func NewTransactionStorage(database *Database) (*TransactionStorage, error) {
	return &TransactionStorage{
		database: database,
	}, nil
}

func (s *TransactionStorage) copyWithNewDatabase(database *Database) *TransactionStorage {
	newTransactionStorage := *s
	newTransactionStorage.database = database

	return &newTransactionStorage
}

func (s *TransactionStorage) beginTransaction(opts TxOptions) (*db.TxController, *TransactionStorage) {
	txController, txDatabase := s.database.BeginTransaction(opts)
	return txController, s.copyWithNewDatabase(txDatabase)
}

func (s *TransactionStorage) executeInTransaction(opts TxOptions, fn func(txStorage *TransactionStorage) error) error {
	return s.database.ExecuteInTransaction(opts, func(txDatabase *Database) error {
		return fn(s.copyWithNewDatabase(txDatabase))
	})
}

// Be careful. For these operations to be correctly spread across multiple transactions:
// (1) they must popagate up any badger errors they encounter (wrapping is okay)
// (2) they must be idempotent, because they might be retried.
func (s *TransactionStorage) updateInMultipleTransactions(operations []dbOperation) (txCount uint, err error) {
	txController, txStorage := s.beginTransaction(TxOptions{})
	defer txController.Rollback(&err)
	txCount = 1

	for _, op := range operations {
		err = op(txStorage)
		if errors.Is(err, bdg.ErrTxnTooBig) {
			// Commit and start new DB transaction
			err = txController.Commit()
			if err != nil {
				return txCount, err
			}
			txController, txStorage = s.beginTransaction(TxOptions{})
			txCount++

			// Retry operation
			err = op(txStorage)
		}
		if err != nil {
			// Either the error was different than bdg.ErrTxnTooBig or retry failed
			return txCount, err
		}
	}

	return txCount, txController.Commit()
}

func (s *TransactionStorage) MarkTransactionsAsPending(txHashes []common.Hash) error {
	return s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		for i := range txHashes {
			var pendingTx stored.PendingTx

			var batchedTx stored.BatchedTx
			err := txStorage.getAndDelete(txHashes[i], &batchedTx)
			if err == nil {
				pendingTx = batchedTx.PendingTx
			} else {
				var failedTx stored.FailedTx
				err = txStorage.getAndDelete(txHashes[i], &failedTx)
				if err != nil {
					return err
				}
				pendingTx = failedTx.PendingTx
			}

			err = txStorage.database.Badger.Insert(txHashes[i], pendingTx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) SetTransactionError(txError models.TxError) error {
	return s.executeInTransaction(TxOptions{ReadOnly: false}, func(txStorage *TransactionStorage) error {
		var pendingTx stored.PendingTx
		err := txStorage.getAndDelete(txError.TxHash, &pendingTx)
		if err != nil {
			return err
		}

		failedTx := stored.NewFailedTxFromError(&pendingTx, &txError.ErrorMessage)
		return txStorage.database.Badger.Insert(txError.TxHash, *failedTx)
	})
}

func (s *TransactionStorage) getAndDelete(key, result interface{}) error {
	// TODO: wrap these errors so it's easier to debug them,
	//       currently there's no way to know where the returned
	//       error came from
	err := s.database.Badger.Get(key, result)
	if err != nil {
		return err
	}

	return s.database.Badger.Delete(key, result)
}

func (s *TransactionStorage) SetTransactionErrors(txErrors ...models.TxError) error {
	errorsCount := len(txErrors)
	if errorsCount == 0 {
		return nil
	}

	operations := make([]dbOperation, errorsCount)
	for i := range txErrors {
		txError := txErrors[i]
		operations[i] = func(txStorage *TransactionStorage) error {
			return txStorage.SetTransactionError(txError)
		}
	}

	dbTxsCount, err := s.updateInMultipleTransactions(operations)
	if err != nil {
		err = fmt.Errorf("storing %d tx error(s) failed during database transaction #%d because of: %w", errorsCount, dbTxsCount, err)
		return errors.WithStack(err)
	}
	log.Debugf("Stored %d tx error(s) in %d database transaction(s)", errorsCount, dbTxsCount)
	return nil
}

func (s *Storage) GetTransactionCount() (*int, error) {
	count := 0
	err := s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		latestBatch, err := txStorage.GetLatestSubmittedBatch()
		if IsNotFoundError(err) {
			return nil
		}
		if err != nil {
			return err
		}

		bhcount, err := txStorage.database.Badger.Count(
			&stored.BatchedTx{},
			bh.Where("CommitmentID.BatchID").Le(latestBatch.ID),
		)
		if err == nil {
			count = int(bhcount)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (s *TransactionStorage) GetTransactionHashesByBatchIDs(batchIDs ...models.Uint256) ([]common.Hash, error) {
	hashes := make([]common.Hash, 0, len(batchIDs)*32)

	var batchedTxs []stored.BatchedTx
	err := s.database.Badger.Find(&batchedTxs, bh.Where("CommitmentID.BatchID").In(bh.Slice(batchIDs)...))
	if err != nil {
		return nil, err
	}

	for i := range batchedTxs {
		hashes = append(hashes, batchedTxs[i].Hash)
	}
	if len(hashes) == 0 {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return hashes, err
}

func (s *TransactionStorage) GetPendingTransactions(txType txtype.TransactionType) (txs models.GenericTransactionArray, err error) {
	err = s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		txs, err = txStorage.unsafeGetPendingTransactions(txType)
		return err
	})
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionStorage) unsafeGetPendingTransactions(txType txtype.TransactionType) (models.GenericTransactionArray, error) {
	var pendingTxs []stored.PendingTx

	err := s.database.Badger.Find(&pendingTxs, bh.Where("TxType").Eq(txType))
	if err != nil {
		return nil, err
	}

	// TODO: the tests never exercise the case where there are multiple pending txns
	//       of different types

	txs := make([]models.GenericTransaction, len(pendingTxs))
	for i := range pendingTxs {
		txs[i] = pendingTxs[i].ToGenericTransaction()
	}

	return models.MakeGenericArray(txs...), nil
}

func (s *TransactionStorage) MarkTransactionsAsIncluded(
	txs models.GenericTransactionArray,
	commitmentID *models.CommitmentID,
) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := 0; i < txs.Len(); i++ {
			tx := txs.At(i)

			// note: the rest of the txn is ignored. We take the existing txn
			//       in our database and record which commitment it belongs
			//       to. We assume that the executed transaction does not
			//       differ from our local records.
			hash := tx.GetBase().Hash

			var pendingTx stored.PendingTx
			err := txStorage.getAndDelete(hash, &pendingTx)
			if err != nil {
				return err
			}

			// this body update is only needed for ToStateID
			pendingTx.Body = stored.NewTxBody(tx)
			batchedTx := stored.NewBatchedTxFromPendingAndCommitment(
				&pendingTx, commitmentID,
			)
			err = txStorage.database.Badger.Insert(hash, *batchedTx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) getBatchedTxByHash(hash common.Hash) (*stored.BatchedTx, error) {
	var batchedTx stored.BatchedTx
	err := s.database.Badger.Get(hash, &batchedTx)
	if err != nil {
		return nil, err
	}

	return &batchedTx, nil
}

func (s *TransactionStorage) getTransactionByHash(hash common.Hash) (models.GenericTransaction, error) {
	batchedTx, err := s.getBatchedTxByHash(hash)
	if err == nil {
		return batchedTx.ToGenericTransaction(), nil
	}

	if err != nil && !errors.Is(err, bh.ErrNotFound) {
		return nil, err
	}

	var pendingTx stored.PendingTx
	err = s.database.Badger.Get(hash, &pendingTx)
	if err == nil {
		return pendingTx.ToGenericTransaction(), nil
	}
	if err != nil && !errors.Is(err, bh.ErrNotFound) {
		return nil, err
	}

	var failedTx stored.FailedTx
	err = s.database.Badger.Get(hash, &failedTx)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	if err != nil {
		return nil, err
	}
	return failedTx.ToGenericTransaction(), nil
}
