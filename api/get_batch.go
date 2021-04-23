package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) GetBatchByHash(hash common.Hash) (*dto.BatchWithCommitments, error) {
	batch, err := a.storage.GetBatchWithAccountRoot(hash)
	if err != nil {
		return nil, err
	}

	commitments, err := a.storage.GetCommitmentsByBatchHash(&hash)
	if err != nil {
		return nil, err
	}
	return &dto.BatchWithCommitments{
		BatchWithAccountRoot: *batch,
		Commitments:          commitments,
	}, nil
}

func (a *API) GetBatchByID(id models.Uint256) (*dto.BatchWithCommitments, error) {
	batch, err := a.storage.GetBatchWithAccountRootByID(id)
	if err != nil {
		return nil, err
	}

	commitments, err := a.storage.GetCommitmentsByBatchID(id)
	if err != nil {
		return nil, err
	}
	return &dto.BatchWithCommitments{
		BatchWithAccountRoot: *batch,
		Commitments:          commitments,
	}, nil
}
