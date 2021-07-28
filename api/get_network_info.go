package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

func (a *API) GetNetworkInfo() (*dto.NetworkInfo, error) {
	networkInfo := dto.NetworkInfo{
		ChainID:         a.client.ChainState.ChainID,
		AccountRegistry: a.client.ChainState.AccountRegistry,
		DeploymentBlock: a.client.ChainState.DeploymentBlock,
		Rollup:          a.client.ChainState.Rollup,
		BlockNumber:     a.storage.GetLatestBlockNumber(),
	}

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

	latestBatch, err := a.storage.GetLatestSubmittedBatch()
	if err != nil && !storage.IsNotFoundError(err) {
		return nil, err
	}
	if latestBatch != nil {
		networkInfo.LatestBatch = &latestBatch.ID
	}

	latestFinalisedBatch, err := a.storage.GetLatestFinalisedBatch(networkInfo.BlockNumber)
	if err != nil && !storage.IsNotFoundError(err) {
		return nil, err
	}
	if latestFinalisedBatch != nil {
		networkInfo.LatestFinalisedBatch = &latestFinalisedBatch.ID
	}

	domain, err := a.storage.GetDomain()
	if err != nil {
		return nil, err
	}
	networkInfo.SignatureDomain = *domain

	return &networkInfo, nil
}
