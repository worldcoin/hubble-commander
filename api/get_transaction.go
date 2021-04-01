package api

import (
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) GetTransaction(hash common.Hash) (*models.TransactionReceipt, error) {
	tx, err := a.storage.GetTransaction(hash)
	if err != nil {
		return nil, err
	}

	status, err := CalculateTransactionStatus(a.storage, tx, commander.LatestBlockNumber)
	if err != nil {
		return nil, err
	}

	returnTx := &models.TransactionReceipt{
		Transaction: *tx,
		Status:      *status,
	}

	return returnTx, nil
}
