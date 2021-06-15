package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (a *API) GetCommitment(id int32) (*dto.Commitment, error) {
	commitment, err := a.storage.GetCommitment(id)
	if err != nil {
		return nil, err
	}

	var transactions interface{}
	switch commitment.Type {
	case txtype.Transfer:
		transactions, err = a.storage.GetTransfersByCommitmentID(id)
	case txtype.Create2Transfer:
		transactions, err = a.storage.GetCreate2TransfersByCommitmentID(id)
	case txtype.MassMigration:
		return nil, dto.ErrNotImplemented
	}
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
