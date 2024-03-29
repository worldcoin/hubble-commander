package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetFailedTransactions(ctx context.Context) (models.GenericTransactionArray, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

	return a.storage.GetAllFailedTransactions()
}
