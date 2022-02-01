package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) Configure(ctx context.Context, params dto.ConfigureParams) error {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return err
	}

	if params.CreateBatches != nil {
		a.enableBatchCreation(*params.CreateBatches)
	}
	if params.AcceptTransactions != nil {
		a.enableTxsAcceptance(*params.AcceptTransactions)
	}
	return nil
}
