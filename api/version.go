package api

func (a *Api) GetVersion() string {
	return a.cfg.Version
}
