package api

import (
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetNetworkInfo() (models.NetworkInfo, error) {
	return models.NetworkInfo{
		ChainState:  a.client.ChainState,
		BlockNumber: commander.LatestBlockNumber,
	}, nil
}
