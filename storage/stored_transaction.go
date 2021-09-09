package storage

import (
	"bytes"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

type TransactionStorage struct {
	database *Database
}

func NewTransactionStorage(database *Database) *TransactionStorage {
	return &TransactionStorage{
		database: database,
	}
}

func (s *TransactionStorage) copyWithNewDatabase(database *Database) *TransactionStorage {
	newTransactionStorage := *s
	newTransactionStorage.database = database

	return &newTransactionStorage
}

func (s *TransactionStorage) BeginTransaction(opts TxOptions) (*db.TxController, *TransactionStorage, error) {
	txController, txDatabase, err := s.database.BeginTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txTransactionStorage := NewTransactionStorage(txDatabase)
	return txController, txTransactionStorage, nil
}

func (s *TransactionStorage) executeInTransaction(opts TxOptions, fn func(txStorage *TransactionStorage) error) error {
	err := s.unsafeExecuteInTransaction(opts, fn)
	if err == bdg.ErrConflict {
		return s.executeInTransaction(opts, fn)
	}
	return err
}

func (s *TransactionStorage) unsafeExecuteInTransaction(opts TxOptions, fn func(txStorage *TransactionStorage) error) error {
	txController, txStorage, err := s.BeginTransaction(opts)
	if err != nil {
		return err
	}
	defer txController.Rollback(&err)

	err = fn(txStorage)
	if err != nil {
		return err
	}

	return txController.Commit()
}

func (s *TransactionStorage) addStoredReceipt(txReceipt *models.StoredReceipt) error {
	return s.database.Badger.Insert(txReceipt.Hash, *txReceipt)
}

func (s *TransactionStorage) getStoredTxWithReceipt(hash common.Hash) (
	storedTx *models.StoredTx, txReceipt *models.StoredReceipt, err error,
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

func (s *TransactionStorage) getStoredTx(hash common.Hash) (*models.StoredTx, error) {
	var storedTx models.StoredTx
	err := s.database.Badger.Get(hash, &storedTx)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("transaction")
	}
	if err != nil {
		return nil, err
	}
	return &storedTx, nil
}

func (s *TransactionStorage) getStoredTxReceipt(hash common.Hash) (*models.StoredReceipt, error) {
	var storedTxReceipt models.StoredReceipt
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
			return NewNotFoundError("transaction")
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
	latestNonce := models.MakeUint256(0)

	err := s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		encodedStateID, err := models.EncodeUint32(&accountStateID)
		if err != nil {
			return err
		}

		indexKey := db.IndexKey(models.StoredTxName, "FromStateID", encodedStateID)
		keyList, err := txStorage.getKeyList(indexKey)
		if err != nil {
			return err
		}
		txHashes, err := decodeKeyListHashes(models.StoredTxPrefix, *keyList)
		if err != nil {
			return err
		}
		if len(txHashes) == 0 {
			return NewNotFoundError("transaction")
		}

		for i := range txHashes {
			tx, err := txStorage.getStoredTx(txHashes[i])
			if err != nil {
				return err
			}
			if tx.Nonce.Cmp(&latestNonce) > 0 {
				latestNonce = tx.Nonce
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &latestNonce, nil
}

func (s *TransactionStorage) MarkTransactionsAsPending(txHashes []common.Hash) error {
	dataType := models.StoredReceipt{}
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

func (s *TransactionStorage) SetTransactionError(txHash common.Hash, errorMessage string) error {
	return s.addStoredReceipt(&models.StoredReceipt{
		Hash:         txHash,
		ErrorMessage: &errorMessage,
	})
}

func (s *Storage) GetTransactionCount() (*int, error) {
	count := new(int)
	err := s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		latestBatch, err := txStorage.GetLatestSubmittedBatch()
		if IsNotFoundError(err) {
			return nil
		}
		if err != nil {
			return err
		}
		seekPrefix := db.IndexKeyPrefix(models.StoredReceiptName, "CommitmentID")
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
	seekPrefix := db.IndexKeyPrefix(models.StoredReceiptName, "CommitmentID")
	err := s.database.Badger.Iterator(seekPrefix, db.ReversePrefetchIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			if validForPrefixes(keyValue(seekPrefix, item.Key()), batchPrefixes) {
				err := item.Value(func(val []byte) error {
					return db.DecodeKeyList(val, &keyList)
				})
				if err != nil {
					return false, err
				}
				txHashes, err := decodeKeyListHashes(models.StoredReceiptPrefix, keyList)
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
		return nil, NewNotFoundError("transaction")
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
