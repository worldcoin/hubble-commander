package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

func (a *API) GetNetworkInfo() (models.NetworkInfo, error) {
	networkInfo := models.NetworkInfo{
		ChainState:  a.client.ChainState,
		BlockNumber: a.storage.GetLatestBlockNumber(),
	}

	latestBatch, err := a.storage.GetLatestBatch()
	if err != nil && !storage.IsNotFoundError(err) {
		return models.NetworkInfo{}, err
	}
	if latestBatch != nil {
		networkInfo.LatestBatch = latestBatch.ID.String()
	}

	latestFinalisedBatch, err := a.storage.GetLatestFinalisedBatch(networkInfo.BlockNumber)
	if err != nil && !storage.IsNotFoundError(err) {
		return models.NetworkInfo{}, err
	}
	if latestFinalisedBatch != nil {
		networkInfo.LatestFinalisedBatch = latestFinalisedBatch.ID.String()
	}

	return networkInfo, nil
}
