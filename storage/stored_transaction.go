package storage

import (
	"bytes"
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
	// We need to "initialize" the indices on fields of pointer type to make them work with bh.Find operations.
	// The problem originates in `indexExists` function in BadgerHold (https://github.com/timshannon/badgerhold/blob/v4.0.1/index.go#L148).
	// Badger assumes that if there is a value for some data type, then there must exist at least one index entry.
	// If you don't index nil values the way we did for stored.TxReceipt.ToStateID it can be the case that there is some
	// StoredTxReceipt stored, but there is no index entry. To work around this we set an empty index entry.
	// See:
	// 	 * stored.TxReceipt Indexes() method
	//   * StoredTransactionTestSuite.TestStoredTxReceipt_FindUsingIndexWorksWhenThereAreOnlyStoredTxReceiptsWithNilToStateID
	err := initializeIndex(database, stored.TxReceiptName, "ToStateID", uint32(0))
	if err != nil {
		return nil, err
	}
	err = initializeIndex(database, stored.TxReceiptName, "CommitmentID", models.CommitmentID{
		BatchID:      models.MakeUint256(0),
		IndexInBatch: 0,
	})
	if err != nil {
		return nil, err
	}

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

func (s *TransactionStorage) addStoredTxReceipt(txReceipt *stored.TxReceipt) error {
	return s.database.Badger.Insert(txReceipt.Hash, *txReceipt)
}

func (s *TransactionStorage) getStoredTxWithReceipt(hash common.Hash) (storedTx *stored.Tx, txReceipt *stored.TxReceipt, err error) {
	err = s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		storedTx, err = txStorage.getStoredTx(hash)
		if err != nil {
			return err
		}

		txReceipt, err = txStorage.getStoredTxReceipt(hash)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return storedTx, txReceipt, nil
}

func (s *TransactionStorage) addStoredTx(tx *stored.Tx) error {
	return s.database.Badger.Insert(tx.Hash, *tx)
}

func (s *TransactionStorage) updateStoredTx(tx *stored.Tx) error {
	return s.database.Badger.Update(tx.Hash, *tx)
}

func (s *TransactionStorage) getStoredTx(hash common.Hash) (*stored.Tx, error) {
	var storedTx stored.Tx
	err := s.database.Badger.Get(hash, &storedTx)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	if err != nil {
		return nil, err
	}
	return &storedTx, nil
}

func (s *TransactionStorage) getStoredTxReceipt(hash common.Hash) (*stored.TxReceipt, error) {
	var storedTxReceipt stored.TxReceipt
	err := s.database.Badger.Get(hash, &storedTxReceipt)
	if err == bh.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &storedTxReceipt, nil
}

func (s *TransactionStorage) MarkTransactionsAsPending(txHashes []common.Hash) error {
	dataType := stored.TxReceipt{}
	return s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		for i := range txHashes {
			err := txStorage.database.Badger.Delete(txHashes[i], dataType)
			if err != nil {
				return err
			}
		}
		return nil
	})
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
			return txStorage.addStoredTxReceipt(&stored.TxReceipt{
				Hash:         txError.TxHash,
				ErrorMessage: &txError.ErrorMessage,
			})
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
		seekPrefix := db.IndexKeyPrefix(stored.TxReceiptName, "CommitmentID")
		err = txStorage.database.Badger.Iterator(seekPrefix, db.PrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
			var commitmentID models.CommitmentID
			err = db.Decode(keyValue(seekPrefix, item.Key()), &commitmentID)
			if err != nil {
				return false, err
			}
			if commitmentID.BatchID.Cmp(&latestBatch.ID) > 0 {
				return false, nil
			}

			var keyList bh.KeyList
			err = item.Value(func(val []byte) error {
				return db.Decode(val, &keyList)
			})
			if err != nil {
				return false, err
			}
			count += len(keyList)
			return false, nil
		})
		if err != nil && err != db.ErrIteratorFinished {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (s *TransactionStorage) GetTransactionHashesByBatchIDs(batchIDs ...models.Uint256) ([]common.Hash, error) {
	batchPrefixes := batchIdsToBatchPrefixes(batchIDs)
	hashes := make([]common.Hash, 0, len(batchIDs)*32)

	var keyList bh.KeyList
	seekPrefix := db.IndexKeyPrefix(stored.TxReceiptName, "CommitmentID")
	err := s.database.Badger.Iterator(seekPrefix, db.ReversePrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
		if validForPrefixes(keyValue(seekPrefix, item.Key()), batchPrefixes) {
			err := item.Value(func(val []byte) error {
				return db.Decode(val, &keyList)
			})
			if err != nil {
				return false, err
			}
			txHashes, err := decodeKeyListHashes(stored.TxReceiptPrefix, keyList)
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
	switch txType {
	case txtype.Transfer:
		return s.GetPendingTransfers()
	case txtype.Create2Transfer:
		return s.GetPendingCreate2Transfers()
	case txtype.MassMigration:
		return s.GetPendingMassMigrations()
	}
	return nil, nil
}

func (s *TransactionStorage) MarkTransactionsAsIncluded(txs models.GenericTransactionArray, commitmentID *models.CommitmentID) error {
	switch txs.Type() {
	case txtype.Transfer:
		return s.MarkTransfersAsIncluded(txs.ToTransferArray(), commitmentID)
	case txtype.Create2Transfer:
		return s.MarkCreate2TransfersAsIncluded(txs.ToCreate2TransferArray(), commitmentID)
	case txtype.MassMigration:
		return s.MarkMassMigrationsAsIncluded(txs.ToMassMigrationArray(), commitmentID)
	}
	return nil
}

func (s *TransactionStorage) getStoredTxFromItem(item *bdg.Item, storedTx *stored.Tx) (bool, error) {
	var hash common.Hash
	err := db.DecodeKey(item.Key(), &hash, stored.TxPrefix)
	if err != nil {
		return false, err
	}
	txReceipt, err := s.getStoredTxReceipt(hash)
	if err != nil {
		return false, err
	}
	if txReceipt != nil {
		return true, nil
	}

	return false, item.Value(storedTx.SetBytes)
}

func decodeKeyListHashes(keyPrefix []byte, keyList bh.KeyList) ([]common.Hash, error) {
	var hash common.Hash
	hashes := make([]common.Hash, 0, len(keyList))
	for i := range keyList {
		err := stored.DecodeHash(keyList[i][len(keyPrefix):], &hash)
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
