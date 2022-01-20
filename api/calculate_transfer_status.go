package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func CalculateTransactionStatus(
	storage *st.Storage,
	transfer *models.TransactionBase,
	latestBlockNumber uint32,
) (*txstatus.TransactionStatus, error) {
	if transfer.ErrorMessage != nil {
		return txstatus.Error.Ref(), nil
	}

	if transfer.CommitmentID == nil {
		return txstatus.Pending.Ref(), nil
	}

	batch, err := storage.GetBatch(transfer.CommitmentID.BatchID)
	if err != nil {
		return nil, err
	}

	if batch.FinalisationBlock == nil {
		return txstatus.Pending.Ref(), nil
	}

	if latestBlockNumber < *batch.FinalisationBlock {
		return txstatus.InBatch.Ref(), nil
	}

	return txstatus.Finalised.Ref(), nil
}
