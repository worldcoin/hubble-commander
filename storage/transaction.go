package storage

import (
	"bytes"
	"fmt"
	"sync/atomic"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	bh "github.com/timshannon/badgerhold/v4"
)

type TransactionStorage struct {
	database *Database

	batchedTxsCount *uint64
}

type dbOperation func(txStorage *TransactionStorage) error

func NewTransactionStorage(database *Database) *TransactionStorage {
	return &TransactionStorage{
		database:        database,
		batchedTxsCount: ref.Uint64(0),
	}
}

func (s *TransactionStorage) copyWithNewDatabase(database *Database) *TransactionStorage {
	newTransactionStorage := *s
	newTransactionStorage.database = database

	return &newTransactionStorage
}

func (s *Storage) initBatchedTxsCounter() (err error) {
	s.batchedTxsCount, err = s.getTransactionCount()
	return err
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
		return s.unsafeAddTransaction(tx)
	})
}

func (s *TransactionStorage) unsafeAddTransaction(tx models.GenericTransaction) error {
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
		return s.database.Badger.Insert(hash, *failedTx)
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
		err = s.database.Badger.Insert(hash, *batchedTx)
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
		return s.database.Badger.Insert(pendingTx.Hash, *pendingTx)
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
		err := s.database.Badger.Delete(*hash, &stored.PendingTx{})
		if errors.Is(err, bh.ErrNotFound) {
			return errors.WithStack(NewNotFoundError("transaction"))
		}
		if err != nil {
			return err
		}
		return s.database.Badger.Insert(newTx.GetBase().Hash, *stored.NewPendingTx(newTx))
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
				err = s.MarkTransactionsAsIncluded(models.GenericArray{txs.At(i)}, txs.At(i).GetBase().CommitmentID)
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

// Be careful. For these operations to be correctly spread across multiple transactions:
// (1) they must popagate up any badger errors they encounter (wrapping is okay)
// (2) they must be idempotent, because they might be retried.
func (s *TransactionStorage) updateInMultipleTransactions(operations []dbOperation) (txCount uint, err error) {
	txController, txStorage := s.beginTransaction(TxOptions{})
	defer txController.Rollback(&err)
	txCount = 1

	for _, op := range operations {
		err = op(txStorage)
		if errors.Is(err, badger.ErrTxnTooBig) {
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
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txHashes {
			err := txStorage.unsafeMarkTransactionAsPending(&txHashes[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) unsafeMarkTransactionAsPending(txHash *common.Hash) error {
	var pendingTx stored.PendingTx

	var batchedTx stored.BatchedTx
	err := s.getAndDelete(*txHash, &batchedTx)
	if err == nil {
		pendingTx = batchedTx.PendingTx
		s.decrementTransactionCount()
	} else {
		var failedTx stored.FailedTx
		err = s.getAndDelete(*txHash, &failedTx)
		if err != nil {
			return err
		}
		pendingTx = failedTx.PendingTx
	}

	return s.database.Badger.Insert(*txHash, pendingTx)
}

func (s *TransactionStorage) SetTransactionError(txError models.TxError) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		var pendingTx stored.PendingTx
		err := txStorage.getAndDelete(txError.TxHash, &pendingTx)
		if err != nil {
			return err
		}

		failedTx := stored.NewFailedTxFromError(&pendingTx, txError.ErrorMessage)
		return txStorage.database.Badger.Insert(txError.TxHash, *failedTx)
	})
}

func (s *TransactionStorage) getAndDelete(key, result interface{}) error {
	err := s.database.Badger.Get(key, result)
	if err != nil {
		return fmt.Errorf("failed to Get item: %w", err)
	}

	err = s.database.Badger.Delete(key, result)
	if err != nil {
		return fmt.Errorf("failed to Delete item: %w", err)
	}

	return nil
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
		return errors.Wrapf(err, "storing %d tx error(s) failed during database transaction #%d", errorsCount, dbTxsCount)
	}
	log.Debugf("Stored %d tx error(s) in %d database transaction(s)", errorsCount, dbTxsCount)
	return nil
}

func (s *TransactionStorage) GetTransactionCount() uint64 {
	return atomic.LoadUint64(s.batchedTxsCount)
}

func (s *TransactionStorage) incrementTransactionCount() {
	atomic.AddUint64(s.batchedTxsCount, 1)
}

func (s *TransactionStorage) decrementTransactionCount() {
	atomic.AddUint64(s.batchedTxsCount, ^uint64(0))
}

func (s *Storage) getTransactionCount() (count *uint64, err error) {
	err = s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		count, err = txStorage.unsafeGetTransactionCount()
		return err
	})
	if err != nil {
		return nil, err
	}
	return count, nil
}

