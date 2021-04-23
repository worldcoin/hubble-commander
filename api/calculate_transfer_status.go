package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func CalculateTransferStatus(
	storage *st.Storage,
	transfer *models.TransactionBase,
	latestBlockNumber uint32,
) (*txstatus.TransactionStatus, error) {
	if transfer.ErrorMessage != nil {
		return txstatus.Error.Ref(), nil
	}

	if transfer.IncludedInCommitment == nil {
		return txstatus.Pending.Ref(), nil
	}

	batch, err := storage.GetBatchByCommitmentID(*transfer.IncludedInCommitment)
	if err != nil {
		return nil, err
	}

	if latestBlockNumber < batch.FinalisationBlock {
		return txstatus.InBatch.Ref(), nil
	}
	return txstatus.Finalised.Ref(), nil
}
