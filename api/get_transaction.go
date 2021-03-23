package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) GetTransaction(hash common.Hash) (*models.TransactionReceipt, error) {
	tx, err := a.storage.GetTransaction(hash)
	if err != nil {
		return nil, err
	}

	status := CalculateTransactionStatus(tx)

	returnTx := &models.TransactionReceipt{
		Transaction: *tx,
		Status:      status,
	}

	return returnTx, nil
}

func CalculateTransactionStatus(tx *models.Transaction) models.TransactionStatus {
	var status models.TransactionStatus

	if tx.IncludedInCommitment == nil {
		status = models.Pending
	} else {
		status = models.Committed
	}

	if tx.ErrorMessage != nil {
		status = models.Error
	}

	return status
}
