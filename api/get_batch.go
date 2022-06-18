package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
)

var getBatchAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(30000, "batch not found"),
}

func (a *API) GetBatchByHash(hash common.Hash) (*dto.BatchWithRootAndCommitments, error) {
	batch, err := a.unsafeGetBatchByHash(hash)
	if err != nil {
		return nil, sanitizeError(err, getBatchAPIErrors)
	}

	return batch, nil
}

func (a *API) unsafeGetBatchByHash(hash common.Hash) (*dto.BatchWithRootAndCommitments, error) {
	// span
	batch, err := a.storage.GetBatchByHash(hash)
	if err != nil {
		return nil, err
	}

	return a.getCommitmentsAndCreateBatchDTO(batch)
}

func (a *API) GetBatchByID(id models.Uint256) (*dto.BatchWithRootAndCommitments, error) {
	batch, err := a.unsafeGetBatchByID(id)
	if err != nil {
		return nil, sanitizeError(err, getBatchAPIErrors)
	}

	return batch, nil
}

func (a *API) unsafeGetBatchByID(id models.Uint256) (*dto.BatchWithRootAndCommitments, error) {
	batch, err := a.storage.GetBatch(id)
	if err != nil {
		return nil, err
	}

	return a.getCommitmentsAndCreateBatchDTO(batch)
}

func (a *API) getCommitmentsAndCreateBatchDTO(batch *models.Batch) (*dto.BatchWithRootAndCommitments, error) {
	minedBlock, err := a.getMinedBlock(batch)
	if err != nil {
		return nil, err
	}

	if batch.Type == batchtype.Genesis {
		return a.createBatchWithCommitments(batch, minedBlock, batchstatus.Finalised.Ref(), nil)
	}

	status := calculateBatchStatus(a.storage.GetLatestBlockNumber(), batch)

	commitments, err := a.storage.GetCommitmentsByBatchID(batch.ID)
	if err != nil {
		return nil, err
	}

	return a.createBatchWithCommitments(batch, minedBlock, status, commitments)
}

func (a *API) createBatchWithCommitments(
	batch *models.Batch,
	minedBlock *uint32,
	status *batchstatus.BatchStatus,
	commitments []models.Commitment,
) (*dto.BatchWithRootAndCommitments, error) {
	switch batch.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return a.createBatchWithTxCommitments(batch, minedBlock, status, commitments)
	case batchtype.MassMigration:
		return a.createBatchWithMMCommitments(batch, minedBlock, status, commitments)
	case batchtype.Deposit:
		return a.createBatchWithDepositCommitments(batch, minedBlock, status, commitments)
	default:
		return dto.MakeBatchWithRootAndCommitments(dto.NewBatch(batch, minedBlock, status), batch.AccountTreeRoot, nil), nil
	}
}

func (a *API) createBatchWithTxCommitments(
	batch *models.Batch,
	minedBlock *uint32,
	status *batchstatus.BatchStatus,
	commitments []models.Commitment,
) (*dto.BatchWithRootAndCommitments, error) {
	batchCommitments := make([]dto.BatchTxCommitment, 0, len(commitments))
	for i := range commitments {
		stateLeaf, err := a.storage.StateTree.Leaf(commitments[i].ToTxCommitment().FeeReceiver)
		if err != nil {
			return nil, err
		}

		batchCommitments = append(batchCommitments, dto.MakeBatchTxCommitment(
			commitments[i].ToTxCommitment(),
			stateLeaf.TokenID,
		))
	}
	return createBatchWithRootAndCommitmentsDTO(batch, minedBlock, status, batchCommitments), nil
}

func (a *API) createBatchWithMMCommitments(
	batch *models.Batch,
	minedBlock *uint32,
	status *batchstatus.BatchStatus,
	commitments []models.Commitment,
) (*dto.BatchWithRootAndCommitments, error) {
	batchCommitments := make([]dto.BatchMMCommitment, 0, len(commitments))
	for i := range commitments {
		batchCommitments = append(batchCommitments, dto.MakeBatchMMCommitment(
			commitments[i].ToMMCommitment(),
		))
	}
	return createBatchWithRootAndCommitmentsDTO(batch, minedBlock, status, batchCommitments), nil
}

func (a *API) createBatchWithDepositCommitments(
	batch *models.Batch,
	minedBlock *uint32,
	status *batchstatus.BatchStatus,
	commitments []models.Commitment,
) (*dto.BatchWithRootAndCommitments, error) {
	batchCommitments := make([]dto.BatchDepositCommitment, 0, len(commitments))
	for i := range commitments {
		batchCommitments = append(batchCommitments, dto.MakeBatchDepositCommitment(
			commitments[i].ToDepositCommitment(),
		))
	}
	return createBatchWithRootAndCommitmentsDTO(batch, minedBlock, status, batchCommitments), nil
}

func (a *API) getMinedBlock(batch *models.Batch) (*uint32, error) {
	if batch.ID.IsZero() {
		return batch.FinalisationBlock, nil
	}

	// Submitted batch
	if batch.FinalisationBlock == nil {
		return nil, nil
	}

	blocks, err := a.client.GetBlocksToFinalise()
	if err != nil {
		return nil, err
	}
	return ref.Uint32(*batch.FinalisationBlock - uint32(*blocks)), nil
}

func calculateBatchStatus(latestBlockNumber uint32, batch *models.Batch) *batchstatus.BatchStatus {
	if batch.FinalisationBlock == nil {
		return batchstatus.Submitted.Ref()
	}

	if latestBlockNumber < *batch.FinalisationBlock {
		return batchstatus.Mined.Ref()
	}

	return batchstatus.Finalised.Ref()
}

func createBatchWithRootAndCommitmentsDTO(
	batch *models.Batch,
	minedBlock *uint32,
	status *batchstatus.BatchStatus,
	commitments interface{},
) *dto.BatchWithRootAndCommitments {
	if batch.FinalisationBlock == nil {
		return dto.MakeBatchWithRootAndCommitments(
			dto.NewSubmittedBatch(batch),
			nil,
			commitments,
		)
	}

	return dto.MakeBatchWithRootAndCommitments(
		dto.NewBatch(batch, minedBlock, status),
		batch.AccountTreeRoot,
		commitments,
	)
}
