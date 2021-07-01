package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (a *API) GetCommitment(id int32) (*dto.Commitment, error) {
	commitment, err := a.storage.GetCommitment(id)
	if err != nil {
		return nil, err
	}

	transactions, err := a.getTransactionsForCommitment(commitment)
	if err != nil {
		return nil, err
	}

	batch, err := a.storage.GetMinedBatch(*commitment.IncludedInBatch)
	if st.IsNotFoundError(err) {
		return nil, st.NewNotFoundError("commitment")
	}
	if err != nil {
		return nil, err
	}

	return &dto.Commitment{
		Commitment:   *commitment,
		Status:       *calculateFinalisedStatus(a.storage.GetLatestBlockNumber(), *batch.FinalisationBlock),
		Transactions: transactions,
	}, nil
}

func (a *API) getTransactionsForCommitment(commitment *models.Commitment) (interface{}, error) {
	switch commitment.Type {
	case txtype.Transfer:
		return a.getTransfersForCommitment(commitment.ID)
	case txtype.Create2Transfer:
		return a.getC2TsForCommitment(commitment.ID)
	default:
		return nil, dto.ErrNotImplemented
	}
}

func (a *API) getTransfersForCommitment(commitmentID int32) ([]dto.TransferForCommitment, error) {
	transfers, err := a.storage.GetTransfersByCommitmentID(commitmentID)
	if err != nil {
		return nil, err
	}
	dtoTransfers := make([]dto.TransferForCommitment, len(transfers))
	for i := range transfers {
		transfer := &transfers[i]
		dtoTransfers[i] = dto.TransferForCommitment{
			TransferForCommitment: transfer,
			ReceiveTime:           dto.NewTimestamp(transfer.ReceiveTime),
		}
	}
	return dtoTransfers, nil
}

func (a *API) getC2TsForCommitment(commitmentID int32) ([]dto.Create2TransferForCommitment, error) {
	transfers, err := a.storage.GetCreate2TransfersByCommitmentID(commitmentID)
	if err != nil {
		return nil, err
	}
	dtoTransfers := make([]dto.Create2TransferForCommitment, len(transfers))
	for i := range transfers {
		transfer := &transfers[i]
		dtoTransfers[i] = dto.Create2TransferForCommitment{
			Create2TransferForCommitment: transfer,
			ReceiveTime:                  dto.NewTimestamp(transfer.ReceiveTime),
		}
	}
	return dtoTransfers, nil
}
