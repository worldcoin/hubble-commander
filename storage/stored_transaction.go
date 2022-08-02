package storage

import (
	"sync/atomic"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	bdg "github.com/dgraph-io/badger/v3"
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

func (s *TransactionStorage) MarkTransactionsAsPending(txIDs []models.CommitmentSlot) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txIDs {
			err := txStorage.unsafeMarkTransactionAsPending(&txIDs[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) unsafeMarkTransactionAsPending(txSlot *models.CommitmentSlot) error {
	var batchedTx stored.BatchedTx
	err := s.getAndDelete(*txSlot, &batchedTx)
	if err != nil {
		return errors.WithStack(err)
	}

	s.decrementTransactionCount()

	pendingTx := batchedTx.PendingTx
	err = s.database.Badger.Insert(pendingTx.Hash, pendingTx)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *TransactionStorage) SetTransactionError(txError models.TxError) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		var pendingTx stored.PendingTx
		err := txStorage.getAndDelete(txError.TxHash, &pendingTx)
		if err != nil {
			return err
		}

		failedTx := stored.NewFailedTxFromError(&pendingTx, txError.ErrorMessage)
		err = txStorage.database.Badger.Insert(txError.TxHash, *failedTx)
		if err != nil {
			return err
		}
		return nil
	})
}

// this must be called from inside a transaction
func (s *TransactionStorage) getAndDelete(key, result interface{}) error {
	err := s.database.Badger.Get(key, result)
	if err != nil {
		return errors.Wrap(err, "failed to get stored transaction")
	}

	err = s.database.Badger.Delete(key, result)
	if err != nil {
		return errors.Wrap(err, "failed to delete item")
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
		bh.Where("ID.BatchID").Le(latestBatch.ID),
	)
	if err != nil {
		return nil, err
	}
	return ref.Uint64(count), nil
}

func (s *Storage) initBatchedTxsCounter() (err error) {
	s.batchedTxsCount, err = s.getTransactionCount()
	return err
}

func (s *TransactionStorage) getTransactionIDsByBatchID(batchID models.Uint256) ([]models.CommitmentSlot, error) {
	slots := make([]models.CommitmentSlot, 0, 32)

	// nolint: gocritic
	seekPrefix := append(stored.BatchedTxPrefix, batchID.Bytes()...)

	// BatchedTx are stored with CommitmentSlot as their primary key: BatchID is the
	// first member of CommitmentSlot which means we effectively have an index on
	// BatchID.

	var id models.CommitmentSlot
	err := s.database.Badger.Iterator(seekPrefix, db.KeyIteratorOpts, func(item *bdg.Item) (bool, error) {
		err := db.DecodeKey(item.Key(), &id, stored.BatchedTxPrefix)
		if err != nil {
			return false, err
		}

		slots = append(slots, id)
		return false, nil
	})
	if err != nil && !errors.Is(err, db.ErrIteratorFinished) {
		return nil, err
	}

	return slots, nil
}

func (s *TransactionStorage) GetTransactionIDsByBatchIDs(batchIDs ...models.Uint256) ([]models.CommitmentSlot, error) {
	ids := make([]models.CommitmentSlot, 0, len(batchIDs)*32)

	for i := range batchIDs {
		txIds, err := s.getTransactionIDsByBatchID(batchIDs[i])
		if err != nil {
			return nil, err
		}
		ids = append(ids, txIds...)
	}

	if len(ids) == 0 {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return ids, nil
}

// TODO: can all our callers just call Storage.GetAllMempoolTransactions()
func (s *TransactionStorage) GetPendingTransactions(txType txtype.TransactionType) (models.GenericArray, error) {
	pendingTxs, err := s.tsGetAllMempoolTransactions()
	if err != nil {
		return nil, err
	}

	txs := make([]models.GenericTransaction, len(pendingTxs))
	for i := range pendingTxs {
		txs[i] = pendingTxs[i].ToGenericTransaction()
	}

	return models.MakeGenericArray(txs...), nil
}

func (s *TransactionStorage) GetAllFailedTransactions() (models.GenericArray, error) {
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

func (s *TransactionStorage) MarkTransactionAsIncluded(
	tx models.GenericTransaction,
	commitmentSlot *models.CommitmentSlot,
) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		return txStorage.unsafeMarkTransactionAsIncluded(tx, commitmentSlot)
	})
}

