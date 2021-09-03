package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (a *API) GetCommitment(id models.CommitmentID) (*dto.Commitment, error) {
	commitment, err := a.storage.GetCommitment(&id)
	if err != nil {
		return nil, err
	}

	transactions, err := a.getTransactionsForCommitment(commitment)
	if err != nil {
		return nil, err
	}

	batch, err := a.storage.GetMinedBatch(commitment.ID.BatchID)
	if st.IsNotFoundError(err) {
		return nil, st.NewNotFoundError("commitment")
	}
	if err != nil {
		return nil, err
	}

	return &dto.Commitment{
		Commitment:   *commitment,
		Status:       *calculateFinalisedStatus(a.storage.GetLatestBlockNumber(), *batch.FinalisationBlock),
		BatchTime:    batch.SubmissionTime,
		Transactions: transactions,
	}, nil
}

func (a *API) getTransactionsForCommitment(commitment *models.Commitment) (interface{}, error) {
	switch commitment.Type {
	case txtype.Transfer:
		return a.getTransfersForCommitment(&commitment.ID)
	case txtype.Create2Transfer:
		return a.getCreate2TransfersForCommitment(&commitment.ID)
	case txtype.Genesis, txtype.MassMigration:
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
