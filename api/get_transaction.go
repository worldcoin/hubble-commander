package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

var getTransactionAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(10000, "transaction not found"),
}

func (a *API) GetTransaction(hash common.Hash) (*dto.TransactionReceipt, error) {
	transaction, err := a.unsafeGetTransaction(hash)
	if err != nil {
		return nil, sanitizeError(err, getTransactionAPIErrors)
	}

	return transaction, nil
}

func (a *API) unsafeGetTransaction(hash common.Hash) (*dto.TransactionReceipt, error) {
	transaction, err := a.storage.GetTransactionWithBatchDetails(hash)
	if err != nil {
		return nil, err
	}

	var transactionBase models.TransactionBase

	switch tx := transaction.Transaction.(type) {
	case *models.Transfer:
		transactionBase = tx.TransactionBase
	case *models.Create2Transfer:
		transactionBase = tx.TransactionBase
	case *models.MassMigration:
		transactionBase = tx.TransactionBase
	}

	status, err := CalculateTransactionStatus(a.storage, &transactionBase, a.storage.GetLatestBlockNumber())
	if err != nil {
		return nil, err
	}

	return &dto.TransactionReceipt{
		TransactionWithBatchDetails: *transaction,
		Status:                      *status,
	}, nil
}
