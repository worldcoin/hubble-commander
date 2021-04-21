package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetTransfers(publicKey *models.PublicKey) ([]models.TransferReceipt, error) {
	transfers, err := a.storage.GetTransfersByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	userTransfers := make([]models.TransferReceipt, 0, len(transfers))
	for i := range transfers {
		status, err := CalculateTransferStatus(a.storage, &transfers[i], a.storage.GetLatestBlockNumber())
		if err != nil {
			return nil, err
		}
		returnTransfer := &models.TransferReceipt{
			Transfer: transfers[i],
			Status:   *status,
		}
		userTransfers = append(userTransfers, *returnTransfer)
	}

	return userTransfers, nil
}
