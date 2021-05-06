package api

import (
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) GetBatchByHash(hash common.Hash) (*dto.BatchWithCommitments, error) {
	batch, err := a.storage.GetBatchWithAccountRoot(hash)
	if err != nil {
		return nil, err
	}
	batch.SubmissionBlock = getSubmissionBlock(batch.FinalisationBlock)

	commitments, err := a.storage.GetCommitmentsByBatchHash(&hash)
	if err != nil {
		return nil, err
	}
	return createBatchWithCommitments(batch, commitments)
}

func (a *API) GetBatchByID(id models.Uint256) (*dto.BatchWithCommitments, error) {
	batch, err := a.storage.GetBatchWithAccountRootByID(id)
	if err != nil {
		return nil, err
	}
	batch.SubmissionBlock = getSubmissionBlock(batch.FinalisationBlock)

	commitments, err := a.storage.GetCommitmentsByBatchID(id)
	if err != nil {
		return nil, err
	}
	return createBatchWithCommitments(batch, commitments)
}

func createBatchWithCommitments(
	batch *models.BatchWithAccountRoot,
	commitments []models.CommitmentWithTokenID,
) (*dto.BatchWithCommitments, error) {
	for i := range commitments {
		commitments[i].LeafHash = commitments[i].CalcLeafHash(batch.AccountTreeRoot)
	}
	return &dto.BatchWithCommitments{
		BatchWithAccountRoot: *batch,
		Commitments:          commitments,
	}, nil
}

func getSubmissionBlock(finalisationBlock uint32) uint32 {
	return finalisationBlock - rollup.BlocksToFinalise
}
