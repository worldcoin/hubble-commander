package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *TransactionStorage) BatchAddMassMigration(txs []models.MassMigration) error {
	return s.BatchAddTransaction(models.MakeMassMigrationArray(txs...))
}

func (s *TransactionStorage) GetMassMigration(hash common.Hash) (*models.MassMigration, error) {
	tx, err := s.getTransactionByHash(hash)
	if err != nil {
		return nil, err
	}
	if tx.Type() != txtype.MassMigration {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	transfer := tx.ToMassMigration()
	return transfer, nil
}

func (s *TransactionStorage) GetPendingMassMigrations() (txs models.MassMigrationArray, err error) {
	genericTxs, err := s.GetPendingTransactions(txtype.MassMigration)
	if err != nil {
		return nil, err
	}
	return genericTxs.ToMassMigrationArray(), nil
}

func (s *TransactionStorage) GetMassMigrationsByCommitmentID(id models.CommitmentID) ([]models.MassMigration, error) {
	var batchedTxs []stored.BatchedTx

	query := bh.Where("CommitmentID").Eq(id).Index("CommitmentID").And("TxType").Eq(txtype.MassMigration)

	err := s.database.Badger.Find(&batchedTxs, query)
	if err != nil {
		return nil, err
	}

	txs := make([]models.MassMigration, len(batchedTxs))
	for i := range batchedTxs {
		txs[i] = *batchedTxs[i].ToMassMigration()
	}

	return txs, nil
}

func (s *TransactionStorage) MarkMassMigrationsAsIncluded(txs []models.MassMigration, commitmentID *models.CommitmentID) error {
	return s.MarkTransactionsAsIncluded(models.MakeMassMigrationArray(txs...), commitmentID)
}
