package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetPendingTransactions(ctx context.Context) (models.GenericTransactionArray, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

	mempoolTxs, err := a.storage.GetAllMempoolTransactions()
	if err != nil {
		return nil, err
	}

	result := make([]models.GenericTransaction, len(mempoolTxs))
	for i := range mempoolTxs {
		result[i] = mempoolTxs[i].ToGenericTransaction()
	}

	return models.MakeGenericArray(result...), nil
}
