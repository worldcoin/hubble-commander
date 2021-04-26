package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetTransactions(publicKey *models.PublicKey) ([]interface{}, error) {
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
		status, err := CalculateTransferStatus(a.storage, &transfers[i].TransactionBase, a.storage.GetLatestBlockNumber())
		if err != nil {
			return nil, err
		}
		userTransfers = append(userTransfers, &dto.TransferReceipt{
			Transfer: transfers[i],
			Status:   *status,
		})
	}

	for i := range create2Transfers {
		status, err := CalculateTransferStatus(a.storage, &create2Transfers[i], a.storage.GetLatestBlockNumber())
		if err != nil {
			return nil, err
		}
		userTransfers = append(userTransfers, &dto.Create2TransferReceipt{
			Create2Transfer: transfers[i],
			Status:          *status,
		})
	}

	return userTransfers, nil
}
