package storage

import (
	"sort"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *TransactionStorage) AddCreate2Transfer(t *models.Create2Transfer) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		if t.CommitmentID != nil || t.ErrorMessage != nil || t.ToStateID != nil {
			err := txStorage.addStoredTxReceipt(models.NewStoredTxReceiptFromCreate2Transfer(t))
			if err != nil {
				return err
			}
		}
		return txStorage.database.Badger.Insert(t.Hash, *models.NewStoredTxFromCreate2Transfer(t))
	})
}

func (s *TransactionStorage) BatchAddCreate2Transfer(txs []models.Create2Transfer) error {
	if len(txs) < 1 {
		return errors.WithStack(ErrNoRowsAffected)
	}

	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			err := txStorage.AddCreate2Transfer(&txs[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Create2Transfer {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return tx.ToCreate2Transfer(txReceipt), nil
}

func (s *TransactionStorage) GetPendingCreate2Transfers() (txs []models.Create2Transfer, err error) {
	err = s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		txs, err = txStorage.unsafeGetPendingCreate2Transfers()
		return err
	})
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionStorage) unsafeGetPendingCreate2Transfers() ([]models.Create2Transfer, error) {
	txs := make([]models.Create2Transfer, 0, 32)
	var storedTx models.StoredTx
	err := s.database.Badger.Iterator(models.StoredTxPrefix, db.KeyIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			skip, err := s.getStoredTxFromItem(item, &storedTx)
			if err != nil || skip {
				return false, err
			}
			if storedTx.TxType == txtype.Create2Transfer {
				txs = append(txs, *storedTx.ToCreate2Transfer(nil))
			}
			return false, nil
		})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}

	sort.SliceStable(txs, func(i, j int) bool {
		return txs[i].Nonce.Cmp(&txs[j].Nonce) < 0
	})

	return txs, nil
}

