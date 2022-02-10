package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getCommitmentAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(20000, "commitment not found"),
}

func (a *API) GetCommitment(id models.CommitmentID) (interface{}, error) {
	commitment, err := a.unsafeGetCommitment(id)
	if err != nil {
		return nil, sanitizeError(err, getCommitmentAPIErrors)
	}

	return commitment, nil
}

func (a *API) unsafeGetCommitment(id models.CommitmentID) (interface{}, error) {
	commitment, err := a.storage.GetCommitment(&id)
	if err != nil {
		return nil, err
	}

	batch, err := a.storage.GetBatch(commitment.GetCommitmentBase().ID.BatchID)
	if err != nil {
		return nil, err
	}

	return a.createCommitmentDTO(commitment, batch)
}

func (a *API) createCommitmentDTO(commitment models.Commitment, batch *models.Batch) (interface{}, error) {
	transactions, err := a.getTransactionsForCommitment(commitment)
	if err != nil {
		return nil, err
	}

	status := calculateBatchStatus(a.storage.GetLatestBlockNumber(), batch)

	switch batch.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return a.createTxCommitmentDTO(commitment, batch, transactions, status)
	case batchtype.MassMigration:
		return dto.NewMMCommitment(commitment.ToMMCommitment(), status, batch.MinedTime, transactions), nil
	case batchtype.Deposit:
		return dto.NewDepositCommitment(commitment.ToDepositCommitment(), status, batch.MinedTime), nil
	default:
		panic("invalid commitment type")
	}
}

func (a *API) getTransactionsForCommitment(commitment models.Commitment) (interface{}, error) {
	commitmentBase := commitment.GetCommitmentBase()
	switch commitmentBase.Type {
	case batchtype.Transfer, batchtype.Create2Transfer, batchtype.MassMigration:
		return a.innerGetTransactionsForCommitment(commitmentBase.ID)
	case batchtype.Deposit:
		return nil, nil
	case batchtype.Genesis:
		return nil, dto.ErrNotImplemented
	}
	return nil, dto.ErrNotImplemented
}

func (a *API) innerGetTransactionsForCommitment(id models.CommitmentID) (interface{}, error) {
	txns, err := a.storage.GetTransactionsByCommitmentID(id)
	if err != nil {
		return nil, err
	}

	txs := make([]interface{}, 0, txns.Len())
	for i := 0; i < txns.Len(); i++ {
		txs = append(txs, dto.MakeTransactionForCommitment(txns.At(i)))
	}
	return txs, nil
}

func (a *API) createTxCommitmentDTO(
	commitment models.Commitment,
	batch *models.Batch,
	transactions interface{},
	status *batchstatus.BatchStatus,
) (interface{}, error) {
	stateLeaf, err := a.storage.StateTree.Leaf(commitment.ToTxCommitment().FeeReceiver)
	if err != nil {
		return nil, err
	}

	commitmentDTO := dto.NewTxCommitment(commitment.ToTxCommitment(), stateLeaf.TokenID, status, batch.MinedTime, transactions)

	return commitmentDTO, nil
}
