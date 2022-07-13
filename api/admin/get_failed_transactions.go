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

	// TODO: should this inspect the mempool? As an admin_ method it's not user-facing
	//       so no need to get this perfect on the first pass

	return a.storage.GetAllFailedTransactions()
}
