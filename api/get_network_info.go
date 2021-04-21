package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

func (a *API) GetNetworkInfo() (*models.NetworkInfo, error) {
	networkInfo := models.NetworkInfo{
		ChainState:  a.client.ChainState,
		BlockNumber: a.storage.GetLatestBlockNumber(),
	}

	latestBatch, err := a.storage.GetLatestBatch()
	if err != nil && !storage.IsNotFoundError(err) {
		return nil, err
	}
	if latestBatch != nil {
		networkInfo.LatestBatch = ref.String(latestBatch.ID.String())
	}

	latestFinalisedBatch, err := a.storage.GetLatestFinalisedBatch(networkInfo.BlockNumber)
	if err != nil && !storage.IsNotFoundError(err) {
		return nil, err
	}
	if latestFinalisedBatch != nil {
		networkInfo.LatestFinalisedBatch = ref.String(latestFinalisedBatch.ID.String())
	}

	return &networkInfo, nil
}
