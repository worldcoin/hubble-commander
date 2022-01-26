package api

type ConfigureParams struct {
	CreateBatches      *bool
	AcceptTransactions *bool
}

func (a *API) Configure(params ConfigureParams) error {
	if params.CreateBatches != nil {
		a.enableBatchCreation(*params.CreateBatches)
	}
	if params.AcceptTransactions != nil {
		a.isAcceptingTransactions = *params.AcceptTransactions
	}
	return nil
}
