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
	case batchtype.Transfer:
		return a.getTransfersForCommitment(commitmentBase.ID)
	case batchtype.Create2Transfer:
		return a.getCreate2TransfersForCommitment(commitmentBase.ID)
	case batchtype.MassMigration:
		return a.getMassMigrationsForCommitment(commitmentBase.ID)
	case batchtype.Deposit:
		return nil, nil
	case batchtype.Genesis:
		return nil, dto.ErrNotImplemented
	}
	return nil, dto.ErrNotImplemented
}

func (a *API) getTransfersForCommitment(id models.CommitmentID) (interface{}, error) {
	transfers, err := a.storage.GetTransfersByCommitmentID(id)
	if err != nil {
		return nil, err
	}

	txs := make([]dto.TransferForCommitment, 0, len(transfers))
	for i := range transfers {
		txs = append(txs, dto.MakeTransferForCommitment(&transfers[i]))
	}
	return txs, nil
}

func (a *API) getCreate2TransfersForCommitment(id models.CommitmentID) (interface{}, error) {
	transfers, err := a.storage.GetCreate2TransfersByCommitmentID(id)
	if err != nil {
		return nil, err
	}

	txs := make([]dto.Create2TransferForCommitment, 0, len(transfers))
	for i := range transfers {
		txs = append(txs, dto.MakeCreate2TransferForCommitment(&transfers[i]))
	}
	return txs, nil
}

func (a *API) getMassMigrationsForCommitment(id models.CommitmentID) (interface{}, error) {
	massMigrations, err := a.storage.GetMassMigrationsByCommitmentID(id)
	if err != nil {
		return nil, err
	}

	txs := make([]dto.MassMigrationForCommitment, 0, len(massMigrations))
	for i := range massMigrations {
		txs = append(txs, dto.MakeMassMigrationForCommitment(&massMigrations[i]))
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
