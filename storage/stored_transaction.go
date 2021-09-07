package storage

import (
	"bytes"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
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

	txTransactionStorage := *s
	txTransactionStorage.database = txDatabase

	return txController, &txTransactionStorage, nil
}

func (s *TransactionStorage) wrapWithTransaction(opts TxOptions, fn func(txStorage *TransactionStorage) error) error {
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

func (s *TransactionStorage) getStoredTxWithReceipt(hash common.Hash) (*models.StoredTx, *models.StoredReceipt, error) {
	txController, txStorage, err := s.BeginTransaction(TxOptions{Badger: true, ReadOnly: true})
	if err != nil {
		return nil, nil, err
	}
	defer txController.Rollback(&err)

	storedTx, err := txStorage.getStoredTx(hash)
	if err != nil {
		return nil, nil, err
	}

	txReceipt, err := txStorage.getStoredTxReceipt(hash)
	if err != nil {
		return nil, nil, err
	}

	err = txController.Commit()
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
			return badger.DecodeKeyList(val, &keyList)
		})
	})
	if err != nil {
		return nil, err
	}
	return &keyList, nil
}

func (s *TransactionStorage) GetLatestTransactionNonce(accountStateID uint32) (*models.Uint256, error) {
	encodedStateID, err := models.EncodeUint32(&accountStateID)
	if err != nil {
		return nil, err
	}

	txController, txStorage, err := s.BeginTransaction(TxOptions{Badger: true, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer txController.Rollback(&err)

	indexKey := badger.IndexKey(models.StoredTxName, "FromStateID", encodedStateID)
	keyList, err := txStorage.getKeyList(indexKey)
	if err != nil {
		return nil, err
	}
	txHashes, err := decodeKeyListHashes(models.StoredTxPrefix, *keyList)
	if err != nil {
		return nil, err
	}
	if len(txHashes) == 0 {
		return nil, NewNotFoundError("transaction")
	}

	latestNonce := models.MakeUint256(0)
	for i := range txHashes {
		var tx *models.StoredTx
		tx, err = txStorage.getStoredTx(txHashes[i])
		if err != nil {
			return nil, err
		}
		if tx.Nonce.Cmp(&latestNonce) > 0 {
			latestNonce = tx.Nonce
		}
	}

	err = txController.Commit()
	if err != nil {
		return nil, err
	}
	return &latestNonce, nil
}

func (s *TransactionStorage) MarkTransactionsAsPending(txHashes []common.Hash) error {
	return s.wrapWithTransaction(TxOptions{Badger: true}, func(txStorage *TransactionStorage) error {
		dataType := models.StoredReceipt{}
		for i := range txHashes {
			err := s.database.Badger.Delete(txHashes[i], dataType)
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
	txController, txStorage, err := s.BeginTransaction(TxOptions{Badger: true, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer txController.Rollback(&err)

	latestBatch, err := txStorage.GetLatestSubmittedBatch()
	if IsNotFoundError(err) {
		return ref.Int(0), nil
	}
	if err != nil {
		return nil, err
	}
	count := 0
	seekPrefix := badger.IndexKeyPrefix(models.StoredReceiptName, "CommitmentID")
	err = txStorage.database.Badger.Iterator(seekPrefix, badger.PrefetchIteratorOpts,
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
				return badger.DecodeKeyList(val, &keyList)
			})
			if err != nil {
				return false, err
			}
			count += len(keyList)
			return false, nil
		})
	if err != nil && err != badger.ErrIteratorFinished {
		return nil, err
	}

	err = txController.Commit()
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (s *TransactionStorage) GetTransactionHashesByBatchIDs(batchIDs ...models.Uint256) ([]common.Hash, error) {
	batchPrefixes := batchIdsToBatchPrefixes(batchIDs)
	hashes := make([]common.Hash, 0, len(batchIDs)*32)

	var keyList bh.KeyList
	seekPrefix := badger.IndexKeyPrefix(models.StoredReceiptName, "CommitmentID")
	err := s.database.Badger.Iterator(seekPrefix, badger.ReversePrefetchIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			if validForPrefixes(keyValue(seekPrefix, item.Key()), batchPrefixes) {
				err := item.Value(func(val []byte) error {
					return badger.DecodeKeyList(val, &keyList)
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
	if err != nil && err != badger.ErrIteratorFinished {
		return nil, err
	}
	if len(hashes) == 0 {
		return nil, NewNotFoundError("transaction")
	}
	return hashes, nil
}

func (s *TransactionStorage) getStoredTxFromItem(item *bdg.Item, storedTx *models.StoredTx) (bool, error) {
	var hash common.Hash
	err := badger.DecodeKey(item.Key(), &hash, models.StoredTxPrefix)
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
		return badger.DecodeKeyList(val, &keyList)
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
