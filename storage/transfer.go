package storage

import (
	"sort"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/utils"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *TransactionStorage) AddTransfer(t *models.Transfer) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		if t.CommitmentID != nil || t.ErrorMessage != nil {
			err := txStorage.addStoredTxReceipt(stored.NewTxReceiptFromTransfer(t))
			if err != nil {
				return err
			}
		}
		return txStorage.addStoredTx(stored.NewTxFromTransfer(t))
	})
}

func (s *TransactionStorage) BatchAddTransfer(txs []models.Transfer) error {
	if len(txs) < 1 {
		return errors.WithStack(ErrNoRowsAffected)
	}

	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			err := txStorage.AddTransfer(&txs[i])
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
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return tx.ToTransfer(txReceipt), nil
}

func (s *TransactionStorage) GetPendingTransfers() (txs models.TransferArray, err error) {
	err = s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		txs, err = txStorage.unsafeGetPendingTransfers()
		return err
	})
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionStorage) unsafeGetPendingTransfers() ([]models.Transfer, error) {
	txs := make([]models.Transfer, 0, 32)
	var storedTx stored.Tx
	err := s.database.Badger.Iterator(stored.TxPrefix, db.KeyIteratorOpts,
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
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}

	sort.SliceStable(txs, func(i, j int) bool {
		return txs[i].Nonce.Cmp(&txs[j].Nonce) < 0
	})

	return txs, nil
}

func (s *TransactionStorage) GetTransfersByCommitmentID(id models.CommitmentID) ([]models.Transfer, error) {
	transfers := make([]models.Transfer, 0, 1)

	err := s.iterateTxsByCommitmentID(id, func(storedTx *stored.Tx, storedTxReceipt *stored.TxReceipt) {
		if storedTx.TxType == txtype.Transfer {
			transfers = append(transfers, *storedTx.ToTransfer(storedTxReceipt))
		}
	})
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func (s *TransactionStorage) iterateTxsByCommitmentID(
	id models.CommitmentID,
	handleTx func(storedTx *stored.Tx, storedTxReceipt *stored.TxReceipt),
) error {
	return s.executeInTransaction(TxOptions{ReadOnly: true}, func(txStorage *TransactionStorage) error {
		return txStorage.unsafeIterateTxsByCommitmentID(id, handleTx)
	})
}

func (s *TransactionStorage) unsafeIterateTxsByCommitmentID(
	id models.CommitmentID,
	handleTx func(storedTx *stored.Tx, storedTxReceipt *stored.TxReceipt),
) error {
	receipts := make([]stored.TxReceipt, 0, 1)
	err := s.database.Badger.Find(
		&receipts,
		bh.Where("CommitmentID").Eq(id).Index("CommitmentID"),
	)
	if err != nil {
		return err
	}

	for i := range receipts {
		storedTx, storedTxReceipt, err := s.getStoredTxWithReceipt(receipts[i].Hash)
		if err != nil {
			return err
		}
		handleTx(storedTx, storedTxReceipt)
	}
	return nil
}

func (s *TransactionStorage) MarkTransfersAsIncluded(txs []models.Transfer, commitmentID *models.CommitmentID) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		for i := range txs {
			txReceipt := stored.NewTxReceiptFromTransfer(&txs[i])
			txReceipt.CommitmentID = commitmentID
			err := txStorage.addStoredTxReceipt(txReceipt)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) GetTransferWithBatchDetails(hash common.Hash) (*models.TransferWithBatchDetails, error) {
	var transfers []models.TransferWithBatchDetails
	err := s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		tx, txReceipt, err := txStorage.getStoredTxWithReceipt(hash)
		if err != nil {
			return err
		}
		if tx.TxType != txtype.Transfer {
			return errors.WithStack(NewNotFoundError("transaction"))
		}

		transfers, err = txStorage.txsToTransfersWithBatchDetails([]stored.Tx{*tx}, []*stored.TxReceipt{txReceipt})
		return err
	})

	if err != nil {
		return nil, err
	}
	return &transfers[0], nil
}

func (s *Storage) GetTransfersByPublicKey(publicKey *models.PublicKey) ([]models.TransferWithBatchDetails, error) {
	var transfers []models.TransferWithBatchDetails
	err := s.ExecuteInTransaction(TxOptions{ReadOnly: true}, func(txStorage *Storage) error {
		txs, txReceipts, err := txStorage.getTransfersByPublicKey(publicKey)
		if err != nil {
			return err
		}
		transfers, err = txStorage.txsToTransfersWithBatchDetails(txs, txReceipts)
		return err
	})
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (s *Storage) getTransfersByPublicKey(publicKey *models.PublicKey) ([]stored.Tx, []*stored.TxReceipt, error) {
	leaves, err := s.GetStateLeavesByPublicKey(publicKey)
	if IsNotFoundError(err) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	stateIDs := utils.ValueToInterfaceSlice(leaves, "StateID")

	txs := make([]stored.Tx, 0, 1)
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

	txReceipts := make([]*stored.TxReceipt, 0, len(txs))
	for i := range txs {
		txReceipt, err := s.getStoredTxReceipt(txs[i].Hash)
		if err != nil {
			return nil, nil, err
		}
		txReceipts = append(txReceipts, txReceipt)
	}

	return txs, txReceipts, nil
}

func (s *Storage) txsToTransfersWithBatchDetails(txs []stored.Tx, txReceipts []*stored.TxReceipt) (
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
