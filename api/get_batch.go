package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) GetBatchByHash(hash common.Hash) (*dto.BatchWithCommitments, error) {
	batch, err := a.storage.GetBatch(hash)
	if err != nil {
		return nil, err
	}

	batchWithCommitments := &dto.BatchWithCommitments{
		Batch: *batch,
	}

	batchWithCommitments.Commitments, err = a.storage.GetCommitmentsByBatchHash(&hash)
	if err != nil {
		return nil, err
	}
	return batchWithCommitments, nil
}

func (a *API) GetBatchByID(id models.Uint256) (*dto.BatchWithCommitments, error) {
	batch, err := a.storage.GetBatchByID(id)
	if err != nil {
		return nil, err
	}

	batchWithCommitments := &dto.BatchWithCommitments{
		Batch: *batch,
	}

	batchWithCommitments.Commitments, err = a.storage.GetCommitmentsByBatchID(id)
	if err != nil {
		return nil, err
	}
	return batchWithCommitments, nil
}
