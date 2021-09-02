package storage

import (
	"sort"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

func (s *TransactionStorage) AddTransfer(t *models.Transfer) error {
	t.SetReceiveTime()
	//TODO-tx: wrap with txn if needed
	return s.addTransfer(t)
}

func (s *TransactionStorage) addTransfer(t *models.Transfer) error {
	if t.CommitmentID != nil || t.ErrorMessage != nil {
		err := s.database.Badger.Insert(t.Hash, models.MakeStoredTxReceiptFromTransfer(t))
		if err != nil {
			return err
		}
	}
	return s.database.Badger.Insert(t.Hash, models.MakeStoredTxFromTransfer(t))
}

// BatchAddTransfer contrary to the AddTransfer method does not set ReceiveTime field on added transfers
func (s *TransactionStorage) BatchAddTransfer(txs []models.Transfer) error {
	if len(txs) < 1 {
		return ErrNoRowsAffected
	}

	for i := range txs {
		err := s.addTransfer(&txs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TransactionStorage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	tx, txReceipt, err := s.getStoredTx(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Transfer {
		return nil, NewNotFoundError("transaction")
	}
	return tx.ToTransfer(txReceipt), nil
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

	if len(txs) <= int(limit) {
		return txs, nil
	}
	return txs[:limit], nil
}

func (s *TransactionStorage) GetTransfersByCommitmentID(id *models.CommitmentID) ([]models.Transfer, error) {
	txReceipts := make([]models.StoredTxReceipt, 0, 32)
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
		// TODO-tx: handle not found err
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
	for i := range txs {
		txReceipt := models.MakeStoredTxReceiptFromTransfer(&txs[i])
		txReceipt.CommitmentID = commitmentID
		err := s.addStoredTxReceipt(&txReceipt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetTransferWithBatchDetails(hash common.Hash) (*models.TransferWithBatchDetails, error) {
	tx, txReceipt, err := s.getStoredTx(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.Transfer {
		return nil, NewNotFoundError("transaction")
	}

	transfers, err := s.transfersToTransfersWithBatchDetails([]models.StoredTx{*tx}, []*models.StoredTxReceipt{txReceipt})
	if err != nil {
		return nil, err
	}
	return &transfers[0], nil
}

func (s *Storage) GetTransfersByPublicKey(publicKey *models.PublicKey) ([]models.TransferWithBatchDetails, error) {
	txs, txReceipts, err := s.getTransfersByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return s.transfersToTransfersWithBatchDetails(txs, txReceipts)
}

func (s *Storage) getTransfersByPublicKey(publicKey *models.PublicKey) (
	[]models.StoredTx, []*models.StoredTxReceipt, error,
) {
	leaves, err := s.GetStateLeavesByPublicKey(publicKey)
	if err != nil && !IsNotFoundError(err) {
		return nil, nil, err
	}
	stateIDs := utils.ValueToInterfaceSlice(leaves, "StateID")

	txs := make([]models.StoredTx, 0, 1)
	err = s.database.Badger.Find(
		&txs,
		bh.Where("ToStateID").In(stateIDs...).Index("ToStateID").
			And("TxType").Eq(txtype.Transfer).
			Or(bh.Where("FromStateID").In(stateIDs...).Index("FromStateID").
				And("TxType").Eq(txtype.Transfer),
			),
	)
	if err != nil {
		return nil, nil, err
	}

	txReceipts := make([]*models.StoredTxReceipt, 0, len(txs))
	for i := range txs {
		txReceipt, err := s.getStoredTxReceipt(txs[i].Hash)
		if err != nil {
			return nil, nil, err
		}
		txReceipts = append(txReceipts, txReceipt)
	}

	return txs, txReceipts, nil
}

func (s *Storage) transfersToTransfersWithBatchDetails(txs []models.StoredTx, txReceipts []*models.StoredTxReceipt) (
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
