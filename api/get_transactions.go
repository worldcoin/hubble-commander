package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getTransactionsAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(10001, "transactions not found"),
}

func (a *API) GetTransactions(publicKey *models.PublicKey) ([]interface{}, error) {
	batch, err := a.unsafeGetTransactions(publicKey)
	if err != nil {
		return nil, sanitizeError(err, getTransactionsAPIErrors, a.cfg.Log.Level)
	}

	return batch, nil
}

func (a *API) unsafeGetTransactions(publicKey *models.PublicKey) ([]interface{}, error) {
	transfers, err := a.storage.GetTransfersByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	create2Transfers, err := a.storage.GetCreate2TransfersByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	userTransfers := make([]interface{}, 0, len(transfers)+len(create2Transfers))
	for i := range transfers {
		receipt, err := a.returnTransferReceipt(&transfers[i])
		if err != nil {
			return nil, err
		}
		userTransfers = append(userTransfers, receipt)
	}

	for i := range create2Transfers {
		receipt, err := a.returnCreate2TransferReceipt(&create2Transfers[i])
		if err != nil {
			return nil, err
		}
		userTransfers = append(userTransfers, receipt)
	}

	return userTransfers, nil
}
