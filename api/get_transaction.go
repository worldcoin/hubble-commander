package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

var getTransactionAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(10000, "transaction not found"),
}

func (a *API) GetTransaction(hash common.Hash) (interface{}, error) {
	transaction, err := a.unsafeGetTransaction(hash)
	if err != nil {
		return nil, sanitizeError(err, getTransactionAPIErrors)
	}

	return transaction, nil
}

func (a *API) unsafeGetTransaction(hash common.Hash) (interface{}, error) {
	transaction, txType, err := a.storage.GetTransactionWithBatchDetails(hash)
	if err != nil {
		return nil, err
	}

	var transactionBase models.TransactionBase

	switch *txType {
	case txtype.Transfer:
		transactionBase = transaction.Transaction.(*models.Transfer).TransactionBase
	case txtype.Create2Transfer:
		transactionBase = transaction.Transaction.(*models.Create2Transfer).TransactionBase
	case txtype.MassMigration:
		transactionBase = transaction.Transaction.(*models.MassMigration).TransactionBase
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
