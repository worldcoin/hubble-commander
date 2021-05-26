package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

func (a *API) GetNetworkInfo() (*dto.NetworkInfo, error) {
	networkInfo := dto.NetworkInfo{
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

	domain, err := a.storage.GetDomain(a.client.ChainState.ChainID)
	if err != nil {
		return nil, err
	}
	networkInfo.SignatureDomain = *domain

	txCount, err := a.storage.GetTransactionsCount()
	if err != nil {
		return nil, err
	}
	networkInfo.TransactionCount = *txCount

	//TODO: sync with other nodes
	accountCount, err := a.storage.GetNextAvailableStateID()
	if err != nil {
		return nil, err
	}
	networkInfo.AccountCount = *accountCount

	return &networkInfo, nil
}
