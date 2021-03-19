package api

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetTransactions(publicKey *models.PublicKey) ([]models.ReturnTransaction, error) {
	accounts, err := a.storage.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}

	userStatesIndexes := make([]models.Uint256, 0, 1)

	for i := range accounts {
		stateLeafs, err := a.storage.GetStateLeafs(accounts[i].AccountIndex)
		if err != nil {
			return nil, err
		}

		for i := range stateLeafs {
			node, err := a.storage.GetStateNodeByHash(stateLeafs[i].DataHash)
			if err != nil {
				return nil, err
			}

			userStatesIndexes = append(userStatesIndexes, models.MakeUint256FromBig(*new(big.Int).SetInt64(int64(node.MerklePath.Path))))
		}
	}

	userTransactions := make([]models.ReturnTransaction, 0, 1)

	for i := range userStatesIndexes {
		transactions, err := a.storage.GetUserTransactions(userStatesIndexes[i])
		if err != nil {
			return nil, err
		}

		for i := range transactions {
			status := CalculateTransactionStatus(&transactions[i])
			returnTx := &models.ReturnTransaction{
				Transaction: transactions[i],
				Status:      status,
			}
			userTransactions = append(userTransactions, *returnTx)
		}
	}

	return userTransactions, nil
}
