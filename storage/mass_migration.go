package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (s *TransactionStorage) AddMassMigration(m *models.MassMigration) error {
	return s.executeInTransaction(TxOptions{}, func(txStorage *TransactionStorage) error {
		if m.CommitmentID != nil || m.ErrorMessage != nil {
			err := txStorage.addStoredTxReceipt(models.NewStoredTxReceiptFromMassMigration(m))
			if err != nil {
				return err
			}
		}
		return txStorage.addStoredTx(models.NewStoredTxFromMassMigration(m))
	})
}

func (s *TransactionStorage) GetMassMigration(hash common.Hash) (*models.MassMigration, error) {
	tx, txReceipt, err := s.getStoredTxWithReceipt(hash)
	if err != nil {
		return nil, err
	}
	if tx.TxType != txtype.MassMigration {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	return tx.ToMassMigration(txReceipt), nil
}
