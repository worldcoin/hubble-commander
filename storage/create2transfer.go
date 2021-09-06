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

func (s *TransactionStorage) AddCreate2Transfer(t *models.Create2Transfer) error {
	t.SetReceiveTime()
	return s.addCreate2Transfer(t)
}

func (s *TransactionStorage) addCreate2Transfer(t *models.Create2Transfer) error {
	if t.CommitmentID != nil || t.ErrorMessage != nil || t.ToStateID != nil {
		err := s.database.Badger.Insert(t.Hash, models.MakeStoredReceiptFromCreate2Transfer(t))
		if err != nil {
			return err
		}
	}
	return s.database.Badger.Insert(t.Hash, models.MakeStoredTxFromCreate2Transfer(t))
}

// BatchAddCreate2Transfer contrary to the AddCreate2Transfer method does not set ReceiveTime field on added transfers
func (s *TransactionStorage) BatchAddCreate2Transfer(txs []models.Create2Transfer) error {
	if len(txs) < 1 {
		return ErrNoRowsAffected
	}

	for i := range txs {
		err := s.addCreate2Transfer(&txs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TransactionStorage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Create2Transfer {
		return nil, NewNotFoundError("transaction")
	}
	return tx.ToCreate2Transfer(txReceipt), nil
}

func (s *TransactionStorage) GetPendingCreate2Transfers(limit uint32) ([]models.Create2Transfer, error) {
	txController, txStorage, err := s.BeginTransaction(TxOptions{Badger: true, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer txController.Rollback(&err)

	txs, err := txStorage.unsafeGetPendingCreate2Transfers(limit)
	if err != nil {
		return nil, err
	}

	err = txController.Commit()
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionStorage) unsafeGetPendingCreate2Transfers(limit uint32) ([]models.Create2Transfer, error) {
	txs := make([]models.Create2Transfer, 0, 32)
	var storedTx models.StoredTx
	err := s.database.Badger.Iterator(models.StoredTxPrefix, badger.KeyIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			var hash common.Hash
			err := badger.DecodeKey(item.Key(), &hash, models.StoredTxPrefix)
			if err != nil {
				return false, err
			}
			txReceipt, err := s.getStoredTxReceipt(hash)
			if err != nil || txReceipt != nil {
				return false, err
			}

			err = item.Value(storedTx.SetBytes)
			if err != nil {
				return false, err
			}
			if storedTx.TxType == txtype.Create2Transfer {
				txs = append(txs, *storedTx.ToCreate2Transfer(nil))
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

func (s *TransactionStorage) GetCreate2TransfersByCommitmentID(id *models.CommitmentID) ([]models.Create2Transfer, error) {
	txReceipts := make([]models.StoredReceipt, 0, 32)
	err := s.database.Badger.Find(
		&txReceipts,
		bh.Where("CommitmentID").Eq(*id).Index("CommitmentID"),
	)
	if err != nil {
		return nil, err
	}

	transfers := make([]models.Create2Transfer, 0, len(txReceipts))
	var tx models.StoredTx
	for i := range txReceipts {
		err = s.database.Badger.Get(txReceipts[i].Hash, &tx)
		if err != nil {
			return nil, err
		}
		if tx.TxType == txtype.Create2Transfer {
			transfers = append(transfers, *tx.ToCreate2Transfer(&txReceipts[i]))
		}
	}
	return transfers, nil
}

func (s *Storage) GetCreate2TransfersByPublicKey(publicKey *models.PublicKey) ([]models.Create2TransferWithBatchDetails, error) {
	txs, txReceipts, err := s.getCreate2TransfersByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return s.create2TransferToTransfersWithBatchDetails(txs, txReceipts)
}

func (s *Storage) getCreate2TransfersByPublicKey(publicKey *models.PublicKey) (
	[]*models.StoredTx, []*models.StoredReceipt, error,
) {
	leaves, err := s.GetStateLeavesByPublicKey(publicKey)
	if err != nil && !IsNotFoundError(err) {
		return nil, nil, err
	}
	stateIDs := utils.ValueToInterfaceSlice(leaves, "StateID")

	txs := make([]models.StoredTx, 0, 1)
	err = s.database.Badger.Find(
		&txs,
		bh.Where("FromStateID").In(stateIDs...).Index("FromStateID").
			And("TxType").Eq(txtype.Create2Transfer),
	)
	if err != nil {
		return nil, nil, err
	}

	receipts := make([]models.StoredReceipt, 0, 1)
	err = s.database.Badger.Find(
		&receipts,
		bh.Where("ToStateID").In(stateIDs...).Index("ToStateID"),
	)
	if err != nil {
		return nil, nil, err
	}

	return s.getMissingStoredTxsData(txs, receipts)
}

func (s *Storage) getMissingStoredTxsData(txs []models.StoredTx, receipts []models.StoredReceipt) (
	[]*models.StoredTx, []*models.StoredReceipt, error,
) {
	hashes := make(map[common.Hash]struct{}, len(txs))

	for i := range txs {
		hashes[txs[i].Hash] = struct{}{}
	}

	for i := range receipts {
		hashes[receipts[i].Hash] = struct{}{}
	}

	resultTxs := make([]*models.StoredTx, 0, len(hashes))
	resultReceipts := make([]*models.StoredReceipt, 0, len(hashes))
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
	for i := range txs {
		txReceipt := models.MakeStoredReceiptFromCreate2Transfer(&txs[i])
		txReceipt.CommitmentID = commitmentID
		err := s.addStoredReceipt(&txReceipt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetCreate2TransferWithBatchDetails(hash common.Hash) (*models.Create2TransferWithBatchDetails, error) {
	tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Create2Transfer {
		return nil, NewNotFoundError("transaction")
	}

	transfers, err := s.create2TransferToTransfersWithBatchDetails([]*models.StoredTx{tx}, []*models.StoredReceipt{txReceipt})
	if err != nil {
		return nil, err
	}
	return &transfers[0], nil
}

func (s *Storage) create2TransferToTransfersWithBatchDetails(txs []*models.StoredTx, txReceipts []*models.StoredReceipt) (
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
