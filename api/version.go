package api

/// GetVersion returns commander's version.
func (a *API) GetVersion() string {
	return a.cfg.Version
}
