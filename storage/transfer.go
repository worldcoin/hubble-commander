package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *TransactionStorage) BatchAddTransfer(txs []models.Transfer) error {
	return s.BatchAddTransaction(models.MakeTransferArray(txs...))
}

// called by a large amount of tests, and nothing else
func (s *TransactionStorage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	tx, err := s.getTransactionByHash(hash)
	if err != nil {
		return nil, err
	}
	if tx.Type() != txtype.Transfer {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	transfer := tx.ToTransfer()
	return transfer, nil
}

func (s *TransactionStorage) GetPendingTransfers() (txs models.TransferArray, err error) {
	genericTxs, err := s.GetPendingTransactions(txtype.Transfer)
	if err != nil {
		return nil, err
	}
	return genericTxs.ToTransferArray(), nil
}

func (s *TransactionStorage) GetTransfersByCommitmentID(id models.CommitmentID) ([]models.Transfer, error) {
	var batchedTxs []stored.BatchedTx

	query := bh.Where("CommitmentID").Eq(id).Index("CommitmentID").And("TxType").Eq(txtype.Transfer)

	err := s.database.Badger.Find(&batchedTxs, query)
	if err != nil {
		return nil, err
	}

	txs := make([]models.Transfer, len(batchedTxs))
	for i := range batchedTxs {
		txs[i] = *batchedTxs[i].ToTransfer()
	}

	return txs, nil
}

func (s *TransactionStorage) MarkTransfersAsIncluded(txs []models.Transfer, commitmentID *models.CommitmentID) error {
	return s.MarkTransactionsAsIncluded(models.MakeTransferArray(txs...), commitmentID)
}
