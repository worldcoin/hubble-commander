package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (a *API) GetPendingBatches(ctx context.Context) ([]dto.PendingBatch, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

	batches, err := a.storage.GetPendingBatches()
	if err != nil {
		return nil, err
	}

	pendingBatches := make([]dto.PendingBatch, 0, len(batches))
	for i := range batches {
		commitments, err := a.getCommitments(&batches[i])
		if err != nil {
			return nil, err
		}
		pendingBatches = append(pendingBatches, dto.PendingBatch{
			ID:              batches[i].ID,
			Type:            batches[i].Type,
			TransactionHash: batches[i].TransactionHash,
			Commitments:     commitments,
		})
	}

	return pendingBatches, nil
}

func (a *API) getCommitments(batch *models.Batch) ([]dto.PendingCommitment, error) {
	commitments, err := a.storage.GetCommitmentsByBatchID(batch.ID)
	if err != nil {
		return nil, err
	}

	dtoCommitments := make([]dto.PendingCommitment, 0, len(commitments))
	for i := range commitments {
		txs, err := a.getTransactionsForCommitment(commitments[i])
		if err != nil {
			return nil, err
		}

		// TODO remove when new primary key for transactions with transaction index is implement
		txQueue := executor.NewTxQueue(txs)

		dtoCommitments = append(dtoCommitments, dto.PendingCommitment{
			Commitment:   commitments[i],
			Transactions: txQueue.PickTxsForCommitment(),
		})
	}
	return dtoCommitments, nil
}

func (a *API) getTransactionsForCommitment(commitment models.Commitment) (models.GenericTransactionArray, error) {
	commitmentBase := commitment.GetCommitmentBase()
	switch commitmentBase.Type {
	case batchtype.Transfer:
		return a.storage.GetTransfersByCommitmentID(commitmentBase.ID)
	case batchtype.Create2Transfer:
		return a.storage.GetCreate2TransfersByCommitmentID(commitmentBase.ID)
	case batchtype.MassMigration:
		return a.storage.GetMassMigrationsByCommitmentID(commitmentBase.ID)
	case batchtype.Deposit:
		return nil, nil
	case batchtype.Genesis:
		return nil, dto.ErrNotImplemented
	}
	return nil, dto.ErrNotImplemented
}
