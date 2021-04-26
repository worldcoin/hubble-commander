package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetTransactions(publicKey *models.PublicKey) ([]dto.TransferReceipt, error) {
	transfers, err := a.storage.GetTransfersByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	userTransfers := make([]dto.TransferReceipt, 0, len(transfers))
	for i := range transfers {
		status, err := CalculateTransferStatus(a.storage, &transfers[i].TransactionBase, a.storage.GetLatestBlockNumber())
		if err != nil {
			return nil, err
		}
		returnTransfer := &dto.TransferReceipt{
			Transfer: transfers[i],
			Status:   *status,
		}
		userTransfers = append(userTransfers, *returnTransfer)
	}

	return userTransfers, nil
}