func (s *TransactionStorage) GetCreate2TransfersByCommitmentID(id *models.CommitmentID) ([]models.Create2Transfer, error) {
	encodeCommitmentID := models.EncodeCommitmentIDPointer(id)
	indexKey := db.IndexKey(models.StoredTxReceiptName, "CommitmentID", encodeCommitmentID)

	var transfers []models.Create2Transfer
	err := s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		// queried Badger directly due to nil index decoding problem
		return txStorage.database.Badger.View(func(txn *bdg.Txn) error {
			hashes, err := getTxHashesByIndexKey(txn, indexKey, models.StoredTxReceiptPrefix)
			if err == bdg.ErrKeyNotFound {
				return nil
			}
			if err != nil {
				return err
			}

			transfers = make([]models.Create2Transfer, 0, len(hashes))
			for i := range hashes {
				storedTx, storedTxReceipt, err := txStorage.getStoredTxWithReceipt(hashes[i])
				if err != nil {
					return err
				}
				if storedTx.TxType == txtype.Create2Transfer {
					transfers = append(transfers, *storedTx.ToCreate2Transfer(storedTxReceipt))
				}
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (s *Storage) GetCreate2TransfersByPublicKey(publicKey *models.PublicKey) ([]models.Create2TransferWithBatchDetails, error) {
	var transfers []models.Create2TransferWithBatchDetails
	err := s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		txs, txReceipts, err := txStorage.getCreate2TransfersByPublicKey(publicKey)
		if err != nil {
			return err
		}
		transfers, err = txStorage.create2TransferToTransfersWithBatchDetails(txs, txReceipts)
		return err
	})
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (s *Storage) getCreate2TransfersByPublicKey(publicKey *models.PublicKey) (
	[]*models.StoredTx, []*models.StoredTxReceipt, error,
) {
	leaves, err := s.GetStateLeavesByPublicKey(publicKey)
	if err != nil && !IsNotFoundError(err) {
		return nil, nil, err
	}

	fromStateIDs := make([]interface{}, 0, len(leaves))
	toStateIDs := make([]uint32, 0, len(leaves))
	for i := range leaves {
		fromStateIDs = append(fromStateIDs, leaves[i].StateID)
		toStateIDs = append(toStateIDs, leaves[i].StateID)
	}

	txs := make([]models.StoredTx, 0, 1)
	err = s.database.Badger.Find(
		&txs,
		bh.Where("FromStateID").In(fromStateIDs...).Index("FromStateID").
			And("TxType").Eq(txtype.Create2Transfer),
	)
	if err != nil {
		return nil, nil, err
	}

	txHashes, err := s.getC2THashesByStateIDs(toStateIDs)
	if err != nil {
		return nil, nil, err
	}

	return s.getMissingStoredTxsData(txs, txHashes)
}

func (s *TransactionStorage) getC2THashesByStateIDs(stateIDs []uint32) ([]common.Hash, error) {
	results := make([]common.Hash, 0, len(stateIDs))
	err := s.database.Badger.View(func(txn *bdg.Txn) error {
		for i := range stateIDs {
			encodedStateID := models.EncodeUint32Pointer(&stateIDs[i])
			indexKey := db.IndexKey(models.StoredTxReceiptName, "ToStateID", encodedStateID)
			hashes, err := getTxHashesByIndexKey(txn, indexKey, models.StoredTxReceiptPrefix)
			if err == bdg.ErrKeyNotFound {
				continue
			}
			if err != nil {
				return err
			}
			results = append(results, hashes...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Storage) getMissingStoredTxsData(txs []models.StoredTx, receiptHashes []common.Hash) (
	[]*models.StoredTx, []*models.StoredTxReceipt, error,
) {
	hashes := make(map[common.Hash]struct{}, len(txs))

	for i := range receiptHashes {
		hashes[receiptHashes[i]] = struct{}{}
	}

	for i := range txs {
		hashes[txs[i].Hash] = struct{}{}
	}

	resultTxs := make([]*models.StoredTx, 0, len(hashes))
	resultReceipts := make([]*models.StoredTxReceipt, 0, len(hashes))
	for hash := range hashes {
		tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
		if err != nil {
			return nil, nil, err
		}
		resultTxs = append(resultTxs, tx)
		resultReceipts = append(resultReceipts, txReceipt)
	}

	return resultTxs, resultReceipts, nil
}

func (s *TransactionStorage) MarkCreate2TransfersAsIncluded(txs []models.Create2Transfer, commitmentID *models.CommitmentID) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			txReceipt := models.NewStoredTxReceiptFromCreate2Transfer(&txs[i])
			txReceipt.CommitmentID = commitmentID
			err := txStorage.addStoredTxReceipt(txReceipt)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) GetCreate2TransferWithBatchDetails(hash common.Hash) (*models.Create2TransferWithBatchDetails, error) {
	var transfers []models.Create2TransferWithBatchDetails
	err := s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		tx, txReceipt, err := txStorage.getStoredTxWithReceipt(hash)
		if err != nil {
			return err
		}
		if tx.TxType != txtype.Create2Transfer {
			return errors.WithStack(NewNotFoundError("transaction"))
		}

		transfers, err = txStorage.create2TransferToTransfersWithBatchDetails([]*models.StoredTx{tx}, []*models.StoredTxReceipt{txReceipt})
		return err
	})
	if err != nil {
		return nil, err
	}
	return &transfers[0], nil
}

func (s *Storage) create2TransferToTransfersWithBatchDetails(txs []*models.StoredTx, txReceipts []*models.StoredTxReceipt) (
	result []models.Create2TransferWithBatchDetails,
	err error,
) {
	result = make([]models.Create2TransferWithBatchDetails, 0, len(txs))
	batchIDs := make(map[models.Uint256]*models.Batch)
	for i := range txs {
		transfer := txs[i].ToCreate2Transfer(txReceipts[i])
		if transfer.CommitmentID == nil {
			result = append(result, models.Create2TransferWithBatchDetails{Create2Transfer: *transfer})
			continue
		}
		batch, ok := batchIDs[transfer.CommitmentID.BatchID]
		if !ok {
			batch, err = s.GetBatch(transfer.CommitmentID.BatchID)
			if err != nil {
				return nil, err
			}
			batchIDs[transfer.CommitmentID.BatchID] = batch
		}

		result = append(result, models.Create2TransferWithBatchDetails{
			Create2Transfer: *transfer,
			BatchHash:       batch.Hash,
			BatchTime:       batch.SubmissionTime,
		})
	}
	return result, nil
}