func (s *TransactionStorage) unsafeMarkTransactionAsIncluded(
	tx models.GenericTransaction,
	commitmentSlot *models.CommitmentSlot,
) error {
	// TODO: this used to delete from the pendingTx table. was it a good idea to have
	//       changed that behavior?

	hash := tx.GetBase().Hash
	addressableValue := tx.GetNonce()
	log.WithFields(log.Fields{
		"hash":  hash,
		"from":  tx.GetFromStateID(),
		"nonce": addressableValue.Uint64(),
	}).Debug("marking transaction as included")

	pendingTx := stored.NewPendingTx(tx)
	// this body update is only needed for ToStateID
	pendingTx.Body = stored.NewTxBody(tx)
	batchedTx := stored.NewBatchedTxFromPendingAndCommitment(
		pendingTx, commitmentSlot,
	)
	err := s.insertBatchedTx(batchedTx)
	if err != nil {
		return err
	}
	s.incrementTransactionCount()

	return nil
}

// Note: This method assumes that transactions were included in the commitment in the same
// order as they are given here.
func (s *TransactionStorage) MarkTransactionsAsIncluded(
	txs models.GenericTransactionArray,
	commitmentID *models.CommitmentID,
) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		if txs.Len() > 255 {
			panic("Commitments cannot have more than 255 transactions")
		}

		for i := 0; i < txs.Len(); i++ {
			tx := txs.At(i)
			commitmentSlot := models.NewCommitmentSlot(*commitmentID, uint8(i))

			err := txStorage.unsafeMarkTransactionAsIncluded(tx, commitmentSlot)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) insertBatchedTx(batchedTx *stored.BatchedTx) error {
	key := batchedTx.ID
	err := s.database.Badger.Insert(key, *batchedTx)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *TransactionStorage) getBatchedTxByHash(hash common.Hash) (*stored.BatchedTx, error) {
	var batchedTx stored.BatchedTx
	err := s.database.Badger.FindOneUsingIndex(&batchedTx, hash, "Hash")
	if err != nil {
		return nil, err
	}

	return &batchedTx, nil
}

func (s *TransactionStorage) tsGetMempoolTransactionByHash(hash common.Hash) (*stored.PendingTx, error) {
	// a hack around the fact that GetMempoolTransactionByHash is defined on Storage,
	// but not TransactionStorage. The long-term fix is to lift this method (and all
	// dependencies) out of TransactionStorage and up into Storage; I don't understand
	// why they were separated. Alt: everything in storage/mempool.go ought to be
	// dropped into TransactionStorage

	storage, err := newStorageFromDatabase(s.database)
	if err != nil {
		return nil, err
	}
	return storage.GetMempoolTransactionByHash(hash)
}

func (s *TransactionStorage) tsGetAllMempoolTransactions() ([]stored.PendingTx, error) {
	storage, err := newStorageFromDatabase(s.database)
	if err != nil {
		return nil, err
	}
	return storage.GetAllMempoolTransactions()
}

func (s *TransactionStorage) getTransactionByHash(hash common.Hash) (models.GenericTransaction, error) {
	batchedTx, err := s.getBatchedTxByHash(hash)
	if err == nil {
		return batchedTx.ToGenericTransaction(), nil
	}

	if err != nil && !errors.Is(err, bh.ErrNotFound) {
		return nil, err
	}

	pendingTx, err := s.tsGetMempoolTransactionByHash(hash)
	if err != nil {
		return nil, err
	}
	if pendingTx != nil {
		return pendingTx.ToGenericTransaction(), nil
	}

	var failedTx stored.FailedTx
	err = s.database.Badger.Get(hash, &failedTx)
	if errors.Is(err, bh.ErrNotFound) {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	if err != nil {
		return nil, err
	}
	return failedTx.ToGenericTransaction(), nil
}
