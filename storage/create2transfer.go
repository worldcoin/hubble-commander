package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *TransactionStorage) BatchAddCreate2Transfer(txs []models.Create2Transfer) error {
	return s.BatchAddTransaction(models.MakeCreate2TransferArray(txs...))
}

func (s *TransactionStorage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	tx, err := s.getTransactionByHash(hash)
	if err != nil {
		return nil, err
	}
	if tx.Type() != txtype.Create2Transfer {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	transfer := tx.ToCreate2Transfer()
	return transfer, nil
}

func (s *TransactionStorage) GetCreate2TransfersByCommitmentID(id models.CommitmentID) ([]models.Create2Transfer, error) {
	batchedTxs := make([]stored.BatchedTx, 0, 32)

	query := bh.Where("CommitmentID").Eq(id).Index("CommitmentID")
	// We're not using `.And("TxType").Eq(txtype.Create2Transfer)` here because of inefficiency in BH implementation

	err := s.database.Badger.Find(&batchedTxs, query)
	if err != nil {
		return nil, err
	}

	txs := make([]models.Create2Transfer, 0, len(batchedTxs))
	for i := range batchedTxs {
		if batchedTxs[i].TxType == txtype.Create2Transfer {
			txs = append(txs, *batchedTxs[i].ToCreate2Transfer())
		}
	}

	return txs, nil
}

func (s *TransactionStorage) MarkCreate2TransfersAsIncluded(txs []models.Create2Transfer, commitmentID *models.CommitmentID) error {
	return s.MarkTransactionsAsIncluded(models.MakeCreate2TransferArray(txs...), commitmentID)
}
