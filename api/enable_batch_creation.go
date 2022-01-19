package api

func (a *API) EnableBatchCreation(enable bool) error {
	a.enableBatchCreation(enable)
	return nil
}
