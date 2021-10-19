package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var getCommitmentAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(20000, "commitment not found"),
}

func (a *API) GetCommitment(id models.CommitmentID) (*dto.Commitment, error) {
	commitment, err := a.unsafeGetCommitment(id)
	if err != nil {
		return nil, sanitizeError(err, getCommitmentAPIErrors)
	}

	return commitment, nil
}

func (a *API) unsafeGetCommitment(id models.CommitmentID) (*dto.Commitment, error) {
	commitment, err := a.storage.GetTxCommitment(&id)
	if err != nil {
		return nil, err
	}

	transactions, err := a.getTransactionsForCommitment(commitment)
	if err != nil {
		return nil, err
	}

	batch, err := a.storage.GetMinedBatch(commitment.ID.BatchID)
	if storage.IsNotFoundError(err) {
		return nil, errors.WithStack(storage.NewNotFoundError("commitment"))
	}
	if err != nil {
		return nil, err
	}

	return &dto.Commitment{
		TxCommitment: *commitment,
		Status:       *calculateFinalisedStatus(a.storage.GetLatestBlockNumber(), *batch.FinalisationBlock),
		BatchTime:    batch.SubmissionTime,
		Transactions: transactions,
	}, nil
}

func (a *API) getTransactionsForCommitment(commitment *models.TxCommitment) (interface{}, error) {
	switch commitment.Type {
	case batchtype.Transfer:
		return a.getTransfersForCommitment(&commitment.ID)
	case batchtype.Create2Transfer:
		return a.getCreate2TransfersForCommitment(&commitment.ID)
	case batchtype.Genesis, batchtype.MassMigration, batchtype.Deposit:
		return nil, dto.ErrNotImplemented
	}
	return nil, dto.ErrNotImplemented
}

func (a *API) getTransfersForCommitment(id *models.CommitmentID) (interface{}, error) {
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

func (a *API) getCreate2TransfersForCommitment(id *models.CommitmentID) (interface{}, error) {
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
