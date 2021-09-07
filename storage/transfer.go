package storage

import (
	"sort"

	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *TransactionStorage) AddTransfer(t *models.Transfer) error {
	return s.wrapWithTransaction(TxOptions{Badger: true}, func(txStorage *TransactionStorage) error {
		if t.CommitmentID != nil || t.ErrorMessage != nil {
			err := txStorage.database.Badger.Insert(t.Hash, models.MakeStoredReceiptFromTransfer(t))
			if err != nil {
				return err
			}
		}
		return txStorage.database.Badger.Insert(t.Hash, models.MakeStoredTxFromTransfer(t))
	})
}

func (s *TransactionStorage) BatchAddTransfer(txs []models.Transfer) error {
	if len(txs) < 1 {
		return ErrNoRowsAffected
	}

	return s.wrapWithTransaction(TxOptions{Badger: true}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			err := s.AddTransfer(&txs[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TransactionStorage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Transfer {
		return nil, NewNotFoundError("transaction")
	}
	return tx.ToTransfer(txReceipt), nil
}

func (s *TransactionStorage) GetPendingTransfers(limit uint32) ([]models.Transfer, error) {
	txController, txStorage, err := s.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return nil, err
	}
	defer txController.Rollback(&err)

	txs, err := txStorage.unsafeGetPendingTransfers(limit)
	if err != nil {
		return nil, err
	}

	err = txController.Commit()
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionStorage) unsafeGetPendingTransfers(limit uint32) ([]models.Transfer, error) {
	txs := make([]models.Transfer, 0, 32)
	var storedTx models.StoredTx
	err := s.database.Badger.Iterator(models.StoredTxPrefix, badger.KeyIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			skip, err := s.getStoredTxFromItem(item, &storedTx)
			if err != nil || skip {
				return false, err
			}
			if storedTx.TxType == txtype.Transfer {
				txs = append(txs, *storedTx.ToTransfer(nil))
			}
			return false, nil
		})
	if err != nil && err != badger.ErrIteratorFinished {
		return nil, err
	}

	sort.SliceStable(txs, func(i, j int) bool {
		return txs[i].Nonce.Cmp(&txs[j].Nonce) < 0
	})

	if len(txs) <= int(limit) {
		return txs, nil
	}
	return txs[:limit], nil
}

func (s *TransactionStorage) GetTransfersByCommitmentID(id *models.CommitmentID) ([]models.Transfer, error) {
	txReceipts := make([]models.StoredReceipt, 0, 32)
	// TODO rework to query the index table manually and get only the valid StoredReceipts, leave a comment explaining
	//  why Find is not used here
	err := s.database.Badger.Find(
		&txReceipts,
		bh.Where("CommitmentID").Eq(*id).Index("CommitmentID"),
	)
	if err != nil {
		return nil, err
	}

	transfers := make([]models.Transfer, 0, len(txReceipts))
	var tx models.StoredTx
	for i := range txReceipts {
		err = s.database.Badger.Get(txReceipts[i].Hash, &tx)
		if err != nil {
			return nil, err
		}
		if tx.TxType == txtype.Transfer {
			transfers = append(transfers, *tx.ToTransfer(&txReceipts[i]))
		}
	}
	return transfers, nil
}

func (s *TransactionStorage) MarkTransfersAsIncluded(txs []models.Transfer, commitmentID *models.CommitmentID) error {
	return s.wrapWithTransaction(TxOptions{Badger: true}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			txReceipt := models.MakeStoredReceiptFromTransfer(&txs[i])
			txReceipt.CommitmentID = commitmentID
			err := s.addStoredReceipt(&txReceipt)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) GetTransferWithBatchDetails(hash common.Hash) (*models.TransferWithBatchDetails, error) {
	txController, txStorage, err := s.BeginTransaction(TxOptions{Badger: true, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer txController.Rollback(&err)

	tx, txReceipt, err := txStorage.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Transfer {
		return nil, NewNotFoundError("transaction")
	}

	transfers, err := txStorage.transfersToTransfersWithBatchDetails([]models.StoredTx{*tx}, []*models.StoredReceipt{txReceipt})
	if err != nil {
		return nil, err
	}

	err = txController.Commit()
	if err != nil {
		return nil, err
	}
	return &transfers[0], nil
}

func (s *Storage) GetTransfersByPublicKey(publicKey *models.PublicKey) ([]models.TransferWithBatchDetails, error) {
	txController, txStorage, err := s.BeginTransaction(TxOptions{Badger: true, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer txController.Rollback(&err)

	txs, txReceipts, err := txStorage.getTransfersByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	results, err := txStorage.transfersToTransfersWithBatchDetails(txs, txReceipts)
	if err != nil {
		return nil, err
	}

	err = txController.Commit()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Storage) getTransfersByPublicKey(publicKey *models.PublicKey) (
	[]models.StoredTx, []*models.StoredReceipt, error,
) {
	leaves, err := s.GetStateLeavesByPublicKey(publicKey)
	if err != nil && !IsNotFoundError(err) {
		return nil, nil, err
	}
	stateIDs := utils.ValueToInterfaceSlice(leaves, "StateID")

	txs := make([]models.StoredTx, 0, 1)
	toStateIDCondition := bh.Where("ToStateID").In(stateIDs...).Index("ToStateID").
		And("TxType").Eq(txtype.Transfer)
	fromStateIDCondition := bh.Where("FromStateID").In(stateIDs...).Index("FromStateID").
		And("TxType").Eq(txtype.Transfer)
	err = s.database.Badger.Find(
		&txs,
		toStateIDCondition.Or(fromStateIDCondition),
	)
	if err != nil {
		return nil, nil, err
	}

	txReceipts := make([]*models.StoredReceipt, 0, len(txs))
	for i := range txs {
		txReceipt, err := s.getStoredTxReceipt(txs[i].Hash)
		if err != nil {
			return nil, nil, err
		}
		txReceipts = append(txReceipts, txReceipt)
	}

	return txs, txReceipts, nil
}

func (s *Storage) transfersToTransfersWithBatchDetails(txs []models.StoredTx, txReceipts []*models.StoredReceipt) (
	result []models.TransferWithBatchDetails,
	err error,
) {
	result = make([]models.TransferWithBatchDetails, 0, len(txs))
	batchIDs := make(map[models.Uint256]*models.Batch)
	for i := range txs {
		transfer := txs[i].ToTransfer(txReceipts[i])
		if transfer.CommitmentID == nil {
			result = append(result, models.TransferWithBatchDetails{Transfer: *transfer})
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

		result = append(result, models.TransferWithBatchDetails{
			Transfer:  *transfer,
			BatchHash: batch.Hash,
			BatchTime: batch.SubmissionTime,
		})
	}
	return result, nil
}
