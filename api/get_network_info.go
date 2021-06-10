package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

func (a *API) GetNetworkInfo() (*dto.NetworkInfo, error) {
	networkInfo := dto.NetworkInfo{
		ChainState:  a.client.ChainState,
		BlockNumber: a.storage.GetLatestBlockNumber(),
	}

	latestBatch, err := a.storage.GetLatestSubmittedBatch()
	if err != nil && !storage.IsNotFoundError(err) {
		return nil, err
	}
	if latestBatch != nil {
		networkInfo.LatestBatch = &latestBatch.Number
	}

	latestFinalisedBatch, err := a.storage.GetLatestFinalisedBatch(networkInfo.BlockNumber)
	if err != nil && !storage.IsNotFoundError(err) {
		return nil, err
	}
	if latestFinalisedBatch != nil {
		networkInfo.LatestFinalisedBatch = &latestFinalisedBatch.Number
	}

	domain, err := a.storage.GetDomain(a.client.ChainState.ChainID)
	if err != nil {
		return nil, err
	}
	networkInfo.SignatureDomain = *domain

	// TODO replace with a more effective approach when we get to a huge number of txs
	txCount, err := a.storage.GetTransactionCount()
	if err != nil {
		return nil, err
	}
	networkInfo.TransactionCount = *txCount

	// TODO this ignores the fact that other nodes can put new accounts in arbitrary state leaves; to be revisited in the future
	accountCount, err := a.storage.GetNextAvailableStateID()
	if err != nil {
		return nil, err
	}
	networkInfo.AccountCount = *accountCount

	return &networkInfo, nil
}
