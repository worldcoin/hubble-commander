package api

func (a *API) GetVersion() string {
	return a.cfg.Version
}
