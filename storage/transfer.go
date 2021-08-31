package storage

import (
	"sort"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *TransactionStorage) AddTransfer(t *models.Transfer) error {
	t.SetReceiveTime()
	return s.database.Badger.Insert(t.Hash, models.MakeStoredTransactionFromTransfer(t))
}

// BatchAddTransfer contrary to the AddTransfer method does not set ReceiveTime field on added transfers
func (s *TransactionStorage) BatchAddTransfer(txs []models.Transfer) error {
	if len(txs) < 1 {
		return ErrNoRowsAffected
	}

	tx, txStorage, err := s.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range txs {
		err = txStorage.database.Badger.Insert(txs[i].Hash, models.MakeStoredTransactionFromTransfer(&txs[i]))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *TransactionStorage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	tx, err := s.getStoredTransaction(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Transfer {
		return nil, NewNotFoundError("transaction")
	}
	return tx.ToTransfer(), nil
}

func (s *TransactionStorage) GetPendingTransfers(limit uint32) ([]models.Transfer, error) {
	txHashes, err := s.getPendingTransactionHashes()
	if IsNotFoundError(err) {
		return []models.Transfer{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tx models.StoredTransaction
	txs := make([]models.Transfer, 0, len(txHashes))
	for i := range txHashes {
		err = s.database.Badger.Get(txHashes[i], &tx)
		if err == bh.ErrNotFound {
			return nil, NewNotFoundError("transaction")
		}
		if err != nil {
			return nil, err
		}
		if tx.TxType == txtype.Transfer && tx.ErrorMessage == nil {
			txs = append(txs, *tx.ToTransfer())
		}
	}

	sort.SliceStable(txs, func(i, j int) bool {
		return txs[i].Nonce.Cmp(&txs[j].Nonce) < 0
	})

	return txs, nil
}

func (s *TransactionStorage) GetTransfersByCommitmentID(id *models.CommitmentID) ([]models.Transfer, error) {
	res := make([]models.StoredTransaction, 0, 32)
	err := s.database.Badger.Find(
		&res,
		bh.Where("CommitmentID").Eq(*id).Index("CommitmentID").
			And("TxType").Eq(txtype.Transfer),
	)
	if err != nil {
		return nil, err
	}

	txs := make([]models.Transfer, 0, len(res))
	for i := range res {
		txs = append(txs, *res[i].ToTransfer())
	}
	return txs, nil
}

func (s *TransactionStorage) MarkTransfersAsIncluded(txs []models.Transfer, commitmentID *models.CommitmentID) error {
	tx, txStorage, err := s.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range txs {
		storedTx := models.MakeStoredTransactionFromTransfer(&txs[i])
		storedTx.CommitmentID = commitmentID
		err = txStorage.updateStoredTransaction(&storedTx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) GetTransferWithBatchDetails(hash common.Hash) (*models.TransferWithBatchDetails, error) {
	tx, err := s.getStoredTransaction(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Transfer {
		return nil, NewNotFoundError("transaction")
	}

	transfers, err := s.transfersToTransfersWithBatchDetails(*tx)
	if err != nil {
		return nil, err
	}
	return &transfers[0], nil
}

func (s *Storage) GetTransfersByPublicKey(publicKey *models.PublicKey) ([]models.TransferWithBatchDetails, error) {
	txs, err := s.getTransactionsByPublicKey(publicKey, txtype.Transfer)
	if err != nil {
		return nil, err
	}
	return s.transfersToTransfersWithBatchDetails(txs...)
}

func (s *Storage) transfersToTransfersWithBatchDetails(txs ...models.StoredTransaction) (
	result []models.TransferWithBatchDetails,
	err error,
) {
	result = make([]models.TransferWithBatchDetails, 0, len(txs))
	batchIDs := make(map[models.Uint256]*models.Batch)
	for i := range txs {
		transfer := txs[i].ToTransfer()
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
