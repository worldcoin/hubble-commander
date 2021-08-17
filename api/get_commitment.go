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
		BatchTime:    batch.SubmissionTime,
		Transactions: transactions,
	}, nil
}

func (a *API) getTransactionsForCommitment(commitment *models.Commitment) (interface{}, error) {
	switch commitment.Type {
	case txtype.Transfer:
		return a.storage.GetTransfersByCommitmentID(commitment.IndexInBatch)
	case txtype.Create2Transfer:
		return a.storage.GetCreate2TransfersByCommitmentID(commitment.IndexInBatch)
	case txtype.Genesis, txtype.MassMigration:
		return nil, dto.ErrNotImplemented
	}
	return nil, dto.ErrNotImplemented
}
