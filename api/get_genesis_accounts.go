package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetGenesisAccounts() *models.GenesisAccounts {
	return &a.client.ChainState.GenesisAccounts
}
