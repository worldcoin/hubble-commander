package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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

func (s *TransactionStorage) MarkTransfersAsIncluded(txs []models.Transfer, commitmentID *models.CommitmentID) error {
	return s.MarkTransactionsAsIncluded(models.MakeTransferArray(txs...), commitmentID)
}