func (s *Storage) unsafeGetTransactionCount() (*uint64, error) {
	latestBatch, err := s.GetLatestSubmittedBatch()
	if IsNotFoundError(err) {
		return ref.Uint64(0), nil
	}
	if err != nil {
		return nil, err
	}

	count, err := s.database.Badger.Count(
		&stored.BatchedTx{},
		bh.Where("CommitmentID.BatchID").Le(latestBatch.ID),
	)
	if err != nil {
		return nil, err
	}
	return ref.Uint64(count), nil
}

func (s *TransactionStorage) GetTransactionHashesByBatchIDs(batchIDs ...models.Uint256) ([]common.Hash, error) {
	hashes := make([]common.Hash, 0, len(batchIDs)*32)

	// We have an index on CommitmentID. It turns out that BatchID is the first
	// member of the BatchID struct, so we effectively have an index on BatchID.
	// We can take advantage of that index by manually iterating.
	// The slow version is: Badger.Find(..., bh.Where("CommitmentID.BatchID").In(...))

	var keyList bh.KeyList
	batchPrefixes := batchIdsToBatchPrefixes(batchIDs)
	seekPrefix := db.IndexKeyPrefix(stored.BatchedTxName, "CommitmentID")
	err := s.database.Badger.Iterator(seekPrefix, db.ReversePrefetchIteratorOpts, func(item *badger.Item) (bool, error) {
		if validForPrefixes(keyValue(seekPrefix, item.Key()), batchPrefixes) {
			err := item.Value(func(val []byte) error {
				return db.Decode(val, &keyList)
			})
			if err != nil {
				return false, err
			}
			txHashes, err := decodeKeyListHashes(stored.BatchedTxPrefix, keyList)
			if err != nil {
				return false, err
			}
			hashes = append(hashes, txHashes...)
		}
		return false, nil
	})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}
	if len(hashes) == 0 {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return hashes, nil
}

func (s *TransactionStorage) GetPendingTransactions(txType txtype.TransactionType) (models.GenericTransactionArray, error) {
	var pendingTxs []stored.PendingTx

	err := s.database.Badger.Find(&pendingTxs, bh.Where("TxType").Eq(txType))
	if err != nil {
		return nil, err
	}

	txs := make([]models.GenericTransaction, len(pendingTxs))
	for i := range pendingTxs {
		txs[i] = pendingTxs[i].ToGenericTransaction()
	}

	return models.MakeGenericArray(txs...), nil
}

func (s *TransactionStorage) GetAllPendingTransactions() (models.GenericTransactionArray, error) {
	var pendingTxs []stored.PendingTx
	err := s.database.Badger.Find(&pendingTxs, &bh.Query{})
	if err != nil {
		return nil, err
	}

	txs := make([]models.GenericTransaction, len(pendingTxs))
	for i := range pendingTxs {
		txs[i] = pendingTxs[i].ToGenericTransaction()
	}

	return models.MakeGenericArray(txs...), nil
}

func (s *TransactionStorage) GetAllFailedTransactions() (models.GenericTransactionArray, error) {
	var failedTxs []stored.FailedTx
	err := s.database.Badger.Find(&failedTxs, nil)
	if err != nil {
		return nil, err
	}

	txs := make([]models.GenericTransaction, len(failedTxs))
	for i := range failedTxs {
		txs[i] = failedTxs[i].ToGenericTransaction()
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
			s.incrementTransactionCount()
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

func decodeKeyListHashes(prefix []byte, keyList bh.KeyList) ([]common.Hash, error) {
	var hash common.Hash
	hashes := make([]common.Hash, 0, len(keyList))
	for i := range keyList {
		err := stored.DecodeHash(keyValue(prefix, keyList[i]), &hash)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

func batchIdsToBatchPrefixes(batchIDs []models.Uint256) [][]byte {
	batchPrefixes := make([][]byte, 0, len(batchIDs))
	for i := range batchIDs {
		batchPrefixes = append(batchPrefixes, batchIDs[i].Bytes())
	}
	return batchPrefixes
}

func keyValue(prefix, key []byte) []byte {
	return key[len(prefix):]
}

func validForPrefixes(s []byte, prefixes [][]byte) bool {
	for i := range prefixes {
		if bytes.HasPrefix(s, prefixes[i]) {
			return true
		}
	}
	return false
}
