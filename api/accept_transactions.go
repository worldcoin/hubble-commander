package api

func (a *API) AcceptTransactions(accept bool) {
	a.disableSendTransaction = !accept
}
