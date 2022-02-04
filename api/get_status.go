package api

import "github.com/Worldcoin/hubble-commander/models/enums/healthstatus"

func (a *API) GetStatus() string {
	if a.isMigrating() {
		return healthstatus.Migrating
	}
	return healthstatus.Ready
}
