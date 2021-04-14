package api

import (
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetTransactions(publicKey *models.PublicKey) ([]models.TransactionReceipt, error) {
	transactions, err := a.storage.GetTransactionsByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	userTransactions := make([]models.TransactionReceipt, 0, len(transactions))
	for i := range transactions {
		status, err := CalculateTransactionStatus(a.storage, &transactions[i], commander.LatestBlockNumber)
		if err != nil {
			return nil, err
		}
		returnTx := &models.TransactionReceipt{
			Transaction: transactions[i],
			Status:      *status,
		}
		userTransactions = append(userTransactions, *returnTx)
	}

	return userTransactions, nil
}
