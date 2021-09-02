package storage

import (
	"sort"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *TransactionStorage) AddCreate2Transfer(t *models.Create2Transfer) error {
	t.SetReceiveTime()
	return s.database.Badger.Insert(t.Hash, models.MakeStoredTransactionFromCreate2Transfer(t))
}

// BatchAddCreate2Transfer contrary to the AddCreate2Transfer method does not set ReceiveTime field on added transfers
func (s *TransactionStorage) BatchAddCreate2Transfer(txs []models.Create2Transfer) error {
	if len(txs) < 1 {
		return ErrNoRowsAffected
	}

	for i := range txs {
		err := s.database.Badger.Insert(txs[i].Hash, models.MakeStoredTransactionFromCreate2Transfer(&txs[i]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TransactionStorage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	tx, err := s.getStoredTransaction(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Create2Transfer {
		return nil, NewNotFoundError("transaction")
	}
	return tx.ToCreate2Transfer(), nil
}

func (s *TransactionStorage) GetPendingCreate2Transfers(limit uint32) ([]models.Create2Transfer, error) {
	txHashes, err := s.getPendingTransactionHashes()
	if IsNotFoundError(err) {
		return []models.Create2Transfer{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tx models.StoredTransaction
	txs := make([]models.Create2Transfer, 0, len(txHashes))
	for i := range txHashes {
		err = s.database.Badger.Get(txHashes[i], &tx)
		if err == bh.ErrNotFound {
			return nil, NewNotFoundError("transaction")
		}
		if err != nil {
			return nil, err
		}
		if tx.TxType == txtype.Create2Transfer && tx.ErrorMessage == nil {
			txs = append(txs, *tx.ToCreate2Transfer())
		}
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
	res := make([]models.StoredTransaction, 0, 32)
	err := s.database.Badger.Find(
		&res,
		bh.Where("CommitmentID").Eq(*id).Index("CommitmentID").
			And("TxType").Eq(txtype.Create2Transfer),
	)
	if err != nil {
		return nil, err
	}

	txs := make([]models.Create2Transfer, 0, len(res))
	for i := range res {
		txs = append(txs, *res[i].ToCreate2Transfer())
	}
	return txs, nil
}

func (s *Storage) GetCreate2TransfersByPublicKey(publicKey *models.PublicKey) ([]models.Create2TransferWithBatchDetails, error) {
	txs, txReceipts, err := s.getTransactionsByPublicKey(publicKey, txtype.Create2Transfer)
	if err != nil {
		return nil, err
	}
	return s.create2TransferToTransfersWithBatchDetails(txs, txReceipts)
}

func (s *TransactionStorage) MarkCreate2TransfersAsIncluded(txs []models.Create2Transfer, commitmentID *models.CommitmentID) error {
	for i := range txs {
		storedTx := models.MakeStoredTransactionFromCreate2Transfer(&txs[i])
		storedTx.CommitmentID = commitmentID
		err := s.updateStoredTransaction(&storedTx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetCreate2TransferWithBatchDetails(hash common.Hash) (*models.Create2TransferWithBatchDetails, error) {
	tx, txReceipt, err := s.getStoredTx(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Create2Transfer {
		return nil, NewNotFoundError("transaction")
	}

	transfers, err := s.create2TransferToTransfersWithBatchDetails([]*models.StoredTx{tx}, []*models.StoredTxReceipt{txReceipt})
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
