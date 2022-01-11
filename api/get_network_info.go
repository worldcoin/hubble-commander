package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

var networkInfoAPIErrors = map[error]*APIError{
	storage.NewNoVacantSubtreeError(0): NewAPIError(
		99000,
		"an error occurred while fetching the account count",
	),
}

func (a *API) GetNetworkInfo() (*dto.NetworkInfo, error) {
	networkInfo, err := a.unsafeGetNetworkInfo()
	if err != nil {
		return nil, sanitizeError(err, networkInfoAPIErrors)
	}

	return networkInfo, nil
}

func (a *API) unsafeGetNetworkInfo() (*dto.NetworkInfo, error) {
	networkInfo := dto.NetworkInfo{
		ChainID:                        a.client.ChainState.ChainID,
		AccountRegistry:                a.client.ChainState.AccountRegistry,
		AccountRegistryDeploymentBlock: a.client.ChainState.AccountRegistryDeploymentBlock,
		TokenRegistry:                  a.client.ChainState.TokenRegistry,
		SpokeRegistry:                  a.client.ChainState.SpokeRegistry,
		DepositManager:                 a.client.ChainState.DepositManager,
		WithdrawManager:                a.client.ChainState.WithdrawManager,
		Rollup:                         a.client.ChainState.Rollup,
		BlockNumber:                    a.storage.GetLatestBlockNumber(),
	}

	// TODO replace with a more effective approach when we get to a huge number of txs
	txCount, err := a.storage.GetTransactionCount()
	if err != nil {
		return nil, err
	}
	networkInfo.TransactionCount = *txCount

	// TODO this ignores the fact that other nodes can put new accounts in arbitrary state leaves; to be revisited in the future
	accountCount, err := a.storage.StateTree.NextAvailableStateID()
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

	domain, err := a.client.GetDomain()
	if err != nil {
		return nil, err
	}
	networkInfo.SignatureDomain = *domain

	return &networkInfo, nil
}
