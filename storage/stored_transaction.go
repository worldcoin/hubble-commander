package storage

import (
	"bytes"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type TransactionStorage struct {
	database *Database
}

func NewTransactionStorage(database *Database) (*TransactionStorage, error) {
	// We need to "initialize" the indices on fields of pointer type to make them work with bh.Find operations.
	// The problem originates in `indexExists` function in BadgerHold (https://github.com/timshannon/badgerhold/blob/v4.0.1/index.go#L148).
	// Badger assumes that if there is a value for some data type, then there must exist at least one index entry.
	// If you don't index nil values the way we did for models.StoredTxReceipt.ToStateID it can be the case that there is some
	// StoredTxReceipt stored, but there is no index entry. To work around this we set an empty index entry.
	// See:
	// 	 * models.StoredTxReceipt Indexes() method
	//   * StoredTransactionTestSuite.TestStoredTxReceipt_FindUsingIndexWorksWhenThereAreOnlyStoredTxReceiptsWithNilToStateID
	err := initializeIndex(database, models.StoredTxReceiptName, "ToStateID", uint32(0))
	if err != nil {
		return nil, err
	}

	return &TransactionStorage{
		database: database,
	}, nil
}

func initializeIndex(database *Database, typeName []byte, indexName string, zeroValue interface{}) error {
	encodedZeroValue, err := db.Encode(zeroValue)
	if err != nil {
		return err
	}
	zeroValueIndexKey := db.IndexKey(typeName, indexName, encodedZeroValue)

	emptyKeyList := make(bh.KeyList, 0)
	encodedEmptyKeyList, err := db.Encode(emptyKeyList)
	if err != nil {
		return err
	}

	return database.Badger.RawUpdate(func(txn *bdg.Txn) error {
		return txn.Set(zeroValueIndexKey, encodedEmptyKeyList)
	})
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

func (s *TransactionStorage) addStoredTxReceipt(txReceipt *models.StoredTxReceipt) error {
	return s.database.Badger.Insert(txReceipt.Hash, *txReceipt)
}

func (s *TransactionStorage) getStoredTxWithReceipt(hash common.Hash) (
	storedTx *models.StoredTx, txReceipt *models.StoredTxReceipt, err error,
) {
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

func (s *TransactionStorage) addStoredTx(tx *models.StoredTx) error {
	return s.database.Badger.Insert(tx.Hash, *tx)
}

func (s *TransactionStorage) getStoredTx(hash common.Hash) (*models.StoredTx, error) {
	var storedTx models.StoredTx
	err := s.database.Badger.Get(hash, &storedTx)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	if err != nil {
		return nil, err
	}
	return &storedTx, nil
}

func (s *TransactionStorage) getStoredTxReceipt(hash common.Hash) (*models.StoredTxReceipt, error) {
	var storedTxReceipt models.StoredTxReceipt
	err := s.database.Badger.Get(hash, &storedTxReceipt)
	if err == bh.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &storedTxReceipt, nil
}

func (s *TransactionStorage) getKeyList(indexKey []byte) (*bh.KeyList, error) {
	var keyList bh.KeyList
	err := s.database.Badger.View(func(txn *bdg.Txn) error {
		item, err := txn.Get(indexKey)
		if err == bdg.ErrKeyNotFound {
			return errors.WithStack(NewNotFoundError("transaction"))
		}
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return db.DecodeKeyList(val, &keyList)
		})
	})
	if err != nil {
		return nil, err
	}
	return &keyList, nil
}

func (s *TransactionStorage) GetLatestTransactionNonce(accountStateID uint32) (*models.Uint256, error) {
	var latestNonce *models.Uint256

	err := s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		indexKey := db.IndexKey(models.StoredTxName, "FromStateID", models.EncodeUint32(accountStateID))
		keyList, err := txStorage.getKeyList(indexKey)
		if err != nil {
			return err
		}
		txHashes, err := decodeKeyListHashes(models.StoredTxPrefix, *keyList)
		if err != nil {
			return err
		}
		if len(txHashes) == 0 {
			return errors.WithStack(NewNotFoundError("transaction"))
		}

		for i := range txHashes {
			tx, receipt, err := txStorage.getStoredTxWithReceipt(txHashes[i])
			if err != nil {
				return err
			}
			if receipt == nil && (latestNonce == nil || tx.Nonce.Cmp(latestNonce) > 0) {
				latestNonce = &tx.Nonce
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if latestNonce == nil {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return latestNonce, nil
}

func (s *TransactionStorage) MarkTransactionsAsPending(txHashes []common.Hash) error {
	dataType := models.StoredTxReceipt{}
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

func (s *TransactionStorage) SetTransactionErrors(txErrors ...models.TxError) (err error) {
	if len(txErrors) == 0 {
		return nil
	}

	txController, txStorage := s.beginTransaction(TxOptions{})
	defer txController.Rollback(&err)

	for i := range txErrors {
		err = txStorage.addStoredTxReceipt(&models.StoredTxReceipt{
			Hash:         txErrors[i].TxHash,
			ErrorMessage: &txErrors[i].ErrorMessage,
		})
		if err == bdg.ErrTxnTooBig {
			err = txController.Commit()
			if err != nil {
				return err
			}
			txController, txStorage = s.beginTransaction(TxOptions{})
		}
		if err != nil {
			return err
		}
	}
	return txController.Commit()
}

func (s *Storage) GetTransactionCount() (*int, error) {
	count := new(int)
	err := s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		latestBatch, err := txStorage.GetLatestSubmittedBatch()
		if IsNotFoundError(err) {
			return nil
		}
		if err != nil {
			return err
		}
		seekPrefix := db.IndexKeyPrefix(models.StoredTxReceiptName, "CommitmentID")
		err = txStorage.database.Badger.Iterator(seekPrefix, db.PrefetchIteratorOpts,
			func(item *bdg.Item) (bool, error) {
				var commitmentID *models.CommitmentID
				commitmentID, err = models.DecodeCommitmentIDPointer(keyValue(seekPrefix, item.Key()))
				if err != nil {
					return false, err
				}
				if commitmentID == nil || commitmentID.BatchID.Cmp(&latestBatch.ID) > 0 {
					return false, nil
				}

				var keyList bh.KeyList
				err = item.Value(func(val []byte) error {
					return db.DecodeKeyList(val, &keyList)
				})
				if err != nil {
					return false, err
				}
				*count += len(keyList)
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
	return count, nil
}

func (s *TransactionStorage) GetTransactionHashesByBatchIDs(batchIDs ...models.Uint256) ([]common.Hash, error) {
	batchPrefixes := batchIdsToBatchPrefixes(batchIDs)
	hashes := make([]common.Hash, 0, len(batchIDs)*32)

	var keyList bh.KeyList
	seekPrefix := db.IndexKeyPrefix(models.StoredTxReceiptName, "CommitmentID")
	err := s.database.Badger.Iterator(seekPrefix, db.ReversePrefetchIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			if validForPrefixes(keyValue(seekPrefix, item.Key()), batchPrefixes) {
				err := item.Value(func(val []byte) error {
					return db.DecodeKeyList(val, &keyList)
				})
				if err != nil {
					return false, err
				}
				txHashes, err := decodeKeyListHashes(models.StoredTxReceiptPrefix, keyList)
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

func (s *TransactionStorage) getStoredTxFromItem(item *bdg.Item, storedTx *models.StoredTx) (bool, error) {
	var hash common.Hash
	err := db.DecodeKey(item.Key(), &hash, models.StoredTxPrefix)
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

func getTxHashesByIndexKey(txn *bdg.Txn, indexKey, typePrefix []byte) ([]common.Hash, error) {
	item, err := txn.Get(indexKey)
	if err != nil {
		return nil, err
	}

	var keyList bh.KeyList
	err = item.Value(func(val []byte) error {
		return db.DecodeKeyList(val, &keyList)
	})
	if err != nil {
		return nil, err
	}

	return decodeKeyListHashes(typePrefix, keyList)
}

func decodeKeyListHashes(keyPrefix []byte, keyList bh.KeyList) ([]common.Hash, error) {
	var hash common.Hash
	hashes := make([]common.Hash, 0, len(keyList))
	for i := range keyList {
		err := models.DecodeDataHash(keyList[i][len(keyPrefix):], &hash)
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
		batchPrefixes = append(batchPrefixes, append([]byte{1}, batchIDs[i].Bytes()...))
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
